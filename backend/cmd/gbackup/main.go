package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath" // Untuk path .env

	"gbackup-new/backend/pkg/database" // Sesuaikan path module baru Anda

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func loadDotEnv() {
	// --- 1. MEMUAT .ENV DARI LOKASI YANG BENAR ---
	// Binary ada di cmd/gbackup, .env ada di ../../.env
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Tidak bisa dapatkan path executable:", err)
		return
	}
	exeDir := filepath.Dir(exePath)

	// Cari .env di: direktori binary, parent, dan parent’s parent
	candidates := []string{
		filepath.Join(exeDir, ".env"),
		filepath.Join(exeDir, "..", ".env"),
		filepath.Join(exeDir, "..", "..", ".env"),
	}

	for _, p := range candidates {
		if _, err := os.Stat(p); err == nil {
			if err := godotenv.Load(p); err == nil {
				fmt.Println("Loaded .env dari:", p)
				return
			}
		}
	}
}

func main() {
	loadDotEnv()
	// --- 2. MEMULAI FONDASI DATABASE ---
	dbInstance := database.Connect()

	if dbInstance == nil {
		log.Fatal("Koneksi DB gagal, instance GORM nil.")
	}

	// --- 3. SETUP WEB SERVER ECHO (Minimal) ---
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "G-Backup New Architecture is Running! DB Connected.")
	})

	fmt.Println("Backend diinisialisasi. Menjalankan Echo server di port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
