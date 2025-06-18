package bets

type BetStatus int

const (
	StatusPending BetStatus = iota
	StatusWon
	StatusLost
)

func (bs BetStatus) String() string {
	switch bs {
	case StatusPending:
		return "PENDING"
	case StatusWon:
		return "WON"
	case StatusLost:
		return "LOST"
	default:
		return "UNKNOWN"
	}
}

type Bet struct {
	PollId              string
	UserId              string
	SelectedOptionIndex int
	BetStatus           BetStatus
}
