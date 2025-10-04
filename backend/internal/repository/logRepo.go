package repository

import (
	"gbackup-system/backend/internal/models"

	"gorm.io/gorm"
)

type LogRepository interface {
	FindAllLogs() ([]models.Log, error)
	CreateLog(log *models.Log) error
}

type logRepositoryImpl struct {
	log *gorm.DB
}

func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepositoryImpl{DB: db}
}

func (r *logRepositoryImpl) FindAllLogs() ([]models.Log, error) {
	var logs []models.Log
	result := r.DB.Preload("Scheduled Job").order("timestamp desc").Find(&logs)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return logs, nil

}
