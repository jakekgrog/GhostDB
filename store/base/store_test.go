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

package base

import (
	"testing"

	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/utils"
)

func TestStore(t *testing.T) {
	conf := config.InitializeConfiguration()

	store := NewStore("LRU")
	store.BuildStore(conf)
	store.RunStore()

	x := store.Execute("get", request.NewRequestFromValues("Key1", "NewValue1", -1))
	utils.AssertEqual(t, x.Message, "CACHE_MISS", "")

	x = store.Execute("put", request.NewRequestFromValues("Key1", "NewValue1", -1))
	utils.AssertEqual(t, x.Status, int32(1), "")
}

func TestStoreQueue(t *testing.T) {
	conf := config.InitializeConfiguration()

<<<<<<< HEAD
	store := NewStore("LRU")
=======
	var store *Store
	store = NewStore("LRU")
>>>>>>> Added Queue command support to store
	store.BuildStore(conf)
	store.RunStore()

	// Enqueue a value
	req := request.NewRequestFromValues("QueueKey", "first", -1)
	res := store.Execute("enqueue", req)
	utils.AssertEqual(t, res.Message, "STORED", "")

	// Enqueu another value
	req = request.NewRequestFromValues("QueueKey", "second", -1)
	res = store.Execute("enqueue", req)
	utils.AssertEqual(t, res.Message, "STORED", "")

	// Dequeue and assert the value is the same as the first queued item
	req = request.NewRequestFromValues("QueueKey", "", -1)
	res = store.Execute("dequeue", req)
	utils.AssertEqual(t, res.Message, "OK", "")
	utils.AssertEqual(t, res.Gobj.Value, "first", "")

	// Dequeue and assert the value is the same as the second queued item
	req = request.NewRequestFromValues("QueueKey", "", -1)
	res = store.Execute("dequeue", req)
	utils.AssertEqual(t, res.Message, "OK", "")
	utils.AssertEqual(t, res.Gobj.Value, "second", "")

	// Dequeue and assert the value is nil (queue empty)
	req = request.NewRequestFromValues("QueueKey", "", -1)
	res = store.Execute("dequeue", req)
	utils.AssertEqual(t, res.Message, "OK", "")
	utils.AssertEqual(t, res.Gobj.Value, nil, "")
}
