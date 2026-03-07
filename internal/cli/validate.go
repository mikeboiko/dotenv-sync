package cli

import (
	"fmt"

	"dotenv-sync/internal/report"
	syncpkg "dotenv-sync/internal/sync"
	"github.com/spf13/cobra"
)

func newValidateCommand(s streams, opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "validate",
		Short: "Verify schema, local env, and provider readiness",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(opts)
			if err != nil {
				return err
			}
			plan, err := syncpkg.PlanValidate(cmd.Context(), cfg, providerFor(cfg))
			if err == nil {
				fmt.Fprintln(s.stdout, report.SummaryLine(report.StatusChecked, cfg.EnvFile, syncpkg.Summarize(plan.Changes), "no blocking issues"))
				return nil
			}
			for _, issue := range plan.Issues {
				fmt.Fprintln(s.stdout, report.ChangeLine(issue.Code, issue.Key, issue.Message))
			}
			for _, change := range plan.Changes {
				if change.ChangeType != "unchanged" {
					fmt.Fprintln(s.stdout, report.ChangeLine(change.ChangeType, change.Key, change.After))
				}
			}
			return err
		},
	}
}
