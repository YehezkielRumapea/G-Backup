package service

import (
	"encoding/json"
	"fmt"
	"gbackup-new/backend/internal/models"
	"gbackup-new/backend/internal/repository"
	"strings"
	"time"
)

type MonitoringService interface {
	UpdateRemoteStatus(remoteName string) error
	GetRemoteStatusList() ([]models.Monitoring, error)
	GetRcloneConfiguredRemotes() ([]string, error)
	GetJobLogs() ([]models.Log, error)
	DiscoverAndSaveRemote() error
	RunRemoteChecks() error
	StartMonitoringDaemon()
	SyncRemotesWithRclone() error
}

type monitoringServiceImpl struct {
	MonitorRepo repository.MonitoringRepository
	LogRepo     repository.LogRepository
	JobRepo     repository.JobRepository
}

const intervalCek = 5 * time.Minute

func NewMonitoringService(mRepo repository.MonitoringRepository, lRepo repository.LogRepository, jRepo repository.JobRepository) MonitoringService {
	return &monitoringServiceImpl{
		MonitorRepo: mRepo,
		LogRepo:     lRepo,
		JobRepo:     jRepo,
	}
}

func (s *monitoringServiceImpl) UpdateRemoteStatus(remoteName string) error {
	fmt.Printf("[Monitoring] Mengecek remote: %s\n", remoteName)

	rcloneArgs := []string{"rclone", "about", remoteName + ":", "--json"}
	result := ExecuteCliJob(rcloneArgs)

	monitor := &models.Monitoring{
		RemoteName:    remoteName,
		LastCheckedAt: time.Now(),
	}

	count, errJob := s.JobRepo.CountJobOnRemote(remoteName)
	if errJob != nil {
		fmt.Printf("⚠️ Gagal hitung job pada %s: %v\n", remoteName, errJob)
		monitor.ActiveJobCount = 0
	} else {
		monitor.ActiveJobCount = count
	}

	// Di dalam UpdateRemoteStatus (blok if !result.Success)

	if !result.Success {
		fmt.Printf("❌ Rclone Error pada %s: %s\n", remoteName, result.ErrorMsg)
		monitor.StatusConnect = "DISCONNECTED"

		// PERBAIKAN: Hapus '&' agar tipenya string (value), bukan *string (pointer)
		monitor.SystemMessage = result.ErrorMsg

		s.MonitorRepo.UpsertRemoteStatus(monitor)
		return fmt.Errorf("gagal terhubung ke %s: %s", remoteName, result.ErrorMsg)
	}

	var rcloneData struct {
		Total uint64 `json:"total"`
		Used  uint64 `json:"used"`
		Free  uint64 `json:"free"`
	}

	if err := json.Unmarshal([]byte(result.Output), &rcloneData); err != nil {
		fmt.Printf("❌ JSON Parse Error: %v\n", err)
		return fmt.Errorf("gagal parsing JSON: %v", err)
	}

	const BytesToGB = 1073741824.0

	monitor.StatusConnect = "CONNECTED"
	monitor.TotalStorageGB = float64(rcloneData.Total) / BytesToGB
	monitor.UsedStorageGB = float64(rcloneData.Used) / BytesToGB
	monitor.FreeStorageGB = float64(rcloneData.Free) / BytesToGB

	usedPercentage := (monitor.UsedStorageGB / monitor.TotalStorageGB) * 100
	const WarningThreshold = 85.0

	if usedPercentage >= WarningThreshold {
		msg := fmt.Sprintf("⚠️ Storage terisi %.1f%%. PERINGATAN!", usedPercentage)
		monitor.SystemMessage = msg
	} else {
		monitor.SystemMessage = ""
	}

	fmt.Printf("✅ %s: %.2f GB / %.2f GB (%.1f%%)\n",
		remoteName, monitor.UsedStorageGB, monitor.TotalStorageGB, usedPercentage)

	return s.MonitorRepo.UpsertRemoteStatus(monitor)
}

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

func (s *monitoringServiceImpl) GetRcloneConfiguredRemotes() ([]string, error) {
	result := ExecuteCliJob([]string{"rclone", "listremotes"})

	if !result.Success {
		return nil, fmt.Errorf("gagal mendapatkan daftar remote: %s", result.ErrorMsg)
	}

	remotes := strings.Split(result.Output, "\n")
	var cleanNames []string

	for _, remote := range remotes {
		name := strings.TrimSpace(remote)
		if len(name) > 0 && strings.HasSuffix(name, ":") {
			cleanNames = append(cleanNames, name[:len(name)-1])
		}
	}
	return cleanNames, nil
}

func (s *monitoringServiceImpl) GetJobLogs() ([]models.Log, error) {
	return s.LogRepo.FindAllLogs()
}

func (s *monitoringServiceImpl) SyncRemotesWithRclone() error {
	fmt.Println("\n" + strings.Repeat("=", 60))
	fmt.Println("[SYNC] Memulai sinkronisasi remote DB dengan rclone.conf")
	fmt.Println(strings.Repeat("=", 60))

	// 1. Ambil remote dari rclone.conf
	fmt.Println("[SYNC] 1. Mengambil daftar remote dari rclone.conf...")
	rcloneRemotes, err := s.GetRcloneConfiguredRemotes()
	if err != nil {
		fmt.Printf("[ERROR] Gagal mengambil remote dari rclone: %v\n", err)
		return fmt.Errorf("gagal mengambil remote dari rclone: %w", err)
	}
	fmt.Printf("[SYNC] ✓ Ditemukan %d remote di rclone.conf: %v\n", len(rcloneRemotes), rcloneRemotes)

	// 2. Ambil remote yang ada di database
	fmt.Println("[SYNC] 2. Mengambil daftar remote dari database...")
	dbRemotes, err := s.MonitorRepo.GetAllRemoteNames()
	if err != nil {
		fmt.Printf("[ERROR] Gagal mengambil remote dari DB: %v\n", err)
		return fmt.Errorf("gagal mengambil remote dari DB: %w", err)
	}
	fmt.Printf("[SYNC] ✓ Ditemukan %d remote di database: %v\n", len(dbRemotes), dbRemotes)

	// 3. Cari remote yang dihapus dari rclone tapi masih ada di DB
	fmt.Println("[SYNC] 3. Mencari remote yang dihapus dari rclone...")
	remoteToDelete := findMissingRemotes(dbRemotes, rcloneRemotes)

	if len(remoteToDelete) > 0 {
		fmt.Printf("[SYNC] ⚠️ Ditemukan %d remote yang sudah dihapus dari rclone: %v\n", len(remoteToDelete), remoteToDelete)

		for _, remoteName := range remoteToDelete {
			fmt.Printf("[SYNC] → Menghapus '%s' dari database...\n", remoteName)
			if err := s.MonitorRepo.DeleteRemoteByName(remoteName); err != nil {
				fmt.Printf("[ERROR] Gagal menghapus '%s': %v\n", remoteName, err)
			} else {
				fmt.Printf("[SYNC] ✓ Berhasil menghapus '%s'\n", remoteName)
			}
		}
	} else {
		fmt.Println("[SYNC] ✓ Tidak ada remote yang perlu dihapus")
	}

	// 4. Cari remote baru dari rclone yang belum di DB
	fmt.Println("[SYNC] 4. Mencari remote baru dari rclone...")
	remoteToAdd := findMissingRemotes(rcloneRemotes, dbRemotes)

	if len(remoteToAdd) > 0 {
		fmt.Printf("[SYNC] ℹ️ Ditemukan %d remote baru dari rclone: %v\n", len(remoteToAdd), remoteToAdd)

		for _, remoteName := range remoteToAdd {
			fmt.Printf("[SYNC] → Menambahkan '%s' ke database...\n", remoteName)
			monitor := &models.Monitoring{
				RemoteName:    remoteName,
				StatusConnect: "DISCONNECTED",
				LastCheckedAt: time.Now(),
			}
			if err := s.MonitorRepo.UpsertRemoteStatus(monitor); err != nil {
				fmt.Printf("[ERROR] Gagal menambah '%s': %v\n", remoteName, err)
			} else {
				fmt.Printf("[SYNC] ✓ Berhasil menambahkan '%s'\n", remoteName)
			}
		}
	} else {
		fmt.Println("[SYNC] ✓ Tidak ada remote baru")
	}

	if len(remoteToDelete) == 0 && len(remoteToAdd) == 0 {
		fmt.Println("[SYNC] ✅ Remote DB sudah sesuai dengan rclone.conf")
	}

	fmt.Println(strings.Repeat("=", 60) + "\n")
	return nil
}

func (s *monitoringServiceImpl) DiscoverAndSaveRemote() error {
	fmt.Println("[Discovery] Menyelaraskan remote dari rclone.conf...")

	if err := s.SyncRemotesWithRclone(); err != nil {
		return fmt.Errorf("gagal sync remote: %w", err)
	}

	remotes, err := s.MonitorRepo.FindAllRemotes()
	if err != nil {
		return fmt.Errorf("gagal mengambil remote dari DB: %w", err)
	}

	if len(remotes) == 0 {
		fmt.Println("[Discovery] ⚠️ Tidak ada remote yang dikonfigurasi")
		return nil
	}

	fmt.Printf("[Discovery] Memulai update status untuk %d remote...\n", len(remotes))

	for _, remote := range remotes {
		go s.UpdateRemoteStatus(remote.RemoteName)
	}

	return nil
}

func (s *monitoringServiceImpl) StartMonitoringDaemon() {
	go func() {
		fmt.Printf("[Daemon] Monitoring Daemon aktif, pengecekan tiap %v\n", intervalCek)

		if err := s.RunRemoteChecks(); err != nil {
			fmt.Printf("⚠️ Initial check error: %v\n", err)
		}

		if err := s.SyncRemotesWithRclone(); err != nil {
			fmt.Printf("❌ ERROR: Gagal sync remote: %v\n", err)
			return
		}

		ticker := time.NewTicker(intervalCek)
		defer ticker.Stop()

		for range ticker.C {
			fmt.Println("[Daemon] Menjalankan remote check...")
			if err := s.RunRemoteChecks(); err != nil {
				fmt.Printf("⚠️ Daemon Error: %v\n", err)
			}
		}
	}()
}

func (s *monitoringServiceImpl) RunRemoteChecks() error {
	remotes, err := s.MonitorRepo.FindAllRemotes()
	if err != nil {
		return err
	}

	for _, remote := range remotes {
		go s.UpdateRemoteStatus(remote.RemoteName)
	}
	return nil
}

func (s *monitoringServiceImpl) FindByRemoteName(remoteName string) (*models.Monitoring, error) {
	return s.MonitorRepo.FindRemoteByName(remoteName)
}

func findMissingRemotes(sourceList, targetList []string) []string {
	missing := []string{}

	targetMap := make(map[string]bool)
	for _, item := range targetList {
		targetMap[item] = true
	}

	for _, item := range sourceList {
		if !targetMap[item] {
			missing = append(missing, item)
		}
	}

	return missing
}
