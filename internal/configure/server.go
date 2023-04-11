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

func (cfg *ServerConfigure) GetContext() Context            { return cfg.context }
func (cfg *ServerConfigure) GetName() string                { return cfg.Name }
func (cfg *ServerConfigure) GetEnable() bool                { return cfg.Enable }
func (cfg *ServerConfigure) GetListenAddress() string       { return cfg.Listen }
func (cfg *ServerConfigure) GetAccessLogPath() string       { return cfg.absPath(cfg.AccessLog) }
func (cfg *ServerConfigure) GetErrorLogPath() string        { return cfg.absPath(cfg.ErrorLog) }
func (cfg *ServerConfigure) GetLogLevel() string            { return cfg.LogLevel }
func (cfg *ServerConfigure) GetIndex() string               { return cfg.absPath(cfg.Index) }
func (cfg *ServerConfigure) GetMultipartMaxMemory() int64   { return cfg.MultipartMaxMemory }
func (cfg *ServerConfigure) GetMultipartTempPath() string   { return cfg.MultipartTempPath }
func (cfg *ServerConfigure) GetProxyNextUpstreamTries() int { return cfg.ProxyNextUpstreamTries }
func (cfg *ServerConfigure) GetPProfEnable() bool           { return cfg.PProfEnable }
func (cfg *ServerConfigure) GetPProfPrefix() string         { return cfg.PProfPrefix }
func (cfg *ServerConfigure) GetEnableTLS() bool             { return cfg.EnableTLS }
func (cfg *ServerConfigure) GetTLSCertFile() string         { return cfg.TLSCertFile }
func (cfg *ServerConfigure) GetTLSKeyFile() string          { return cfg.TLSKeyFile }
func (cfg *ServerConfigure) GetConfig() *ModuleConfig       { return &ModuleConfig{m: cfg.Config} }

func (cfg *ServerConfigure) GetProxyConnectTimeout() time.Duration {
	return time.Duration(cfg.ProxyConnectTimeout) * time.Second
}

func (cfg *ServerConfigure) GetProxySendTimeout() time.Duration {
	return time.Duration(cfg.ProxySendTimeout) * time.Second
}

func (cfg *ServerConfigure) GetProxyReadTimeout() time.Duration {
	return time.Duration(cfg.ProxyReadTimeout) * time.Second
}

func (cfg *ModuleConfig) GetInt(key string) int {
	v, ok := cfg.m[key]
	if !ok {
		return 0
	}

	val, yes := v.(int)
	if yes {
		return val
	}
	return 0
}

func (cfg *ModuleConfig) GetBool(key string) bool {
	v, ok := cfg.m[key]
	if !ok {
		return false
	}

	val, yes := v.(bool)
	if yes {
		return val
	}
	return false
}

func (cfg ModuleConfig) GetString(key string) string {
	v, ok := cfg.m[key]
	if !ok {
		return ""
	}

	val, yes := v.(string)
	if yes {
		return val
	}
	return ""
}

func (cfg ModuleConfig) GetStringArray(key string) []string {
	v, ok := cfg.m[key]
	if !ok {
		return []string{}
	}

	items, ok := v.([]interface{})
	if !ok {
		return []string{}
	}

	s := []string{}
	for _, item := range items {
		v, yes := item.(string)
		if !yes {
			return []string{}
		}
		s = append(s, v)
	}
	return s
}
