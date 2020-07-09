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
	"testing"

	"github.com/ghostdb/ghostdb-cache-node/utils"
	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
)

func TestLru(t *testing.T) {
	var config config.Configuration = config.InitializeConfiguration()
	cache := NewLRU(config)
	cache.Size = int32(2);
	cache.Put(request.NewRequestFromValues("England", "London", -1))
	cache.Put(request.NewRequestFromValues("Ireland", "Dublin", -1))
	
	// HEAD -> Dublin -> London

	cache.Put(request.NewRequestFromValues("America", "Washington", -1)) // England:London evicted here

	// HEAD -> Washington -> Dublin

	v1 := cache.Get(request.NewRequestFromValues("England", "", -1)) // Should be a cache miss
	utils.AssertEqual(t, "CACHE_MISS", v1.Message, "")

	// Ireland should be next to be evicted
	// If we 'Get' Ireland then it should be considered MRU
	// And 'America' Should now be LRU
	v2 := cache.Get(request.NewRequestFromValues("Ireland", "", -1))
	utils.AssertEqual(t, "Dublin", v2.Gobj.Value, "")
	
	// HEAD -> Dublin -> Washington

	cache.Put(request.NewRequestFromValues("France", "Paris", -1)) // America should be evicted here
	
	// HEAD -> Paris -> Dublin

	v3 := cache.Get(request.NewRequestFromValues("America", "", -1)) // Should be a cache miss
	utils.AssertEqual(t, CACHE_MISS, v3.Message, "")
	
	cache.Put(request.NewRequestFromValues("Italy", "Rome", -1)) // Ireland should be evicted here

	// HEAD -> Rome -> Paris
	
	v4 := cache.Get(request.NewRequestFromValues("France", "", -1))
	utils.AssertEqual(t, "Paris", v4.Gobj.Value, "")

	// HEAD -> Paris -> Rome

	message := cache.Add(request.NewRequestFromValues("France", "Paris", -1))
	utils.AssertEqual(t, NOT_STORED, message.Message, "")

	message = cache.Add(request.NewRequestFromValues("Poland", "Warsaw", -1))
	utils.AssertEqual(t, STORED, message.Message, "")

	message = cache.Delete(request.NewRequestFromValues("Poland", "", -1))
	utils.AssertEqual(t, REMOVED, message.Message, "")

	message = cache.CountKeys(request.NewRequestFromValues("Key1", "", -1))
	utils.AssertEqual(t, message.Gobj.Value.(int32) > 0, true, "")

	message = cache.Delete(request.NewRequestFromValues("USA", "", -1))
	utils.AssertEqual(t, NOT_FOUND, message.Message, "")

	message = cache.Put(request.NewRequestFromValues("England", "London", -1))
	utils.AssertEqual(t, STORED, message.Message, "")

	message = cache.Put(request.NewRequestFromValues("England", "London", -1))
	utils.AssertEqual(t, STORED, message.Message, "")

	message = cache.Flush(request.NewRequestFromValues("Key1", "", -1))
	utils.AssertEqual(t, FLUSHED, message.Message, "")

	message = cache.CountKeys(request.NewRequestFromValues("Key1", "", -1))
	utils.AssertEqual(t, message.Gobj.Value.(int32), int32(0), "")
}
