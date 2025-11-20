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
	GetJobLogs() ([]models.Log, error)
	DiscoverAndSaveRemote() error // Untuk UI Logs\
	startMonitoringDaemon()
	RunRemoteChecks() error
}

// monitoringServiceImpl adalah struct implementasi
type monitoringServiceImpl struct {
	MonitorRepo repository.MonitoringRepository
	LogRepo     repository.LogRepository // Dependency untuk GetJobLogs
	JobRepo     repository.JobRepository
}

const intervalCek = 30 * time.Minute

// NewMonitoringService adalah constructor untuk DI
func NewMonitoringService(mRepo repository.MonitoringRepository, lRepo repository.LogRepository, jRepo repository.JobRepository) MonitoringService {
	return &monitoringServiceImpl{
		MonitorRepo: mRepo,
		LogRepo:     lRepo,
		JobRepo:     jRepo,
	}
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

	count, errJob := s.JobRepo.CountJobOnRemote(remoteName)
	if errJob != nil {
		fmt.Printf("gagal menghitung job pada remote %s: %v\n", remoteName, errJob)
		monitor.ActiveJobCount = 0
	} else {
		monitor.ActiveJobCount = count
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

	usedPercentage := (monitor.UsedStorageGB / monitor.TotalStorageGB) * 100
	const WarningThreshold = 85.0

	// Early warn
	if usedPercentage >= WarningThreshold {
		monitor.SystemMessage = fmt.Sprintf("PERINGATAN: Storage terisi %.1f%%. Ruang kritis!", usedPercentage)
	} else if monitor.SystemMessage != "" {
		monitor.SystemMessage = ""

	}
	// 5. Update Database (melalui Repository Upsert)
	return s.MonitorRepo.UpsertRemoteStatus(monitor)
}

// GetRemoteStatusList: Mengambil data status dari DB untuk UI
func (s *monitoringServiceImpl) GetRemoteStatusList() ([]models.Monitoring, error) {
	remotes, err := s.MonitorRepo.FindAllRemotes()
	if err != nil {
		return nil, err
	}

	for i := range remotes {
		count, _ := s.JobRepo.CountJobOnRemote(remotes[i].RemoteName)
		remotes[i].ActiveJobCount = count
	}

	return remotes, nil
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

func (s *monitoringServiceImpl) DiscoverAndSaveRemote() error {
	remoteNames, err := s.GetRcloneConfiguredRemotes()
	if err != nil {
		return fmt.Errorf("gagal mengambil remote")
	}

	for _, name := range remoteNames {
		monitor := &models.Monitoring{
			RemoteName:    name,
			StatusConnect: "DISCONNECTED",
			LastCheckedAt: time.Now(),
		}

		if err := s.MonitorRepo.UpsertRemoteStatus(monitor); err != nil {
			return fmt.Errorf("gagal menyimpan remote %s ke DB: %w", name, err)
		}
	}

	for _, name := range remoteNames {
		// Jalankan di Goroutine agar tidak memblokir startup
		go s.UpdateRemoteStatus(name)
	}
	fmt.Printf("[MONITORING] Berhasil menemukan dan menyimpan %d remote.\n", len(remoteNames))
	return nil
}

func (s *monitoringServiceImpl) FindByRemoteName(remoteName string) (*models.Monitoring, error) {
	return s.MonitorRepo.FindRemoteByName(remoteName)
}

func (r *monitoringServiceImpl) RunRemoteChecks() error {
	remotes, err := r.MonitorRepo.FindAllRemotes()
	if err != nil {
		return err
	}

	for _, remote := range remotes {
		// Jalankan UpdateRemoteStatus di goroutine agar pengecekan Rclone berjalan paralel
		go r.UpdateRemoteStatus(remote.RemoteName)
	}
	return nil
}

func (r *monitoringServiceImpl) startMonitoringDaemon() {
	go func() {
		fmt.Printf("Monitoring Daemon aktif, pengecekan tiap %s\n", intervalCek)
		for {
			if err := r.RunRemoteChecks(); err != nil {
				fmt.Printf("⚠️ Daemon Error: %v\n", err)
			}
			time.Sleep(intervalCek)
		}
	}()
}
