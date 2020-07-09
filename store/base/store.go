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

package base

import (
	"log"
	"time"
	
	"github.com/ghostdb/ghostdb-cache-node/store/cache"
	"github.com/ghostdb/ghostdb-cache-node/store/lru"
	"github.com/ghostdb/ghostdb-cache-node/store/crawlers"
	"github.com/ghostdb/ghostdb-cache-node/store/persistence"
	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/store/response"
	"github.com/ghostdb/ghostdb-cache-node/store/monitor"
)

// LRU store commands
const (
	STORE_GET = "get"
	STORE_PUT = "put"
	STORE_ADD = "add"
	STORE_DELETE = "delete"
	STORE_FLUSH = "flush"
	STORE_NODE_SIZE = "nodeSize"
	STORE_APP_METRICS = "getAppMetrics"
)

// Policy types
const (
	LRU_TYPE = "LRU"   // Least recently used
	LFU_TYPE = "LFU"   // Least frequenty used
	MRU_TYPE = "MRU"   // Most recently used
	ARC_TYPE = "ARC"   // Adaptive Replacement Cache
	TLRU_TYPE = "TLRU" // Time-aware Least Recently Used
)

type HandlerType func(request.CacheRequest) response.CacheResponse

type BaseStore interface {
	// Maps commands to functions
	registerHandlers() map[string]interface{}
	// Create a store cache from a policy type
	newCacheFromPolicy(p string) interface{}
	// Build store from a simple store
	BuildStore(config.Configuration)
	// Build store from snapshot
	BuildStoreFromSnapshot()
	// Execute a given command
	Execute(cmd string, args interface{}) interface{}
	// Create snapshot of stores cache
	CreateSnapshot()
	// RunStore starts the store schedulers
	RunStore()
}

type Store struct {
	Conf   config.Configuration
	policy   string
	Cache    cache.Cache
	commands map[string]interface{}
	crawlerScheduler *crawlers.CrawlerScheduler
	snapshotScheduler *persistence.SnapshotScheduler
	appMetrics *monitor.AppMetrics
}

func NewStore(policy string) *Store {
	return &Store{
		policy: policy,
	}
}

func (store *Store) Execute(cmd string, args request.CacheRequest) response.CacheResponse {
	// Handle App metrics
	if cmd == "getAppMetrics" {
		return monitor.GetAppMetrics()
	}

	// Handle all other commands
	if _, ok := store.commands[cmd]; !ok {
		return response.BadCommandResponse(cmd)
	}
	execResult := store.commands[cmd].(func(request.CacheRequest) response.CacheResponse)(args)
	if store.Conf.PersistenceAOF {
		writeAof(cmd, &args)
	}
	// CHECK RESPONSE AND SEND TO APP METRICS
	monitor.WriteMetrics(store.appMetrics, cmd, execResult)
	return execResult
}

func writeAof(cmd string, args *request.CacheRequest) {
	if isWriteOp(cmd) {
		var gobj = args.Gobj
		persistence.WriteBuffer(cmd, gobj.Key, gobj.Value, gobj.TTL)
	}
}

func isWriteOp(cmd string) bool {
	writeOps := map[string]bool {
		STORE_ADD: true,
		STORE_PUT: true,
		STORE_DELETE: true,
		STORE_FLUSH: true,
	}
	return writeOps[cmd]
}

func (store *Store) CreateSnapshot() {
	_, err := persistence.CreateSnapshot(&store.Cache, &store.Conf)
	if err != nil {
		log.Println("Failed to create PIT Snapshot of the stores cache!")
	}
	log.Println("Sucessfully created PIT Snapshot of the stores cache!")
}

func (store *Store) BuildStore(conf config.Configuration) {
	store.Conf = conf
	store.Cache = store.newCacheFromPolicy(store.policy)
	store.commands = make(map[string]interface{})
	store.commands = store.registerHandlers()
	store.crawlerScheduler = crawlers.NewCrawlerScheduler(conf.CrawlerInterval)
	store.snapshotScheduler = persistence.NewSnapshotScheduler(conf.SnapshotInterval)
	store.appMetrics = monitor.NewAppMetrics(time.Duration(store.Conf.AppMetricInterval), true)
}

func (baseStore *Store) registerHandlers() map[string]interface{} {
	return map[string]interface{} {
		STORE_GET: baseStore.Cache.Get,
		STORE_PUT: baseStore.Cache.Put,
		STORE_ADD: baseStore.Cache.Add,
		STORE_DELETE: baseStore.Cache.Delete,
		STORE_FLUSH: baseStore.Cache.Flush,
		STORE_NODE_SIZE: baseStore.Cache.CountKeys,
	}
}

func (store *Store) BuildStoreFromSnapshot(bs *[]byte) {
	// FUTURE: Switch to handle building for specified Cache types
	c, _ := persistence.BuildCacheFromSnapshot(bs)
	store.Cache = &(c)
}

func (store *Store) BuildStoreFromAof() {
	maxAofByteSize := store.Conf.AofMaxBytes
	persistence.RebootAof(&store.Cache, maxAofByteSize)
}

func (store *Store) RunStore() {
	go crawlers.StartCrawlers(&store.Cache, store.crawlerScheduler)
	if store.Conf.SnapshotEnabled {
		go persistence.StartSnapshotter(&store.Cache, &store.Conf, store.snapshotScheduler)
	} else if store.Conf.PersistenceAOF {
		if ok, _ := persistence.AofExists(); ok {
			go persistence.RebootAof(&store.Cache, store.Conf.AofMaxBytes)
		} else {
			go persistence.BootAOF(&store.Cache, store.Conf.AofMaxBytes)
		}
	}
}

func (store *Store) StopStore() {
	go crawlers.StopScheduler(store.crawlerScheduler)
	if store.Conf.SnapshotEnabled {
		go persistence.StopSnapshotter(store.snapshotScheduler)
	}
}

func (store *Store) newCacheFromPolicy(policy string) cache.Cache {
	switch policy {
	case LRU_TYPE:
		return lru.NewLRU(store.Conf)
	default:
		return nil
	}
}
