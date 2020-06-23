package ghost_http

import (
	"encoding/json"
	"net/http"
	
	"github.com/valyala/fasthttp"
	"github.com/ghostdb/ghostdb-cache-node/store/base"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
)

var (
	GHOST_GET = "/get"
	GHOST_PUT = "/put"
	GHOST_ADD = "/add"
	GHOST_DELETE = "/delete"
	GHOST_FLUSH = "/flush"
	GHOST_SNITCH = "/getSnitchMetrics"
	GHOST_WATCHDOG = "/getWatchdogMetrics"
	GHOST_PING = "/ping"
	GHOST_NODE_SIZE = "/nodeSize"
)

var store *base.Store

// NodeConfig configures the store for the server
func NodeConfig(s *base.Store) {
	store = s
}

// Router passes control to handlers
func Router() {
	routes := func(ctx *fasthttp.RequestCtx) {
		var req request.CacheRequest
		var path = ctx.Path()
		var cmd = string(path[1:])
		var body = ctx.PostBody()

		if err := json.Unmarshal(body, &req); err != nil {
			ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
			ctx.SetStatusCode(422)
			if err := json.NewEncoder(ctx).Encode(err); err != nil {
				panic(err)
			}
		}

		var res = store.Execute(cmd, req)
		
		ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.SetStatusCode(http.StatusOK)

		if err := json.NewEncoder(ctx).Encode(res); err != nil {
			panic(err)
		}
	}
	fasthttp.ListenAndServe(":7991", routes)
}