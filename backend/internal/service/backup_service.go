package service

import (
	"fmt"
	"gbackup-new/backend/internal/models"
	"gbackup-new/backend/internal/repository"
	"time"
)

// BackupService interface (Kontrak)
type BackupService interface {
	CreateJobAndDispatch(job *models.ScheduledJob) error
	TriggerManualJob(jobID uint) error
}

type backupServiceImpl struct {
	JobRepo repository.JobRepository
	LogRepo repository.LogRepository
}

func NewBackupService(jRepo repository.JobRepository, lRepo repository.LogRepository) BackupService {
	return &backupServiceImpl{JobRepo: jRepo, LogRepo: lRepo}
}

// ----------------------------------------------------
// FUNGSI UTAMA ORKESTRASI & DISPATCH
// ----------------------------------------------------

// CreateJobAndDispatch: Menerima Job dari Controller, menyimpan, dan/atau memicu
func (s *backupServiceImpl) CreateJobAndDispatch(job *models.ScheduledJob) error {

	// (Tambahkan logic enkripsi DbPass di sini jika diperlukan sebelum Create)

	// 1. SELALU SIMPAN JOB KE DATABASE (sebagai Template)
	if err := s.JobRepo.Create(job); err != nil {
		return fmt.Errorf("gagal menyimpan job template: %w", err)
	}

	// 2. JALANKAN JIKA MANUAL
	if job.ScheduleCron == "" {
		fmt.Printf("[DISPATCHER] Job %s (Manual) disimpan (ID: %d) dan dipicu langsung.\n", job.JobName, job.ID)
		// Memulai eksekusi di background (Goroutine)
		go s.executeJobLifecycle(*job)
		// Job Terjadwal (Auto Backup)
		fmt.Printf("[DISPATCHER] Job %d (%s) disimpan untuk Scheduler.\n", job.ID, job.JobName)
	}

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
	// ...
	// Contoh sederhana:
	args := []string{
		"rclone",
		job.RcloneMode,
		job.SourcePath,
		fmt.Sprintf("%s:%s", job.RemoteName, job.DestinationPath),
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
	s.JobRepo.UpdateLastRunStatus(job.ID, time.Now(), dbStatus)
}
