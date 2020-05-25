package watchDog

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"os/user"
	"sync"
	"sync/atomic"
	"time"
)

// WatchDogLogFilePath Log file path
const (
	WatchDogLogFilePath = "/ghostdb/ghostdb_watchDog.log"
	WatchDogTempFileName = "/ghostdb/ghostdb_watchDog_tmp.log"
	MaxWatchDogLogSize = 10
)

// WatchDog struct is used to record cache events
type WatchDog struct {
	// TotalRequests is the cumulative number
	// of requests to the cache node.
	TotalRequests  uint64

	// GetRequests is the cumulative number
	// of Get requests to the cache node.
	GetRequests    uint64

	// PutRequests is the cumulative number
	// of Put requests to the cache node.
	PutRequests    uint64

	// AddRequests it the cumulative number
	// of Add requests to the cache node.
	AddRequests    uint64

	// DeleteRequests is the cumulative number
	// of Delete requests to the cache node.
	DeleteRequests uint64

	// FlushRequests is the cumulative number
	// of Flush requests received by the cache node.
	FlushRequests  uint64

	// CacheMiss is the cumulative number of cache misses
	// encountered by the cache node.
	CacheMiss uint64

	// Stored is the cumulative number of key-value pairs
	// successfully stored into the cache node.
	Stored    uint64

	// Not stored is the cumulative number of key-value
	// pairs unsuccessfully stored into the cache node.
	NotStored uint64

	// Removed is the cumulative number of key-value pairs
	// removed from the cache node. This includes key-value
	// pairs removed by the cache crawlers.
	Removed   uint64

	// NotFound is the cumulative number of key-value pairs
	// not found in the cache during a deletion.
	NotFound  uint64

	// Flushed is the cumulative number of key-value pairs
	// removed from the cache node by flushes
	Flushed   uint64

	// ErrFlush is the cumulative number of errors received
	// when attempting to flush the cache node.
	ErrFlush  uint64

	// Mux is a mutex lock.
	Mux            sync.Mutex

	// WriteInterval is the interval for writing log entries.
	WriteInterval  time.Duration

	// EntryTimestamp is a bool representing whether or not to
	// include timestamps on the log entries.
	EntryTimestamp bool
}

// ReadWatchDog struct is used to Unmarshal log entries
type ReadWatchDog struct {
	Timestamp      string `json:"Timestamp"`
	TotalRequests  uint64 
	GetRequests    uint64 
	PutRequests    uint64 
	AddRequests    uint64 
	DeleteRequests uint64 
	FlushRequests  uint64 

	CacheMiss uint64 
	Stored    uint64 `json: "-"`
	NotStored uint64
	Removed   uint64 `json: "-"`
	NotFound  uint64 
	Flushed   uint64 `json: "-"`
	ErrFlush  uint64 
}

// Boot instantiates a watchdog log struct and its corresponding log file
func Boot(writeInterval time.Duration, entryTimestamp bool) *WatchDog {
	var watchDog WatchDog

	watchDog.WriteInterval = writeInterval
	watchDog.EntryTimestamp = entryTimestamp

	usr, _ := user.Current()
	configPath := usr.HomeDir

	// Create application metrics file
	_, err := os.Create(configPath + WatchDogLogFilePath)
	if err != nil {
		fmt.Println(err) // Allows the CI runner to test successfully (Update when test_config is working)
	}
	// _, err := os.OpenFile(configPath+WatchDogLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	go Dump(&watchDog)

	return &watchDog
}

// ErrFlush is a setter that increments
// its corresponding struct field by one
func ErrFlush(appMetrics *WatchDog) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.ErrFlush, 1)
	defer appMetrics.Mux.Unlock()
}

// Flushed is a setter that increments
// its corresponding struct field by one
func Flushed(appMetrics *WatchDog) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.Flushed, 1)
	defer appMetrics.Mux.Unlock()
}

// NotFound is a setter that increments
// its corresponding struct field by one
func NotFound(appMetrics *WatchDog) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.NotFound, 1)
	defer appMetrics.Mux.Unlock()
}

// Removed is a setter that increments
// its corresponding struct field by one
func Removed(appMetrics *WatchDog) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.Removed, 1)
	defer appMetrics.Mux.Unlock()
}

// NotStored is a setter that increments
// its corresponding struct field by one
func NotStored(appMetrics *WatchDog) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.NotStored, 1)
	defer appMetrics.Mux.Unlock()
}

// Stored is a setter that increments
// its corresponding struct field by one
func Stored(appMetrics *WatchDog) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.Stored, 1)
	defer appMetrics.Mux.Unlock()
}

// CacheMiss is a setter that increments
// its corresponding struct field by one
func CacheMiss(appMetrics *WatchDog) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.CacheMiss, 1)
	defer appMetrics.Mux.Unlock()
}

// GetHit is a setter that increments
// its corresponding struct field every time a endpoint is hit
// it also increments a total value
func GetHit(appMetrics *WatchDog) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.GetRequests, 1)
	atomic.AddUint64(&appMetrics.TotalRequests, 1)
	defer appMetrics.Mux.Unlock()
}

// FlushHit is a setter that increments
// its corresponding struct field every time a endpoint is hit
// it also increments a total value
func FlushHit(appMetrics *WatchDog) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.FlushRequests, 1)
	atomic.AddUint64(&appMetrics.TotalRequests, 1)
	defer appMetrics.Mux.Unlock()
}

// AddHit is a setter that increments
// its corresponding struct field every time a endpoint is hit
// it also increments a total value
func AddHit(appMetrics *WatchDog) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.AddRequests, 1)
	atomic.AddUint64(&appMetrics.TotalRequests, 1)
	defer appMetrics.Mux.Unlock()
}

// DeleteHit is a setter that increments
// its corresponding struct field every time a endpoint is hit
// it also increments a total value
func DeleteHit(appMetrics *WatchDog) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.DeleteRequests, 1)
	atomic.AddUint64(&appMetrics.TotalRequests, 1)
	defer appMetrics.Mux.Unlock()
}

// PutHit is a setter that increments
// its corresponding struct field every time a endpoint is hit
// it also increments a total value
func PutHit(appMetrics *WatchDog) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.PutRequests, 1)
	atomic.AddUint64(&appMetrics.TotalRequests, 1)
	defer appMetrics.Mux.Unlock()
}

// Dump writes the contents of the watchdog struct to the watchdog log file
func Dump(appMetrics *WatchDog) {
	usr, _ := user.Current()
	configPath := usr.HomeDir

	var total string
	for {
		time.Sleep(appMetrics.WriteInterval * time.Second)

		needRotate, err := logMustRotate(configPath + WatchDogLogFilePath)
		if err != nil {
			log.Fatalf("failed to check if log needs to be rotated: %s", err.Error())
		}
		if needRotate {
			nBytes, err := Rotate()
			if err != nil {
				log.Fatalf("failed to rotate watchdog log file: %s", err.Error())
			}
			log.Printf("successfully rotated watchdog log file: %d bytes rotated", nBytes)
		}

		file, err := os.OpenFile(configPath + WatchDogLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0755)
		if err != nil {
			fmt.Println(err) // Allows the CI runner to test successfully (Update when test_config is working)
		}
		defer file.Close()

		if appMetrics.EntryTimestamp {
			total = fmt.Sprintf(`{"Timestamp": "%s", "TotalRequsts": %d, `, time.Now().Format(time.RFC3339), appMetrics.TotalRequests)
		} else {
			total = fmt.Sprintf(`{"TotalRequests": %d, `, appMetrics.TotalRequests)
		}

		getMetrics := fmt.Sprintf(`"GetRequests": %d, "CacheMiss": %d, `, appMetrics.GetRequests, appMetrics.CacheMiss)
		putMetrics := fmt.Sprintf(`"PutRequests": %d, `, appMetrics.PutRequests)
		addMetrics := fmt.Sprintf(`"AddRequsets": %d, "NotStored": %d, `, appMetrics.AddRequests, appMetrics.NotStored)
		deleteMetrics := fmt.Sprintf(`"DeleteRequests": %d, "NotFound": %d, `, appMetrics.DeleteRequests, appMetrics.NotFound)
		flushMetrics := fmt.Sprintf(`"FlushRequests": %d, "ErrFlush": %d}`+"\n", appMetrics.FlushRequests, appMetrics.ErrFlush)

		file.WriteString(total + getMetrics + putMetrics + addMetrics + deleteMetrics + flushMetrics)	
	}
}

// GetWatchdogMetrics reads the Watchdog log
// unmarshals each entry and appends it to a slice
func GetWatchdogMetrics() []ReadWatchDog {
	usr, _ := user.Current()
	configPath := usr.HomeDir

	file, err := os.OpenFile(configPath + WatchDogLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open watchdog log file: %s", err.Error())
	}
	var entries []ReadWatchDog
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry ReadWatchDog
		line := scanner.Text()
		json.Unmarshal([]byte(line), &entry)
		entries = append(entries, entry)
	}
	return entries
}

func logMustRotate(logfile string) (bool, error) {
	fi, err := os.Stat(logfile)
	if err != nil {
		return false, err
	}
	// get the size
	size := fi.Size()
	if size >= MaxWatchDogLogSize {
		return true, nil
	}
	return false, nil
}

// Rotate rotates the main watchdog log by copying the contents of the 
// log file to a temporary log and clearing the main log.
// If there is data in the temp log file, it is cleared. 
func Rotate() (int64, error) {
	usr, _ := user.Current()
	configPath := usr.HomeDir

	src := configPath + WatchDogLogFilePath
	tmp := configPath + WatchDogTempFileName

	// Check if tmp file exists
	exists, err := tmpFileExists(tmp)
	if err != nil {
		return 0, fmt.Errorf("Error when checking for temp log existence: %s", err.Error())
	}

	// If it exists, clear it
	if exists {
		_, err := cleanFile(tmp)
		if err != nil {
			return 0, fmt.Errorf("failed to clean temp log")
		}
	}

	// Open the source file (main log)
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, fmt.Errorf("failed to stat log file")
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	// Open the tmp log (or create if it doesn't exist)
	dst, err := os.OpenFile(tmp, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return 0, fmt.Errorf("failed to open %s temporary log file", tmp)
	}
	defer dst.Close()

	// Copy the contents of main log to tmp log
	nBytes, err := io.Copy(dst, source)

	if err != nil {
		return 0, fmt.Errorf("failed to copy log to temp log")
	}

	// clear the main log
	_, err = cleanFile(src)

	if err != nil {
		return 0, fmt.Errorf("failed to clean watchdog log file")
	}

	return nBytes, err
}

func tmpFileExists(tmpFilename string) (bool, error) {
	if _, err := os.Stat(tmpFilename); os.IsNotExist(err) {
		dst, err := os.OpenFile(tmpFilename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return false, fmt.Errorf("failed to open %s temporary log file", tmpFilename)
		}
		defer dst.Close()
	}
	return true, nil
}

func cleanFile(filePath string) (bool, error) {
	err := os.Remove(filePath)

	if err != nil {
		return false, err
	}

	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false, err
	}
	defer file.Close()

	return true, err
}
