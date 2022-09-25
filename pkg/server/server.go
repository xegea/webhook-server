package server

import (
	"net/http"

	"github.com/go-redis/redis"
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
	case "/get":
		s.GetRequestHandler(w, r)
	case "/post":
		s.SaveRequestHandler(w, r)
	default:
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
	}
}

func (s *Server) GetRequestHandler(w http.ResponseWriter, r *http.Request) {

}

func (s *Server) SaveRequestHandler(w http.ResponseWriter, r *http.Request) {

}
