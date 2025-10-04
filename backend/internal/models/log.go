package models

import "time"

// Log merepresentasikan satu catatan eksekusi dari setiap Job.
type Log struct {
	ID uint `gorm:"primaryKey"`

	// Foreign Key (Pointer membuatnya nullable di MariaDB)
	JobID   uint
	JobName string
	// Status Operasi
	OperationType string `gorm:"type:enum('BACKUP', 'RESTORE', 'MANUAL_BACKUP');not null"`
	Status        string `gorm:"type:enum('SUCCESS', 'FAIL', 'ERROR');not null"`

	// Data Penting untuk Debugging dan Logging
	Message     string `gorm:"type:text"`
	DurationSec int    `gorm:"nullable"`          // Durasi eksekusi
	Checksum    string `gorm:"size:255;nullable"` // Checksum transfer

	// KRUSIAL: Menyimpan konfigurasi untuk Job Manual/Restore (jika JobID NULL)
	ConfigSnapshot string `gorm:"type:json;nullable"`

	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP"`

	// Relasi
	ScheduledJob ScheduledJob `gorm:"foreignKey:JobID"`
}
