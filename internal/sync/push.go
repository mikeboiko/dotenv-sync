package sync

import (
	"context"
	"errors"
	"fmt"
	"os"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/envfile"
	"dotenv-sync/internal/provider"
	"dotenv-sync/internal/report"
)

func PlanPush(ctx context.Context, cfg config.Config, prov provider.PushProvider) (Plan, provider.EnvPayload, error) {
	schema, err := envfile.ParseFile(cfg.SchemaFile, envfile.KindSchema)
	if err != nil {
		return Plan{}, provider.EnvPayload{}, report.NewAppError("E004", report.ExitOperational, "schema file missing", "push needs .env.example as schema context", "create .env.example or run 'ds init'", err)
	}
	local, err := envfile.ParseFile(cfg.EnvFile, envfile.KindLocal)
	if err != nil {
		if os.IsNotExist(err) {
			return Plan{}, provider.EnvPayload{}, report.NewAppError("E006", report.ExitOperational, "local env file missing", "push needs .env as the upload source", "create .env or choose --env", err)
		}
		return Plan{}, provider.EnvPayload{}, report.NewAppError("E006", report.ExitOperational, "local env file could not be read", "push cannot load the current .env", "fix the file and retry", err)
	}
	return PlanPushDocs(ctx, cfg, schema, local, prov)
}

func PlanPushDocs(ctx context.Context, cfg config.Config, schema, local envfile.Document, prov provider.PushProvider) (Plan, provider.EnvPayload, error) {
	plan := Plan{Mode: "push", Schema: schema, LocalEnv: local, Config: cfg}
	if !cfg.UsesNoteJSON() {
		return plan, provider.EnvPayload{}, report.NewAppError("E009", report.ExitOperational, "ds push requires storage_mode: note_json", "push cannot safely write into field-based Bitwarden items", "set storage_mode: note_json in .envsync.yaml, migrate the repo item, and retry", nil)
	}
	plan.Issues = append(plan.Issues, collectDocumentIssues(schema)...)
	plan.Issues = append(plan.Issues, collectDocumentIssues(local)...)
	if len(plan.Issues) > 0 {
		return plan, provider.EnvPayload{}, issueAsAppError(plan.Issues[0], "push cannot upload the current .env")
	}
	status, err := prov.CheckReadiness(ctx)
	if err != nil {
		return plan, provider.EnvPayload{}, err
	}
	plan.ProviderStatus = status
	if status.Code != "" {
		return plan, provider.EnvPayload{}, report.NewAppError(status.Code, report.ExitOperational, status.Problem, status.Impact, status.Action, nil)
	}
	current, err := prov.LoadEnvPayload(ctx)
	if err != nil {
		return plan, provider.EnvPayload{}, err
	}
	target := provider.EnvPayload{
		ItemName:    cfg.ItemName,
		StorageMode: cfg.StorageMode,
		Exists:      current.Exists,
		Format:      NoteJSONFormat,
		Password:    current.Password,
		Env:         CanonicalDocumentEnv(local),
	}
	plan.Changes = buildPushChanges(schema.AssignmentMap(), target.Env, current.Env)
	plan.WriteRequired = !NoteJSONEqual(target.Env, current.Env) || !current.Exists
	if !current.Exists && len(target.Env) == 0 {
		plan.WriteRequired = false
	}
	return plan, target, nil
}

func buildPushChanges(schemaKeys map[string]envfile.EnvironmentLine, localEnv, providerEnv map[string]string) []ChangeRecord {
	keys := unionKeys(localEnv, providerEnv)
	changes := make([]ChangeRecord, 0, len(keys))
	for _, key := range keys {
		localValue, hasLocal := localEnv[key]
		providerValue, hasProvider := providerEnv[key]
		if _, inSchema := schemaKeys[key]; hasLocal && !inSchema {
			changes = append(changes, ChangeRecord{
				Key:        key,
				ChangeType: "extra",
				Before:     report.RedactValue(providerValue),
				After:      report.RedactValue(localValue),
				File:       "provider",
				Message:    "present in .env but not in .env.example",
			})
			continue
		}
		switch {
		case hasLocal && !hasProvider:
			changes = append(changes, ChangeRecord{Key: key, ChangeType: "add", File: "provider", After: report.RedactValue(localValue), Message: "will be added to Bitwarden"})
		case !hasLocal && hasProvider:
			changes = append(changes, ChangeRecord{Key: key, ChangeType: "update", File: "provider", Before: report.RedactValue(providerValue), After: report.RedactValue(""), Message: "will be removed from Bitwarden"})
		case localValue != providerValue:
			changes = append(changes, ChangeRecord{Key: key, ChangeType: "update", File: "provider", Before: report.RedactValue(providerValue), After: report.RedactValue(localValue), Message: "will be updated in Bitwarden"})
		default:
			changes = append(changes, ChangeRecord{Key: key, ChangeType: "unchanged", File: "provider", After: report.RedactValue(localValue), Message: "already current"})
		}
	}
	return changes
}

func unionKeys(left, right map[string]string) []string {
	set := map[string]struct{}{}
	for key := range left {
		set[key] = struct{}{}
	}
	for key := range right {
		set[key] = struct{}{}
	}
	return report.SortedKeys(set)
}

func issueAsAppError(issue ValidationIssue, impact string) error {
	problem := issue.Message
	if issue.Key != "" {
		problem = fmt.Sprintf("%s: %s", issue.Message, issue.Key)
	}
	return report.NewAppError(issue.Code, report.ExitOperational, problem, impact, issue.Action, nil)
}

func PreserveAppError(err error, fallback func(error) error) error {
	if err == nil {
		return nil
	}
	var appErr *report.AppError
	if errors.As(err, &appErr) {
		return err
	}
	return fallback(err)
}
