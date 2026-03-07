package cli

import (
	"fmt"

	"dotenv-sync/internal/report"
	syncpkg "dotenv-sync/internal/sync"
	"github.com/spf13/cobra"
)

func newDiffCommand(s streams, opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "diff",
		Short: "Preview redacted changes without writing",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(opts)
			if err != nil {
				return err
			}
			plan, err := syncpkg.PlanDiff(cmd.Context(), cfg, providerFor(cfg))
			for _, change := range plan.Changes {
				fmt.Fprintln(s.stdout, report.ChangeLine(change.ChangeType, change.Key, change.After))
			}
			if err == nil {
				fmt.Fprintln(s.stdout, report.SummaryLine(report.StatusUnchanged, cfg.EnvFile, syncpkg.Summarize(plan.Changes), "already up to date"))
			}
			return err
		},
	}
}
