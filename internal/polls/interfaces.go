package polls

type PollService interface {
	CreatePoll(title string, options []string) (*Poll, error)
	ClosePoll(pollID string)
	SelectOutcome(poll *Poll, outcomeIndex int) error
	GetPollById(id string) (*Poll, error)
}
