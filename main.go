// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package main

import (
	"github.com/kubearmor/kubearmor-action/common"
	ctrl "github.com/kubearmor/kubearmor-action/pkg/controller"
	"github.com/kubearmor/kubearmor-action/pkg/controller/client"
	exe "github.com/kubearmor/kubearmor-action/utils/exec"
	osi "github.com/kubearmor/kubearmor-action/utils/os"

	"github.com/sethvargo/go-githubactions"
)

func main() {
	action := githubactions.New()
	// 1. Get the inputs
	oldAppName := action.GetInput("old-app-image-name")
	if oldAppName == "" {
		action.Fatalf("old-app-image-name cannot be empty")
	}
	newAppName := action.GetInput("new-app-image-name")
	if newAppName == "" {
		action.Fatalf("new-app-image-name cannot be empty")
	}
	filepath := action.GetInput("filepath")

	action.Infof("oldAppName: %v, newAppName: %v, filepath: %v", oldAppName, newAppName, filepath)

	// Create the k8s client
	client, err := client.NewK8sClient()
	if err != nil {
		action.Fatalf("failed to create k8s client: %v", err)
		return
	}
	// Create the app namespace
	err = client.CreateNamespace(common.AppNamespace)
	if err != nil {
		action.Fatalf("failed to create namespace: %v", err)
		return
	}
	defer func() {
		err := client.DeleteNamespace(common.AppNamespace)
		if err != nil {
			action.Fatalf("failed to delete namespace: %v", err)
			return
		}
	}()
	// 2. Deploy the old app
	// Create fileHelper
	oldAppFilehelper := osi.NewFileHelper(common.OldAppTemplateFilePath)
	// Replace the old app image name
	oldAppObj, err := oldAppFilehelper.ReplaceImageName(common.OldAppImagePlaceholderName, oldAppName)
	if err != nil {
		action.Fatalf("failed to replace image name: %v", err)
		return
	}
	action.Infof("old-app: %v", oldAppObj)
	// Create the old app controller
	oldAppCtrl := ctrl.NewApp(oldAppObj)
	// Create old app
	err = oldAppCtrl.Create(client)
	if err != nil {
		action.Fatalf("failed to create app: %v", err)
		return
	}
	action.Infof("Create old app successfully!")
	// Wait for the old app to be running
	action.Infof("Wait for the old app to be running...")
	err = client.WaitAllPodRunning()
	if err != nil {
		action.Fatalf("failed to wait for the old app to be running: %v", err)
		return
	}
	// Check the pods
	res, err := exe.RunSimpleCmd("kubectl get po -A")
	action.Infof("res(get pods):\n %v", res)
	if err != nil {
		action.Fatalf("failed to get pods: %v", err)
		return
	}
	// Delete the old app
	err = oldAppCtrl.Delete(client)
	if err != nil {
		action.Fatalf("failed to delete app: %v", err)
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
	// Create fileHelper
	newAppFilehelper := osi.NewFileHelper(common.NewAppTemplateFilePath)
	// Replace the new app image name
	newAppObj, err := newAppFilehelper.ReplaceImageName(common.NewAppImagePlaceholderName, newAppName)
	if err != nil {
		action.Fatalf("failed to replace image name: %v", err)
		return
	}
	action.Infof("new-app: %v", newAppObj)
	// Create the new app controller
	newAppCtrl := ctrl.NewApp(newAppObj)
	// Create new app
	err = newAppCtrl.Create(client)
	if err != nil {
		action.Fatalf("failed to create app: %v", err)
		return
	}
	action.Infof("Create new app successfully!")
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
	// err = newAppCtrl.Delete(client)
	// if err != nil {
	// 	action.Fatalf("failed to delete app: %v", err)
	// 	return
	// }

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
