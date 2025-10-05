package repository

import (
	"gbackup-system/backend/internal/models"

	"gorm.io/gorm"
)

type LogRepository interface {
	FindAllLogs() ([]models.Log, error)
	CreateLog(log *models.Log) error
}

type LogRepositoryImpl struct {
	LogRepo *gorm.DB
}

func (r *LogRepositoryImpl) CreateLog(log *models.Log) error {
	result := r.LogRepo.Create(log)
	return result.Error
}

func NewLogRepository(db *gorm.DB) LogRepository {
	return &LogRepositoryImpl{LogRepo: db}
}

func (r *LogRepositoryImpl) FindAllLogs() ([]models.Log, error) {
	var logs []models.Log
	result := r.LogRepo.Preload("ScheduledJob").Order("timestamp desc").Find(&logs)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return logs, nil

}
