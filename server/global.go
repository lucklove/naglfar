package server

import (
	"sync"
)

type TimeRange struct {
	Start int64 `json:"start"`
	Stop  int64 `json:"stop"`
}

type ThresholdRange struct {
	Top    int64 `json:"top"`
	Bottom int64 `json:"bottom"`
}

type GlobalStore struct {
	mu           sync.Mutex
	ChangePoints map[string]map[string][]TimeRange     `json:"change_points"`
	Threshold    map[string]map[string]*ThresholdRange `json:"threshold"`
}

func NewGlobalStore() *GlobalStore {
	return &GlobalStore{
		ChangePoints: make(map[string]map[string][]TimeRange),
	}
}

func (gs *GlobalStore) HasChangePoint(fragment string) bool {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	return len(gs.ChangePoints[fragment]) > 0
}

func (gs *GlobalStore) SetChangePoint(fragment, event string, r TimeRange) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.ChangePoints[fragment] == nil {
		gs.ChangePoints[fragment] = make(map[string][]TimeRange)
	}
	for _, point := range gs.ChangePoints[fragment][event] {
		if point.Start == r.Start {
			return
		}
	}
	gs.ChangePoints[fragment][event] = append(gs.ChangePoints[fragment][event], r)
}

func (gs *GlobalStore) GetChangePoint(fragment, event string) []TimeRange {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.ChangePoints[fragment] == nil {
		return nil
	}
	return append(make([]TimeRange, 0), gs.ChangePoints[fragment][event]...)
}

func (gs *GlobalStore) HasThreshold(fragment string) bool {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	return len(gs.Threshold[fragment]) > 0
}

func (gs *GlobalStore) GetThreshold(fragment, event string) *ThresholdRange {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.Threshold[fragment] == nil {
		return nil
	}
	return gs.Threshold[fragment][event]
}

func (gs *GlobalStore) SetThreshold(fragment, event string, r *ThresholdRange) {
	gs.mu.Lock()
	defer gs.mu.Unlock()

	if gs.Threshold[fragment] == nil {
		gs.Threshold[fragment] = make(map[string]*ThresholdRange)
	}
	gs.Threshold[fragment][event] = r
}
