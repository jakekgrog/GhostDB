package base

// import (
// 	"bufio"
// 	"os"
// 	"os/user"
// 	"testing"
// 	"time"

// 	"github.com/ghostdb/ghostdb-cache-node/utils"
// 	"github.com/ghostdb/ghostdb-cache-node/config"
// 	"github.com/ghostdb/ghostdb-cache-node/store/request"
// )

// func TestWatchDog(t *testing.T) {
// 	var store *Store
	
// 	var config config.Configuration = config.InitializeConfiguration()
// 	cache = NewLRU(config)

// 	//Delete pre-existing metrics
// 	usr, _ := user.Current()
// 	configPath := usr.HomeDir
// 	os.Remove(configPath + WatchDogLogFilePath)
// 	os.Remove(configPath + "/ghostDBPersistence.log")

// 	cache = NewLRU(config)

// 	cache.Add(request.NewRequestFromValues("Key1", "Value1", -1))
// 	cache.Add(request.NewRequestFromValues("Key2", "Value1", -1))
// 	cache.Add(request.NewRequestFromValues("Key3", "Value1", -1))
// 	cache.Put(request.NewRequestFromValues("Key1", "Value2", -1))
// 	cache.Put(request.NewRequestFromValues("Key4", "Value1", -1))
// 	cache.Get(request.NewRequestFromValues("Key1", "", -1))
// 	cache.Get(request.NewRequestFromValues("Key2", "", -1))
// 	cache.Get(request.NewRequestFromValues("Key5", "", -1))
// 	cache.Delete(request.NewRequestFromValues("Key1", "", -1))
// 	cache.Delete(request.NewRequestFromValues("Key1", "", -1))
// 	cache.Flush(request.NewRequestFromValues("Key1", "", -1))
// 	cache.Flush(request.NewRequestFromValues("Key1", "", -1))
// 	time.Sleep(11 * time.Second)

// 	utils.AssertEqual(t, utils.FileExists(configPath+WatchDogLogFilePath), true, "")
// 	utils.AssertEqual(t, utils.FileNotEmpty(configPath+WatchDogLogFilePath), true, "")

// 	file, err := os.Open(configPath + WatchDogLogFilePath)
// 	if err != nil {
// 		panic(err)
// 	}
// 	scanner := bufio.NewScanner(file)

// 	//Bug: scanner.Scan() doesn't move to next token
// 	scanner.Scan()
// 	scanner.Scan()
// 	metrics := scanner.Text()
// 	expectedOutput := `{"TotalHits": 12, "TotalGets": 3, "CacheMiss": 1, "TotalPuts": 2, "TotalAdds": 3, "NotStored": 0, "TotalDeletes": 2, "NotFound": 1, "TotalFlushes": 2, "ErrFlush": 2}`
// 	utils.AssertEqual(t, metrics, expectedOutput, "")

	
// 	return
// }