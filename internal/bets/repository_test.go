package bets

import (
	"betting-discord-bot/internal/storage"
	"os"
	"slices"
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
	for _, testcase := range tests {
		t.Run(testcase.name, func(t *testing.T) {
			// ARRANGE: Set up the repository and db
			repo, cleanup := testcase.setup(t)
			t.Cleanup(cleanup)

			testBetRepositories(t, repo)
		})

	}
}

func testBetRepositories(t *testing.T, repo BetRepository) {
	t.Run("it should save and get a bet from the repo", func(t *testing.T) {

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

	t.Run("it should get all bets from a user", func(t *testing.T) {
		// ACT: Create multiple bets for the same user.
		userID := "user789"
		bets := []Bet{
			{PollId: "poll1", UserId: userID, SelectedOptionIndex: 0, BetStatus: Pending},
			{PollId: "poll2", UserId: userID, SelectedOptionIndex: 1, BetStatus: Pending},
			{PollId: "poll3", UserId: userID, SelectedOptionIndex: 0, BetStatus: Pending},
		}
		for _, bet := range bets {
			if err := repo.Save(&bet); err != nil {
				t.Fatalf("Failed to save bet: %v", err)
			}
		}

		// ACT: Retrieve all bets for the user.
		retrievedBets, retrieveErr := repo.GetBetsFromUser(userID)
		if retrieveErr != nil {
			t.Fatalf("Failed to get bets from user: %v", retrieveErr)
		}

		// ASSERT: Check that the retrieved bets match the original.
		if len(retrievedBets) != len(bets) {
			t.Fatalf("Expected %d bets for user %s, got %d", len(bets), userID, len(retrievedBets))
		}

		for i, bet := range retrievedBets {
			if slices.Contains(bets, bet) {
				continue
			}

			t.Errorf("Retrieved bet at index %d does not match original: got %+v, want %+v", i, bet, bets[i])
		}
	})

	t.Run("it should get all bets from a poll", func(t *testing.T) {
		// ACT: Create multiple bets for the same poll.
		pollID := "poll456"
		bets := []Bet{
			{PollId: pollID, UserId: "user1", SelectedOptionIndex: 0, BetStatus: Pending},
			{PollId: pollID, UserId: "user2", SelectedOptionIndex: 1, BetStatus: Pending},
			{PollId: pollID, UserId: "user3", SelectedOptionIndex: 0, BetStatus: Pending},
		}
		for _, bet := range bets {
			if err := repo.Save(&bet); err != nil {
				t.Fatalf("Failed to save bet: %v", err)
			}
		}

		// ACT: Retrieve all bets for the poll.
		retrievedBets, retrieveErr := repo.GetBetsByPollId(pollID)
		if retrieveErr != nil {
			t.Fatalf("Failed to get bets by PollId: %v", retrieveErr)
		}

		// ASSERT: Check that the retrieved bets match the original.
		if len(retrievedBets) != len(bets) {
			t.Fatalf("Expected %d bets for poll %s, got %d", len(bets), pollID, len(retrievedBets))
		}

		for i, bet := range retrievedBets {
			if slices.Contains(bets, bet) {
				continue
			}

			t.Errorf("Retrieved bet at index %d does not match original: got %+v, want %+v", i, bet, bets[i])
		}
	})
}
