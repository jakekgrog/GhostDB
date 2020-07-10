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

package cache

import (
	"github.com/ghostdb/ghostdb-cache-node/store/lru"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/store/response"
)

// Cache is an interface for the cache object
type Cache interface {
	// Put will add a key/value pair to the cache, possibly
	// overwriting an existing key/value pair. Put will evict
	// a key/value pair if the cache is full.
	Put(reqObj request.CacheRequest) response.CacheResponse

	// Get will fetch a key/value pair from the cache
	Get(reqObj request.CacheRequest) response.CacheResponse

	// Add will add a key/value pair to the cache if the key
	// does not exist already. It will not evict a key/value pair
	// from the cache. If the cache is full, the key/value pair does
	// not get added.
	Add(reqObj request.CacheRequest) response.CacheResponse

	// Delete removes a key/value pair from the cache
	// Returns NOT_FOUND if the key does not exist.
	Delete(reqObj request.CacheRequest) response.CacheResponse

	// DeleteByKey functions the same as Delete, however it is
	// used in various locations to reduce the cost of allocating
	// request objects for internal deletion mechanisms
	// e.g. the cache crawlers.
	DeleteByKey(key string) response.CacheResponse

	// Flush removes all key/value pairs from the cache even if they
	// have not expired
	Flush(request.CacheRequest) response.CacheResponse

	// CountKeys return the number of keys in the cache
	CountKeys(request.CacheRequest) response.CacheResponse

	// GetHashtableReference is for internal use by crawlers and AOF
	GetHashtableReference() *map[string]*lru.Node
}
