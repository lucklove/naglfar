package server

import (
	"context"
	"net/http"

	"github.com/gorilla/mux"
	lu "github.com/lucklove/tidb-log-parser/utils"
	du "github.com/pingcap/diag/pkg/utils"
)

func (s *Server) logs(ctx context.Context, r *http.Request) (map[string]interface{}, error) {
	start, err := du.ParseTime(mux.Vars(r)["start"])
	if err != nil {
		return nil, err
	}
	stop, err := du.ParseTime(mux.Vars(r)["stop"])
	if err != nil {
		return nil, err
	}
	fragment := mux.Vars(r)["fid"]
	event := mux.Vars(r)["eid"]

	logs, err := s.client.GetLog(ctx, fragment, start, stop, nil, event)
	if err != nil {
		return nil, err
	}

	stats := make(map[string]lu.StringSet)
	for _, log := range logs {
		for _, f := range log.Fields {
			if stats[f.Name] == nil {
				stats[f.Name] = lu.NewStringSet()
			}
			stats[f.Name].Insert(f.Value)
		}
	}
	ss := map[string]int{}
	for f, s := range stats {
		ss["f_"+f] = len(s)
	}

	xs := make([]map[string]interface{}, 0)
	for _, log := range logs {
		l := map[string]interface{}{
			"event_id":  event,
			"timestamp": log.Header.DateTime.Unix(),
			"level":     log.Header.Level,
			"message":   log.Message,
		}
		for _, f := range log.Fields {
			l["f_"+f.Name] = f.Value
		}
		xs = append(xs, l)
	}
	return map[string]interface{}{
		"logs":  xs,
		"stats": ss,
	}, nil
}
