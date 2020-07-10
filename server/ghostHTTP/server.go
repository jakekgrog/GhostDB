/*
 * Copyright (c) 2020, Jake Grogan
 * All rights reserved.
 *
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 *
 *  * Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 *
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 *
 *  * Neither the name of the copyright holder nor the names of its
 *    contributors may be used to endorse or promote products derived from
 *    this software without specific prior written permission.
 *
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
 */

package ghostHTTP

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ghostdb/ghostdb-cache-node/store/base"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/store/response"
	"github.com/ghostdb/ghostdb-cache-node/systemMonitor"
	"github.com/valyala/fasthttp"
)

var (
	GHOST_GET         = "/get"
	GHOST_PUT         = "/put"
	GHOST_ADD         = "/add"
	GHOST_DELETE      = "/delete"
	GHOST_FLUSH       = "/flush"
	GHOST_SYS_METRICS = "/getSysMetrics"
	GHOST_APP_METRICS = "/getAppMetrics"
	GHOST_PING        = "/ping"
	GHOST_NODE_SIZE   = "/nodeSize"
)

var store *base.Store

// NodeConfig configures the store for the server
func NodeConfig(s *base.Store) {
	store = s
}

// Router passes control to handlers
func Router() {
	routes := func(ctx *fasthttp.RequestCtx) {
		req := new(request.CacheRequest)
		path := ctx.Path()
		cmd := string(path[1:])
		body := ctx.PostBody()

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
			res = systemMonitor.GetSysMetrics()
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
	err := fasthttp.ListenAndServe(":7991", routes)
	if err != nil{
		log.Fatal(err)
	}
}
