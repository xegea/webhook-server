package server

import (
	"net/http"

	"github.com/xegea/webhook_server/pkg/config"
)

type Server struct {
	Config config.Config
}

func NewServer(cfg config.Config) Server {
	svr := Server{
		Config: cfg,
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
