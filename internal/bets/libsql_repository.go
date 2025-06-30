package bets

import "database/sql"

type libSQLRepository struct {
}

func NewLibSQLRepository(db *sql.DB) BetRepository {
	return &libSQLRepository{}
}

func (repo libSQLRepository) Save(bet *Bet) error {
	//TODO implement me
	panic("implement me")
}

func (repo libSQLRepository) GetByPollIdAndUserId(pollID string, userID string) (*Bet, error) {
	//TODO implement me
	panic("implement me")
}

func (repo libSQLRepository) GetBetsFromUser(userID string) ([]Bet, error) {
	//TODO implement me
	panic("implement me")
}

func (repo libSQLRepository) GetBetsByPollId(pollID string) ([]Bet, error) {
	//TODO implement me
	panic("implement me")
}

func (repo libSQLRepository) UpdateBet(bet *Bet) error {
	//TODO implement me
	panic("implement me")
}
