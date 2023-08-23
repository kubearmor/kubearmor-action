// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package visualisation

// SummaryData Structure
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

// ProcessData Structure
type ProcessData struct {
	Source      string `json:"Source,omitempty"`
	Destination string `json:"Destination,omitempty"`
	Count       string `json:"Count,omitempty"`
	UpdatedTime string `json:"UpdatedTime,omitempty"`
	Status      string `json:"Status,omitempty"`
}

// FileData Structure
type FileData struct {
	Source      string `json:"Source,omitempty"`
	Destination string `json:"Destination,omitempty"`
	Count       string `json:"Count,omitempty"`
	UpdatedTime string `json:"UpdatedTime,omitempty"`
	Status      string `json:"Status,omitempty"`
}

// IngressConnection Structure
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

// EgressConnection Structure
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

// VisualSysData Structure
type VisualSysData struct {
	Name        string                       `json:"Name"`
	Namespace   string                       `json:"Namespace"`
	AppName     string                       `json:"AppName"`
	Labels      []string                     `json:"Labels,omitempty"`
	ProcessData map[string]map[string]string `json:"Process,omitempty"`
	FileData    map[string]map[string]string `json:"File,omitempty"`
	NetworkData map[string]map[string]string `json:"Network,omitempty"`
}

// VisualNetworkData Structure
type VisualNetworkData struct {
	/**
	 * 1. NsIps, eg.:
	 	package "namespace: default" {
			[pod/sd-ran-consensus-1] #Lightblue
			[pod/sd-ran-consensus-2] #Lightblue
	 	}
	*/
	NsIps map[string][]string // Array of :namespace -> [Ip1 #color, Ip2 #color, ...]
	/**
	 * 2. Connections, eg.:
	 	[pod/sd-ran-consensus-1] -[#blue]-> [pod/sd-ran-consensus-2] : TCP/8080
		[pod/calico-node-ztkhd] -[#blue]-> [pod/sd-ran-consensus-2] : TCP/3550
	*/
	Connections []string // Array of: [source] -[#blue]-> [destination] : protocol/port
}
