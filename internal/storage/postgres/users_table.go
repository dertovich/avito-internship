package postgres

import (
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type UsersTable struct {
	db *sql.DB
}

func NewUsersTable(db *sql.DB) (*UsersTable, error) {
	const op = "storage.postgres.NewUsersTable"

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS users(
		id SERIAL PRIMARY KEY,
		segments TEXT[]
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &UsersTable{db: db}, nil
}

func (p *Postgres) CreateUser(user_id int, segments []string) error {
	const op = "storage.postgres.users_table.CreateUser"

	stmt, err := p.usersTable.db.Prepare("INSERT INTO users(id, segments) VALUES(?,?)")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec(user_id, pq.Array(segments))
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

	stmt, err := p.usersTable.db.Prepare(
		"UPDATE users SET segments = array_append(segments, $1) WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, segment := range segments {
		_, err = stmt.Exec(segment, user_id)
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

	stmt, err := p.usersTable.db.Prepare(
		"UPDATE users SET segments = array_remove(segments, $1) WHERE id = $2")
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for _, segment := range segments {
		_, err = stmt.Exec(segment, user_id)
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

	var segments []string
	err := p.usersTable.db.QueryRow(
		"SELECT segments FROM users WHERE id = $1", user_id).Scan(pq.Array(&segments))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return segments, nil
}
