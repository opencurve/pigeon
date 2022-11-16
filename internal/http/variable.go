/*
 *  Copyright (c) 2022 NetEase Inc.
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 */

/*
 * Project: Pigeon
 * Created Date: 2022-09-28
 * Author: Jingli Chen (Wine93)
 */

package http

import (
	"github.com/gin-gonic/gin"
	"github.com/opencurve/pigeon/internal/utils"
)

type Variable struct {
	RemoteAddr string
	StartTime  int64
	ServerAddr string
	Index      string
	RequestURI string
	LogAttach  string
}

func NewVariable(server *HTTPServer, ctx *gin.Context) *Variable {
	cfg := server.cfg
	return &Variable{
		RemoteAddr: ctx.RemoteIP(),
		StartTime:  utils.UnixMilli(),
		ServerAddr: cfg.GetListenAddress(),
		RequestURI: ctx.Request.URL.RequestURI(),
		Index:      cfg.GetIndex(),
		LogAttach:  "-",
	}
}
