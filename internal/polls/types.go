package polls

type Poll struct {
	ID      string
	Title   string
	Options []string
	Status  PollStatus
	Outcome OutcomeStatus
}

type PollStatus int

const (
	Open PollStatus = iota
	Closed
)

type OutcomeStatus int

const (
	Option1 OutcomeStatus = iota
	Option2
	Pending
)
