package base

import (
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/utils"
	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
)

func TestSnapshotScheduler(t *testing.T) {
	var config config.Configuration = config.InitializeConfiguration()
	config.SnapshotInterval = int32(5)
	var store *Store
	store = NewStore("LRU")
	store.BuildStore(config)
	store.RunStore()

	store.Execute("put", request.NewRequestFromValues("England", "London", 2))
	store.Execute("put", request.NewRequestFromValues("Italy", "Rome", -1))

	utils.AssertEqual(t, store.Execute("nodeSize", request.NewEmptyRequest()).Gobj.Value, int32(2), "")

	time.Sleep(6 * time.Second)

	store.StopStore()
}
