package server

import (
	"fmt"
	"net/http"
)

// Server http server entity
type Server struct {
	handler http.Handler
	port    int
}

// NewServer creates new instance of http server
func NewServer(port int, handler http.Handler) Server {
	return Server{
		handler: handler,
		port:    port,
	}
}

// Start - starts server
func (s Server) Start() {
	go http.ListenAndServe(fmt.Sprintf(":%d", s.port), s.handler)
}
