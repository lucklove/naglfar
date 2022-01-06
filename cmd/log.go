package cmd

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/lucklove/naglfar/pkg/client"
	"github.com/lucklove/tidb-log-parser/event"
	"github.com/lucklove/tidb-log-parser/parser"
	du "github.com/pingcap/diag/pkg/utils"
	"github.com/pingcap/tiup/pkg/tui"
	"github.com/spf13/cobra"
)

func newLogCommand() *cobra.Command {
	filters := []string{}
	begin := ""
	end := ""

	cmd := &cobra.Command{
		Use:   "log <fragment> [events]",
		Short: "get log detail of specified fragment",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Help()
			}

			c := client.New()
			defer c.Close()

			start, err := du.ParseTime(begin)
			if err != nil {
				start = time.Now().Add(-time.Hour * 24 * 30)
			}
			stop, err := du.ParseTime(end)
			if err != nil {
				stop = time.Now()
			}

			logs, err := c.GetLog(cmd.Context(), args[0], start, stop, filters, args[1:]...)
			if err != nil {
				return err
			}
			sort.Slice(logs, func(i, j int) bool {
				return logs[i].Header.DateTime.Before(logs[j].Header.DateTime)
			})

			em, err := event.NewEventManager(event.ComponentTiDB)
			if err != nil {
				return err
			}

			table := [][]string{{"ID", "time", "level", "message", "fields"}}
			for _, log := range logs {
				table = append(table, []string{
					fmt.Sprintf("%d", em.GetLogEventID(&log)),
					log.Header.DateTime.Format(time.RFC3339),
					string(log.Header.Level),
					log.Message,
					fields(log.Fields),
				})
			}
			tui.PrintTable(table, true)

			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&filters, "filter", "f", nil, "filter fields values")
	cmd.Flags().StringVarP(&begin, "begin", "b", begin, "specific begin time")
	cmd.Flags().StringVarP(&end, "end", "e", end, "specific end time")

	return cmd
}

func fields(fs []parser.LogField) string {
	xs := []string{}
	sort.Slice(fs, func(i, j int) bool {
		return fs[i].Name < fs[j].Name
	})
	for _, f := range fs {
		xs = append(xs, fmt.Sprintf("%s=%s", f.Name, f.Value))
	}
	return strings.Join(xs, ",")
}
