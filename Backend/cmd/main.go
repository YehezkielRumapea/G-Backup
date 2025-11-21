package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	// Import semua package yang dibutuhkan
	"gbackup-new/backend/internal/handler" // Dibutuhkan untuk AutoMigrate
	"gbackup-new/backend/internal/repository"
	"gbackup-new/backend/internal/service"
	"gbackup-new/backend/pkg/database"

	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echomid "github.com/labstack/echo/v4/middleware"
)

const DefaultAdminUsername = "admin"
const DefaultAdminPassword = "admin123" // Ganti ini di produksi

func loadDotEnv() {
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

	// fallback: coba relative ke CWD (kalau kebetulan cocok)
	_ = godotenv.Load(".env", "../.env")
}
func main() {
	// --- 1. MEMUAT .ENV ---
	loadDotEnv()

	// Ambil Kunci Rahasia JWT (Wajib)
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
	// remoteRepo := repository.NewRemoteRepository(dbInstance) // Untuk "Add Remote"

	// 3.2. Inisialisasi Service (Lapisan Logika Bisnis)
	authSvc := service.NewAuthService(userRepo, jwtSecretKey)
	monitorSvc := service.NewMonitoringService(monitorRepo, logRepo, jobRepo)
	backupSvc := service.NewBackupService(jobRepo, logRepo, monitorRepo, monitorSvc) // Orkestrator 3 Fase
	schedulerSvc := service.NewSchedulerService(jobRepo, backupSvc)
	browserSvc := service.NewBrowserService(browserRepo)
	// remoteSvc := service.NewRemoteService(remoteRepo) // Service "Add Remote"

	// 3.3. Inisialisasi Handler (Controller Layer)
	authHandler := handler.NewAuthHandler(authSvc)
	monitorHandler := handler.NewMonitoringHandler(monitorSvc, schedulerSvc, logRepo)
	jobHandler := handler.NewJobHandler(schedulerSvc, backupSvc, jobRepo) // Membutuhkan BackupSvc untuk TriggerManual
	backupHandler := handler.NewBackupHandler(backupSvc)
	restoreHandler := handler.NewRestoreHandler(backupSvc)
	browserHandler := handler.NewBrowserHandler(browserSvc)
	setupHandler := handler.NewSetupHandler(authSvc)
	// remoteHandler := handler.NewRemoteHandler(remoteSvc)

	// --- 4. SEEDING ADMIN AWAL ---
	// if err := authSvc.RegisterAdmin(DefaultAdminUsername, DefaultAdminPassword); err != nil {
	// fmt.Printf("Admin Seeding Status: %v\n", err)
	// } else {
	// fmt.Println("✅ Admin user berhasil dibuat.")
	// }

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

	// ===========================================
	// RUTE PUBLIK UTAMA (TIDAK DILINDUNGI JWT)
	// ===========================================

	// Rute Publik (Login)
	e.POST("/api/v1/auth/login", authHandler.Login)

	// ⭐ KOREKSI DI SINI: GUNAKAN PATH LENGKAP /api/v1/setup AGAR TIDAK TERKENA JWT MIDDLEWARE
	setupGroup := e.Group("/api/v1/setup")
	setupGroup.GET("/status", setupHandler.GetSetupStatus)          // Cek apakah admin sudah ada
	setupGroup.POST("/register", setupHandler.RegisterInitialAdmin) // Registrasi admin pertama

	// ===========================================
	// RUTE PRIVAT (DILINDUNGI JWT)
	// ===========================================

	// Rute Privat (Dilindungi JWT)
	r := e.Group("/api/v1")
	r.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSecretKey),
	})) // Menerapkan Middleware JWT

	// Authentication Endpoints (Catatan: Ini sudah di bawah JWT, hanya bisa diakses setelah login)
	// Jika Anda ingin user bisa menggunakan setupHandler.Login (ya	ng biasanya dipakai untuk register)
	// pada rute ini setelah setup selesai, pastikan Anda menggunakan authHandler.Login di sini,
	// atau pindahkan semua rute Auth ke public/r.
	// Saat ini authHandler.Login sudah ada di atas (publik), jadi rute ini di bawah JWT tidak efisien.
	// authGroup := r.Group("/auth")
	// authGroup.POST("/login", setupHandler.Login) // Hapus atau ganti dengan r.POST("/auth/login", authHandler.Login) jika perlu

	// Rute Monitoring dan Logs
	r.GET("/monitoring/remotes", monitorHandler.GetRemoteStatusList)
	r.GET("/monitoring/logs", monitorHandler.GetJobLogs)
	r.GET("/monitoring/jobs", monitorHandler.GetScheduledJobs)

	// Rute Job Management (Create, List, Trigger)
	r.GET("/jobs/scheduled", jobHandler.GetScheduledJobs)    // List Job Monitoring
	r.GET("/jobs/script/:id", jobHandler.GetJobScript)       // Pratinjau Script
	r.POST("/jobs/trigger/:id", jobHandler.TriggerManualJob) // Tombol "Run Now"
	r.GET("/jobs/manual", jobHandler.GetManualJob)
	r.DELETE("/jobs/delete/:id", jobHandler.DeleteJob)
	r.PUT("/jobs/update/:id", jobHandler.UpdateJob)
	r.GET("/jobs/:id", jobHandler.GetJobByID)

	// Rute Aksi
	r.POST("/jobs/new", backupHandler.CreateNewJob)        // Create Backup (Manual/Auto)
	r.POST("/jobs/restore", restoreHandler.TriggerRestore) // Create Restore
	r.GET("/browser/files", browserHandler.ListFiles)
	r.GET("/browser/remotes", browserHandler.GetAvailableRemotes)
	r.GET("/browser/info", browserHandler.GetFileInfo)

	// Rute Konfigurasi
	// r.POST("/remotes/new", remoteHandler.AddNewRemote) // Add New Remote

	// --- 7. START DAEMONS (Goroutines) ---
	schedulerSvc.StartDaemon() // Memulai CRON Daemon
	monitorSvc.StartMonitoringDaemon()
	// (Tambahkan Goroutine untuk auto-update monitoring jika perlu)

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
