package polls

import "errors"

type PollService interface {
	CreatePoll(title string, options []string) (Poll, error)
	ClosePoll(pollID string) error
	SelectOutcome(pollID string, outcomeIndex OutcomeStatus) error
	GetPollById(id string) (Poll, error)
	GetOpenPolls() ([]Poll, error)
}

type PollRepository interface {
	Save(poll *poll) error
	GetById(id string) (*poll, error)
	GetOpenPolls() ([]*poll, error)
	Update(poll *poll) error
	Delete(pollID string) error
}

var ErrPollIsAlreadyClosed = errors.New("poll is already closed")
