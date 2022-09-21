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

package http

import (
	"github.com/gin-gonic/gin"
)

var _ IRoutes = (*router)(nil)

type (
	HandlerFunc func(r *Request) bool

	HandlersChain []HandlerFunc

	RouterGroup interface {
		Group(string, ...HandlerFunc) RouterGroup
		IRoutes
	}

	IRoutes interface {
		Route(string, ...HandlerFunc)
	}

	router struct {
		server   *HTTPServer
		group    *gin.RouterGroup
		handlers HandlersChain
	}
)

func (r *router) Group(relativePath string, handlers ...HandlerFunc) RouterGroup {
	return &router{
		server:   r.server,
		group:    r.group.Group(relativePath),
		handlers: append(r.handlers, handlers...),
	}
}

func (r *router) Route(relativePath string, handlers ...HandlerFunc) {
	r.handlers = append(r.handlers, handlers...)
	r.group.Any(relativePath, r.warpHandlers())
}

func (r *router) warpHandlers() gin.HandlerFunc {
	return func(c *gin.Context) {
		request := NewRequest(c, r.server.configure)
		for _, handler := range r.handlers {
			if !handler(request) {
				break
			}
		}
		request.Finalize()
	}
}
