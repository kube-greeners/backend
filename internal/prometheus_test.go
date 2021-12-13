package internal

import (
	"testing"
	"time"
)

func auxTestParseIntervalValid(functionParam string, expectedValue time.Duration, t *testing.T) {
	output, err := parseInterval(functionParam)
	if err != nil {
		t.Error("Unexpected error for param " + functionParam)
	}
	if output != expectedValue {
		t.Error("Unexpected value for param " + functionParam)
	}
}

func auxTestParseIntervalFailing(functionParam string, t *testing.T) {
	output, err := parseInterval(functionParam)
	if err == nil {
		t.Error("Error expected for param " + functionParam)
	}
	if output != 0 {
		t.Error("Error expected for param " + functionParam + ", output should be 0")
	}

}

func TestParseIntervalValidDay(t *testing.T) {
	auxTestParseIntervalValid("2d", 2*24*time.Hour, t)
}

func TestParseIntervalValidHour(t *testing.T) {
	auxTestParseIntervalValid("4h", 4*time.Hour, t)
}

func TestParseIntervalValidMinute(t *testing.T) {
	auxTestParseIntervalValid("1m", time.Minute, t)
}

func TestParseIntervalValidWeek(t *testing.T) {
	auxTestParseIntervalValid("10w", 10*24*7*time.Hour, t)
}

func TestParseIntervalValidSecond(t *testing.T) {
	auxTestParseIntervalValid("500s", 500*time.Second, t)
}

func TestParseIntervalFailingNullValue(t *testing.T) {
	auxTestParseIntervalFailing("0s", t)
}
func TestParseIntervalFailingNegativeValue(t *testing.T) {
	auxTestParseIntervalFailing("-2s", t)
}

func TestParseIntervalFailingMissingValue(t *testing.T) {
	auxTestParseIntervalFailing("s", t)
}

func TestParseIntervalFailingFloatValue(t *testing.T) {
	auxTestParseIntervalFailing("4.5h", t)
}

func TestParseIntervalFailingEmpty(t *testing.T) {
	auxTestParseIntervalFailing("", t)
}

func TestParseIntervalFailingWrongFactor(t *testing.T) {
	auxTestParseIntervalFailing("2p", t)
}
