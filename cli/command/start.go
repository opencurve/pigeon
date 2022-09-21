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
	"github.com/opencurve/pigeon/internal/configure"
	"github.com/opencurve/pigeon/internal/core"
	cliutils "github.com/opencurve/pigeon/internal/utils"
	"github.com/spf13/cobra"
	"golang.org/x/sync/errgroup"
)

type startOptions struct {
	filename string
	prefix   string
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
	flags.StringVarP(&options.filename, "conf", "c", "pigeon.yaml", "Specify pigeon configure file")

	return cmd
}

func runStart(pigeon *core.Pigeon, options startOptions) error {
	/*
		cfg, err := pigeon.ParseConfigure(options.filename)
		if err != nil {
			return err
		}
	*/
	cfg := &configure.Configure{}

	var g errgroup.Group
	servers := pigeon.Servers()
	for _, s := range servers {
		err := s.Init(cfg)
		if err != nil {
			return err
		}

		g.Go(func() error {
			return s.Start()
		})
	}

	return g.Wait()
}
