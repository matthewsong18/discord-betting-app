package bets

import "errors"

type BetService interface {
	CreateBet(pollID string, userID string, selectedOptionIndex int) (*Bet, error)
	GetBet(pollID string, userID string) (*Bet, error)
	UpdateBetsByPollId(pollID string) error
	GetBetsFromUser(userID string) ([]Bet, error)
}

type BetRepository interface {
	Save(bet *Bet) error
	GetByPollIdAndUserId(pollID string, userID string) (*Bet, error)
	GetBetsFromUser(userID string) ([]Bet, error)
	GetBetsByPollId(pollID string) ([]Bet, error)
	UpdateBet(bet *Bet) error
}

// Errors related to bets

var ErrBetNotFound = errors.New("bet not found")
var ErrUserAlreadyBet = errors.New("user already bet")
