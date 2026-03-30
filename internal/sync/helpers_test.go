package sync

import (
	"context"

	"dotenv-sync/internal/provider"
)

type fakeProvider struct {
	status      provider.Status
	resolutions map[string]provider.Resolution
	payload     provider.EnvPayload
	loadErr     error
	storeErr    error
	stored      provider.EnvPayload
	loadCalls   int
	storeCalls  int
}

func (f fakeProvider) Name() string { return "fake" }
func (f fakeProvider) CheckReadiness(context.Context) (provider.Status, error) {
	if f.status.Provider == "" {
		f.status.Provider = "bitwarden"
		f.status.CLIInstalled = true
		f.status.Authenticated = true
		f.status.Unlocked = true
	}
	return f.status, nil
}
func (f fakeProvider) Resolve(_ context.Context, key, ref string) (provider.Resolution, error) {
	if res, ok := f.resolutions[key]; ok {
		return res, nil
	}
	if res, ok := f.resolutions[ref]; ok {
		return res, nil
	}
	return provider.Resolution{Key: key, Ref: ref, Source: "missing", IssueCode: "E005"}, nil
}
func (f fakeProvider) ResolveMany(ctx context.Context, refs map[string]string) (map[string]provider.Resolution, error) {
	result := map[string]provider.Resolution{}
	for key, ref := range refs {
		res, _ := f.Resolve(ctx, key, ref)
		res.Key = key
		res.Ref = ref
		result[key] = res
	}
	return result, nil
}

func (f *fakeProvider) LoadEnvPayload(context.Context) (provider.EnvPayload, error) {
	f.loadCalls++
	if f.loadErr != nil {
		return provider.EnvPayload{}, f.loadErr
	}
	return f.payload, nil
}

func (f *fakeProvider) StoreEnvPayload(_ context.Context, payload provider.EnvPayload) (provider.WriteResult, error) {
	f.storeCalls++
	if f.storeErr != nil {
		return provider.WriteResult{}, f.storeErr
	}
	f.stored = payload
	return provider.WriteResult{ItemName: payload.ItemName, Created: !payload.Exists, Updated: payload.Exists}, nil
}
