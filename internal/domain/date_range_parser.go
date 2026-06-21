package domain

import (
	"errors"
	"fmt"
	"math"
	"strings"
	"time"
)

// ErrInvalidDateFormat is returned by ParseDateRangeEpoch when one of the
// inputs does not match the yyyy-MM-dd shape. Use errors.Is to check it.
var ErrInvalidDateFormat = errors.New("invalid date format, expected yyyy-MM-dd")

// ParseDateRangeEpoch translates two optional date strings (yyyy-MM-dd, UTC)
// into a closed millisecond-epoch range [from 00:00, to 23:59:59.999].
//
// Returns nil when both inputs are blank — the caller should treat that as
// "no filter applied". Wraps ErrInvalidDateFormat on parse failures.
func ParseDateRangeEpoch(from, to string) (*Int64Range, error) {
	f := strings.TrimSpace(from)
	t := strings.TrimSpace(to)
	if f == "" && t == "" {
		return nil, nil
	}

	var fromEpoch int64 = 0
	if f != "" {
		d, err := time.Parse("2006-01-02", f)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidDateFormat, f)
		}
		fromEpoch = d.UTC().UnixMilli()
	}

	var toEpoch int64 = math.MaxInt64
	if t != "" {
		d, err := time.Parse("2006-01-02", t)
		if err != nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidDateFormat, t)
		}
		// End of the requested day: next-day midnight UTC minus 1ms.
		toEpoch = d.UTC().AddDate(0, 0, 1).UnixMilli() - 1
	}

	return &Int64Range{Min: fromEpoch, Max: toEpoch}, nil
}
