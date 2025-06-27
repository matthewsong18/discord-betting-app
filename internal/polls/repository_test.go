package polls

import (
	"betting-discord-bot/internal/storage"
	"github.com/google/uuid"
	"os"
	"testing"
)

func TestPollRepositoryImplementations(t *testing.T) {
	// This "table" defines all the implementations we want to test.
	testCases := []struct {
		name  string
		setup func(t *testing.T) (PollRepository, func())
	}{
		{
			name: "InMemoryRepository",
			setup: func(t *testing.T) (PollRepository, func()) {
				repo := NewMemoryRepository()
				teardown := func() {}
				return repo, teardown
			},
		},
		{
			name: "LibSQLRepository",
			setup: func(t *testing.T) (PollRepository, func()) {
				dbPath := "PollsTest.db"
				db, err := storage.InitializeDatabase(dbPath, "")
				if err != nil {
					t.Fatalf("Failed to initialize test database: %v", err)
				}
				repo := NewLibSQLRepository(db)

				teardown := func() {
					if err := db.Close(); err != nil {
						t.Errorf("Failed to close database: %v", err)
					}
					if err := os.Remove(dbPath); err != nil {
						t.Errorf("Failed to remove test database file %s: %v", dbPath, err)
					}
				}
				return repo, teardown
			},
		},
	}

	// Now, we loop through our implementations.
	for _, testCase := range testCases {
		// Use t.Run to create a subtest for each implementation.
		// This gives us clean, organized test output.
		t.Run(testCase.name, func(t *testing.T) {
			// ARRANGE: Call the specific setup function for this implementation.
			repo, teardown := testCase.setup(t)
			t.Cleanup(teardown)

			// ACT & ASSERT: Run the universal test suite against the configured repo.
			testPollRepository(t, repo)
		})
	}
}

func testPollRepository(t *testing.T, repo PollRepository) {
	t.Helper()

	t.Run("it should save and retrieve the poll", func(t *testing.T) {
		// ARRANGE: Create a new poll to save
		pollToSave := &Poll{
			ID:    uuid.New().String(),
			Title: "First Poll",
			Options: []string{
				"Option 1",
				"Option 2",
			},
			Outcome: 0,
			Status:  Open,
		}

		// ACT: Save the poll
		err := repo.Save(pollToSave)
		if err != nil {
			t.Fatalf("Save() returned an unexpected error: %v", err)
		}

		// ACT: Get the poll back
		retrievedPoll, err := repo.GetById(pollToSave.ID)
		if err != nil {
			t.Fatalf("GetById() returned an unexpected error: %v", err)
		}

		// ASSERT
		if retrievedPoll.ID != pollToSave.ID {
			t.Errorf("Expected poll ID %s, but got %s", pollToSave.ID, retrievedPoll.ID)
		}
		if retrievedPoll.Title != pollToSave.Title {
			t.Errorf("Expected poll title %q, but got %q", pollToSave.Title, retrievedPoll.Title)
		}
		if retrievedPoll.Status != pollToSave.Status {
			t.Errorf("Expected poll status %v, but got %v", pollToSave.Status, retrievedPoll.Status)
		}
		if retrievedPoll.Outcome != pollToSave.Outcome {
			t.Errorf("Expected poll outcome %v, but got %v", pollToSave.Outcome, retrievedPoll.Outcome)
		}
	})

	t.Run("it should update the poll", func(t *testing.T) {
		pollToUpdate := &Poll{
			ID:     uuid.NewString(),
			Title:  "Poll to Update",
			Status: Open,
		}

		// ACT: Save the poll first
		if err := repo.Save(pollToUpdate); err != nil {
			t.Fatalf("Save() returned an unexpected error: %v", err)
		}

		// ACT: Update the poll
		pollToUpdate.Title = "Updated Poll Title"
		if err := repo.Update(pollToUpdate); err != nil {
			t.Fatalf("Update() returned an unexpected error: %v", err)
		}

		// ACT: Retrieve the updated poll
		retrievedPoll, err := repo.GetById(pollToUpdate.ID)
		if err != nil {
			t.Fatalf("GetById() returned an unexpected error: %v", err)
		}

		// ASSERT
		if retrievedPoll.Title != "Updated Poll Title" {
			t.Errorf("Expected updated poll title %q, but got %q", "Updated Poll Title", retrievedPoll.Title)
		}

	})
}
