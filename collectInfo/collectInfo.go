package collectInfo

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

type output struct {
	CustomerName string `json:"customerName"`
	UcsdInfo     struct {
		NodeType       string   `json:"nodeType"`
		CurrentVersion string   `json:"currentVersion"`
		BaseVersion    string   `json:"baseVersion"`
		UpgradeHistory []string `json:"upgradeHistory"`
		VCPU           string   `json:"vCPU"`
		RAM            string   `json:"ram"`
	} `json:"ucsdInfo"`
	UcsdObjects struct {
		TotalVMs    int `json:"totalVMs"`
		ActiveVMs   int `json:"activeVMs"`
		Clouds      int `json:"clouds"`
		Groups      int `json:"groups"`
		Users       int `json:"users"`
		Vdc         int `json:"vdc"`
		Catalogs    int `json:"catalogs"`
		TotalEsx    int `json:"totalEsx"`
		ActiveEsx   int `json:"activeEsx"`
		HypervHosts int `json:"hypervHosts"`
	} `json:"ucsdObjects"`
	UcsdInfa []struct {
		AccountName string `json:"accountName"`
		AccounType  string `json:"accounType"`
		Version     string `json:"version"`
		Model       string `json:"model"`
		Category    string `json:"category"`
		PodName     string `json:"podName"`
	} `json:"ucsdInfa"`
	ConnectorPacks []struct {
		Name              string `json:"name"`
		Version           string `json:"version"`
		DownloadedVersion string `json:"downloadedVersion"`
	} `json:"connectorPacks"`
	TasksInfo []struct {
		Name  string `json:"name"`
		Type  string `json:"type"`
		Count int    `json:"count"`
	} `json:"tasksInfo"`
}

func CollectInformation(customerName, summaryFile, diagsFile, outputFile string) {

	var s = new(output)

	s.CustomerName = customerName

	updateSummaryInfo(summaryFile, s)

	updateDiagsInfo(diagsFile, s)
	responseA, _ := json.Marshal(s)

	log.Println("Writing output to " + outputFile)
	file, err := os.Create(outputFile)
	defer file.Close()
	if err == nil {
		file.Write(responseA)
	}

	log.Println(string(responseA))

}

func updateSummaryInfo(summaryFileName string, s *output) {
	dat, err := ioutil.ReadFile(summaryFileName)
	if err != nil {
		fmt.Println(err.Error())
	}

	inputString := string(dat)

	i := strings.Split(inputString, "\n")

	for index, val := range i {
		if strings.Contains(val, ":") {
			tokens := strings.Split(val, ":")
			if tokens[0] == "Node Type" {
				value := strings.TrimSpace(tokens[1])
				s.UcsdInfo.NodeType = value

			}
			if tokens[0] == "Current Version" {
				value := strings.TrimSpace(tokens[1])
				s.UcsdInfo.CurrentVersion = value
			}

			if tokens[0] == "Base Version" {
				value := strings.TrimSpace(tokens[1])
				s.UcsdInfo.BaseVersion = value

			}
			if tokens[0] == "Upgrade History" {
				var history []string
				upgradeHistory := strings.TrimSpace(i[index+1])
				for len(upgradeHistory) != 0 {
					history = append(history, upgradeHistory)
					upgradeHistory = i[index+1]
					index = index + 1
					upgradeHistory = strings.TrimSpace(i[index+1])
				}
				s.UcsdInfo.UpgradeHistory = history
			}
			if tokens[0] == "vCPU, RAM" {
				value := strings.TrimSpace(tokens[1])
				fmt.Println(value)
				if len(value) != 0 && strings.Contains(value, ",") {
					cpuRAM := strings.Split(value, ",")
					vcpu := strings.TrimSpace(cpuRAM[0])
					ram := strings.TrimSpace(cpuRAM[1])
					s.UcsdInfo.VCPU = vcpu
					s.UcsdInfo.RAM = ram
				}
			}
		}
	}
}

func updateDiagsInfo(diagsFileName string, s *output) {
	dat, err := ioutil.ReadFile(diagsFileName)
	if err != nil {
		fmt.Println(err.Error())
	}

	inputString := string(dat)

	i := strings.Split(inputString, "\n")

	for _, val := range i {

		if strings.Contains(val, ":") {
			tokens := strings.Split(val, ":")
			if tokens[0] == "Total VMs" {
				value := strings.TrimSpace(tokens[1])

				intVal, err := strconv.Atoi(value)
				if err == nil {
					s.UcsdObjects.TotalVMs = intVal
				}

			}
			if tokens[0] == "Active VMs" {
				value := strings.TrimSpace(tokens[1])

				intVal, err := strconv.Atoi(value)
				if err == nil {
					s.UcsdObjects.ActiveVMs = intVal
				}
			}

			if tokens[0] == "Number of Clouds" {
				value := strings.TrimSpace(tokens[1])

				intVal, err := strconv.Atoi(value)
				if err == nil {
					s.UcsdObjects.Clouds = intVal
				}

			}
			if tokens[0] == "Number of Groups" {
				value := strings.TrimSpace(tokens[1])

				intVal, err := strconv.Atoi(value)
				if err == nil {
					s.UcsdObjects.Groups = intVal
				}

			}

			if tokens[0] == "Number of Groups" {
				value := strings.TrimSpace(tokens[1])

				intVal, err := strconv.Atoi(value)
				if err == nil {
					s.UcsdObjects.Groups = intVal
				}

			}

			if tokens[0] == "Number of Users" {
				value := strings.TrimSpace(tokens[1])

				intVal, err := strconv.Atoi(value)
				if err == nil {
					s.UcsdObjects.Users = intVal
				}

			}

			if tokens[0] == "Number of vDCs" {
				value := strings.TrimSpace(tokens[1])

				intVal, err := strconv.Atoi(value)
				if err == nil {
					s.UcsdObjects.Vdc = intVal
				}

			}

			if tokens[0] == "Number of Catalogs" {
				value := strings.TrimSpace(tokens[1])

				intVal, err := strconv.Atoi(value)
				if err == nil {
					s.UcsdObjects.Catalogs = intVal
				}
			}

			if tokens[0] == "VMware Total ESX Servers" {
				value := strings.TrimSpace(tokens[1])

				intVal, err := strconv.Atoi(value)
				if err == nil {
					s.UcsdObjects.TotalEsx = intVal
				}
			}

			if tokens[0] == "VMware Active ESX Servers" {
				value := strings.TrimSpace(tokens[1])

				intVal, err := strconv.Atoi(value)
				if err == nil {
					s.UcsdObjects.ActiveEsx = intVal
				}
			}

			if tokens[0] == "HyperV Total Hosts" {
				value := strings.TrimSpace(tokens[1])

				intVal, err := strconv.Atoi(value)
				if err == nil {
					s.UcsdObjects.HypervHosts = intVal
				}
			}
		}
	}
}
