package lru_cache

import (
	"testing"

	"github.com/ghostdb/ghostdb-cache-node/internal/ghost_config"
	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
	"github.com/ghostdb/ghostdb-cache-node/cache/utils"
)

func TestLru(t *testing.T) {
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	cache := lru_cache.NewLRU(config)
	cache.Size = int32(2);
	cache.Put("England", "London", -1)
	cache.Put("Ireland", "Dublin", -1)

	// HEAD -> Dublin -> London

	cache.Put("America", "Washington", -1) // England:London evicted here

	// HEAD -> Washington -> Dublin

	v1 := cache.Get("England") // Should be a cache miss
	utils.AssertEqual(t, "CACHE_MISS", v1, "")

	// Ireland should be next to be evicted
	// If we 'Get' Ireland then it should be considered MRU
	// And 'America' Should now be LRU
	v2 := cache.Get("Ireland")
	utils.AssertEqual(t, "Dublin", v2, "")
	
	// HEAD -> Dublin -> Washington

	cache.Put("France", "Paris", -1) // America should be evicted here
	
	// HEAD -> Paris -> Dublin

	v3 := cache.Get("America") // Should be a cache miss
	utils.AssertEqual(t, lru_cache.CACHE_MISS, v3, "")
	
	cache.Put("Italy", "Rome", -1) // Ireland should be evicted here

	// HEAD -> Rome -> Paris
	
	v4 := cache.Get("France")
	utils.AssertEqual(t, "Paris", v4, "")

	// HEAD -> Paris -> Rome

	message := cache.Add("France", "paris", -1)
	utils.AssertEqual(t, lru_cache.NOT_STORED, message, "")

	message = cache.Add("Poland", "Warsaw", -1)
	utils.AssertEqual(t, lru_cache.STORED, message, "")

	message = cache.Delete("Poland")
	utils.AssertEqual(t, lru_cache.REMOVED, message, "")

	utils.AssertEqual(t, cache.CountKeys() > 0, true, "")

	message = cache.Delete("USA")
	utils.AssertEqual(t, lru_cache.NOT_FOUND, message, "")

	message = cache.Flush()
	utils.AssertEqual(t, lru_cache.FLUSHED, message, "")

	utils.AssertEqual(t, cache.CountKeys(), int32(0), "")
}
