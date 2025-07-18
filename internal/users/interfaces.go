package users

type UserService interface {
	CreateUser(discordID string) (*User, error)
	GetWinLoss(userID string) (*WinLoss, error)
}

type UserRepository interface {
	Save(user *User) error
	GetByID(id string) (*User, error)
	GetByDiscordID(discordID string) (*User, error)
	Delete(discordID string) error
}
