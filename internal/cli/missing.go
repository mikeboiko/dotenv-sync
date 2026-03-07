package cli

import (
	"fmt"

	"dotenv-sync/internal/report"
	syncpkg "dotenv-sync/internal/sync"
	"github.com/spf13/cobra"
)

func newMissingCommand(s streams, opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "missing",
		Short: "List unresolved provider-backed schema keys",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(opts)
			if err != nil {
				return err
			}
			plan, err := syncpkg.PlanMissing(cmd.Context(), cfg, providerFor(cfg))
			for _, change := range plan.Changes {
				if change.ChangeType == "missing" {
					fmt.Fprintln(s.stdout, report.MissingLine(change.Key))
				}
			}
			if err == nil {
				fmt.Fprintln(s.stdout, report.SummaryLine(report.StatusChecked, "provider", syncpkg.Summarize(plan.Changes), "all provider-managed keys resolved"))
			}
			return err
		},
	}
}
