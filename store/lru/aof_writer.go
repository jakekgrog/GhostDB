package lru

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

var buffer bytes.Buffer

// WriteBuffer writes cache command in log format
func WriteBuffer(verb string, key string, value string, ttl int64) {
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
