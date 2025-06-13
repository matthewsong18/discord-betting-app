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
	}

	return poll, nil
}

func notExactlyTwo(options []string) bool {
	return len(options) != 2
}

func (s *service) ClosePoll(poll *Poll) {
	poll.IsOpen = false
}

var _ PollService = (*service)(nil)
