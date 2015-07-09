// Package cachingfacade implements a dsb.DSBFacade proxy which caches results for a given period
package cachingfacade

import (
	"fmt"
	"github.com/clausthrane/futfut/datasources/dsb"
	"github.com/clausthrane/futfut/models"
	"github.com/clausthrane/futfut/services"
	"github.com/pmylund/go-cache"
	"log"
	"os"
	"time"
)

const defaultCacheTime = 0

var logger = log.New(os.Stdout, " ", log.Ldate|log.Lmicroseconds|log.Lshortfile)

// A hard cache for all stations
var cachedStationList *models.StationList
var stationByIDMap map[string]*models.StationList

func init() {
	stationByIDMap = make(map[string]*models.StationList)
}

type cachingfacade struct {
	localcache *cache.Cache
	facade     dsb.DSBFacade
}

// New returns a new caching dsb.DSBFacade. All requests to the dsb.DSBFacade interface are
// proxied to the given 'facade' constructor argument
func New(facade dsb.DSBFacade) dsb.DSBFacade {
	c := cache.New(10*time.Minute, 20*time.Second)
	return &cachingfacade{c, facade}
}

func (c *cachingfacade) GetStation(stationid string) (chan *models.StationList, chan error) {
	success, failure := make(chan *models.StationList), make(chan error)
	go func() {
		if stationByIDMap[stationid] == nil {
			downstreamSucces, downstreamFailure := c.facade.GetStation(stationid)
			select {
			case result := <-downstreamSucces:
				stationByIDMap[stationid] = result
			case err := <-downstreamFailure:
				failure <- err
			case <-time.After(time.Second * 30):
				failure <- services.NewServiceTimeoutError("Failed to get all trains in time")
			}
		}
		success <- stationByIDMap[stationid]
	}()
	return success, failure
}

func (c *cachingfacade) GetStations() (chan *models.StationList, chan error) {
	//logger.Println("Intercepting request for GetStatins()")
	success, failure := make(chan *models.StationList), make(chan error)
	go func() {
		if cachedStationList == nil {
			logger.Println("Sending request downstream")
			downstreamSuccess, downstreamErr := c.facade.GetStations()
			select {
			case cachedStationList = <-downstreamSuccess:
			case err := <-downstreamErr:
				failure <- err
			case <-time.After(time.Second * 30):
				failure <- services.NewServiceTimeoutError("Failed to get all trains in time")
			}
		}
		success <- cachedStationList
	}()
	return success, failure
}

func (c *cachingfacade) GetAllTrains() (chan *models.TrainEventList, chan error) {
	return c.facade.GetAllTrains()
}

type trainCacheEntry struct {
	list *models.TrainEventList
	err  error
}

// GetTrains searches the embedded cache for a result otherwise deletates the request to the underlying
// facade.
func (c *cachingfacade) GetTrains(key string, value string) (chan *models.TrainEventList, chan error) {
	logger.Printf("Intercepting request for GetTrains(%s, %s)", key, value)
	cacheKey := toCacheKey("GetTrains", key, value)
	success, failure := make(chan *models.TrainEventList), make(chan error)
	go func() {
		if cachedItem, found := c.localcache.Get(cacheKey); found {
			v := cachedItem.(*trainCacheEntry)
			if v.err != nil {
				failure <- v.err
			} else {
				success <- v.list
			}
		} else {
			privSuccess, privFailure := c.facade.GetTrains(key, value)
			select {
			case result := <-privSuccess:
				success <- c.cacheSuccess(cacheKey, result)
			case err := <-privFailure:
				failure <- c.cacheFalure(cacheKey, err)
			case <-time.After(time.Second * 30):
				failure <- services.NewServiceTimeoutError("Failed to get all trains in time")
			}
		}
	}()
	return success, failure
}

func (c *cachingfacade) cacheSuccess(cacheKey string, value *models.TrainEventList) *models.TrainEventList {
	c.localcache.Set(cacheKey, &trainCacheEntry{value, nil}, defaultCacheTime)
	return value
}

func (c *cachingfacade) cacheFalure(cacheKey string, err error) error {
	c.localcache.Set(cacheKey, &trainCacheEntry{nil, err}, defaultCacheTime)
	return err
}

// toCacheKey creates a key in reverse order of arguments, as it probably better for indexing
func toCacheKey(methodName string, key string, value string) string {
	return fmt.Sprintf("%s:%s:%s", value, key, methodName)
}
