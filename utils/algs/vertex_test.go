package graph

import (
	"github.com/clausthrane/futfut/models"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestToVertex(t *testing.T) {
	assert := assert.New(t)
	in := []models.TrainEvent{e1, e2, e3}

	lout := ToVertices(in)

	assert.Equal(6, len(lout), "Each train event has two timestamps -arrival and departure- both should be edges")

	/*
		logger.Println("inspecting list")
		for _, v := range lout {
			logger.Println(v.String())
		}

		logger.Println("inspecting map")
		for k, _ := range mout {
			logger.Println(k)
			value := mout[k]
			logger.Println(value)
		}
	*/
	some := lout[1]

	assert.Equal(MAX_DISTANCE, some.TimeFromSource(), "Should be initialized")
	some.SetTimeFromSource(42)
	assert.Equal(int64(42), some.TimeFromSource(), "Can be updated")

	/*
		mappedSome := mout[some.HashKey()]

		assert.Equal(mappedSome.HashKey(), some.HashKey())
		(*mappedSome).SetTimeFromSource(0)
		assert.Equal(int64(0), some.TimeFromSource())
	*/
}

func TestUpdateTimeFromSource(t *testing.T) {
	assert := assert.New(t)

	v := NewVertex(time.Now(), &e1)
	assert.Equal(MAX_DISTANCE, v.TimeFromSource())

	v.SetTimeFromSource(42)
	assert.Equal(int64(42), v.TimeFromSource())
}

func TestSetPrev(t *testing.T) {
	assert := assert.New(t)

	v := NewVertex(time.Now(), &e1)
	assert.Nil(v.GetPrev())

	u := NewVertex(time.Now(), &e2)
	v.SetPrev(u)

	assert.Equal(u, v.GetPrev())

	p := v
	p = p.GetPrev()
	assert.Equal(u, p)
}
