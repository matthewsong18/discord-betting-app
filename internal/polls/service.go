package polls

import (
	"errors"
	"github.com/google/uuid"
)

type service struct {
}

func NewService() PollService {
	return &service{}
}

func (s *service) CreatePoll(title string, options []string) (*Poll, error) {
	if notExactlyTwo(options) {
		return nil, errors.New("poll must have exactly two options")
	}

	poll := &Poll{
		ID:      uuid.New().String(),
		Title:   title,
		Options: options,
		IsOpen:  true,
		Outcome: -1, // -1 indicates no outcome selected yet
	}

	return poll, nil
}

func notExactlyTwo(options []string) bool {
	return len(options) != 2
}

func (s *service) ClosePoll(poll *Poll) {
	poll.IsOpen = false
}

func (s *service) SelectOutcome(poll *Poll, outcomeIndex int) error {
	if outcomeIndex < 0 || outcomeIndex >= len(poll.Options) {
		return errors.New("invalid outcome index")
	}
	poll.Outcome = outcomeIndex

	return nil
}

var _ PollService = (*service)(nil)
