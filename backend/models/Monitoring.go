package models

import "time"

// Monitoring merepresentasikan satu baris di tabel 'monitoring'
type Monitoring struct {
	ID uint `gorm:"primaryKey"`

	// Nama Remote GDrive
	RemoteName string `gorm:"size:100;unique;not null"`

	// Status Koneksi
	StatusConnect string `gorm:"type:enum('CONNECTED','DISCONNECTED');not null"`

	// Metrik Storage (Menggunakan float64 untuk presisi GB/TB)
	TotalStorageGB float64 `gorm:"default:0"`
	UsedStorageGB  float64 `gorm:"default:0"`
	FreeStorageGB  float64 `gorm:"default:0"` // Kolom yang bisa dihitung ulang

	// Timestamp
	LastCheckedAt time.Time
}

// Catatan:
// 1. Field JobCount akan ditambahkan ke struct terpisah (Response DTO)
//    di lapisan Service/Handler setelah dihitung.
// 2. Struct ini diganti nama menjadi 'Monitoring' agar konsisten dengan tabel DB.
