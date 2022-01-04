package server

import (
	"net/http"

	"github.com/lucklove/naglfar/pkg/client"
)

type Server struct {
	router http.Handler
	client *client.Client
}

func New() *Server {
	s := &Server{
		client: client.New(),
	}
	s.router = router(s)
	return s
}

func (s *Server) Run(address string) error {
	// do some magic
	go func() {
	}()

	return http.ListenAndServe(address, s.router)
}
