package repository

import (
	"errors"
	"gbackup-system/backend/internal/models"

	"gorm.io/gorm"
)

type UserRepository interface {
	CreateUser(user *models.User) error
	FindByUsername(Username string) (*models.User, error)
}

type UserRepositoryImpl struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &UserRepositoryImpl{DB: db}
}

// Menambah User baru ke DB
func (r *UserRepositoryImpl) CreateUser(user *models.User) error {
	result := r.DB.Create(user)
	return result.Error
}

// Mencari User berdasarkan Username untuk Auth
func (r *UserRepositoryImpl) FindByUsername(Username string) (*models.User, error) {
	var user models.User

	result := r.DB.Where("usesrname = ?", Username).First(&user)

	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, result.Error
	}

	return &user, nil
}
