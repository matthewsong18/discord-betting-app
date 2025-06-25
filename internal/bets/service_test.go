package bets

import (
	"betting-discord-bot/internal/polls"
	"testing"
)

func TestCreateBet(t *testing.T) {
	pollService := polls.NewService()
	betService := NewService(pollService)

	poll, err := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})
	if err != nil {
		t.Fatal("Failed to create poll:", err)
	}

	pollId := poll.ID
	userId := "12345"
	selectedOptionIndex := 0
	bet, err1 := betService.CreateBet(pollId, userId, selectedOptionIndex)

	if err1 != nil {
		t.Fatal("CreateBet returned an unexpected error:", err1)
	}

	if bet.PollId != pollId {
		t.Errorf("Expected bet to be associated with poll %s, but got %s", pollId, bet.PollId)
	}

	if bet.UserId != userId {
		t.Errorf("Expected bet to be associated with user %s, but got %s", userId, bet.UserId)
	}

	if bet.SelectedOptionIndex != selectedOptionIndex {
		t.Errorf("Expected bet to select option %d, but got %d", selectedOptionIndex, bet.SelectedOptionIndex)
	}
}

func TestInvalidOption(t *testing.T) {
	pollService := polls.NewService()
	betService := NewService(pollService)
	pollId := "12345"
	userId := "12345"
	selectedOptionIndex := -1 // Invalid index
	_, err := betService.CreateBet(pollId, userId, selectedOptionIndex)

	if err == nil {
		t.Fatal("Expected CreateBet to return an error for invalid option index, but got nil")
	}

	if err.Error() != "invalid option index" {
		t.Errorf("Expected error message 'invalid option index', but got '%s'", err.Error())
	}
}

func TestPreventingMultipleBetsPerPoll(t *testing.T) {
	pollService := polls.NewService()
	betService := NewService(pollService)

	poll, _ := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})

	// Create the first bet for the poll
	pollId := poll.ID
	userId := "12345"
	selectedOptionIndex := 0

	_, _ = betService.CreateBet(pollId, userId, selectedOptionIndex)

	// Attempt to create a second bet for the same poll
	_, err := betService.CreateBet(pollId, userId, selectedOptionIndex)

	if err == nil {
		t.Fatal("Expected an error when creating a second bet for the same poll, but got nil")
	}

	if err.Error() != "user bet already exists for this poll" {
		t.Errorf("Expected error message 'user bet already exists for this poll', but got '%s'", err.Error())
	}
}

func TestCannotBetOnClosedPoll(t *testing.T) {
	pollService := polls.NewService()
	betService := NewService(pollService)

	poll, err := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})
	if err != nil {
		t.Fatal("Failed to create poll:", err)
	}
	pollService.ClosePoll(poll.ID)

	// Attempt to create a bet on a closed poll
	_, err = betService.CreateBet(poll.ID, "12345", 0)
	if err == nil {
		t.Fatal("Expected an error when betting on a closed poll, but got nil")
	}

	// Check if the error message is as expected
	if err.Error() != "cannot bet on a closed poll" {
		t.Errorf("Expected error message 'cannot bet on a closed poll', but got '%s'", err.Error())
	}
}

func TestGetBetOutcome(t *testing.T) {
	// Check if the bet outcome is correctly retrieved

	pollService := polls.NewService()
	betService := NewService(pollService)
	poll, err := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})
	if err != nil {
		t.Fatal("Failed to create poll:", err)
	}

	pollId := poll.ID
	userId := "12345"
	selectedOptionIndex := 0
	bet, err1 := betService.CreateBet(pollId, userId, selectedOptionIndex)

	if err1 != nil {
		t.Fatal("CreateBet returned an unexpected error:", err1)
	}

	if bet.BetStatus != Pending {
		t.Fatalf("Expected bet status to be 'PENDING', but got '%s'", bet.BetStatus)
	}

	pollService.ClosePoll(poll.ID)
	err = pollService.SelectOutcome(poll.ID, selectedOptionIndex)
	if err != nil {
		t.Fatal("SelectOutcome returned an unexpected error:", err)
	}

	betService.UpdateBetsByPollId(*poll)
	bet, err2 := betService.GetBet(pollId, userId)

	if err2 != nil {
		t.Fatal("GetBet returned an unexpected error:", err2)
	}

	if bet.BetStatus != Won {
		t.Errorf("Expected bet status to be 'WON', but got '%s'", bet.BetStatus)
	}

	if bet.SelectedOptionIndex != selectedOptionIndex {
		t.Errorf("Expected bet to select option %d, but got %d", selectedOptionIndex, bet.SelectedOptionIndex)
	}
}

func TestGettingUserBets(t *testing.T) {
	pollService := polls.NewService()
	betService := NewService(pollService)

	poll, createPollErr := pollService.CreatePoll("Test Poll", []string{"Option 1", "Option 2"})
	if createPollErr != nil {
		t.Fatal("Failed to create poll:", createPollErr)
	}

	userID := "12345"
	bet, createBetErr := betService.CreateBet(poll.ID, userID, 0)
	if createBetErr != nil {
		t.Fatal("Failed to create bet:", createBetErr)
	}

	bets, getBetsErr := betService.GetBetsFromUser(userID)
	if getBetsErr != nil {
		t.Fatal("Failed to get bets from user:", getBetsErr)
	}

	if len(bets) != 1 {
		t.Errorf("Expected 1 bet for user %s, but got %d", userID, len(bets))
	}

	if bets[0] != *bet || bets[0].PollId != poll.ID || bets[0].UserId != userID {
		t.Errorf("Expected bet %v, but got %v", bet, bets[0])
	}
}
