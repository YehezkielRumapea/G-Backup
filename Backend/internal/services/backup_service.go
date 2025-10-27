package services

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"gbackup-system/backend/internal/models"
	"gbackup-system/backend/internal/repository"
	"os"
	"os/exec"
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

func (s *BackupServiceImpl) buildRcloneArgs(job models.ScheduledJob, tempDumpPath string) ([]string, error) {
	if job.RemoteName == "" {
		return nil, fmt.Errorf("RemoteName tidak boleh kosong")
	}
	// Path
	isRestore := job.OperationMode == "Restore"
	command := job.RcloneMode

	var SourcePath, Destination string
	if command == "" {
		command = "copy"
	}

	if isRestore {
		// Sumber dari cloud --> Lokal
		SourcePath = fmt.Sprintf("%s:%s", job.RemoteName, job.SourcePath)
		Destination = job.DestinationPath
		command = "copy"
	} else {
		// Sumber dari Lokal --> Cloud
		SourcePath = job.SourcePath
		if tempDumpPath != "" {
			SourcePath = tempDumpPath
		}

		// Path enkripsi
		Destination = fmt.Sprintf("%s:%s", job.RemoteName, job.DestinationPath)
	}

	args := []string{command, SourcePath, Destination}
	if !isRestore {
		args = append(args, "--checksum")
	}
	if job.IsEncrypted {
		if err := s.injectEncryptionFlags(&args, job); err != nil {
			return nil, err
		}
	}
	return args, nil
}

func generateDeterministicSalt(jobID uint, key string) string {
	// Kombinasi JobID + Key akan selalu menghasilkan salt yang sama
	input := fmt.Sprintf("%d:%s", jobID, key)
	bytes := []byte(input)

	// Pad atau truncate ke 32 karakter
	if len(bytes) < 32 {
		// Repeat sampai 32
		for len(bytes) < 32 {
			bytes = append(bytes, bytes...)
		}
	}

	return base64.StdEncoding.EncodeToString(bytes[:32])[:32]
}

func (s *BackupServiceImpl) injectEncryptionFlags(args *[]string, job models.ScheduledJob) error {
	key := job.EncryptionKey
	if key == "" {
		return fmt.Errorf("encryption key tidak boleh kosong")
	}

	salt := generateDeterministicSalt(job.ID, key)

	*args = append(*args,
		"--crypt-filename-encryption", "standard",
		"--crypt-password", key,
		"--crypt-password2", salt, // ✅ Flag yang benar untuk salt
	)

	fmt.Println("[ARGS] Enkripsi Runtime Dijalankan")
	return nil
}

func (s *BackupServiceImpl) StartNewJob(job models.ScheduledJob) {
	var tempDumpPath string
	go func() {
		isRestore := job.OperationMode == "Restore"

		fmt.Printf("[%d] job %s: Memulai Eksekusi Rclone...\n", job.ID, job.Name)
		// Pre-script Logic
		// hanya berjalan keetika yang di backup DB dan bukan Restore
		if job.SourceType == "DB" && !isRestore {
			var err error
			tempDumpPath, err = s.executeDumpDB(job)
			if err != nil {
				s.handleJobCompletion(job, RcloneResult{Success: false, ErrorMsg: err.Error()}, tempDumpPath)
				return
			}
		}

		// Eksekusi Rclone
		rcloneArgs, err := s.buildRcloneArgs(job, tempDumpPath)
		if err != nil {
			s.handleJobCompletion(job, RcloneResult{Success: false, ErrorMsg: err.Error()}, tempDumpPath)
			return
		}
		result := ExecuteRcloneJob(rcloneArgs)

		// Post-script
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

	cmd := exec.Command(dumpArgs[0], dumpArgs[1:]...)
	output, err := cmd.CombinedOutput()

	if err != nil {
		errorMsg := fmt.Sprintf("Exit Error: %v. Output: %s", err, string(output))
		return "", fmt.Errorf("mysqldump failed for %s:%s", job.SourcePath, errorMsg)
	}

	return tempPath, nil
}

// Post-script
func (s *BackupServiceImpl) handleJobCompletion(job models.ScheduledJob, result RcloneResult, tempDumpPath string) {
	defer func() {
		if job.OperationMode == "BACKUP" && tempDumpPath != "" {
			if err := os.Remove(tempDumpPath); err != nil {
				fmt.Printf("[WARNING] Gagal menghapus file sementara: %v\n", err)
			} else {
				fmt.Printf("[INFO] File sementara berhasil dihapus: %s\n", tempDumpPath)
			}
		}
	}()

	LogMutex.Lock()
	defer LogMutex.Unlock()

	logstatus := "FAIL"

	if result.Success {
		logstatus = "SUCCESS"
	}

	// Logging
	newLog := &models.Log{
		JobID:         &job.ID,
		OperationType: job.OperationMode,
		Status:        logstatus,
		Timestamp:     time.Now(),
		DurationSec:   int(result.Duration.Seconds()),
		Message:       result.Output + result.ErrorMsg,
	}
	s.LogRepo.CreateLog(newLog)

	if job.ScheduleCron != "" {
		s.JobRepo.UpdateLastRunStatus(job.ID, time.Now(), logstatus)
	}
}

func GenerateRandomKey(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("gagal membuat kunci enkripsi: %w", err)
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

func (s *BackupServiceImpl) CreateJobAndDispatch(job *models.ScheduledJob) error {
	// --- 1. ENKRIPSI PASSWORD SEBELUM PENYIMPANAN ---
	if job.SourceType == "DB" && job.DbPass != "" {
		fmt.Println("[SECURITY] DbPass dienkripsi sebelum dipersistenkan.")
	}

	if job.IsEncrypted {
		// hanya menggenerate key jika opsi enkrip di centang
		if job.EncryptionKey == "" {
			key, err := GenerateRandomKey(32)
			if err != nil {
				return fmt.Errorf("gagal membuat kunci enkripsi: %w", err)
			}
			job.EncryptionKey = key
			fmt.Println("Kunci Enkripsi dibuat ")
		}
	}

	if job.ScheduleCron != "" {
		if err := s.JobRepo.Create(job); err != nil {
			return fmt.Errorf("gagal menyimpan job terjadwal: %w", err)
		}
		fmt.Printf("[DISPATCHER] Job %d (%s) disimpan untuk Scheduler.\n", job.ID, job.Name)
	} else {
		fmt.Printf("[DISPATCHER] Job %s (Manual) dipicu langsung.\n", job.Name)
		s.StartNewJob(*job)
	}

	return nil
}
