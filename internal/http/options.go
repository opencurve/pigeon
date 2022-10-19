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
 * Created Date: 2022-09-22
 * Author: Jingli Chen (Wine93)
 */

package http

import (
	"time"
)

type (
	ProxyOption func(*PorxyOptions)

	PorxyOptions struct {
		Scheme         string
		Address        string
		Method         string
		Uri            string
		Args           map[string]string
		Headers        map[string]string
		Body           interface{}
		ConnectTimeout time.Duration
		ReadTimeout    time.Duration
	}
)

func (r *Request) WithScheme(scheme string) ProxyOption {
	return func(options *PorxyOptions) {
		options.Scheme = scheme
	}
}

func (r *Request) WithMethod(method string) ProxyOption {
	return func(options *PorxyOptions) {
		options.Method = method
	}
}

func (r *Request) WithURI(uri string) ProxyOption {
	return func(options *PorxyOptions) {
		options.Uri = uri
	}
}

func (r *Request) WithArguments(args map[string]string) ProxyOption {
	return func(options *PorxyOptions) {
		options.Args = args
	}
}

func (r *Request) WithHeaders(headers map[string]string) ProxyOption {
	return func(options *PorxyOptions) {
		options.Headers = headers
	}
}

func (r *Request) WithBody(body interface{}) ProxyOption {
	return func(options *PorxyOptions) {
		options.Body = body
	}
}
