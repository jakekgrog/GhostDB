package linked_list

import (
	"testing"

	"github.com/ghostdb/ghostdb-cache-node/cache/utils"
)

func TestListOperations(t *testing.T) {
	dll := InitList()
	Insert(dll, "Ireland", "Dublin", -1)
	Insert(dll, "Italy", "Rome", -1)

	n1, _ := Insert(dll, "Germany", "Berlin", -1)
	utils.AssertEqual(t, n1.Key, "Germany", "")
	utils.AssertEqual(t, n1.Value, "Berlin", "")
	utils.AssertEqual(t, n1.TTL, int64(-1), "")

	n, _ := Insert(dll, "France", "Paris", -1)
	utils.AssertEqual(t, n.Key != "Paris", true, "")
	utils.AssertEqual(t, n.Value != "France", true, "")
	utils.AssertEqual(t, n.TTL, int64(-1), "")

	n, _ = Insert(dll, "Belgium", "Brussels", -1)
	utils.AssertEqual(t, n.Next.Key, "France", "")
	utils.AssertEqual(t, n.Prev.Key, "", "")

	n, _ = RemoveNode(dll, n1)
	utils.AssertEqual(t, n.Key, "Germany", "")
	utils.AssertEqual(t, n.Value, "Berlin", "")
	utils.AssertEqual(t, n.TTL, int64(-1), "")

	Insert(dll, n1.Key, n1.Value, -1)
}
