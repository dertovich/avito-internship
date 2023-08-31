package storage

import (
	"errors"
)

var (
	ErrUserNotFound    = errors.New("User not found")
	ErrSegmentExists   = errors.New("Segment is exists")
	ErrSegmentNotFound = errors.New("Segment not found")
)
