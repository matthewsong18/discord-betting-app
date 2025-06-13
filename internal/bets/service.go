package bets

import "errors"

type service struct{}

func NewService() BetService {
	return &service{}
}

func (s service) CreateBet(pollId string, userId int, selectedOptionIndex int) (*Bet, error) {
	if selectedOptionIndex < 0 || selectedOptionIndex > 2 {
		return nil, errors.New("invalid option index")
	}

	bet := &Bet{
		PollId:              pollId,
		UserId:              userId,
		SelectedOptionIndex: selectedOptionIndex,
	}

	return bet, nil
}
