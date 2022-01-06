package server

import (
	"net/http"

	"github.com/lucklove/naglfar/pkg/client"
)

type Server struct {
	router http.Handler
	client *client.Client
	store  *GlobalStore
}

func New() *Server {
	s := &Server{
		client: client.New(),
		store:  &GlobalStore{},
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
