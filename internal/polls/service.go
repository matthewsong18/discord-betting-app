package polls

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
)

const NoOutcomeYet = -1

type service struct {
	pollRepo PollRepository
}

func NewService(pollRepo PollRepository) PollService {
	return &service{
		pollRepo,
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
		Status:  Open,
		Outcome: NoOutcomeYet,
	}

	err := s.pollRepo.Save(poll)
	if err != nil {
		return nil, err
	}

	return poll, nil
}

func notExactlyTwo(options []string) bool {
	return len(options) != 2
}

func (s *service) ClosePoll(pollID string) error {
	poll, err := s.GetPollById(pollID)
	if err != nil {
		return fmt.Errorf("failed to get poll by ID: %w", err)
	}

	poll.Status = Closed
	if err := s.pollRepo.Update(poll); err != nil {
		return fmt.Errorf("failed to update poll status: %w", err)
	}

	return nil
}

func (s *service) SelectOutcome(pollID string, outcomeIndex int) error {
	poll, err := s.pollRepo.GetById(pollID)
	if err != nil {
		return fmt.Errorf("failed to get poll by ID: %w", err)
	}

	if outcomeIndex < 0 || outcomeIndex >= len(poll.Options) {
		return errors.New("invalid outcome index")
	}

	poll.Outcome = outcomeIndex

	if err := s.pollRepo.Update(poll); err != nil {
		return fmt.Errorf("failed to update poll outcome: %w", err)
	}

	return nil
}

func (s *service) GetPollById(id string) (*Poll, error) {
	poll, err := s.pollRepo.GetById(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get poll by ID: %w", err)
	}

	return poll, nil
}

var _ PollService = (*service)(nil)
