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

package persistence

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/store/cache"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
)

var buffer bytes.Buffer

var configPath string

const (
	logFile       = "/ghostDBPersistence.log"
	tempLog       = "/temp_ghostDBPersistence.log"
	writeInterval = 1
)

var tmpBuffer bytes.Buffer

type logFormat struct {
	Time  string `json:"Time"`
	Verb  string `json:"Verb"`
	Key   string `json:"Key"`
	Value string `json:"Value"`
	TTL   string `json:"TTL"`
}

/*
	TODO:
		BootAOF takes a reference to a store cache
		and populates the cache.
*/

// BootAOF Reads from log if it exists
// otherwise creates and writes to one
func BootAOF(cache *cache.Cache, maxAOFSize int64) {
	CreateAOF(getLogPath())
	go flushBuffer(cache, maxAOFSize)
}

func RebootAof(cache *cache.Cache, maxAofSize int64) {
	BuildCacheFromAof(cache, getLogPath())
	go flushBuffer(cache, maxAofSize)
}

func AofExists() (bool, error) {
	logPath := getLogPath()
	_, err := os.Stat(logPath)
	if err != nil {
		return false, err
	}
	return true, err
}

func getLogPath() string {
	configPath, _ = os.UserConfigDir()
	return configPath + logFile
}

func getTempLogPath() string {
	configPath, _ = os.UserConfigDir()
	return configPath + tempLog
}

func CreateAOF(logPath string) {
	file, err := os.Create(logPath)
	if err != nil {
		log.Fatalf("failed to create AOF log file: %s", err.Error())
	}
	buf := bufio.NewWriter(file)
	buf.WriteString("---Created: " + time.Now().Format(time.RFC850) + "---\n")
	buf.Flush()
	file.Close()
}

func flushBuffer(cache *cache.Cache, maxAOFSize int64) {
	for {
		time.Sleep(writeInterval * time.Second)
		if GetAOFSize() > maxAOFSize {
			go appendBufferContent(true)
			go reduceAOF(cache)
		}
		go appendBufferContent(false)
	}
}

func appendBufferContent(dualWrite bool) {
	file, err := os.OpenFile(configPath+logFile, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		log.Fatalf("failed to write to AOF log file: %s", err.Error())
	}
	for _, v := range GetBufferBytes() {
		file.WriteString(string(v))
	}
	file.Close()
	if dualWrite {
		tmpBuffer.Write(GetBufferBytes())
	}
	FlushBuffer()
}

func reduceAOF(cache *cache.Cache) {
	CreateAOF(getTempLogPath())
	for k, v := range *((*cache).GetHashtableReference()) {
		timeStamp := time.Now().Format(time.RFC850)
		entry := fmt.Sprintf(`{"Time":"%s", "Verb":"add", "Key":"%s", "Value":"%s", "TTL":"%d"}`+"\n", timeStamp, k, v.Value, v.TTL)
		tmpBuffer.WriteString(entry)
	}
	file, err := os.OpenFile(configPath+tempLog, os.O_APPEND|os.O_WRONLY, 0o600)
	if err != nil {
		log.Fatalf("failed to create temporary AOF log file for AOF reduction: %s", err.Error())
	}
	for _, v := range tmpBuffer.Bytes() {
		file.WriteString(string(v))
	}
	file.Close()
	tmpBuffer.Reset()
	os.Rename(configPath+tempLog, configPath+logFile)
}

func GetAOFSize() int64 {
	file, err := os.Stat(configPath + logFile)
	if err != nil {
		log.Fatalf("failed to retrieve AOF log file information: %s", err.Error())
	}
	return file.Size()
}

// WriteBuffer writes cache command in log format
func WriteBuffer(verb string, key string, value interface{}, ttl int64) {
	timeStamp := time.Now().Format(time.RFC850)
	var event string
	if strings.Compare(verb, "flush") == 0 {
		event = fmt.Sprintf(`{"Time":"%s", "Verb":"%s", "Key":"NA", "Value":"NA", "TTL":"-1"}`+"\n", timeStamp, verb)
	} else {
		event = fmt.Sprintf(`{"Time":"%s", "Verb":"%s", "Key":"%s", "Value":"%s", "TTL":"%d"}`+"\n", timeStamp, verb, key, value, ttl)
	}
	buffer.WriteString(event)
}

// GetBuffer returns buffer
func GetBufferBytes() []byte {
	return buffer.Bytes()
}

func GetBufferString() string {
	return buffer.String()
}

func FlushBuffer() {
	buffer.Reset()
}

// BuildCache parses a pre-existing AOF
// rebuilds cache using contents
func BuildCacheFromAof(cache *cache.Cache, logPath string) {
	file, err := os.Open(logPath)
	if err != nil {
		log.Fatalf("failed to open AOF log file: %s", err.Error())
	}
	scanner := bufio.NewScanner(file)
	scanner.Scan() // Ignore creation date
	for scanner.Scan() {
		var lf logFormat
		aofEntry := []byte(scanner.Text())
		err = json.Unmarshal(aofEntry, &lf)
		if err != nil {
			// If line is incomplete ignore it
			continue
		}

		// Convert the log entry to a cache object
		n, err := strconv.ParseInt(lf.TTL, 10, 64)
		cacheRequest := request.NewRequestFromValues(lf.Key, lf.Value, n)

		switch lf.Verb {
		case "flush":
			(*cache).Flush(request.NewEmptyRequest())
		case "put":
			if err != nil {
				log.Fatalf("failed to parse AOF log entry: %s", err.Error())
			}
			(*cache).Put(cacheRequest)
		case "add":
			if err != nil {
				log.Fatalf("failed to parse AOF log entry: %s", err.Error())
			}
			(*cache).Add(cacheRequest)
		case "delete":
			(*cache).DeleteByKey(cacheRequest.Gobj.Key)
		}
	}
}

func logEntryToCacheRequest(logEntry *logFormat) request.CacheRequest {
	n, _ := strconv.ParseInt(logEntry.TTL, 10, 64)
	return request.NewRequestFromValues(logEntry.Key, logEntry.Value, n)
}
