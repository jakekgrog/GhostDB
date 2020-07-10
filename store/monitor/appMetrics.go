/*
 * Copyright (c) 2020, Jake Grogan
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 *
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 *  * Neither the name of the copyright holder nor the names of its
 *    contributors may be used to endorse or promote products derived from
 *    this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package monitor

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/user"
	"sync"
	"sync/atomic"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/constants"
	"github.com/ghostdb/ghostdb-cache-node/store/response"
	"github.com/ghostdb/ghostdb-cache-node/utils"
)

// AppMetricsLogFilePath Log file path
const (
	AppMetricsLogFilePath  = "/ghostdb/ghostdb_appMetrics.log"
	AppMetricsTempFileName = "/ghostdb/ghostdb_appMetrics_tmp.log"
	MaxAppMetricsLogSize   = 500000
)

// AppMetrics struct is used to record cache events
type AppMetrics struct {
	// TotalRequests is the cumulative number
	// of requests to the cache node.
	TotalRequests uint64

	// GetRequests is the cumulative number
	// of Get requests to the cache node.
	GetRequests uint64

	// PutRequests is the cumulative number
	// of Put requests to the cache node.
	PutRequests uint64

	// AddRequests it the cumulative number
	// of Add requests to the cache node.
	AddRequests uint64

	// DeleteRequests is the cumulative number
	// of Delete requests to the cache node.
	DeleteRequests uint64

	// FlushRequests is the cumulative number
	// of Flush requests received by the cache node.
	FlushRequests uint64

	// CacheMiss is the cumulative number of cache misses
	// encountered by the cache node.
	CacheMiss uint64

	// Stored is the cumulative number of key-value pairs
	// successfully stored into the cache node.
	Stored uint64

	// Not stored is the cumulative number of key-value
	// pairs unsuccessfully stored into the cache node.
	NotStored uint64

	// Removed is the cumulative number of key-value pairs
	// removed from the cache node. This includes key-value
	// pairs removed by the cache crawlers.
	Removed uint64

	// NotFound is the cumulative number of key-value pairs
	// not found in the cache during a deletion.
	NotFound uint64

	// Flushed is the cumulative number of key-value pairs
	// removed from the cache node by flushes
	Flushed uint64

	// ErrFlush is the cumulative number of errors received
	// when attempting to flush the cache node.
	ErrFlush uint64

	// Mux is a mutex lock.
	Mux sync.Mutex

	// WriteInterval is the interval for writing log entries.
	WriteInterval time.Duration

	// EntryTimestamp is a bool representing whether or not to
	// include timestamps on the log entries.
	EntryTimestamp bool
}

// ReadAppMetrics struct is used to Unmarshal log entries
type ReadAppMetrics struct {
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

// Boot instantiates a appMetrics log struct and its corresponding log file
func NewAppMetrics(writeInterval time.Duration, entryTimestamp bool) *AppMetrics {
	var appMetrics AppMetrics

	appMetrics.WriteInterval = writeInterval
	appMetrics.EntryTimestamp = entryTimestamp

	usr, _ := user.Current()
	configPath := usr.HomeDir

	// Create application metrics file
	file, err := os.OpenFile(configPath+AppMetricsLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o755)
	if err != nil {
		fmt.Println(err) // Allows the CI runner to test successfully (Update when test_config is working)
	}
	defer file.Close()
	// _, err := os.OpenFile(configPath+AppMetricsLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

	go Dump(&appMetrics)

	return &appMetrics
}

func WriteMetrics(appMetrics *AppMetrics, cmd string, resp response.CacheResponse) {
	switch cmd {
	case constants.STORE_GET:
		GetHit(appMetrics)
		if resp.Status != 1 {
			CacheMiss(appMetrics)
		}
	case constants.STORE_PUT:
		PutHit(appMetrics)
		if resp.Status != 1 {
			NotStored(appMetrics)
		} else {
			Stored(appMetrics)
		}
	case constants.STORE_ADD:
		AddHit(appMetrics)
		if resp.Status != 1 {
			NotStored(appMetrics)
		} else {
			Stored(appMetrics)
		}
	case constants.STORE_DELETE:
		DeleteHit(appMetrics)
		if resp.Status != 1 {
			NotFound(appMetrics)
		} else {
			Removed(appMetrics)
		}
	case constants.STORE_FLUSH:
		FlushHit(appMetrics)
		if resp.Status != 1 {
			ErrFlush(appMetrics)
		} else {
			Flushed(appMetrics)
		}
	}
}

// ErrFlush is a setter that increments
// its corresponding struct field by one
func ErrFlush(appMetrics *AppMetrics) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.ErrFlush, 1)
	defer appMetrics.Mux.Unlock()
}

// Flushed is a setter that increments
// its corresponding struct field by one
func Flushed(appMetrics *AppMetrics) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.Flushed, 1)
	defer appMetrics.Mux.Unlock()
}

// NotFound is a setter that increments
// its corresponding struct field by one
func NotFound(appMetrics *AppMetrics) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.NotFound, 1)
	defer appMetrics.Mux.Unlock()
}

// Removed is a setter that increments
// its corresponding struct field by one
func Removed(appMetrics *AppMetrics) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.Removed, 1)
	defer appMetrics.Mux.Unlock()
}

// NotStored is a setter that increments
// its corresponding struct field by one
func NotStored(appMetrics *AppMetrics) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.NotStored, 1)
	defer appMetrics.Mux.Unlock()
}

// Stored is a setter that increments
// its corresponding struct field by one
func Stored(appMetrics *AppMetrics) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.Stored, 1)
	defer appMetrics.Mux.Unlock()
}

// CacheMiss is a setter that increments
// its corresponding struct field by one
func CacheMiss(appMetrics *AppMetrics) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.CacheMiss, 1)
	defer appMetrics.Mux.Unlock()
}

// GetHit is a setter that increments
// its corresponding struct field every time a endpoint is hit
// it also increments a total value
func GetHit(appMetrics *AppMetrics) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.GetRequests, 1)
	atomic.AddUint64(&appMetrics.TotalRequests, 1)
	defer appMetrics.Mux.Unlock()
}

// FlushHit is a setter that increments
// its corresponding struct field every time a endpoint is hit
// it also increments a total value
func FlushHit(appMetrics *AppMetrics) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.FlushRequests, 1)
	atomic.AddUint64(&appMetrics.TotalRequests, 1)
	defer appMetrics.Mux.Unlock()
}

// AddHit is a setter that increments
// its corresponding struct field every time a endpoint is hit
// it also increments a total value
func AddHit(appMetrics *AppMetrics) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.AddRequests, 1)
	atomic.AddUint64(&appMetrics.TotalRequests, 1)
	defer appMetrics.Mux.Unlock()
}

// DeleteHit is a setter that increments
// its corresponding struct field every time a endpoint is hit
// it also increments a total value
func DeleteHit(appMetrics *AppMetrics) {
	appMetrics.Mux.Lock()
	atomic.AddUint64(&appMetrics.DeleteRequests, 1)
	atomic.AddUint64(&appMetrics.TotalRequests, 1)
	defer appMetrics.Mux.Unlock()
}

// PutHit is a setter that increments
// its corresponding struct field every time a endpoint is hit
// it also increments a total value
func PutHit(appMetrics *AppMetrics) {
	appMetrics.Mux.Lock()
	defer appMetrics.Mux.Unlock()
	atomic.AddUint64(&appMetrics.PutRequests, 1)
	atomic.AddUint64(&appMetrics.TotalRequests, 1)
}

// Dump writes the contents of the appMetrics struct to the appMetrics log file
func Dump(appMetrics *AppMetrics) {
	usr, _ := user.Current()
	configPath := usr.HomeDir

	var total string
	for {
		time.Sleep(appMetrics.WriteInterval * time.Second)

		needRotate, err := utils.LogMustRotate(configPath+AppMetricsLogFilePath, MaxAppMetricsLogSize)
		if err != nil {
			log.Fatalf("failed to check if log needs to be rotated: %s", err.Error())
		}
		if needRotate {
			nBytes, err := utils.Rotate(configPath+AppMetricsLogFilePath, configPath+AppMetricsTempFileName)
			if err != nil {
				log.Fatalf("failed to rotate appMetrics log file: %s", err.Error())
			}
			log.Printf("successfully rotated appMetrics log file: %d bytes rotated", nBytes)
		}

		file, err := os.OpenFile(configPath+AppMetricsLogFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o755)
		if err != nil {
			fmt.Println(err) // Allows the CI runner to test successfully (Update when test_config is working)
		}

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
		file.Close()
	}
}

// GetAppMetrics reads the AppMetrics log
// unmarshals each entry and appends it to a slice
func GetAppMetrics() response.CacheResponse {
	usr, _ := user.Current()
	configPath := usr.HomeDir

	file, err := os.Open(configPath + AppMetricsLogFilePath)
	if err != nil {
		log.Fatalf("failed to open appMetrics log file: %s", err.Error())
	}
	defer file.Close()

	var entries []ReadAppMetrics
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry ReadAppMetrics
		line := scanner.Text()
		json.Unmarshal([]byte(line), &entry)
		entries = append(entries, entry)
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}

	res := response.NewResponseFromValue(entries)
	return res
}
