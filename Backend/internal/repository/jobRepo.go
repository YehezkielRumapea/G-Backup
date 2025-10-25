package repository

import (
	"fmt"
	"gbackup-system/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// Kontrak Interface
type JobRepository interface {
	Create(job *models.ScheduledJob) error
	FindActiveJobs() ([]models.ScheduledJob, error)
	UpdateLastRunStatus(JobID uint, lastRunTime time.Time, status string) error
	UpdateEncryptionSalt(jobID uint, salt string) error
}

type JobRepositoryImpl struct {
	JobRepo *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &JobRepositoryImpl{JobRepo: db}
}

// Job Active untuk Schedular
func (r *JobRepositoryImpl) Create(Job *models.ScheduledJob) error {
	result := r.JobRepo.Create(Job)
	return result.Error
}

func (r *JobRepositoryImpl) FindActiveJobs() ([]models.ScheduledJob, error) {
	var jobs []models.ScheduledJob
	result := r.JobRepo.Where("is_Active= ? ", true).Find(&jobs)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return jobs, nil
}

func (r *JobRepositoryImpl) UpdateLastRunStatus(jobID uint, lastRun time.Time, status string) error {
	// ✅ Gunakan WHERE untuk cegah race condition
	result := r.JobRepo.Model(&models.ScheduledJob{}).
		Where("id = ? AND (status_queue != ? OR status_queue IS NULL)", jobID, "RUNNING").
		Updates(map[string]interface{}{
			"last_run":     lastRun,
			"status_queue": status,
		})

	if result.RowsAffected == 0 {
		return fmt.Errorf("job sudah running atau tidak ditemukan")
	}

	return result.Error
}

func (r *JobRepositoryImpl) UpdateEncryptionSalt(jobID uint, salt string) error {
	return r.JobRepo.Model(&models.ScheduledJob{}).
		Where("id = ?", jobID).
		Update("encryption_salt", salt).Error
}
