// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package visualisation

type SummaryData struct {
	DeploymentName    string              `json:"DeploymentName"`
	PodName           string              `json:"PodName"`
	ClusterName       string              `json:"ClusterName"`
	Namespace         string              `json:"Namespace"`
	Label             string              `json:"Label"`
	ProcessData       []ProcessData       `json:"ProcessData,omitempty"`
	FileData          []FileData          `json:"FileData,omitempty"`
	IngressConnection []IngressConnection `json:"IngressConnection,omitempty"`
	EgressConnection  []EgressConnection  `json:"EgressConnection,omitempty"`
}

type ProcessData struct {
	Source      string `json:"Source,omitempty"`
	Destination string `json:"Destination,omitempty"`
	Count       string `json:"Count,omitempty"`
	UpdatedTime string `json:"UpdatedTime,omitempty"`
	Status      string `json:"Status,omitempty"`
}

type FileData struct {
	Source      string `json:"Source,omitempty"`
	Destination string `json:"Destination,omitempty"`
	Count       string `json:"Count,omitempty"`
	UpdatedTime string `json:"UpdatedTime,omitempty"`
	Status      string `json:"Status,omitempty"`
}

type IngressConnection struct {
	Protocol    string `json:"Protocol,omitempty"`
	Command     string `json:"Command,omitempty"`
	IP          string `json:"IP,omitempty"`
	Port        string `json:"Port,omitempty"`
	Labels      string `json:"Labels,omitempty"`
	Namespace   string `json:"Namespace,omitempty"`
	Count       string `json:"Count,omitempty"`
	UpdatedTime string `json:"UpdatedTime,omitempty"`
}

type EgressConnection struct {
	Protocol    string `json:"Protocol,omitempty"`
	Command     string `json:"Command,omitempty"`
	IP          string `json:"IP,omitempty"`
	Port        string `json:"Port,omitempty"`
	Labels      string `json:"Labels,omitempty"`
	Namespace   string `json:"Namespace,omitempty"`
	Count       string `json:"Count,omitempty"`
	UpdatedTime string `json:"UpdatedTime,omitempty"`
}

type VisualSysData struct {
	Name        string                       `json:"Name"`
	Namespace   string                       `json:"Namespace"`
	Labels      []string                     `json:"Labels,omitempty"`
	ProcessData map[string]map[string]string `json:"Process,omitempty"`
	FileData    map[string]string            `json:"File,omitempty"`
	NetworkData map[string]map[string]string `json:"Network,omitempty"`
}

type VisualNetworkData struct {
	/**
	 * 1. NsLabels, eg.:
	 	package "namespace: default" {
			[app: emailservice] #Lightblue
			[app: paymentservice] #Lightblue
	 	}
	*/
	NsLabels []map[string][]string // Array of :namespace -> [label1 #color, label2 #color, ...]
	/**
	 * 2. Connections, eg.:
	 	[app: checkoutservice] -[#blue]-> [app: emailservice] : TCP/8080
		[app: recommendationservice] -[#blue]-> [app: productcatalogservice] : TCP/3550
	*/
	Connections []string // Array of: [source] -[#blue]-> [destination] : protocol/port
}
