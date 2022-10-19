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
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/opencurve/pigeon/internal/configure"
	"github.com/opencurve/pigeon/internal/utils"
	"github.com/opencurve/pigeon/pkg/log"
	"go.uber.org/zap"
)

type (
	Request struct {
		Context *gin.Context
		server  *HTTPServer
		ctx     map[string]interface{}
		Var     *Variable

		// request
		Method     string
		Scheme     string
		Host       string
		Uri        string
		Args       map[string]string
		HeadersIn  map[string]string
		BodyReader io.ReadCloser

		// response
		Status     int
		HeadersOut map[string]string
		content    content
	}
)

func NewRequest(c *gin.Context, server *HTTPServer) *Request {
	request := c.Request
	// request scheme
	scheme := utils.Choose(request.TLS == nil, "http", "https")
	// request headers
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
	// response headers
	version := server.cfg.GetContext().Version
	headersOut := map[string]string{
		"Server": "pigeon/" + version,
	}

	return &Request{
		Context: c,
		server:  server,
		ctx:     map[string]interface{}{},
		Var:     NewVariable(server),

		Method:     request.Method,
		Scheme:     scheme,
		Host:       request.URL.Host,
		Uri:        request.URL.Path,
		Args:       args,
		HeadersIn:  headers,
		BodyReader: request.Body,

		Status:     -1,
		HeadersOut: headersOut,
	}
}

func (r *Request) GetConfig() *configure.ModuleConfig {
	return r.server.cfg.GetConfig()
}

func (r *Request) SetModuleCtx(name string, ctx interface{}) {
	r.ctx[name] = ctx
}

func (r *Request) GetModuleCtx(name string) interface{} {
	return r.ctx[name]
}

func (r *Request) NextHandler() bool {
	return true
}

func (r *Request) GetURLParam(key string) string {
	return r.Context.Param(key)
}

func (r *Request) GetMultipartForm() (*multipart.Form, error) {
	return r.Context.MultipartForm()
}

func (r *Request) GetFormFile(key string) (multipart.File, *multipart.FileHeader, error) {
	return r.Context.Request.FormFile(key)
}

func (r *Request) BindBody(any interface{}) error {
	c := r.Context
	return c.ShouldBind(any)
}

func (r *Request) BindArgument(any interface{}) error {
	return r.Context.ShouldBindQuery(any)
}

func (r *Request) ProxyPass(address string, opts ...ProxyOption) bool {
	cfg := r.server.cfg
	options := PorxyOptions{
		Method:      r.Method,
		Scheme:      r.Scheme,
		Address:     address,
		Uri:         r.Uri,
		Args:        r.Args,
		Headers:     r.HeadersIn,
		Body:        r.BodyReader,
		ReadTimeout: cfg.GetProxyReadTimeout(),
	}
	for _, opt := range opts {
		opt(&options)
	}

	proxy := NewProxy(options)
	resp, err := proxy.Do()
	if err != nil {
		r.Status = http.StatusBadGateway
		r.Logger().Error("proxy pass failed", log.Field("error", err))
		return false
	}

	// handle response
	r.Status = resp.StatusCode
	for k, v := range resp.Header {
		r.HeadersOut[k] = v[0]
	}
	r.content = &Reader{
		reader: resp.Body,
		size:   resp.ContentLength,
		ctype:  r.HeadersOut["Content-Type"],
	}

	return false
}

func (r *Request) SendString(message string) bool {
	r.content = &Message{message: message}
	return false
}

func (r *Request) SendJSON(m interface{}) bool {
	r.content = &JSON{m: m}
	return false
}

func (r *Request) SendFile(filename string) bool {
	r.content = &File{filename: filename}
	return false
}

func (r *Request) SendBuffer(reader io.Reader, size int64) bool {
	r.content = &Buffer{reader: reader, size: size}
	return false
}

func (r *Request) send(buffer *Buffer) {
	var err error
	w := r.Context.Writer
	if buffer.size > 0 {
		_, err = io.CopyN(w, buffer.reader, buffer.size)
	} else {
		// FIXME: dangerous!!!
		_, err = io.Copy(w, buffer.reader)
	}
	if err != nil && err != io.EOF {
		panic(err)
	}
}

func (r *Request) Exit(code int, message ...string) bool {
	r.Status = code
	if len(message) > 0 {
		r.content = &Message{
			message: strings.Join(message, ""),
		}
	}
	return false
}

func (r *Request) Logger() *zap.Logger {
	return r.server.errorLogger
}

func (r *Request) log() {
	ctx := r.Context
	format := []string{
		ctx.RemoteIP(),
		ctx.Request.Method,
		ctx.Request.URL.RequestURI(),
		ctx.Request.Proto,
		strconv.Itoa(r.Status),
		fmt.Sprintf("%.3f", float64(utils.UnixMilli()-r.Var.StartTime)/1000),
		ctx.Request.UserAgent(),
		r.Var.LogAttach,
	}
	r.server.accessLogger.Info(strings.Join(format, " "))
}

func (r *Request) Finalize() {
	defer r.log() // access log
	ctx := r.Context

	// response status
	if r.Status < 100 {
		r.Status = 200
	}
	ctx.Status(r.Status)

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
		ctx.String(r.Status, content.(*Message).message)
	case *JSON:
		ctx.JSON(r.Status, content.(*JSON).m)
	case *File:
		ctx.File(content.(*File).filename)
	case *Reader:
		reader := content.(*Reader)
		ctx.DataFromReader(r.Status, reader.size, reader.ctype, reader.reader, nil)
	case *Buffer:
		buffer := content.(*Buffer)
		r.send(buffer)
	}
}
