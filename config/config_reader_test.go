package config

import (
	"testing"

	"github.com/ghostdb/ghostdb-cache-node/utils"
)

func TestInitializeConfiguration(t *testing.T) {
	conf := InitializeConfiguration()
	utils.AssertEqual(t, conf.KeyspaceSize, int32(65536), "")
	utils.AssertEqual(t, conf.SysMetricInterval, int32(300), "")
	utils.AssertEqual(t, conf.AppMetricInterval, int32(300), "")
	utils.AssertEqual(t, conf.DefaultTTL, int32(-1), "")
	utils.AssertEqual(t, conf.CrawlerInterval, int32(300), "")
	utils.AssertEqual(t, conf.SnapshotInterval, int32(3600), "")
	utils.AssertEqual(t, conf.PersistenceAOF, false, "")
	utils.AssertEqual(t, conf.EntryTimestamp, true, "")
	utils.AssertEqual(t, conf.EnableEncryption, true, "")
	utils.AssertEqual(t, conf.Passphrase, "SUPPLY_ME", "")
}
