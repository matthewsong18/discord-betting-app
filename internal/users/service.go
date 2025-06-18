package users

import (
	"betting-discord-bot/internal/bets"
	"fmt"
	"github.com/google/uuid"
)

type service struct {
	betService bets.BetService
}

func NewService(betService bets.BetService) UserService {
	return &service{
		betService: betService,
	}
}

func (service service) CreateUser(discordID string) (*User, error) {
	user := &User{
		ID:        uuid.NewString(),
		DiscordID: discordID,
	}
	return user, nil
}

func (service service) GetWinLoss(userID string) (*WinLoss, error) {
	winLoss := &WinLoss{
		Wins:   0,
		Losses: 0,
	}

	betList, betListErr := service.betService.GetBetsFromUser(userID)
	if betListErr != nil {
		return nil, fmt.Errorf("failed to get bets for user %s: %w", userID, betListErr)
	}

	for _, bet := range betList {
		switch bet.BetStatus {
		case bets.StatusWon:
			winLoss.Wins++
		case bets.StatusLost:
			winLoss.Losses++
		case bets.StatusPending:
			continue
		}
	}

	return winLoss, nil
}

var _ UserService = (*service)(nil)
