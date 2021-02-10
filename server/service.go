package server

import (
	"encoding/json"
	"log"
	"net/http"
	"fmt"

	"github.com/valyala/fasthttp"
	"github.com/ghostdb/ghostdb-cache-node/store/base"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/store/response"
	"github.com/ghostdb/ghostdb-cache-node/system_monitor"
)

var (
    // HTTPAddr contains the port used by this node
    HTTPAddr string
)

// Service is a type to be used by the raft consensus protocol
// consists of a base http address and a store in the finite state machine
type Service struct {
	addr  string
	store *base.Store
}

// NewService is used to initialize a new service struct
// parameters: addr (a string of a http address), store (a store of node details)
// returns: *Service (a newly initialized service struct)
func NewService(addr string, store *base.Store) *Service {
	return &Service{
		addr:  addr,
		store: store,
	}
}

func (service *Service) Start() {
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
		} else if cmd == "join" {
			handleJoin(ctx, service.store)
		} else {
			fmt.Println(service.store.RaftDir)
			res = service.store.Execute(cmd, *req)
		}
		
		ctx.Response.Header.Set("Content-Type", "application/json; charset=UTF-8")
		ctx.SetStatusCode(http.StatusOK)

		if err := json.NewEncoder(ctx).Encode(res); err != nil {
			panic(err)
		}
	}

	HTTPAddr = service.addr
	log.Println("Serving...")
	fasthttp.ListenAndServe(HTTPAddr, routes)
}

type JoinRequest struct {
	Addr   string `json:"addr"`
	Id interface{} `json:"id"`
}

func handleJoin(ctx *fasthttp.RequestCtx, store *base.Store) {
	m := make(map[string]string)
	if err := json.Unmarshal(ctx.PostBody(), &m); err != nil {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	if len(m) != 2 {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	remoteAddr, ok := m["addr"]
	if !ok {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	nodeID, ok := m["id"]
	if !ok {
		ctx.SetStatusCode(http.StatusBadRequest)
		return
	}

	if err := store.Join(nodeID, remoteAddr); err != nil {
		panic(err)
	}
}