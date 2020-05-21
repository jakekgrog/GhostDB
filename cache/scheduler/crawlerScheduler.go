package scheduler

import (
	"time"

	"github.com/ghostdb/ghostdb-cache-node/cache/crawler"
	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
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

// StartCrawlers will periodically run crawlers on the cache
// until the ticker is stopped.
func StartCrawlers(cache *lru_cache.LRUCache, scheduler *CrawlerScheduler) {
	ticker := time.NewTicker(scheduler.Interval)

	for {
		select {
			case <- ticker.C:
				go crawler.StartCrawl(cache)
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
