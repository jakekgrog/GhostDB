package crawlers

import (
	"time"

	"github.com/ghostdb/ghostdb-cache-node/store/cache"
	"github.com/ghostdb/ghostdb-cache-node/store/lru"
)

// CrawlerScheduler represents a scheduler for cache crawlers
type CrawlerScheduler struct {
	Interval time.Duration
	stop     chan bool
}

// NewCrawlerScheduler initializes a new Crawler Scheduler
func NewCrawlerScheduler(interval int32) *CrawlerScheduler {
	scheduler := &CrawlerScheduler{
		Interval: time.Duration(interval) * time.Second,
		stop:     make(chan bool),
	}

	return scheduler
}

/*
	StartCrawlers will start the cache crawler.

	The crawler is periodically run on the cache until the ticker
	is stopped.
	StartCrawlers is generic in nature and will run the appropriate
	crawler for the stores underlying cache policy implementation.

	POLICY TO CRAWLER MAP:
		1) LRU Policy -> startDllCrawler
		
*/
func StartCrawlers(cache *cache.Cache, scheduler *CrawlerScheduler) {
	ticker := time.NewTicker(scheduler.Interval)
	switch (*cache).(type) {
	case *lru.LRUCache:
		startDllCrawler((*cache).(*lru.LRUCache), ticker, scheduler)
	}
}

// startDllCrawler starts a crawler that crawls a doubly linked list
func startDllCrawler(cache *lru.LRUCache, ticker *time.Ticker, scheduler *CrawlerScheduler) {
	for {
		select {
			case <- ticker.C:
				go StartCrawl(cache)
			case <- scheduler.stop:
				ticker.Stop()
				return
		}
	}
}

// StopScheduler will stop the crawler scheduler by passing
// a boolean to the scheduler channel.
func StopScheduler(scheduler *CrawlerScheduler) {
	scheduler.stop <- true
}
