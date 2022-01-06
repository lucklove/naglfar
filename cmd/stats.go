package cmd

import (
	"fmt"
	"math"
	"os"
	"path"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/lucklove/naglfar/pkg/client"
	"github.com/lucklove/tidb-log-parser/event"
	"github.com/lucklove/tidb-log-parser/store"
	"github.com/lucklove/tidb-log-parser/utils"
	"github.com/pingcap/tiup/pkg/tui"
	"github.com/spf13/cobra"
)

func newFieldStatsCommand() *cobra.Command {
	filters := []string{}

	cmd := &cobra.Command{
		Use:   "fstats",
		Short: "naglfar stats <fragment> <event> <field>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return cmd.Help()
			}

			c := client.New()
			defer c.Close()

			n := time.Now()
			fmap, err := c.GetFieldStats(cmd.Context(), args[0], n.Add(-time.Hour*24*30), n, filters, args[1], args[2])
			if err != nil {
				return err
			}

			table := [][]string{{"Field", "Value", "Count"}}
			for f, cnt := range fmap {
				table = append(table, []string{args[2], f, fmt.Sprintf("%d", cnt)})
			}
			tui.PrintTable(table, true)

			return nil
		},
	}

	cmd.Flags().StringSliceVarP(&filters, "filter", "f", nil, "filter fields values")

	return cmd
}

func newStatsCommand() *cobra.Command {
	withFields := false
	filters := []string{}

	cmd := &cobra.Command{
		Use:   "stats",
		Short: "naglfar stats <fragment>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}

			c := client.New()
			defer c.Close()

			n := time.Now()
			countMap, messageMap, err := c.GetStats(cmd.Context(), args[0], n.Add(-time.Hour*24*30), n)
			if err != nil {
				return err
			}

			stats := make(map[string]map[string]utils.StringSet)
			if withFields {
				logs, err := c.GetLog(cmd.Context(), args[0], n.Add(-time.Hour*24*30), n, filters, args[1:]...)
				if err != nil {
					return err
				}
				em, err := event.NewEventManager(event.ComponentTiDB)
				if err != nil {
					return err
				}
				for _, log := range logs {
					id := fmt.Sprintf("%d", em.GetLogEventID(&log))
					if stats[id] == nil {
						stats[id] = make(map[string]utils.StringSet)
					}
					for _, f := range log.Fields {
						if stats[id][f.Name] == nil {
							stats[id][f.Name] = utils.NewStringSet()
						}
						stats[id][f.Name].Insert(f.Value)
					}
				}
			}

			var count int64
			for _, cnt := range countMap {
				count += cnt
			}

			home, err := os.UserHomeDir()
			if err != nil {
				return err
			}
			store, err := store.NewSQLiteStorage(path.Join(home, ".tiup/storage/naglfar/log.db"), "tidb")
			if err != nil {
				return err
			}

			wm := make(map[string]float64)
			lfc, err := store.LogFragmentCount()
			if err != nil {
				return err
			}
			eids := []string{}
			for eid, cnt := range countMap {
				eids = append(eids, eid)
				id, err := strconv.Atoi(eid)
				if err != nil {
					return err
				}
				ec, err := store.EventCount(uint(id))
				if err != nil {
					return err
				}
				wm[eid] = float64(cnt) / float64(count) * math.Log(float64(lfc)/float64(ec+1))
			}

			sort.Slice(eids, func(i, j int) bool {
				return wm[eids[i]] > wm[eids[j]]
			})

			header := []string{"ID", "weight", "count", "template"}
			if withFields {
				header = append(header, "fields")
			}
			table := [][]string{header}
			for _, eid := range eids {
				row := []string{
					eid,
					fmt.Sprintf("%f", wm[eid]),
					fmt.Sprintf("%d", countMap[eid]),
					messageMap[eid],
				}
				if withFields {
					row = append(row, renderFields(stats[eid]))
				}
				table = append(table, row)
			}
			tui.PrintTable(table, true)

			return nil
		},
	}

	cmd.Flags().BoolVar(&withFields, "with-fields", false, "print fields stats")
	cmd.Flags().StringSliceVarP(&filters, "filter", "f", nil, "filter fields values")

	return cmd
}

func renderFields(fmap map[string]utils.StringSet) string {
	if len(fmap) == 0 {
		return ""
	}
	xs := []string{}
	for k, v := range fmap {
		xs = append(xs, fmt.Sprintf("%s(%d)", k, len(v)))
	}
	return strings.Join(xs, ",")
}
