package repository

import (
	"gbackup-new/backend/internal/models"

	"gorm.io/gorm"
)

type RemoteRepository interface {
	Create(remote *models.Remote) error
}

type remoteRepositoryImpl struct {
	DB *gorm.DB
}

func NewRemoteRepository(db *gorm.DB) RemoteRepository {
	return &remoteRepositoryImpl{DB: db}
}

func (r *remoteRepositoryImpl) Create(remote *models.Remote) error {
	result := r.DB.Create(remote)
	return result.Error
}
