package sync

import (
	"context"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/provider"
	"dotenv-sync/internal/report"
)

func PlanMissing(ctx context.Context, cfg config.Config, prov provider.Provider) (Plan, error) {
	plan, _, err := PlanForward(ctx, cfg, prov)
	if err != nil {
		return plan, err
	}
	missing := false
	for _, change := range plan.Changes {
		if change.ChangeType == "missing" {
			missing = true
			break
		}
	}
	if missing {
		return plan, report.SilentExit(report.ExitValidation)
	}
	return plan, nil
}
