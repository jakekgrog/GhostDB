/*
 * Copyright (c) 2020, Jake Grogan
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 *
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 *  * Neither the name of the copyright holder nor the names of its
 *    contributors may be used to endorse or promote products derived from
 *    this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

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
		case <-ticker.C:
			go StartCrawl(cache)
		case <-scheduler.stop:
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
