package append_only_file

import (
	"bufio"
	"encoding/json"
	"os"
	"strconv"

	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
)

//BuildCache parses a pre-existing AOF
//rebuilds cache using contents
func BuildCache(cache *lru_cache.LRUCache, logPath string) {
	file, err := os.Open(logPath)
	if err != nil {
		panic(err)
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
				panic(err)
			}
			cache.Put(lf.Key, lf.Value, n)
		case "add":
			n, err := strconv.ParseInt(lf.TTL, 10, 64)
			if err != nil {
				panic(err)
			}
			cache.Add(lf.Key, lf.Value, n)
		case "delete":
			cache.Delete(lf.Key)
		}
	}
}
