package cli

import (
	"errors"
	"fmt"

	"dotenv-sync/internal/report"
	syncpkg "dotenv-sync/internal/sync"
	"github.com/spf13/cobra"
)

func newPushCommand(s streams, opts *rootOptions) *cobra.Command {
	var dryRun bool
	cmd := &cobra.Command{
		Use:   "push",
		Short: "Upload the current .env into Bitwarden",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(opts)
			if err != nil {
				return err
			}
			prov := pushProviderFor(cfg)
			plan, target, err := syncpkg.PlanPush(cmd.Context(), cfg, prov)
			var appErr *report.AppError
			if err != nil && errors.As(err, &appErr) {
				return err
			}
			for _, change := range plan.Changes {
				if dryRun {
					fmt.Fprintln(s.stdout, report.ChangeLine(change.ChangeType, change.Key, change.After))
				}
			}
			if err != nil {
				return err
			}
			summary := syncpkg.Summarize(plan.Changes)
			targetName := "bitwarden:" + target.ItemName
			if dryRun {
				if plan.WriteRequired {
					fmt.Fprintln(s.stdout, report.SummaryLine(report.StatusChecked, targetName, summary, "dry-run"))
				} else {
					fmt.Fprintln(s.stdout, report.SummaryLine(report.StatusUnchanged, targetName, summary, "already up to date"))
				}
				return nil
			}
			if !plan.WriteRequired {
				fmt.Fprintln(s.stdout, report.SummaryLine(report.StatusUnchanged, targetName, summary, "already up to date"))
				return nil
			}
			if _, err := prov.StoreEnvPayload(cmd.Context(), target); err != nil {
				return err
			}
			fmt.Fprintln(s.stdout, report.SummaryLine(report.StatusWritten, targetName, summary, ""))
			return nil
		},
	}
	cmd.Flags().BoolVar(&dryRun, "dry-run", false, "preview without writing to Bitwarden")
	return cmd
}
