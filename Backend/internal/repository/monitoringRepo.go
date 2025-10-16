package repository

import (
	"gbackup-system/backend/internal/models"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type MonitoringRepository interface {
	UpsertRemoteStatus(Monitoring *models.Monitoring) error
	FindAllRemote() ([]*models.Monitoring, error)
	FindAllLogs() ([]*models.Log, error)
}

type MonitoringRepositoryImpl struct {
	DB *gorm.DB
}

func NewMoniRepository(db *gorm.DB) MonitoringRepository {
	return &MonitoringRepositoryImpl{DB: db}
}

func (r *MonitoringRepositoryImpl) UpsertRemoteStatus(monitoring *models.Monitoring) error {
	result := r.DB.Clauses(clause.OnConflict{
		Columns: []clause.Column{{Name: "RemoteName"}},
		DoUpdates: clause.AssignmentColumns([]string{
			"StatusConnect",
			"TotalStorageGB",
			"UsedStorageGB",
			"FreeStorageGB",
			"LastCheckedAt",
		}),
	}).Create(monitoring)
	return result.Error
}

func (r *MonitoringRepositoryImpl) FindAllRemote() ([]*models.Monitoring, error) {
	var remotes []*models.Monitoring
	result := r.DB.Find(&remotes)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return remotes, nil
}

func (r *MonitoringRepositoryImpl) FindAllLogs() ([]*models.Log, error) {
	var logs []*models.Log
	result := r.DB.Find(&logs)
	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}
	return logs, nil
}
