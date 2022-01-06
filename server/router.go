package server

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/pingcap/fn"
)

func router(s *Server) http.Handler {
	r := mux.NewRouter()

	r.Handle("/api/v1/fragments", fn.Wrap(s.listFragment)).Methods("GET")
	r.Handle("/api/v1/{start}/{stop}/events/{id}/fragments", fn.Wrap(s.searchFragmentByEvent)).Methods("GET")
	r.Handle("/api/v1/{start}/{stop}/fragments/{id}/logs/stats", fn.Wrap(s.stats)).Methods("GET")
	r.Handle("/api/v1/{start}/{stop}/fragments/{fid}/events/{eid}/fields/{field}/logs/stats", fn.Wrap(s.fieldStats)).Methods("GET")
	r.Handle("/api/v1/{start}/{stop}/fragments/{fid}/events/{eid}/logs", fn.Wrap(s.logs)).Methods("GET")
	r.Handle("/api/v1/{start}/{stop}/fragments/{id}/logs/trend", fn.Wrap(s.trend)).Methods("GET")
	r.Handle("/api/v1/{start}/{stop}/fragments/{fid}/events/{eid}/fields/{field}/logs/trend", fn.Wrap(s.fieldTrend)).Methods("GET")

	return r
}