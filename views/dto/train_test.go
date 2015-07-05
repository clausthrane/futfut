package dto

import (
	"github.com/clausthrane/futfut/models"
	"github.com/stretchr/testify/assert"
	"regexp"
	"strconv"
	"testing"
	"time"
)

func TestUnixtime(t *testing.T) {
	assert := assert.New(t)
	ms := 1436023620000 / 1000
	ut := time.Unix(int64(ms), 0)
	assert.Equal("2015-07-04 15:27:00 +0000 UTC", ut.UTC().String())
}

func TestTimeRegex(t *testing.T) {
	assert := assert.New(t)
	r, _ := regexp.Compile("([0-9]+)")
	i, _ := strconv.Atoi(r.FindString("/Date(1436023620000)/"))
	assert.Equal(1436023620000, i)

}

func TestConvertTrain(t *testing.T) {
	assert := assert.New(t)

	in := &models.Train{
		"fc63491b-8b93-4d50-9482-39cdc31f5024",
		"8600761",
		"S-tog",
		"Hundige",
		8600769,
		"2",
		"/Date(1435936635503)/",
		"12246",
		"86",
		"",
		"/Date(1435936635503)/",
		"",
		"",
		false,
		"A",
		"Syd",
		"12",
	}

	out, err := NewTrainConverter().ConvertTrain(in)
	assert.Nil(err, "no errors expected")

	assert.Equal(12246, out.TrainNumber, "should have been mapped")
	assert.Equal(8600761, out.StationId, "should have been mapped")
	assert.Equal("Hundige", out.DestinationName, "should have been mapped")
}
