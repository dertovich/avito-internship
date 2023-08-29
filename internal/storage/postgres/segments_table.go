package postgres

import (
	"avito-internship/internal/storage"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

type SegmentsTable struct {
	db *sql.DB
}

func NewSegmentsTable(db *sql.DB) (*SegmentsTable, error) {
	const op = "storage.postgres.NewSegmentsTable"

	stmt, err := db.Prepare(`
	CREATE TABLE IF NOT EXISTS segments(
		id SERIAL PRIMARY KEY,
		name TEXT NOT NULL UNIQUE
	);
	`)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	_, err = stmt.Exec()
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	return &SegmentsTable{db: db}, nil
}

func (p *Postgres) CreateSegment(segmentToCreate string) (int64, error) {
	const op = "storage.postgres.segments_table.CreateSegment"

	stmt, err := p.segmentsTable.db.Prepare("INSERT INTO segments(id, name) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(segmentToCreate)
	if err != nil {
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrSegmentExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}
	return id, nil
}

func (p *Postgres) DeleteSegment(segmentToDelete string) (int64, error) {
	const op = "storage.postgres.segments_table.DeleteSegment"

	stmt, err := p.segmentsTable.db.Prepare("DELETE FROM segments WHERE name = $1")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(segmentToDelete)
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return rowsAffected, nil
}
