package server

import (
	"net/http"

	"github.com/guisithos/save-my-read/internal/interfaces/http/handlers"
)

type Server struct {
	router    *http.ServeMux
	jwtSecret string
}

func NewServer(authHandler *handlers.AuthHandler) *Server {
	s := &Server{
		router: http.NewServeMux(),
	}

	// Register routes
	s.router.HandleFunc("/api/auth/register", authHandler.Register)
	s.router.HandleFunc("/api/auth/login", authHandler.Login)

	return s
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *Server) Router() http.Handler {
	return s.router
}
