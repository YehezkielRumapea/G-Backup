package models

import "time"

// ScheduledJob merepresentasikan satu Job terjadwal di tabel 'scheduled_jobs'
type ScheduledJob struct {
	ID     uint   `gorm:"primaryKey;type:int unsigned"`
	UserID uint   `gorm:"index;type:int unsigned;not null"`
	Name   string `gorm:"column:job_name;size:100;not null"`

	// Konfigurasi Sumber & Kredensial
	SourceType string `gorm:"column:job_type;type:enum('FILE','DB');size:10;not null"`
	SourcePath string `gorm:"size:255;not null"`
	DbUser     string `gorm:"size:50"`
	DbPass     string `gorm:"size:255"` // Akan disimpan terenkripsi di DB jika Job Auto

	// Remote dan Enkripsi
	RemoteName      string `gorm:"size:100;not null"`
	DestinationPath string `gorm:"size:255;not null"`
	IsEncrypted     bool   `gorm:"default:false"`
	EncryptionKey   string `gorm:"size:255;nullable"` // Kunci Enkripsi Runtime

	// RemoteID: Dipertahankan untuk relasi yang mungkin tidak digunakan secara langsung
	RemoteID      uint   `gorm:"index;type:int unsigned"`
	OperationMode string `gorm:"type:enum('BACKUP','MANUAL_BACKUP');not null;default:'BACKUP'"`

	// Penjadwalan dan Status Antrian
	ScheduleCron string `gorm:"size:50;not null"`
	Priority     int    `gorm:"default:5"`
	StatusQueue  string `gorm:"type:enum('PENDING','RUNNING','COMPLETED','FAILED');default:'PENDING'"`
	IsActive     bool   `gorm:"default:true"`
	RcloneMode   string `gorm:"column:rclone_mode;type:enum('copy','sync','move');default:'copy'"`

	// Timestamp dan Relasi
	LastRun   *time.Time `gorm:"column:last_run_at;nullable"`
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relasi (GORM)
	User User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Logs []Log `gorm:"foreignKey:JobID;references:ID"`
}
