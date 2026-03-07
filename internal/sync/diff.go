package sync

import (
	"context"
	"sort"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/provider"
	"dotenv-sync/internal/report"
)

func PlanDiff(ctx context.Context, cfg config.Config, prov provider.Provider) (Plan, error) {
	plan, _, err := PlanForward(ctx, cfg, prov)
	if err != nil {
		return plan, err
	}
	sort.Slice(plan.Changes, func(i, j int) bool { return plan.Changes[i].Key < plan.Changes[j].Key })
	if summary := Summarize(plan.Changes); summary.Added+summary.Updated+summary.Missing+summary.Extra == 0 {
		return plan, nil
	}
	return plan, report.SilentExit(report.ExitValidation)
}
