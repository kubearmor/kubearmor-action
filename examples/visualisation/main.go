// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package main

import (
	"fmt"

	"github.com/kubearmor/kubearmor-action/common"
	visual "github.com/kubearmor/kubearmor-action/pkg/visualisation"
)

func main() {
	// jsonFile is the name of the JSON file that you shared with me
	oldJSONFile := common.GetWorkDir() + "/test/testdata/old-summary-data.json"
	newJSONFile := common.GetWorkDir() + "/test/testdata/new-summary-data.json"

	// sd := visual.ParseSummaryData(jsonFile)
	// vnd := visual.ParseNetworkData(sd)
	// fmt.Println(vnd)
	appName := "wordpress"
	err := visual.ConvertNetworkJSONToImage(oldJSONFile, newJSONFile, "net.png", appName)
	if err != nil {
		fmt.Println("Network-Visualisation Error:", err)
	}
	err = visual.ConvertSysJSONToImage(newJSONFile, "sys.png", appName)
	if err != nil {
		fmt.Println("System-Visualisation Error:", err)
	}
}
