package bets

import (
	"database/sql"
	"errors"
	"fmt"
)

type libSQLRepository struct {
	db *sql.DB
}

func NewLibSQLRepository(db *sql.DB) BetRepository {
	return &libSQLRepository{
		db: db,
	}
}

func (repo libSQLRepository) Save(bet *Bet) error {
	query := `INSERT INTO bets (poll_id, user_id, selected_option_index, bet_status) VALUES (?, ?, ?, ?)`

	preparedStatement, preparedErr := repo.db.Prepare(query)
	if preparedErr != nil {
		return fmt.Errorf("error while preparing save bet statement: %w", preparedErr)
	}

	_, execErr := preparedStatement.Exec(bet.PollId, bet.UserId, bet.SelectedOptionIndex, bet.BetStatus)
	if execErr != nil {
		return fmt.Errorf("error while executing save bet statement: %w", execErr)
	}

	return nil
}

func (repo libSQLRepository) GetByPollIdAndUserId(pollID string, userID string) (*Bet, error) {
	query := "SELECT poll_id, user_id, selected_option_index, bet_status FROM bets WHERE poll_id = ? AND user_id = ?"
	preparedStatement, preparedErr := repo.db.Prepare(query)
	if preparedErr != nil {
		return nil, fmt.Errorf("error while preparing get bet by poll_id and user_id statement: %w", preparedErr)
	}

	var bet Bet
	row := preparedStatement.QueryRow(pollID, userID)
	if scanErr := row.Scan(&bet.PollId, &bet.UserId, &bet.SelectedOptionIndex, &bet.BetStatus); scanErr != nil {
		if errors.Is(scanErr, sql.ErrNoRows) {
			return nil, ErrBetNotFound
		}
		return nil, fmt.Errorf("error while scanning bet: %w", scanErr)
	}

	return &bet, nil
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
