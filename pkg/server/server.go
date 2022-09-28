package server

import (
	"encoding/json"
	"fmt"
	"io"

	"log"
	"net/http"
	"strings"
	"time"

	"github.com/Jeffail/gabs"
	"github.com/go-redis/redis"
	"github.com/google/uuid"

	"github.com/xegea/webhook_server/pkg/config"
)

type Server struct {
	Config config.Config
	Client redis.Client
}

func NewServer(cfg config.Config, client redis.Client) Server {
	svr := Server{
		Config: cfg,
		Client: client,
	}
	return svr
}

type Request struct {
	Url     string          `json:"url,omitempty"`
	Host    string          `json:"host,omitempty"`
	Method  string          `json:"method,omitempty"`
	Body    any             `json:"body,omitempty"`
	Headers json.RawMessage `json:"headers,omitempty"`
}

type ErrorResponse struct {
	Code int
	Desc string
}

type HttpResponse struct {
	Data    json.RawMessage
	Error   ErrorResponse
	Success bool
}

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "/token":
		switch r.Method {
		case "POST":
			s.CreateTokenHandler(w, r)
		case "GET":
			s.GetTokenHandler(w, r)
		}
	case "/ping":
		s.PingRequestHandler(w, r)
	case "/get":
		s.GetRequestHandler(w, r)
	case "/webhook":
		s.SaveRequestHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (s *Server) CreateTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	id := uuid.New()

	//https://stackoverflow.com/a/55052845/2147883
	d := json.NewDecoder(r.Body)
	d.DisallowUnknownFields() // catch unwanted fields
	b := struct {
		Url *string `json:"url"` // pointer so we can test for field absence
	}{}
	err := d.Decode(&b)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if b.Url == nil {
		http.Error(w, "missing field 'url' from JSON object", http.StatusBadRequest)
		return
	}

	err = s.Client.Set(fmt.Sprint(id), *b.Url, time.Hour*1).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

	res := struct {
		Token string `json:"token"`
	}{}
	res.Token = id.String()

	json.NewEncoder(w).Encode(res)
}

func (s *Server) GetTokenHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	token := r.URL.Query().Get("t")
	_, err := uuid.Parse(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	url, err := s.Client.Get(token).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	res := struct {
		Url string `json:"url"`
	}{}
	res.Url = url

	json.NewEncoder(w).Encode(res)
}

func (s *Server) PingRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	req := parseRequest(r)

	b, err := json.MarshalIndent(req, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, string(b))
}

func (s *Server) GetRequestHandler(w http.ResponseWriter, r *http.Request) {

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
	}

	token := r.URL.Query().Get("t")
	_, err := uuid.Parse(token)
	if err != nil {
		log.Print(err)
	}

	val, err := s.Client.Get(token).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := parseRequest(r)

	err = json.Unmarshal([]byte(val), &req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	b, err := json.MarshalIndent(req, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-type", "application/json")
	fmt.Fprint(w, string(b))
}

func (s *Server) SaveRequestHandler(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()

	req := parseRequest(r)

	b, err := json.MarshalIndent(req, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	err = s.Client.Set(fmt.Sprint(id), string(b), time.Hour*1).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, id)
}

func parseRequest(r *http.Request) Request {
	var req Request = Request{
		Url:    r.URL.Path,
		Host:   r.Host,
		Method: r.Method,
	}

	var row json.RawMessage = []byte(`{`)
	for h, v := range r.Header {
		row = append(row, []byte(`"`+h+`":"`+strings.Join(v, ",")+`",`)...)
	}
	row = append(row[:len(row)-1], []byte(`}`)...)

	req.Headers = row

	b, _ := io.ReadAll(r.Body)

	if len(b) > 0 && strings.Contains(r.Header.Get("Content-Type"), "application/json") {
		jsonParsed, err := gabs.ParseJSON(b)
		if err != nil {
			panic(err)
		}
		b = jsonParsed.EncodeJSON()
	}

	req.Body = string(b)
	return req
}

// func jsonError(err error, code int) HttpResponse {
// 	return HttpResponse{
// 		Error: ErrorResponse{
// 			Code: code,
// 			Desc: err.Error(),
// 		},
// 		Success: false,
// 	}
// }
