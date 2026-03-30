package bitwarden

import (
	"context"
	"errors"
	"strings"
	"sync"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/provider"
	"dotenv-sync/internal/report"
	syncpkg "dotenv-sync/internal/sync"
)

type Adapter struct {
	client        *RBWClient
	cfg           config.Config
	mu            sync.Mutex
	cache         map[string]provider.Resolution
	notePayload   provider.EnvPayload
	notePayloadOK bool
}

func NewAdapter(cfg config.Config) *Adapter {
	return &Adapter{client: NewRBWClient(), cfg: cfg, cache: map[string]provider.Resolution{}}
}

func (a *Adapter) Name() string { return "bitwarden" }

func (a *Adapter) CheckReadiness(ctx context.Context) (provider.Status, error) {
	return checkReadinessWithClient(ctx, a.client)
}

func (a *Adapter) Resolve(ctx context.Context, key, ref string) (provider.Resolution, error) {
	if a.cfg.UsesNoteJSON() {
		return a.resolveNoteJSON(ctx, key)
	}
	return a.resolveField(ctx, key, ref)
}

func (a *Adapter) resolveField(ctx context.Context, key, ref string) (provider.Resolution, error) {
	fieldName := ref
	if fieldName == "" {
		fieldName = key
	}
	itemName := a.itemName()
	cacheKey := itemName + "|" + fieldName
	a.mu.Lock()
	if cached, ok := a.cache[cacheKey]; ok {
		a.mu.Unlock()
		cached.Key = key
		cached.Ref = cacheKey
		return cached, nil
	}
	a.mu.Unlock()
	out, err := a.client.Run(ctx, "get", "--field", fieldName, itemName)
	resolution := provider.Resolution{Key: key, Ref: cacheKey}
	if err != nil {
		lower := strings.ToLower(err.Error() + " " + out)
		if strings.Contains(lower, "not found") || strings.Contains(lower, "missing") || strings.Contains(lower, "no item") || strings.Contains(lower, "no such field") || strings.Contains(lower, "field not found") {
			resolution.Source = "missing"
			resolution.IssueCode = "E005"
		} else {
			resolution.Source = "error"
			resolution.IssueCode = "E003"
		}
	} else {
		resolution.Source = "provider"
		resolution.Value = out
	}
	a.mu.Lock()
	a.cache[cacheKey] = resolution
	a.mu.Unlock()
	return resolution, nil
}

func (a *Adapter) ResolveMany(ctx context.Context, refs map[string]string) (map[string]provider.Resolution, error) {
	if a.cfg.UsesNoteJSON() {
		return a.resolveManyNoteJSON(ctx, refs)
	}
	results := make(map[string]provider.Resolution, len(refs))
	for key, ref := range refs {
		res, err := a.Resolve(ctx, key, ref)
		if err != nil {
			return nil, err
		}
		results[key] = res
	}
	return results, nil
}

func (a *Adapter) LoadEnvPayload(ctx context.Context) (provider.EnvPayload, error) {
	if !a.cfg.UsesNoteJSON() {
		return provider.EnvPayload{ItemName: a.itemName(), StorageMode: a.cfg.StorageMode, Env: map[string]string{}}, report.NewAppError("E009", report.ExitOperational, "ds push requires storage_mode: note_json", "push cannot safely write into field-based Bitwarden items", "set storage_mode: note_json in .envsync.yaml, migrate the repo item, and retry", nil)
	}
	if payload, ok := a.cachedNotePayload(); ok {
		return payload, nil
	}
	itemName := a.itemName()
	rawItem, err := a.client.GetRawItem(ctx, itemName)
	if err != nil {
		if errors.Is(err, ErrItemNotFound) {
			payload := provider.EnvPayload{ItemName: itemName, StorageMode: a.cfg.StorageMode, Env: map[string]string{}}
			a.cacheNotePayload(payload)
			return payload, nil
		}
		return provider.EnvPayload{}, report.NewAppError("E003", report.ExitOperational, "provider payload could not be loaded", "sync and push cannot read the repo-scoped Bitwarden item", "check rbw and retry", err)
	}
	envelope, err := syncpkg.ParseNoteJSON(rawItem.Notes)
	if err != nil {
		return provider.EnvPayload{}, report.NewAppError("E010", report.ExitOperational, "provider note_json payload is malformed", "sync and push cannot trust the repo-scoped Bitwarden payload", "repair or recreate the Bitwarden item notes and retry", err)
	}
	payload := provider.EnvPayload{
		ItemName:    itemName,
		StorageMode: a.cfg.StorageMode,
		Exists:      true,
		Format:      envelope.Format,
		Notes:       rawItem.Notes,
		Password:    rawItem.Password,
		Env:         syncpkg.CanonicalEnvMap(envelope.Env),
	}
	a.cacheNotePayload(payload)
	return payload, nil
}

func (a *Adapter) StoreEnvPayload(ctx context.Context, payload provider.EnvPayload) (provider.WriteResult, error) {
	notes, err := syncpkg.RenderNoteJSON(payload.Env)
	if err != nil {
		return provider.WriteResult{}, report.NewAppError("E010", report.ExitOperational, "provider note_json payload is malformed", "push cannot serialize the current .env into the repo-scoped payload", "fix the local env values and retry", err)
	}
	itemName := payload.ItemName
	if strings.TrimSpace(itemName) == "" {
		itemName = a.itemName()
	}
	if payload.Exists {
		if err := a.client.EditItem(ctx, itemName, payload.Password, notes); err != nil {
			return provider.WriteResult{}, report.NewAppError("E003", report.ExitOperational, "Bitwarden write failed", "push could not update the repo-scoped provider payload", "check rbw and retry", err)
		}
	} else {
		if err := a.client.AddItem(ctx, itemName, payload.Password, notes); err != nil {
			return provider.WriteResult{}, report.NewAppError("E003", report.ExitOperational, "Bitwarden write failed", "push could not create the repo-scoped provider payload", "check rbw and retry", err)
		}
	}
	if err := a.client.Sync(ctx); err != nil {
		return provider.WriteResult{}, report.NewAppError("E003", report.ExitOperational, "Bitwarden sync failed", "push could not refresh the local provider cache", "run 'rbw sync' and retry", err)
	}
	a.cacheNotePayload(provider.EnvPayload{
		ItemName:    itemName,
		StorageMode: a.cfg.StorageMode,
		Exists:      true,
		Format:      syncpkg.NoteJSONFormat,
		Notes:       notes,
		Password:    payload.Password,
		Env:         syncpkg.CanonicalEnvMap(payload.Env),
	})
	return provider.WriteResult{ItemName: itemName, Created: !payload.Exists, Updated: payload.Exists}, nil
}

func (a *Adapter) resolveManyNoteJSON(ctx context.Context, refs map[string]string) (map[string]provider.Resolution, error) {
	payload, err := a.LoadEnvPayload(ctx)
	if err != nil {
		return nil, err
	}
	results := make(map[string]provider.Resolution, len(refs))
	for key := range refs {
		results[key] = resolutionFromPayload(payload, key, a.itemName())
	}
	return results, nil
}

func (a *Adapter) resolveNoteJSON(ctx context.Context, key string) (provider.Resolution, error) {
	payload, err := a.LoadEnvPayload(ctx)
	if err != nil {
		return provider.Resolution{}, err
	}
	return resolutionFromPayload(payload, key, a.itemName()), nil
}

func resolutionFromPayload(payload provider.EnvPayload, key, itemName string) provider.Resolution {
	ref := itemName + "::notes::" + key
	value, ok := payload.Env[key]
	if !ok {
		return provider.Resolution{Key: key, Ref: ref, Source: "missing", IssueCode: "E005"}
	}
	return provider.Resolution{Key: key, Ref: ref, Source: "provider", Value: value}
}

func (a *Adapter) itemName() string {
	itemName := strings.TrimSpace(a.cfg.ItemName)
	if itemName == "" {
		return "dotenv-sync"
	}
	return itemName
}

func (a *Adapter) cachedNotePayload() (provider.EnvPayload, bool) {
	a.mu.Lock()
	defer a.mu.Unlock()
	if !a.notePayloadOK {
		return provider.EnvPayload{}, false
	}
	return clonePayload(a.notePayload), true
}

func (a *Adapter) cacheNotePayload(payload provider.EnvPayload) {
	cloned := clonePayload(payload)
	a.mu.Lock()
	defer a.mu.Unlock()
	a.notePayload = cloned
	a.notePayloadOK = true
	if a.cfg.UsesNoteJSON() {
		a.cache = map[string]provider.Resolution{}
		for key, value := range cloned.Env {
			cacheKey := cloned.ItemName + "::notes::" + key
			a.cache[cacheKey] = provider.Resolution{Key: key, Ref: cacheKey, Source: "provider", Value: value}
		}
	}
}

func clonePayload(payload provider.EnvPayload) provider.EnvPayload {
	payload.Env = syncpkg.CanonicalEnvMap(payload.Env)
	return payload
}
