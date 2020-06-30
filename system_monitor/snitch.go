package system_monitor

import (
	"bufio"
	"encoding/json"
	"log"
	"os"
	"os/user"
	"runtime"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/utils"
	"github.com/ghostdb/ghostdb-cache-node/store/response"
)

const (
	SnitchLogFileName = "/ghostdb/ghostdb_snitch.log"
	SnitchTempFileName = "/ghostdb/ghostdb_snitch_tmp.log"
	MaxSnitchLogSize = 10000000
)

type Snitch struct {
	Timestamp    string

	// Alloc is bytes of allocated heap objects.
	Alloc        uint64

	// TotalAlloc is cumulative bytes allocated for heap objects.
	//
	// TotalAlloc increases as heap objects are allocated, but
	// unlike Alloc and HeapAlloc, it does not decrease when
	// objects are freed.
	TotalAlloc   uint64

	// Sys is the total bytes of memory obtained from the OS.
	//
	// Sys is the sum of the XSys fields below. Sys measures the
	// virtual address space reserved by the Go runtime for the
	// heap, stacks, and other internal data structures. It's
	// likely that not all of the virtual address space is backed
	// by physical memory at any given moment, though in general
	// it all was at some point.
	Sys          uint64

	// Lookups is the number of pointer lookups performed by the
	// runtime.
	//
	// This is primarily useful for debugging runtime internals.
	Lookups      uint64

	// Mallocs is the cumulative count of heap objects allocated.
	Mallocs      uint64

	// Frees is the cumulative count of heap objects freed.
	Frees        uint64

	// The number of live objects is Mallocs - Frees.
	LiveObjects  uint64

	// HeapAlloc is bytes of allocated heap objects.
	//
	// "Allocated" heap objects include all reachable objects, as
	// well as unreachable objects that the garbage collector has
	// not yet freed. Specifically, HeapAlloc increases as heap
	// objects are allocated and decreases as the heap is swept
	// and unreachable objects are freed. Sweeping occurs
	// incrementally between GC cycles, so these two processes
	// occur simultaneously, and as a result HeapAlloc tends to
	// change smoothly (in contrast with the sawtooth that is
	// typical of stop-the-world garbage collectors).
	HeapAlloc    uint64

	// HeapSys is bytes of heap memory obtained from the OS.
	//
	// HeapSys measures the amount of virtual address space
	// reserved for the heap. This includes virtual address space
	// that has been reserved but not yet used, which consumes no
	// physical memory, but tends to be small, as well as virtual
	// address space for which the physical memory has been
	// returned to the OS after it became unused (see HeapReleased
	// for a measure of the latter).
	//
	// HeapSys estimates the largest size the heap has had.
	HeapSys      uint64

	// HeapIdle is bytes in idle (unused) spans.
	//
	// Idle spans have no objects in them. These spans could be
	// (and may already have been) returned to the OS, or they can
	// be reused for heap allocations, or they can be reused as
	// stack memory.
	//
	// HeapIdle minus HeapReleased estimates the amount of memory
	// that could be returned to the OS, but is being retained by
	// the runtime so it can grow the heap without requesting more
	// memory from the OS. If this difference is significantly
	// larger than the heap size, it indicates there was a recent
	// transient spike in live heap size.
	HeapIdle     uint64

	// HeapInuse is bytes in in-use spans.
	//
	// In-use spans have at least one object in them. These spans
	// can only be used for other objects of roughly the same
	// size.
	//
	// HeapInuse minus HeapAlloc estimates the amount of memory
	// that has been dedicated to particular size classes, but is
	// not currently being used. This is an upper bound on
	// fragmentation, but in general this memory can be reused
	// efficiently.
	HeapInuse    uint64

	// HeapReleased is bytes of physical memory returned to the OS.
	//
	// This counts heap memory from idle spans that was returned
	// to the OS and has not yet been reacquired for the heap.
	HeapReleased uint64

	// StackInuse is bytes in stack spans.
	//
	// In-use stack spans have at least one stack in them. These
	// spans can only be used for other stacks of the same size.
	//
	// There is no StackIdle because unused stack spans are
	// returned to the heap (and hence counted toward HeapIdle).
	StackInuse   uint64

	// StackSys is bytes of stack memory obtained from the OS.
	//
	// StackSys is StackInuse, plus any memory obtained directly
	// from the OS for OS thread stacks (which should be minimal).
	StackSys     uint64

	// GCSys is bytes of memory in garbage collection metadata.
	GCSys        uint64

	// NextGC is the target heap size of the next GC cycle.
	//
	// The garbage collector's goal is to keep HeapAlloc ≤ NextGC.
	// At the end of each GC cycle, the target for the next cycle
	// is computed based on the amount of reachable data and the
	// value of GOGC.
	NextGC       uint64

	// LastGC is the time the last garbage collection finished, as
	// nanoseconds since 1970 (the UNIX epoch).
	LastGC       uint64

	// PauseTotalNs is the cumulative nanoseconds in GC
	// stop-the-world pauses since the program started.
	//
	// During a stop-the-world pause, all goroutines are paused
	// and only the garbage collector can run.
	PauseTotalNs uint64

	// NumGC is the number of completed GC cycles.
	NumGC        uint32

	// NumGoroutine returns the number of goroutines that currently exist.
	NumGoroutine int
}

func StartSnitchMonitor() {
	var snitch Snitch
	var rt runtime.MemStats
	runtime.ReadMemStats(&rt)

	snitch.Timestamp = time.Now().Format(time.RFC3339)
	snitch.NumGoroutine = runtime.NumGoroutine()
	snitch.Alloc = rt.Alloc
	snitch.TotalAlloc = rt.TotalAlloc
	snitch.Sys = rt.Sys
	snitch.Lookups = rt.Lookups
	snitch.Mallocs = rt.Mallocs
	snitch.Frees = rt.Frees
	snitch.LiveObjects = snitch.Mallocs - rt.Frees
	snitch.HeapAlloc = rt.HeapAlloc
	snitch.HeapSys = rt.HeapSys
	snitch.HeapIdle = rt.HeapIdle
	snitch.HeapInuse = rt.HeapInuse
	snitch.HeapReleased = rt.HeapReleased
	snitch.StackInuse = rt.StackInuse
	snitch.StackSys = rt.StackSys
	snitch.GCSys = rt.GCSys
	snitch.NextGC = rt.NextGC
	snitch.LastGC = rt.LastGC
	snitch.PauseTotalNs = rt.PauseTotalNs
	snitch.NumGC = rt.NumGC

	metrics, _ := json.Marshal(snitch)
	metrics = append(metrics, "\n"...)

	usr, _ := user.Current()
	configPath := usr.HomeDir

	needRotate, err := utils.LogMustRotate(configPath + SnitchLogFileName, MaxSnitchLogSize)
	if err != nil {
		log.Fatalf("failed to check if log needs to be rotated: %s", err.Error())
	}
	if needRotate {
		nBytes, err := utils.Rotate(configPath + SnitchLogFileName, configPath + SnitchTempFileName)
		if err != nil {
			log.Fatalf("failed to rotate snitch log file: %s", err.Error())
		}
		log.Printf("successfully rotated snitch log file: %d bytes rotated", nBytes)
	}

	file, err := os.OpenFile(configPath + SnitchLogFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatalf("failed to open snitch log file: %s", err.Error())
	}
	defer file.Close()

	// Underscore shows bytes written (FOR USE IN WATCHDOG METRICS)
	file.Write(metrics)
	if err != nil {
		log.Fatalf("failed to write to snitch log: %s", err.Error())
	}

	// Sets the finalizer associated with the object
	// allowing it to be released back to the heap
	runtime.SetFinalizer(&snitch, func(snitch *Snitch) {})
}

// GetSnitchMetrics reads the snitch log file and
// returns the entries in the log file as a Snitch
// object array.
func GetSnitchMetrics() response.CacheResponse {
	usr, _ := user.Current()
	configPath := usr.HomeDir

	file, err := os.Open(configPath + SnitchLogFileName)
	if err != nil {
		log.Fatalf("failed to open snitch log file: %s", err.Error())
	}
	defer file.Close()

	var data []Snitch
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		var entry Snitch
		line := scanner.Text()
		json.Unmarshal([]byte(line), &entry)
		data = append(data, entry)
	}
	
	res := response.NewResponseFromValue(data)
	return res
}
