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

	"github.com/ghostdb/ghostdb-cache-node/utils"
	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/store/persistence"
)

func TestSerializer(t *testing.T) {
	var config config.Configuration = config.InitializeConfiguration()

	var store *Store
	store = NewStore("LRU")
	store.BuildStore(config)

	store.Execute("put", request.NewRequestFromValues("Italy", "Rome", -1))
	store.Execute("put", request.NewRequestFromValues("England", "London", 2))

	encryptionEnabled := config.EnableEncryption
	passphrase := "SUPPLY_PASSPHRASE"

	store.CreateSnapshot()

	bytes := persistence.ReadSnapshot(encryptionEnabled, passphrase)

	c, err := persistence.BuildCacheFromSnapshot(bytes)
	if err != nil {
		panic(err)
	}

	store.Cache = &c

	val := store.Execute("get", request.NewRequestFromValues("England", "", -1))

	utils.AssertEqual(t, val.Gobj.Key, "London", "")
	utils.AssertEqual(t, store.Execute("nodeSize", request.NewEmptyRequest()), int32(2), "")

	// Test the config was rebuilt correctly.
	utils.AssertEqual(t, store.Conf.KeyspaceSize, int32(65536), "")
	utils.AssertEqual(t, store.Conf.SysMetricInterval, int32(300), "")
	utils.AssertEqual(t, store.Conf.AppMetricInterval, int32(300), "")
	utils.AssertEqual(t, store.Conf.DefaultTTL, int32(-1), "")
	utils.AssertEqual(t, store.Conf.CrawlerInterval, int32(300), "")
	utils.AssertEqual(t, store.Conf.SnapshotInterval, int32(3600), "")
	utils.AssertEqual(t, store.Conf.PersistenceAOF, false, "")
	utils.AssertEqual(t, store.Conf.EntryTimestamp, true, "")
	utils.AssertEqual(t, store.Conf.EnableEncryption, true, "")
}
