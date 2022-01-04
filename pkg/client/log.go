package client

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/lucklove/tidb-log-parser/parser"
)

func (c *Client) GetLog(ctx context.Context, frag string, start, stop time.Time, events ...string) ([]parser.LogEntry, error) {
	queryAPI := c.client.QueryAPI(c.orgID)

	tr, err := queryAPI.Query(ctx, fmt.Sprintf(`
		from(bucket: "%s")
			|> range(start: %s, stop: %s) %s
	`, frag, start.Format(time.RFC3339), stop.Format(time.RFC3339), buildEventFilter(events)))
	if err != nil {
		return nil, err
	}
	defer tr.Close()

	logs := []parser.LogEntry{}
	for tr.Next() {
		rec := tr.Record()

		log := parser.LogEntry{
			Header: parser.LogHeader{
				DateTime: rec.Time(),
				Level:    parser.LogLevel(rec.Values()["level"].(string)),
			},
			Message: rec.Values()["message"].(string),
		}
		for k, v := range rec.Values() {
			if strings.HasPrefix(k, "f_") {
				log.Fields = append(log.Fields, parser.LogField{
					Name:  k[2:],
					Value: v.(string),
				})
			}
		}
		logs = append(logs, log)
	}

	return logs, nil
}
