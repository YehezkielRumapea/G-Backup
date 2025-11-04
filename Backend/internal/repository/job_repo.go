package repository

import (
	"errors"
	"fmt"
	"gbackup-new/backend/internal/models" // Sesuaikan path module
	"time"

	"gorm.io/gorm"
)

// JobRepository mendefinisikan kontrak untuk operasi Job
type JobRepository interface {
	Create(job *models.ScheduledJob) error
	FindJobByID(jobID uint) (*models.ScheduledJob, error)
	FindAllActiveJobs() ([]models.ScheduledJob, error)
	UpdateLastRunStatus(jobID uint, lastRunTime time.Time, status string) error
}

type jobRepositoryImpl struct {
	DB *gorm.DB
}

func NewJobRepository(db *gorm.DB) JobRepository {
	return &jobRepositoryImpl{DB: db}
}

// Create: Menyimpan template Job baru
func (r *jobRepositoryImpl) Create(job *models.ScheduledJob) error {
	result := r.DB.Create(job)
	return result.Error
}

// FindJobByID: Mengambil satu job berdasarkan ID (untuk preview script / trigger manual)
func (r *jobRepositoryImpl) FindJobByID(jobID uint) (*models.ScheduledJob, error) {
	var job models.ScheduledJob
	result := r.DB.First(&job, jobID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("job ID %d tidak ditemukan", jobID)
		}
		return nil, result.Error
	}
	return &job, nil
}

// FindAllActiveJobs: Mengambil job yang siap dieksekusi oleh Scheduler
func (r *jobRepositoryImpl) FindAllActiveJobs() ([]models.ScheduledJob, error) {
	var jobs []models.ScheduledJob
	// Ambil Job yang aktif dan BUKAN Job Manual (karena Job Manual tidak punya cron)
	result := r.DB.Where("is_active = ? AND schedule_cron IS NOT NULL", true).Find(&jobs)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return jobs, nil
}

// UpdateLastRunStatus: Mengupdate status (PENDING, RUNNING, COMPLETED, FAILED_PRE...)
func (r *jobRepositoryImpl) UpdateLastRunStatus(jobID uint, lastRunTime time.Time, status string) error {
	result := r.DB.Model(&models.ScheduledJob{}).
		Where("id = ?", jobID).
		Updates(map[string]interface{}{
			"last_run_at":  lastRunTime,
			"status_queue": status,
		})
	return result.Error
}
