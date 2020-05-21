package scheduler

import (
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/internal/ghost_config"
	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
	"github.com/ghostdb/ghostdb-cache-node/cache/scheduler"
	"github.com/ghostdb/ghostdb-cache-node/cache/utils"
)

func TestSnapshotScheduler(t *testing.T) {
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()

	var cache *lru_cache.LRUCache
	cache = lru_cache.NewLRU(config)

	sch := scheduler.NewSnapshotScheduler(int32(5))

	utils.AssertEqual(t, sch.Interval, time.Duration(int32(5)) * time.Second, "")

	go scheduler.StartSnapshotter(cache, sch)

	cache.Put("England", "London", 2)
	cache.Put("Italy", "Rome", -1)

	time.Sleep(6 * time.Second)

	scheduler.StopSnapshotter(sch)
}
