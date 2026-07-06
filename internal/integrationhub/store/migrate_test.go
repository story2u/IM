package store

import (
	"strings"
	"testing"
)

func TestEmbeddedMigrationDefinesCoreTables(t *testing.T) {
	content, err := migrationFiles.ReadFile("migrations/001_init.sql")
	if err != nil {
		t.Fatalf("read migration: %v", err)
	}
	sql := string(content)
	for _, table := range []string{
		"integration_channels",
		"integration_channel_accounts",
		"integration_conversations",
		"integration_messages",
		"integration_inbound_events",
		"integration_outbound_commands",
		"integration_sop_workflows",
		"integration_ai_policies",
		"integration_audit_logs",
	} {
		if !strings.Contains(sql, "CREATE TABLE IF NOT EXISTS "+table) {
			t.Fatalf("migration missing table %s", table)
		}
	}
}

func TestSplitSQLIgnoresBlankStatements(t *testing.T) {
	statements := splitSQL("CREATE TABLE a(id text); ;\nCREATE TABLE b(id text);")
	if len(statements) != 2 {
		t.Fatalf("statement count = %d, want 2", len(statements))
	}
}
