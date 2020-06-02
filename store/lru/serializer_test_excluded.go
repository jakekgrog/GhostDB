package lru

import (
	"testing"

	"github.com/ghostdb/ghostdb-cache-node/utils"
	"github.com/ghostdb/ghostdb-cache-node/config"
)

func TestSerializer(t *testing.T) {
	var config config.Configuration = config.InitializeConfiguration()

	var cache *LRUCache
	cache = NewLRU(config)

	cache.Put("Italy", "Rome", -1)
	cache.Put("England", "London", 2)

	encryptionEnabled := config.EnableEncryption
	passphrase := "SUPPLY_PASSPHRASE"

	CreateSnapshot(cache, encryptionEnabled, passphrase)

	bytes := ReadSnapshot(encryptionEnabled, passphrase)

	c, err := BuildCacheFromSnapshot(bytes)

	if err != nil {
		panic(err)
	}

	val := c.Get("England")

	utils.AssertEqual(t, val, "London", "")
	utils.AssertEqual(t, c.Size, int32(65536), "")
	utils.AssertEqual(t, c.Full, false, "")
	utils.AssertEqual(t, c.CountKeys(), int32(2), "")

	// Test the config was rebuilt correctly.
	utils.AssertEqual(t, c.Config.KeyspaceSize, int32(65536), "")
	utils.AssertEqual(t, c.Config.SnitchMetricInterval, int32(300), "")
	utils.AssertEqual(t, c.Config.WatchdogMetricInterval, int32(300), "")
	utils.AssertEqual(t, c.Config.DefaultTTL, int32(-1), "")
	utils.AssertEqual(t, c.Config.CrawlerInterval, int32(300), "")
	utils.AssertEqual(t, c.Config.SnapshotInterval, int32(3600), "")
	utils.AssertEqual(t, c.Config.PersistenceAOF, false, "")
	utils.AssertEqual(t, c.Config.EntryTimestamp, true, "")
	utils.AssertEqual(t, c.Config.EnableEncryption, true, "")
}
