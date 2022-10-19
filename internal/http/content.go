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
 * Created Date: 2022-09-27
 * Author: Jingli Chen (Wine93)
 */

package http

import (
	"io"
)

type content interface {
	data()
}

type (
	Message struct {
		message string
	}

	JSON struct {
		m interface{}
	}

	File struct {
		filename string
	}

	Reader struct {
		reader io.Reader
		size   int64
		ctype  string // Content Type
	}

	Buffer struct {
		reader io.Reader
		size   int64
	}
)

func (_ *Message) data() {}
func (_ *JSON) data()    {}
func (_ *File) data()    {}
func (_ *Reader) data()  {}
func (_ *Buffer) data()  {}
