// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package controller

import (
	"context"

	"github.com/zhy76/kubearmor-action/common"
	"github.com/zhy76/kubearmor-action/pkg/controller/client"

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
