package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"

	// Import semua package yang dibutuhkan
	"gbackup-new/backend/internal/handler" // Dibutuhkan untuk AutoMigrate
	"gbackup-new/backend/internal/repository"
	"gbackup-new/backend/internal/service"
	"gbackup-new/backend/middleware"
	"gbackup-new/backend/pkg/database"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
	echomid "github.com/labstack/echo/v4/middleware"
)

const DefaultAdminUsername = "admin"
const DefaultAdminPassword = "admin123" // Ganti ini di produksi

func main() {
	// --- 1. MEMUAT .ENV ---
	err := godotenv.Load(filepath.Join("..", "..", ".env")) // Path: backend/.env
	if err != nil {
		log.Fatal("KRITIS: Gagal memuat file .env di root /backend.")
	}

	// Ambil Kunci Rahasia JWT (Wajib)
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("FATAL: JWT_SECRET_KEY environment variable tidak diatur.")
	}

	// --- 2. KONEKSI DATABASE ---
	dbInstance := database.Connect()
	if dbInstance == nil {
		log.Fatal("Koneksi DB gagal, instance GORM nil.")
	}

	// --- 3. DEPENDENCY INJECTION (DI) ---

	// 3.1. Inisialisasi Repository (Lapisan Data Access)
	userRepo := repository.NewUserRepository(dbInstance)
	jobRepo := repository.NewJobRepository(dbInstance)
	logRepo := repository.NewLogRepository(dbInstance)
	monitorRepo := repository.NewMonitoringRepository(dbInstance)
	// remoteRepo := repository.NewRemoteRepository(dbInstance) // Untuk "Add Remote"

	// 3.2. Inisialisasi Service (Lapisan Logika Bisnis)
	authSvc := service.NewAuthService(userRepo)
	monitorSvc := service.NewMonitoringService(monitorRepo, logRepo)
	backupSvc := service.NewBackupService(jobRepo, logRepo) // Orkestrator 3 Fase
	schedulerSvc := service.NewSchedulerService(jobRepo, backupSvc)
	// remoteSvc := service.NewRemoteService(remoteRepo) // Service "Add Remote"

	// 3.3. Inisialisasi Handler (Controller Layer)
	authHandler := handler.NewAuthHandler(authSvc)
	monitorHandler := handler.NewMonitoringHandler(monitorSvc)
	jobHandler := handler.NewJobHandler(schedulerSvc, backupSvc) // Membutuhkan BackupSvc untuk TriggerManual
	backupHandler := handler.NewBackupHandler(backupSvc)
	restoreHandler := handler.NewRestoreHandler(backupSvc)
	// remoteHandler := handler.NewRemoteHandler(remoteSvc)

	// --- 4. SEEDING ADMIN AWAL ---
	if err := authSvc.RegisterAdmin(DefaultAdminUsername, DefaultAdminPassword); err != nil {
		fmt.Printf("Admin Seeding Status: %v\n", err)
	} else {
		fmt.Println("✅ Admin user berhasil dibuat.")
	}

	// --- 5. SETUP WEB SERVER ECHO ---
	e := echo.New()
	e.Use(echomid.Logger())
	e.Use(echomid.Recover())
	e.Use(echomid.CORSWithConfig(echomid.CORSConfig{
		AllowOrigins: []string{"http://localhost:5173"}, // URL Frontend
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodDelete},
	}))

	// --- 6. SETUP ROUTING ---
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "G-Backup New Architecture is Running!")
	})

	// Rute Publik (Login)
	e.POST("/api/v1/auth/login", authHandler.Login)

	// Rute Privat (Dilindungi JWT)
	r := e.Group("/api/v1")
	r.Use(middleware.JWTGuard()) // Menerapkan Middleware JWT

	// Rute Monitoring dan Logs
	r.GET("/monitoring/remotes", monitorHandler.GetRemoteStatusList)
	r.GET("/monitoring/logs", monitorHandler.GetJobLogs)

	// Rute Job Management (Create, List, Trigger)
	r.GET("/jobs/scheduled", jobHandler.GetScheduledJobs)    // List Job Monitoring
	r.GET("/jobs/script/:id", jobHandler.GetJobScript)       // Pratinjau Script
	r.POST("/jobs/trigger/:id", jobHandler.TriggerManualJob) // Tombol "Run Now"

	// Rute Aksi
	r.POST("/jobs/new", backupHandler.CreateNewJob)        // Create Backup (Manual/Auto)
	r.POST("/jobs/restore", restoreHandler.TriggerRestore) // Create Restore

	// Rute Konfigurasi
	r.POST("/remotes/new", remoteHandler.AddNewRemote) // Add New Remote

	// --- 7. START DAEMONS (Goroutines) ---
	schedulerSvc.StartDaemon() // Memulai CRON Daemon
	// (Tambahkan Goroutine untuk auto-update monitoring jika perlu)

	fmt.Println("\nBackend diinisialisasi. Menjalankan Echo server di port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
