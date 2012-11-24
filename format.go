package nut

import (
	"time"
)

func parseString(input string) (string, bool) {
	return "", false
}

func parseInt(input string) (int64, bool) {
	return 0, false
}

func parseBool(input string) (bool, bool) {
	return false, false
}

func parseDuration(input string) (time.Duration, bool) {
	return time.Duration(0), false
}

func parseTime(input string) (*time.Time, bool) {
	return nil, false
}
