package cmd

import (
	"fmt"
	"time"

	"github.com/lucklove/naglfar/pkg/client"
	"github.com/spf13/cobra"
)

func newTrendCommand() *cobra.Command {
	field := ""

	cmd := &cobra.Command{
		Use:   "trend",
		Short: "naglfar trend <fragment> [events]",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Help()
			}

			c := client.New()
			defer c.Close()

			n := time.Now()
			if field != "" && len(args) == 2 {
				trends, err := c.GetFieldTrend(cmd.Context(), args[0], n.Add(-time.Hour*24*30), n, args[1], field)
				if err != nil {
					return err
				}
				for _, trend := range trends {
					fmt.Println(trend.EventID, trend.Name, len(trend.Points))
				}
			} else {
				trends, err := c.GetTrend(cmd.Context(), args[0], n.Add(-time.Hour*24*30), n, args[1:]...)
				if err != nil {
					return err
				}
				for _, trend := range trends {
					fmt.Println(trend.EventID, trend.Name, len(trend.Points))
				}
			}

			return nil
		},
	}

	cmd.Flags().StringVarP(&field, "field", "f", "", "Specify the field to group")

	return cmd
}
