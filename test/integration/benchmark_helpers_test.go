package integration_test

import (
	"os"
	"path/filepath"
	"testing"
)

func loadLargeSchema(tb testing.TB) string {
	tb.Helper()
	data, err := os.ReadFile(filepath.Join(repoRootFromTB(tb), "test", "testdata", "env", "large-schema.env.example"))
	if err != nil {
		tb.Fatal(err)
	}
	return string(data)
}
