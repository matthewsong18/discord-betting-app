package users

type UserService interface {
	CreateUser(discordID string) (*User, error)
	GetWinLoss(userID string) (*WinLoss, error)
}
