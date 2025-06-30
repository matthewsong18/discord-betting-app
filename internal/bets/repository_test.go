package bets

import (
	"betting-discord-bot/internal/storage"
	"testing"
)

func TestRepositories(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) (repo BetRepository, cleanup func())
	}{
		// TODO: test cases
		{
			name: "InMemoryRepository",
			setup: func(t *testing.T) (BetRepository, func()) {
				repo := NewMemoryRepository()
				return repo, func() {
					// Cleanup if necessary
				}
			},
		},
		{
			name: "LibSqlRepository",
			setup: func(t *testing.T) (BetRepository, func()) {
				dbPath := "BetsTest.db"
				db, initErr := storage.InitializeDatabase(dbPath, "")
				if initErr != nil {
					return nil, nil
				}

				repo := NewLibSQLRepository(db)
				return repo, func() {
					if err := db.Close(); err != nil {
						t.Fatalf("Failed to close LibSqlRepository: %v", err)
					}
				}
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}
