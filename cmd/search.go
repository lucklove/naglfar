package cmd

import (
	"fmt"
	"time"

	"github.com/lucklove/naglfar/pkg/client"
	du "github.com/pingcap/diag/pkg/utils"
	"github.com/spf13/cobra"
)

func newSearchCommand() *cobra.Command {
	begin := ""
	end := ""

	cmd := &cobra.Command{
		Use:   "search <log-id>",
		Short: "search specified log in imported fragments",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
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
			frags, err := c.Search(cmd.Context(), args[0], start, stop)
			if err != nil {
				return err
			}

			for _, frag := range frags {
				fmt.Println(frag)
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&begin, "begin", "b", begin, "specific begin time")
	cmd.Flags().StringVarP(&end, "end", "e", end, "specific end time")

	return cmd
}
