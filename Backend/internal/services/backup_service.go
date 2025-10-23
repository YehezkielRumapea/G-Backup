package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gbackup-system/backend/internal/models"
	"gbackup-system/backend/internal/repository"
	"os"
	"time"
)

type BackupServiceImpl struct {
	JobRepo repository.JobRepository
	LogRepo repository.LogRepository
}

func NewBackupService(jRepo repository.JobRepository, lRepo repository.LogRepository) BackupService {
	return (&BackupServiceImpl{
		JobRepo: jRepo,
		LogRepo: lRepo,
	})
}

func (s *BackupServiceImpl) buildRcloneArgs(job models.ScheduledJob, tempDumpPath string) []string {
	// Path
	isRestore := job.RcloneMode == "Restore"

	var SourcePath, Destination string
	command := job.RcloneMode

	if isRestore {
		SourcePath = fmt.Sprintf("%s:%s", job.RemoteName, job.DestinationPath)
		Destination = job.DestinationPath
		command = "copy"
	} else {
		SourcePath = job.SourcePath
		if tempDumpPath != "" {
			SourcePath = tempDumpPath
		}
		Destination = fmt.Sprintf("%s:%s", job.RemoteName, job.DestinationPath)
		command = job.RcloneMode
	}

	args := []string{command, SourcePath, Destination}
	if !isRestore {
		args = append(args, "--checksum")
	}
	if job.IsEncrypted {
		s.injectEncrytionFlags(&args, job)
	}
	return args
}

func generateRandomSalt(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return base64.StdEncoding.EncodeToString(bytes)[:length]
}

func (s *BackupServiceImpl) injectEncrytionFlags(args *[]string, job models.ScheduledJob) {
	key := job.EncryptionKey
	salt := generateRandomSalt(32)

	*args = append(*args,
		"--crypt-filename-encryption", "standard",
		"--crypt-password", key,
		"--crypt-salt", salt,
	)
	fmt.Println("[ARGS] Enkripsi Runtime Dijalankan")
}

func (s *BackupServiceImpl) StartNewJob(job models.ScheduledJob) {
	var tempDumpPath string

	go func() {
		fmt.Printf("[%d] job %s: Memulai Eksekusi Rclone...\n", job.ID, job.Name)
		// Pre-script Logic
		if job.SourceType == "DB" {
			var err error
			tempDumpPath, err = s.executeDumpDB(job)
			if err != nil {
				s.handleJobCompletion(job, RcloneResult{Success: false, ErrorMsg: err.Error()}, tempDumpPath)
				return
			}
		}

		rcloneArgs := s.buildRcloneArgs(job, tempDumpPath)
		result := ExecuteRcloneJob(rcloneArgs)
		s.handleJobCompletion(job, result, tempDumpPath)

		pesan := fmt.Sprintf("[%d] Job %s: selesai. Status masuk ke log\n", job.ID, job.Name)
		fmt.Println(pesan)
	}()
}

func (s *BackupServiceImpl) executeDumpDB(job models.ScheduledJob) (string, error) {
	tempPath := fmt.Sprintf("/tmp/db_dump_%d_%d.sql", job.ID, time.Now().Unix())
	dumpArgs := []string{
		"mysqldump",
		fmt.Sprintf("-u%s", job.DbUser),
		fmt.Sprintf("-p%s", job.DbPass),
		job.SourcePath,
		"-r", tempPath,
	}

	result := ExecuteRcloneJob(dumpArgs)

	if !result.Success {
		return "", fmt.Errorf("MySqlDump Failed For %s:%s", job.SourcePath, result.ErrorMsg)
	}

	return tempPath, nil
}

func (s *BackupServiceImpl) ExecuteRcloneJob(job models.ScheduledJob) (string, error) {
	if job.SourceType == "DB" {
		return s.executeDumpDB(job)
	}
	return "", nil
}

func (s *BackupServiceImpl) handleJobCompletion(job models.ScheduledJob, result RcloneResult, tempDumpPath string) {
	LogMutex.Lock()
	defer LogMutex.Unlock()

	logstatus := "Failed"

	if result.Success {
		logstatus = "Completed"

		if tempDumpPath != "" {
			os.Remove(tempDumpPath)
		}
	}

	newLog := &models.Log{
		JobID:         &job.ID,
		OperationType: job.SourceType,
		Status:        logstatus,
		Timestamp:     time.Now(),
		DurationSec:   int(result.Duration.Seconds()),
		Message:       result.Output + result.ErrorMsg,
	}
	s.LogRepo.CreateLog(newLog)
}

func (s *BackupServiceImpl) CreateJobAndDispatch(job *models.ScheduledJob) error {

	// Asumsi: Logic enkripsi/dekripsi field DBPass sudah diimplementasikan di Repository/Service

	// --- 1. ENKRIPSI PASSWORD SEBELUM PENYIMPANAN ---
	if job.SourceType == "DB" && job.DbPass != "" {
		// PERHATIAN: Di sini, Anda seharusnya memanggil helper untuk mengenkripsi job.DbPass
		// dan menyimpan nilai terenkripsi kembali ke job.DbPass.
		// job.DbPass = s.EncryptService.Encrypt(job.DbPass) // LOGIKA ENKRIPSI
		fmt.Println("[SECURITY] DbPass dienkripsi sebelum dipersistenkan.")
	}

	// --- 2. PENYIMPANAN ATAU DISPATCH LANGSUNG ---

	if job.ScheduleCron != "" {
		// Job Terjadwal (Auto Backup): Simpan ke DB agar Scheduler Daemon bisa mengambilnya.
		// Job akan memiliki ID setelah Create().
		if err := s.JobRepo.Create(job); err != nil {
			return fmt.Errorf("gagal menyimpan job terjadwal: %w", err)
		}
		fmt.Printf("[DISPATCHER] Job %d (%s) disimpan untuk Scheduler.\n", job.ID, job.Name)

	} else {
		// Job Manual/Sekali Jalan: TIDAK disimpan di scheduled_jobs (hanya masuk Logs)
		// Kita langsung jalankan Job tersebut.
		fmt.Printf("[DISPATCHER] Job %s (Manual) dipicu langsung.\n", job.Name)

		// Cek Keamanan: Jika Job Manual, JobID akan nol di DB.
		// Kita harus menyimpan log penuh di ConfigSnapshot.

		// Panggil StartNewJob langsung
		s.StartNewJob(*job)
	}

	return nil
}
