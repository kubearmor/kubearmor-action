// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package controller

import (
	"context"
	"fmt"
	"kubearmor-action/common"
	"kubearmor-action/pkg/controller/client"
	"kubearmor-action/pkg/utils"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

// A controller to manage the lifecycle of a single app

type App struct {
	obj *unstructured.Unstructured
}

// NewApp method creates a new app
func NewApp(obj *unstructured.Unstructured) *App {
	return &App{
		obj: obj,
	}
}

// CreatePod method creates a pod in k8s cluster
func (app *App) Create(c *client.Client) error {
	// TODO: Check if the namespace exists
	// Define the GroupVersionResource
	gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	// Create the pod
	_, err := c.DynamicClient.Resource(gvr).Namespace(common.AppNamespace).Create(context.Background(), app.obj, metav1.CreateOptions{})
	return err
}

func (app *App) Delete(c *client.Client) error {
	// Define the GroupVersionResource
	gvr := schema.GroupVersionResource{Version: "v1", Resource: "pods"}
	// Delete the pod
	err := c.DynamicClient.Resource(gvr).Namespace(common.AppNamespace).Delete(context.Background(), app.obj.GetName(), metav1.DeleteOptions{})
	if err != nil {
		return err
	}
	return err
}

func (app *App) WaitAllAppRunning(c *client.Client) error {
	time.Sleep(30 * time.Second)
	err := utils.Retry(10, 5*time.Second, func() error {
		flag, err := c.CheckAllPodsReadyUnderNamespace(common.AppNamespace)
		if err != nil {
			return err
		}
		if !flag {
			err = c.OutputNotReadyPodInfo()
			if err != nil {
				return err
			}
			return fmt.Errorf("app not ready")
		}
		return nil
	})
	return err
}
