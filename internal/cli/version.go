package cli

import (
	"fmt"

	"dotenv-sync/pkg/dotenvsync"
	"github.com/spf13/cobra"
)

func newVersionCommand(s streams) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print detailed version information",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			_, _ = fmt.Fprintln(s.stdout, dotenvsync.Detailed())
			return nil
		},
	}
}
