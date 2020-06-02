package lru

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/config"
)

//BuildCache parses a pre-existing AOF
//rebuilds cache using contents
func BuildCacheFromAof(cache *LRUCache, logPath string) {
	file, err := os.Open(logPath)
	if err != nil {
		log.Fatalf("failed to open AOF log file: %s", err.Error())
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan() //Ignore creation date
	for scanner.Scan() {
		var lf logFormat
		aofEntry := []byte(scanner.Text())
		err = json.Unmarshal(aofEntry, &lf)
		if err != nil {
			// If line is incomplete ignore it
			continue
		}

		switch lf.Verb {
		case "flush":
			cache.Flush()
		case "put":
			n, err := strconv.ParseInt(lf.TTL, 10, 64)
			if err != nil {
				log.Fatalf("failed to parse AOF log entry: %s", err.Error())
			}
			cache.Put(lf.Key, lf.Value, n)
		case "add":
			n, err := strconv.ParseInt(lf.TTL, 10, 64)
			if err != nil {
				log.Fatalf("failed to parse AOF log entry: %s", err.Error())
			}
			cache.Add(lf.Key, lf.Value, n)
		case "delete":
			cache.Delete(lf.Key)
		}
	}
}

// BuildCache rebuilds the cache from the byte stream of the snapshot
func BuildCacheFromSnapshot(bs *[]byte) (*LRUCache, error) {
	// Create a new cache instance.
	var cache LRUCache

	// Create a new configuration object
	var config config.Configuration = config.InitializeConfiguration()
	
	// Unmarshal the byte stream and update the new cache object with the result.
	err := json.Unmarshal(*bs, &cache)
	
	if err != nil {
		log.Fatalf("failed to rebuild cache from snapshot: %s", err.Error())
	}

	cache.Config = config

	// Create a new doubly linked list object
	ll := InitList()

	// Populate the caches hashtable and doubly linked list with the values 
	// from the unmarshalled byte stream
	for _, v := range cache.Hashtable {
		n, err := Insert(ll, v.Key, v.Value, v.TTL)
		if err != nil {
			return &LRUCache{}, err
		}
		cache.Hashtable[v.Key] = n
	}

	// Reset the watchdog
	wdMetricInterval := time.Duration(config.WatchdogMetricInterval)
	cache.WatchDog = Boot(wdMetricInterval, config.EntryTimestamp)

	cache.DLL = ll

	return &cache, nil
}