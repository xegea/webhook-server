package server

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
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

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
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

type Request struct {
	Url     string          `json:"url,omitempty"`
	Host    string          `json:"host,omitempty"`
	Method  string          `json:"method,omitempty"`
	Body    any             `json:"body,omitempty"`
	Headers json.RawMessage `json:"headers,omitempty"`
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
	w.Header().Set("Content-type", "application/json")

	if err := r.ParseForm(); err != nil {
		log.Printf("Error parsing form: %s", err)
	}

	token := r.Form.Get("t")
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

	fmt.Fprint(w, string(b))
}

func (s *Server) SaveRequestHandler(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()
	fmt.Println(id.String())

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

	b, _ := ioutil.ReadAll(r.Body)

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
