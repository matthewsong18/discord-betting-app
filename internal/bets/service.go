package bets

import (
	"errors"
	"fmt"

	"betting-discord-bot/internal/polls"
)

type service struct {
	pollService polls.PollService
	betRepo     BetRepository
}

func NewService(pollService polls.PollService, betRepo BetRepository) BetService {
	return &service{
		pollService: pollService,
		betRepo:     betRepo,
	}
}

func (betService *service) CreateBet(pollID string, userID string, selectedOptionIndex int) (Bet, error) {
	if selectedOptionIndex < 0 || selectedOptionIndex > 2 {
		return nil, errors.New("invalid option index")
	}

	poll, err := betService.pollService.GetPollById(pollID)
	if err != nil {
		return nil, err
	}

	if poll.GetStatus() == polls.Closed {
		return nil, ErrPollIsClosed
	}

	if err := checkIfUserAlreadyBetOnPoll(pollID, userID, betService); err != nil {
		return nil, err
	}

	bet := &bet{
		PollID:              pollID,
		UserID:              userID,
		SelectedOptionIndex: selectedOptionIndex,
		BetStatus:           Pending,
	}

	if err := betService.betRepo.Save(bet); err != nil {
		return nil, fmt.Errorf("failed to save bet: %w", err)
	}

	return bet, nil
}

func checkIfUserAlreadyBetOnPoll(pollId string, userId string, s *service) error {
	bet, err := s.betRepo.GetByPollIdAndUserId(pollId, userId)
	if err != nil {
		return nil
	}
	if bet != nil {
		return ErrUserAlreadyBet
	}

	return nil
}

func (betService *service) GetBet(pollID string, userID string) (Bet, error) {
	if bet, err := betService.betRepo.GetByPollIdAndUserId(pollID, userID); err != nil {
		return nil, fmt.Errorf("failed to get bet: %w", err)
	} else {
		return bet, nil
	}
}

func (betService *service) UpdateBetsByPollId(pollID string) error {
	poll, err := betService.pollService.GetPollById(pollID)
	if err != nil {
		return fmt.Errorf("failed to get poll by ID: %w", err)
	}

	betList, err := betService.betRepo.GetBetsByPollId(pollID)
	if err != nil {
		return fmt.Errorf("failed to get bets for poll: %w", err)
	}

	pollResult := int(poll.GetOutcome())
	for _, bet := range betList {
		if bet.SelectedOptionIndex == pollResult {
			bet.BetStatus = Won
		} else {
			bet.BetStatus = Lost
		}

		if err := betService.betRepo.UpdateBet(bet); err != nil {
			return fmt.Errorf("failed to update bet: %w", err)
		}
	}

	return nil
}

func (betService *service) GetBetsFromUser(userID string) ([]Bet, error) {
	userBets, err := betService.betRepo.GetBetsFromUser(userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get bets for user: %w", err)
	}

	var bets []Bet
	for _, bet := range userBets {
		bets = append(bets, bet)
	}

	return bets, nil
}

var _ BetService = (*service)(nil)
