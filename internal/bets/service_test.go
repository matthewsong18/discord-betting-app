package bets

import (
	"betting-discord-bot/internal/polls"
	"testing"
)

func TestCreateBet(t *testing.T) {
	service := NewService(nil)

	pollId := "12345"
	userId := 12345
	selectedOptionIndex := 0
	bet, err := service.CreateBet(pollId, userId, selectedOptionIndex)

	if err != nil {
		t.Fatal("CreateBet returned an unexpected error:", err)
	}

	if bet.PollId != pollId {
		t.Errorf("Expected bet to be associated with poll %s, but got %s", pollId, bet.PollId)
	}

	if bet.UserId != userId {
		t.Errorf("Expected bet to be associated with user %d, but got %d", userId, bet.UserId)
	}

	if bet.SelectedOptionIndex != selectedOptionIndex {
		t.Errorf("Expected bet to select option %d, but got %d", selectedOptionIndex, bet.SelectedOptionIndex)
	}
}

func TestInvalidOption(t *testing.T) {
	service := NewService(nil)
	pollId := "12345"
	userId := 12345
	selectedOptionIndex := -1 // Invalid index
	_, err := service.CreateBet(pollId, userId, selectedOptionIndex)

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
	userId := 12345
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

	}
}
