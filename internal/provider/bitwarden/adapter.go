package bitwarden

import (
	"context"
	"strings"
	"sync"

	"dotenv-sync/internal/config"
	"dotenv-sync/internal/provider"
)

type Adapter struct {
	client *RBWClient
	cfg    config.Config
	mu     sync.Mutex
	cache  map[string]provider.Resolution
}

func NewAdapter(cfg config.Config) *Adapter {
	return &Adapter{client: NewRBWClient(), cfg: cfg, cache: map[string]provider.Resolution{}}
}

func (a *Adapter) Name() string { return "bitwarden" }

func (a *Adapter) CheckReadiness(ctx context.Context) (provider.Status, error) {
	return checkReadinessWithClient(ctx, a.client)
}

func (a *Adapter) Resolve(ctx context.Context, key, ref string) (provider.Resolution, error) {
	if ref == "" {
		ref = key
	}
	a.mu.Lock()
	if cached, ok := a.cache[ref]; ok {
		a.mu.Unlock()
		cached.Key = key
		cached.Ref = ref
		return cached, nil
	}
	a.mu.Unlock()
	out, err := a.client.Run(ctx, "get", ref)
	resolution := provider.Resolution{Key: key, Ref: ref}
	if err != nil {
		lower := strings.ToLower(err.Error() + " " + out)
		if strings.Contains(lower, "not found") || strings.Contains(lower, "missing") || strings.Contains(lower, "no item") {
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
	a.cache[ref] = resolution
	a.mu.Unlock()
	return resolution, nil
}

func (a *Adapter) ResolveMany(ctx context.Context, refs map[string]string) (map[string]provider.Resolution, error) {
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
