package bets

import "testing"

func TestCreateBet(t *testing.T) {
	service := NewService()

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
	service := NewService()
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
