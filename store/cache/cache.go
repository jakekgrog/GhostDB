package cache

import (
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/store/response"
	"github.com/ghostdb/ghostdb-cache-node/store/lru"
)

// Cache is an interface for the cache object
type Cache interface {
	// Put will add a key/value pair to the cache, possibly
	// overwriting an existing key/value pair. Put will evict
	// a key/value pair if the cache is full.
	Put(reqObj request.CacheRequest) response.CacheResponse

	// Get will fetch a key/value pair from the cache
	Get(reqObj request.CacheRequest) response.CacheResponse

	// Add will add a key/value pair to the cache if the key
	// does not exist already. It will not evict a key/value pair
	// from the cache. If the cache is full, the key/value pair does
	// not get added.
	Add(reqObj request.CacheRequest) response.CacheResponse

	// Delete removes a key/value pair from the cache
	// Returns NOT_FOUND if the key does not exist.
	Delete(reqObj request.CacheRequest) response.CacheResponse

	// DeleteByKey functions the same as Delete, however it is
	// used in various locations to reduce the cost of allocating
	// request objects for internal deletion mechanisms 
	// e.g. the cache crawlers.
	DeleteByKey(key string) response.CacheResponse

	// Flush removes all key/value pairs from the cache even if they
	// have not expired
	Flush(request.CacheRequest) response.CacheResponse

	// CountKeys return the number of keys in the cache
	CountKeys(request.CacheRequest) response.CacheResponse

	// GetHashtableReference is for internal use by crawlers and AOF
	GetHashtableReference() *map[string]*lru.Node
}