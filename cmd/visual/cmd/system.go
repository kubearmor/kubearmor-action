// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package cmd

import (
	"fmt"
	"path/filepath"

	visual "github.com/kubearmor/kubearmor-action/pkg/visualisation"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var systemCmd = &cobra.Command{
	Use:     "system",
	Short:   "system subcommand is a command to visualization system behaviors.",
	Example: "visual system -f [json file name] --app [app name] -o [png file name]",
	Run: func(cmd *cobra.Command, args []string) {
		b := cmd.Flags().Changed("file")
		if b == false {
			klog.Fatalf("Error: 'file' flag is not set")
		}
		fmt.Println("file:", jsonFile)
		jsonFile, err := filepath.Abs(jsonFile)
		if err != nil {
			klog.Fatalf("Error: getting absolute path of 'file' flag: %v", err)
		}
		fmt.Println("app name:", appName)
		err = visual.ConvertSysJSONToImage(jsonFile, sysOutput, appName)
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(systemCmd)

	flags := systemCmd.PersistentFlags()
	flags.StringVarP(&jsonFile, "file", "f", "", "karmor summary JSON file name")
	flags.StringVarP(&appName, "app", "", "", "filter app name, if you want to visualize specific app")
	flags.StringVarP(&sysOutput, "output", "o", "sys.png", "output image file name")

	if err := systemCmd.MarkPersistentFlagRequired("file"); err != nil {
		klog.Fatalf("Error: marking 'file' flag as required: %v", err)
	}
}
