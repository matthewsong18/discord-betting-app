package bets

import "betting-discord-bot/internal/polls"

type BetService interface {
	CreateBet(pollId string, userId int, selectedOptionIndex int) (*Bet, error)
	GetBet(pollId string, userId int) (*Bet, error)
	UpdateBetsByPollId(poll polls.Poll)
}
