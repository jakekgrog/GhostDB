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
	"bufio"
	"os"
	"os/user"
	"testing"
	"time"

	"github.com/ghostdb/ghostdb-cache-node/config"
	"github.com/ghostdb/ghostdb-cache-node/store/monitor"
	"github.com/ghostdb/ghostdb-cache-node/store/request"
	"github.com/ghostdb/ghostdb-cache-node/utils"
)

func TestAppMetrics(t *testing.T) {
	var config config.Configuration = config.InitializeConfiguration()

	var store *Store
	store = NewStore("LRU")
	store.BuildStore(config)
	store.RunStore()

	// Delete pre-existing metrics
	usr, _ := user.Current()
	configPath := usr.HomeDir
	os.Remove(configPath + monitor.AppMetricsLogFilePath)
	os.Remove(configPath + "/ghostDBPersistence.log")

	store.Execute("add", request.NewRequestFromValues("Key1", "Value1", -1))
	store.Execute("add", request.NewRequestFromValues("Key2", "Value1", -1))
	store.Execute("add", request.NewRequestFromValues("Key3", "Value1", -1))
	store.Execute("put", request.NewRequestFromValues("Key1", "Value2", -1))
	store.Execute("put", request.NewRequestFromValues("Key4", "Value1", -1))
	store.Execute("get", request.NewRequestFromValues("Key1", "", -1))
	store.Execute("get", request.NewRequestFromValues("Key2", "", -1))
	store.Execute("get", request.NewRequestFromValues("Key5", "", -1))
	store.Execute("delete", request.NewRequestFromValues("Key1", "", -1))
	store.Execute("delete", request.NewRequestFromValues("Key1", "", -1))
	store.Execute("flush", request.NewRequestFromValues("Key1", "", -1))
	store.Execute("flush", request.NewRequestFromValues("Key1", "", -1))
	time.Sleep(11 * time.Second)

	utils.AssertEqual(t, utils.FileExists(configPath+monitor.AppMetricsLogFilePath), true, "")
	utils.AssertEqual(t, utils.FileNotEmpty(configPath+monitor.AppMetricsLogFilePath), true, "")

	file, err := os.Open(configPath + monitor.AppMetricsLogFilePath)
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(file)

	// Bug: scanner.Scan() doesn't move to next token
	scanner.Scan()
	scanner.Scan()
	metrics := scanner.Text()
	expectedOutput := `{"TotalHits": 12, "TotalGets": 3, "CacheMiss": 1, "TotalPuts": 2, "TotalAdds": 3, "NotStored": 0, "TotalDeletes": 2, "NotFound": 1, "TotalFlushes": 2, "ErrFlush": 2}`
	utils.AssertEqual(t, metrics, expectedOutput, "")
}
