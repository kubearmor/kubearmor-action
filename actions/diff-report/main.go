// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package main

import (
	"github.com/kubearmor/kubearmor-action/pkg/controller/client"
	exe "github.com/kubearmor/kubearmor-action/utils/exec"

	"github.com/sethvargo/go-githubactions"
)

func main() {
	action := githubactions.New()
	// 1. Get the inputs
	oldApp := action.GetInput("old-app")
	if oldApp == "" {
		action.Fatalf("old-app cannot be empty")
	}
	newApp := action.GetInput("new-app-url")
	if newApp == "" {
		action.Fatalf("new-app cannot be empty")
	}
	filepath := action.GetInput("filepath")

	action.Infof("oldApp: %v, newApp: %v, filepath: %v", oldApp, newApp, filepath)

	// Create the k8s client
	client, err := client.NewK8sClient()
	if err != nil {
		action.Fatalf("failed to create k8s client: %v", err)
		return
	}

	// 2. Deploy the old app
	res, err := exe.RunSimpleCmd("kubectl apply -f " + oldApp)
	action.Infof("res(apply old-app):\n %v", res)
	if err != nil {
		action.Fatalf("failed to apply old-app: %v", err)
		return
	}
	action.Infof("Apply old app successfully!")
	// Wait for the old app to be running
	action.Infof("Wait for the old app to be running...")
	err = client.WaitAllPodRunning()
	if err != nil {
		action.Fatalf("failed to wait for the old app to be running: %v", err)
		return
	}
	// Check the pods
	res, err = exe.RunSimpleCmd("kubectl get po -A")
	action.Infof("res(get pods):\n %v", res)
	if err != nil {
		action.Fatalf("failed to get pods: %v", err)
		return
	}
	// Delete the old app
	res, err = exe.RunSimpleCmd("kubectl delete -f " + oldApp)
	action.Infof("res(delete old-app):\n %v", res)
	if err != nil {
		action.Fatalf("failed to delete old-app: %v", err)
		return
	}

	// 3. Save the old app's baseline file
	// res, err = utils.RunSimpleCmd("karmor summary -n " + common.AppNamespace + " -p app-old > " + filepath + "/baseline")
	// action.Infof("res(old-app): %v", res)
	// if err != nil {
	// 	action.Fatalf("failed to run karmor summary for old-app: %v", err)
	// 	return
	// }

	// 4. Deploy the new app
	res, err = exe.RunSimpleCmd("kubectl apply -f " + newApp)
	action.Infof("res(apply new-app):\n %v", res)
	if err != nil {
		action.Fatalf("failed to apply new-app: %v", err)
		return
	}
	action.Infof("Apply new app successfully!")
	// Wait for the new app to be running
	action.Infof("Wait for the new app to be running...")
	err = client.WaitAllPodRunning()
	if err != nil {
		action.Fatalf("failed to wait for the new app to be running: %v", err)
		return
	}
	// Check the pods
	res, err = exe.RunSimpleCmd("kubectl get po -A")
	action.Infof("res(get pods):\n %v", res)
	if err != nil {
		action.Fatalf("failed to get pods: %v", err)
		return
	}
	// Delete the new app
	res, err = exe.RunSimpleCmd("kubectl delete -f " + newApp)
	action.Infof("res(delete new-app):\n %v", res)
	if err != nil {
		action.Fatalf("failed to delete new-app: %v", err)
		return
	}

	// 5. Save the new app's updated file
	// res, err = utils.RunSimpleCmd("karmor summary -n " + common.AppNamespace + " -p app-new > " + filepath + "/updated")
	// action.Infof("res(new-app): %v", res)
	// if err != nil {
	// 	action.Fatalf("failed to run karmor summary for new-app: %v", err)
	// 	return
	// }

	// 6. Compare the baseline file and updated file
	// res, err = utils.RunSimpleCmd("diff " + filepath + "/baseline " + filepath + "/updated || true")
	// action.Infof("res(diff): %v", res)
	// if err != nil {
	// 	action.Fatalf("failed to run diff: %v", err)
	// 	return
	// }
}
