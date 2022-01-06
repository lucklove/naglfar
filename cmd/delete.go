package cmd

import (
	"github.com/lucklove/naglfar/pkg/client"
	"github.com/spf13/cobra"
)

func newDeleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "naglfar delete <fragment>",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}

			c := client.New()
			defer c.Close()

			return c.DeleteFragment(cmd.Context(), args[0])
		},
	}

	return cmd
}
