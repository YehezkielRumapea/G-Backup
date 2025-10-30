package models

import "time"

// Log merepresentasikan satu catatan eksekusi
type Log struct {
	ID    uint  `gorm:"primaryKey;type:int unsigned"`
	JobID *uint `gorm:"index;type:int unsigned"` // Nullable

	// Status (Enum disesuaikan dengan SQL baru)
	Status string `gorm:"type:enum('SUCCESS', 'FAIL_PRE_SCRIPT', 'FAIL_RCLONE', 'FAIL_POST_SCRIPT', 'ERROR');not null"`

	ConfigSnapshot string    `gorm:"type:json;nullable"`
	Message        string    `gorm:"type:text"` // Output mentah CLI
	DurationSec    int       `gorm:"column:duration_sec;nullable"`
	Timestamp      time.Time `gorm:"column:timestamp;default:CURRENT_TIMESTAMP"`

	ScheduledJob ScheduledJob `gorm:"foreignKey:JobID"`
}
