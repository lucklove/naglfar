package cmd

import (
	"fmt"
	"time"

	"github.com/lucklove/naglfar/pkg/client"
	"github.com/spf13/cobra"
)

func newSearchCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "search",
		Short: "naglfar search <event-id>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}

			c := client.New()
			defer c.Close()

			n := time.Now()
			frags, err := c.Search(cmd.Context(), args[0], n.Add(-time.Hour*24*30), n)
			if err != nil {
				return err
			}

			for _, frag := range frags {
				fmt.Println(frag)
			}
			return nil
		},
	}

	return cmd
}
