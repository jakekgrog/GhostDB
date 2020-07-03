package response

import (
	"fmt"
	
	"github.com/ghostdb/ghostdb-cache-node/store/object"
)

const (
	INVALID_COMMAND_ERR = "INVALID_COMMAND_ERR"
)

type CacheResponse struct {
	// GhostDB Cache Object
	Gobj    object.CacheObject
	// Status of the result - 0 or 1 - fail or succeed
	Status  int32
	// Message contains any textual data that is to 
	// be sent back
	Message string
	// Error message returned if something went wrong
	// during command execution
	Error   string
}

func NewResponseFromValue(value interface{}) CacheResponse{
	return CacheResponse{
		Gobj: object.NewCacheObjectFromValue(value),
		Status: 1,
		Message: "OK",
		Error: "",
	}
}

func NewCacheMissResponse() CacheResponse {
	return CacheResponse {
		Gobj: object.NewEmptyCacheObject(),
		Status: 0,
		Message: "CACHE_MISS",
		Error: "",
	}
}

func NewResponseFromMessage(msg string, status int32) CacheResponse {
	return CacheResponse {
		Gobj: object.NewEmptyCacheObject(),
		Status: status,
		Message: msg,
		Error: "",
	}
}

func BadCommandResponse(cmd string) CacheResponse {
	return CacheResponse {
		Gobj: object.NewEmptyCacheObject(),
		Status: 0,
		Message: fmt.Sprintf("Command '%s' is not a recognized command", cmd),
		Error: INVALID_COMMAND_ERR,
	}
}

func NewPingResponse() CacheResponse {
	return CacheResponse {
		Gobj: object.NewEmptyCacheObject(),
		Status: 1,
		Message: "Pong!",
		Error: "",
	}
}