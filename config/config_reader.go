package config

import (
	"encoding/json"
	"io/ioutil"
)

const (
	CONFIG_FILE_LINUX = "/etc/ghostdb/ghostdbConf.json"

	DEFAULT_KEYSPACE_SIZE            = 65536
	DEFAULT_SNITCH_METRIC_INTERVAL   = 300  // 5 minutes
	DEFAULT_WATCHDOG_METRIC_INTERVAL = 300  // 5 minutes
	DEFAULT_TTL                      = -1   // Never Expire
	DEFAULT_CRAWLER_INTERVAL         = 300  // 5 minutes
	DEFAULT_SNAPSHOT_INTERVAL        = 3600 // 1 hour
	DEFAULT_SNAPSHOT_ENABLED         = true
	DEFAULT_AOF_PERSISTENCE          = false
	DEFAULT_AOF_MAX_BYTES            = 50000000 // 50MB
	DEFAULT_ENTRY_TIMESTAMP          = true  // Enable timestamps in watchdog logs
	DEFAULT_ENABLE_ENCRYPTION        = true
	DEFAULT_PASSPHRASE               = "SUPPLY_ME"
)

type Configuration struct {
	// KeyspaceSize represents the maximum number
	// of key-value pairs allowed in the cache
	KeyspaceSize           int32

	// SnitchMetricInterval represents the interval
	// in seconds of when the snitch logging subsystem
	// should output to the log file.
	SnitchMetricInterval   int32

	// WatchdogMetricInterval represents the interval
	// in seconds of when the watchdog logging subsystem
	// should output to the log file.
	WatchdogMetricInterval int32

	// DefaultTTL represents the default time-to-live
	// for key-value pairs in the cache. If set to -1
	// the key-value pair will never expire
	DefaultTTL             int32

	// CrawlerInterval is the time, in seconds, of how
	// often the crawlers should crawl the cache to remove
	// stale key-value pairs.
	CrawlerInterval        int32

	// SnapshotInterval is the time, in seconds, of how
	// often the snapshot scheduler should create a
	// point-in-time snapshot of the cache.
	SnapshotInterval       int32

	// Snapshot enabled is a bool which enables snapshotting.
	// It takes precedence over AOF persistence.
	SnapshotEnabled        bool

	// PersistenceAOF is a bool which enables AOF persistence.
	// Snapshots take precedence if both are enabled.
	PersistenceAOF         bool

	// AOFMaxBytes represents the maximum number of bytes to be
	// written to the append-only-file log before an AOF log
	// rotation takes place.
	AofMaxBytes            int64

	// EntryTimestamp is a boolean which enables timestamping for
	// entries in the watchdog logging subsystem.
	EntryTimestamp         bool

	// EnableEncryption is a bool that enables encryption of
	// snapshots using 128-bit AES.
	EnableEncryption       bool

	// Passphrase is the passphrase to be used for snapshot encryption
	// should it be enabled.
	Passphrase             string
}

// InitializeConfiguration initializes the cache configuration object
func InitializeConfiguration() Configuration {
	var config, err = InitializeFromConfig()
	if err != nil {
		config.SetDefaultParams()
	}
	return config
}

// SetDefaultParams will set default parameters for the cache configuration
// object if the initializer was unable to initialize from a config file.
func (conf *Configuration) SetDefaultParams() {
	conf.KeyspaceSize = DEFAULT_KEYSPACE_SIZE
	conf.SnitchMetricInterval = DEFAULT_SNITCH_METRIC_INTERVAL
	conf.WatchdogMetricInterval = DEFAULT_WATCHDOG_METRIC_INTERVAL
	conf.DefaultTTL = DEFAULT_TTL
	conf.CrawlerInterval = DEFAULT_CRAWLER_INTERVAL
	conf.SnapshotInterval = DEFAULT_SNAPSHOT_INTERVAL
	conf.SnapshotEnabled = DEFAULT_SNAPSHOT_ENABLED
	conf.PersistenceAOF = DEFAULT_AOF_PERSISTENCE
	conf.AofMaxBytes = DEFAULT_AOF_MAX_BYTES
	conf.EntryTimestamp = DEFAULT_ENTRY_TIMESTAMP
	conf.EnableEncryption = DEFAULT_ENABLE_ENCRYPTION
	conf.Passphrase = DEFAULT_PASSPHRASE
}

// InitializeFromConfig initializes a configuration object from
// a configuration file.
func InitializeFromConfig() (Configuration, error) {
	var config Configuration

	file, err := ioutil.ReadFile(CONFIG_FILE_LINUX)
	if err != nil {
		return Configuration{}, err
	}

	err = json.Unmarshal(file, &config)
	if err != nil {
		return Configuration{}, err
	}

	return config, nil
}
