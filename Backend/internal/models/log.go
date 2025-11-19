package models

import "time"

// Log merepresentasikan satu catatan eksekusi
type Log struct {
	ID    uint  `gorm:"primaryKey;type:int unsigned"`
	JobID *uint `gorm:"column:job_id;index;type:int unsigned"`
	// Status (Enum disesuaikan dengan SQL baru)
	Status           string    `gorm:"type:enum('SUCCESS', 'FAIL_PRE_SCRIPT', 'FAIL_RCLONE', 'FAIL_POST_SCRIPT', 'FAIL_SOURCE_CHECK', 'ERROR');not null"`
	ConfigSnapshot   *string   `gorm:"type:json;nullable"`
	Message          string    `gorm:"type:text"` // Output mentah CLI
	DurationSec      int       `gorm:"column:duration_sec;nullable"`
	TransferredBytes int64     `gorm:"column:transferred_bytes;default:0"`
	Timestamp        time.Time `gorm:"column:timestamp;default:CURRENT_TIMESTAMP"`

	ScheduledJob ScheduledJob `gorm:"foreignKey:JobID"`
}
