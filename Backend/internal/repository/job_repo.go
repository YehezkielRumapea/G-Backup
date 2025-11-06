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
	UpdateJobActivity(JobID uint, isActive bool) error
	CountJobOnRemote(remoteName string) (int64, error)
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
	result := r.DB.Where("schedule_cron IS NOT NULL AND schedule_cron != ?", "").Find(&jobs)

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

func (r *jobRepositoryImpl) UpdateJobActivity(jobID uint, isActive bool) error {
	result := r.DB.Model(&models.ScheduledJob{}).
		Where("id = ?", jobID).
		Update("is_active", isActive)
	return result.Error
}

func (r *jobRepositoryImpl) CountJobOnRemote(remoteName string) (int64, error) {
	var count int64

	err := r.DB.Model(&models.ScheduledJob{}).
		Where("remote_name = ? AND schedule_cron != ? AND is_active = ?",
			remoteName, "", true).
		Count(&count).Error

	if err != nil {
		// Log error atau tangani sesuai kebijakan aplikasi
		return 0, fmt.Errorf("gagal menghitung job terjadwal untuk remote %s: %w", remoteName, err)
	}

	return count, nil
}
