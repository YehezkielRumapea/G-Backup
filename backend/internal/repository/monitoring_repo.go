package repository

import (
	"gbackup-new/backend/internal/models" // Sesuaikan path module

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// MonitoringRepository mendefinisikan kontrak (APA yang bisa dilakukan)
type MonitoringRepository interface {
	UpsertRemoteStatus(monitor *models.Monitoring) error
	FindAllRemotes() ([]models.Monitoring, error)
}

// monitoringRepositoryImpl adalah implementasi (BAGAIMANA dilakukan)
type monitoringRepositoryImpl struct {
	DB *gorm.DB
}

// NewMonitoringRepository adalah constructor (Factory)
func NewMonitoringRepository(db *gorm.DB) MonitoringRepository {
	return &monitoringRepositoryImpl{DB: db}
}

// ----------------------------------------------------
// FUNGSI IMPLEMENTASI
// ----------------------------------------------------

// UpsertRemoteStatus: Melakukan Update OR Insert (UPSERT)
// Ini dipanggil oleh Service setelah mendapat output dari 'rclone about'.
func (r *monitoringRepositoryImpl) UpsertRemoteStatus(monitor *models.Monitoring) error {
	// Logic Kritis: Jika remote_name sudah ada, update. Jika belum, insert.
	result := r.DB.Clauses(clause.OnConflict{
		// 1. Kolom yang diuji konflik: remote_name harus unik.
		Columns: []clause.Column{{Name: "remote_name"}},

		// 2. Aksi jika konflik terjadi: UPDATE kolom-kolom ini.
		DoUpdates: clause.AssignmentColumns([]string{
			"status_connect",
			"total_storage_gb",
			"used_storage_gb",
			"free_storage_gb",
			"last_checked_at",
		}),
	}).Create(monitor)

	return result.Error
}

// FindAllRemotes: Mengambil semua data monitoring untuk UI Remote Monitoring.
func (r *monitoringRepositoryImpl) FindAllRemotes() ([]models.Monitoring, error) {
	var remotes []models.Monitoring

	result := r.DB.Find(&remotes)

	if result.Error != nil && result.Error != gorm.ErrRecordNotFound {
		return nil, result.Error
	}

	return remotes, nil
}
