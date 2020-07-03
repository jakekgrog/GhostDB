package ghost_http

import (
	"encoding/json"
	"log"
	"net/http"
	
	"github.com/valyala/fasthttp"
	"github.com/ghostdb/ghostdb-cache-node/store/base"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/store/response"
	"github.com/ghostdb/ghostdb-cache-node/system_monitor"
)

var (
	GHOST_GET = "/get"
	GHOST_PUT = "/put"
	GHOST_ADD = "/add"
	GHOST_DELETE = "/delete"
	GHOST_FLUSH = "/flush"
	GHOST_SYS_METRICS = "/getSysMetrics"
	GHOST_APP_METRICS = "/getAppMetrics"
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
		var req = new(request.CacheRequest)
		var path = ctx.Path()
		var cmd = string(path[1:])
		var body = ctx.PostBody()

		if err := json.Unmarshal(body, &req); err != nil {
			log.Println(err)
			ctx.Request.Header.Set("Content-Type", "application/json; charset=UTF-8")
			ctx.SetStatusCode(422)
			if err := json.NewEncoder(ctx).Encode(err); err != nil {
				panic(err)
			}
		}

		var res response.CacheResponse
		// Handle SysMet
		if cmd == "getSysMetrics" {
			res = system_monitor.GetSysMetrics()
		} else if cmd == "ping" {
			res = response.NewPingResponse()
		} else {
			res = store.Execute(cmd, *req)
		}
		
		ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.SetStatusCode(http.StatusOK)

		if err := json.NewEncoder(ctx).Encode(res); err != nil {
			panic(err)
		}
	}
	fasthttp.ListenAndServe(":7991", routes)
}
