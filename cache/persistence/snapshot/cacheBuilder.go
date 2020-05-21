package snapshot

import (
	"encoding/json"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/internal/ghost_config"
	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
	"github.com/ghostdb/ghostdb-cache-node/cache/linked_list"
	"github.com/ghostdb/ghostdb-cache-node/internal/monitoring/watchDog"
)

// BuildCache rebuilds the cache from the byte stream of the snapshot
func BuildCache(bs *[]byte) (*lru_cache.LRUCache, error) {
	// Create a new cache instance.
	var cache lru_cache.LRUCache

	// Create a new configuration object
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	
	// Unmarshal the byte stream and update the new cache object with the result.
	err := json.Unmarshal(*bs, &cache)
	
	if err != nil {
		panic(err)
	}

	cache.Config = config

	// Create a new doubly linked list object
	ll := linked_list.InitList()

	// Populate the caches hashtable and doubly linked list with the values 
	// from the unmarshalled byte stream
	for _, v := range cache.Hashtable {
		n, err := linked_list.Insert(ll, v.Key, v.Value, v.TTL)
		if err != nil {
			return &lru_cache.LRUCache{}, err
		}
		cache.Hashtable[v.Key] = n
	}

	// Reset the watchdog
	wdMetricInterval := time.Duration(config.WatchdogMetricInterval)
	cache.WatchDog = watchDog.Boot(wdMetricInterval, config.EntryTimestamp)

	cache.DLL = ll

	return &cache, nil
}
