// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package cmd

import (
	"fmt"
	"path/filepath"

	visual "github.com/kubearmor/kubearmor-action/pkg/visualisation"
	"github.com/kubearmor/kubearmor-action/utils"
	"github.com/spf13/cobra"
	"k8s.io/klog"
)

var networkCmd = &cobra.Command{
	Use:     "network",
	Short:   "network subcommand is a command to visualization network connection behaviors differences.",
	Example: "visual network --old [old json file name] --new [new json file name] -app [app name] -o [png file name]",
	Run: func(cmd *cobra.Command, args []string) {
		a := cmd.Flags().Changed("old")
		if a == false {
			klog.Fatalf("Error: 'old' flag is not set")
		}
		b := cmd.Flags().Changed("new")
		if b == false {
			klog.Fatalf("Error: 'new' flag is not set")
		}

		fmt.Println("old file:", oldFile)
		// Check is URL
		var err error
		if !utils.CheckIsURL(oldFile) {
			oldFile, err = filepath.Abs(oldFile)
			if err != nil {
				klog.Fatalf("Error: getting absolute path of 'oldFile' flag: %v", err)
			}
		}

		fmt.Println("new file:", newFile)
		// Check is URL
		if !utils.CheckIsURL(newFile) {
			newFile, err = filepath.Abs(newFile)
			if err != nil {
				klog.Fatalf("Error: getting absolute path of 'newFile' flag: %v", err)
			}
		}

		fmt.Println("app name:", appName)
		err = visual.ConvertNetworkJSONToImage(oldFile, newFile, netOutput, appName)
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)

	flags := networkCmd.PersistentFlags()
	flags.StringVarP(&oldFile, "old", "", "", "old karmor summary JSON file name")
	flags.StringVarP(&newFile, "new", "", "", "new karmor summary JSON file name")
	flags.StringVarP(&appName, "app", "", "", "filter app name, if you want to visualize specific app")
	flags.StringVarP(&netOutput, "output", "o", "net.png", "output image file name")

	if err := networkCmd.MarkPersistentFlagRequired("old"); err != nil {
		klog.Fatalf("Error: marking 'old' flag as required: %v", err)
	}
	if err := networkCmd.MarkPersistentFlagRequired("new"); err != nil {
		klog.Fatalf("Error: marking 'new' flag as required: %v", err)
	}
}
