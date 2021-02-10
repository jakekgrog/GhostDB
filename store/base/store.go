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

	"fmt"
	"net"
	"os"
	"regexp"
	"strings"
	"encoding/json"
	
	"github.com/hashicorp/raft"
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

const (
	retainSnapshotCount = 2
	raftTimeout = 10*time.Second
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
	// Join joins a node, identified by nodeID and located at addr, to this store.
	// The node must be ready to respond to Raft communications at that address.
	Join(nodeID string, addr string) error
	// Open opens the store. If enableSingle is set, and there are no existing peers,
	// then this node becomes the first node, and therefore leader, of the cluster.
	// localID should be the server identifier for this node.
	Open(enableSingle bool, localID string) error
}

type Store struct {
	RaftDir            string
	RaftBind           string
	Raft               *raft.Raft
	ServerID           string
	NumericalID        int
	PeersLength        int

	Conf               config.Configuration
	policy             string
	Cache              cache.Cache
	commands           map[string]interface{}
	crawlerScheduler   *crawlers.CrawlerScheduler
	snapshotScheduler  *persistence.SnapshotScheduler
	appMetrics         *monitor.AppMetrics
}

// Command is the struct used by the replication log.
// All write commands can be written to the replciation log
// in this format
type Command struct {
	Cmd  string
	Args request.CacheRequest
}

func NewStore(policy string) *Store {
	return &Store{
		policy: policy,
	}
}

func (store *Store) Execute(cmd string, args request.CacheRequest) response.CacheResponse {
	// All commands that are not write commands don't need to call Apply() on the store.
	// We can handle them as before.
	if cmd == "get" || cmd == "getAppMetrics"{
		// Handle get
		if cmd == "get" {
			if _, ok := store.commands[cmd]; !ok {
				return response.BadCommandResponse(cmd)
			}
			
			execResult := store.commands[cmd].(func(request.CacheRequest) response.CacheResponse)(args)
			if store.Conf.PersistenceAOF {
				writeAof(cmd, &(args))
			}
			return execResult
		}
		// Handle getAppMetrics
		return response.BadCommandResponse(cmd)
	} else {
		// All write commands need to be applied to the replication log
		c := &Command{
			Cmd: cmd,
			Args: args,
		}

		b, err := json.Marshal(c)
		if err != nil {
			return response.BadCommandResponse(cmd)
		}

		applyFuture := store.Raft.Apply(b, raftTimeout)
		if err := applyFuture.Error(); err != nil {
			return response.NewResponseFromMessage("Error commiting to raft cluster", 500)
		}

		res, ok := applyFuture.Response().(response.CacheResponse)
		if !ok {
			return response.NewResponseFromMessage("Error commiting to raft cluster 2", 500)
		}

		return res
	}
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

func (store *Store) Join(nodeID string, addr string) error {
	fmt.Printf("received join request for remote node %s at %s\n", nodeID, addr)

	configFuture := store.Raft.GetConfiguration()
	if err := configFuture.Error(); err != nil {
		fmt.Printf("failed to get raft configuration: %v", err)
		return err
	}

	for _, srv := range configFuture.Configuration().Servers {
		// If a node already exists with either the joining node's ID or address,
		// that node may need to be removed from the config first.
		if srv.ID == raft.ServerID(nodeID) || srv.Address == raft.ServerAddress(addr) {
			// However if *both* the ID and the address are the same, then nothing -- not even
			// a join operation -- is needed.
			if srv.Address == raft.ServerAddress(addr) && srv.ID == raft.ServerID(nodeID) {
				fmt.Printf("node %s at %s already member of cluster, ignoring join request", nodeID, addr)
				return nil
			}

			future := store.Raft.RemoveServer(srv.ID, 0, 0)
			if err := future.Error(); err != nil {
				return fmt.Errorf("error removing existing node %s at %s: %s", nodeID, addr, err)
			}
		}
	}

	// Add a voter. The voters will decide who is master when a new leader election is called.
	f := store.Raft.AddVoter(raft.ServerID(nodeID), raft.ServerAddress(addr), 0, 0)
	if f.Error() != nil {
		return f.Error()
	}
	fmt.Printf("node %s at %s joined successfully", nodeID, addr)
	return nil
}

func (store *Store) Open(enableSingle bool, localID string) error {
	config := raft.DefaultConfig()
	config.LocalID = raft.ServerID(localID)
	store.ServerID = localID
	
	addr, err := net.ResolveTCPAddr("tcp", store.RaftBind)
	if err != nil {
		return err
	}
	
	transport, err := raft.NewTCPTransport(store.RaftBind, addr, 3, raftTimeout, os.Stderr)
	if err != nil {
		return err
	}

	snapshots, err := raft.NewFileSnapshotStore(store.RaftDir, retainSnapshotCount, os.Stderr)
	if err != nil {
		return fmt.Errorf("file snapshot store: %s", err)
	}

	var logStore raft.LogStore
	var stableStore raft.StableStore
	logStore = raft.NewInmemStore()
	stableStore = raft.NewInmemStore()

	ra, err := raft.NewRaft(config, (*fsm)(store), logStore, stableStore, snapshots, transport)
	if err != nil {
		return fmt.Errorf("new raft: %s", err)
	}
	store.Raft = ra

	if enableSingle {
		configuration := raft.Configuration{
			Servers: []raft.Server{
				{
					ID:      config.LocalID,
					Address: transport.LocalAddr(),
				},
			},
		}
		ra.BootstrapCluster(configuration)
	}

	return nil
}

type fsm Store

// GetNumericalID is used to get the numerical ID of a node from the list of peers
// parameters: ID (a string identifier of the node), peers (an array of current nodes in the cluster)
// returns: int (the numeric value of a node's indentifier), otherwise -1
func GetNumericalID(ID string, peers []string) int {
	for i, value := range peers {
		if value == ID {
			return i
		}
	}
	return -1
}


func PeersList(rawConfig string) []string {
	peers := []string{}
	re := regexp.MustCompile(`ID:[0-9A-z]*`)
	for _, peer := range re.FindAllString(rawConfig, -1) {
		peers = append(peers, strings.Replace(peer, "ID:", "", -1))
	}
	return peers
}

// Apply applies a Raft log entry to the key-value store.
func (f *fsm) Apply(l *raft.Log) interface{} {
	// Handle all other commands
	var c Command
	
	if err := json.Unmarshal(l.Data, &c); err != nil {
		panic(fmt.Sprintf("failed to unmarshal command: %s", err.Error()))
	}

	if _, ok := f.commands[c.Cmd]; !ok {
		return response.BadCommandResponse(c.Cmd)
	}
	
	execResult := f.commands[c.Cmd].(func(request.CacheRequest) response.CacheResponse)(c.Args)
	if f.Conf.PersistenceAOF {
		writeAof(c.Cmd, &(c.Args))
	}
	// CHECK RESPONSE AND SEND TO APP METRICS
	monitor.WriteMetrics(f.appMetrics, c.Cmd, execResult)
	
	return execResult
}
