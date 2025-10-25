package models

import "time"

// Log merepresentasikan satu catatan eksekusi dari setiap Job.
type Log struct {
	ID uint `gorm:"primaryKey;type:int unsigned"`

	// KRITIS: Pointer (*uint) untuk Nullability Job Manual/Restore
	JobID *uint `gorm:"index;type:int unsigned"`

	// Status Operasi
	OperationType string `gorm:"type:enum('BACKUP', 'RESTORE', 'MANUAL_BACKUP');not null"`
	Status        string `gorm:"type:enum('SUCCESS', 'FAIL', 'ERROR');not null"`

	// Data Penting untuk Debugging
	Message     string `gorm:"type:text"`
	DurationSec int    `gorm:"column:duration_sec;nullable"`
	Checksum    string `gorm:"size:255;nullable"`

	// KRUSIAL: Menyimpan konfigurasi Job Manual/Restore
	ConfigSnapshot *string `gorm:"type:json;nullable"`

	Timestamp time.Time `gorm:"column:timestamp;default:CURRENT_TIMESTAMP"`
}
