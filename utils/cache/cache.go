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

type cachingfacade struct {
	localcache *cache.Cache
	facade     dsb.DSBFacade
}

func New(facade dsb.DSBFacade) dsb.DSBFacade {
	c := cache.New(30*time.Second, 10*time.Second)
	return &cachingfacade{c, facade}
}

func (c *cachingfacade) GetStations() (chan *models.StationList, chan error) {
	return c.facade.GetStations()
}

type trainCacheEntry struct {
	list *models.TrainList
	err  error
}

// GetTrains searches the embedded cache for a result otherwise deletates the request to the underlying
// facade.
func (c *cachingfacade) GetTrains(key string, value string) (chan *models.TrainList, chan error) {
	logger.Println("Intercepting request")
	cacheKey := toCacheKey("GetTrains", key, value)
	success, failure := make(chan *models.TrainList), make(chan error)
	go func() {
		if cachedItem, found := c.localcache.Get(cacheKey); found {
			v := cachedItem.(*trainCacheEntry)
			if v.err != nil {
				failure <- v.err
			} else {
				success <- v.list
			}
			logger.Println("Returning cached result")
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

func (c *cachingfacade) cacheSuccess(cacheKey string, value *models.TrainList) *models.TrainList {
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
