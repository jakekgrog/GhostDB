package constants

import (
	"testing"

	"github.com/ghostdb/ghostdb-cache-node/utils"
)

func TestConstants(t *testing.T) {
	utils.AssertEqual(t, STORE_GET, "get", "")
	utils.AssertEqual(t, STORE_PUT, "put", "")
	utils.AssertEqual(t, STORE_ADD, "add", "")
	utils.AssertEqual(t, STORE_DELETE, "delete", "")
	utils.AssertEqual(t, STORE_FLUSH, "flush", "")
	utils.AssertEqual(t, STORE_NODE_SIZE, "nodeSize", "")
	utils.AssertEqual(t, STORE_APP_METRICS, "getAppMetrics", "")

	utils.AssertEqual(t, LRU_TYPE, "LRU", "")
	utils.AssertEqual(t, LFU_TYPE, "LFU", "")
	utils.AssertEqual(t, MRU_TYPE, "MRU", "")
	utils.AssertEqual(t, ARC_TYPE, "ARC", "")
	utils.AssertEqual(t, TLRU_TYPE, "TLRU", "")
}