// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "visual",
	Short:   "visual is a command to visualization system or network behaviors.",
	Example: "visual system -f [json file name] --app [app name] -o [png file name]\nvisual network --old [old json file name] --new [new json file name] -app [app name] -o [png file name]",
}

// Execute executes the root command.
func Execute() {
	rootCmd.Execute() // #nosec
}
