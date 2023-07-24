// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package main

import (
	"github.com/kubearmor/kubearmor-action/pkg/controller/client"

	"github.com/sethvargo/go-githubactions"
)

func main() {
	action := githubactions.New()

	// Create the k8s client
	client, err := client.NewK8sClient()
	if err != nil {
		action.Fatalf("failed to create k8s client: %v", err)
		return
	}

	// Wait for all pods to be running
	action.Infof("Wait for all pods to be running...")
	err = client.WaitAllPodRunning()
	if err != nil {
		action.Fatalf("failed to wait for all pods to be running: %v", err)
		return
	}
}
