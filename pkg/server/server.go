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
	Config   config.Config
	RedisCli redis.Client
}

func NewServer(cfg config.Config, client redis.Client) Server {
	svr := Server{
		Config:   cfg,
		RedisCli: client,
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

func (s Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	resourcePath := "/" + strings.Split(r.URL.Path, "/")[1]
	switch resourcePath {
	case "/ping":
		s.PingRequestHandler(w, r)
	case "/get":
		s.GetRequestHandler(w, r)
	case "/webhook":
		s.SaveRequestHandler(w, r)
	case "/url":
		switch r.Method {
		case "POST":
			s.CreateTokenHandler(w, r)
		case "GET":
			s.GetTokenHandler(w, r)
		}
	case "/pop":
		{
			token := strings.Split(r.URL.Path, "/")[2]
			_, err := uuid.Parse(token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			s.PopRequestHandler(w, r)
		}
	case "/resp":
		{
			token := strings.Split(r.URL.Path, "/")[2]
			_, err := uuid.Parse(token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			s.SaveResponseHandler(w, r)
		}
	default:
		{
			token := strings.Split(r.URL.Path, "/")[1]
			_, err := uuid.Parse(token)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
				return
			}
			s.PushRequestHandler(w, r)
		}
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

	err = s.RedisCli.Set("token:"+fmt.Sprint(id), *b.Url, 1*time.Hour).Err()
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

	url, err := s.RedisCli.Get("token:" + token).Result()
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

func (s *Server) PushRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html")

	path := strings.Split(r.URL.Path, "/")
	token := path[1]
	_, err := uuid.Parse(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req := parseRequest(r)

	b, err := json.Marshal(req)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = s.RedisCli.LPush("request:"+token, b).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.RedisCli.Expire("request:"+token, 1*time.Hour)

	timer := time.NewTimer(time.Duration(5000) * time.Second)

	fmt.Println()
	fmt.Print("waiting for response...")
	respCh := make(chan string)

	go func() {
		for {
			time.Sleep(1 * time.Second)
			res, err := s.RedisCli.LPop("response:" + token).Result()
			if err != nil {
				fmt.Print(".")
			}
			respCh <- res
		}
	}()
	select {
	case <-timer.C:
		http.Error(w, err.Error(), http.StatusGatewayTimeout)
	case res := <-respCh:
		fmt.Fprint(w, res)
	}
}

func (s *Server) PopRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	path := strings.Split(r.URL.Path, "/")
	token := path[2]
	_, err := uuid.Parse(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	req, err := s.RedisCli.LPop("request:" + token).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	var request Request
	err = json.Unmarshal([]byte(req), &request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	json.NewEncoder(w).Encode(request)
}

func (s *Server) SaveResponseHandler(w http.ResponseWriter, r *http.Request) {

	token := strings.Split(r.URL.Path, "/")[2]
	_, err := uuid.Parse(token)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	b, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	_, err = s.RedisCli.LPush("response:"+token, b).Result()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	s.RedisCli.Expire("request:"+token, 1*time.Hour)
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

	val, err := s.RedisCli.Get(token).Result()
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

	err = s.RedisCli.Set(fmt.Sprint(id), string(b), 1*time.Hour).Err()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	fmt.Fprint(w, id)
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

func parseRequest(r *http.Request) Request {
	var req Request = Request{
		Url:    r.URL.Path,
		Host:   r.Host,
		Method: r.Method,
	}

	row := make(map[string]string)
	for h, v := range r.Header {
		row[h] = strings.Join(v, ",")
	}
	jsonData, err := json.Marshal(row)
	if err != nil {
		fmt.Printf("Error: %s", err.Error())
	}
	req.Headers = jsonData

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
