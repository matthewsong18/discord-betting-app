package users

type user struct {
	ID        string
	DiscordID string
}

type User interface {
	GetID() string
	GetDiscordID() string
}

func (u *user) GetID() string        { return u.ID }
func (u *user) GetDiscordID() string { return u.DiscordID }

type WinLoss struct {
	Wins   int
	Losses int
}
