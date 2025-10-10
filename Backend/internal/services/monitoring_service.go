package services

import (
	"encoding/json"
	"fmt"
	"gbackup-system/backend/internal/models"
	"gbackup-system/backend/internal/repository"
	"time"
)

type MonitoringService interface {
	UpdateRemoteStatus(remoteName string) error
	GetRemoteStatusList() ([]*models.Monitoring, error)
}

type MonitoringServiceImpl struct {
	MonitorRepo repository.MonitoringRepository
}

func NewMonitoringService(mRepo repository.MonitoringRepository) MonitoringService {
	return &MonitoringServiceImpl{MonitorRepo: mRepo}
}

func (s *MonitoringServiceImpl) UpdateRemoteStatus(remoteName string) error {
	rcloneArgs := []string{"about", remoteName + ":", "--json"}

	result := ExecuteRcloneJob(rcloneArgs)

	monitor := &models.Monitoring{
		RemoteName:    remoteName,
		LastCheckedAt: time.Now(),
	}

	if !result.Success {
		monitor.StatusConnect = "Disconnected"
		s.MonitorRepo.UpsertRemoteStatus(monitor)
		return fmt.Errorf("Gagal Terhubung ke %s, %s", remoteName, result.ErrorMsg)
	}

	var rcloneData struct {
		Total uint64 `json:"total"`
		Used  uint64 `json:"used"`
	}

	if err := json.Unmarshal([]byte(result.Output), &rcloneData); err != nil {
		return fmt.Errorf("Gagal Parsing Json Output Rclone: %v", err)
	}

	const BytesToGb = 1073741824.0 // konversi data GB ke bit
	totalGB := float64(rcloneData.Total) / BytesToGb
	usedGB := float64(rcloneData.Used) / BytesToGb

	monitor.StatusConnect = "Connected"
	monitor.TotalStorageGB = totalGB
	monitor.UsedStorageGB = usedGB
	monitor.FreeStorageGB = totalGB - usedGB

	return s.MonitorRepo.UpsertRemoteStatus(monitor)
}

func (s *MonitoringServiceImpl) GetRemoteStatusList() ([]*models.Monitoring, error) {
	return s.MonitorRepo.FindAllRemote()
}
