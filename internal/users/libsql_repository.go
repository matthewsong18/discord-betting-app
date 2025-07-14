package users

import (
	"database/sql"
	"errors"
	"fmt"
)

type libsqlRepository struct {
	db *sql.DB
}

func NewLibSQLRepository(db *sql.DB) UserRepository {
	return &libsqlRepository{db}
}

func (repo *libsqlRepository) Save(user *User) error {
	query := `INSERT INTO users (id, discord_id) VALUES (?, ?)`
	_, err := repo.db.Exec(query, user.ID, user.DiscordID)
	if err != nil {
		return fmt.Errorf("error saving user: %w", err)
	}

	return nil
}

func (repo *libsqlRepository) GetByID(id string) (*User, error) {
	query := `SELECT id, discord_id FROM users WHERE id = ?`
	row := repo.db.QueryRow(query, id)

	var user User
	err := row.Scan(&user.ID, &user.DiscordID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	return &user, nil
}

func (repo *libsqlRepository) GetByDiscordID(discordID string) (*User, error) {
	query := `SELECT id, discord_id FROM users WHERE discord_id = ?`
	row := repo.db.QueryRow(query, discordID)

	var user User
	err := row.Scan(&user.ID, &user.DiscordID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error retrieving user: %w", err)
	}
	return &user, nil
}

func (repo *libsqlRepository) Delete(discordID string) error {
	query := "DELETE FROM users WHERE discord_id = ?"
	result, err := repo.db.Exec(query, discordID)
	if err != nil {
		return fmt.Errorf("error deleting user with discord_id %s: %w", discordID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("error checking rows affected for discord_id %s: %w", discordID, err)
	}

	if rowsAffected == 0 {
		return ErrUserNotFound
	}

	return nil
}

var _ UserRepository = (*libsqlRepository)(nil)
