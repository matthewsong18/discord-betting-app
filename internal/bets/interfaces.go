package bets

type BetService interface {
	CreateBet(pollID string, userID string, selectedOptionIndex int) (*Bet, error)
	GetBet(pollID string, userID string) (*Bet, error)
	UpdateBetsByPollId(pollID string) error
	GetBetsFromUser(userID string) ([]Bet, error)
}
