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
 * Created Date: 2022-09-21
 * Author: Jingli Chen (Wine93)
 */

package configure

import (
	"path"
	"path/filepath"
	"time"

	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
)

type Context struct {
	Version string
	Prefix  string
}

/*
 * global:
 *   pid: pigeon.pid
 *   access_log: pigeon_access.log
 *   error_log: pigeon_error.log
 *   log_level: error
 *   index: html
 *   multipart_max_memory: 8388608
 *   multipart_temp_path: /dev/shm
 *   proxy_connect_timeout: 3
 *   proxy_send_timeout: 60
 *   proxy_read_timeout: 60
 *   proxy_next_upstream_tries: 0
 *   config:
 *     enable: true
 *
 * servers:
 *   - name: server1
 *     listen: 127.0.0.1:8000
 *   - name: server2
 *     listen: 127.0.0.1:8001
 *     config:
 *       enable: false
 *
 * upstreams:
 *   - name: upstream1
 *     check_interval: 1
 *     servers:
 *        - 127.0.0.1:9000
 *        - 127.0.0.1:9001
 *   - name: upstream2
 *     servers:
 *        - 127.0.0.1:9000
 *        - 127.0.0.1:9001
 */
type (
	Global struct {
		PidPath      string `mapstructure:"pid" default:"logs/pigeon.pid"`
		CloseTimeout int64  `mapstructure:"close_timeout" default:"60"`
		AbortTimeout int64  `mapstructure:"abort_timeout" default:"60"`

		AccessLog string `mapstructure:"access_log" default:"logs/pigeon_access.log"`
		ErrorLog  string `mapstructure:"error_log" default:"logs/pigeon_error.log"`
		LogLevel  string `mapstructure:"log_level" default:"error"`
		Index     string `mapstructure:"index" default:"html"`

		MultipartMaxMemory int64  `mapstructure:"multipart_max_memory" default:"1048576"`
		MultipartTempPath  string `mapstructure:"multipart_temp_path" default:"/tmp"`

		ProxyConnectTimeout    int `mapstructure:"proxy_connect_timeout" default:"3"`
		ProxySendTimeout       int `mapstructure:"proxy_send_timeout" default:"60"`
		ProxyReadTimeout       int `mapstructure:"proxy_read_timeout" default:"60"`
		ProxyNextUpstreamTries int `mapstructure:"proxy_next_upstream_tries" default:"0"`

		Config map[string]interface{} `mapstructure:"config"`
	}

	Server struct {
		context Context

		Name      string `mapstructure:"name" default:"localhost"`
		Enable    bool   `mapstructure:"enable" default:"true"`
		Listen    string `mapstructure:"listen" default:":8000"`
		AccessLog string `mapstructure:"access_log" default:"logs/pigeon_access.log"`
		ErrorLog  string `mapstructure:"error_log" default:"logs/pigeon_error.log"`
		LogLevel  string `mapstructure:"log_level" default:"error"`
		Index     string `mapstructure:"index" default:"html"`

		MultipartMaxMemory int64  `mapstructure:"multipart_max_memory" default:"1048576"`
		MultipartTempPath  string `mapstructure:"multipart_temp_path" default:"/tmp"`

		ProxyConnectTimeout    int `mapstructure:"proxy_connect_timeout" default:"3"`
		ProxySendTimeout       int `mapstructure:"proxy_send_timeout" default:"60"`
		ProxyReadTimeout       int `mapstructure:"proxy_read_timeout" default:"60"`
		ProxyNextUpstreamTries int `mapstructure:"proxy_next_upstream_tries" default:"0"`

		PProfEnable bool   `mapstructure:"pprof_enable" default:"false"`
		PProfPrefix string `mapstructure:"pprof_prefix" default:"/debug/pprof"`

		Config map[string]interface{} `mapstructure:"config"`
	}

	Upstream struct {
		Name    string   `mapstructure:"name"`
		Servers []string `mapstructure:"servers"`
	}

	Configure struct {
		context Context

		Global Global `mapstructure:"global"`

		Servers []Server `mapstructure:"servers"`

		Upstreams []Upstream `mapstructure:"upstreams"`
	}
)

func Default(ctx Context) *Configure {
	cfg := &Configure{context: ctx}
	defaults.SetDefaults(&cfg.Global)
	return cfg
}

func Parse(filename string, ctx Context) (*Configure, error) {
	parser := viper.NewWithOptions(viper.KeyDelimiter("::"))
	parser.SetConfigFile(filename)
	parser.SetConfigType("yaml")
	err := parser.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Configure{context: ctx}
	err = parser.Unmarshal(cfg)
	if err != nil {
		return nil, err
	}

	defaults.SetDefaults(&cfg.Global)
	for i := range cfg.Servers {
		server := &cfg.Servers[i]
		cfg.merge(server)
		server.context = ctx
		defaults.SetDefaults(server)
	}
	return cfg, nil
}

func newIfNil(config map[string]interface{}) map[string]interface{} {
	if config == nil {
		return map[string]interface{}{}
	}
	return config
}

func (cfg *Configure) merge(server *Server) {
	global := cfg.Global
	if len(server.AccessLog) == 0 {
		server.AccessLog = global.AccessLog
	}
	if len(server.ErrorLog) == 0 {
		server.ErrorLog = global.ErrorLog
	}
	if len(server.LogLevel) == 0 {
		server.LogLevel = global.LogLevel
	}
	if len(server.Index) == 0 {
		server.Index = global.Index
	}
	if server.MultipartMaxMemory == 0 {
		server.MultipartMaxMemory = global.MultipartMaxMemory
	}
	if len(server.MultipartTempPath) == 0 {
		server.MultipartTempPath = global.MultipartTempPath
	}
	if server.ProxyConnectTimeout == 0 {
		server.ProxyConnectTimeout = global.ProxyConnectTimeout
	}
	if server.ProxySendTimeout == 0 {
		server.ProxySendTimeout = global.ProxySendTimeout
	}
	if server.ProxyReadTimeout == 0 {
		server.ProxyReadTimeout = global.ProxyReadTimeout
	}
	if server.ProxyNextUpstreamTries == 0 {
		server.ProxyNextUpstreamTries = global.ProxyNextUpstreamTries
	}

	gconfig := newIfNil(global.Config)
	sconfig := newIfNil(server.Config)
	for k, v := range gconfig {
		if sconfig[k] == nil {
			sconfig[k] = v
		}
	}
	cfg.Global.Config = gconfig
	server.Config = sconfig
}

func (cfg *Configure) absPath(filename string) string {
	if filepath.IsAbs(filename) {
		return filename
	}
	return path.Join(cfg.context.Prefix, filename)
}

func (cfg *Configure) GetPidFile() string {
	return cfg.absPath(cfg.Global.PidPath)
}

func (cfg *Configure) GetCloseTimeout() time.Duration {
	return time.Duration(cfg.Global.CloseTimeout) * time.Second
}

func (cfg *Configure) GetAbortTimeout() time.Duration {
	return time.Duration(cfg.Global.AbortTimeout) * time.Second
}

func (cfg *Configure) GetErrorLogPath() string {
	return cfg.absPath(cfg.Global.ErrorLog)
}

func (cfg *Configure) GetDefaultServer() *ServerConfigure {
	server := &Server{context: cfg.context}
	defaults.SetDefaults(server)
	return server
}

func (cfg *Configure) GetServer(name string) *ServerConfigure {
	for _, server := range cfg.Servers {
		if server.Name == name {
			return &server
		}
	}
	return nil
}

func (cfg *Configure) GetUpstreams() []Upstream {
	return cfg.Upstreams
}
