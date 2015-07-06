package models

import (
	"github.com/clausthrane/futfut/utils"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHumanReadableDateConversion(t *testing.T) {
	assert := assert.New(t)
	in := "/Date(1436068203504)/"
	out := humanReadableDate(in)
	assert.Equal("2015-07-05 03:50:03 +0000 UTC", out.UTC().String())
}

func TestHumanReadableArrivalDate(t *testing.T) {
	assert := assert.New(t)

	train := TrainEvent{"event_id1",
		"station_id2",
		"type",
		"end_dest",
		123,
		"track1",
		"",
		"Train2",
		"86",
		"/Date(1436068203504)/",
		"",
		"",
		"",
		false,
		"",
		"",
		"",
	}

	out := train.HumanReadableArrivalDate()
	assert.Equal("2015-07-05 03:50:03 +0000 UTC", utils.ParseDateString(out).UTC().String())
}
