package domain

import (
	"errors"
	"testing"
)

func TestParseDateRangeEpoch_BothNilOrBlank(t *testing.T) {
	cases := []struct{ from, to string }{
		{"", ""},
		{"   ", ""},
		{"", "  "},
		{"   ", "   "},
	}
	for _, c := range cases {
		got, err := ParseDateRangeEpoch(c.from, c.to)
		if err != nil {
			t.Fatalf("ParseDateRangeEpoch(%q,%q) returned error: %v", c.from, c.to, err)
		}
		if got != nil {
			t.Fatalf("ParseDateRangeEpoch(%q,%q) = %+v, want nil", c.from, c.to, got)
		}
	}
}

func TestParseDateRangeEpoch_OnlyFrom(t *testing.T) {
	// 2024-01-15 00:00:00 UTC = 1705276800000 ms
	const expectedFrom int64 = 1705276800000
	got, err := ParseDateRangeEpoch("2024-01-15", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil range")
	}
	if got.Min != expectedFrom {
		t.Errorf("Min = %d, want %d", got.Min, expectedFrom)
	}
	const wantMax int64 = 1<<63 - 1
	if got.Max != wantMax {
		t.Errorf("Max = %d, want MaxInt64 (%d)", got.Max, wantMax)
	}
}

func TestParseDateRangeEpoch_OnlyTo(t *testing.T) {
	// 2024-01-15 -> end of day = 2024-01-16 00:00:00 UTC - 1ms = 1705363199999
	const expectedTo int64 = 1705363199999
	got, err := ParseDateRangeEpoch("", "2024-01-15")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil range")
	}
	if got.Min != 0 {
		t.Errorf("Min = %d, want 0", got.Min)
	}
	if got.Max != expectedTo {
		t.Errorf("Max = %d, want %d", got.Max, expectedTo)
	}
}

func TestParseDateRangeEpoch_Both(t *testing.T) {
	got, err := ParseDateRangeEpoch("2024-01-01", "2024-01-31")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if got == nil {
		t.Fatal("expected non-nil range")
	}
	const expectedFrom int64 = 1704067200000  // 2024-01-01 00:00:00 UTC
	const expectedTo int64 = 1706745599999    // 2024-01-31 23:59:59.999 UTC
	if got.Min != expectedFrom {
		t.Errorf("Min = %d, want %d", got.Min, expectedFrom)
	}
	if got.Max != expectedTo {
		t.Errorf("Max = %d, want %d", got.Max, expectedTo)
	}
}

func TestParseDateRangeEpoch_InvalidFormat(t *testing.T) {
	_, err := ParseDateRangeEpoch("15-01-2024", "")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidDateFormat) {
		t.Errorf("error %v is not ErrInvalidDateFormat", err)
	}
}

func TestParseDateRangeEpoch_InvalidTo(t *testing.T) {
	_, err := ParseDateRangeEpoch("2024-01-01", "not-a-date")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, ErrInvalidDateFormat) {
		t.Errorf("error %v is not ErrInvalidDateFormat", err)
	}
}
