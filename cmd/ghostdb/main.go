package main

import (
	"os"

	"github.com/ghostdb/ghostdb-cache-node/api"
	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
	"github.com/ghostdb/ghostdb-cache-node/cache/scheduler"
	"github.com/ghostdb/ghostdb-cache-node/internal/ghost_config"
	"github.com/ghostdb/ghostdb-cache-node/internal/monitoring/snitch"
	"github.com/ghostdb/ghostdb-cache-node/internal/monitoring/watchDog"
	"github.com/ghostdb/ghostdb-cache-node/cache/persistence/snapshot"
	"github.com/ghostdb/ghostdb-cache-node/cache/persistence/append_only_file"
)

// Node configuration file
var config ghost_config.Configuration

// Main cache object
var cache *lru_cache.LRUCache

// Schedulers
var crawlerScheduler *scheduler.CrawlerScheduler
var snitchScheduler *snitch.SnitchScheduler
var snapshotScheduler *scheduler.SnapshotScheduler

func init() {
	config = ghost_config.InitializeConfiguration()

	configPath, _ := os.UserConfigDir()

	// Create snitch and watchdog logfiles if they do not exist
	snitchFile, err := os.OpenFile(configPath+snitch.SnitchLogFileName, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer snitchFile.Close()

	watchdogFile, err := os.OpenFile(configPath + watchDog.WatchDogLogFilePath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		panic(err)
	}
	defer watchdogFile.Close()

	// Build the cache from a snapshot if snaps enabled.
	// If the snapshot does not exist, then build a new cache.
	if config.SnapshotEnabled {
		if _, err := os.Stat(configPath + snapshot.SnitchLogFilename); err == nil {
			bytes := snapshot.ReadSnapshot(config.EnableEncryption, config.Passphrase)
			cache, _ = snapshot.BuildCache(bytes)
		} else {
			cache = lru_cache.NewLRU(config) 
		}
	} else {
		cache = lru_cache.NewLRU(config)
		if config.PersistenceAOF {
			append_only_file.BootAOF(cache, config.AofMaxBytes)
		}
	}

	crawlerScheduler = scheduler.NewCrawlerScheduler(config.CrawlerInterval)
	snitchScheduler = snitch.NewSnitchScheduler(config.SnitchMetricInterval)
	snapshotScheduler = scheduler.NewSnapshotScheduler(config.SnapshotInterval)
}

func main() {
	go scheduler.StartCrawlers(cache, crawlerScheduler)
	go snitch.StartSnitch(snitchScheduler)
	if config.SnapshotEnabled {
		go scheduler.StartSnapshotter(cache, snapshotScheduler)
	}
	api.NodeConfig(cache)
	api.Router()
}
