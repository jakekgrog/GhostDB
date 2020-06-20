package base

import (
	"testing"

	// "github.com/ghostdb/ghostdb-cache-node/store/lru"
	// "github.com/ghostdb/ghostdb-cache-node/utils"
	//"github.com/ghostdb/ghostdb-cache-node/config"
	// "github.com/ghostdb/ghostdb-cache-node/store/request"
)

func TestSerializer(t *testing.T) {
	// var config config.Configuration = config.InitializeConfiguration()

	// var store *Store
	// store = NewStore("LRU")

	// cache.Put(request.NewRequestFromValues("Italy", "Rome", -1))
	// cache.Put(request.NewRequestFromValues("England", "London", 2))

	// encryptionEnabled := config.EnableEncryption
	// passphrase := "SUPPLY_PASSPHRASE"

	// CreateSnapshot(*cache, &config)

	// bytes := ReadSnapshot(encryptionEnabled, passphrase)

	// c, err := BuildCacheFromSnapshot(bytes)

	// if err != nil {
	// 	panic(err)
	// }

	// val := c.Get(request.NewRequestFromValues("England", "", -1))

	// utils.AssertEqual(t, val.Gobj.Key, "London", "")
	// utils.AssertEqual(t, c.Size, int32(65536), "")
	// utils.AssertEqual(t, c.Full, false, "")
	// utils.AssertEqual(t, c.CountKeys(request.NewEmptyRequest()), int32(2), "")

	// // Test the config was rebuilt correctly.
	// utils.AssertEqual(t, c.Config.KeyspaceSize, int32(65536), "")
	// utils.AssertEqual(t, c.Config.SnitchMetricInterval, int32(300), "")
	// utils.AssertEqual(t, c.Config.WatchdogMetricInterval, int32(300), "")
	// utils.AssertEqual(t, c.Config.DefaultTTL, int32(-1), "")
	// utils.AssertEqual(t, c.Config.CrawlerInterval, int32(300), "")
	// utils.AssertEqual(t, c.Config.SnapshotInterval, int32(3600), "")
	// utils.AssertEqual(t, c.Config.PersistenceAOF, false, "")
	// utils.AssertEqual(t, c.Config.EntryTimestamp, true, "")
	// utils.AssertEqual(t, c.Config.EnableEncryption, true, "")
}
