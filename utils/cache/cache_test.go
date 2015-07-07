package cachingfacade

import (
	"fmt"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/tests/mockdsbfacade"
	"github.com/pmylund/go-cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"testing"
	"time"
)

func TestCachingFacade(t *testing.T) {
	assert := assert.New(t)

	testKey := "key"
	testValue := "value"

	fail := make(chan error, 1)
	succ := make(chan *models.TrainEventList, 1)

	callCount := 0
	remoteAPIMock := new(mockfacade.MockDSB)
	remoteAPIMock.On("GetTrains", mock.Anything, mock.Anything).Return(succ, fail).Run(func(args mock.Arguments) {
		callCount += 1
	})

	facadeUnderTest := New(remoteAPIMock)

	succ1, fail1 := facadeUnderTest.GetTrains(testKey, testValue)

	// Pretend to get stuff from the wire
	succ <- &models.TrainEventList{}

	assertSuccess(t, succ1, fail1)

	succ2, fail2 := facadeUnderTest.GetTrains(testKey, testValue)
	assertSuccess(t, succ2, fail2)

	assert.Equal(1, callCount)
}

func assertSuccess(t *testing.T, success chan *models.TrainEventList, failure chan error) {
	select {
	case <-success:
		t.Log("Got success")
	case <-failure:
		t.Fail()
	case <-time.After(time.Second * 2):
		t.Fail()
	}
}

func TestCache(t *testing.T) {
	c := cache.New(2*time.Second, 1*time.Second)

	c.Set("foo", "bar", 0)

	if v, found := c.Get("foo"); found {
		fmt.Println(v)
	}

	time.Sleep(2 * time.Second)

	if _, found := c.Get("foo"); !found {
		fmt.Println("gone!")
	}
}
