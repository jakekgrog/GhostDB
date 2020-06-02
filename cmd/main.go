package main

import (
	"log"
	"os"
	"os/user"
	"os/signal"
	"syscall"

	"github.com/ghostdb/ghostdb-cache-node/server/ghost_http"
	"github.com/ghostdb/ghostdb-cache-node/store/lru"
	"github.com/ghostdb/ghostdb-cache-node/config"
)

// Node configuration file
var conf config.Configuration

// Main cache object
var cache *lru.LRUCache

// Schedulers
var crawlerScheduler *lru.CrawlerScheduler
var snitchScheduler *lru.SnitchScheduler
var snapshotScheduler *lru.SnapshotScheduler

func init() {
	conf = config.InitializeConfiguration()

	usr, _ := user.Current()
	configPath := usr.HomeDir
	log.Println("LOG PATH: "+configPath)

	err := os.Mkdir(configPath+"/ghostdb", 0777)

	// Create snitch and watchdog logfiles if they do not exist
	snitchFile, err := os.OpenFile(configPath + lru.SnitchLogFileName, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("Failed to create or read snitch log file: %s", err.Error())
		panic(err)
	}
	defer snitchFile.Close()

	watchdogFile, err := os.OpenFile(configPath + lru.WatchDogLogFilePath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to create or read watchdog log file: %s", err.Error())
	}
	defer watchdogFile.Close()

	// Build the cache from a snapshot if snaps enabled.
	// If the snapshot does not exist, then build a new cache.
	if conf.SnapshotEnabled {
		if _, err := os.Stat(configPath + lru.SnitchLogFilename); err == nil {
			bytes := lru.ReadSnapshot(conf.EnableEncryption, conf.Passphrase)
			cache, _ = lru.BuildCacheFromSnapshot(bytes)
			log.Println("successfully booted from snapshot...")
		} else {
			cache = lru.NewLRU(conf) 
			log.Println("successfully booted new cache...")
		}
	} else {
		cache = lru.NewLRU(conf)
		if conf.PersistenceAOF {
			lru.BootAOF(cache, conf.AofMaxBytes)
			log.Println("successfully booted from AOF...")
		}
		log.Println("successfully booted new cache...")
	}

	crawlerScheduler = lru.NewCrawlerScheduler(conf.CrawlerInterval)
	snitchScheduler = lru.NewSnitchScheduler(conf.SnitchMetricInterval)
	snapshotScheduler = lru.NewSnapshotScheduler(conf.SnapshotInterval)
}

func main() {
	go lru.StartCrawlers(cache, crawlerScheduler)
	log.Println("successfully started crawler lru...")
	go lru.StartSnitch(snitchScheduler)
	log.Println("successfully started snitch monitor...")
	if conf.SnapshotEnabled {
		go lru.StartSnapshotter(cache, snapshotScheduler)
		log.Println("successfully started snapshot lru...")
	}
	ghost_http.NodeConfig(cache)
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
