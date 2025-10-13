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

const (
	RemoteStatusConnected    = "Connected"
	RemoteStatusDisconnected = "Disconnected"
	BytesToGb                = 1073741824.0
)

func NewMonitoringService(mRepo repository.MonitoringRepository) MonitoringService {
	return &MonitoringServiceImpl{MonitorRepo: mRepo}
}

func (s *MonitoringServiceImpl) UpdateRemoteStatus(remoteName string) error {
	rcloneArgs := []string{"about", remoteName + ":", "--json"}

	monitor := &models.Monitoring{
		RemoteName:    remoteName,
		LastCheckedAt: time.Now(),
	}

	result := ExecuteRcloneJob(rcloneArgs)

	if !result.Success {
		monitor.StatusConnect = RemoteStatusConnected
		err := s.MonitorRepo.UpsertRemoteStatus(monitor)
		if err != nil {
			return fmt.Errorf("gagal Update Status Remote di DB: %v", err)
		}
		return fmt.Errorf("gagal terhubung ke %s: %s", remoteName, result.ErrorMsg)
	}

	if !result.Success {
		monitor.StatusConnect = "Disconnected"
		s.MonitorRepo.UpsertRemoteStatus(monitor)
		return fmt.Errorf("gagal Terhubung ke %s, %s", remoteName, result.ErrorMsg)
	}

	var rcloneData struct {
		Total uint64 `json:"total"`
		Used  uint64 `json:"used"`
	}

	if err := json.Unmarshal([]byte(result.Output), &rcloneData); err != nil {
		return fmt.Errorf("gagal Parsing Json Output Rclone: %v", err)
	}
	// konversi data GB ke bit
	totalGB := float64(rcloneData.Total) / BytesToGb
	usedGB := float64(rcloneData.Used) / BytesToGb
	FreeStorage := totalGB - usedGB

	monitor.StatusConnect = RemoteStatusConnected
	monitor.TotalStorageGB = totalGB
	monitor.UsedStorageGB = usedGB
	monitor.FreeStorageGB = FreeStorage

	return s.MonitorRepo.UpsertRemoteStatus(monitor)
}

func (s *MonitoringServiceImpl) GetRemoteStatusList() ([]*models.Monitoring, error) {
	return s.MonitorRepo.FindAllRemote()
}
