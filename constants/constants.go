package constants

const (
	// STORE COMMANDS
	STORE_GET = "get"
	STORE_PUT = "put"
	STORE_ADD = "add"
	STORE_DELETE = "delete"
	STORE_FLUSH = "flush"
	STORE_NODE_SIZE = "nodeSize"
	STORE_APP_METRICS = "getAppMetrics"

	// STORE POLICY TYPES
	LRU_TYPE = "LRU"   // Least recently used
	LFU_TYPE = "LFU"   // Least frequenty used
	MRU_TYPE = "MRU"   // Most recently used
	ARC_TYPE = "ARC"   // Adaptive Replacement Cache
	TLRU_TYPE = "TLRU" // Time-aware Least Recently Used
)