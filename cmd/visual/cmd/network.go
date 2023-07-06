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

var networkCmd = &cobra.Command{
	Use:     "network",
	Short:   "network subcommand is a command to visualization network connection behaviors.",
	Example: "visual network -f [json file name] -o [png file name]",
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
		err = visual.ConvertNetworkJSONToImage(jsonFile, output)
		if err != nil {
			fmt.Println("Error:", err)
		}
	},
}

func init() {
	rootCmd.AddCommand(networkCmd)

	flags := networkCmd.PersistentFlags()
	flags.StringVarP(&jsonFile, "file", "f", "", "karmor summary JSON file name")
	flags.StringVarP(&output, "output", "o", "net.png", "PNG file name")

	if err := networkCmd.MarkPersistentFlagRequired("file"); err != nil {
		klog.Fatalf("Error: marking 'file' flag as required: %v", err)
	}
}
