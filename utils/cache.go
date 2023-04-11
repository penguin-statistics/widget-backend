package utils

import (
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

// AcceptableFailCount describes how many fails would be tolerated until the server refuses to serve any more data
const AcceptableFailCount = 10

// DefaultTimeout describes the default timeout for the underlying http.Client that the updater uses
const DefaultTimeout = time.Second * 20

// Updater describes a func type which will be used to call and get data update from
type Updater func(instance *Cache) (interface{}, error)

// CacheConfig describes the configuration to provide when initializing a cache instance
type CacheConfig struct {
	// Name is the name of this cache instance
	Name string

	// Server is the server of this cache instance
	Server string

	// Interval is the update Interval of cache
	Interval time.Duration

	// Updater is the function to update of underlying cache instance
	Updater Updater
}

// Cache is the thread-safe cache instance that automatically updates data every Interval with Updater
type Cache struct {
	// == data ==

	// Updated is to represent the last updated time
	Updated *time.Time

	// content is the actual content of the cache; nullable
	content interface{}

	// FailCount is the counter for how many times does the cache failed to update (after retries)
	FailCount int

	// == instance ==
	// Name is the name of this cache instance
	Name string

	// Server is the server of this cache instance
	Server string

	// Interval is the update interval of cache
	Interval time.Duration

	// Client is the http.Client
	Client *http.Client

	// update is the function to update of underlying cache instance
	Updater Updater

	// Logger is the logger of such instance
	Logger *logrus.Entry

	mutex *sync.RWMutex
}

// Update calls the underlying cache's Updater and updates cache status accordingly
func (c *Cache) Update() error {
	c.Logger.Debugln("updating instance with new data...")
	data, err := c.Updater(c)
	if err != nil {
		return err
	}

	c.mutex.Lock()

	now := time.Now()
	c.Updated = &now
	c.content = data
	c.FailCount = 0

	c.mutex.Unlock()

	c.Logger.Infoln("successfully updated instance with new data")

	return nil
}

// Ready indicates whether the cache is ready to serve valid data. Notice that if the existing data have exceeded the FailCount limit, this check would also fail
func (c *Cache) Ready() bool {
	if c.content != nil && c.FailCount <= AcceptableFailCount {
		return true
	}
	return false
}

// NewCache creates a new Cache with CacheConfig provided
func NewCache(config CacheConfig) *Cache {
	loggerName := "cache:" + config.Name
	if config.Server != "" {
		loggerName += ":" + config.Server
	}

	instance := &Cache{
		Name:     config.Name,
		Server:   config.Server,
		Interval: config.Interval,
		Client: &http.Client{
			Timeout: DefaultTimeout,
		},
		Updater: config.Updater,
		Logger:  NewLogger(loggerName),

		mutex: &sync.RWMutex{},
	}

	onFetchErr := func(err error) {
		if err != nil {
			instance.Logger.Errorln("error occurred when trying to fetch new data", err)

			instance.mutex.Lock()
			instance.FailCount++
			instance.mutex.Unlock()
		}
	}

	instance.Logger.Debugln("instance created")

	ticker := time.NewTicker(config.Interval)

	onFetchErr(instance.Update())

	go func() {
		for range ticker.C {
			onFetchErr(instance.Update())
		}
	}()

	return instance
}

// Content retrieves the cache content with proper mutex handling
func (c *Cache) Content() interface{} {
	c.mutex.RLock()
	defer c.mutex.RUnlock()
	return c.content
}
