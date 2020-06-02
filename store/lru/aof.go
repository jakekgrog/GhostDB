package lru

import (
	"bufio"
	"bytes"
	"fmt"
	"log"
	"os"
	"time"
)

var configPath string

const logFile = "/ghostDBPersistence.log"
const tempLog = "/temp_ghostDBPersistence.log"
const writeInterval = 1

var tmpBuffer bytes.Buffer

type logFormat struct {
	Time  string `json:"Time"`
	Verb  string `json:"Verb"`
	Key   string `json:"Key"`
	Value string `json:"Value"`
	TTL   string `json:"TTL"`
}

// BootAOF Reads from log if it exists
// otherwise creates and writes to one
func BootAOF(cache *LRUCache, maxAOFSize int64) {
	configPath, _ = os.UserConfigDir()
	_, err := os.Stat(configPath + logFile)
	if err == nil {
		BuildCacheFromAof(cache, configPath+logFile)
		go flushBuffer(cache, maxAOFSize)
	} else {
		createAOF(configPath + logFile)
		go flushBuffer(cache, maxAOFSize)
	}
}

func createAOF(logPath string) {
	file, err := os.Create(logPath)
	if err != nil {
		log.Fatalf("failed to create AOF log file: %s", err.Error())
	}
	buf := bufio.NewWriter(file)
	buf.WriteString("---Created: " + time.Now().Format(time.RFC850) + "---\n")
	buf.Flush()
	file.Close()
}

func flushBuffer(cache *LRUCache, maxAOFSize int64) {
	for {
		time.Sleep(writeInterval * time.Second)
		if getAOFSize() > maxAOFSize {
			go appendBufferContent(true)
			go reduceAOF(cache)
		}
		go appendBufferContent(false)
	}
}

func appendBufferContent(dualWrite bool) {
	file, err := os.OpenFile(configPath+logFile, os.O_APPEND|os.O_WRONLY, 0600)
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

func reduceAOF(cache *LRUCache) {
	createAOF(configPath + tempLog)
	for k, v := range cache.Hashtable {
		timeStamp := time.Now().Format(time.RFC850)
		entry := fmt.Sprintf(`{"Time":"%s", "Verb":"add", "Key":"%s", "Value":"%s", "TTL":"%d"}`+"\n", timeStamp, k, v.Value, v.TTL)
		tmpBuffer.WriteString(entry)
	}
	file, err := os.OpenFile(configPath+tempLog, os.O_APPEND|os.O_WRONLY, 0600)
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

func getAOFSize() int64 {
	file, err := os.Stat(configPath + logFile)
	if err != nil {
		log.Fatalf("failed to retrieve AOF log file information: %s", err.Error())
	}
	return file.Size()
}
