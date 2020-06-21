package base

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/utils"
	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/store/persistence"
)

func TestAOF(t *testing.T) {
	configPath, _ := os.UserConfigDir()

	var config config.Configuration = config.InitializeConfiguration()
	const aofMaxBytes = 300

	var err = os.Remove(configPath + "/ghostDBPersistence.log")
	if err != nil {
		return
	}

	var store *Store
	store = NewStore("LRU")
	store.BuildStore(config)
	store.RunStore()

	
	store.Execute("add", request.NewRequestFromValues("Key1", "Value1", -1))
	store.Execute("add", request.NewRequestFromValues("Key2", "Value2", -1))

	for i := 1; i <= 100; i++ {
		store.Execute("put", request.NewRequestFromValues("Key1", "NewValue1", -1))
		time.Sleep(10 * time.Millisecond)
	}

	// Give go routine time to flush/write buffer & rewrite file
	time.Sleep(2 * time.Second)
	// Check file has shrunk below max size
	if persistence.GetAOFSize() >= aofMaxBytes {
		fmt.Println(persistence.GetAOFSize())
		t.Error("AOF Size exceeded threshold")
	}
	// Simulate cache restart
	var newStore *Store
	newStore = NewStore("LRU")
	newStore.BuildStore(config)
	newStore.RunStore()

	utils.AssertEqual(t, newStore.Execute("get", request.NewRequestFromValues("Key1", "", -1)), "NewValue1", "")
}
