package polls

type Poll struct {
	ID      string
	Title   string
	Options []string
	Status  PollStatus
	Outcome int
}

type PollStatus int

const (
	Open PollStatus = iota
	Closed
)
