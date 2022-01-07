package cmd

import (
	"github.com/lucklove/naglfar/server"
	"github.com/spf13/cobra"
)

func newServerCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "server <web-root>",
		Short: "bootstrap server",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return cmd.Help()
			}

			s := server.New(args[0])

			return s.Run(":2048")
		},
	}

	return cmd
}
