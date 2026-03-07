package cli

import (
	"fmt"

	"dotenv-sync/internal/report"
	"github.com/spf13/cobra"
)

func newDoctorCommand(s streams, opts *rootOptions) *cobra.Command {
	return &cobra.Command{
		Use:   "doctor",
		Short: "Check config and rbw prerequisites",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := loadConfig(opts)
			if err != nil {
				return err
			}
			fmt.Fprintln(s.stdout, report.SummaryLine(report.StatusChecked, "config", report.Summary{}, cfg.ConfigFile))
			status, err := providerFor(cfg).CheckReadiness(cmd.Context())
			if err != nil {
				return report.NewAppError("E003", report.ExitOperational, "rbw readiness failed", "provider-backed commands cannot continue", "check rbw and retry", err)
			}
			if status.Code != "" {
				return report.NewAppError(status.Code, report.ExitOperational, status.Problem, status.Impact, status.Action, nil)
			}
			fmt.Fprintln(s.stdout, report.SummaryLine(report.StatusChecked, "provider", report.Summary{}, "rbw ready"))
			return nil
		},
	}
}
