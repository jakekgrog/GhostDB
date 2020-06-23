package base

import (
	"bufio"
	"os"
	"os/user"
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/utils"
	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/store/monitor"
)

func TestWatchDog(t *testing.T) {
	var config config.Configuration = config.InitializeConfiguration()
	
	var store *Store
	store = NewStore("LRU")
	store.BuildStore(config)
	store.RunStore()

	//Delete pre-existing metrics
	usr, _ := user.Current()
	configPath := usr.HomeDir
	os.Remove(configPath + monitor.WatchDogLogFilePath)
	os.Remove(configPath + "/ghostDBPersistence.log")

	store.Execute("add", request.NewRequestFromValues("Key1", "Value1", -1))
	store.Execute("add", request.NewRequestFromValues("Key2", "Value1", -1))
	store.Execute("add", request.NewRequestFromValues("Key3", "Value1", -1))
	store.Execute("put", request.NewRequestFromValues("Key1", "Value2", -1))
	store.Execute("put", request.NewRequestFromValues("Key4", "Value1", -1))
	store.Execute("get", request.NewRequestFromValues("Key1", "", -1))
	store.Execute("get", request.NewRequestFromValues("Key2", "", -1))
	store.Execute("get", request.NewRequestFromValues("Key5", "", -1))
	store.Execute("delete", request.NewRequestFromValues("Key1", "", -1))
	store.Execute("delete", request.NewRequestFromValues("Key1", "", -1))
	store.Execute("flush", request.NewRequestFromValues("Key1", "", -1))
	store.Execute("flush", request.NewRequestFromValues("Key1", "", -1))
	time.Sleep(11 * time.Second)

	utils.AssertEqual(t, utils.FileExists(configPath+monitor.WatchDogLogFilePath), true, "")
	utils.AssertEqual(t, utils.FileNotEmpty(configPath+monitor.WatchDogLogFilePath), true, "")

	file, err := os.Open(configPath + monitor.WatchDogLogFilePath)
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