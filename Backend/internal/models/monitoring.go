package models

import "time"

// Monitoring merepresentasikan satu baris di tabel 'monitoring'
type Monitoring struct {
	ID             uint      `gorm:"primaryKey;type:int unsigned"`
	RemoteName     string    `gorm:"column:remote_name;size:100;unique;not null"`
	StatusConnect  string    `gorm:"column:status_connect;type:enum('CONNECTED','DISCONNECTED');not null"`
	TotalStorageGB float64   `gorm:"column:total_storage_gb;default:0"`
	UsedStorageGB  float64   `gorm:"column:used_storage_gb;default:0"`
	FreeStorageGB  float64   `gorm:"column:free_storage_gb;default:0"`
	LastCheckedAt  time.Time `gorm:"column:last_checked_at"`
}
