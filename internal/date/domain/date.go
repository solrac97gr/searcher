package domain

import "time"

// DateFormatter exposes methods for date manipulation.
type DateFormatter interface {
	FromISO8601String(string) (*time.Time, error)
}
