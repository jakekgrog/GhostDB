package lru_cache

import (
	"bufio"
	"os"
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/internal/ghost_config"
	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
	"github.com/ghostdb/ghostdb-cache-node/cache/utils"
)

const WatchDogLogFilePath = "/ghostdb/ghostdb_watchDog.log"

func TestWatchDog(t *testing.T) {
	var cache lru_cache.Cache
	
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	cache = lru_cache.NewLRU(config)

	//Delete pre-existing metrics
	configPath, _ := os.UserConfigDir()
	os.Remove(configPath + WatchDogLogFilePath)
	os.Remove(configPath + "/ghostDBPersistence.log")

	cache = lru_cache.NewLRU(config)

	cache.Add("Key1", "Value1", -1)
	cache.Add("Key2", "Value1", -1)
	cache.Add("Key3", "Value1", -1)
	cache.Put("Key1", "Value2", -1)
	cache.Put("Key4", "Value1", -1)
	cache.Get("Key1")
	cache.Get("Key2")
	cache.Get("Key5")
	cache.Delete("Key1")
	cache.Delete("Key1")
	cache.Flush()
	cache.Flush()
	time.Sleep(11 * time.Second)

	file, err := os.Open(configPath + WatchDogLogFilePath)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)

	//Bug: scanner.Scan() doesn't move to next token
	scanner.Scan()
	scanner.Scan()
	metrics := scanner.Text()
	expectedOutput := `{"TotalHits": 12, "TotalGets": 3, "CacheMiss": 1, "TotalPuts": 2, "TotalAdds": 3, "NotStored": 0, "TotalDeletes": 2, "NotFound": 1, "TotalFlushes": 2, "ErrFlush": 2}`
	utils.AssertEqual(t, metrics, expectedOutput, "")
	return
}
