package sync

import (
	"context"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/provider"
	"dotenv-sync/internal/report"
)

func PlanValidate(ctx context.Context, cfg config.Config, prov provider.Provider) (Plan, error) {
	plan, _, err := PlanForward(ctx, cfg, prov)
	if err != nil {
		return plan, err
	}
	for _, change := range plan.Changes {
		switch change.ChangeType {
		case "add", "update", "missing", "extra":
			return plan, report.SilentExit(report.ExitValidation)
		}
	}
	return plan, nil
}
