package server

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lucklove/naglfar/pkg/client"
	"github.com/lucklove/tidb-log-parser/event"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

type Server struct {
	router http.Handler
	client *client.Client
	store  *GlobalStore
}

func New() *Server {
	s := &Server{
		client: client.New(),
		store:  NewGlobalStore(),
	}
	s.router = router(s)
	return s
}

func (s *Server) Run(address string) error {
	// do some magic
	go func() {
		for {
			s.buildChangePoint(context.TODO())
			s.buildThreshold(context.TODO())
			s.buildSimilarity(context.TODO())
			time.Sleep(60 * time.Second)
		}
	}()

	return http.ListenAndServe(address, s.router)
}

func (s *Server) similarity(ctx context.Context, r *http.Request) (map[string]float64, error) {
	fragment := mux.Vars(r)["fid"]

	return s.store.GetSimilarity(fragment), nil
}

func (s *Server) threshhold(ctx context.Context, r *http.Request) (*ThresholdRange, error) {
	fragment := mux.Vars(r)["fid"]
	event := mux.Vars(r)["eid"]

	return s.store.GetThreshold(fragment, event), nil
}

func (s *Server) changePoints(ctx context.Context, r *http.Request) ([]TimeRange, error) {
	fragment := mux.Vars(r)["fid"]
	event := mux.Vars(r)["eid"]

	return s.store.GetChangePoint(fragment, event), nil
}

func (s *Server) buildSimilarity(ctx context.Context) {
	log.Info("build similarity")
	frags, err := s.client.ListFragments(ctx)
	if err != nil {
		log.Error("list fragment failed", zap.Error(err))
	}
	if len(frags) < 2 {
		log.Info("no enought fragments found")
		return
	}
	em, err := event.NewEventManager(event.ComponentTiDB)
	if err != nil {
		return
	}
	for _, frag1 := range frags {
		if len(s.store.GetSimilarity(frag1)) == len(frags)-1 {
			log.Info("similarity existed, ignore", zap.String("fragment", frag1))
			continue
		}
		logs1, err := s.client.GetLog(ctx, frag1, time.Now().Add(-time.Hour*24*30), time.Now(), nil)
		if err != nil {
			continue
		}
		for _, frag2 := range frags {
			if frag2 == frag1 {
				continue
			}
			logs2, err := s.client.GetLog(ctx, frag2, time.Now().Add(-time.Hour*24*30), time.Now(), nil)
			if err != nil {
				continue
			}
			inbuf := bytes.NewBuffer(nil)
			outbuf := bytes.NewBuffer(nil)
			for _, log := range logs1 {
				fmt.Fprintf(inbuf, "%d,%d\n", log.Header.DateTime.Unix(), em.GetLogEventID(&log))
			}
			fmt.Fprintf(inbuf, "\n")
			for _, log := range logs2 {
				fmt.Fprintf(inbuf, "%d,%d\n", log.Header.DateTime.Unix(), em.GetLogEventID(&log))
			}
			cmd := exec.CommandContext(ctx, "/usr/bin/python3", "/root/logdeep/demo/LogSimilarity.py")
			cmd.Env = append(cmd.Env, "PYTHONPATH=/root/logdeep")
			cmd.Stdin = inbuf
			cmd.Stdout = outbuf
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Error("error run command", zap.Error(err))
				continue
			}

			scanner := bufio.NewScanner(outbuf)
			for scanner.Scan() {
				line := scanner.Text()
				if !strings.Contains(line, "word2vec") {
					continue
				}
				xs := strings.Split(line, ":")
				if len(xs) != 2 {
					continue
				}
				var sim float64
				fmt.Sscanf(xs[1], "%f", &sim)
				s.store.SetSimilarity(frag1, frag2, sim)
			}
			if err := scanner.Err(); err != nil {
				log.Error("scan error", zap.Error(err))
			} else {
				log.Info("build similarity success", zap.String("fragment", frag1))
			}
		}
	}
}

func (s *Server) buildThreshold(ctx context.Context) {
	log.Info("build treshhold")
	frags, err := s.client.ListFragments(ctx)
	if err != nil {
		log.Error("list fragment failed", zap.Error(err))
	}
	em, err := event.NewEventManager(event.ComponentTiDB)
	if err != nil {
		return
	}
	for _, frag := range frags {
		if s.store.HasThreshold(frag) {
			log.Info("treshhold existed, ignore", zap.String("fragment", frag))
			continue
		}
		log.Info("build treshhold for fragment", zap.String("fragment", frag))
		logs, err := s.client.GetLog(ctx, frag, time.Now().Add(-time.Hour*24*30), time.Now(), nil)
		if err != nil {
			continue
		}
		inbuf := bytes.NewBuffer(nil)
		outbuf := bytes.NewBuffer(nil)
		for _, log := range logs {
			fmt.Fprintf(inbuf, "%d,%d\n", log.Header.DateTime.Unix(), em.GetLogEventID(&log))
		}
		cmd := exec.CommandContext(ctx, "/usr/bin/python3", "/root/logdeep/demo/thresholdShow.py")
		cmd.Env = append(cmd.Env, "PYTHONPATH=/root/logdeep")
		cmd.Stdin = inbuf
		cmd.Stdout = outbuf
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Error("error run command", zap.Error(err))
			continue
		}
		scanner := bufio.NewScanner(outbuf)
		for scanner.Scan() {
			var eid, top, bottom int64
			if _, err := fmt.Sscanf(scanner.Text(), "%d,%d,%d", &eid, &bottom, &top); err != nil {
				continue
			}
			s.store.SetThreshold(frag, fmt.Sprintf("%d", eid), &ThresholdRange{Top: top, Bottom: bottom})
		}
		if err := scanner.Err(); err != nil {
			log.Error("scan error", zap.Error(err))
		} else {
			log.Info("build treshhold success", zap.String("fragment", frag))
		}
	}
}

func (s *Server) buildChangePoint(ctx context.Context) {
	log.Info("build change point")
	frags, err := s.client.ListFragments(ctx)
	if err != nil {
		log.Error("list fragment failed", zap.Error(err))
	}
	em, err := event.NewEventManager(event.ComponentTiDB)
	if err != nil {
		return
	}
	for _, frag := range frags {
		if s.store.HasChangePoint(frag) {
			log.Info("change point existed, ignore", zap.String("fragment", frag))
			continue
		}
		log.Info("build change point for fragment", zap.String("fragment", frag))
		logs, err := s.client.GetLog(ctx, frag, time.Now().Add(-time.Hour*24*30), time.Now(), nil)
		if err != nil {
			continue
		}
		inbuf := bytes.NewBuffer(nil)
		outbuf := bytes.NewBuffer(nil)
		for _, log := range logs {
			fmt.Fprintf(inbuf, "%d,%d\n", log.Header.DateTime.Unix(), em.GetLogEventID(&log))
		}
		cmd := exec.CommandContext(ctx, "/usr/bin/python3", "/root/logdeep/demo/SuddenChangeDetection.py")
		cmd.Env = append(cmd.Env, "PYTHONPATH=/root/logdeep")
		cmd.Stdin = inbuf
		cmd.Stdout = outbuf
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Error("error run command", zap.Error(err))
			continue
		}
		scanner := bufio.NewScanner(outbuf)
		for scanner.Scan() {
			var eid, start, stop int64
			if _, err := fmt.Sscanf(scanner.Text(), "%d,%d,%d", &eid, &start, &stop); err != nil {
				continue
			}
			s.store.SetChangePoint(frag, fmt.Sprintf("%d", eid), TimeRange{Start: start, Stop: stop})
		}
		if err := scanner.Err(); err != nil {
			log.Error("scan error", zap.Error(err))
		} else {
			log.Info("build change point success", zap.String("fragment", frag))
		}
	}
}
