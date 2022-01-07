package cmd

import (
	"github.com/lucklove/naglfar/server"
	"github.com/spf13/cobra"
)

func newServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server",
		Short: "bootstrap server",
		RunE: func(cmd *cobra.Command, args []string) error {
			s := server.New()

			return s.Run(":2048")
		},
	}

	return cmd
}
