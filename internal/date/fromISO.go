package date

import (
	"fmt"
	"time"
)

// FromISO8601String takes a date string in ISO8601 format and returns a
// *time.Time in UTC.
func (f Formatter) FromISO8601String(s string) (*time.Time, error) {
	parsedDate, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return nil, fmt.Errorf("validation_error: %s", err.Error())
	}

	parsedDate = parsedDate.UTC()

	return &parsedDate, nil
}
