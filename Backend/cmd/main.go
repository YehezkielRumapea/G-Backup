package cmd

import (
	"fmt"
	"net/http"
	"os"

	"gbackup-system/backend/database"
	"gbackup-system/backend/internal/controllers" // Perbaikan: Gunakan 'handler'
	"gbackup-system/backend/internal/repository"
	"gbackup-system/backend/internal/services"

	// Frameworks dan Middleware
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	echomid "github.com/labstack/echo/v4/middleware"
)

const DefaultAdminUsername = "admin"
const DefaultAdminPassword = "admin123"

func main() {
	// Setup DB dan Connection
	dbInstance := database.Connect()

	jwtSecretKey := os.Getenv("JWT_SECRET")
	if jwtSecretKey == "" {
		fmt.Println("fatal: env key not found")
	}

	// Dependensi Injextion
	// 1. Inisiasliasi Repo (Data Access)
	userRepo := repository.NewUserRepository(dbInstance)
	jobRepo := repository.NewJobRepository(dbInstance)
	MonitorRepo := repository.NewMoniRepository(dbInstance)
	LogRepo := repository.NewLogRepository(dbInstance)

	// 2. Inisiaslisasi Servicews
	authSvc := services.NewAuthService(userRepo)
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

	// Start Server
	fmt.Println("backup sudah diinisialisasi. Menjalankan Echo server di port 8080")
	e.Logger.Fatal(e.Start(":8080"))
}
