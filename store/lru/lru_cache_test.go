package lru

import (
	"testing"

	"github.com/ghostdb/ghostdb-cache-node/utils"
	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
)

func TestLru(t *testing.T) {
	var config config.Configuration = config.InitializeConfiguration()
	cache := NewLRU(config)
	cache.Size = int32(2);
	cache.Put(request.NewRequestFromValues("England", "London", -1))
	cache.Put(request.NewRequestFromValues("Ireland", "Dublin", -1))
	
	// HEAD -> Dublin -> London

	cache.Put(request.NewRequestFromValues("America", "Washington", -1)) // England:London evicted here

	// HEAD -> Washington -> Dublin

	v1 := cache.Get(request.NewRequestFromValues("England", "", -1)) // Should be a cache miss
	utils.AssertEqual(t, "CACHE_MISS", v1.Message, "")

	// Ireland should be next to be evicted
	// If we 'Get' Ireland then it should be considered MRU
	// And 'America' Should now be LRU
	v2 := cache.Get(request.NewRequestFromValues("Ireland", "", -1))
	utils.AssertEqual(t, "Dublin", v2.Gobj.Value, "")
	
	// HEAD -> Dublin -> Washington

	cache.Put(request.NewRequestFromValues("France", "Paris", -1)) // America should be evicted here
	
	// HEAD -> Paris -> Dublin

	v3 := cache.Get(request.NewRequestFromValues("America", "", -1)) // Should be a cache miss
	utils.AssertEqual(t, CACHE_MISS, v3.Message, "")
	
	cache.Put(request.NewRequestFromValues("Italy", "Rome", -1)) // Ireland should be evicted here

	// HEAD -> Rome -> Paris
	
	v4 := cache.Get(request.NewRequestFromValues("France", "", -1))
	utils.AssertEqual(t, "Paris", v4.Gobj.Value, "")

	// HEAD -> Paris -> Rome

	message := cache.Add(request.NewRequestFromValues("France", "Paris", -1))
	utils.AssertEqual(t, NOT_STORED, message.Message, "")

	message = cache.Add(request.NewRequestFromValues("Poland", "Warsaw", -1))
	utils.AssertEqual(t, STORED, message.Message, "")

	message = cache.Delete(request.NewRequestFromValues("Poland", "", -1))
	utils.AssertEqual(t, REMOVED, message.Message, "")

	message = cache.CountKeys(request.NewRequestFromValues("Key1", "", -1))
	utils.AssertEqual(t, message.Gobj.Value.(int32) > 0, true, "")

	message = cache.Delete(request.NewRequestFromValues("USA", "", -1))
	utils.AssertEqual(t, NOT_FOUND, message.Message, "")

	message = cache.Put(request.NewRequestFromValues("England", "London", -1))
	utils.AssertEqual(t, STORED, message.Message, "")

	message = cache.Put(request.NewRequestFromValues("England", "London", -1))
	utils.AssertEqual(t, STORED, message.Message, "")

	message = cache.Flush(request.NewRequestFromValues("Key1", "", -1))
	utils.AssertEqual(t, FLUSHED, message.Message, "")

	message = cache.CountKeys(request.NewRequestFromValues("Key1", "", -1))
	utils.AssertEqual(t, message.Gobj.Value.(int32), int32(0), "")
}
