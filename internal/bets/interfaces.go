package bets

type BetService interface {
	CreateBet(pollId string, userId int, selectedOptionIndex int) (*Bet, error)
}
