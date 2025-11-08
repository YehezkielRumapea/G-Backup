package repository

import (
	"errors"                              // Diperlukan untuk error handling GORM
	"gbackup-new/backend/internal/models" // Sesuaikan dengan path module Anda

	"gorm.io/gorm"
)

// UserRepository mendefinisikan kontrak (APA yang bisa dilakukan)
type UserRepository interface {
	CreateUser(user *models.User) error
	FindByUsername(username string) (*models.User, error)
}

// userRepositoryImpl adalah implementasi (BAGAIMANA dilakukan)
type userRepositoryImpl struct {
	DB *gorm.DB
}

// NewUserRepository adalah constructor (Factory) untuk Dependency Injection
// (Dipanggil dari main.go)
func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepositoryImpl{DB: db}
}

// ----------------------------------------------------
// FUNGSI IMPLEMENTASI
// ----------------------------------------------------

// CreateUser: Menjalankan SQL INSERT untuk user baru (Admin Seeder)
func (r *userRepositoryImpl) CreateUser(user *models.User) error {
	result := r.DB.Create(user)
	return result.Error // Mengembalikan error jika ada (misalnya, username sudah ada)
}

// FindByUsername: Mencari user berdasarkan username untuk proses login
func (r *userRepositoryImpl) FindByUsername(username string) (*models.User, error) {
	var user models.User

	// Menjalankan SQL SELECT ... WHERE username = ?
	result := r.DB.Where("username = ?", username).First(&user)

	// --- Error Handling Kritis ---
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			// User-nya memang tidak ada.
			return nil, nil
		}
		// Error database fatal
		return nil, result.Error
	}

	return &user, nil
}
