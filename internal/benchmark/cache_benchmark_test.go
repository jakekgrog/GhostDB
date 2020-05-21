package benchmark

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/internal/ghost_config"
	"github.com/ghostdb/ghostdb-cache-node/cache/lru_cache"
)

func BenchmarkPutToCacheParallel(b *testing.B) {
	writeToCachePutParallel(b)
}

func BenchmarkPutToCache(b *testing.B) {
	b.ReportAllocs()
	
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	cache := lru_cache.NewLRU(config)
	cache.Size = int32(b.N)
	
	for i := 0; i < b.N; i++ {
		cache.Put(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1)
	}
}

func BenchmarkAddToCacheParallel(b *testing.B) {
	writeToCacheAddParallel(b)
}

func BenchmarkAddToCache(b *testing.B) {
	b.ReportAllocs()
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	cache := lru_cache.NewLRU(config)
	cache.Size = int32(b.N)
	for i := 0; i < b.N; i++ {
		cache.Add(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1)
	}
}

func BenchmarkGetFromCacheParallel(b *testing.B) {
	getFromCacheParallel(b)
}

func BenchmarkGetFromCache(b *testing.B) {
	b.StopTimer()
	
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	cache := lru_cache.NewLRU(config)
	cache.Size = int32(b.N)
	
	for i := 0; i < b.N; i++ {
		cache.Put(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1)
	}

	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cache.Get(fmt.Sprintf("key-%d", i))
	}
}

func BenchmarkDeleteFromCacheParallel(b *testing.B) {
	deleteFromCacheParallel(b)
}

func BenchmarkDeleteFromCache(b *testing.B) {
	b.StopTimer()
	
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	cache := lru_cache.NewLRU(config)
	cache.Size = int32(b.N)

	for i := 0; i < b.N; i++ {
		cache.Put(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1)
	}

	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cache.Delete(fmt.Sprintf("key-%d", i))
	}
}

func BenchmarkFlushCache(b *testing.B) {
	b.StopTimer()
	
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	cache := lru_cache.NewLRU(config)
	cache.Size = int32(b.N)

	for i := 0; i < b.N; i++ {
		cache.Put(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1)
	}

	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cache.Flush()
	}
}

func writeToCachePutParallel(b *testing.B) {
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	cache := lru_cache.NewLRU(config)
	cache.Size = int32(b.N)

	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Int()
		counter := 0
		
		b.ReportAllocs()
		for pb.Next() {
			cache.Put(fmt.Sprintf("key-%d-%d", id, counter), fmt.Sprintf("value-%d-%d", counter, id), -1)
			counter = counter + 1
		}
	})
}

func writeToCacheAddParallel(b *testing.B) {
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	cache := lru_cache.NewLRU(config)
	cache.Size = int32(b.N)

	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Int()
		counter := 0
		
		b.ReportAllocs()
		for pb.Next() {
			cache.Add(fmt.Sprintf("key-%d-%d", id, counter), fmt.Sprintf("value-%d-%d", counter, id), -1)
			counter = counter + 1
		}
	})
}

func getFromCacheParallel(b *testing.B) {
	b.StopTimer()
	
	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	cache := lru_cache.NewLRU(config)
	cache.Size = int32(b.N)

	for i := 0; i < b.N; i++ {
		cache.Put(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1)
	}

	b.StartTimer()

	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		
		b.ReportAllocs()
		for pb.Next() {
			cache.Get(fmt.Sprintf("key-%d", i))
			i = i + 1
		}
	})
}

func deleteFromCacheParallel(b *testing.B) {
	b.StopTimer()

	var config ghost_config.Configuration = ghost_config.InitializeConfiguration()
	cache := lru_cache.NewLRU(config)
	cache.Size = int32(b.N)

	for i := 0; i < b.N; i++ {
		cache.Put(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1)
	}

	b.StartTimer()

	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		
		b.ReportAllocs()
		for pb.Next() {
			cache.Delete(fmt.Sprintf("key-%d", i))
			i = i + 1
		}
	})
}