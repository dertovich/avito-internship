package postgres

import (
	"avito-internship/internal/storage"
	"fmt"
)

func (p *Postgres) SegmentExists(segment string) (bool, error) {
	const op = "storage.postgres.segments_table.SegmentExists"

	var res bool
	err := p.segmentsTable.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM segments WHERE name = $1)", segment).Scan(&res)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (p *Postgres) UserExists(user_id int64) (bool, error) {
	const op = "storage.postgres.users_table.UserExists"

	var res bool
	err := p.usersTable.QueryRow(
		"SELECT EXISTS (SELECT 1 FROM users WHERE id = $1)", user_id).Scan(&res)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return res, nil
}

func (p *Postgres) validateUserAndSegment(user_id int64, segments []string) error {
	userExists, err := p.UserExists(user_id)
	if err != nil {
		return err
	}
	if !userExists {
		return storage.ErrUserNotFound
	}

	for _, segment := range segments {
		segmentExists, err := p.SegmentExists(segment)
		if err != nil {
			return err
		}
		if !segmentExists {
			return storage.ErrSegmentNotFound
		}
	}

	return nil
}

func (p *Postgres) validateSegments(segments []string) ([]string, error) {
	// Check that the list of segments is not empty.
	if len(segments) == 0 {
		return nil, fmt.Errorf("segments must not be empty")
	}

	// Check that the list of segments does not contain any duplicates.
	uniqueSegments := make(map[string]bool)
	for _, segment := range segments {
		if uniqueSegments[segment] {
			return nil, fmt.Errorf("segments must not contain any duplicates")
		}
		uniqueSegments[segment] = true
	}

	return segments, nil
}
