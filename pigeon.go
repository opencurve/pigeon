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
 * Created Date: 2022-09-20
 * Author: Jingli Chen (Wine93)
 */

package pigeon

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/opencurve/pigeon/cli/command"
	"github.com/opencurve/pigeon/internal/configure"
	"github.com/opencurve/pigeon/internal/consts"
	"github.com/opencurve/pigeon/internal/core"
	"github.com/opencurve/pigeon/internal/http"
	"github.com/opencurve/pigeon/pkg/log"
	"go.uber.org/zap"
)

type (
	Request    = http.Request
	HTTPServer = http.HTTPServer
	Configure  = configure.ServerConfigure
	JSON       = gin.H
	Logger     = zap.Logger
)

var (
	NewHTTPServer = http.NewHTTPServer
	Field         = log.Field

	HTTP_METHOD_GET    = consts.HTTP_METHOD_GET
	HTTP_METHOD_HEAD   = consts.HTTP_METHOD_HEAD
	HTTP_METHOD_PUT    = consts.HTTP_METHOD_PUT
	HTTP_METHOD_POST   = consts.HTTP_METHOD_POST
	HTTP_METHOD_DELETE = consts.HTTP_METHOD_DELETE

	HTTP_STATUS_CONTINUE              = consts.HTTP_STATUS_CONTINUE
	HTTP_STATUS_OK                    = consts.HTTP_STATUS_OK
	HTTP_STATUS_CREATED               = consts.HTTP_STATUS_CREATED
	HTTP_STATUS_MOVED_PERMANENTLY     = consts.HTTP_STATUS_MOVED_PERMANENTLY
	HTTP_STATUS_MOVED_TEMPORARILY     = consts.HTTP_STATUS_MOVED_TEMPORARILY
	HTTP_STATUS_NOT_MODIFIED          = consts.HTTP_STATUS_NOT_MODIFIED
	HTTP_STATUS_BAD_REQUEST           = consts.HTTP_STATUS_BAD_REQUEST
	HTTP_STATUS_UNAUTHORIZED          = consts.HTTP_STATUS_UNAUTHORIZED
	HTTP_STATUS_FORBIDDEN             = consts.HTTP_STATUS_FORBIDDEN
	HTTP_STATUS_NOT_FOUND             = consts.HTTP_STATUS_NOT_FOUND
	HTTP_STATUS_NOT_ALLOWED           = consts.HTTP_STATUS_NOT_ALLOWED
	HTTP_STATUS_INTERNAL_SERVER_ERROR = consts.HTTP_STATUS_INTERNAL_SERVER_ERROR
	HTTP_STATUS_BAD_GATEWAY           = consts.HTTP_STATUS_BAD_GATEWAY
	HTTP_STATUS_SERVICE_UNAVAILABLE   = consts.HTTP_STATUS_SERVICE_UNAVAILABLE
	HTTP_STATUS_GATEWAY_TIMEOUT       = consts.HTTP_STATUS_GATEWAY_TIMEOUT
)

func Serve(servers ...*HTTPServer) {
	pigeon := core.NewPigeon(servers)
	cli := command.NewPigeonCommand(pigeon)
	if cli.Execute() != nil {
		os.Exit(1)
	}
}
