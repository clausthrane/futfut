package dto

import (
	"github.com/clausthrane/futfut/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

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
		"",
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
	assert.Equal(8600761, out.CurrentStationId, "should have been mapped")
	assert.Equal("Hundige", out.DestinationName, "should have been mapped")
}
