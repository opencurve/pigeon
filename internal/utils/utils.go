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

package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func FirstOne(name ...string) string {
	if len(name) == 0 {
		return ""
	}
	return name[0]
}

func Str2Int(s string) (int, bool) {
	v, err := strconv.Atoi(s)
	return v, err == nil
}

func Str2Bool(s string) (bool, bool) { // value, ok
	v, err := strconv.ParseBool(s)
	return v, err == nil
}

func UnixSec() int64 {
	return time.Now().UTC().Unix()
}

func UnixMilli() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func MakeArgument(args map[string]string) string {
	out := []string{}
	for k, v := range args {
		out = append(out, strings.Join([]string{k, v}, "="))
	}
	return strings.Join(out, "&")
}

func GetBinaryDir() string {
	ex, err := os.Executable()
	if err != nil {
		return "/tmp"
	}
	return filepath.Dir(ex)
}

func Choose(ok bool, first, second string) string {
	if ok {
		return first
	}
	return second
}

func FileExist(filepath string) bool {
	_, err := os.Stat(filepath)
	return err == nil
}

func ReadFile(filename string) (string, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func WriteFile(filename, data string, mode int) error {
	file, err := os.OpenFile(filename, os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.FileMode(mode))
	if err != nil {
		return err
	}
	defer file.Close()

	n, err := file.WriteString(data)
	if err != nil {
		return err
	} else if n != len(data) {
		return fmt.Errorf("write abort")
	}

	return nil
}
