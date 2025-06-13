package polls

type PollService interface {
	CreatePoll(title string, options []string) (*Poll, error)
	ClosePoll(poll *Poll)
	SelectOutcome(poll *Poll, outcomeIndex int) error
}
