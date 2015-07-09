// Package utils implements a few helper functions
package utils

import (
	"time"
)

// DATE_LAYOUT is the layout produced by the time.Time String()
// this constant is usefil for parsing a string back into a time.Time
const DATE_LAYOUT = "2006-01-02 15:04:05 -0700 MST"

// SubmitAsync starts a new go routine and submits the error on the given channel
func SubmitAsync(e error, c chan error) {
	go func() {
		c <- e
	}()
}

// ParseDateString will return a time.Time from a string which is the result of a
// time.Time String() call. I.e useful for roundtrip conversion
func ParseDateString(date string) time.Time {
	t, _ := time.Parse(DATE_LAYOUT, date)
	return t
}

func AddWithoutOverflow(a int64, b int64) int64 {
	sum := a + b
	if sum < a {
		return a
	} else {
		return sum
	}
}
