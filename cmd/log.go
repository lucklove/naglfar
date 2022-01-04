package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/lucklove/naglfar/pkg/client"
	"github.com/lucklove/tidb-log-parser/parser"
	"github.com/lucklove/tidb-log-parser/utils"
	"github.com/spf13/cobra"
)

func newLogCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "log",
		Short: "naglfar log <fragment> [events]",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Help()
			}

			c := client.New()
			defer c.Close()

			now := time.Now()
			logs, err := c.GetLog(cmd.Context(), args[0], now.Add(-time.Hour*24*30), now, args[1:]...)
			if err != nil {
				return err
			}

			stats := make(map[string]utils.StringSet)
			for _, log := range logs {
				for _, f := range log.Fields {
					if stats[f.Name] == nil {
						stats[f.Name] = utils.NewStringSet()
					}
					stats[f.Name].Insert(f.Value)
				}
			}
			for k, s := range stats {
				fmt.Printf("%s=%d\n", k, len(s))
			}

			return nil
		},
	}

	return cmd
}

func fields(fs []parser.LogField) string {
	xs := []string{}
	for _, f := range fs {
		xs = append(xs, fmt.Sprintf("%s=%s", f.Name, f.Value))
	}
	return strings.Join(xs, ",")
}
