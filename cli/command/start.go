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

package command

import (
	"net/http"
	"os"

	"github.com/Wine93/grace/gracehttp"
	"github.com/opencurve/pigeon/internal/configure"
	"github.com/opencurve/pigeon/internal/core"
	cliutils "github.com/opencurve/pigeon/internal/utils"
	utils "github.com/opencurve/pigeon/internal/utils"
	daemon "github.com/sevlyar/go-daemon"
	"github.com/spf13/cobra"
)

type startOptions struct {
	filename string
}

func NewStartCommand(pigeon *core.Pigeon) *cobra.Command {
	var options startOptions

	cmd := &cobra.Command{
		Use:   "start [OPTIONS]",
		Short: "Start pigeon",
		Args:  cliutils.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runStart(pigeon, options)
		},
		DisableFlagsInUseLine: true,
	}

	flags := cmd.Flags()
	flags.StringVarP(&options.filename, "conf", "c", pigeon.DefaultConfFile(), "Specify pigeon configure file")

	return cmd
}

func parse(pigeon *core.Pigeon, filename string) (*configure.Configure, error) {
	ctx := configure.Context{
		Version: pigeon.Version(),
		Prefix:  pigeon.GetPrefix(),
	}
	if !utils.FileExist(filename) {
		return configure.Default(ctx), nil
	}
	return configure.Parse(filename, ctx)
}

func runStart(pigeon *core.Pigeon, options startOptions) error {
	// 1. parse configure file
	cfg, err := parse(pigeon, options.filename)
	if err != nil {
		return err
	}

	// 2. init server by configure
	servers := []*http.Server{}
	for _, server := range pigeon.Servers() {
		err := server.Init(cfg)
		if err != nil {
			return err
		}

		if server.Enable() {
			servers = append(servers, server.Server())
		}
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
	return gracehttp.ServeWithOptions(servers,
		gracehttp.StopTimeout(cfg.GetCloseTimeout()),
		gracehttp.KillTimeout(cfg.GetAbortTimeout()))
}
