// internal/bets/repository_test.go

package bets

import (
	"betting-discord-bot/internal/storage"
	"os"
	"strings"
	"testing"
)

// setupLibSQL is a helper function specifically for the LibSQL implementation.
func setupLibSQL(t *testing.T) (BetRepository, func()) {
	t.Helper()

	// Sanitize the test name to create a clean, unique filename for each test run.
	sanitizedTestName := strings.ReplaceAll(t.Name(), "/", "_")
	dbPath := sanitizedTestName + ".db"

	// Proactively remove any old database file from a previous failed run.
	_ = os.Remove(dbPath)

	db, err := storage.InitializeDatabase(dbPath, "")
	if err != nil {
		t.Fatalf("Failed to initialize test database: %v", err)
	}

	repo := NewLibSQLRepository(db)

	teardown := func() {
		if err := db.Close(); err != nil {
			t.Fatal("failed to close database")
		}
		if err := os.Remove(dbPath); err != nil {
			t.Fatal("failed to remove database file")
		}
	}

	return repo, teardown
}

// setupInMemory is a helper function for the in-memory implementation.
func setupInMemory(t *testing.T) (BetRepository, func()) {
	t.Helper()
	repo := NewMemoryRepository()
	teardown := func() {
		// No cleanup needed for the in-memory version
	}
	return repo, teardown
}

// TestBetRepositoryImplementations is the main entry point for testing all BetRepository implementations.
func TestBetRepositoryImplementations(t *testing.T) {
	// This table defines all the implementations we want to test.
	implementations := []struct {
		name  string
		setup func(t *testing.T) (BetRepository, func())
	}{
		{name: "InMemoryRepository", setup: setupInMemory},
		{name: "LibSQLRepository", setup: setupLibSQL},
	}

	// This table defines all the behavioral tests we want to run against each implementation.
	testCases := []struct {
		name string
		run  func(t *testing.T, repo BetRepository)
	}{
		{"it should save and get a bet", testSaveAndGet},
		{"it should get all bets from a user", testGetAllBetsFromUser},
		{"it should get all bets from a poll", testGetAllBetsFromPoll},
	}

	// Loop through each implementation and run each test against it. Did this
	// because I needed the setup/teardown to run for each individual test instead of
	// once per implementation.
	for _, impl := range implementations {
		t.Run(impl.name, func(t *testing.T) {
			for _, tc := range testCases {
				t.Run(tc.name, func(t *testing.T) {
					repo, cleanup := impl.setup(t)
					t.Cleanup(cleanup)

					// Run the actual test logic.
					tc.run(t, repo)
				})
			}
		})
	}
}

func testSaveAndGet(t *testing.T, repo BetRepository) {
	// ARRANGE
	bet := &Bet{
		PollId:              "poll123",
		UserId:              "user456",
		SelectedOptionIndex: 0,
		BetStatus:           Pending,
	}

	// ACT & ASSERT (Save)
	if err := repo.Save(bet); err != nil {
		t.Fatalf("Failed to save bet: %v", err)
	}

	// ACT & ASSERT (Get)
	retrievedBet, err := repo.GetByPollIdAndUserId(bet.PollId, bet.UserId)
	if err != nil {
		t.Fatalf("Failed to get bet by PollId and UserId: %v", err)
	}
	if retrievedBet == nil {
		t.Fatal("Retrieved bet is nil, expected a valid bet")
	}
	if retrievedBet.PollId != bet.PollId || retrievedBet.UserId != bet.UserId {
		t.Errorf("Retrieved bet does not match original: got %+v, want %+v", retrievedBet, bet)
	}
}

func testGetAllBetsFromUser(t *testing.T, repo BetRepository) {
	// ARRANGE
	userID := "user789"
	bets := []Bet{
		{PollId: "poll1", UserId: userID, SelectedOptionIndex: 0, BetStatus: Pending},
		{PollId: "poll2", UserId: userID, SelectedOptionIndex: 1, BetStatus: Pending},
	}
	for _, bet := range bets {
		if err := repo.Save(&bet); err != nil {
			t.Fatalf("Failed to save bet: %v", err)
		}
	}

	// ACT
	retrievedBets, err := repo.GetBetsFromUser(userID)
	if err != nil {
		t.Fatalf("Failed to get bets from user: %v", err)
	}

	// ASSERT
	if len(retrievedBets) != len(bets) {
		t.Fatalf("Expected %d bets for user %s, got %d", len(bets), userID, len(retrievedBets))
	}
}

func testGetAllBetsFromPoll(t *testing.T, repo BetRepository) {
	// ARRANGE
	pollID := "poll456"
	bets := []Bet{
		{PollId: pollID, UserId: "user1", SelectedOptionIndex: 0, BetStatus: Pending},
		{PollId: pollID, UserId: "user2", SelectedOptionIndex: 1, BetStatus: Pending},
	}
	for _, bet := range bets {
		if err := repo.Save(&bet); err != nil {
			t.Fatalf("Failed to save bet: %v", err)
		}
	}

	// ACT
	retrievedBets, err := repo.GetBetsByPollId(pollID)
	if err != nil {
		t.Fatalf("Failed to get bets by PollId: %v", err)
	}

	// ASSERT
	if len(retrievedBets) != len(bets) {
		t.Fatalf("Expected %d bets for poll %s, got %d", len(bets), pollID, len(retrievedBets))
	}
}
