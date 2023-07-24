// SPDX-License-Identifier: Apache-2.0
// Copyright 2023 Authors of KubeArmor

package visualisation

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/kubearmor/kubearmor-action/common"
	"github.com/kubearmor/kubearmor-action/utils"
	exe "github.com/kubearmor/kubearmor-action/utils/exec"
	osi "github.com/kubearmor/kubearmor-action/utils/os"
	"k8s.io/klog"
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
		klog.Errorf("Error: %v", err)
		return nil
	}

	// summaryData is a slice of SummaryData objects that will hold the parsed JSON data
	var summaryDatas []*SummaryData

	// json.Unmarshal is a function that takes a JSON byte array and a pointer to a variable and parses the JSON data into that variable
	err = json.Unmarshal(data, &summaryDatas)

	// if err is not nil, it means there was an error during the parsing process
	if err != nil {
		// print the error message and exit the program
		klog.Errorf("Error: %v", err)
		return nil
	}

	// if err is nil, it means the parsing process was successful
	return summaryDatas
}

// ParseSysData parses the summary data and returns a VisualSysData object
func ParseSysData(summaryDatas []*SummaryData, appName string) *VisualSysData {
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
		// filter by appName
		if appName != "" {
			if !strings.Contains(sd.PodName, appName) {
				continue
			}
			vs.AppName = appName
		}
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

// Convert a VisualNetworkData object to a PlantUML format and save it in net.puml
func ConvertVndToPlantUML(vnd *VisualNetworkData, appName string) error {
	// Create a file named net.puml
	file, err := os.Create(PWD + "net.puml") // #nosec
	if err != nil {
		return err
	}
	defer file.Close()

	// Write the @startuml tag
	_, err = file.WriteString("@startuml\n")
	if err != nil {
		return err
	}

	// Loop over the NsLabels and write the package and labels
	for ns, ips := range vnd.NsIps {
		_, err = file.WriteString(fmt.Sprintf("package \"namespace: %s\" {\n", ns))
		if err != nil {
			return err
		}
		for _, ip := range ips {
			color := "Lightblue"
			if strings.Contains(ip, appName) {
				color = "Orange"
			}
			_, err = file.WriteString(fmt.Sprintf("[%s] #%s\n", ip, color))
			if err != nil {
				return err
			}
		}
		_, err = file.WriteString("}\n")
		if err != nil {
			return err
		}
	}

	// Loop over the Connections and write the arrows and ports
	for _, conn := range vnd.Connections {
		_, err = file.WriteString(fmt.Sprintf("%s", conn))
		if err != nil {
			return err
		}
	}

	// Write the @enduml tag
	_, err = file.WriteString("@enduml\n")
	if err != nil {
		return err
	}

	return nil
}

type connectionKey struct {
	src      string
	dst      string
	protocol string
	port     string
}
type connectionValue struct {
	// red: changed, yellow: TCPv6, blue: TCP, green: UDP, grey: others
	edgeColor string
	// -1: deleted, 0: no change, 1: added
	kind int
}

func ParseNetworkData(sd_olds, sd_news []*SummaryData, appName string) *VisualNetworkData {
	if len(sd_news) == 0 {
		return nil
	}
	nsips_old := make(map[string][]string, 0)
	nsips_new := make(map[string][]string, 0)
	ips := make(map[string]bool, 0)
	cds_old := make(map[connectionKey]connectionValue, 0)
	cds_new := make(map[connectionKey]connectionValue, 0)
	cds_merged := make(map[connectionKey]connectionValue, 0)

	vn := &VisualNetworkData{}
	vn.NsIps = make(map[string][]string)
	for _, sd_old := range sd_olds {
		// Get Namespace Labels
		getNsIps(sd_old, nsips_old)
		// Get Different Network Connections
		getDiffConnectionData(sd_old, cds_old, appName)
	}

	for _, sd_new := range sd_news {
		// Get Namespace Labels
		getNsIps(sd_new, nsips_new)
		// Get Different Network Connections
		getDiffConnectionData(sd_new, cds_new, appName)
	}

	// merge the connections
	for k, v := range cds_new {
		// fmt.Println(k, v)
		if _, ok := cds_old[k]; ok {
			// if exists in both, means unchanged
			cds_merged[k] = connectionValue{edgeColor: v.edgeColor, kind: 0}
		} else {
			// if exists in new, means added
			cds_merged[k] = connectionValue{edgeColor: "red", kind: 1}
		}
		// add ips to the ips map, to filter nsips
		ips[k.src] = true
		ips[k.dst] = true
	}
	// check for deleted connections
	for k := range cds_old {
		// fmt.Println(k, v)
		if _, ok := cds_new[k]; !ok {
			// if exists in old, but not in new, means deleted
			cds_merged[k] = connectionValue{edgeColor: "red", kind: -1}
		}
		// add ips to the ips map, to filter nsips
		ips[k.src] = true
		ips[k.dst] = true
	}
	// write the connections to the vn
	for k, v := range cds_merged {
		// unchanged
		if v.kind == 0 {
			edge := fmt.Sprintf("[%s] -[#%s]-> [%s] : %s/%s\n", k.src, v.edgeColor, k.dst, k.protocol, k.port)
			vn.Connections = append(vn.Connections, edge)
		}
		// added
		if v.kind == 1 {
			edge := fmt.Sprintf("[%s] -[#%s]-> [%s] : ++%s/%s\n", k.src, v.edgeColor, k.dst, k.protocol, k.port)
			vn.Connections = append(vn.Connections, edge)
		}
		// deleted
		if v.kind == -1 {
			edge := fmt.Sprintf("[%s] -[#%s]..> [%s] : --%s/%s\n", k.src, v.edgeColor, k.dst, k.protocol, k.port)
			vn.Connections = append(vn.Connections, edge)
		}
	}
	// filter nsips by ips
	for ns, ipss := range nsips_old {
		for _, ip := range ipss {
			if _, ok := ips[ip]; ok {
				vn.NsIps[ns] = append(vn.NsIps[ns], ip)
			}
		}
	}
	for ns, ipss := range nsips_new {
		for _, ip := range ipss {
			if _, ok := ips[ip]; ok {
				vn.NsIps[ns] = append(vn.NsIps[ns], ip)
			}
		}
	}
	return vn
}

func getNsIps(summaryData *SummaryData, nsips map[string][]string) {
	if summaryData.PodName == "" {
		return
	}

	pod := "pod/" + summaryData.PodName

	// if namespace exists
	if _, ok := nsips[summaryData.Namespace]; ok {
		for _, ip := range nsips[summaryData.Namespace] {
			if ip == pod {
				return
			}
		}
		nsips[summaryData.Namespace] = append(nsips[summaryData.Namespace], pod)
		return
	}
	// if namespace does not exist
	nsips[summaryData.Namespace] = append([]string(nil), pod)
}

func getEdgeColor(protocol string) string {
	if protocol == "TCPv6" {
		return "orange"
	} else if protocol == "TCP" {
		return "blue"
	} else if protocol == "UDP" {
		return "green"
	} else {
		return "grey"
	}
}

func getDiffConnectionData(sd *SummaryData, cds map[connectionKey]connectionValue, appName string) {
	if sd.PodName == "" {
		return
	}

	for _, net := range sd.IngressConnection {
		if net.IP == "" || net.IP == common.LOCALHOST {
			continue
		}
		dst := "pod/" + sd.PodName
		src := net.IP
		ck := connectionKey{src: src, dst: dst, protocol: net.Protocol, port: net.Port}
		cv := connectionValue{edgeColor: getEdgeColor(net.Protocol), kind: 0}
		// filter by appName
		if appName != "" && !strings.Contains(src, appName) && !strings.Contains(dst, appName) {
			continue
		}
		cds[ck] = cv
	}
	for _, net := range sd.EgressConnection {
		if net.IP == "" || net.IP == common.LOCALHOST {
			continue
		}
		dst := net.IP
		src := "pod/" + sd.PodName
		ck := connectionKey{src: src, dst: dst, protocol: net.Protocol, port: net.Port}
		cv := connectionValue{edgeColor: getEdgeColor(net.Protocol), kind: 0}
		// filter by appName
		if appName != "" && !strings.Contains(src, appName) && !strings.Contains(dst, appName) {
			continue
		}
		cds[ck] = cv
	}
}

// func getConnectionData(summaryData *SummaryData, vn *VisualNetworkData) {
// 	if summaryData.Label == "" {
// 		return
// 	}
// 	for _, net := range summaryData.IngressConnection {
// 		if net.Labels == "" {
// 			continue
// 		}
// 		edgeColor := getEdgeColor(net.Protocol)
// 		dst := summaryData.Label
// 		src := net.Labels
// 		edge := fmt.Sprintf("[%s] -[#%s]-> [%s] : %s/%s\n", src, edgeColor, dst, net.Protocol, net.Port)
// 		vn.Connections = append(vn.Connections, edge)
// 	}
// 	for _, net := range summaryData.EgressConnection {
// 		if net.Labels == "" {
// 			continue
// 		}
// 		edgeColor := getEdgeColor(net.Protocol)
// 		dst := net.Labels
// 		src := summaryData.Label
// 		edge := fmt.Sprintf("[%s] -[#%s]-> [%s] : %s/%s\n", src, edgeColor, dst, net.Protocol, net.Port)
// 		vn.Connections = append(vn.Connections, edge)
// 	}
// 	vn.Connections = utils.RemoveDuplication(vn.Connections)
// }

// ConvertSysJSONToImage converts the summary system JSON data to a plantuml image
func ConvertSysJSONToImage(jsonFile string, output string, appName string) error {
	klog.Infoln("Cheking Dependencies...")
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

	// get summary data from json file
	klog.Infoln("Parsing Summary Data...")
	sd := ParseSummaryData(jsonFile)
	if sd == nil {
		return fmt.Errorf("Error: SummaryData is nil")
	}

	// parse visual sys data from summary data
	klog.Infoln("Parsing Visual System Data...")
	vsd := ParseSysData(sd, appName)
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

	klog.Infoln("Creating PlantUML File...")
	// Create plantuml file
	fw := osi.NewFileWriter(PWD + "/sys.puml")
	err = fw.WriteFile(sys_puml)
	if err != nil {
		return err
	}

	klog.Infoln("Creating Image...")
	s, err := exe.RunSimpleCmd("java -jar -DPLANTUML_LIMIT_SIZE=100000 -Xmx8096m " + PWD + "./plantuml.jar " + PWD + "/sys.puml -output ./")
	klog.Infoln(s)
	if err != nil {
		return err
	}

	klog.Infoln("Removing PlantUML File...")
	err = osi.RemoceFile(PWD + "/sys.puml")
	if err != nil {
		return err
	}
	_, err = exe.RunSimpleCmd("mv " + PWD + "/sys.png " + common.GetWorkDir() + "/" + output)
	if err != nil {
		return err
	}
	klog.Infoln("Image Created Successfully!")
	return nil
}

// ConvertNetworkJSONToImage converts the summary network JSON data to a plantuml image
func ConvertNetworkJSONToImage(jsonFile_old string, jsonFile_new string, output string, appName string) error {
	klog.Info("Cheking Dependencies...")
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

	// get old summary data from old json file
	klog.Infoln("Parsing Old Summary Data...")
	sd_olds := ParseSummaryData(jsonFile_old)
	if sd_olds == nil {
		return fmt.Errorf("Error: Old SummaryData is nil")
	}

	// get new summary data from new json file
	klog.Infoln("Parsing New Summary Data...")
	sd_news := ParseSummaryData(jsonFile_new)
	if sd_news == nil {
		return fmt.Errorf("Error: New SummaryData is nil")
	}

	// parse visual network connections data from summary data
	klog.Infoln("Parsing Visual Network Connections Data...")
	vnd := ParseNetworkData(sd_olds, sd_news, appName)
	if vnd == nil {
		return fmt.Errorf("Error: VisualNetworkData is nil")
	}

	// Create plantuml file
	klog.Infoln("Creating PlantUML File...")
	err := ConvertVndToPlantUML(vnd, appName)
	if err != nil {
		return err
	}

	klog.Infoln("Creating Image...")
	s, err := exe.RunSimpleCmd("java -jar -DPLANTUML_LIMIT_SIZE=100000 -Xmx8096m " + PWD + "./plantuml.jar " + PWD + "/net.puml -output ./")
	klog.Infoln(s)
	if err != nil {
		return err
	}

	klog.Infoln("Removing PlantUML File...")
	err = osi.RemoceFile(PWD + "/net.puml")
	if err != nil {
		return err
	}
	_, err = exe.RunSimpleCmd("mv " + PWD + "/net.png " + common.GetWorkDir() + "/" + output)
	if err != nil {
		return err
	}
	klog.Infoln("Image Created Successfully!")
	return nil
}
