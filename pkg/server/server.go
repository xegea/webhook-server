package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

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
	Url    string        `json:"url,omitempty"`
	Host   string        `json:"host,omitempty"`
	Method string        `json:"method,omitempty"`
	Body   io.ReadCloser `json:"body,omitempty"`
}

func (s *Server) PingRequestHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")

	var req Request = Request{
		Url:    r.URL.Path,
		Host:   r.Host,
		Method: r.Method,
		Body:   r.Body,
	}

	b, err := json.MarshalIndent(req, "", "\t")
	if err != nil {
		log.Print(err)
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
		// Internal server error
		log.Print(err)
	}

	var req Request = Request{
		Url:    r.URL.Path,
		Host:   r.Host,
		Method: r.Method,
		Body:   r.Body,
	}

	err = json.Unmarshal([]byte(val), &req)
	if err != nil {
		log.Print(err)
	}

	b, err := json.MarshalIndent(req, "", "\t")
	if err != nil {
		log.Print(err)
	}

	fmt.Fprint(w, string(b))
}

func (s *Server) SaveRequestHandler(w http.ResponseWriter, r *http.Request) {
	id := uuid.New()
	fmt.Println(id.String())

	var req Request = Request{
		Url:    r.URL.Path,
		Host:   r.Host,
		Method: r.Method,
		Body:   r.Body,
	}

	b, err := json.MarshalIndent(req, "", "\t")
	if err != nil {
		panic(err)
	}

	err = s.Client.Set(fmt.Sprint(id), string(b), 3600).Err()
	if err != nil {
		panic(err)
	}
}
