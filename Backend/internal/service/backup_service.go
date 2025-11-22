package service

import (
	"encoding/json"
	"fmt"
	"gbackup-new/backend/internal/models"
	"gbackup-new/backend/internal/repository"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
	"strings"
	"time"
)

// BackupService interface (Kontrak)
type BackupService interface {
	CreateJobAndDispatch(job *models.ScheduledJob) error
	TriggerManualJob(jobID uint) error
	DeleteJob(JobId uint) error
	UpdateJob(jobID uint, updatedJob *models.ScheduledJob) error
	GetJobByID(jobID uint) (*models.ScheduledJob, error)
}

type backupServiceImpl struct {
	MonitorRepo repository.MonitoringRepository
	JobRepo     repository.JobRepository
	LogRepo     repository.LogRepository
	MonitorSvc  MonitoringService
}

type RcloneFileInfo struct {
	Name    string    `json:"Name"`    // Nama file/folder
	Size    int64     `json:"Size"`    // Ukuran dalam bytes
	ModTime time.Time `json:"ModTime"` // Waktu modifikasi
	IsDir   bool      `json:"IsDir"`   // True jika folder
}

func NewBackupService(
	jRepo repository.JobRepository,
	lRepo repository.LogRepository,
	mRepo repository.MonitoringRepository,
	mSvc MonitoringService,
) BackupService {
	return &backupServiceImpl{
		JobRepo:     jRepo,
		LogRepo:     lRepo,
		MonitorRepo: mRepo,
		MonitorSvc:  mSvc,
	}
}

const MinFreeGB = 1.0

// ----------------------------------------------------
// FUNGSI UTAMA ORKESTRASI & DISPATCH
// ----------------------------------------------------

// CreateJobAndDispatch: Menerima Job dari Controller, menyimpan, dan/atau memicu
func (s *backupServiceImpl) CreateJobAndDispatch(job *models.ScheduledJob) error {

	if job.OperationMode == "" {
		// Karena default DB Anda adalah BACKUP, kita set di Go.
		// Jika niatnya Restore, handler HARUS sudah menyetelnya.
		job.OperationMode = "BACKUP"
	}
	// [END PERBAIKAN KRITIS]

	// (Logika RcloneMode sudah diatasi di handler/service sebelumnya)
	if job.RcloneMode == "" {
		job.RcloneMode = "copy"
	}

	// 1. SELALU SIMPAN JOB KE DATABASE (sebagai Template)
	if err := s.JobRepo.Create(job); err != nil {
		return fmt.Errorf("gagal menyimpan job template: %w", err)
	}

	// 2. JALANKAN JIKA MANUAL
	if job.ScheduleCron == "" {
		fmt.Printf("[DISPATCHER] Job %s (Manual) disimpan (ID: %d) dan dipicu langsung.\n", job.JobName, job.ID)
		// Memulai eksekusi di background (Goroutine)
		go s.executeJobLifecycle(*job)
		return nil
		// Job Terjadwal (Auto Backup)
		// fmt.Printf("[DISPATCHER] Job %d (%s) disimpan untuk Scheduler.\n", job.ID, job.JobName)
	}
	fmt.Printf("[DISPATCHER] Job %d (%s) disimpan untuk Scheduler.\n", job.ID, job.JobName)
	return nil
}

// TriggerManualJob: Memicu Job yang sudah ada di DB
func (s *backupServiceImpl) TriggerManualJob(jobID uint) error {
	job, err := s.JobRepo.FindJobByID(jobID)
	if err != nil {
		return err
	}

	// Langsung eksekusi di background
	go s.executeJobLifecycle(*job)
	return nil
}

// ----------------------------------------------------
// FUNGSI EKSEKUSI 3 FASE (INTI)
// ----------------------------------------------------

// executeJobLifecycle: Menjalankan Pre-Script, Rclone, dan Post-Script
func (s *backupServiceImpl) executeJobLifecycle(job models.ScheduledJob) {
	fmt.Printf("[WORKER %d] Job %s: Memulai Eksekusi 3 Fase...\n", job.ID, job.JobName)

	// Set Status RUNNING (Locking)
	s.JobRepo.UpdateLastRunStatus(job.ID, time.Now(), "RUNNING")

	var finalResult RcloneResult
	var finalStatus string

	// 🚨 LANGKAH PENCEGAHAN (Hanya untuk Job Terjadwal) 🚨
	if job.OperationMode == "BACKUP" && job.ScheduleCron != "" {
		const MinFreeGB = 1.0 // Batas aman yang dibutuhkan

		// 1. Dapatkan status remote
		monitor, err := s.MonitorRepo.FindRemoteByName(job.RemoteName)
		if err == nil && monitor != nil {

			// 2. Hitung/Estimasi ukuran sumber
			sourceSizeGB, _ := s.CalculateSourceSizeGB(job.SourcePath)
			requiredSpace := sourceSizeGB + MinFreeGB

			// 3. Bandingkan dan SUSPEND
			if monitor.FreeStorageGB < requiredSpace {
				errorMsg := fmt.Sprintf("⛔ Job ditangguhkan: Ruang di %s (%.2f GB) tidak cukup untuk data ini (%.2f GB).",
					job.RemoteName, monitor.FreeStorageGB, requiredSpace)

				// Catat kegagalan dan SUSPEND Job
				s.handleJobCompletion(job, RcloneResult{Success: false, ErrorMsg: errorMsg}, "FAIL_STORAGE")
				s.JobRepo.UpdateJobActivity(job.ID, false) // ✅ SUSPEND JOB
				fmt.Printf("⛔ [WORKER %d] Job Terjadwal DITANGGUHKAN karena ruang penuh.\n", job.ID)
				return
			}
		}
		// Jika gagal mendapatkan status monitor, biarkan Job berjalan (risiko kecil).
	}

	// --- FASE 1: PRE-SCRIPT ---
	if job.PreScript != "" {
		fmt.Printf("[WORKER %d] Menjalankan Pre-Script...\n", job.ID)
		// "Hardening" script user
		hardenedPreScript := fmt.Sprintf("set -eo pipefail; \n%s", job.PreScript)
		preScriptArgs := []string{"bash", "-c", hardenedPreScript}

		result := ExecuteCliJob(preScriptArgs)
		if !result.Success {
			fmt.Printf("❌ [WORKER %d] Pre-Script GAGAL.\n", job.ID)
			finalResult = result
			finalStatus = "FAIL_PRE_SCRIPT"
			s.handleJobCompletion(job, finalResult, finalStatus)
			return // Hentikan eksekusi
		}
	}

	// Timestamp
	// ============================================================
	// 🆕 FASE 1.5: TIMESTAMP & ROUND ROBIN
	// ============================================================
	var runtimeDestPath string

	if job.OperationMode == "BACKUP" {
		timestamp := time.Now().Format("20060102_150405")

		sourceInfo, err := os.Stat(job.SourcePath)
		if err != nil {
			errorMsg := fmt.Sprintf("Failed to stat source path: %v", err)
			s.handleJobCompletion(job, RcloneResult{Success: false, ErrorMsg: errorMsg}, "FAIL_SOURCE_CHECK")
			return
		}

		var newDestinationName string

		var isSourceDir = sourceInfo.IsDir()

		if isSourceDir {
			// Folder → hasil tetap folder
			folderName := filepath.Base(job.SourcePath)
			newDestinationName = fmt.Sprintf("%s_%s", folderName, timestamp)
		} else {
			// FILE → hasil tetap file
			fileName := filepath.Base(job.SourcePath)
			ext := filepath.Ext(fileName)
			nameWithoutExt := strings.TrimSuffix(fileName, ext)
			newDestinationName = fmt.Sprintf("%s_%s%s", nameWithoutExt, timestamp, ext)
		}

		originalDestPath := job.DestinationPath
		runtimeDestPath = filepath.Join(originalDestPath, newDestinationName)

		fmt.Printf("[WORKER %d] 🎯 Runtime destination: %s:%s\n", job.ID, job.RemoteName, runtimeDestPath)
		fmt.Printf("[WORKER %d] 💾 DB destination (unchanged): %s\n", job.ID, job.DestinationPath)

		// Round Robin Cleanup
		fmt.Printf("[WORKER %d] 🔄 Checking for old backups...\n", job.ID)
		if err := s.CleanupOldBackups(job.RemoteName, originalDestPath); err != nil {
			fmt.Printf("⚠️ [WORKER %d] Cleanup warning: %v\n", job.ID, err)
		}
	} else {
		runtimeDestPath = job.DestinationPath
	}
	// ============================================================

	// --- FASE 2: RCLONE EXECUTION ---
	fmt.Printf("[WORKER %d] Menjalankan Rclone...\n", job.ID)
	rcloneArgs := s.buildRcloneArgs(job, runtimeDestPath)
	resultRclone := ExecuteCliJob(rcloneArgs)

	if !resultRclone.Success {
		fmt.Printf("❌ [WORKER %d] Rclone GAGAL.\n", job.ID)
		finalResult = resultRclone
		finalStatus = "FAIL_RCLONE"
		s.handleJobCompletion(job, finalResult, finalStatus)
		return // Hentikan eksekusi
	}

	lines := strings.Split(resultRclone.Output, "\n")
	var transferStatus string
	for _, line := range lines {
		if strings.Contains(line, "Transferred:") {
			transferStatus = strings.TrimSpace(line)
		}
	}

	if transferStatus != "" {
		fmt.Printf("📊 [WORKER %d] Stats: %s\n", job.ID, transferStatus)
		// Tambahkan info ini ke output agar tersimpan di DB Log
		resultRclone.Output = fmt.Sprintf("%s\n\n%s", transferStatus, resultRclone.Output)
	}

	// --- FASE 3: POST-SCRIPT ---
	if job.PostScript != "" {
		fmt.Printf("[WORKER %d] Menjalankan Post-Script...\n", job.ID)
		hardenedPostScript := fmt.Sprintf("set -eo pipefail; \n%s", job.PostScript)
		postScriptArgs := []string{"bash", "-c", hardenedPostScript}

		resultPost := ExecuteCliJob(postScriptArgs)
		if !resultPost.Success {
			fmt.Printf("❌ [WORKER %d] Post-Script GAGAL.\n", job.ID)
			finalResult = resultPost
			finalStatus = "FAIL_POST_SCRIPT"
			s.handleJobCompletion(job, finalResult, finalStatus)
			return // Hentikan eksekusi
		}
	}

	// --- FASE 4: SUKSES ---
	fmt.Printf("✅ [WORKER %d] Job Selesai.\n", job.ID)
	finalResult = resultRclone // Log output Rclone jika sukses
	finalStatus = "SUCCESS"
	s.handleJobCompletion(job, finalResult, finalStatus)
}

// ----------------------------------------------------
// FUNGSI HELPER (COMMAND GENERATION & LOGGING)
// ----------------------------------------------------

// buildRcloneArgs: Menyusun command Rclone
func (s *backupServiceImpl) buildRcloneArgs(job models.ScheduledJob, runtimeDestPath string) []string {
	isRestore := job.OperationMode == "RESTORE"
	command := strings.ToLower(job.RcloneMode)

	// Tentukan apakah sumber adalah file atau folder
	sourceInfo, err := os.Stat(job.SourcePath)
	isSourceDir := false
	if err == nil {
		isSourceDir = sourceInfo.IsDir()
	}

	var SourcePath, Destination string

	if isRestore {
		// RESTORE SELALU copy folder atau file dari remote ke lokal
		SourcePath = fmt.Sprintf("%s:%s", job.RemoteName, job.SourcePath)
		Destination = job.DestinationPath
		command = "copy" // restore selalu copy
	} else {
		// BACKUP
		SourcePath = job.SourcePath
		Destination = fmt.Sprintf("%s:%s", job.RemoteName, runtimeDestPath)

		if !isSourceDir {
			// File → gunakan copyto agar tidak dianggap folder
			command = "copyto"
		} else {
			// Folder → tetap copy
			command = "copy"
		}
	}

	args := []string{
		"rclone",
		command,
		SourcePath,
		Destination,
		"--checksum",
		"--no-traverse",
		"--progress",    // ✅ TAMBAHAN: Tampilkan progress dengan stats
		"--stats", "5s", // ✅ UPDATE: Print stats setiap 5 detik (bukan 0)
		"--stats-log-level", "INFO", // ✅ UPDATE: Change to INFO level
		"--human-readable",
	}

	return args
}

// handleJobCompletion: Logika Logging dan Final Status Update
func (s *backupServiceImpl) handleJobCompletion(job models.ScheduledJob, result RcloneResult, status string) {
	LogMutex.Lock()
	defer LogMutex.Unlock()
	stats := parseRcloneStats(result.Output)
	logMessage := result.ErrorMsg

	// 1. Catat ke tabel Logs (dengan TransferredBytes yang sudah di-parse)
	newLog := &models.Log{
		JobID:            &job.ID,
		Status:           status,
		Message:          logMessage,
		DurationSec:      int(result.Duration.Seconds()),
		TransferredBytes: result.TransferredBytes, // ✅ UPDATED: dari parsed bytes
		Timestamp:        time.Now(),
	}

	fmt.Printf("[LOG DEBUG] Saving ID: %d | Status: %s | Bytes: %d\n", job.ID, status, result.TransferredBytes)

	s.LogRepo.CreateLog(newLog)

	// 2. Update status akhir di tabel scheduled_jobs
	dbStatus := status
	if status == "SUCCESS" {
		dbStatus = "COMPLETED"
	}

	if job.ScheduleCron == "" {
		dbStatus = "PENDING"
	}

	s.JobRepo.UpdateLastRunStatus(job.ID, time.Now(), dbStatus)

	// ✅ TAMBAHAN: Log ringkasan
	if status == "SUCCESS" {
		fmt.Printf("✅ [COMPLETE] Job %d: Transferred %.2f GB in %d seconds (Speed: %s)\n",
			job.ID,
			result.Duration,
			stats.Speed,
		)
	} else {
		fmt.Printf("❌ [FAILED] Job %d: %s\n", job.ID, status)
	}
}

func (s *backupServiceImpl) CalculateSourceSizeGB(path string) (float64, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("gagal stat path %s: %w", path, err)
	}
	const BytesToGB = 1073741824.0
	return float64(info.Size()) / BytesToGB, nil
}

func (s *backupServiceImpl) DeleteJob(JobID uint) error {
	fmt.Printf("[AUDIT] User meminta penghapusan Job ID: %d\n", JobID)

	// Panggil Repository untuk menghapus
	if err := s.JobRepo.DeleteJob(JobID); err != nil {
		return fmt.Errorf("gagal menghapus job ID %d: %w", JobID, err)
	}
	return nil
}

func (s *backupServiceImpl) GetJobByID(jobID uint) (*models.ScheduledJob, error) {
	job, err := s.JobRepo.FindJobByID(jobID)
	if err != nil {
		return nil, err
	}
	return job, nil
}
func (s *backupServiceImpl) UpdateJob(jobID uint, updatedJob *models.ScheduledJob) error {
	fmt.Printf("[UPDATE] Memperbarui Job ID: %d\n", jobID)

	// 1. Cek apakah job exist
	existingJob, err := s.JobRepo.FindJobByID(jobID)
	if err != nil {
		return fmt.Errorf("job tidak ditemukan: %w", err)
	}

	// 2. ✅ Build update map (hanya field yang ada)
	updates := make(map[string]interface{})

	if updatedJob.JobName != "" {
		updates["job_name"] = updatedJob.JobName
	}
	if updatedJob.OperationMode != "" {
		updates["operation_mode"] = updatedJob.OperationMode
	}
	if updatedJob.RcloneMode != "" {
		updates["rclone_mode"] = updatedJob.RcloneMode
	}
	if updatedJob.SourcePath != "" {
		updates["source_path"] = updatedJob.SourcePath
	}
	if updatedJob.DestinationPath != "" {
		updates["destination_path"] = updatedJob.DestinationPath
	}
	if updatedJob.RemoteName != "" {
		updates["remote_name"] = updatedJob.RemoteName
	}

	// ✅ Allow empty string untuk script (untuk clear script)
	updates["pre_script"] = updatedJob.PreScript
	updates["post_script"] = updatedJob.PostScript

	// ✅ Schedule cron bisa kosong (untuk ubah jadi manual job)
	updates["schedule_cron"] = updatedJob.ScheduleCron

	// ✅ IsActive harus explicit check
	if updatedJob.IsActive != existingJob.IsActive {
		updates["is_active"] = updatedJob.IsActive
	}

	updates["updated_at"] = time.Now()

	// 3. ✅ Validate minimal 1 field (exclude updated_at)
	if len(updates) <= 1 {
		return fmt.Errorf("tidak ada field yang diubah")
	}

	// 4. Update ke database
	if err := s.JobRepo.UpdateJob(jobID, updates); err != nil {
		return fmt.Errorf("gagal update job: %w", err)
	}

	fmt.Printf("[UPDATE] Job %d berhasil diperbarui (%d fields)\n", jobID, len(updates)-1)
	return nil
}

func (s *backupServiceImpl) CleanupOldBackups(remoteName, destinationPath string) error {
	const MAX_RETENTION = 10

	fmt.Printf("[Round Robin] Checking backups in %s:%s...\n", remoteName, destinationPath)

	listCmd := fmt.Sprintf("rclone lsjson %s:%s", remoteName, destinationPath)
	output, err := exec.Command("bash", "-c", listCmd).Output()
	if err != nil {
		return fmt.Errorf("failed to list remote files: %w", err)
	}

	var files []RcloneFileInfo
	if err := json.Unmarshal(output, &files); err != nil {
		return fmt.Errorf("failed to parse rclone output: %w", err)
	}

	var backupItems []RcloneFileInfo
	timestampPattern := regexp.MustCompile(`_\d{8}_\d{6}`)

	for _, f := range files {
		if timestampPattern.MatchString(f.Name) {
			backupItems = append(backupItems, f)
		}
	}

	currentCount := len(backupItems)
	fmt.Printf("[Round Robin] Found %d backup items (limit: %d)\n", currentCount, MAX_RETENTION)

	if currentCount < MAX_RETENTION {
		fmt.Printf("[Round Robin] No cleanup needed (%d/%d)\n", currentCount, MAX_RETENTION)
		return nil
	}

	sort.Slice(backupItems, func(i, j int) bool {
		return backupItems[i].ModTime.Before(backupItems[j].ModTime)
	})

	itemsToDelete := currentCount - MAX_RETENTION + 1

	fmt.Printf("[Round Robin] Deleting %d oldest backup(s)...\n", itemsToDelete)

	var deletedItems []string
	for i := 0; i < itemsToDelete && i < len(backupItems); i++ {
		itemToDelete := backupItems[i]
		fullPath := fmt.Sprintf("%s:%s/%s", remoteName, destinationPath, itemToDelete.Name)

		var deleteCmd string
		if itemToDelete.IsDir {
			deleteCmd = fmt.Sprintf("rclone purge %s", fullPath)
		} else {
			deleteCmd = fmt.Sprintf("rclone delete %s", fullPath)
		}

		if err := exec.Command("bash", "-c", deleteCmd).Run(); err != nil {
			fmt.Printf("⚠️  [Round Robin] Failed to delete %s: %v\n", itemToDelete.Name, err)
			continue
		}

		deletedItems = append(deletedItems, itemToDelete.Name)

		itemType := "file"
		sizeStr := fmt.Sprintf("%.2f MB", float64(itemToDelete.Size)/(1024*1024))
		if itemToDelete.IsDir {
			itemType = "folder"
			sizeStr = "folder"
		}

		fmt.Printf("🗑️  [Round Robin] Deleted %s: %s (date: %s, size: %s)\n",
			itemType,
			itemToDelete.Name,
			itemToDelete.ModTime.Format("2006-01-02 15:04:05"),
			sizeStr,
		)
	}

	fmt.Printf("✅ [Round Robin] Cleanup complete. Deleted %d item(s). Space available for new backup.\n", len(deletedItems))
	return nil
}
