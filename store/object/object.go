package object

type CacheObject struct {
	Key   string `json:"Key"`
	Value interface{} `json:"Value"`
	TTL   int64 `json:"TTL"`
}

func NewCacheObjectFromValue(value interface{}) CacheObject{
	return CacheObject {
		Key: "",
		Value: value,
		TTL: -1,
	}
}

func NewCacheObjectFromParams(key string, value interface{}, ttl int64) CacheObject {
	return CacheObject {
		Key: key,
		Value: value,
		TTL: ttl,
	}
}

func NewEmptyCacheObject() CacheObject {
	return CacheObject {
		Key: "",
		Value: nil,
		TTL: -1,
	}
}