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
	"github.com/sethvargo/go-githubactions"
	"k8s.io/klog"
)

var (
	// PWD is the current working directory
	PWD = common.GetWorkDir() + "/pkg/visualisation/"
)

// ParseSummaryData parses the summary data and returns a slice of SummaryData objects
func ParseSummaryData(path string) []*SummaryData {
	// data is a byte array that will hold the content of the JSON file
	var data []byte

	// err is an error variable that will store any error that occurs during the reading or parsing process
	var err error

	// ReadFile reads the file from the given address and returns it as a []byte array.
	// It can handle both remote urls and local paths.
	data, err = utils.ReadFile(path) // #nosec

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
		handlePsfileSet(sd, vs, "Process")
		// Get Files Accessed
		handlePsfileSet(sd, vs, "File")
		// Get Network Data
		handleNetworkSet(sd, vs)
	}
	return vs
}

// getLabel gets the label from the summary data and appends it to the VisualSysData object
func getLabel(summaryData *SummaryData, vs *VisualSysData) {
	vs.Labels = append(vs.Labels, summaryData.Label)
}

// handlePsfileSet handles the process and file data and appends it to the VisualSysData object
func handlePsfileSet(summaryData *SummaryData, vs *VisualSysData, kind string) {
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

// handleNetworkSet handles the network data and appends it to the VisualSysData object
func handleNetworkSet(summaryData *SummaryData, vs *VisualSysData) {
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

// ConvertVndToPlantUML Convert a VisualNetworkData object to a PlantUML format and save it in net.puml
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
			if appName != "" && strings.Contains(ip, appName) {
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

// ParseNetworkData parses the summary data and returns a VisualNetworkData object
func ParseNetworkData(sdOlds, sdNews []*SummaryData, appName string) *VisualNetworkData {
	if len(sdNews) == 0 {
		return nil
	}
	nsipsOld := make(map[string][]string, 0)
	nsipsNew := make(map[string][]string, 0)
	ips := make(map[string]bool, 0)
	cdsOld := make(map[connectionKey]connectionValue, 0)
	cdsNew := make(map[connectionKey]connectionValue, 0)
	cdsMerged := make(map[connectionKey]connectionValue, 0)

	vn := &VisualNetworkData{}
	vn.NsIps = make(map[string][]string)
	// if sdOlds == nil || len(sdOlds) == 0
	if sdOlds == nil || len(sdOlds) == 0 {
		for _, sdNew := range sdNews {
			// Get Namespace Labels
			getNsIps(sdNew, nsipsNew)
			// Get Different Network Connections
			getDiffConnectionData(sdNew, cdsNew, appName)
		}

		// merge the connections
		for k, v := range cdsNew {
			// unchanged
			if v.kind == 0 {
				edge := fmt.Sprintf("[%s] -[#%s]-> [%s] : %s/%s\n", k.src, v.edgeColor, k.dst, k.protocol, k.port)
				vn.Connections = append(vn.Connections, edge)
			}
			// add ips to the ips map, to filter nsips
			ips[k.src] = true
			ips[k.dst] = true
		}

		// filter nsips by ips
		for ns, ipss := range nsipsNew {
			for _, ip := range ipss {
				if _, ok := ips[ip]; ok {
					vn.NsIps[ns] = append(vn.NsIps[ns], ip)
				}
			}
		}
		return vn
	}
	// if sdOlds != nil && len(sdOlds) != 0
	for _, sdOld := range sdOlds {
		// Get Namespace Labels
		getNsIps(sdOld, nsipsOld)
		// Get Different Network Connections
		getDiffConnectionData(sdOld, cdsOld, appName)
	}

	for _, sdNew := range sdNews {
		// Get Namespace Labels
		getNsIps(sdNew, nsipsNew)
		// Get Different Network Connections
		getDiffConnectionData(sdNew, cdsNew, appName)
	}

	// merge the connections
	for k, v := range cdsNew {
		// fmt.Println(k, v)
		if _, ok := cdsOld[k]; ok {
			// if exists in both, means unchanged
			cdsMerged[k] = connectionValue{edgeColor: v.edgeColor, kind: 0}
		} else {
			// if exists in new, means added
			cdsMerged[k] = connectionValue{edgeColor: "red", kind: 1}
		}
		// add ips to the ips map, to filter nsips
		ips[k.src] = true
		ips[k.dst] = true
	}
	// check for deleted connections
	for k := range cdsOld {
		// fmt.Println(k, v)
		if _, ok := cdsNew[k]; !ok {
			// if exists in old, but not in new, means deleted
			cdsMerged[k] = connectionValue{edgeColor: "red", kind: -1}
		}
		// add ips to the ips map, to filter nsips
		ips[k.src] = true
		ips[k.dst] = true
	}
	// write the connections to the vn
	for k, v := range cdsMerged {
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
	for ns, ipss := range nsipsOld {
		for _, ip := range ipss {
			if _, ok := ips[ip]; ok {
				vn.NsIps[ns] = append(vn.NsIps[ns], ip)
			}
		}
	}
	for ns, ipss := range nsipsNew {
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
	sysPuml := append(start, jsonData...)
	sysPuml = append(sysPuml, end...)
	// fmt.Printf("%+v\n", string(sysPuml))

	klog.Infoln("Creating PlantUML File...")
	// Create plantuml file
	fw := osi.NewFileWriter(PWD + "/sys.puml")
	err = fw.WriteFile(sysPuml)
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
	err = osi.RemoveFile(PWD + "/sys.puml")
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
func ConvertNetworkJSONToImage(jsonFileOld string, jsonFileNew string, output string, appName string) error {
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
	var sdOlds []*SummaryData
	if jsonFileOld != "" {
		klog.Infoln("Parsing Old Summary Data...")
		sdOlds = ParseSummaryData(jsonFileOld)
		// handle nil
		if sdOlds == nil {
			githubactions.Warningf("Warning: Old summary report file path is invalid!")
		}
	}
	// get new summary data from new json file
	klog.Infoln("Parsing New Summary Data...")
	sdNews := ParseSummaryData(jsonFileNew)
	if sdNews == nil {
		return fmt.Errorf("Error: New SummaryData is nil")
	}

	// parse visual network connections data from summary data
	klog.Infoln("Parsing Visual Network Connections Data...")
	vnd := ParseNetworkData(sdOlds, sdNews, appName)
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
	err = osi.RemoveFile(PWD + "/net.puml")
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
