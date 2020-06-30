package request

import (
	"github.com/ghostdb/ghostdb-cache-node/store/object"
)

type CacheRequest struct {
	Gobj object.CacheObject `json:"Gobj"`
}

func NewRequestFromValues(key string, value interface{}, ttl int64) CacheRequest {
	return CacheRequest{
		Gobj: object.NewCacheObjectFromParams(key, value, ttl),
	}
}

func NewEmptyRequest() CacheRequest {
	return CacheRequest{
		Gobj: object.NewEmptyCacheObject(),
	}
}