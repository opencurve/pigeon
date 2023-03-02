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
	"net/http"
	"net/url"

	"github.com/imroc/req/v3"
)

var (
	client *req.Client
)

func init() {
	client = req.NewClient()
}

type Proxy struct {
	options PorxyOptions
}

func NewProxy(options PorxyOptions) *Proxy {
	return &Proxy{options: options}
}

func (p *Proxy) makeURL() string {
	options := p.options
	return (&url.URL{
		Scheme:   options.Scheme,
		Host:     options.Address,
		Path:     options.Uri,
		RawQuery: options.Args,
	}).String()
}

func (p *Proxy) Do() (*http.Response, error) {
	options := p.options
	resp, err := client.R().
		SetHeaders(options.Headers).
		SetBody(options.Body).
		Send(options.Method, p.makeURL())
	if err != nil {
		return nil, err
	}

	return resp.Response, nil
}
