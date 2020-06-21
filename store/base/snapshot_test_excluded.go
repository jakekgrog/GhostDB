package base

import (
	"testing"

	"github.com/ghostdb/ghostdb-cache-node/utils"
	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/store/persistence"
)

func TestSerializer(t *testing.T) {
	var config config.Configuration = config.InitializeConfiguration()

	var store *Store
	store = NewStore("LRU")
	store.BuildStore(config)

	store.Execute("put", request.NewRequestFromValues("Italy", "Rome", -1))
	store.Execute("put", request.NewRequestFromValues("England", "London", 2))

	encryptionEnabled := config.EnableEncryption
	passphrase := "SUPPLY_PASSPHRASE"

	store.CreateSnapshot()

	bytes := persistence.ReadSnapshot(encryptionEnabled, passphrase)

	c, err := persistence.BuildCacheFromSnapshot(bytes)
	if err != nil {
		panic(err)
	}

	store.Cache = &c

	val := store.Execute("get", request.NewRequestFromValues("England", "", -1))

	utils.AssertEqual(t, val.Gobj.Key, "London", "")
	utils.AssertEqual(t, store.Execute("nodeSize", request.NewEmptyRequest()), int32(2), "")

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
