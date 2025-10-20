package models

import "time"

type ScheduledJob struct {
	ID     uint   `gorm:"primaryKey"`
	UserID uint   `gorm:"not null"`
	Name   string `gorm:"size:100;not null"`

	// Konfigurasi Sumber
	SourceType string `gorm:"type:enum('FILE','DB');size:10;not null"`
	SourcePath string `gorm:"size:255;not null"` // Path Lokal atau Nama DB
	DbUser     string // Dibutuhkan untuk mysqldump (Password di-handle via env/config)
	DbPass     string // Dibutuhkan untuk mysqldump

	// Remote dan Enkripsi
	RemoteID      uint   `gorm:"not null"` // Foreign Key ke tabel Remote (jika Anda membuatnya)
	RemoteName    string `gorm:"size:100;not null"`
	IsEncrypted   bool   `gorm:"default:false"`
	EncryptionKey string `gorm:"size:255;nullable"` // Kunci Enkripsi Runtime

	// Penjadwalan dan Status Antrian
	Schedule    string `gorm:"size:50;not null"` // Cron Expression (Tidak NULL karena ini tabel Scheduled)
	Priority    int    `gorm:"default:5"`
	StatusQueue string `gorm:"type:enum('PENDING','RUNNING','COMPLETED','FAILED');default:PENDING"`
	IsActive    bool   `gorm:"default:true"`
	RcloneMode  string `gorm:"type:enum('copy','sync','move');default:'copy'"`

	// Konfigurasi Tujuan
	DestinationPath string `gorm:"size:255;not null"` // Path di Remote
	JobType         string `gorm:"type:enum('FILE','DB');not null"`
	// Timestamp
	LastRun *time.Time // Waktu Job terakhir kali berjalan (bisa NULL)

	// Bawaan GORM
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relasi (Asumsi Anda membuat struct Remote)
	Monitoring Monitoring `gorm:"foreignKey:RemoteID"`
	User       User       `gorm:"foreignKey:UserID"`
}
