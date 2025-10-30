package service

import (
	"encoding/json"
	"fmt"
	"gbackup-new/backend/internal/models"
	"gbackup-new/backend/internal/repository"
	"strings" // Diperlukan untuk parsing output listremotes
	"time"
)

// MonitoringService interface (Kontrak)
type MonitoringService interface {
	UpdateRemoteStatus(remoteName string) error
	GetRemoteStatusList() ([]models.Monitoring, error)
	GetRcloneConfiguredRemotes() ([]string, error) // Untuk Startup Discovery
	GetJobLogs() ([]models.Log, error)             // Untuk UI Logs
}

// monitoringServiceImpl adalah struct implementasi
type monitoringServiceImpl struct {
	MonitorRepo repository.MonitoringRepository
	LogRepo     repository.LogRepository // Dependency untuk GetJobLogs
}

// NewMonitoringService adalah constructor untuk DI
func NewMonitoringService(mRepo repository.MonitoringRepository, lRepo repository.LogRepository) MonitoringService {
	return &monitoringServiceImpl{MonitorRepo: mRepo, LogRepo: lRepo}
}

// ----------------------------------------------------
// FUNGSI IMPLEMENTASI
// ----------------------------------------------------

// UpdateRemoteStatus: Menggunakan Executor untuk memperbarui status remote
func (s *monitoringServiceImpl) UpdateRemoteStatus(remoteName string) error {
	// 1. Generate command
	rcloneArgs := []string{"rclone", "about", remoteName + ":", "--json"}

	// 2. Eksekusi command Rclone (menggunakan Executor)
	result := ExecuteCliJob(rcloneArgs)

	monitor := &models.Monitoring{
		RemoteName:    remoteName,
		LastCheckedAt: time.Now(),
	}

	if !result.Success {
		// --- LOGIKA KEGAGALAN KONEKSI ---
		monitor.StatusConnect = "DISCONNECTED"
		s.MonitorRepo.UpsertRemoteStatus(monitor) // Tetap update DB (status DISCONNECTED)

		// Opsional: Catat error sistem ini ke tabel logs
		// s.LogRepo.CreateLog(...)

		return fmt.Errorf("gagal terhubung ke %s: %s", remoteName, result.ErrorMsg)
	}

	// --- LOGIKA SUKSES ---

	// Struct Golang untuk Parsing Output JSON Rclone
	var rcloneData struct {
		Total uint64 `json:"total"`
		Used  uint64 `json:"used"`
		Free  uint64 `json:"free"`
	}

	// 3. Parsing Output JSON
	if err := json.Unmarshal([]byte(result.Output), &rcloneData); err != nil {
		return fmt.Errorf("gagal parsing JSON output Rclone: %v", err)
	}

	// 4. Konversi Data (Bytes ke GB)
	const BytesToGB = 1073741824.0 // 1024^3

	monitor.StatusConnect = "CONNECTED"
	monitor.TotalStorageGB = float64(rcloneData.Total) / BytesToGB
	monitor.UsedStorageGB = float64(rcloneData.Used) / BytesToGB
	monitor.FreeStorageGB = float64(rcloneData.Free) / BytesToGB

	// 5. Update Database (melalui Repository Upsert)
	return s.MonitorRepo.UpsertRemoteStatus(monitor)
}

// GetRemoteStatusList: Mengambil data status dari DB untuk UI
func (s *monitoringServiceImpl) GetRemoteStatusList() ([]models.Monitoring, error) {
	return s.MonitorRepo.FindAllRemotes()
}

// GetRcloneConfiguredRemotes: Mengambil daftar remote dari rclone.conf (Untuk Startup Discovery)
func (s *monitoringServiceImpl) GetRcloneConfiguredRemotes() ([]string, error) {
	result := ExecuteCliJob([]string{"rclone", "listremotes"})

	if !result.Success {
		return nil, fmt.Errorf("gagal mendapatkan daftar remote dari rclone.conf: %s", result.ErrorMsg)
	}

	remotes := strings.Split(result.Output, "\n")
	var cleanNames []string

	for _, remote := range remotes {
		name := strings.TrimSpace(remote)
		if len(name) > 0 && strings.HasSuffix(name, ":") {
			cleanNames = append(cleanNames, name[:len(name)-1]) // Hapus titik dua
		}
	}
	return cleanNames, nil
}

// GetJobLogs: Mengambil riwayat log dari LogRepository
func (s *monitoringServiceImpl) GetJobLogs() ([]models.Log, error) {
	return s.LogRepo.FindAllLogs()
}
