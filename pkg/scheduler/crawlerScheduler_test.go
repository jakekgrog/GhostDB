package scheduler

import (
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/internal/ghost_config"
	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
	"github.com/ghostdb/ghostdb-cache-node/cache/scheduler"
	"github.com/ghostdb/ghostdb-cache-node/cache/utils"
)

func TestCrawlerScheduler(t *testing.T) {
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()

	var cache *lru_cache.LRUCache
	cache = lru_cache.NewLRU(config)

	sch := scheduler.NewCrawlerScheduler(int32(5))

	utils.AssertEqual(t, sch.Interval, time.Duration(int32(5)) * time.Second, "")

	go scheduler.StartCrawlers(cache, sch)

	cache.Put("England", "London", 2)
	cache.Put("Italy", "Rome", -1)

	cache.Mux.Lock()
	utils.AssertEqual(t, cache.Count, int32(2), "")
	cache.Mux.Unlock()

	time.Sleep(6 * time.Second)

	cache.Mux.Lock()
	utils.AssertEqual(t, cache.Count, int32(1), "")
	cache.Mux.Unlock()

	scheduler.StopScheduler(sch)
}
