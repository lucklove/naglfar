package cmd

import (
	"github.com/spf13/cobra"
)

func Command() *cobra.Command {
	cmd := &cobra.Command{
		Use: `naglfar`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
		SilenceUsage: true,
	}

	cmd.AddCommand(
		newImportCommand(),
		newListCommand(),
		newStatsCommand(),
		newTrendCommand(),
		newSearchCommand(),
		newLogCommand(),
		newFieldStatsCommand(),
		newServerCommand(),
	)

	return cmd
}
