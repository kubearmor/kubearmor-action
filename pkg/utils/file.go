// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package utils

import (
	"fmt"
	"os"

	"strings"

	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/util/yaml"
)

type FileHelper struct {
	filePath string
}

// NewFileHelper creates a new instance of the FileHelper structure
func NewFileHelper(filePath string) *FileHelper {
	return &FileHelper{
		filePath: filePath,
	}
}

// ReplaceImageName replaces the image name in the specified file
func (f *FileHelper) ReplaceImageName(imagePlaceholderName, targetImageName string) (*unstructured.Unstructured, error) {
	// Read the file content
	data, err := os.ReadFile(f.filePath)
	if err != nil {
		return &unstructured.Unstructured{}, fmt.Errorf("failed to read file: %w", err)
	}

	// Convert the file contents to json format
	jsonData, err := yaml.ToJSON(data)
	if err != nil {
		return &unstructured.Unstructured{}, fmt.Errorf("failed to convert yaml to json: %w", err)
	}

	// Parse json data into an instance of an Unstructured structure
	obj := &unstructured.Unstructured{}
	err = obj.UnmarshalJSON(jsonData)
	if err != nil {
		return &unstructured.Unstructured{}, fmt.Errorf("failed to unmarshal json: %w", err)
	}

	// Replace the image field in an Unstructured structure
	// Get spec.containers fields and convert to type []interface{}
	containers, found, err := unstructured.NestedSlice(obj.Object, "spec", "containers")
	if err != nil {
		return &unstructured.Unstructured{}, fmt.Errorf("failed to get containers field: %w", err)
	}
	if !found {
		return &unstructured.Unstructured{}, fmt.Errorf("containers field not found")
	}

	// The first container is accessed by index and converted to the map[string]interface{} type
	container, ok := containers[0].(map[string]interface{})
	if !ok {
		return &unstructured.Unstructured{}, fmt.Errorf("failed to convert container to map[string]interface{}")
	}

	// Gets the image field and converts it to string type
	imageValue, found, err := unstructured.NestedString(container, "image")
	if err != nil {
		return &unstructured.Unstructured{}, fmt.Errorf("failed to get image field: %w", err)
	}
	if !found {
		return &unstructured.Unstructured{}, fmt.Errorf("image field not found")
	}

	// Replace the value of the image field
	newImageValue := strings.Replace(imageValue, imagePlaceholderName, targetImageName, 1)

	// Set a new value for the image field
	err = unstructured.SetNestedField(container, newImageValue, "image")
	if err != nil {
		return &unstructured.Unstructured{}, fmt.Errorf("failed to set image field: %w", err)
	}

	// Set the new value of the spec.containers field
	err = unstructured.SetNestedField(obj.Object, containers, "spec", "containers")
	if err != nil {
		return &unstructured.Unstructured{}, fmt.Errorf("failed to set containers field: %w", err)
	}
	return obj, nil
}
