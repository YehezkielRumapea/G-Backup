package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"gbackup-system/backend/database"
	"gbackup-system/backend/internal/controllers" // Perbaikan: Gunakan 'handler'
	"gbackup-system/backend/internal/repository"
	"gbackup-system/backend/internal/services"
	"path/filepath"

	// Frameworks dan Middleware
	"github.com/joho/godotenv"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echomid "github.com/labstack/echo/v4/middleware"
)

const DefaultAdminUsername = "admin"
const DefaultAdminPassword = "admin123"

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

	loadDotEnv()

	jwtSecretKey := os.Getenv("JWT_SECRET_KEY")
	if jwtSecretKey == "" {
		log.Fatal("fatal: env key not found")
	}

	// Setup DB dan Connection
	dbInstance := database.Connect()

	// Dependensi Injextion
	// 1. Inisiasliasi Repo (Data Access)
	userRepo := repository.NewUserRepository(dbInstance)
	jobRepo := repository.NewJobRepository(dbInstance)
	MonitorRepo := repository.NewMoniRepository(dbInstance)
	LogRepo := repository.NewLogRepository(dbInstance)

	// 2. Inisiaslisasi Servicews
	authSvc := services.NewAuthService(userRepo)
	authSvc.SetSecretKey(jwtSecretKey)
	monitorSvc := services.NewMonitoringService(MonitorRepo, LogRepo)
	backupSvc := services.NewBackupService(jobRepo, LogRepo)
	schedulerSvc := services.NewSchedulerService(jobRepo, backupSvc)

	// 3. Injection service ke handler
	authHandler := controllers.NewAuthHandler(authSvc)
	monitorHandler := controllers.NewMonitoringHandler(monitorSvc)
	jobHandler := controllers.NewJobHandler(schedulerSvc, jobRepo)
	backupHandler := controllers.NewBackupHandler(backupSvc)   // Post/jobs/new
	restoreHandler := controllers.NewRestoreHandler(backupSvc) //Post/jobs/restore
	// Job Handler di inisialisasi disini

	// 4. Seeding Admin awal
	// Logic pendaftaran admin tunggal
	if err := authSvc.RegisterAdmin(DefaultAdminUsername, DefaultAdminPassword); err != nil {
		fmt.Printf("admin seeding status: %v\n", err)
	} else {
		fmt.Println("Admin user berhasil dibuat")
	}

	// SetUp Web server Echo
	e := echo.New()

	// 1. Middleware Dasar
	e.Use(echomid.Logger())
	e.Use(echomid.Recover())

	// MiddleWare CORS, Untuk koms ke FrontEnd
	e.Use(echomid.CORSWithConfig(echomid.CORSConfig{
		AllowOrigins: []string{"FrontEnd_URL"}, // Url FrontEnd
		AllowMethods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
	}))

	// Handler Health Check
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "G-Backup System is running")
	})

	// SetUp Routing
	// 1. Route Publik
	e.POST("/api/v1/auth/login", authHandler.Login)
	// 2. Route Privat
	r := e.Group("/api/v1")
	r.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(jwtSecretKey),
	}))

	// Route Monitoring and Logs
	r.GET("/monitoring/remotes", monitorHandler.GetRemoteStatusList)
	r.GET("/monitoring/logs", monitorHandler.GetJobLogs)

	// Job Management Route
	r.GET("/jobs/scheduled", jobHandler.GetScheduledJob)   // List Job
	r.POST("/jobs/new", backupHandler.CreateNewJob)        // Create/Dispatch Backup/Manual
	r.POST("/jobs/restore", restoreHandler.TriggerRestore) // Create/Dispatch Restore Job

	// Scheduler Daemon
	// Mulai Goutine untuk Cron-Job
	schedulerSvc.StartDaemon()

	go func() {
		time.Sleep(2 * time.Second)
		fmt.Println("Inisialisasi status remote monitoring...")

		// Ambil remote dan rclone.conf
		remoteNames, err := monitorSvc.GetRcloneConfiguredRemotes()
		if err != nil {
			fmt.Printf("Gagal mendapatkan remote dari rclone: %v\n", err)
			return
		}

		if len(remoteNames) == 0 {
			fmt.Println("Tidak ada remote yang dikonfigurasi di rclone.")
			return
		}

		fmt.Println("Status Remote yang terdeteksi dari rclone akan dimasukan ke DB:", remoteNames)
		for _, name := range remoteNames {
			monitorSvc.UpdateRemoteStatus(name)
		}
		fmt.Println("Inisialisasi status remote monitoring selesai.")
	}()

	// Start Server
	fmt.Println("backup sudah diinisialisasi. Menjalankan Echo server di port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
