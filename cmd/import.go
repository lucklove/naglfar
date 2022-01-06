package cmd

import (
	"fmt"
	"os"
	"time"

	influxdb2 "github.com/influxdata/influxdb-client-go/v2"
	http2 "github.com/influxdata/influxdb-client-go/v2/api/http"
	"github.com/lucklove/naglfar/pkg/client"
	"github.com/lucklove/tidb-log-parser/event"
	"github.com/lucklove/tidb-log-parser/parser"
	"github.com/spf13/cobra"
)

func newImportCommand() *cobra.Command {
	fragment := ""

	cmd := &cobra.Command{
		Use:   "import",
		Short: "naglfar import <log-file> --as <fragment>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}

			c := client.New()
			defer c.Close()

			// get non-blocking write client
			writeAPI, err := c.ImportFragment(cmd.Context(), fragment)
			if err != nil {
				return err
			}
			writeAPI.SetWriteFailedCallback(func(batch string, error http2.Error, retryAttempts uint) bool {
				panic(batch)
			})

			file, err := os.Open(args[0])
			if err != nil {
				return err
			}
			defer file.Close()

			parser := parser.NewStreamParser(file)
			em, err := event.NewEventManager(event.ComponentTiDB)
			if err != nil {
				return err
			}

			cnt := 0
			for {
				log, err := parser.Next()
				if log == nil && err == nil {
					break
				}
				if log == nil || err != nil {
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
			return nil
		},
	}

	cmd.Flags().StringVarP(&fragment, "as", "", "", "Specify the fragment name")
	cmd.MarkFlagRequired("as")

	return cmd
}
