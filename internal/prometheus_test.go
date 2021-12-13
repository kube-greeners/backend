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

func TestParseIntervalValid(t *testing.T) {

	auxTestParseIntervalValid("2d", 2*24*time.Hour, t)
	auxTestParseIntervalValid("4h", 4*time.Hour, t)
	auxTestParseIntervalValid("1m", time.Minute, t)
	auxTestParseIntervalValid("10w", 10*24*7*time.Hour, t)
	auxTestParseIntervalValid("500s", 500*time.Second, t)
}

func TestParseIntervalFailing(t *testing.T) {
	auxTestParseIntervalFailing("0s", t)
	auxTestParseIntervalFailing("-2s", t)

	auxTestParseIntervalFailing("s", t)
	auxTestParseIntervalFailing("4.5h", t)
	auxTestParseIntervalFailing("s", t)
	auxTestParseIntervalFailing("", t)
	auxTestParseIntervalFailing("2p", t)
}
