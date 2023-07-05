// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package visualisation

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/kubearmor-action/common"
	"github.com/kubearmor-action/utils"
	exe "github.com/kubearmor-action/utils/exec"
	osi "github.com/kubearmor-action/utils/os"
)

var (
	PWD = common.GetWorkDir() + "/pkg/visualisation/"
)

// SummaryData is a struct that represents the karmor summary JSON data
func ParseSummaryData(path string) []*SummaryData {
	// data is a byte array that will hold the content of the JSON file
	var data []byte

	// err is an error variable that will store any error that occurs during the reading or parsing process
	var err error

	// os.ReadFile is a function that takes a file name and returns its content as a byte array and an error value
	data, err = os.ReadFile(path) // #nosec

	// if err is not nil, it means there was an error during the reading process
	if err != nil {
		// print the error message and exit the program
		fmt.Println("Error:", err)
		return nil
	}

	// summaryData is a slice of SummaryData objects that will hold the parsed JSON data
	var summaryDatas []*SummaryData

	// json.Unmarshal is a function that takes a JSON byte array and a pointer to a variable and parses the JSON data into that variable
	err = json.Unmarshal(data, &summaryDatas)

	// if err is not nil, it means there was an error during the parsing process
	if err != nil {
		// print the error message and exit the program
		fmt.Println("Error:", err)
		return nil
	}

	// if err is nil, it means the parsing process was successful
	return summaryDatas
}

// ParseSysData parses the summary data and returns a VisualSysData object
func ParseSysData(summaryDatas []*SummaryData) *VisualSysData {
	if len(summaryDatas) == 0 {
		return nil
	}
	vs := &VisualSysData{
		Name:      "sys-" + utils.GetUUID(),
		Namespace: summaryDatas[0].Namespace,
		Labels:    make([]string, 0),
	}
	vs.ProcessData = make(map[string]map[string]string)
	vs.FileData = make(map[string]string)
	vs.NetworkData = make(map[string]map[string]string)
	for _, sd := range summaryDatas {
		getLabel(sd, vs)
		// Get Processes Produced
		handle_psfile_set(sd, vs, "Process")
		// Get Files Accessed
		handle_psfile_set(sd, vs, "File")
		// Get Network Data
		handle_network_set(sd, vs)
	}
	return vs
}

// getLabel gets the label from the summary data and appends it to the VisualSysData object
func getLabel(summaryData *SummaryData, vs *VisualSysData) {
	vs.Labels = append(vs.Labels, summaryData.Label)
}

// handle_psfile_set handles the process and file data and appends it to the VisualSysData object
func handle_psfile_set(summaryData *SummaryData, vs *VisualSysData, kind string) {
	if kind == "Process" {
		for _, ps := range summaryData.ProcessData {
			if _, ok := vs.ProcessData[ps.Source]; !ok {
				vs.ProcessData[ps.Source] = make(map[string]string)
			}
			vs.ProcessData[ps.Source][ps.Destination] = "o"
		}
	} else if kind == "File" {
		for _, file := range summaryData.FileData {
			vs.FileData[file.Destination] = "o"
		}
	}
}

// handle_network_set handles the network data and appends it to the VisualSysData object
func handle_network_set(summaryData *SummaryData, vs *VisualSysData) {
	for _, net := range summaryData.IngressConnection {
		if _, ok := vs.NetworkData[net.Protocol]; !ok {
			vs.NetworkData[net.Protocol] = make(map[string]string)
		}
		vs.NetworkData[net.Protocol][net.Command] = "o"
	}
	for _, net := range summaryData.EgressConnection {
		if _, ok := vs.NetworkData[net.Protocol]; !ok {
			vs.NetworkData[net.Protocol] = make(map[string]string)
		}
		vs.NetworkData[net.Protocol][net.Command] = "o"
	}
}

// ConvertJsonToImage converts the JSON data to a plantuml image
func ConvertJsonToImage(jsonFile string, output string) error {
	// get summary data from json file
	fmt.Println("Parsing Summary Data...")
	sd := ParseSummaryData(jsonFile)
	if sd == nil {
		return fmt.Errorf("Error: SummaryData is nil")
	}

	// parse visual sys data from summary data
	fmt.Println("Parsing Visual Sys Data...")
	vsd := ParseSysData(sd)
	if vsd == nil {
		return fmt.Errorf("Error: VisualSysData is nil")
	}
	jsonData, err := json.MarshalIndent(vsd, "", "    ")
	if err != nil {
		return err
	}
	// fmt.Printf("JSON: %+v", string(jsonData))
	start := []byte("@startjson\n")
	end := []byte("\n@endjson")
	sys_puml := append(start, jsonData...)
	sys_puml = append(sys_puml, end...)
	// fmt.Printf("%+v\n", string(sys_puml))

	fmt.Println("Cheking Dependencies...")
	// Check if java is installed
	_, b := exe.CheckCmdIsExist("java")
	if !b {
		return fmt.Errorf("Error: java not installed")
	}
	// Check if plantuml.jar is installed
	b = osi.IsFileExist(PWD + "/plantuml.jar")
	if !b {
		return fmt.Errorf("Error: plantuml.jar not installed")
	}

	fmt.Println("Creating PlantUML File...")
	// Create plantuml file
	fw := osi.NewFileWriter(PWD + "/sys.puml")
	err = fw.WriteFile(sys_puml)
	if err != nil {
		return err
	}

	fmt.Println("Creating Image...")
	s, err := exe.RunSimpleCmd("java -jar " + PWD + "./plantuml.jar " + PWD + "/sys.puml -output ./")
	fmt.Println(s)
	if err != nil {
		return err
	}

	fmt.Println("Removing PlantUML File...")
	err = osi.RemoceFile(PWD + "/sys.puml")
	if err != nil {
		return err
	}
	_, err = exe.RunSimpleCmd("mv " + PWD + "/sys.png " + common.GetWorkDir() + "/" + output)
	if err != nil {
		return err
	}
	fmt.Println("Image Created Successfully!")
	return nil
}
