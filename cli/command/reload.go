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
	"fmt"
	"strconv"
	"syscall"

	"github.com/opencurve/pigeon/internal/core"
	cliutils "github.com/opencurve/pigeon/internal/utils"
	utils "github.com/opencurve/pigeon/internal/utils"
	"github.com/spf13/cobra"
)

type reloadOptions struct {
	filename string
}

func NewReloadCommand(pigeon *core.Pigeon) *cobra.Command {
	var options reloadOptions

	cmd := &cobra.Command{
		Use:   "reload",
		Short: "Reload pigeon",
		Args:  cliutils.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			return runReload(pigeon, options)
		},
		DisableFlagsInUseLine: true,
	}

	flags := cmd.Flags()
	flags.StringVarP(&options.filename, "conf", "c", pigeon.DefaultConfFile(), "Specify pigeon configure file")

	return cmd
}

func getPid(pigeon *core.Pigeon, filename string) (int, error) {
	cfg, err := parse(pigeon, filename)
	if err != nil {
		return 0, err
	}

	pidfile := cfg.GetPidFile()
	data, _ := utils.ReadFile(pidfile)
	return strconv.Atoi(data)
}

func runReload(pigeon *core.Pigeon, options reloadOptions) error {
	pid, err := getPid(pigeon, options.filename)
	if err != nil {
		return fmt.Errorf("read pid file failed: %v", err)
	}
	return syscall.Kill(pid, syscall.SIGUSR2)
}
