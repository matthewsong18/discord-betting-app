package bets

import "betting-discord-bot/internal/polls"

type BetService interface {
	CreateBet(pollID string, userID string, selectedOptionIndex int) (*Bet, error)
	GetBet(pollID string, userID string) (*Bet, error)
	UpdateBetsByPollId(poll polls.Poll)
	GetBetsFromUser(userID string) ([]Bet, error)
}
