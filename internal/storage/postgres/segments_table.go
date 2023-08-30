package postgres

import (
	"avito-internship/internal/storage"
	"database/sql"
	"fmt"

	"github.com/lib/pq"
)

func NewSegmentsTable(db *sql.DB) (*sql.DB, error) {
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

	return db, nil
}

func (p *Postgres) CreateSegment(segmentToCreate string) (int64, error) {
	const op = "storage.postgres.segments_table.CreateSegment"

	tx, err := p.segmentsTable.Begin()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to begin transaction: %w", op, err)
	}

	stmt, err := p.segmentsTable.Prepare("INSERT INTO segments(name) VALUES(DEFAULT, $1)")
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.Exec(segmentToCreate)
	if err != nil {
		tx.Rollback()
		pqErr, ok := err.(*pq.Error)
		if ok && pqErr.Code == "23505" {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrSegmentExists)
		}
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		tx.Rollback()
		return 0, fmt.Errorf("%s: failed to get last insert id: %w", op, err)
	}

	err = tx.Commit()
	if err != nil {
		return 0, fmt.Errorf("%s: failed to commit transaction: %w", op, err)
	}

	return id, nil
}

func (p *Postgres) DeleteSegment(segmentToDelete string) (int64, error) {
	const op = "storage.postgres.segments_table.DeleteSegment"

	stmt, err := p.segmentsTable.Prepare("DELETE FROM segments WHERE name = $1")
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
