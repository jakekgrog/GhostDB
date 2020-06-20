package lru

// import (
// 	"bufio"
// 	"encoding/json"
// 	"log"
// 	"os"
// 	"strconv"

// 	"github.com/ghostdb/ghostdb-cache-node/store/request"
// )

// //BuildCache parses a pre-existing AOF
// //rebuilds cache using contents
// func BuildCacheFromAof(cache *LRUCache, logPath string) {
// 	file, err := os.Open(logPath)
// 	if err != nil {
// 		log.Fatalf("failed to open AOF log file: %s", err.Error())
// 	}
// 	scanner := bufio.NewScanner(file)
// 	scanner.Scan() //Ignore creation date
// 	for scanner.Scan() {
// 		var lf logFormat
// 		aofEntry := []byte(scanner.Text())
// 		err = json.Unmarshal(aofEntry, &lf)
// 		if err != nil {
// 			// If line is incomplete ignore it
// 			continue
// 		}

// 		// Convert the log entry to a cache object
// 		n, err := strconv.ParseInt(lf.TTL, 10, 64)
// 		cacheRequest := request.NewRequestFromValues(lf.Key, lf.Value, n)

// 		switch lf.Verb {
// 		case "flush":
// 			cache.Flush()
// 		case "put":
// 			if err != nil {
// 				log.Fatalf("failed to parse AOF log entry: %s", err.Error())
// 			}
// 			cache.Put(cacheRequest)
// 		case "add":
// 			if err != nil {
// 				log.Fatalf("failed to parse AOF log entry: %s", err.Error())
// 			}
// 			cache.Add(cacheRequest)
// 		case "delete":
// 			cache.DeleteByKey(cacheRequest.Gobj.Key)
// 		}
// 	}
// }

// func logEntryToCacheRequest(logEntry *logFormat) request.CacheRequest {
// 	n, _ := strconv.ParseInt(logEntry.TTL, 10, 64)
// 	return request.NewRequestFromValues(logEntry.Key, logEntry.Value, n)
// }