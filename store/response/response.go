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

package response

import (
	"fmt"
	
	"github.com/ghostdb/ghostdb-cache-node/store/object"
)

const (
	INVALID_COMMAND_ERR = "INVALID_COMMAND_ERR"
)

type CacheResponse struct {
	// GhostDB Cache Object
	Gobj    object.CacheObject
	// Status of the result - 0 or 1 - fail or succeed
	Status  int32
	// Message contains any textual data that is to 
	// be sent back
	Message string
	// Error message returned if something went wrong
	// during command execution
	Error   string
}

func NewResponseFromValue(value interface{}) CacheResponse{
	return CacheResponse{
		Gobj: object.NewCacheObjectFromValue(value),
		Status: 1,
		Message: "OK",
		Error: "",
	}
}

func NewCacheMissResponse() CacheResponse {
	return CacheResponse {
		Gobj: object.NewEmptyCacheObject(),
		Status: 0,
		Message: "CACHE_MISS",
		Error: "",
	}
}

func NewResponseFromMessage(msg string, status int32) CacheResponse {
	return CacheResponse {
		Gobj: object.NewEmptyCacheObject(),
		Status: status,
		Message: msg,
		Error: "",
	}
}

func BadCommandResponse(cmd string) CacheResponse {
	return CacheResponse {
		Gobj: object.NewEmptyCacheObject(),
		Status: 0,
		Message: fmt.Sprintf("Command '%s' is not a recognized command", cmd),
		Error: INVALID_COMMAND_ERR,
	}
}

func NewPingResponse() CacheResponse {
	return CacheResponse {
		Gobj: object.NewEmptyCacheObject(),
		Status: 1,
		Message: "Pong!",
		Error: "",
	}
}