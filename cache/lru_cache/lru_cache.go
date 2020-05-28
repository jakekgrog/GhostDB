package lru_cache

import (
	"log"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/internal/ghost_config"
	"github.com/ghostdb/ghostdb-cache-node/internal/monitoring/watchDog"
	"github.com/ghostdb/ghostdb-cache-node/cache/linked_list"
)

const (
	CACHE_MISS = "CACHE_MISS"
	STORED     = "STORED"
	NOT_STORED = "NOT_STORED"
	REMOVED    = "REMOVED"
	NOT_FOUND  = "NOT_FOUND"
	FLUSHED    = "FLUSH"
	ERR_FLUSH  = "ERR_FLUSH"
)

// Cache is an interface for the cache object
type Cache interface {
	// Put will add a key/value pair to the cache, possibly
	// overwriting an existing key/value pair. Put will evict
	// a key/value pair if the cache is full.
	Put(k string, v string, ttl int64) string

	// Get will fetch a key/value pair from the cache
	Get(k string) string

	// Add will add a key/value pair to the cache if the key
	// does not exist already. It will not evict a key/value pair
	// from the cache. If the cache is full, the key/value pair does
	// not get added.
	Add(k string, v string, ttl int64) string

	// Delete removes a key/value pair from the cache
	// Returns NOT_FOUND if the key does not exist.
	Delete(k string) string

	// Flush removes all key/value pairs from the cache even if they
	// have not expired
	Flush() string

	// CountKeys return the number of keys in the cache
	CountKeys() int32
}

// LRUCache represents a cache object
type LRUCache struct {
	// Size represents the maximum number of allowable
	// key-value pairs in the cache.
	Size      int32

	// Count records the number of key-value pairs
	// currently in the cache.
	Count     int32

	// Full tracks if Count is equal to Size
	Full      bool

	// DLL is a doubly linked list containing all key-value pairs
	DLL       *linked_list.List `json:"omitempty"`

	// Hashtable maps to nodes in the doubly linked list
	Hashtable map[string]*linked_list.Node

	// Config is the user configuration for the cache node.
	// This is instantiated at startup and persists for the lifetime
	// of the node.
	Config    ghost_config.Configuration

	// Watchdog is the application metrics logging subsystem
	WatchDog  *watchDog.WatchDog
	
	// Mux is a mutex lock
	Mux       sync.Mutex
}

// NewLRU will initialize the cache
func NewLRU(config ghost_config.Configuration) *LRUCache {
	wdMetricInterval := time.Duration(config.WatchdogMetricInterval)
	return &LRUCache{
		Size:      config.KeyspaceSize,
		Count:     int32(0),
		Full:      false,
		DLL:       linked_list.InitList(),
		Hashtable: newHashtable(),
		Config:    config,
		WatchDog:  watchDog.Boot(wdMetricInterval, config.EntryTimestamp),
	}
}

func newHashtable() map[string]*linked_list.Node {
	return make(map[string]*linked_list.Node)
}

// Get will fetch a key/value pair from the cache
func (cache *LRUCache) Get(key string) string {

	cache.Mux.Lock()
	watchDog.GetHit(cache.WatchDog)
	cache.Mux.Unlock()

	cache.Mux.Lock()
	nodeToGet := cache.Hashtable[key]
	cache.Mux.Unlock()

	if nodeToGet == nil {
		cache.Mux.Lock()
		watchDog.CacheMiss(cache.WatchDog)
		cache.Mux.Unlock()
		return CACHE_MISS
	}

	cache.Mux.Lock()
	n, _ := linked_list.RemoveNode(cache.DLL, nodeToGet)
	cache.Mux.Unlock()

	cache.Mux.Lock()
	node, _ := linked_list.Insert(cache.DLL, n.Key, n.Value, n.TTL)
	cache.Mux.Unlock()

	cache.Mux.Lock()
	cache.Hashtable[key] = node
	cache.Mux.Unlock()

	return node.Value
}

// Put will add a key/value pair to the cache, possibly
// overwriting an existing key/value pair. Put will evict
// a key/value pair if the cache is full.
func (cache *LRUCache) Put(key string, value string, ttl int64) string {

	cache.Mux.Lock()
	watchDog.PutHit(cache.WatchDog)
	cache.Mux.Unlock()

	if !cache.Full {
		inCache := keyInCache(cache, key)

		newNode, _ := linked_list.Insert(cache.DLL, key, value, ttl)

		insertIntoHashtable(cache, key, newNode)

		if !inCache {
			cache.Mux.Lock()
			atomic.AddInt32(&cache.Count, 1)
			cache.Mux.Unlock()
		}

		if cache.Count == cache.Size {
			cache.Full = true
		}

	} else {
		// SPECIAL CASE: Just update the value
		inCache := keyInCache(cache, key)
		if inCache {
			// Get the value node
			node, _ := cache.Hashtable[key]
	
			// Update the value
			node.Value = value
			return STORED
		} else {
			n, _ := linked_list.RemoveLast(cache.DLL)

			deleteFromHashtable(cache, n.Key)

			newNode, _ := linked_list.Insert(cache.DLL, key, value, ttl)
			insertIntoHashtable(cache, key, newNode)
		}
	}

	if cache.Config.PersistenceAOF {
		WriteBuffer("put", key, value, ttl)
	}
	return STORED
}

func deleteFromHashtable(cache *LRUCache, key string) {
	cache.Mux.Lock()
	defer cache.Mux.Unlock()
	delete(cache.Hashtable, key)
}

// Add will add a key/value pair to the cache if the key
// does not exist already. It will not evict a key/value pair
// from the cache. If the cache is full, the key/value pair does
// not get added.
func (cache *LRUCache) Add(key string, value string, ttl int64) string {

	cache.Mux.Lock()
	watchDog.AddHit(cache.WatchDog)
	cache.Mux.Unlock()

	cache.Mux.Lock()
	_, ok := cache.Hashtable[key]
	cache.Mux.Unlock()
	if ok {
		cache.Mux.Lock()
		watchDog.NotStored(cache.WatchDog)
		cache.Mux.Unlock()
		return NOT_STORED
	}
	if !cache.Full {
		inCache := keyInCache(cache, key)

		newNode, _ := linked_list.Insert(cache.DLL, key, value, ttl)

		insertIntoHashtable(cache, key, newNode)

		if !inCache {
			atomic.AddInt32(&cache.Count, 1)
		}

		if cache.Count == cache.Size {
			cache.Full = true
		}
	} else {
		n, _ := linked_list.RemoveLast(cache.DLL)
		deleteFromHashtable(cache, n.Key)

		newNode, _ := linked_list.Insert(cache.DLL, key, value, ttl)
		insertIntoHashtable(cache, key, newNode)
	}
	if cache.Config.PersistenceAOF {
		WriteBuffer("add", key, value, ttl)
	}
	cache.Mux.Lock()
	watchDog.Stored(cache.WatchDog)
	cache.Mux.Unlock()
	return STORED
}

// Delete removes a key/value pair from the cache
// Returns NOT_FOUND if the key does not exist.
func (cache *LRUCache) Delete(key string) string {

	cache.Mux.Lock()
	watchDog.DeleteHit(cache.WatchDog)
	cache.Mux.Unlock()

	cache.Mux.Lock()
	_, ok := cache.Hashtable[key]
	cache.Mux.Unlock()
	if ok {
		cache.Mux.Lock()
		nodeToRemove := cache.Hashtable[key]
		cache.Mux.Unlock()

		if nodeToRemove == nil {
			cache.Mux.Lock()
			watchDog.NotFound(cache.WatchDog)
			cache.Mux.Unlock()
			return NOT_FOUND
		}

		deleteFromHashtable(cache, nodeToRemove.Key)
		_, err := linked_list.RemoveNode(cache.DLL, nodeToRemove)

		if err != nil {
			log.Println("failed to remove key-value pair")
		}

		cache.Mux.Lock()
		atomic.AddInt32(&cache.Count, -1)
		cache.Mux.Unlock()

		cache.Full = false
		if cache.Config.PersistenceAOF {
			WriteBuffer("delete", key, "NA", -1)
		}
		cache.Mux.Lock()
		watchDog.Removed(cache.WatchDog)
		cache.Mux.Unlock()
		return REMOVED
	}
	cache.Mux.Lock()
	watchDog.NotFound(cache.WatchDog)
	cache.Mux.Unlock()
	return NOT_FOUND
}

// Flush removes all key/value pairs from the cache even if they have not expired
func (cache *LRUCache) Flush() string {

	cache.Mux.Lock()
	watchDog.FlushHit(cache.WatchDog)
	cache.Mux.Unlock()

	for k := range cache.Hashtable {
		n, _ := linked_list.RemoveLast(cache.DLL)
		if n == nil {
			break
		}
		delete(cache.Hashtable, k)
		cache.Mux.Lock()
		atomic.AddInt32(&cache.Count, -1)
		cache.Mux.Unlock()
	}
	
	if cache.Count == int32(0) {
		if cache.Config.PersistenceAOF {
			WriteBuffer("flush", "NA", "NA", -1)
		}
		cache.Mux.Lock()
		watchDog.Flushed(cache.WatchDog)
		cache.Mux.Unlock()
		return FLUSHED
	}
	cache.Mux.Lock()
	watchDog.ErrFlush(cache.WatchDog)
	cache.Mux.Unlock()
	return ERR_FLUSH
}

// CountKeys return the number of keys in the cache
func (cache *LRUCache) CountKeys() int32 {
	return cache.Count
}

func insertIntoHashtable(cache *LRUCache, key string, node *linked_list.Node) {
	cache.Mux.Lock()
	defer cache.Mux.Unlock()
	cache.Hashtable[key] = node
}

func keyInCache(cache *LRUCache, key string) (bool) {
	cache.Mux.Lock()
	defer cache.Mux.Unlock()
	_, ok := cache.Hashtable[key] 
	if ok {
		return true
	}
	return false
}