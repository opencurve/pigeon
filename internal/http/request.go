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
	"io"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opencurve/pigeon/internal/configure"
)

type Buffer struct {
	Data []byte
	Err  error
}

type (
	Request struct {
		Context   *gin.Context
		configure *configure.Configure

		// request
		Method    string
		Scheme    string
		Host      string
		Uri       string
		Args      map[string]string
		HeadersIn map[string]string
		Body      io.ReadCloser

		// response
		Status     int
		HeadersOut map[string]string

		content content
	}
)

func NewRequest(c *gin.Context, cfg *configure.Configure) *Request {
	// request headers
	request := c.Request
	headers := map[string]string{}
	for k := range request.Header {
		headers[k] = c.GetHeader(k)
	}
	// request arguments
	args := map[string]string{}
	values := request.URL.Query()
	for key := range values {
		args[key] = values.Get(key)
	}

	return &Request{
		Context:   c,
		configure: cfg,

		Method:    request.Method,
		Scheme:    request.URL.Scheme,
		Host:      request.URL.Host,
		Uri:       request.RequestURI,
		Args:      args,
		HeadersIn: headers,
		Body:      request.Body,

		Status:     -1,
		HeadersOut: map[string]string{},
	}
}

func (r *Request) ReadBody(opts ...ReadOption) <-chan Buffer {
	options := DefaultReadOptions
	for _, opt := range opts {
		opt(options)
	}

	ch := make(chan Buffer)
	go func() {
		defer close(ch)
		defer r.Body.Close()
		buffer := make([]byte, options.BufferSize)
		for {
			_, err := r.Body.Read(buffer)
			ch <- Buffer{Data: buffer, Err: err}
			if err != nil {
				break
			}
		}
	}()

	return ch
}

func (r *Request) Bind(any interface{}) error {
	return r.Context.ShouldBind(any)
}

func (r *Request) ProxyPass() bool {
	return false
}

func (r *Request) SendString(message string) bool {
	r.content = &Message{message: message}
	return false
}

func (r *Request) SendJSON(m gin.H) bool {
	r.content = &JSON{m: m}
	return false
}

func (r *Request) SendFile(filename string) bool {
	r.content = &File{filename: filename}
	return false
}

func (r *Request) NextHandler() bool {
	return true
}

func (r *Request) Exit(code int, message ...string) bool {
	r.Status = code
	if len(message) > 0 {
		r.content = &Message{message: strings.Join(message, "")}
	}
	return false
}

func (r *Request) Finalize() {
	ctx := r.Context

	// response status
	status := 200
	if r.Status != -1 {
		status = r.Status
	}
	ctx.Status(status)

	// response headers
	for k, v := range r.HeadersOut {
		ctx.Header(k, v)
	}

	// response body
	content := r.content
	if content == nil {
		return
	}

	switch content.(type) {
	case *Message:
		ctx.String(status, content.(*Message).message)
	case *JSON:
		ctx.JSON(status, content.(*JSON).m)
	case *File:
		ctx.File(content.(*File).filename)
	}
}
