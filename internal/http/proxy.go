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
	"net"
	"net/http"
	"net/url"

	"github.com/opencurve/pigeon/internal/utils"
)

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
		RawPath:  options.Uri,
		RawQuery: utils.MakeArgument(options.Args),
	}).String()
}

func (p *Proxy) Do() (*http.Response, error) {
	options := p.options

	client := &http.Client{
		Timeout: options.ReadTimeout,
		Transport: &http.Transport{
			DialContext: (&net.Dialer{
				Timeout: options.ConnectTimeout,
			}).DialContext,
		},
	}

	req, err := http.NewRequest(options.Method, p.makeURL(), options.Body)
	if err != nil {
		return nil, err
	}

	for k, v := range options.Headers {
		req.Header.Set(k, v)
	}
	return client.Do(req)
}
