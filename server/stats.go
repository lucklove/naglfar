package server

import (
	"context"
	"math"
	"net/http"
	"os"
	"path"
	"sort"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/lucklove/tidb-log-parser/store"
	"github.com/pingcap/diag/pkg/utils"
)

func (s *Server) fieldStats(ctx context.Context, r *http.Request) (map[string]int64, error) {
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

	return s.client.GetFieldStats(ctx, fragment, start, stop, nil, event, field)
}

func (s *Server) stats(ctx context.Context, r *http.Request) ([]map[string]interface{}, error) {
	start, err := utils.ParseTime(mux.Vars(r)["start"])
	if err != nil {
		return nil, err
	}
	stop, err := utils.ParseTime(mux.Vars(r)["stop"])
	if err != nil {
		return nil, err
	}
	fragment := mux.Vars(r)["id"]

	countMap, messageMap, err := s.client.GetStats(ctx, fragment, start, stop)
	if err != nil {
		return nil, err
	}
	var count int64
	for _, cnt := range countMap {
		count += cnt
	}

	home, err := os.UserHomeDir()
	if err != nil {
		return nil, err
	}
	store, err := store.NewSQLiteStorage(path.Join(home, ".tiup/storage/naglfar/log.db"), "tidb")
	if err != nil {
		return nil, err
	}

	wm := make(map[string]float64)
	lfc, err := store.LogFragmentCount()
	if err != nil {
		return nil, err
	}
	eids := []string{}
	for eid, cnt := range countMap {
		eids = append(eids, eid)
		id, err := strconv.Atoi(eid)
		if err != nil {
			return nil, err
		}
		ec, err := store.EventCount(uint(id))
		if err != nil {
			return nil, err
		}
		wm[eid] = float64(cnt) / float64(count) * math.Log(float64(lfc)/float64(ec+1))
	}

	sort.Slice(eids, func(i, j int) bool {
		return wm[eids[i]] > wm[eids[j]]
	})

	xs := make([]map[string]interface{}, 0)
	for _, eid := range eids {
		xs = append(xs, map[string]interface{}{
			"event_id": eid,
			"weight":   wm[eid],
			"count":    countMap[eid],
			"message":  messageMap[eid],
		})
	}
	return xs, nil
}
