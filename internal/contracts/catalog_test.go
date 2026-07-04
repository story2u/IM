package contracts

import (
	"path/filepath"
	"runtime"
	"testing"
)

func TestLoadCatalogReadsProjectSchemas(t *testing.T) {
	root := projectContractRoot(t)
	catalog, err := LoadCatalog(root)
	if err != nil {
		t.Fatalf("LoadCatalog() error = %v", err)
	}
	if len(catalog) != 6 {
		t.Fatalf("expected 6 project schemas, got %d", len(catalog))
	}
	if err := RequireSchemas(catalog, "agent-heartbeat.schema.json", "connector-delivery-receipt.schema.json", "connector-inbound-event.schema.json", "connector-outbound-message.schema.json", "task-create.schema.json", "task-status.schema.json"); err != nil {
		t.Fatalf("RequireSchemas() error = %v", err)
	}
}

func projectContractRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("runtime.Caller failed")
	}
	repoRoot := filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
	return filepath.Join(repoRoot, "contracts", "v1")
}
