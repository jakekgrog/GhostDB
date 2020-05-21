package append_only_file

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
)

var configPath string

const log = "/ghostDBPersistence.log"
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
func BootAOF(cache *lru_cache.LRUCache, maxAOFSize int64) {
	configPath, _ = os.UserConfigDir()
	_, err := os.Stat(configPath + log)
	if err == nil {
		BuildCache(cache, configPath+log)
		go flushBuffer(cache, maxAOFSize)
	} else {
		createAOF(configPath + log)
		go flushBuffer(cache, maxAOFSize)
	}
}

func createAOF(logPath string) {
	file, err := os.Create(logPath)
	if err != nil {
		panic(err)
	}
	buf := bufio.NewWriter(file)
	buf.WriteString("---Created: " + time.Now().Format(time.RFC850) + "---\n")
	buf.Flush()
	file.Close()
}

func flushBuffer(cache *lru_cache.LRUCache, maxAOFSize int64) {
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
	file, err := os.OpenFile(configPath+log, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	for _, v := range lru_cache.GetBufferBytes() {
		file.WriteString(string(v))
	}
	file.Close()
	if dualWrite {
		tmpBuffer.Write(lru_cache.GetBufferBytes())
	}
	lru_cache.FlushBuffer()
}

func reduceAOF(cache *lru_cache.LRUCache) {
	createAOF(configPath + tempLog)
	for k, v := range cache.Hashtable {
		timeStamp := time.Now().Format(time.RFC850)
		entry := fmt.Sprintf(`{"Time":"%s", "Verb":"add", "Key":"%s", "Value":"%s", "TTL":"%d"}`+"\n", timeStamp, k, v.Value, v.TTL)
		tmpBuffer.WriteString(entry)
	}
	file, err := os.OpenFile(configPath+tempLog, os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}
	for _, v := range tmpBuffer.Bytes() {
		file.WriteString(string(v))
	}
	file.Close()
	tmpBuffer.Reset()
	os.Rename(configPath+tempLog, configPath+log)
}

func getAOFSize() int64 {
	file, err := os.Stat(configPath + log)
	if err != nil {
		panic(err)
	}
	return file.Size()
}
