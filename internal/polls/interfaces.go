package polls

import "errors"

type PollService interface {
	CreatePoll(title string, options []string) (*Poll, error)
	ClosePoll(pollID string) error
	SelectOutcome(pollID string, outcomeIndex OutcomeStatus) error
	GetPollById(id string) (*Poll, error)
	GetOpenPolls() ([]*Poll, error)
}

type PollRepository interface {
	Save(poll *Poll) error
	GetById(id string) (*Poll, error)
	GetOpenPolls() ([]*Poll, error)
	Update(poll *Poll) error
	Delete(pollID string) error
}

var ErrPollIsAlreadyClosed = errors.New("poll is already closed")
