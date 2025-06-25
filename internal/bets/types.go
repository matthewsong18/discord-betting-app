package bets

type BetStatus int

const (
	Pending BetStatus = iota
	Won
	Lost
)

func (bs BetStatus) String() string {
	switch bs {
	case Pending:
		return "PENDING"
	case Won:
		return "WON"
	case Lost:
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
