package models

import "time"

// User merepresentasikan satu baris di tabel 'users' untuk otentikasi.
type User struct {
	ID           uint   `gorm:"primaryKey"`
	Username     string `gorm:"size:50;unique;not null"`
	PasswordHash string `gorm:"size:255;not null"`
	IsActive     bool   `gorm:"default:true"` // Untuk mengontrol status user

	// Timestamp standar GORM
	CreatedAt time.Time
	UpdatedAt time.Time
}
