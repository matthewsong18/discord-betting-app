package users

import (
	"betting-discord-bot/internal/bets"
	"betting-discord-bot/internal/polls"
	"errors"
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

func TestGetUserByDiscordID(t *testing.T) {
	userRepo := NewMemoryRepository()
	userService := NewService(userRepo, nil)

	user, err := userService.CreateUser("12345")
	if err != nil {
		t.Fatalf("CreateUser returned an unexpected error: %v", err)
	}

	retrievedUser, err := userService.GetUserByDiscordID("12345")
	if err != nil {
		t.Fatalf("GetUserByDiscordID returned an unexpected error: %v", err)
	}

	if retrievedUser.ID != user.ID {
		t.Errorf("Expected ID to be '%s', got '%s'", user.ID, retrievedUser.ID)
	}
}

func TestDeleteUser(t *testing.T) {
	userRepo := NewMemoryRepository()
	userService := NewService(userRepo, nil)

	user, err := userService.CreateUser("12345")
	if err != nil {
		t.Fatalf("CreateUser returned an unexpected error: %v", err)
	}

	err = userService.DeleteUser(user.DiscordID)
	if err != nil {
		t.Fatalf("DeleteUser returned an unexpected error: %v", err)
	}

	_, err = userService.GetUserByDiscordID(user.DiscordID)
	if err == nil {
		t.Fatalf("Expected GetUserByDiscordID to return an error after deletion")
	}

	if errors.Is(err, ErrUserNotFound) {
		t.Fatalf("Expected GetUserByDiscordID to return ErrUserNotFound after deletion")
	}
}
