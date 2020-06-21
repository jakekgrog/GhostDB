package main

import (
	"log"
	"os"
	"os/user"
	"os/signal"
	"syscall"

	"github.com/ghostdb/ghostdb-cache-node/server/ghost_http"
	"github.com/ghostdb/ghostdb-cache-node/store/base"
	"github.com/ghostdb/ghostdb-cache-node/store/persistence"
	"github.com/ghostdb/ghostdb-cache-node/store/monitor"
	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/system_monitor"
)

// Node configuration file
var conf config.Configuration

// Main cache object
var store *base.Store

// Schedulers
//var crawlerScheduler *lru.CrawlerScheduler
var snitchScheduler *system_monitor.SnitchScheduler
//var snapshotScheduler *lru.SnapshotScheduler

func init() {
	conf = config.InitializeConfiguration()

	usr, _ := user.Current()
	configPath := usr.HomeDir
	log.Println("LOG PATH: "+configPath)

	err := os.Mkdir(configPath+"/ghostdb", 0777)
	if err != nil {
		log.Fatalf("Failed to create GhostDB configuration directory")
	}

	// Create snitch and watchdog logfiles if they do not exist
	snitchFile, err := os.OpenFile(configPath + system_monitor.SnitchLogFileName, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0777)
	if err != nil {
		log.Fatalf("Failed to create or read snitch log file: %s", err.Error())
		panic(err)
	}
	defer snitchFile.Close()

	watchdogFile, err := os.OpenFile(configPath + monitor.WatchDogLogFilePath, os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("Failed to create or read watchdog log file: %s", err.Error())
	}
	defer watchdogFile.Close()


	store = base.NewStore("LRU") // FUTURE: Read store type from config
	store.BuildStore(conf)


	// Build the cache from a snapshot if snaps enabled.
	// If the snapshot does not exist, then build a new cache.
	if conf.SnapshotEnabled {
		if _, err := os.Stat(configPath + system_monitor.SnitchLogFileName); err == nil {
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

	// crawlerScheduler = lru.NewCrawlerScheduler(conf.CrawlerInterval)
	snitchScheduler = system_monitor.NewSnitchScheduler(conf.SnitchMetricInterval)
	// snapshotScheduler = lru.NewSnapshotScheduler(conf.SnapshotInterval)
}

func main() {
	// go lru.StartCrawlers(store.Cache, crawlerScheduler)
	// log.Println("successfully started crawler lru...")
	go system_monitor.StartSnitch(snitchScheduler)
	log.Println("successfully started snitch monitor...")
	// if conf.SnapshotEnabled {
	// 	go lru.StartSnapshotter(store.Cache, snapshotScheduler)
	// 	log.Println("successfully started snapshot lru...")
	// }
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
