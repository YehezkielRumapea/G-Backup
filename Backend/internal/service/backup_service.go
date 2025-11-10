package service

import (
	"fmt"
	"gbackup-new/backend/internal/models"
	"gbackup-new/backend/internal/repository"
	"os"
	"strings"
	"time"
)

// BackupService interface (Kontrak)
type BackupService interface {
	CreateJobAndDispatch(job *models.ScheduledJob) error
	TriggerManualJob(jobID uint) error
	DeleteJob(JobId uint) error
}

type backupServiceImpl struct {
	MonitorRepo repository.MonitoringRepository
	JobRepo     repository.JobRepository
	LogRepo     repository.LogRepository
	MonitorSvc  MonitoringService
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

	// --- FASE 2: RCLONE EXECUTION ---
	fmt.Printf("[WORKER %d] Menjalankan Rclone...\n", job.ID)
	rcloneArgs := s.buildRcloneArgs(job)
	resultRclone := ExecuteCliJob(rcloneArgs)

	if !resultRclone.Success {
		fmt.Printf("❌ [WORKER %d] Rclone GAGAL.\n", job.ID)
		finalResult = resultRclone
		finalStatus = "FAIL_RCLONE"
		s.handleJobCompletion(job, finalResult, finalStatus)
		return // Hentikan eksekusi
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
func (s *backupServiceImpl) buildRcloneArgs(job models.ScheduledJob) []string {
	// (Implementasi logic inversi path Restore/Backup di sini)
	isRestore := job.OperationMode == "RESTORE"
	var SourcePath, Destination string
	command := strings.ToLower(job.RcloneMode)

	if isRestore {
		SourcePath = fmt.Sprintf("%s:%s", job.RemoteName, job.SourcePath)
		Destination = job.DestinationPath
		command = "copy"

		fmt.Printf("[DEBUG INVERSION] Mode di Struct: %s | isRestore: %t | Cmd: %s | Source: %s\n",
			job.OperationMode,
			isRestore,
			command,
			SourcePath)
	} else {
		SourcePath = job.SourcePath
		Destination = fmt.Sprintf("%s:%s", job.RemoteName, job.DestinationPath)
	}
	args := []string{
		"rclone",
		command,
		SourcePath,
		Destination,
		"--checksum", // Flag keamanan
	}
	return args
}

// handleJobCompletion: Logika Logging dan Final Status Update
func (s *backupServiceImpl) handleJobCompletion(job models.ScheduledJob, result RcloneResult, status string) {
	LogMutex.Lock()
	defer LogMutex.Unlock()

	// 1. Catat ke tabel Logs
	newLog := &models.Log{
		JobID:       &job.ID,
		Status:      status,
		Message:     result.Output + result.ErrorMsg,
		DurationSec: int(result.Duration.Seconds()),
		Timestamp:   time.Now(),
	}
	s.LogRepo.CreateLog(newLog)

	// 2. Update status akhir di tabel scheduled_jobs
	dbStatus := status
	if status == "SUCCESS" {
		dbStatus = "COMPLETED"
	}

	if job.ScheduleCron == "" {
		// Status PENDING/IDLE cocok untuk menunggu trigger manual berikutnya
		dbStatus = "PENDING" // Atau "IDLE", tergantung konvensi Anda
	}

	s.JobRepo.UpdateLastRunStatus(job.ID, time.Now(), dbStatus)
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
