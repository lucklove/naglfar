package server

import (
	"context"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
	"github.com/lucklove/naglfar/pkg/client"
	"github.com/pingcap/diag/pkg/utils"
)

func (s *Server) trend(ctx context.Context, r *http.Request) ([]client.Trend, error) {
	events := strings.Split(r.FormValue("events"), ",")
	start, err := utils.ParseTime(mux.Vars(r)["start"])
	if err != nil {
		return nil, err
	}
	stop, err := utils.ParseTime(mux.Vars(r)["stop"])
	if err != nil {
		return nil, err
	}
	fragment := mux.Vars(r)["id"]

	return s.client.GetTrend(ctx, fragment, start, stop, events...)
}

func (s *Server) fieldTrend(ctx context.Context, r *http.Request) ([]client.Trend, error) {
	start, err := utils.ParseTime(mux.Vars(r)["start"])
	if err != nil {
		return nil, err
	}
	stop, err := utils.ParseTime(mux.Vars(r)["stop"])
	if err != nil {
		return nil, err
	}
	fragment := mux.Vars(r)["fid"]
	event := mux.Vars(r)["eid"]
	field := mux.Vars(r)["field"]

	return s.client.GetFieldTrend(ctx, fragment, start, stop, event, field)
}
