package persistence

import (
	"time"

	"github.com/ghostdb/ghostdb-cache-node/store/cache"
	"github.com/ghostdb/ghostdb-cache-node/config"
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

/**
* StartSnapshotter will start the cache snapshot scheduler
* 
* Snapshots are periodically taken of the stores cache until the ticker
* is stopped.
* StartSnapshotter is generic in nature and will runthe appropriate
* snapshotter for the stores underlying cache policy implementation
*
*/
func StartSnapshotter(cache *cache.Cache, conf *config.Configuration, scheduler *SnapshotScheduler) {
	ticker := time.NewTicker(scheduler.Interval)

	for {
		select {
		case <-ticker.C:
			go CreateSnapshot(cache, conf)
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
