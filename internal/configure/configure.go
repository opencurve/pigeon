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
	"github.com/spf13/viper"
)

/*
 * global:
 *   prefix: /tmp
 *   access_log: pigeon_access.log
 *   error_log: pigeon_error.log
 *
 * servers:
 *   - name: server1
 *     listen: 127.0.0.1:8000
 *   - name: server2
 *     listen: 127.0.0.1:8001
 *
 * upstreams:
 *   - name: upstream1
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
		Prefix    string `mapstructure:"prefix"`
		AccessLog string `mapstructure:"access_log"`
		ErrorLog  string `mapstructure:"error_log"`
	}

	Server struct {
		Name      string `mapstructure:"name"`
		Listen    string `mapstructure:"listen"`
		AccessLog string `mapstructure:"access_log"`
		ErrorLog  string `mapstructure:"error_log"`
	}

	Upstream struct {
		Name    string   `mapstructure:"name"`
		Servers []string `mapstructure:"servers"`
	}

	Configure struct {
		Global Global `mapstructure:"global"`

		Servers   []Server   `mapstructure:"servers"`
		Upstreams []Upstream `mapstructure:"upstreams"`
	}
)

func Parse(filename string) (*Configure, error) {
	parser := viper.NewWithOptions(viper.KeyDelimiter("::"))
	parser.SetConfigFile(filename)
	parser.SetConfigType("yaml")
	err := parser.ReadInConfig()
	if err != nil {
		return nil, err
	}

	cfg := &Configure{}
	err = parser.Unmarshal(cfg)
	return cfg, err
}
