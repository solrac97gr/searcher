package date

import (
	"github.com/solrac97gr/searcher/internal/date/domain"
)

// Formatter implements domain.DateFormatter.
type Formatter struct{}

// NewFormatter returns a new *date.Formatter implementing domain.DateFormatter.
func NewFormatter() (domain.DateFormatter, error) {
	return &Formatter{}, nil
}
