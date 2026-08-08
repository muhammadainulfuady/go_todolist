package database

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"github.com/sirupsen/logrus"
)

const schemaMigrationsDDL = `
CREATE TABLE IF NOT EXISTS schema_migrations (
    version     VARCHAR(20)  NOT NULL PRIMARY KEY,
    filename    VARCHAR(255) NOT NULL,
    applied_at  DATETIME     NOT NULL DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci`

type migration struct {
	Version  int
	Filename string
}

// RunMigrations mengeksekusi semua file *.up.sql di dir secara berurutan
// (000001 -> 000002 -> ...) yang belum tercatat di schema_migrations.
func RunMigrations(db *sql.DB, dir string) error {
	if _, err := db.Exec(schemaMigrationsDDL); err != nil {
		return fmt.Errorf("gagal membuat tabel schema_migrations: %w", err)
	}

	files, err := migrationFiles(dir, ".up.sql", false)
	if err != nil {
		return err
	}

	applied, err := appliedVersions(db)
	if err != nil {
		return err
	}

	for _, m := range files {
		if applied[m.Version] {
			logrus.WithField("file", m.Filename).Info("skip (sudah diterapkan)")
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, m.Filename))
		if err != nil {
			return fmt.Errorf("gagal membaca %s: %w", m.Filename, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("gagal eksekusi migrasi %s: %w", m.Filename, err)
		}

		if _, err := db.Exec(
			"INSERT INTO schema_migrations (version, filename) VALUES (?, ?)",
			fmt.Sprintf("%06d", m.Version), m.Filename,
		); err != nil {
			return fmt.Errorf("gagal mencatat versi %d: %w", m.Version, err)
		}

		logrus.WithField("file", m.Filename).Info("diterapkan")
	}

	logrus.Info("Migrasi up selesai.")
	return nil
}

// DropMigrations mengeksekusi file *.down.sql secara terbalik
// (000003 -> 000002 -> ...) lalu menghapus catatan dari schema_migrations.
func DropMigrations(db *sql.DB, dir string) error {
	if _, err := db.Exec(schemaMigrationsDDL); err != nil {
		return fmt.Errorf("gagal membuat tabel schema_migrations: %w", err)
	}

	files, err := migrationFiles(dir, ".down.sql", true)
	if err != nil {
		return err
	}

	applied, err := appliedVersions(db)
	if err != nil {
		return err
	}

	for _, m := range files {
		if !applied[m.Version] {
			logrus.WithField("file", m.Filename).Info("skip (belum diterapkan)")
			continue
		}

		content, err := os.ReadFile(filepath.Join(dir, m.Filename))
		if err != nil {
			return fmt.Errorf("gagal membaca %s: %w", m.Filename, err)
		}

		if _, err := db.Exec(string(content)); err != nil {
			return fmt.Errorf("gagal eksekusi rollback %s: %w", m.Filename, err)
		}

		if _, err := db.Exec(
			"DELETE FROM schema_migrations WHERE version = ?",
			fmt.Sprintf("%06d", m.Version),
		); err != nil {
			return fmt.Errorf("gagal menghapus catatan versi %d: %w", m.Version, err)
		}

		logrus.WithField("file", m.Filename).Info("di-rollback")
	}

	logrus.Info("Migrasi down selesai.")
	return nil
}

func migrationFiles(dir, suffix string, descending bool) ([]migration, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("gagal membaca direktori %s: %w", dir, err)
	}

	var files []migration
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), suffix) {
			continue
		}
		versionStr := strings.SplitN(e.Name(), "_", 2)[0]
		version, err := strconv.Atoi(versionStr)
		if err != nil {
			return nil, fmt.Errorf("prefix versi tidak valid pada %s: %w", e.Name(), err)
		}
		files = append(files, migration{Version: version, Filename: e.Name()})
	}

	sort.Slice(files, func(i, j int) bool {
		if descending {
			return files[i].Version > files[j].Version
		}
		return files[i].Version < files[j].Version
	})

	return files, nil
}

func appliedVersions(db *sql.DB) (map[int]bool, error) {
	rows, err := db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return nil, fmt.Errorf("gagal query schema_migrations: %w", err)
	}
	defer rows.Close()

	applied := make(map[int]bool)
	for rows.Next() {
		var v string
		if err := rows.Scan(&v); err != nil {
			return nil, err
		}
		version, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		applied[version] = true
	}
	return applied, rows.Err()
}
