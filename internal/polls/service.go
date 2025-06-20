package polls

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

type service struct {
	pollList []Poll
}

func NewService() PollService {
	return &service{
		pollList: make([]Poll, 0),
	}
}

func (s *service) CreatePoll(title string, options []string) (*Poll, error) {
	if notExactlyTwo(options) {
		return nil, errors.New("poll must have exactly two options")
	}

	// Create a new poll
	poll := &Poll{
		ID:      uuid.New().String(),
		Title:   title,
		Options: options,
		IsOpen:  true,
		Outcome: -1, // -1 indicates no outcome selected yet
	}

	// Add a copy of the poll to the service's poll list
	s.pollList = append(s.pollList, *poll)

	// Get the poll copy added to the list
	pointerToPollCopy := &s.pollList[len(s.pollList)-1]

	return pointerToPollCopy, nil
}

func notExactlyTwo(options []string) bool {
	return len(options) != 2
}

func (s *service) ClosePoll(pollID string) {
	for i, storedPoll := range s.pollList {
		if storedPoll.ID == pollID {
			s.pollList[i].IsOpen = false
			return
		}
	}
}

func (s *service) SelectOutcome(pollID string, outcomeIndex int) error {
	poll, err := s.GetPollById(pollID)
	if err != nil {
		return fmt.Errorf("failed to get poll by ID: %w", err)
	}

	if outcomeIndex < 0 || outcomeIndex >= len(poll.Options) {
		return errors.New("invalid outcome index")
	}

	poll.Outcome = outcomeIndex

	return nil
}

func (s *service) GetPollById(id string) (*Poll, error) {
	for i, poll := range s.pollList {
		// Skip if the poll ID does not match
		if poll.ID != id {
			continue
		}

		return &s.pollList[i], nil
	}
	return nil, errors.New("poll not found")
}

var _ PollService = (*service)(nil)
