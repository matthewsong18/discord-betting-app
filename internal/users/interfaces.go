package users

import "errors"

type UserService interface {
	CreateUser(discordID string) (*User, error)
	GetUserByDiscordID(discordID string) (*User, error)
	DeleteUser(discordID string) error
	GetWinLoss(userID string) (*WinLoss, error)
}

type UserRepository interface {
	Save(user *User) error
	GetByID(id string) (*User, error)
	GetByDiscordID(discordID string) (*User, error)
	Delete(discordID string) error
}

var ErrUserNotFound = errors.New("user not found")
