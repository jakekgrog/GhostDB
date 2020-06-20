package ghost_http

import (
	// "encoding/json"
	// "errors"
	// "net/http"
	// "strconv"
	// "strings"

	// "github.com/ghostdb/ghostdb-cache-node/store/lru"
	// "github.com/valyala/fasthttp"
)

// // Global pointer to cache object for server
// var cache *lru.LRUCache

// // Represents a key-value pair body object
// type kvPair struct {
// 	Key   string `json:"Key"`
// 	Value string `json:"Value"`
// 	TTL   int64  `json:"TTL"`
// }

// // Represents response message object
// type message struct {
// 	Message string `json:"Message"`
// }

// // Handle GET requests
// func getValueHandler(ctx *fasthttp.RequestCtx) {
// 	var kv kvPair

// 	body := ctx.PostBody()

// 	if err := json.Unmarshal(body, &kv); err != nil {
// 		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 		ctx.SetStatusCode(422)
// 		if err := json.NewEncoder(ctx).Encode(err); err != nil {
// 			panic(err)
// 		}
// 	}

// 	if strings.Compare("", kv.Key) == 0 {
// 		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 		ctx.SetStatusCode(http.StatusBadRequest)
// 		if err := json.NewEncoder(ctx).Encode(kv); err != nil {
// 			panic(err)
// 		}
// 		return
// 	}

// 	kv.Value = cache.Get(kv.Key)

// 	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 	ctx.SetStatusCode(http.StatusOK)

// 	if err := json.NewEncoder(ctx).Encode(kv); err != nil {
// 		panic(err)
// 	}
// }

// // Handle FLUSH requests
// func flushCacheHandler(ctx *fasthttp.RequestCtx) {
// 	var msg message

// 	msg.Message = cache.Flush()
// 	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 	ctx.SetStatusCode(http.StatusOK)
// 	if err := json.NewEncoder(ctx).Encode(msg); err != nil {
// 		panic(err)
// 	}
// }

// // Handle ADD requests
// func addValueHandler(ctx *fasthttp.RequestCtx) {
// 	var kv kvPair
// 	var msg message
// 	body := ctx.PostBody()

// 	if err := json.Unmarshal(body, &kv); err != nil {
// 		ctx.Request.Header.Set("Content-Type", "application/jsonl charset=UTF-8")
// 		ctx.SetStatusCode(422)
// 		if err := json.NewEncoder(ctx).Encode(err); err != nil {
// 			panic(err)
// 		}
// 	}

// 	if (strings.Compare("", kv.Key) == 0) || (strings.Compare("", kv.Value) == 0) {
// 		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 		ctx.SetStatusCode(http.StatusBadRequest)
// 		if err := json.NewEncoder(ctx).Encode(kv); err != nil {
// 			panic(err)
// 		}
// 		return
// 	}

// 	msg.Message = cache.Add(kv.Key, kv.Value, kv.TTL)
// 	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 	ctx.SetStatusCode(http.StatusOK)
// 	if err := json.NewEncoder(ctx).Encode(msg); err != nil {
// 		panic(err)
// 	}
// 	return
// }

// // Handle DELETE requests
// func deleteValueHandler(ctx *fasthttp.RequestCtx) {
// 	var kv kvPair
// 	var msg message

// 	body := ctx.PostBody()

// 	if err := json.Unmarshal(body, &kv); err != nil {
// 		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 		ctx.SetStatusCode(422)
// 		if err := json.NewEncoder(ctx).Encode(err); err != nil {
// 			panic(err)
// 		}
// 	}

// 	if strings.Compare("", kv.Key) == 0 {
// 		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 		ctx.SetStatusCode(http.StatusBadRequest)
// 		if err := json.NewEncoder(ctx).Encode(kv); err != nil {
// 			panic(err)
// 		}
// 		return
// 	}

// 	msg.Message = cache.Delete(kv.Key)

// 	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 	ctx.SetStatusCode(http.StatusOK)
// 	if err := json.NewEncoder(ctx).Encode(msg); err != nil {
// 		panic(err)
// 	}
// 	return
// }

// // Handle PUT requests
// func putValueHandler(ctx *fasthttp.RequestCtx) {
// 	var kv kvPair
// 	body := ctx.PostBody()

// 	if err := json.Unmarshal(body, &kv); err != nil {
// 		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 		ctx.SetStatusCode(422)
// 		if err := json.NewEncoder(ctx).Encode(err); err != nil {
// 			panic(err)
// 		}
// 	}

// 	if (strings.Compare("", kv.Key) == 0) || (strings.Compare("", kv.Value) == 0) {
// 		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 		ctx.SetStatusCode(http.StatusBadRequest)
// 		if err := json.NewEncoder(ctx).Encode(kv); err != nil {
// 			panic(err)
// 		}
// 		return
// 	}

// 	cache.Put(kv.Key, kv.Value, kv.TTL)

// 	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 	ctx.SetStatusCode(http.StatusOK)
// 	if err := json.NewEncoder(ctx).Encode(kv); err != nil {
// 		panic(err)
// 	}
// 	return
// }

// // Handle SNITCH requests
// func getSnitchMetricsHandler(ctx *fasthttp.RequestCtx) {
// 	var msg message

// 	bytes, _ := json.Marshal(lru.GetSnitchMetrics())
// 	msg.Message = string(bytes)

// 	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 	ctx.SetStatusCode(http.StatusOK)

// 	if err := json.NewEncoder(ctx).Encode(msg); err != nil {
// 		panic(err)
// 	}
// }

// // Handle WATCHDOG requests
// func getWatchdogMetricsHandler(ctx *fasthttp.RequestCtx) {
// 	var msg message

// 	bytes, _ := json.Marshal(lru.GetWatchdogMetrics())
// 	msg.Message = string(bytes)

// 	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 	ctx.SetStatusCode(http.StatusOK)

// 	if err := json.NewEncoder(ctx).Encode(msg); err != nil {
// 		panic(err)
// 	}
// }

// // Handle PING requests
// func pingHandler(ctx *fasthttp.RequestCtx) {
// 	ctx.SetStatusCode(200)
// 	return
// }

// // Handle NODE_SIZE requests
// func nodeSizeHandler(ctx *fasthttp.RequestCtx) {
// 	var msg message

// 	count := cache.CountKeys();
// 	msg.Message = string(strconv.FormatInt(int64(count), 10))

// 	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
// 	ctx.SetStatusCode(http.StatusOK)

// 	if err := json.NewEncoder(ctx).Encode(msg); err != nil {
// 		panic(err)
// 	}
// }

// func purgeValueHandler(w http.ResponseWriter, r *http.Request) error {
// 	return errors.New("Not Implemented")
// }

// func batchGetValuesHandler(w http.ResponseWriter, r *http.Request) error {
// 	return errors.New("Not Implemented")
// }

// func batchPutValuesHandler(w http.ResponseWriter, r *http.Request) error {
// 	return errors.New("Not Implemented")
// }

// // NodeConfig configures the cache for the server
// func NodeConfig(cacheObject *lru.LRUCache) {
// 	cache = cacheObject
// }
