package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pingcap/diag/pkg/utils"
)

func (s *Server) listFragment(ctx context.Context) ([]string, error) {
	return s.client.ListFragments(ctx)
}

func (s *Server) searchFragmentByEvent(ctx context.Context, r *http.Request) ([]string, error) {
	start, err := utils.ParseTime(mux.Vars(r)["start"])
	if err != nil {
		return nil, err
	}
	stop, err := utils.ParseTime(mux.Vars(r)["stop"])
	if err != nil {
		return nil, err
	}
	event := mux.Vars(r)["id"]
	frags, err := s.client.Search(ctx, event, start, stop)
	if err != nil {
		return nil, err
	}
	return frags, nil
}
