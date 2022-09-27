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
	"github.com/opencurve/pigeon/internal/core"
	"github.com/opencurve/pigeon/internal/http"
	"github.com/opencurve/pigeon/pkg/log"
)

type (
	Request    = http.Request
	HTTPServer = http.HTTPServer
	JSON       = gin.H
)

var (
	NewHTTPServer = http.NewHTTPServer
	Field        = log.Field
)

func Serve(servers ...*HTTPServer) {
	pigeon := core.NewPigeon(servers)
	cli := command.NewPigeonCommand(pigeon)
	if cli.Execute() != nil {
		os.Exit(1)
	}
}
