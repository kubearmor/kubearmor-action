// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"path/filepath"
	"time"

	"github.com/zhy76/kubearmor-action/common"
	"github.com/zhy76/kubearmor-action/utils"

	"github.com/pkg/errors"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

const (
	ReadyStatus = "Ready"
	TRUE        = "True"
	FALSE       = "False"
)

type Client struct {
	// ClientSet is a kubernetes clientset.
	ClientSet *kubernetes.Clientset
	// DynamicClient is a dynamic client.
	DynamicClient *dynamic.DynamicClient
}

// NamespacePod is a namespace and its pods.
type NamespacePod struct {
	// Namespace object
	Namespace v1.Namespace
	// PodList is a list of Pods.
	PodList *v1.PodList
}

// EventPod contains events of a pod.
type EventPod struct {
	// Reason is the reason this event was generated.
	Reason string
	// Message is a human-readable description of the status of this operation.
	Message string
	// Count is the number of times this event has occurred.
	Count int32
	// Type of this event
	Type string
	// Action is what action was taken/failed regarding to the Regarding object.
	Action string
	// Namespace is the namespace this event applies to.
	Namespace string
}

func NewK8sClient() (*Client, error) {
	kubeconfig := filepath.Join(common.DefaultKubeConfigDir(), "config")
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	}

	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, errors.Wrap(err, "failed to build kube config")
	}

	// create the clientset
	clientSet, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create clientset")
	}

	// create the dynamic client
	dynamicClient, err := dynamic.NewForConfig(config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create dynamic client")
	}

	return &Client{
		ClientSet:     clientSet,
		DynamicClient: dynamicClient,
	}, nil
}

// Creates a namespace.
func (c *Client) CreateNamespace(name string) error {
	// create a namespace
	namespace := &v1.Namespace{
		ObjectMeta: metav1.ObjectMeta{
			Name: name,
		},
	}
	_, err := c.ClientSet.CoreV1().Namespaces().Create(context.Background(), namespace, metav1.CreateOptions{})
	return err
}

func (c *Client) DeleteNamespace(name string) error {
	return c.ClientSet.CoreV1().Namespaces().Delete(context.Background(), name, metav1.DeleteOptions{})
}

// listNamespaces returns a list of all namespaces.
func (c *Client) listNamespaces() (*v1.NamespaceList, error) {
	namespaceList, err := c.ClientSet.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to get namespaces")
	}
	return namespaceList, nil
}

// ListAllNamespacesPods returns a list of all namespaces and pods.
func (c *Client) ListAllNamespacesPods() ([]*NamespacePod, error) {
	namespaceList, err := c.listNamespaces()
	if err != nil {
		return nil, err
	}
	var namespacePodList []*NamespacePod
	for _, ns := range namespaceList.Items {
		pods, err := c.ClientSet.CoreV1().Pods(ns.Name).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return nil, errors.Wrapf(err, "failed to get all namespace pods")
		}
		namespacePod := NamespacePod{
			Namespace: ns,
			PodList:   pods,
		}
		namespacePodList = append(namespacePodList, &namespacePod)
	}

	return namespacePodList, nil
}

// Check if all pods are ready
func (c *Client) CheckAllPodsReady() (bool, error) {
	namespacePodList, err := c.ListAllNamespacesPods()
	if err != nil {
		return false, err
	}
	for _, namespacePod := range namespacePodList {
		if namespacePod.PodList == nil {
			continue
		}
		// pods.Items maybe nil
		if len(namespacePod.PodList.Items) == 0 {
			continue
		}
		for _, pod := range namespacePod.PodList.Items {
			// pod.Status.ContainerStatus == nil because of pod contain initcontainer
			if len(pod.Status.ContainerStatuses) == 0 {
				continue
			}
			if !pod.Status.ContainerStatuses[0].Ready {
				fmt.Printf("pod %s is not ready\n", pod.Name)
				return false, nil
			}
		}
	}
	return true, nil
}

// GetPodLog returns the log of a pod.
func (c *Client) GetPodLog(namespace, podName string) (string, error) {
	req := c.ClientSet.CoreV1().Pods(namespace).GetLogs(podName, &v1.PodLogOptions{})
	podLogs, err := req.Stream(context.TODO())
	if err != nil {
		return "", err
	}

	defer func() {
		_ = podLogs.Close()
	}()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, podLogs)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// GetPodEvents returns the events of a pod.
func (c *Client) GetPodEvents(namespace, podName string) ([]v1.Event, error) {
	events, err := c.ClientSet.CoreV1().
		Events(namespace).
		List(context.TODO(), metav1.ListOptions{FieldSelector: "involvedObject.name=" + podName, TypeMeta: metav1.TypeMeta{Kind: "Pod"}})
	if err != nil {
		return nil, err
	}
	return events.Items, nil
}

// GetPodReadyStatus returns the ready status of a pod.
func (c *Client) getPodReadyStatus(pod v1.Pod) bool {
	for _, condition := range pod.Status.Conditions {
		if condition.Type != ReadyStatus {
			continue
		}
		if condition.Status == TRUE {
			return true
		}
	}
	return false
}

// GetNotReadyPodEvent returns the events of not ready pods.
func (c *Client) GetNotReadyPodEvent() (map[string][]EventPod, error) {
	namespacePodList, err := c.ListAllNamespacesPods()
	if err != nil {
		return nil, err
	}
	result := make(map[string][]EventPod)
	for _, podNamespace := range namespacePodList {
		for _, pod := range podNamespace.PodList.Items {
			if c.getPodReadyStatus(pod) {
				continue
			}
			events, err := c.GetPodEvents(podNamespace.Namespace.Name, pod.Name)
			if err != nil {
				return nil, err
			}
			var eventpods []EventPod
			for _, event := range events {
				eventpods = append(eventpods, EventPod{
					Reason:    event.Reason,
					Message:   event.Message,
					Count:     event.Count,
					Type:      event.Type,
					Action:    event.Action,
					Namespace: event.Namespace,
				})
			}
			result[pod.Name] = append(result[pod.Name], eventpods...)
		}
	}
	return result, nil
}

// If exist pod not ready, show pod events and logs
func (c *Client) OutputNotReadyPodInfo() error {
	podEvents, err := c.GetNotReadyPodEvent()
	if err != nil {
		return err
	}
	for podName, events := range podEvents {
		fmt.Println("=========================================================================================================================================")
		fmt.Println("PodName: " + podName)
		fmt.Println("******************************************************Events*****************************************************************************")
		var namespace string
		for _, event := range events {
			namespace = event.Namespace
			fmt.Printf("Reason: %s\n", event.Reason)
			fmt.Printf("Message: %s\n", event.Message)
			fmt.Printf("Count: %v\n", event.Count)
			fmt.Printf("Type: %s\n", event.Type)
			fmt.Printf("Action: %s\n", event.Action)
			fmt.Println("------------------------------------------------------------------------------------------------------------------------------------")
		}
		log, err := c.GetPodLog(namespace, podName)
		if err != nil {
			return err
		}
		fmt.Println("********************************************************Log*****************************************************************************")
		fmt.Println(log)
		fmt.Println("=========================================================================================================================================")
	}
	return nil
}

func (c *Client) WaitAllPodRunning() error {
	time.Sleep(30 * time.Second)
	err := utils.Retry(10, 5*time.Second, func() error {
		flag, err := c.CheckAllPodsReady()
		if err != nil {
			return err
		}
		if !flag {
			err = c.OutputNotReadyPodInfo()
			if err != nil {
				return err
			}
			return fmt.Errorf("pods not ready")
		}
		return nil
	})
	return err
}
