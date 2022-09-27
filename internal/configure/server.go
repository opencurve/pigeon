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

package configure

import (
	"fmt"
)

type ServerConfigure struct {
	Name       string
	ListenIP   string
	ListenPort uint
	AccessLog  string
	ErrorLog   string
	LogLevel   string
}

func DefaultServer() *ServerConfigure {
	return &ServerConfigure{
		Name:       "default",
		ListenIP:   "0.0.0.0",
		ListenPort: 8000,
		AccessLog:  "pigeon_access.log",
		ErrorLog:   "pigeon_error.log",
		LogLevel:   "error",
	}
}
func (cfg *ServerConfigure) GetListenAddress() string {
	return fmt.Sprintf("%s:%d", cfg.ListenIP, cfg.ListenPort)
}
func (cfg *ServerConfigure) GetAccessLogPath() string { return cfg.AccessLog }
func (cfg *ServerConfigure) GetErrorLogPath() string  { return cfg.ErrorLog }
func (cfg *ServerConfigure) GetLogLevel() string  { return cfg.LogLevel }
