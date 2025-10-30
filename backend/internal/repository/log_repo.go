package repository

import (
	"gbackup-new/backend/internal/models"

	"gorm.io/gorm"
)

// LogRepository mendefinisikan kontrak
type LogRepository interface {
	CreateLog(log *models.Log) error
	FindAllLogs() ([]models.Log, error)
}

type logRepositoryImpl struct {
	DB *gorm.DB
}

func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepositoryImpl{DB: db}
}

// CreateLog: Mencatat hasil eksekusi (dipanggil oleh Service)
func (r *logRepositoryImpl) CreateLog(log *models.Log) error {
	result := r.DB.Create(log)
	return result.Error
}

// FindAllLogs: Mengambil riwayat untuk UI Logs Page
func (r *logRepositoryImpl) FindAllLogs() ([]models.Log, error) {
	var logs []models.Log
	// Preload "ScheduledJob" untuk mengambil JobName (jika JobID tidak NULL)
	result := r.DB.Preload("ScheduledJob").Order("timestamp desc").Find(&logs)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return logs, nil
}
