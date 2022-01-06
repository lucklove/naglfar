package client

import (
	"context"
	"fmt"
	"strings"
	"time"
)

type Trend struct {
	EventID string  `json:"event_id"`
	Name    string  `json:"name"`
	Points  []Point `json:"points"`
}

type Point struct {
	Timestamp int64 `json:"timestamp"`
	Value     int64 `json:"value"`
}

func (c *Client) GetFieldTrend(ctx context.Context, frag string, start, stop time.Time, event string, field string) ([]Trend, error) {
	queryAPI := c.client.QueryAPI(c.orgID)

	if !strings.HasPrefix(field, "f_") {
		field = "f_" + field
	}

	tr, err := queryAPI.Query(ctx, fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s) %s
			|> filter(fn: (r) => r["_measurement"] =~ /[0-9]{5}/) 
			|> group(columns: ["%s"])
			|> window(every: 5m)
			|> sum()
			|> duplicate(column: "_stop", as: "_time")
			|> window(every: inf)
	`, frag, start.Format(time.RFC3339), stop.Format(time.RFC3339), buildEventFilter([]string{event}), field))
	if err != nil {
		return nil, err
	}
	defer tr.Close()

	tm := make(map[string]Trend)
	for tr.Next() {
		rec := tr.Record()
		name := ""
		if n, ok := rec.Values()[field].(string); ok {
			name = n
		}
		t := Trend{
			EventID: rec.Measurement(),
			Name:    name,
		}
		if mt, ok := tm[name]; ok {
			t = mt
		}
		t.Points = append(t.Points, Point{
			Timestamp: rec.Time().Unix(),
			Value:     rec.Value().(int64),
		})
		tm[name] = t
	}

	ts := []Trend{}
	for _, t := range tm {
		ts = append(ts, t)
	}
	return ts, tr.Err()
}

func (c *Client) GetTrend(ctx context.Context, frag string, start, stop time.Time, events ...string) ([]Trend, error) {
	queryAPI := c.client.QueryAPI(c.orgID)

	tr, err := queryAPI.Query(ctx, fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s) %s
			|> filter(fn: (r) => r["_measurement"] =~ /[0-9]{5}/) 
			|> group(columns: ["_measurement", "name"])
			|> window(every: 5m)
			|> sum()
			|> duplicate(column: "_stop", as: "_time")
			|> window(every: inf)
	`, frag, start.Format(time.RFC3339), stop.Format(time.RFC3339), buildEventFilter(events)))
	if err != nil {
		return nil, err
	}
	defer tr.Close()

	tm := make(map[string]Trend)
	for tr.Next() {
		rec := tr.Record()
		t := Trend{
			EventID: rec.Measurement(),
			Name:    rec.Values()["name"].(string),
		}
		if mt, ok := tm[rec.Measurement()]; ok {
			t = mt
		}
		t.Points = append(t.Points, Point{
			Timestamp: rec.Time().Unix(),
			Value:     rec.Value().(int64),
		})
		tm[rec.Measurement()] = t
	}

	ts := []Trend{}
	for _, t := range tm {
		ts = append(ts, t)
	}
	return ts, tr.Err()
}
