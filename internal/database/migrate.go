package database

import (
	"database/sql"
	"fmt"
	"os"
)

const DefaultSchemaPath = "database/schema.sql"

func ApplySchema(db *sql.DB, schemaPath string) error {
	if schemaPath == "" {
		schemaPath = DefaultSchemaPath
	}

	schema, err := os.ReadFile(schemaPath)
	if err != nil {
		return fmt.Errorf("read schema file: %w", err)
	}

	if _, err := db.Exec(string(schema)); err != nil {
		return fmt.Errorf("apply schema: %w", err)
	}

	return nil
}
