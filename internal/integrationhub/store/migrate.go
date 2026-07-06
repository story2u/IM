package store

import (
	"context"
	"crypto/sha256"
	"embed"
	"encoding/hex"
	"errors"
	"fmt"
	"path/filepath"
	"sort"
	"strings"

	"github.com/jackc/pgx/v5"
)

//go:embed migrations/*.sql
var migrationFiles embed.FS

func (s *PostgresStore) Migrate(ctx context.Context) error {
	if _, err := s.pool.Exec(ctx, `
CREATE TABLE IF NOT EXISTS integration_schema_migrations (
    version TEXT PRIMARY KEY,
    checksum TEXT NOT NULL,
    applied_at TIMESTAMPTZ NOT NULL DEFAULT now()
)`); err != nil {
		return err
	}
	entries, err := migrationFiles.ReadDir("migrations")
	if err != nil {
		return err
	}
	names := make([]string, 0, len(entries))
	for _, entry := range entries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".sql") {
			names = append(names, entry.Name())
		}
	}
	sort.Strings(names)
	for _, name := range names {
		path := filepath.Join("migrations", name)
		content, err := migrationFiles.ReadFile(path)
		if err != nil {
			return err
		}
		if err := s.applyMigration(ctx, name, string(content)); err != nil {
			return err
		}
	}
	return nil
}

func (s *PostgresStore) applyMigration(ctx context.Context, version string, sqlText string) error {
	checksum := sha256.Sum256([]byte(sqlText))
	checksumHex := hex.EncodeToString(checksum[:])
	var existing string
	err := s.pool.QueryRow(ctx, `SELECT checksum FROM integration_schema_migrations WHERE version = $1`, version).Scan(&existing)
	if err == nil {
		if existing != checksumHex {
			return fmt.Errorf("migration %s checksum changed", version)
		}
		return nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		return err
	}
	tx, err := s.pool.Begin(ctx)
	if err != nil {
		return err
	}
	defer rollback(ctx, tx)
	for _, statement := range splitSQL(sqlText) {
		if _, err := tx.Exec(ctx, statement); err != nil {
			return fmt.Errorf("apply migration %s: %w", version, err)
		}
	}
	if _, err := tx.Exec(ctx, `INSERT INTO integration_schema_migrations (version, checksum) VALUES ($1, $2)`, version, checksumHex); err != nil {
		return err
	}
	return tx.Commit(ctx)
}

func splitSQL(sqlText string) []string {
	parts := strings.Split(sqlText, ";")
	statements := make([]string, 0, len(parts))
	for _, part := range parts {
		statement := strings.TrimSpace(part)
		if statement == "" {
			continue
		}
		statements = append(statements, statement)
	}
	return statements
}
