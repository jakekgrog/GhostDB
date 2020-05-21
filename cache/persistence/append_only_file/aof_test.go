package append_only_file

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/internal/ghost_config"
	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
	"github.com/ghostdb/ghostdb-cache-node/cache/utils"
)

func TestAOF(t *testing.T) {
	configPath, _ = os.UserConfigDir()
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	var cache *lru_cache.LRUCache
	const aofMaxBytes = 300
	cache = lru_cache.NewLRU(config)

	var err = os.Remove(configPath + "/ghostDBPersistence.log")
	if err != nil {
		return
	}

	BootAOF(cache, aofMaxBytes)

	cache.Add("Key1", "Value1", -1)
	cache.Add("Key2", "Value2", -1)

	for i := 1; i <= 100; i++ {
		cache.Put("Key1", "NewValue1", -1)
		time.Sleep(10 * time.Millisecond)
	}

	// Give go routine time to flush/write buffer & rewrite file
	time.Sleep(2 * time.Second)
	// Check file has shrunk below max size
	if getAOFSize() >= aofMaxBytes {
		fmt.Println(getAOFSize())
		t.Error("AOF Size exceeded threshold")
	}
	// Simulate cache restart
	newCache := lru_cache.NewLRU(config)
	// AOF will see pre-existing log and attempt to rebuild
	BootAOF(newCache, aofMaxBytes)

	utils.AssertEqual(t, newCache.Get("Key1"), "NewValue1", "")
}
