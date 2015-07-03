package dto

import (
	"github.com/clausthrane/futfut/models"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestConvertStation(t *testing.T) {
	assert := assert.New(t)

	in := &models.Station{"a", "b", "c", "d", "e"}
	out := NewStationConverter().ConvertStation(in)

	assert.Equal("b", out.Name, "Name should be copied")
	assert.Equal("e", out.CountryCode, "CountryCode is mapped from CountryName")
	assert.Equal("c", out.Id, "Id should be the UIC")
}

func TestConvertStationList(t *testing.T) {
	assert := assert.New(t)
	list := []models.Station{models.Station{"a", "b", "c", "d", "e"},
		models.Station{"r", "s", "t", "u", "v"}}

	in := &models.StationList{list}
	out := NewStationConverter().ConvertStationList(in)

	assert.Equal(2, out.Count)
	assert.Equal("s", out.Stations[1].Name, "Name should be copied")
}
