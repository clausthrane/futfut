package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestParseDateString(t *testing.T) {
	assert := assert.New(t)
	date := "2015-07-05 03:50:03 +0000 UTC"
	assert.Equal(date, ParseDateString(date).UTC().String())
}
