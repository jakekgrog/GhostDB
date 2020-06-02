package lru

import (
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/utils"
	"github.com/ghostdb/ghostdb-cache-node/config"
)

func TestCrawlerScheduler(t *testing.T) {
	var config config.Configuration = config.InitializeConfiguration()

	var cache *LRUCache
	cache = NewLRU(config)

	sch := NewCrawlerScheduler(int32(5))

	utils.AssertEqual(t, sch.Interval, time.Duration(int32(5)) * time.Second, "")

	go StartCrawlers(cache, sch)

	cache.Put("England", "London", 2)
	cache.Put("Italy", "Rome", -1)

	cache.Mux.Lock()
	utils.AssertEqual(t, cache.Count, int32(2), "")
	cache.Mux.Unlock()

	time.Sleep(6 * time.Second)

	cache.Mux.Lock()
	utils.AssertEqual(t, cache.Count, int32(1), "")
	cache.Mux.Unlock()

	StopScheduler(sch)
}
