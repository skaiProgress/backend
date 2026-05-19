package migrate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Runner applies SQL migrations from a filesystem directory.
type Runner struct {
	m *migrate.Migrate
}

// NewRunner creates a migration runner for the given database URL and migrations path.
func NewRunner(databaseURL, migrationsPath string) (*Runner, error) {
	absPath, err := filepath.Abs(migrationsPath)
	if err != nil {
		return nil, fmt.Errorf("resolve migrations path: %w", err)
	}

	sourceURL := "file://" + filepath.ToSlash(absPath)
	m, err := migrate.New(sourceURL, databaseURL)
	if err != nil {
		return nil, fmt.Errorf("create migrator: %w", err)
	}

	return &Runner{m: m}, nil
}

// ResolveMigrationsPath returns MIGRATIONS_PATH or a sensible default.
func ResolveMigrationsPath() string {
	if p := os.Getenv("MIGRATIONS_PATH"); p != "" {
		return p
	}
	if _, err := os.Stat("/app/migrations"); err == nil {
		return "/app/migrations"
	}
	return "migrations"
}

// Up applies all pending migrations.
func (r *Runner) Up() error {
	if err := r.m.Up(); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// Down rolls back one migration.
func (r *Runner) Down() error {
	if err := r.m.Steps(-1); err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return err
	}
	return nil
}

// Status returns current migration version information.
func (r *Runner) Status() (version uint, dirty bool, err error) {
	version, dirty, err = r.m.Version()
	if err != nil {
		return 0, false, err
	}
	return version, dirty, nil
}

// Close releases migration resources.
func (r *Runner) Close() error {
	sourceErr, dbErr := r.m.Close()
	if sourceErr != nil {
		return sourceErr
	}
	return dbErr
}
