package platform

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

type MigrationInfo struct {
	Version int
	Name    string
	UpPath  string
}

func RunMigrations(ctx context.Context, pool *pgxpool.Pool, migrationsDir string) error {
	slog.Info("checking database migrations", "directory", migrationsDir)

	// Ensure schema_migrations table exists (standard golang-migrate format)
	_, err := pool.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version BIGINT PRIMARY KEY,
			dirty BOOLEAN NOT NULL
		);
	`)
	if err != nil {
		return fmt.Errorf("failed to ensure schema_migrations table: %w", err)
	}

	// Read current version
	var currentVersion int64 = 0
	var dirty bool = false
	row := pool.QueryRow(ctx, "SELECT version, dirty FROM schema_migrations LIMIT 1")
	if err := row.Scan(&currentVersion, &dirty); err != nil {
		// table is empty, start at 0
		currentVersion = 0
	}

	if dirty {
		return fmt.Errorf("database migration state is dirty at version %d; manual intervention required", currentVersion)
	}

	// Discover .up.sql files
	files, err := os.ReadDir(migrationsDir)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var pending []MigrationInfo
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".up.sql") {
			continue
		}
		parts := strings.SplitN(f.Name(), "_", 2)
		if len(parts) < 2 {
			continue
		}
		version, err := strconv.Atoi(parts[0])
		if err != nil {
			continue
		}
		if int64(version) > currentVersion {
			pending = append(pending, MigrationInfo{
				Version: version,
				Name:    f.Name(),
				UpPath:  filepath.Join(migrationsDir, f.Name()),
			})
		}
	}

	sort.Slice(pending, func(i, j int) bool {
		return pending[i].Version < pending[j].Version
	})

	if len(pending) == 0 {
		slog.Info("database schema is up to date", "current_version", currentVersion)
		return nil
	}

	for _, m := range pending {
		slog.Info("applying migration", "version", m.Version, "file", m.Name)
		content, err := os.ReadFile(m.UpPath)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", m.Name, err)
		}

		tx, err := pool.Begin(ctx)
		if err != nil {
			return fmt.Errorf("failed to begin transaction for migration %d: %w", m.Version, err)
		}

		if _, err := tx.Exec(ctx, string(content)); err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("failed executing migration %d: %w", m.Version, err)
		}

		// Update schema_migrations
		_, err = tx.Exec(ctx, `
			INSERT INTO schema_migrations (version, dirty)
			VALUES ($1, false)
			ON CONFLICT (version) DO UPDATE SET dirty = false;
		`, m.Version)
		if err != nil {
			_ = tx.Rollback(ctx)
			return fmt.Errorf("failed updating schema_migrations for version %d: %w", m.Version, err)
		}

		// Clean up older records so only the latest remains (like golang-migrate)
		_, _ = tx.Exec(ctx, `DELETE FROM schema_migrations WHERE version < $1`, m.Version)

		if err := tx.Commit(ctx); err != nil {
			return fmt.Errorf("failed committing migration %d: %w", m.Version, err)
		}

		slog.Info("successfully applied migration", "version", m.Version)
	}

	return nil
}
