package lru

import (
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/utils"
)

func TestCrawler(t *testing.T) {
	var cache *LRUCache

	var config config.Configuration = config.InitializeConfiguration()
	cache = NewLRU(config)

	cache.Put("England", "London", 5)
	cache.Put("Italy", "Rome", -1)
	cache.Put("Ireland", "Dublin", 11)
	time.Sleep(10 * time.Second)
	go StartCrawl(cache)
	time.Sleep(2 * time.Second)

	// Node with key "England" should be considered stale after 10 seconds
	// and should therefore be evicted by the crawler.
	// Node with key "Italy" should not be evicted as its TTL is set to
	// never expire (-1).
	// Node with key "Ireland" should not be evicted as it's TTL has not expired.
	cache.Mux.Lock()
	utils.AssertEqual(t, cache.Count, int32(2), "")
	cache.Mux.Unlock()
}
