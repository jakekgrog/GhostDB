package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
	"github.com/ghostdb/ghostdb-cache-node/internal/monitoring/snitch"
	"github.com/ghostdb/ghostdb-cache-node/internal/monitoring/watchDog"
	"github.com/valyala/fasthttp"
)

var cache *lru_cache.LRUCache

type kvPair struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
	TTL   int64  `json:"TTL"`
}

type message struct {
	Message string `json:"Message"`
}

// Router passes control to handlers
func Router() {
	routes := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case "/get":
			getValue(ctx)
		case "/put":
			putValue(ctx)
		case "/add":
			addValue(ctx)
		case "/delete":
			deleteValue(ctx)
		case "/flush":
			flushCache(ctx)
		case "/getSnitchMetrics":
			getSnitchMetrics(ctx)
		case "/getWatchdogMetrics":
			getWatchdogMetrics(ctx)
		case "/ping":
			ping(ctx)
		case "/nodeSize":
			nodeSize(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
	fasthttp.ListenAndServe(":7991", routes)
}

func getValue(ctx *fasthttp.RequestCtx) {
	var kv kvPair

	body := ctx.PostBody()

	if err := json.Unmarshal(body, &kv); err != nil {
		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.SetStatusCode(422)
		if err := json.NewEncoder(ctx).Encode(err); err != nil {
			panic(err)
		}
	}

	if strings.Compare("", kv.Key) == 0 {
		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.SetStatusCode(http.StatusBadRequest)
		if err := json.NewEncoder(ctx).Encode(kv); err != nil {
			panic(err)
		}
		return
	}

	kv.Value = cache.Get(kv.Key)

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.SetStatusCode(http.StatusOK)

	if err := json.NewEncoder(ctx).Encode(kv); err != nil {
		panic(err)
	}
}

func flushCache(ctx *fasthttp.RequestCtx) {
	var msg message

	msg.Message = cache.Flush()
	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.SetStatusCode(http.StatusOK)
	if err := json.NewEncoder(ctx).Encode(msg); err != nil {
		panic(err)
	}
}

func addValue(ctx *fasthttp.RequestCtx) {
	var kv kvPair
	var msg message
	body := ctx.PostBody()

	if err := json.Unmarshal(body, &kv); err != nil {
		ctx.Request.Header.Set("Content-Type", "application/jsonl charset=UTF-8")
		ctx.SetStatusCode(422)
		if err := json.NewEncoder(ctx).Encode(err); err != nil {
			panic(err)
		}
	}

	if (strings.Compare("", kv.Key) == 0) || (strings.Compare("", kv.Value) == 0) {
		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.SetStatusCode(http.StatusBadRequest)
		if err := json.NewEncoder(ctx).Encode(kv); err != nil {
			panic(err)
		}
		return
	}

	msg.Message = cache.Add(kv.Key, kv.Value, kv.TTL)
	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.SetStatusCode(http.StatusOK)
	if err := json.NewEncoder(ctx).Encode(msg); err != nil {
		panic(err)
	}
	return
}

func deleteValue(ctx *fasthttp.RequestCtx) {
	var kv kvPair
	var msg message

	body := ctx.PostBody()

	if err := json.Unmarshal(body, &kv); err != nil {
		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.SetStatusCode(422)
		if err := json.NewEncoder(ctx).Encode(err); err != nil {
			panic(err)
		}
	}

	if strings.Compare("", kv.Key) == 0 {
		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.SetStatusCode(http.StatusBadRequest)
		if err := json.NewEncoder(ctx).Encode(kv); err != nil {
			panic(err)
		}
		return
	}

	msg.Message = cache.Delete(kv.Key)

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.SetStatusCode(http.StatusOK)
	if err := json.NewEncoder(ctx).Encode(msg); err != nil {
		panic(err)
	}
	return
}

func putValue(ctx *fasthttp.RequestCtx) {
	var kv kvPair
	body := ctx.PostBody()

	if err := json.Unmarshal(body, &kv); err != nil {
		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.SetStatusCode(422)
		if err := json.NewEncoder(ctx).Encode(err); err != nil {
			panic(err)
		}
	}

	if (strings.Compare("", kv.Key) == 0) || (strings.Compare("", kv.Value) == 0) {
		ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.SetStatusCode(http.StatusBadRequest)
		if err := json.NewEncoder(ctx).Encode(kv); err != nil {
			panic(err)
		}
		return
	}

	cache.Put(kv.Key, kv.Value, kv.TTL)

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.SetStatusCode(http.StatusOK)
	if err := json.NewEncoder(ctx).Encode(kv); err != nil {
		panic(err)
	}
	return
}

func getSnitchMetrics(ctx *fasthttp.RequestCtx) {
	var msg message

	bytes, _ := json.Marshal(snitch.GetSnitchMetrics())
	msg.Message = string(bytes)

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.SetStatusCode(http.StatusOK)

	if err := json.NewEncoder(ctx).Encode(msg); err != nil {
		panic(err)
	}
}

func getWatchdogMetrics(ctx *fasthttp.RequestCtx) {
	var msg message

	bytes, _ := json.Marshal(watchDog.GetWatchdogMetrics())
	msg.Message = string(bytes)

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.SetStatusCode(http.StatusOK)

	if err := json.NewEncoder(ctx).Encode(msg); err != nil {
		panic(err)
	}
}

func ping(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(200)
	return
}

func nodeSize(ctx *fasthttp.RequestCtx) {
	var msg message

	count := cache.CountKeys();
	msg.Message = string(strconv.FormatInt(int64(count), 10))

	ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
	ctx.SetStatusCode(http.StatusOK)

	if err := json.NewEncoder(ctx).Encode(msg); err != nil {
		panic(err)
	}
}

func purgeValue(w http.ResponseWriter, r *http.Request) error {
	return errors.New("Not Implemented")
}

func batchGetValues(w http.ResponseWriter, r *http.Request) error {
	return errors.New("Not Implemented")
}

func batchPutValues(w http.ResponseWriter, r *http.Request) error {
	return errors.New("Not Implemented")
}

// NodeConfig configures the cache for the server
func NodeConfig(cacheObject *lru_cache.LRUCache) {
	cache = cacheObject
}
