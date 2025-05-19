package services

import (
	"encoding/json"
	"fmt"
	"io"
	"math"
	"net/http"
	"os"
)

type JobStatus struct {
	Stop   bool   `json:"Stop"`
	Status string `json:"Status"`
}

func GetSystemHealthUpdate() (healthyCount int, unhealthyCount int, err error) {
	baseURL := os.Getenv("NOMAD_BASE_URL")
	if baseURL == "" {
		fmt.Println("NOMAD_BASE_URL is not set")
		return 0, 0, fmt.Errorf("NOMAD_BASE_URL is not set")
	}

	nomadTOKEN := os.Getenv("NOMAD_TOKEN")
	if nomadTOKEN == "" {
		fmt.Println("NOMAD_TOKEN is not set")
		return 0, 0, fmt.Errorf("NOMAD_TOKEN is not set")
	}

	externalAPIURL := fmt.Sprintf("%s/v1/jobs/statuses?namespace=*", baseURL)

	request, err := http.NewRequest("GET", externalAPIURL, nil)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	request.Header.Set("X-Nomad-Token", nomadTOKEN)
	response, err := http.DefaultClient.Do(request)

	// response, err := http.Get(externalAPIURL)

	if err != nil {
		return 0, 0, fmt.Errorf("failed to fetch data from external API: %v", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		return 0, 0, fmt.Errorf("failed to fetch data from external API: %v, %s", response.Status, string(bodyBytes))
	}

	var jobStatuses []JobStatus
	// bodyBytes, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	return 0, 0, fmt.Errorf("failed to read response body: %v", err)
	// }
	// if err := json.Unmarshal(bodyBytes, &jobStatuses); err != nil {
	// 	return 0, 0, fmt.Errorf("failed to decode response: %v", err)
	// }

	if err := json.NewDecoder(response.Body).Decode(&jobStatuses); err != nil {
		return 0, 0, fmt.Errorf("failed to decode response: %v", err)
	}

	for _, job := range jobStatuses {
		if !job.Stop && job.Status == "pending" {
			unhealthyCount++
		} else {
			healthyCount++
		}
	}

	return healthyCount, unhealthyCount, nil
}

// Structs for API responses
type NodeResponse struct {
	ID string `json:"ID"`
}

type StatsResponse struct {
	AllocDirStats struct {
		UsedPercent float64 `json:"UsedPercent"`
	} `json:"AllocDirStats"`
	Memory struct {
		Used  int64 `json:"Used"`
		Total int64 `json:"Total"`
	} `json:"Memory"`
	CPU []struct {
		TotalPercent float64 `json:"TotalPercent"`
	} `json:"CPU"`
}

func GetSystemUsage() (storageUsage float64, memoryUsage float64, averageCPUUsage float64, err error) {
	baseURL := os.Getenv("NOMAD_BASE_URL")
	if baseURL == "" {
		fmt.Println("NOMAD_BASE_URL is not set")
		return 0, 0, 0, fmt.Errorf("NOMAD_BASE_URL is not set")
	}

	nomadTOKEN := os.Getenv("NOMAD_TOKEN")
	if nomadTOKEN == "" {
		fmt.Println("NOMAD_TOKEN is not set")
		return 0, 0, 0, fmt.Errorf("NOMAD_TOKEN is not set")
	}

	getNodesAPIURL := fmt.Sprintf("%s/v1/nodes", baseURL)

	request, err := http.NewRequest("GET", getNodesAPIURL, nil)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	request.Header.Set("X-Nomad-Token", nomadTOKEN)
	response, err := http.DefaultClient.Do(request)

	// response, err := http.Get(getNodesAPIURL)

	if err != nil {
		return 0, 0, 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		return 0, 0, 0, fmt.Errorf("failed to fetch data from external API: %v, %s", response.Status, string(bodyBytes))
	}

	var nodes []NodeResponse
	// id, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	return 0, 0, 0, err
	// }
	// if err := json.Unmarshal(id, &nodes); err != nil {
	// 	return 0, 0, 0, err
	// }

	if err := json.NewDecoder(response.Body).Decode(&nodes); err != nil {
		return 0, 0, 0, err
	}

	if len(nodes) == 0 {
		return 0, 0, 0, fmt.Errorf("no nodes found")
	}

	nodeID := nodes[0].ID

	getUsageDetatilsAPIURL := fmt.Sprintf("%s/v1/client/stats?node_id=%s", baseURL, nodeID)

	request, err = http.NewRequest("GET", getUsageDetatilsAPIURL, nil)
	if err != nil {
		return 0, 0, 0, fmt.Errorf("failed to create HTTP request: %v", err)
	}
	request.Header.Set("X-Nomad-Token", nomadTOKEN)
	response, err = http.DefaultClient.Do(request)

	// response, err = http.Get(getUsageDetatilsAPIURL)
	if err != nil {
		return 0, 0, 0, err
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(response.Body)
		return 0, 0, 0, fmt.Errorf("failed to fetch data from external API: %v, %s", response.Status, string(bodyBytes))
	}

	var stats StatsResponse
	// bodyBytes, err := io.ReadAll(response.Body)
	// if err != nil {
	// 	return 0, 0, 0, err
	// }
	// if err := json.Unmarshal(bodyBytes, &stats); err != nil {
	// 	return 0, 0, 0, err
	// }
	if err := json.NewDecoder(response.Body).Decode(&stats); err != nil {
		return 0, 0, 0, err
	}

	storageUsage = math.Round(stats.AllocDirStats.UsedPercent*100) / 100

	if stats.Memory.Total > 0 {
		memoryUsage = (float64(stats.Memory.Used) / float64(stats.Memory.Total)) * 100
		memoryUsage = math.Round(memoryUsage*100) / 100
	}

	var totalCPUUsage float64
	numCores := len(stats.CPU)

	if numCores == 0 {
		return 0, 0, 0, fmt.Errorf("no CPU data found")
	}

	for _, cpu := range stats.CPU {
		totalCPUUsage += cpu.TotalPercent
	}

	averageCPUUsage = totalCPUUsage / float64(numCores)
	averageCPUUsage = math.Round(averageCPUUsage*100) / 100

	return storageUsage, memoryUsage, averageCPUUsage, nil
}
