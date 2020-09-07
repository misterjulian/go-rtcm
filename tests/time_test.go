package tests

import (
	"testing"
	"time"

	"github.com/misterjulian/go-rtcm/timestamp"
)

func TestDF034(t *testing.T) {
	beforeUtcZero := time.Date(2019, 2, 7, 23, 41, 40, 0, time.UTC)
	if beforeUtcZero != timestamp.DF034(9700000, beforeUtcZero) {
		t.Errorf("DF034 timestamp incorrect before UTC 0: %v", timestamp.DF034(9700000, beforeUtcZero))
	}

	afterUtcZero := time.Date(2019, 2, 8, 00, 44, 44, 0, time.UTC)
	if afterUtcZero != timestamp.DF034(13484000, afterUtcZero) {
		t.Errorf("DF034 timestamp incorrect after UTC 0: %v", timestamp.DF034(13484000, afterUtcZero))
	}
}

func TestDF386(t *testing.T) {
	beforeUtcZero := time.Date(2019, 2, 7, 23, 41, 40, 0, time.UTC)
	if beforeUtcZero != timestamp.DF386(9700, beforeUtcZero) {
		t.Errorf("DF386 timestamp incorrect before UTC 0: %v", timestamp.DF386(9700, beforeUtcZero))
	}

	afterUtcZero := time.Date(2019, 2, 8, 00, 44, 44, 0, time.UTC)
	if afterUtcZero != timestamp.DF386(13484, afterUtcZero) {
		t.Errorf("DF386 timestamp incorrect after UTC 0: %v", timestamp.DF386(13484, afterUtcZero))
	}
}
