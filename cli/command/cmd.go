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

	"github.com/opencurve/pigeon/internal/core"
	cliutils "github.com/opencurve/pigeon/internal/utils"
	"github.com/spf13/cobra"
)

type rootOptions struct{}

func addSubCommands(cmd *cobra.Command, pigeon *core.Pigeon) {
	cmd.AddCommand(
		NewStartCommand(pigeon),   // pigeon start
		NewStopCommand(pigeon),    // pigeon stop
		NewRestartCommand(pigeon), // pigeon restart
	)
}

func setupRootCommand(cmd *cobra.Command) {
	cmd.SetVersionTemplate("Pigeon v{{.Version}}\n")
	cliutils.SetFlagErrorFunc(cmd)
	cliutils.SetHelpTemplate(cmd)
	cliutils.SetUsageTemplate(cmd)
}

func NewPigeonCommand(pigeon *core.Pigeon) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "pigeon [OPTIONS] COMMAND [ARGS...]",
		Short:   "An easy-to-use, flexible HTTP web framework based on Gin",
		Version: pigeon.Version(),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				return cliutils.ShowHelp(pigeon.Err())(cmd, args)
			}

			return fmt.Errorf("pigeon: '%s' is not a pigeon command.\n"+
				"See 'pigeon --help'", args[0])
		},
		SilenceUsage:          true, // silence usage when an error occurs
		DisableFlagsInUseLine: true,
	}

	cmd.Flags().BoolP("version", "v", false, "Print version information and quit")
	cmd.PersistentFlags().BoolP("help", "h", false, "Print usage")

	addSubCommands(cmd, pigeon)
	setupRootCommand(cmd)

	return cmd
}
