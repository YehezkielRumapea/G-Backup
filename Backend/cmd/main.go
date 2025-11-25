package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	// Import semua package yang dibutuhkan
	"gbackup-new/backend/internal/handler"
	"gbackup-new/backend/internal/repository"
	"gbackup-new/backend/internal/service"
	"gbackup-new/backend/pkg/database"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echomid "github.com/labstack/echo/v4/middleware"
)

func loadDotEnv() {
	exePath, err := os.Executable()
	if err != nil {
		fmt.Println("Tidak bisa dapatkan path executable:", err)
		return
	}
	exeDir := filepath.Dir(exePath)

	// Searching .env di lokal
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

	_ = godotenv.Load(".env", "../.env")
}
func main() {
	// --- 1. MEMUAT .ENV ---
	loadDotEnv()

	// Ambil Kunci JWT (Wajib)
	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("FATAL: JWT_SECRET_KEY ada")
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
	browserRepo := repository.NewBrowserRepository()

	// 3.2. Inisialisasi Service (Lapisan Logika Bisnis)
	authSvc := service.NewAuthService(userRepo, jwtSecretKey)
	monitorSvc := service.NewMonitoringService(monitorRepo, logRepo, jobRepo)
	backupSvc := service.NewBackupService(jobRepo, logRepo, monitorRepo, monitorSvc) // Orkestrator 3 Fase
	schedulerSvc := service.NewSchedulerService(jobRepo, backupSvc)
	browserSvc := service.NewBrowserService(browserRepo)

	// 3.3. Inisialisasi Handler (Controller Layer)
	authHandler := handler.NewAuthHandler(authSvc)
	monitorHandler := handler.NewMonitoringHandler(monitorSvc, schedulerSvc, logRepo)
	jobHandler := handler.NewJobHandler(schedulerSvc, backupSvc, jobRepo) // Membutuhkan BackupSvc untuk TriggerManual
	backupHandler := handler.NewBackupHandler(backupSvc)
	restoreHandler := handler.NewRestoreHandler(backupSvc)
	browserHandler := handler.NewBrowserHandler(browserSvc)
	setupHandler := handler.NewSetupHandler(authSvc)

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
	setupGroup := e.Group("/api/v1/setup")
	setupGroup.GET("/status", setupHandler.GetSetupStatus)          // Cek apakah admin sudah ada
	setupGroup.POST("/register", setupHandler.RegisterInitialAdmin) // Registrasi admin pertama

	// Rute Privat
	r := e.Group("/api/v1")
	r.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSecretKey),
	})) // Menerapkan Middleware JWT

	// Rute Monitoring dan Logs
	r.GET("/monitoring/remotes", monitorHandler.GetRemoteStatusList)
	r.GET("/monitoring/logs", monitorHandler.GetJobLogs)
	r.GET("/monitoring/jobs", monitorHandler.GetScheduledJobs)
	r.GET("/monitoring/drivemail", monitorHandler.GetRemotes)

	// Rute Job Management (Create, List, Trigger)
	r.GET("/jobs/scheduled", jobHandler.GetScheduledJobs)    // List Job Monitoring
	r.GET("/jobs/script/:id", jobHandler.GetJobScript)       // Pratinjau Script
	r.POST("/jobs/trigger/:id", jobHandler.TriggerManualJob) // Tombol "Run Now"
	r.GET("/jobs/manual", jobHandler.GetManualJob)           // List Job Manual
	r.DELETE("/jobs/delete/:id", jobHandler.DeleteJob)       // Hapus Job
	r.PUT("/jobs/update/:id", jobHandler.UpdateJob)          // Update Job
	r.GET("/jobs/:id", jobHandler.GetJobByID)                // Get Job by ID
	r.GET("/jobs/alljobs", monitorHandler.GetAllJobs)        // Get All Jobs (scheduled + Manual)

	// Rute Aksi
	r.POST("/jobs/new", backupHandler.CreateNewJob)               // Create Backup (Manual/Auto)
	r.POST("/jobs/restore", restoreHandler.TriggerRestore)        // Create Restore
	r.GET("/browser/files", browserHandler.ListFiles)             //browse files
	r.GET("/browser/remotes", browserHandler.GetAvailableRemotes) //list remotes
	r.GET("/browser/info", browserHandler.GetFileInfo)            //file info

	// --- 7. START DAEMONS (Goroutines) ---
	schedulerSvc.StartDaemon() // Memulai CRON Daemon
	monitorSvc.StartMonitoringDaemon()

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("[Startup] Inisialisasi status remote monitoring...")

		if err := monitorSvc.SyncRemotesWithRclone(); err != nil {
			fmt.Printf("❌ ERROR: Gagal sync remote: %v\n", err)
			return
		}

		remotes, err := monitorSvc.GetRemoteStatusList()
		if err != nil {
			fmt.Printf("❌ ERROR: Gagal mengambil remote: %v\n", err)
			return
		}

		if len(remotes) > 0 {
			fmt.Printf("✅ Menemukan %d remote, memulai update status...\n", len(remotes))
			for _, remote := range remotes {
				go monitorSvc.UpdateRemoteStatus(remote.RemoteName)
			}
		} else {
			fmt.Println("⚠️ Tidak ada remote yang dikonfigurasi di rclone.conf")
		}

		fmt.Println("✅ Inisialisasi remote monitoring selesai.")
	}()

	fmt.Println("\nBackend diinisialisasi. Menjalankan Echo server di port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
