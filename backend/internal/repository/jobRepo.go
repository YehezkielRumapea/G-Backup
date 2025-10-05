package repository

import (
	"gbackup-system/backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// Kontrak Interface
type JobRepository interface {
	Create(job *models.ScheduledJob) error
	FindActiveJobs() ([]models.ScheduledJob, error)
	UpdateLastRunStatus(JobID uint, lastRunTime time.Time, status string) error
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

func (r *JobRepositoryImpl) UpdateLastRunStatus(JobID uint, lastRunTime time.Time, status string) error {
	result := r.JobRepo.Model(&models.ScheduledJob{}).Where("id =?", JobID).Updates(map[string]interface{}{
		"last_run_at":  lastRunTime,
		"status_queue": status,
	})
	return result.Error
}
