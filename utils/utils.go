package utils

import (
	"time"
)

const DATE_LAYOUT = "2006-01-02 15:04:05 -0700 MST"

func SubmitAsync(e error, c chan error) {
	go func() {
		c <- e
	}()
}

func ParseDateString(date string) time.Time {
	t, _ := time.Parse(DATE_LAYOUT, date)
	return t
}
