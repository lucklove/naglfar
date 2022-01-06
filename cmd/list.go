package cmd

import (
	"fmt"

	"github.com/lucklove/naglfar/pkg/client"
	"github.com/spf13/cobra"
)

func newListCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "list imported fragment list",
		RunE: func(cmd *cobra.Command, args []string) error {
			c := client.New()
			defer c.Close()

			frags, err := c.ListFragments(cmd.Context())
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
