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
	"log"
	"os"
	"os/signal"
	"os/user"
	"syscall"

	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/server/ghost_http"
	"github.com/ghostdb/ghostdb-cache-node/store/base"
	"github.com/ghostdb/ghostdb-cache-node/store/monitor"
	"github.com/ghostdb/ghostdb-cache-node/store/persistence"
	"github.com/ghostdb/ghostdb-cache-node/system_monitor"
)

// Node configuration file
var conf config.Configuration

// Main cache object
var store *base.Store

// Schedulers
var sysMetricsScheduler *system_monitor.SysMetricsScheduler

func init() {
	conf = config.InitializeConfiguration()

	usr, _ := user.Current()
	configPath := usr.HomeDir
	log.Println("LOG PATH: " + configPath)

	err := os.Mkdir(configPath+"/ghostdb", 0o777)
	if err != nil {
		log.Printf("Failed to create GhostDB configuration directory")
	}

	// Create sysMetrics and appMetrics logfiles if they do not exist
	sysMetricsFile, err := os.OpenFile(configPath+system_monitor.SysMetricsLogFilename, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o777)
	if err != nil {
		log.Fatalf("Failed to create or read sysMetrics log file: %s", err.Error())
		panic(err)
	}
	defer sysMetricsFile.Close()

	appMetricsFile, err := os.OpenFile(configPath+monitor.AppMetricsLogFilePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o644)
	if err != nil {
		log.Fatalf("Failed to create or read appMetrics log file: %s", err.Error())
	}
	defer appMetricsFile.Close()

	store = base.NewStore("LRU") // FUTURE: Read store type from config
	store.BuildStore(conf)

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

	sysMetricsScheduler = system_monitor.NewSysMetricsScheduler(conf.SysMetricInterval)
}

func main() {
	go system_monitor.StartSysMetrics(sysMetricsScheduler)
	log.Println("successfully started sysMetrics monitor...")
	ghost_http.NodeConfig(store)
	log.Println("successfully started GhostDB Node server...")
	log.Println("GhostDB started successfully...")

	t := make(chan os.Signal)
	signal.Notify(t, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-t
		log.Println("exiting...")
		os.Exit(1)
	}()

	ghost_http.Router()
}
