package lru

// import (
// 	"fmt"
// 	"os"
// 	"testing"
// 	"time"

// 	"github.com/ghostdb/ghostdb-cache-node/utils"
// 	"github.com/ghostdb/ghostdb-cache-node/config"
// 	"github.com/ghostdb/ghostdb-cache-node/store/request"
// )

// func TestAOF(t *testing.T) {
// 	configPath, _ = os.UserConfigDir()
// 	var config config.Configuration = config.InitializeConfiguration()
// 	var cache *LRUCache
// 	const aofMaxBytes = 300
// 	cache = NewLRU(config)

// 	var err = os.Remove(configPath + "/ghostDBPersistence.log")
// 	if err != nil {
// 		return
// 	}

// 	BootAOF(cache, aofMaxBytes)

// 	cache.Add(request.NewRequestFromValues("Key1", "Value1", -1))
// 	cache.Add(request.NewRequestFromValues("Key2", "Value2", -1))

// 	for i := 1; i <= 100; i++ {
// 		cache.Put(request.NewRequestFromValues("Key1", "NewValue1", -1))
// 		time.Sleep(10 * time.Millisecond)
// 	}

// 	// Give go routine time to flush/write buffer & rewrite file
// 	time.Sleep(2 * time.Second)
// 	// Check file has shrunk below max size
// 	if getAOFSize() >= aofMaxBytes {
// 		fmt.Println(getAOFSize())
// 		t.Error("AOF Size exceeded threshold")
// 	}
// 	// Simulate cache restart
// 	newCache := NewLRU(config)
// 	// AOF will see pre-existing log and attempt to rebuild
// 	BootAOF(newCache, aofMaxBytes)

// 	utils.AssertEqual(t, newCache.Get(request.NewRequestFromValues("Key1", "", -1)), "NewValue1", "")
// }
