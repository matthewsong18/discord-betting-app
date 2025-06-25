package polls

type PollService interface {
	CreatePoll(title string, options []string) (*Poll, error)
	ClosePoll(pollID string)
	SelectOutcome(pollID string, outcomeIndex int) error
	GetPollById(id string) (*Poll, error)
}

type PollRepository interface {
	Save(poll Poll) error
	GetById(id string) (Poll, error)
	Update(poll Poll) error
	Delete(pollID string) error
}
