package db

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/jackc/pgx/v5/pgxpool"
)

func RunMigrations(conn *pgxpool.Pool) error {
	files, err := filepath.Glob("migrations/*.sql")
	if err != nil {
		return err
	}

	// files MUST start with number prefix
	sort.Strings(files)

	for _, file := range files {
		fmt.Printf("Running migration: %s\n", file)

		sqlBytes, err := os.ReadFile(file)
		if err != nil {
			return err
		}

		_, err = conn.Exec(context.Background(), string(sqlBytes))
		if err != nil {
			return fmt.Errorf("migration failed (%s): %w", file, err)
		}
	}

	return nil
}
