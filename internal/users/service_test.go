package users

import (
	"betting-discord-bot/internal/bets"
	"betting-discord-bot/internal/polls"
	"testing"
)

func TestCreateUser(t *testing.T) {
	pollMemoryRepo := polls.NewMemoryRepository()
	pollService := polls.NewService(pollMemoryRepo)
	betService := bets.NewService(pollService, nil)
	userRepo := NewMemoryRepository()
	userService := NewService(userRepo, betService)

	user, err := userService.CreateUser("12345")
	if err != nil {
		t.Fatalf("CreateUser returned an unexpected error: %v", err)
	}

	if user.DiscordID != "12345" {
		t.Errorf("Expected DiscordID to be '12345', got '%s'", user.DiscordID)
	}
}

func TestGetUserWinLoss(t *testing.T) {
	// ARRANGE (Outer Scope)
	pollMemoryRepo := polls.NewMemoryRepository()
	pollService := polls.NewService(pollMemoryRepo)
	betRepo := bets.NewMemoryRepository()
	betService := bets.NewService(pollService, betRepo)
	userRepo := NewMemoryRepository()
	userService := NewService(userRepo, betService)

	user, createUserErr := userService.CreateUser("12345")
	if createUserErr != nil {
		t.Fatalf("Setup failed: Could not create user: %v", createUserErr)
	}

	// Subtest 1: initial state
	t.Run("it should return zero wins and losses for a new user", func(t *testing.T) {
		// ACT
		winLoss, err := userService.GetWinLoss(user.ID)
		if err != nil {
			t.Fatalf("GetWinLoss returned an unexpected error: %v", err)
		}

		// ASSERT
		if winLoss.Wins != 0 || winLoss.Losses != 0 {
			t.Errorf("Expected win-loss to be (0, 0), got (%d, %d)", winLoss.Wins, winLoss.Losses)
		}
	})

	// Subtest 2: after a winning bet
	t.Run("it should count one win after a correct bet is resolved", func(t *testing.T) {
		// ARRANGE (Inner Scope)
		poll, _ := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})
		_, createBetErr := betService.CreateBet(poll.ID, user.ID, 0)
		if createBetErr != nil {
			t.Fatal("CreateBet returned an unexpected error: ", createBetErr)
		}
		pollOutcomeErr := pollService.SelectOutcome(poll.ID, 0)
		if pollOutcomeErr != nil {
			t.Fatal("SelectOutcome returned an unexpected error: ", pollOutcomeErr)
		}
		if err := betService.UpdateBetsByPollId(poll.ID); err != nil {
			t.Fatal("UpdateBetsByPollId returned an unexpected error: ", err)
		}

		// ACT
		winLoss, err := userService.GetWinLoss(user.ID)
		if err != nil {
			t.Fatalf("GetWinLoss returned an unexpected error: %v", err)
		}

		// ASSERT
		if winLoss.Wins != 1 || winLoss.Losses != 0 {
			t.Errorf("Expected win-loss to be (1, 0), got (%d, %d)", winLoss.Wins, winLoss.Losses)
		}
	})

	// Subtest 3: after a losing bet
	t.Run("it should count one loss after an incorrect bet is resolved", func(t *testing.T) {
		// ARRANGE (Inner Scope)
		poll, _ := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})
		_, createBetErr := betService.CreateBet(poll.ID, user.ID, 0)
		if createBetErr != nil {
			t.Fatal("CreateBet returned an unexpected error: ", createBetErr)
		}
		pollOutcomeErr := pollService.SelectOutcome(poll.ID, 1)
		if pollOutcomeErr != nil {
			t.Fatal("SelectOutcome returned an unexpected error: ", pollOutcomeErr)
		}
		if err := betService.UpdateBetsByPollId(poll.ID); err != nil {
			t.Fatal("UpdateBetsByPollId returned an unexpected error: ", err)
		}

		// ACT
		winLoss, err := userService.GetWinLoss(user.ID)
		if err != nil {
			t.Fatalf("GetWinLoss returned an unexpected error: %v", err)
		}

		// ASSERT
		if winLoss.Wins != 1 || winLoss.Losses != 1 {
			t.Errorf("Expected win-loss to be (1, 1), got (%d, %d)", winLoss.Wins, winLoss.Losses)
		}
	})

	t.Run("it should not count pending bets in win-loss", func(t *testing.T) {
		// ARRANGE (Inner Scope)
		poll, _ := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})
		_, createBetErr := betService.CreateBet(poll.ID, user.ID, 0)
		if createBetErr != nil {
			t.Fatal("CreateBet returned an unexpected error: ", createBetErr)
		}

		// ACT
		winLoss, err := userService.GetWinLoss(user.ID)
		if err != nil {
			t.Fatalf("GetWinLoss returned an unexpected error: %v", err)
		}

		// ASSERT
		if winLoss.Wins != 1 || winLoss.Losses != 1 {
			t.Errorf("Expected win-loss to be (1, 1), got (%d, %d)", winLoss.Wins, winLoss.Losses)
		}
	})
}
