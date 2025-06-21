package main

import (
	"betting-discord-bot/internal/storage"
	"fmt"
	"log"
	"os"
)

func run() (err error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		log.Fatal("DB_PATH environment variable is not set")
	}

	encryptionKey := os.Getenv("ENCRYPTION_KEY")
	if encryptionKey == "" {
		log.Println("ENCRYPTION_KEY environment variable is not set, using unencrypted database")
	}

	db, initDbError := storage.InitializeDatabase(dbPath, encryptionKey)

	if initDbError != nil {
		return fmt.Errorf("failed to initialize database: %w", initDbError)
	}

	log.Println("Database initialized successfully")

	// Initialize internal services here

	// Run discord logic here

	log.Println("Application started successfully")

	defer func() {
		if closeError := db.Close(); closeError != nil {
			fmt.Println("Error closing database", closeError)
			if err == nil {
				err = closeError
			}
		}
	}()

	return nil
}

func main() {
	if err := run(); err != nil {
		log.Fatalf("application failed to start: %v", err)
	}
}
