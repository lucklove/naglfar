package main

import (
	"fmt"
	"os"
	"strings"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	http2 "github.com/influxdata/influxdb-client-go/v2/api/http"
	"github.com/lucklove/tidb-log-parser/event"
	"github.com/lucklove/tidb-log-parser/parser"
)

func assertNil(err error) {
	if err != nil {
		panic(err)
	}
}

func ignore(log *parser.LogEntry) bool {
	ignoreFiles := []string{}
	ignoreMessages := []string{}

	for _, f := range ignoreFiles {
		if f == log.Header.File {
			return true
		}
	}
	for _, m := range ignoreMessages {
		if strings.Contains(log.Message, m) {
			return true
		}
	}
	return false
}

func main() {
	// Create a client
	// You can generate an API Token from the "API Tokens Tab" in the UI
	client := influxdb2.NewClient("http://hackhost:8086", "XSIB6UoPvDWd7XOmwR3KiFL9C_Yf5vbOx0bjFTBu41xNE9svgzBmw4bytKxemGy2ZbkmjhzwUDyDuVGwKqFnjg==")

	// get non-blocking write client
	writeAPI := client.WriteAPI("Manual Piolot", "tidb-cluster-03")
	writeAPI.SetWriteFailedCallback(func(batch string, error http2.Error, retryAttempts uint) bool {
		panic(batch)
	})

	parser := parser.NewStreamParser(os.Stdin)
	em, err := event.NewEventManager(event.ComponentTiDB)
	assertNil(err)

	cnt := 0
	for {
		log, err := parser.Next()
		if log == nil && err == nil {
			break
		}
		if log == nil || err != nil {
			continue
		}
		if ignore(log) {
			continue
		}
		event := em.GetRuleByLog(log)
		if event == nil {
			continue
		}

		cnt = (cnt + 1) % 1000
		p := influxdb2.NewPointWithMeasurement(fmt.Sprintf("%d", event.ID)).
			AddTag("name", event.Name).
			AddTag("message", log.Message).
			AddTag("level", string(log.Header.Level)).
			AddField("count", 1).
			SetTime(log.Header.DateTime.Add(time.Nanosecond * time.Duration(cnt)))
		for _, f := range log.Fields {
			p = p.AddTag("f_"+f.Name, f.Value)
		}
		// write point asynchronously
		writeAPI.WritePoint(p)
	}

	writeAPI.Flush()
	// always close client at the end
	defer client.Close()
}
