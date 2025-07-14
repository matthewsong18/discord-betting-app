package users

import (
	"betting-discord-bot/internal/bets"
	"betting-discord-bot/internal/polls"
	"testing"
)

func TestCreateUser(t *testing.T) {
	pollMemoryRepo := polls.NewMemoryRepository()
	pollService := polls.NewService(pollMemoryRepo)
	betService := bets.NewService(pollService, nil)
	userRepo := NewMemoryRepository()
	userService := NewService(userRepo, betService)

	user, err := userService.CreateUser("12345")
	if err != nil {
		t.Fatalf("CreateUser returned an unexpected error: %v", err)
	}

	if user.DiscordID != "12345" {
		t.Errorf("Expected DiscordID to be '12345', got '%s'", user.DiscordID)
	}
}
