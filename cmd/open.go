package cmd

import (
	"fmt"

	"github.com/lucklove/naglfar/pkg/browser"
	"github.com/spf13/cobra"
)

func newOpenCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "open <fragment> [event]",
		Short: "open a browser to view logs",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return cmd.Help()
			}

			if len(args) == 1 {
				browser.Open(fmt.Sprintf("http://hackhost:2048/index.html?fid=%s", args[0]))
			} else {
				browser.Open(fmt.Sprintf("http://hackhost:2048/index.html?fid=%s&eid=%s", args[0], args[1]))
			}

			return nil
		},
	}

	return cmd
}
