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

package core

import (
	"io"
	"os"
	"path"

	"github.com/opencurve/pigeon/internal/http"
	"github.com/opencurve/pigeon/internal/utils"
)

type Pigeon struct {
	prefix string

	servers []*http.HTTPServer

	in  io.Reader
	out io.Writer
	err io.Writer
}

func NewPigeon(servers []*http.HTTPServer) *Pigeon {
	return &Pigeon{
		prefix: utils.GetBinaryDir(),

		servers: servers,

		in:  os.Stdin,
		out: os.Stdout,
		err: os.Stderr,
	}
}

func (p *Pigeon) Version() string                  { return Version }
func (pigeon *Pigeon) In() io.Reader               { return pigeon.in }
func (pigeon *Pigeon) Out() io.Writer              { return pigeon.out }
func (pigeon *Pigeon) Err() io.Writer              { return pigeon.err }
func (pigeon *Pigeon) Servers() []*http.HTTPServer { return pigeon.servers }
func (pigeon *Pigeon) SetPrefix(prefix string)     { pigeon.prefix = prefix }
func (pigeon *Pigeon) GetPrefix() string           { return pigeon.prefix }
func (pigeon *Pigeon) DefaultConfFile() string     { return path.Join(pigeon.prefix, "conf/pigeon.yaml") }
