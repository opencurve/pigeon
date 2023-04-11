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

package log

import (
	zaplog "github.com/pingcap/log"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type M struct {
	Key   string
	Value string
}

func New(level, filename string) (*zap.Logger, error) {
	logger, _, err := zaplog.InitLogger(&zaplog.Config{
		Level:            level,
		File:             zaplog.FileLogConfig{Filename: filename},
		Format:           "text",
		DisableTimestamp: false,
	}, zap.AddStacktrace(zapcore.FatalLevel))

	return logger, err
}

func Field(key string, val interface{}) zap.Field {
	switch val.(type) {
	case bool:
		return zap.Bool(key, val.(bool))
	case string:
		return zap.String(key, val.(string))
	case []byte:
		return zap.String(key, string(val.([]byte)))
	case int:
		return zap.Int(key, val.(int))
	case int64:
		return zap.Int64(key, val.(int64))
	case uint16:
		return zap.Uint16(key, val.(uint16))
	case uint32:
		return zap.Uint32(key, val.(uint32))
	case uint64:
		return zap.Uint64(key, val.(uint64))
	case float64:
		return zap.Float64(key, val.(float64))
	case error:
		return zap.String(key, val.(error).Error())
	}
	return zap.Skip()
}
