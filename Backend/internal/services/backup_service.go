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
	return (&BackupServiceImpl{JobRepo: jRepo, LogRepo: lRepo})
}

func (s *BackupServiceImpl) buildRcloneArgs(job models.ScheduledJob, tempDumpPath string) []string {
	// Path
	sourcePathRclone := job.SourcePath
	if tempDumpPath != "" {
		sourcePathRclone = tempDumpPath
	}
	destination := fmt.Sprintf("%s:%s", job.RemoteName, job.DestinationPath)
	command := job.RcloneMode

	// Arg Rclone

	args := []string{command, sourcePathRclone, destination}
	args = append(args, "--checksum")

	// if job.RcloneMode == "sync" {
	// 	fmt.Println("[ARGS] Mode Sync aktif (Overwrite)")
	// } else {
	// 	fmt.Println("[ARGS] Mode Copy aktif (Append)")
	// }

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
		if job.JobType == "DB" {
			var err error
			tempDumpPath, err = s.ExecuteRcloneJob(job)
			if err != nil {
				s.handleJobCompletion(job, RcloneResult{Success: false, ErrorMsg: err.Error()}, tempDumpPath)
				return
			}
		}

		rcloneArgs := ExecuteRcloneJob(rcloneArgs)
		s.handleJobCompletion(job, result, tempDumpPath)

		fmt.Sprintf("[%d] Job %s: selesai. Status masuk ke log\n", job.ID, job.Name)
	}()
}

func (s *BackupServiceImpl) executeDumpDB(job models.ScheduledJob) (string, error) {
	tempPath := fmt.Sprintf("/tmp/db_dump_%d_%d.sql", job.ID, time.Now().Unix())

	dumpArgs := []string{
		fmt.Sprintf("-u%s", job.DbUser),
		fmt.Sprintf("-p%s", job.Dbpass),
		job.SourcePath,
		"-r", tempPath,
	}

	result := ExecuteRcloneJob(dumpArgs)

	if !result.Success {
		return "", fmt.Errorf("MySqlDump Failed For %s:%s", job.SourcePath, result.ErrorMsg)
	}

	return tempPath, nil
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

	newLog := &models.Log{}
}
