package lru

import (
	"time"
)

// StartCrawl crawls the cache and evicts stale data
func StartCrawl(cache *LRUCache) {
	markedKeys := mark(cache)
	sweep(cache, markedKeys)
	return
}

// Traverse the cache and mark key-value pair nodes
// for removal.
func mark(cache *LRUCache) []string {
	markedKeys := []string{}

	node, _ := GetLastNode(cache.DLL)

	// List is empty
	if node == nil {
		return []string{}
	}

	// Crawl until node.Prev is nil i.e. the Head Node
	for ok := true; ok; ok = !(node.Prev == nil) {
		node.Mux.Lock()

		if node.TTL != -1 {
			now := time.Now().Unix()

			if node.CreatedAt+node.TTL < now {
				markedKeys = append(markedKeys, node.Key)
			}
		}
		node.Mux.Unlock()
		node = node.Prev
	}

	return markedKeys
}

// Sweep the cache removing the marked nodes
func sweep(cache *LRUCache, keys []string) {
	for _, key := range keys {
		cache.Delete(key)
	}
	return
}
