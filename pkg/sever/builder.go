package server

import (
	"log"
	"net/http"
	"time"
)

type Server struct {
	mux    *http.ServeMux
	server *http.Server
}

func NewBuilder() *Server {
	return &Server{
		mux: http.NewServeMux(),
	}
}

func (s *Server) WithRoute(path string, handler http.HandlerFunc) *Server {
	s.mux.Handle(path, handler)
	return s
}

func (s *Server) WithTimeout(timeout int) *Server {
	s.server = &http.Server{
		Addr:         ":8080",
		Handler:      s.mux,
		ReadTimeout:  time.Duration(timeout) * time.Second,
		WriteTimeout: time.Duration(timeout) * time.Second,
	}
	return s
}

func (s *Server) Build() *Server {
	return s
}

func (s *Server) Start() error {
	log.Println("Starting server on :8080")
	return s.server.ListenAndServe()
}

func (s *Server) Shutdown() error {
	log.Println("Shutting down server")
	return s.server.Close()
}
