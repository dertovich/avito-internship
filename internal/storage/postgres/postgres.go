package postgres

import (
	"database/sql"
	"fmt"
)

type Postgres struct {
	segmentsTable *sql.DB
	usersTable    *sql.DB
}

func New(postgresPath string) (*Postgres, error) {
	const op = "storage.postgres.New"

	db, err := sql.Open("postgres", postgresPath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	segmentsTable, err := NewSegmentsTable(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	usersTable, err := NewUsersTable(db)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &Postgres{segmentsTable: segmentsTable, usersTable: usersTable}, nil
}
