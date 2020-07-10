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
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/store/persistence"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/utils"
)

func TestAOF(t *testing.T) {
	configPath, _ := os.UserConfigDir()

	var config config.Configuration = config.InitializeConfiguration()
	const aofMaxBytes = 300

	err := os.Remove(configPath + "/ghostDBPersistence.log")
	if err != nil {
		return
	}

	store := NewStore("LRU")
	store.BuildStore(config)
	store.RunStore()

	store.Execute("add", request.NewRequestFromValues("Key1", "Value1", -1))
	store.Execute("add", request.NewRequestFromValues("Key2", "Value2", -1))

	for i := 1; i <= 100; i++ {
		store.Execute("put", request.NewRequestFromValues("Key1", "NewValue1", -1))
		time.Sleep(10 * time.Millisecond)
	}

	// Give go routine time to flush/write buffer & rewrite file
	time.Sleep(2 * time.Second)
	// Check file has shrunk below max size
	if persistence.GetAOFSize() >= aofMaxBytes {
		fmt.Println(persistence.GetAOFSize())
		t.Error("AOF Size exceeded threshold")
	}
	// Simulate cache restart
	newStore := NewStore("LRU")
	newStore.BuildStore(config)
	newStore.RunStore()

	utils.AssertEqual(t, newStore.Execute("get", request.NewRequestFromValues("Key1", "", -1)), "NewValue1", "")
}
