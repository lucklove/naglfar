package cmd

import (
	"fmt"
	"math"
	"os"
	"path"
	"sort"
	"strconv"
	"time"

	"github.com/lucklove/naglfar/pkg/client"
	"github.com/lucklove/tidb-log-parser/store"
	"github.com/spf13/cobra"
)

func newFieldStatsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "fstats",
		Short: "naglfar stats <fragment> <event>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 3 {
				return cmd.Help()
			}

			c := client.New()
			defer c.Close()

			n := time.Now()
			fmap, err := c.GetFieldStats(cmd.Context(), args[0], n.Add(-time.Hour*24*30), n, args[1], args[2])
			if err != nil {
				return err
			}
			for f, cnt := range fmap {
				fmt.Println(f, cnt)
			}

			return nil
		},
	}

	return cmd
}

func newStatsCommand() *cobra.Command {
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

			for _, eid := range eids {
				fmt.Printf("%f\t%d\t%s\n", wm[eid], countMap[eid], messageMap[eid])
			}

			return nil
		},
	}

	return cmd
}
