package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

func NewUsersTable(db *sql.DB) (*sql.DB, error) {
	const op = "storage.postgres.NewUsersTable"

	tx, err := db.Begin()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := tx.Prepare(`
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		segments pq.StringArray
	);
	`)
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return db, nil
}

func (p *Postgres) CreateUser(user_id int64, segments []string) error {
	const op = "storage.postgres.users_table.CreateUser"

	segments, err := p.validateSegments(segments)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	stmt, err := p.usersTable.Prepare("INSERT INTO users(id, segments) VALUES(?,?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(user_id, pq.StringArray(segments))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (p *Postgres) AddUserToSegment(user_id int64, segments []string) error {
	const op = "storage.postgres.users_table.AddUserToSegment"

	if err := p.validateUserAndSegment(user_id, segments); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, segment := range segments {
		_, err := p.usersTable.Exec(
			"UPDATE users SET segments = pq.ArrayAppend(segments, $1) WHERE id = $2",
			segment,
			user_id,
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (p *Postgres) RemoveSegmentsFromUser(user_id int64, segments []string) error {
	const op = "storage.postgres.users_table.RemoveSegmentsFromUser"

	if err := p.validateUserAndSegment(user_id, segments); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, segment := range segments {
		_, err := p.usersTable.Exec(
			"UPDATE users SET segments = pq.ArrayDelete(segments, $1) WHERE id = $2",
			segment,
			user_id,
		)
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
	}

	return nil
}

func (p *Postgres) ShowActiveSegmentUser(user_id int64) ([]string, error) {
	const op = "storage.postgres.users_table.ShowActiveSegmentUser"

	if _, err := p.UserExists(user_id); err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var segments pq.StringArray
	err := p.usersTable.QueryRow(
		"SELECT segments FROM users WHERE id = $1", user_id).Scan(&segments)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	var activeSegments []string
	for _, segment := range segments {
		activeSegments = append(activeSegments, segment)
	}

	return activeSegments, nil
}
