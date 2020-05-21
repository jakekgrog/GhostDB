package scheduler

import (
	"time"

	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
	"github.com/ghostdb/ghostdb-cache-node/cache/persistence/snapshot"
)

// SnapshotScheduler represents a scheduler for the cache snapshotter
type SnapshotScheduler struct {
	Interval time.Duration
	stop     chan bool
}

// NewSnapshotScheduler initializes a new Snapshot Scheduler
func NewSnapshotScheduler(interval int32) *SnapshotScheduler {
	scheduler := &SnapshotScheduler{
		Interval: time.Duration(interval) * time.Second,
		stop:     make(chan bool),
	}
	return scheduler
}

// StartSnapshotter starts the snapshot scheduler
func StartSnapshotter(cache *lru_cache.LRUCache, scheduler *SnapshotScheduler) {
	ticker := time.NewTicker(scheduler.Interval)

	for {
		select {
		case <-ticker.C:
			go snapshot.CreateSnapshot(cache, cache.Config.EnableEncryption, cache.Config.Passphrase)
		case <-scheduler.stop:
			ticker.Stop()
			return
		}
	}
}

// StopSnapshotter stops the Snapshotter by passing
// a bool to the scheduler channel
func StopSnapshotter(scheduler *SnapshotScheduler) {
	scheduler.stop <- true
}
