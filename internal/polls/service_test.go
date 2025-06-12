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
