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

func (betService *service) CreateBet(pollID string, userID string, selectedOptionIndex int) (*Bet, error) {
	if selectedOptionIndex < 0 || selectedOptionIndex > 2 {
		return nil, errors.New("invalid option index")
	}

	poll, err := betService.pollService.GetPollById(pollID)
	if err != nil {
		return nil, err
	}

	if poll.Status == polls.Closed {
		return nil, errors.New("cannot bet on a closed poll")
	}

	err = checkIfUserAlreadyBetOnPoll(pollID, userID, betService)

	if err != nil {
		return nil, err
	}

	bet := &Bet{
		PollId:              pollID,
		UserId:              userID,
		SelectedOptionIndex: selectedOptionIndex,
		BetStatus:           Pending,
	}

	// Add the bet to the poll's bet list
	betService.betList = append(betService.betList, *bet)

	return bet, nil
}

func checkIfUserAlreadyBetOnPoll(pollId string, userId string, s *service) error {
	for _, bet := range s.betList {
		// Skip if the bet is not for the specified poll or user
		if bet.PollId != pollId || bet.UserId != userId {
			continue
		}

		return errors.New("user bet already exists for this poll")
	}
	return nil
}

func (betService *service) GetBet(pollID string, userID string) (*Bet, error) {
	for _, bet := range betService.betList {
		if bet.PollId == pollID && bet.UserId == userID {
			return &bet, nil
		}
	}
	return nil, errors.New("bet not found for the specified poll and user")
}

func (betService *service) UpdateBetsByPollId(poll polls.Poll) {
	for i, bet := range betService.betList {
		if bet.PollId == poll.ID {
			if bet.SelectedOptionIndex == poll.Outcome {
				betService.betList[i].BetStatus = Won
			} else {
				betService.betList[i].BetStatus = Lost
			}
		}
	}
}

func (betService *service) GetBetsFromUser(userID string) ([]Bet, error) {
	var userBets []Bet

	for _, bet := range betService.betList {
		if bet.UserId == userID {
			userBets = append(userBets, bet)
		}
	}

	return userBets, nil
}

var _ BetService = (*service)(nil)
