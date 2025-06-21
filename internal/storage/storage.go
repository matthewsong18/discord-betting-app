package storage

import (
	"database/sql"
	"fmt"
	"log"
	"net/url"

	_ "github.com/tursodatabase/go-libsql"
)

func InitializeDatabase(dbPath, encryptionKey string) (*sql.DB, error) {
	dataSourceName := buildDatabaseConnection(dbPath, encryptionKey)

	log.Printf("Attempting sql.Open with DSN: %s", dataSourceName)

	db, err := sql.Open("libsql", dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	err = db.Ping()
	if err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	log.Println("Database ping successful!")

	schemaStatements := buildSchema()
	for _, statement := range schemaStatements {
		if _, err := db.Exec(statement); err != nil {
			return nil, fmt.Errorf("failed to execute schema statement: %w", err)
		}
	}

	return db, nil
}

func buildDatabaseConnection(dbPath string, encryptionKey string) string {
	query := make(url.Values)

	if encryptionKey != "" {
		query.Set("encryptionKey", encryptionKey)
	}

	dataSourceNameURL := url.URL{
		Scheme:   "file",
		Opaque:   dbPath,
		RawQuery: query.Encode(),
	}

	dataSourceName := dataSourceNameURL.String()
	return dataSourceName
}

func buildSchema() []string {
	return []string{
		`CREATE TABLE IF NOT EXISTS polls (
			id TEXT PRIMARY KEY,
			title TEXT,
			is_open BOOLEAN,
			outcome INTEGER
		);`,
		`CREATE TABLE IF NOT EXISTS bets (
			id TEXT PRIMARY KEY,
			poll_id TEXT,
			user_id TEXT,
			outcome_index INTEGER,
			status TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS poll_options (
            poll_id TEXT,
            option_index INTEGER,
            option_text TEXT,
            PRIMARY KEY (poll_id, option_index)
        );`,
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			discord_id TEXT UNIQUE
		);`,
	}
}
