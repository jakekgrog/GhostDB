package lru

import (
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/utils"
	"github.com/ghostdb/ghostdb-cache-node/config"
)

func TestSnapshotScheduler(t *testing.T) {
	var config config.Configuration = config.InitializeConfiguration()

	var cache *LRUCache
	cache = NewLRU(config)

	sch := NewSnapshotScheduler(int32(5))

	utils.AssertEqual(t, sch.Interval, time.Duration(int32(5)) * time.Second, "")

	go StartSnapshotter(cache, sch)

	cache.Put("England", "London", 2)
	cache.Put("Italy", "Rome", -1)

	time.Sleep(6 * time.Second)

	StopSnapshotter(sch)
}
