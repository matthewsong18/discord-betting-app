package polls

import "testing"

func TestCreatePoll(t *testing.T) {
	pollMemoryRepo := NewMemoryRepository()
	service := NewService(pollMemoryRepo)

	title := "Which team will win first map?"
	options := []string{"Team A", "Team B"}

	poll, err := service.CreatePoll(title, options)

	if err != nil {
		t.Fatalf("CreatePoll returned an unexpected error: %v", err)
	}

	if poll.ID == "" {
		t.Error("Expected poll ID to be set, but it was empty")
	}

	if poll.Title != title {
		t.Errorf("Expected poll title to be '%s', but got '%s'", title, poll.Title)
	}

	for i, option := range options {
		if poll.Options[i] != option {
			t.Errorf("Expected option %d to be '%s', but got '%s'", i, option, poll.Options[i])
		}
	}

	if poll.Status != Open {
		t.Error("Expected poll to be open, but it was closed")
	}

}

func TestExactlyTwoOptions(t *testing.T) {
	pollMemoryRepo := NewMemoryRepository()
	service := NewService(pollMemoryRepo)

	title := "Which team will win first map?"
	options := []string{"Team A", "Team B", "Team C"}

	_, err := service.CreatePoll(title, options)

	if err == nil {
		t.Fatal("Expected CreatePoll to return an error for more than two options, but it did not")
	}

	expectedError := "poll must have exactly two options"
	if err.Error() != expectedError {
		t.Errorf("Expected error message '%s', but got '%s'", expectedError, err.Error())
	}
}

func TestClosePoll(t *testing.T) {
	pollMemoryRepo := NewMemoryRepository()
	service := NewService(pollMemoryRepo)

	poll, err := createDefaultTestPoll(service)

	if err != nil {
		t.Fatal("CreatePoll returned an unexpected error:", err)
	}
	if poll.Status != Open {
		t.Fatal("Expected poll to be open after creation, but it was closed")
	}

	if err := service.ClosePoll(poll.ID); err != nil {
		t.Fatal("ClosePoll returned an unexpected error:", err)
	}

	updatedPoll, updateError := service.GetPollById(poll.ID)
	if updateError != nil {
		t.Fatal("GetPollById returned an unexpected error:", updateError)
	}

	if updatedPoll.Status != Closed {
		t.Error("Expected poll to be closed after ClosePoll, but it was still open")
	}
}

func TestSelectOutcome(t *testing.T) {
	// Setup
	pollMemoryRepo := NewMemoryRepository()
	service := NewService(pollMemoryRepo)

	poll, err := createDefaultTestPoll(service)
	if err != nil {
		t.Fatal("CreatePoll returned an unexpected error", err)
	}

	// Test selecting an outcome
	teamAIndex := 0
	err = service.SelectOutcome(poll.ID, teamAIndex)
	if err != nil {
		t.Fatal("SelectOutcome returned an unexpected error", err)
	}

	// Get the updated poll
	poll, err = service.GetPollById(poll.ID)
	if err != nil {
		t.Fatal("GetPollById returned an unexpected error:", err)
	}

	// Verify the outcome
	if poll.Outcome != teamAIndex {
		t.Errorf("Expected selected outcome to be '%d', but got '%d'", teamAIndex, poll.Outcome)
	}
}

func createDefaultTestPoll(service PollService) (Poll, error) {
	title := "Which team will win first map?"
	options := []string{"Team A", "Team B"}
	poll, err := service.CreatePoll(title, options)
	return poll, err
}

func TestGetPollById(t *testing.T) {
	// Testing that the GetPollById method retrieves the exact poll that was created by CreatePoll instead of a copy.
	pollMemoryRepo := NewMemoryRepository()
	pollService := NewService(pollMemoryRepo)

	poll, err := createDefaultTestPoll(pollService)
	if err != nil {
		t.Fatal("CreatePoll returned an unexpected error:", err)
	}

	retrievedPoll, err := pollService.GetPollById(poll.ID)
	if err != nil {
		t.Fatal("GetPollById returned an unexpected error:", err)
	}

	if retrievedPoll.ID != poll.ID {
		t.Errorf("Expected retrieved poll to be equal to created poll, but they differ")
	}
}
