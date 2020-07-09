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
