package server

import (
	"github.com/valyala/fasthttp"
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

// Router passes control to handlers
func Router() {
	routes := func(ctx *fasthttp.RequestCtx) {
		switch string(ctx.Path()) {
		case GHOST_GET:
			getValueHandler(ctx)
		case GHOST_PUT:
			putValueHandler(ctx)
		case GHOST_ADD:
			addValueHandler(ctx)
		case GHOST_DELETE:
			deleteValueHandler(ctx)
		case GHOST_FLUSH:
			flushCacheHandler(ctx)
		case GHOST_SNITCH:
			getSnitchMetricsHandler(ctx)
		case GHOST_WATCHDOG:
			getWatchdogMetricsHandler(ctx)
		case GHOST_PING:
			pingHandler(ctx)
		case GHOST_NODE_SIZE:
			nodeSizeHandler(ctx)
		default:
			ctx.Error("not found", fasthttp.StatusNotFound)
		}
	}
	fasthttp.ListenAndServe(":7991", routes)
}