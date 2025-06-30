package bets

import (
	"betting-discord-bot/internal/storage"
	"os"
	"testing"
)

func TestRepositories(t *testing.T) {
	tests := []struct {
		name  string
		setup func(t *testing.T) (repo BetRepository, cleanup func())
	}{
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

					if err := os.RemoveAll(dbPath); err != nil {
						t.Fatalf("Failed to remove test database file %s: %v", dbPath, err)
					}
				}
			},
		},
	}
	for _, test := range tests {
		t.Run("it should save and get a bet from the repo", func(t *testing.T) {
			// ARRANGE: Call the specific setup function for this implementation.
			repo, cleanup := test.setup(t)
			t.Cleanup(cleanup)

			// ACT: Create a bet and save it to the repository.
			bet := &Bet{
				PollId:              "poll123",
				UserId:              "user456",
				SelectedOptionIndex: 0,
				BetStatus:           Pending,
			}

			if err := repo.Save(bet); err != nil {
				t.Fatalf("Failed to save bet: %v", err)
			}

			// ACT: Retrieve the bet by PollId and UserId.
			retrievedBet, err := repo.GetByPollIdAndUserId(bet.PollId, bet.UserId)
			if err != nil {
				t.Fatalf("Failed to get bet by PollId and UserId: %v", err)
			}

			// ASSERT: Check that the retrieved bet matches the original.
			if retrievedBet == nil {
				t.Fatal("Retrieved bet is nil, expected a valid bet")
			}

			if retrievedBet.PollId != bet.PollId ||
				retrievedBet.UserId != bet.UserId ||
				retrievedBet.SelectedOptionIndex != bet.SelectedOptionIndex ||
				retrievedBet.BetStatus != bet.BetStatus {

				t.Errorf("Retrieved bet does not match original: got %+v, want %+v", retrievedBet, bet)
			}
		})
	}
}
