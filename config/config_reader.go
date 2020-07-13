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

package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	ConfigFileLinux = "/etc/ghostdb/ghostdbConf.json"

	DefaultKeyspaceSize      = 65536
	DefaultSysMetricInterval = 300  // 5 minutes
	DefaultAppMetricInterval = 300  // 5 minutes
	DefaultTTL               = -1   // Never Expire
	DefaultCrawlerInterval   = 300  // 5 minutes
	DefaultSnapshotInterval  = 3600 // 1 hour
	DefaultSnapshotEnabled   = true
	DefaultAofPersistence    = false
	DefaultAofMaxBytes       = 50000000 // 50MB
	DefaultEntryTimestamp    = true     // Enable timestamps in appMetrics logs
	DefaultEnableEncryption  = true
	DefaultPassphrase        = "SUPPLY_ME"
)

type Configuration struct {
	// KeyspaceSize represents the maximum number
	// of key-value pairs allowed in the cache
	KeyspaceSize int32

	// SysMetricInterval represents the interval
	// in seconds of when the sys metrics logging subsystem
	// should output to the log file.
	SysMetricInterval int32

	// AppMetricInterval represents the interval
	// in seconds of when the app metrics logging subsystem
	// should output to the log file.
	AppMetricInterval int32

	// DefaultTTL represents the default time-to-live
	// for key-value pairs in the cache. If set to -1
	// the key-value pair will never expire
	DefaultTTL int32

	// CrawlerInterval is the time, in seconds, of how
	// often the crawlers should crawl the cache to remove
	// stale key-value pairs.
	CrawlerInterval int32

	// SnapshotInterval is the time, in seconds, of how
	// often the snapshot scheduler should create a
	// point-in-time snapshot of the cache.
	SnapshotInterval int32

	// Snapshot enabled is a bool which enables snapshotting.
	// It takes precedence over AOF persistence.
	SnapshotEnabled bool

	// PersistenceAOF is a bool which enables AOF persistence.
	// Snapshots take precedence if both are enabled.
	PersistenceAOF bool

	// AOFMaxBytes represents the maximum number of bytes to be
	// written to the append-only-file log before an AOF log
	// rotation takes place.
	AofMaxBytes int64

	// EntryTimestamp is a boolean which enables timestamping for
	// entries in the appMetrics logging subsystem.
	EntryTimestamp bool

	// EnableEncryption is a bool that enables encryption of
	// snapshots using 128-bit AES.
	EnableEncryption bool

	// Passphrase is the passphrase to be used for snapshot encryption
	// should it be enabled.
	Passphrase string
}

// InitializeConfiguration initializes the cache configuration object
func InitializeConfiguration() Configuration {
	config, err := InitializeFromConfig()
	if err != nil {
		config.SetDefaultParams()
	}
	return config
}

// SetDefaultParams will set default parameters for the cache configuration
// object if the initializer was unable to initialize from a config file.
func (conf *Configuration) SetDefaultParams() {
	conf.KeyspaceSize = DefaultKeyspaceSize
	conf.SysMetricInterval = DefaultSysMetricInterval
	conf.AppMetricInterval = DefaultAppMetricInterval
	conf.DefaultTTL = DefaultTTL
	conf.CrawlerInterval = DefaultCrawlerInterval
	conf.SnapshotInterval = DefaultSnapshotInterval
	conf.SnapshotEnabled = DefaultSnapshotEnabled
	conf.PersistenceAOF = DefaultAofPersistence
	conf.AofMaxBytes = DefaultAofMaxBytes
	conf.EntryTimestamp = DefaultEntryTimestamp
	conf.EnableEncryption = DefaultEnableEncryption
	conf.Passphrase = DefaultPassphrase
}

// InitializeFromConfig initializes a configuration object from
// a configuration file.
func InitializeFromConfig() (Configuration, error) {
	var config Configuration

	file, err := ioutil.ReadFile(ConfigFileLinux)
	if err != nil {
		return Configuration{}, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return Configuration{}, err
	}

	return config, nil
}
