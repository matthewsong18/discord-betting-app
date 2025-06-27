package polls

import (
	"database/sql"
	"errors"
	"fmt"
)

type libSQLRepository struct {
	db *sql.DB
}

func NewLibSQLRepository(db *sql.DB) PollRepository {
	return &libSQLRepository{db: db}
}

func (repo *libSQLRepository) Save(poll *Poll) error {
	if err := saveToPollsTable(poll, repo); err != nil {
		return fmt.Errorf("save polls table failed: %w", err)
	}

	if err := saveToOptionsTable(poll, repo); err != nil {
		return fmt.Errorf("save options table failed: %w", err)
	}

	return nil
}

func saveToPollsTable(poll *Poll, repo *libSQLRepository) error {
	query := "INSERT INTO polls (id, title, status) VALUES (?, ?, ?)"
	preparedStatement, prepareError := repo.db.Prepare(query)
	if prepareError != nil {
		return fmt.Errorf("error while preparing statement: %w", prepareError)
	}
	if result, execErr := preparedStatement.Exec(poll.ID, poll.Title, poll.Status); execErr != nil {
		return fmt.Errorf("error while executing statement: %w", execErr)
	} else {
		rowsAffected, _ := result.RowsAffected()
		if rowsAffected == 0 {
			return fmt.Errorf("no rows were affected by the insert operation")
		}
	}
	return nil
}

func saveToOptionsTable(poll *Poll, repo *libSQLRepository) error {
	query := "INSERT INTO poll_options (poll_id, option_index, option_text) VALUES (?, ?, ?)"
	preparedStatement, prepareError := repo.db.Prepare(query)
	if prepareError != nil {
		return fmt.Errorf("error while preparing statement: %w", prepareError)
	}

	for index, option := range poll.Options {
		if result, execErr := preparedStatement.Exec(poll.ID, index, option); execErr != nil {
			return fmt.Errorf("error while executing statement for option %d: %w", index, execErr)
		} else {
			rowsAffected, _ := result.RowsAffected()
			if rowsAffected == 0 {
				return fmt.Errorf("no rows were affected by the insert operation for option %d", index)
			}
		}
	}

	return nil
}

func (repo *libSQLRepository) GetById(id string) (*Poll, error) {
	poll, pollErr := getFromPollTable(id, repo)
	if pollErr != nil {
		return nil, fmt.Errorf("error while getting poll from polls table: %w", pollErr)
	}

	var optionsErr error
	poll.Options, optionsErr = getFromOptionsTable(id, repo)
	if optionsErr != nil {
		return nil, fmt.Errorf("error while getting options from options table: %w", optionsErr)
	}

	return poll, nil
}

func getFromPollTable(id string, repo *libSQLRepository) (*Poll, error) {
	query := "SELECT id, title, status FROM polls WHERE id = ?"
	preparedStatement, err := repo.db.Prepare(query)
	if err != nil {
		return nil, fmt.Errorf("error while preparing statement: %w", err)
	}

	row := preparedStatement.QueryRow(id)
	poll := &Poll{}
	if err := row.Scan(&poll.ID, &poll.Title, &poll.Status); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("poll with id %s not found", id)
		}
		return nil, fmt.Errorf("error while scanning row: %w", err)
	}
	return poll, nil
}

func getFromOptionsTable(pollID string, repo *libSQLRepository) ([]string, error) {
	query := "SELECT option_text FROM poll_options WHERE poll_id = ? ORDER BY option_index"
	preparedStatement, preparedErr := repo.db.Prepare(query)
	if preparedErr != nil {
		return nil, fmt.Errorf("error while preparing statement: %w", preparedErr)
	}

	rows, rowErr := preparedStatement.Query(pollID)
	if rowErr != nil {
		return nil, fmt.Errorf("error while executing query: %w", rowErr)
	}

	var options []string
	for rows.Next() {
		var option string
		if err := rows.Scan(&option); err != nil {
			return nil, fmt.Errorf("error while scanning row: %w", err)
		}
		options = append(options, option)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error while iterating over rows: %w", err)
	}

	return options, nil
}

func (repo *libSQLRepository) Update(poll *Poll) error {
	//TODO implement me
	panic("implement me")
}

func (repo *libSQLRepository) Delete(pollID string) error {
	//TODO implement me
	panic("implement me")
}
