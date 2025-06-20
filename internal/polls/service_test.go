package polls

import "testing"

func TestCreatePoll(t *testing.T) {
	service := NewService()

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

	if !poll.IsOpen {
		t.Error("Expected poll to be open, but it was closed")
	}

}

func TestExactlyTwoOptions(t *testing.T) {
	service := NewService()

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
	service := NewService()

	poll, err := createDefaultTestPoll(service)

	if err != nil {
		t.Fatal("CreatePoll returned an unexpected error:", err)
	}
	if !poll.IsOpen {
		t.Fatal("Expected poll to be open after creation, but it was closed")
	}

	service.ClosePoll(poll.ID)

	updatedPoll, updateError := service.GetPollById(poll.ID)
	if updateError != nil {
		t.Fatal("GetPollById returned an unexpected error:", updateError)
	}

	if updatedPoll.IsOpen {
		t.Error("Expected poll to be closed after ClosePoll, but it was still open")
	}
}

func TestSelectOutcome(t *testing.T) {
	service := NewService()

	poll, err := createDefaultTestPoll(service)

	if err != nil {
		t.Fatal("CreatePoll returned an unexpected error", err)
	}

	teamAIndex := 0
	err = service.SelectOutcome(poll, teamAIndex)
	if err != nil {
		t.Fatal("SelectOutcome returned an unexpected error", err)
	}

	if poll.Outcome != teamAIndex {
		t.Errorf("Expected selected outcome to be '%d', but got '%d'", teamAIndex, poll.Outcome)
	}
}

func createDefaultTestPoll(service PollService) (*Poll, error) {
	title := "Which team will win first map?"
	options := []string{"Team A", "Team B"}
	poll, err := service.CreatePoll(title, options)
	return poll, err
}
