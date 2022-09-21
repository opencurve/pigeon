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

package consts

func init() { /* do nothing */ }

const (
	HTTP_SCHEME_HTTP  = "http"
	HTTP_SCHEME_HTTPS = "https"

	HTTP_METHOD_GET    = "GET"
	HTTP_METHOD_HEAD   = "HEAD"
	HTTP_METHOD_PUT    = "PUT"
	HTTP_METHOD_POST   = "POST"
	HTTP_METHOD_DELETE = "DELETE"

	HTTP_STATUS_CONTINUE              = 100
	HTTP_STATUS_OK                    = 200
	HTTP_STATUS_CREATED               = 201
	HTTP_STATUS_MOVED_PERMANENTLY     = 301
	HTTP_STATUS_MOVED_TEMPORARILY     = 302
	HTTP_STATUS_NOT_MODIFIED          = 304
	HTTP_STATUS_BAD_REQUEST           = 400
	HTTP_STATUS_UNAUTHORIZED          = 401
	HTTP_STATUS_FORBIDDEN             = 403
	HTTP_STATUS_NOT_FOUND             = 404
	HTTP_STATUS_NOT_ALLOWED           = 405
	HTTP_STATUS_INTERNAL_SERVER_ERROR = 500
	HTTP_STATUS_BAD_GATEWAY           = 502
	HTTP_STATUS_SERVICE_UNAVAILABLE   = 503
	HTTP_STATUS_GATEWAY_TIMEOUT       = 504
)
