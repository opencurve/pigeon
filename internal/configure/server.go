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
	"path"
	"path/filepath"
	"time"

	"github.com/opencurve/pigeon/internal/utils"
)

type (
	ServerConfigure = Server

	ModuleConfig struct {
		m map[string]interface{}
	}
)

func (cfg *ServerConfigure) absPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return path.Join(cfg.context.Prefix, filename)
}

func (cfg *ServerConfigure) GetContext() Context {
	return cfg.context
}

func (cfg *ServerConfigure) GetName() string {
	return cfg.Name
}

func (cfg *ServerConfigure) GetEnable() bool {
	return cfg.Enable
}

func (cfg *ServerConfigure) GetListenAddress() string {
	return cfg.Listen
}

func (cfg *ServerConfigure) GetAccessLogPath() string {
	return cfg.absPath(cfg.AccessLog)
}

func (cfg *ServerConfigure) GetErrorLogPath() string {
	return cfg.absPath(cfg.ErrorLog)
}

func (cfg *ServerConfigure) GetLogLevel() string {
	return cfg.LogLevel
}

func (cfg *ServerConfigure) GetIndex() string {
	return cfg.absPath(cfg.Index)
}

func (cfg *ServerConfigure) GetProxyConnectTimeout() time.Duration {
	return time.Duration(cfg.ProxyConnectTimeout) * time.Second
}

func (cfg *ServerConfigure) GetProxySendTimeout() time.Duration {
	return time.Duration(cfg.ProxySendTimeout) * time.Second
}

func (cfg *ServerConfigure) GetProxyReadTimeout() time.Duration {
	return time.Duration(cfg.ProxyReadTimeout) * time.Second
}

func (cfg *ServerConfigure) GetProxyNextUpstreamTries() int {
	return cfg.ProxyNextUpstreamTries
}

func (cfg *ServerConfigure) GetConfig() *ModuleConfig {
	return &ModuleConfig{m: cfg.Config}
}

func (cfg *ModuleConfig) GetBool(key string) bool {
	v, ok := cfg.m[key]
	if !ok {
		return false
	}

	val, ok := utils.Str2Bool(v.(string))
	if ok {
		return val
	}
	return false
}

func (cfg ModuleConfig) GetString(key string) string {
	v, ok := cfg.m[key]
	if !ok {
		return ""
	}
	return v.(string)
}
