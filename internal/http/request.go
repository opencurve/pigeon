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
	request := c.Request
	headers := map[string]string{}
	for k := range c.Request.Header {
		headers[k] = c.GetHeader(k)
	}
	args := map[string]string{}
	for _, param := range c.Params {
		args[param.Key] = param.Value
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
		r.content = &Message{message: message[0]}
	}
	return false
}

func (r *Request) Finalize() {
	ctx := r.Context

	// status code
	status := 200
	if r.Status != -1 {
		ctx.Status(status)
	}

	// response headers
	for k, v := range r.HeadersOut {
		ctx.Header(k, v)
	}

	content := r.content
	if content != nil {
		switch content.(type) {
		case *Message:
			ctx.String(status, content.(*Message).message)
		case *File:
			ctx.File(content.(*File).filename)
		}
	}
}
