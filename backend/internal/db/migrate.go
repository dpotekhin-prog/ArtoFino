package db

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"sort"
	"strings"
)

// RunMigrations executes all .sql files inside migrations folder.
func RunMigrations(db *sql.DB, migrationsPath string) error {
	files, err := ioutil.ReadDir(migrationsPath)
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	var migrationFiles []string
	for _, f := range files {
		if strings.HasSuffix(f.Name(), ".sql") {
			migrationFiles = append(migrationFiles, f.Name())
		}
	}

	// Sort files to ensure correct order (0001, 0002, ...)
	sort.Strings(migrationFiles)

	for _, file := range migrationFiles {
		path := filepath.Join(migrationsPath, file)

		content, err := ioutil.ReadFile(path)
		if err != nil {
			return fmt.Errorf("failed to read migration file %s: %w", file, err)
		}

		log.Printf("Running migration: %s", file)

		_, err = db.Exec(string(content))
		if err != nil {
			return fmt.Errorf("failed to execute migration %s: %w", file, err)
		}
	}

	log.Println("All migrations applied successfully")
	return nil
}
