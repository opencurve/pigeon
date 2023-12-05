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
	"fmt"
	"io"
	nethttp "net/http"
	"os"
	"path"
	"strconv"
	"syscall"

	"github.com/Wine93/grace/gracehttp"
	"github.com/opencurve/pigeon/internal/configure"
	"github.com/opencurve/pigeon/internal/http"
	"github.com/opencurve/pigeon/internal/utils"
	"github.com/sevlyar/go-daemon"
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

func (pigeon *Pigeon) Shutdown() {
	for _, server := range pigeon.servers {
		server.Shutdown()
	}
}

func (pigeon *Pigeon) Start(filename string) (err error) {
	// 1. parse configure file
	cfg, err := pigeon.parse(filename)
	if err != nil {
		return err
	}

	// 2. init server by configure
	servers := []*nethttp.Server{}
	for _, server := range pigeon.Servers() {
		err := server.Init(cfg)
		if err != nil {
			return err
		} else if !server.Enable() {
			continue
		}
		servers = append(servers, server.Server())
	}

	// 3. start a daemon
	context := &daemon.Context{
		PidFileName: cfg.GetPidFile(),
		LogFileName: cfg.GetErrorLogPath(),
	}
	child, _ := context.Reborn() // NOTE: it only run once
	if child != nil {            // parent process
		return nil
	}

	// 4. write pid to file
	fi, err := os.OpenFile(cfg.GetPidFile(), os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	pidFile := daemon.NewLockFile(fi)
	err = pidFile.WritePid()
	if err != nil {
		return err
	}

	// 5. start server in child process
	defer func() { pigeon.Shutdown() }()
	err = gracehttp.ServeWithOptions(servers,
		gracehttp.StopTimeout(cfg.GetCloseTimeout()),
		gracehttp.KillTimeout(cfg.GetAbortTimeout()))
	return err
}

func (pigeon *Pigeon) parse(filename string) (*configure.Configure, error) {
	ctx := configure.Context{
		Version: pigeon.Version(),
		Prefix:  pigeon.GetPrefix(),
	}
	if !utils.FileExist(filename) {
		return configure.Default(ctx), nil
	}
	return configure.Parse(filename, ctx)
}

func (pigeon *Pigeon) Stop(filename string) error {
	pid, err := pigeon.getPid(filename)
	if err != nil {
		return fmt.Errorf("read pid file failed: %v", err)
	}
	return syscall.Kill(pid, syscall.SIGTERM)
}

func (pigeon *Pigeon) getPid(filename string) (int, error) {
	cfg, err := pigeon.parse(filename)
	if err != nil {
		return 0, err
	}

	pidfile := cfg.GetPidFile()
	data, _ := utils.ReadFile(pidfile)
	return strconv.Atoi(data)
}

func (pigeon *Pigeon) Reload(filename string) error {
	pid, err := pigeon.getPid(filename)
	if err != nil {
		return fmt.Errorf("read pid file failed: %v", err)
	}
	return syscall.Kill(pid, syscall.SIGUSR2)
}
