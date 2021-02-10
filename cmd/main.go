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

package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
	"os/signal"
	"syscall"
	"time"
	"net/http"
	"bytes"
	"encoding/json"

	"github.com/ghostdb/ghostdb-cache-node/server"
	"github.com/ghostdb/ghostdb-cache-node/store/base"
	"github.com/ghostdb/ghostdb-cache-node/store/persistence"
	"github.com/ghostdb/ghostdb-cache-node/store/monitor"
	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/system_monitor"
)

const (
	DefaultRaftAddr = ":11000"
	DefaultHTTPAddr = ":7991"
	retainSnapshotCount = 2
	raftTimeout = 10 * time.Second
)

var (
	httpAddr string
	raftAddr string
	joinAddr string
	nodeID   string
)

// Node configuration file
var conf config.Configuration

// Main cache object
var store *base.Store

// Schedulers
var sysMetricsScheduler *system_monitor.SysMetricsScheduler

func init() {
	flag.StringVar(&httpAddr, "http", DefaultHTTPAddr, "Set HTTP bind address")
	flag.StringVar(&raftAddr, "raft", DefaultRaftAddr, "Set Raft bind address")
	flag.StringVar(&joinAddr, "join", "", "Set join address, if any")
	flag.StringVar(&nodeID, "id", "", "Node ID")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <raft-data-path> \n", os.Args[0])
		flag.PrintDefaults()
	}

	conf = config.InitializeConfiguration()

	usr, _ := user.Current()
	configPath := usr.HomeDir

	err := os.Mkdir(configPath+"/ghostdb", 0777)
	if err != nil {
		log.Printf("Failed to create GhostDB configuration directory")
	}

	// Create sysMetrics and appMetrics logfiles if they do not exist
	sysMetricsFile, err := os.OpenFile(configPath + system_monitor.SysMetricsLogFilename, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("Failed to create or read sysMetrics log file: %s", err.Error())
		panic(err)
	}
	defer sysMetricsFile.Close()

	appMetricsFile, err := os.OpenFile(configPath + monitor.AppMetricsLogFilePath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to create or read appMetrics log file: %s", err.Error())
	}
	defer appMetricsFile.Close()

	sysMetricsScheduler = system_monitor.NewSysMetricsScheduler(conf.SysMetricInterval)
}

func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "No Raft storage directory specified\n")
		os.Exit(1)
	}

	// Ensure Raft sotrage exists
	raftDir := flag.Arg(0)
	if raftDir == "" {
		fmt.Fprintf(os.Stderr, "No Raft storage director specified\n")
		os.Exit(1)
	}
	os.MkdirAll(raftDir, 0700)

	go system_monitor.StartSysMetrics(sysMetricsScheduler)
	log.Println("successfully started sysMetrics monitor...")

	store = base.NewStore("LRU") // FUTURE: Read store type from config
	store.BuildStore(conf)

	store.RaftDir = raftDir
	store.RaftBind = raftAddr
	if err := store.Open(joinAddr == "", nodeID); err != nil {
		log.Fatalf("failed to open store: %s", err.Error())
	} 

	usr, _ := user.Current()
	configPath := usr.HomeDir
	// Build the cache from a snapshot if snaps enabled.
	// If the snapshot does not exist, then build a new cache.
	if conf.SnapshotEnabled {
		if _, err := os.Stat(configPath + persistence.GetSnapshotFilename()); err == nil {
			bytes := persistence.ReadSnapshot(conf.EnableEncryption, conf.Passphrase)
			store.BuildStoreFromSnapshot(bytes)
			log.Println("successfully booted from snapshot...")
		} else { 
			log.Println("successfully booted new cache...")
		}
	} else {
		if conf.PersistenceAOF {
			if ok, _ := persistence.AofExists(); ok {
				persistence.RebootAof(&store.Cache, conf.AofMaxBytes)
			}
			store.BuildStoreFromAof()
			persistence.BootAOF(&store.Cache, conf.AofMaxBytes)
			log.Println("successfully booted from AOF...")
		}
		log.Println("successfully booted new cache...")
	}

	store.RunStore()

	log.Println("Starting store...")

	sysMetricsScheduler = system_monitor.NewSysMetricsScheduler(conf.SysMetricInterval)

	service := server.NewService(httpAddr, store)
	go service.Start()

	log.Println("Starting service...")

	if joinAddr != "" {
		if err := join(joinAddr, raftAddr, nodeID); err != nil {
			log.Fatalf("failed to join node at %s: %s", joinAddr, err.Error())
		}
	}

	log.Println("started successfully ...")

	t := make(chan os.Signal)
	signal.Notify(t, os.Interrupt, syscall.SIGTERM)
	<-t

	log.Println("exiting ...")
	
}

func join(joinAddr, raftAddr, nodeID string) error {
	b, err := json.Marshal(map[string]string{"addr": raftAddr, "id": nodeID})
	if err != nil {
		return err
	}
	resp, err := http.Post(fmt.Sprintf("http://%s/join", joinAddr), "application-type/json", bytes.NewReader(b))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}