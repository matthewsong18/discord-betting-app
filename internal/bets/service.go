package bets

import (
	"betting-discord-bot/internal/polls"
	"errors"
)

type service struct {
	pollService polls.PollService
	betList     []Bet
}

func NewService(pollService polls.PollService) BetService {
	return &service{
		pollService: pollService,
		betList:     make([]Bet, 0),
	}
}

func (s *service) CreateBet(pollId string, userId int, selectedOptionIndex int) (*Bet, error) {
	if selectedOptionIndex < 0 || selectedOptionIndex > 2 {
		return nil, errors.New("invalid option index")
	}

	poll, err := s.pollService.GetPollById(pollId)
	if err != nil {
		return nil, err
	}

	if poll.IsOpen == false {
		return nil, errors.New("cannot bet on a closed poll")
	}

	err = checkIfUserAlreadyBetOnPoll(pollId, userId, s)

	if err != nil {
		return nil, err
	}

	bet := &Bet{
		PollId:              pollId,
		UserId:              userId,
		SelectedOptionIndex: selectedOptionIndex,
	}

	// Add the bet to the poll's bet list
	s.betList = append(s.betList, *bet)

	return bet, nil
}

func checkIfUserAlreadyBetOnPoll(pollId string, userId int, s *service) error {
	for _, bet := range s.betList {
		// Skip if the bet is not for the specified poll or user
		if bet.PollId != pollId || bet.UserId != userId {
			continue
		}

		return errors.New("user bet already exists for this poll")
	}
	return nil
}
