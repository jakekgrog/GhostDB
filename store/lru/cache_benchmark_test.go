/*
 * Copyright (c) 2020, Jake Grogan
 * All rights reserved.
 * 
 * Redistribution and use in source and binary forms, with or without
 * modification, are permitted provided that the following conditions are met:
 * 
 *  * Redistributions of source code must retain the above copyright notice, this
 *    list of conditions and the following disclaimer.
 * 
 *  * Redistributions in binary form must reproduce the above copyright notice,
 *    this list of conditions and the following disclaimer in the documentation
 *    and/or other materials provided with the distribution.
 * 
 *  * Neither the name of the copyright holder nor the names of its
 *    contributors may be used to endorse or promote products derived from
 *    this software without specific prior written permission.
 * 
 * THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
 * AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
 * IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
 * DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
 * FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
 * DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
 * SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
 * CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
 * OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
 * OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.
*/

package lru

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
)

func BenchmarkPutToCacheParallel(b *testing.B) {
	writeToCachePutParallel(b)
}

func BenchmarkPutToCache(b *testing.B) {
	b.ReportAllocs()
	
	var config config.Configuration = config.InitializeConfiguration()
	cache := NewLRU(config)
	cache.Size = int32(b.N)
	
	for i := 0; i < b.N; i++ {
		cache.Put(request.NewRequestFromValues(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1))
	}
}

func BenchmarkAddToCacheParallel(b *testing.B) {
	writeToCacheAddParallel(b)
}

func BenchmarkAddToCache(b *testing.B) {
	b.ReportAllocs()
	var config config.Configuration = config.InitializeConfiguration()
	cache := NewLRU(config)
	cache.Size = int32(b.N)
	for i := 0; i < b.N; i++ {
		cache.Add(request.NewRequestFromValues(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1))
	}
}

func BenchmarkGetFromCacheParallel(b *testing.B) {
	getFromCacheParallel(b)
}

func BenchmarkGetFromCache(b *testing.B) {
	b.StopTimer()
	
	var config config.Configuration = config.InitializeConfiguration()
	cache := NewLRU(config)
	cache.Size = int32(b.N)
	
	for i := 0; i < b.N; i++ {
		cache.Put(request.NewRequestFromValues(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1))
	}

	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cache.Get(request.NewRequestFromValues(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1))
	}
}

func BenchmarkDeleteFromCacheParallel(b *testing.B) {
	deleteFromCacheParallel(b)
}

func BenchmarkDeleteFromCache(b *testing.B) {
	b.StopTimer()
	
	var config config.Configuration = config.InitializeConfiguration()
	cache := NewLRU(config)
	cache.Size = int32(b.N)

	for i := 0; i < b.N; i++ {
		cache.Put(request.NewRequestFromValues(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1))
	}

	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cache.Delete(request.NewRequestFromValues(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1))
	}
}

func BenchmarkFlushCache(b *testing.B) {
	b.StopTimer()
	
	var config config.Configuration = config.InitializeConfiguration()
	cache := NewLRU(config)
	cache.Size = int32(b.N)

	for i := 0; i < b.N; i++ {
		cache.Put(request.NewRequestFromValues(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1))
	}

	b.StartTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		cache.Flush(request.NewRequestFromValues("Key1", "", -1))
	}
}

func writeToCachePutParallel(b *testing.B) {
	var config config.Configuration = config.InitializeConfiguration()
	cache := NewLRU(config)
	cache.Size = int32(b.N)

	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Int()
		counter := 0
		
		b.ReportAllocs()
		for pb.Next() {
			
			cache.Put(request.NewRequestFromValues(fmt.Sprintf("key-%d-%d", id, counter), fmt.Sprintf("value-%d-%d", counter, id), -1))
			counter = counter + 1
		}
	})
}

func writeToCacheAddParallel(b *testing.B) {
	var config config.Configuration = config.InitializeConfiguration()
	cache := NewLRU(config)
	cache.Size = int32(b.N)

	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		id := rand.Int()
		counter := 0
		
		b.ReportAllocs()
		for pb.Next() {
			cache.Add(request.NewRequestFromValues(fmt.Sprintf("key-%d-%d", id, counter), fmt.Sprintf("value-%d-%d", counter, id), -1))
			counter = counter + 1
		}
	})
}

func getFromCacheParallel(b *testing.B) {
	b.StopTimer()
	
	var config config.Configuration = config.InitializeConfiguration()
	cache := NewLRU(config)
	cache.Size = int32(b.N)

	for i := 0; i < b.N; i++ {
		cache.Put(request.NewRequestFromValues(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1))
	}

	b.StartTimer()

	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		
		b.ReportAllocs()
		for pb.Next() {
			cache.Get(request.NewRequestFromValues(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1))
			i = i + 1
		}
	})
}

func deleteFromCacheParallel(b *testing.B) {
	b.StopTimer()

	var config config.Configuration = config.InitializeConfiguration()
	cache := NewLRU(config)
	cache.Size = int32(b.N)

	for i := 0; i < b.N; i++ {
		cache.Put(request.NewRequestFromValues(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1))
	}

	b.StartTimer()

	rand.Seed(time.Now().Unix())

	b.RunParallel(func(pb *testing.PB) {
		i := 0
		
		b.ReportAllocs()
		for pb.Next() {
			cache.Delete(request.NewRequestFromValues(fmt.Sprintf("key-%d", i), fmt.Sprintf("value-%d", i), -1))
			i = i + 1
		}
	})
}