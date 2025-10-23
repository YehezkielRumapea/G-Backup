package controllers

import (
	"fmt"
	"gbackup-system/backend/internal/models"
	"gbackup-system/backend/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// BackupRequest DTO: Struct untuk menerima input konfigurasi Job Baru
type BackupRequest struct {
	JobName         string `json:"job_name" validate:"required"`
	SourcePath      string `json:"source_path" validate:"required"`
	RemoteName      string `json:"remote_name" validate:"required"`
	SourceType      string `json:"job_type" validate:"required"`         // FILE atau DB
	DestinationPath string `json:"destination_path" validate:"required"` // Path di cloud
	ScheduleCron    string `json:"schedule_cron"`                        // CRON string (kosong jika Manual)
	IsEncrypted     bool   `json:"is_encrypted"`

	// Kredensial DB yang sensitif (Disimpan terenkripsi di DB jika Job Auto)
	DbUser string `json:"db_user"`
	DbPass string `json:"db_pass"`
}

type BackupHandler struct {
	BackupSvc services.BackupService
}

func NewBackupHandler(svc services.BackupService) *BackupHandler {
	return &BackupHandler{BackupSvc: svc}
}

// CreateNewJob: Endpoint POST /api/v1/jobs/new
func (h *BackupHandler) CreateNewJob(c echo.Context) error {
	var req BackupRequest

	// 1. Binding dan Validasi Input
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format input JSON tidak valid."})
	}

	// 2. Dapatkan User ID (Asumsi Middleware JWT menyuntikkan UserID)
	// Logika: UserID diekstrak dari JWT payload
	userID := uint(1) // Placeholder untuk user admin tunggal

	// 3. Tentukan Mode Operasi
	opMode := "BACKUP"
	if req.ScheduleCron == "" {
		opMode = "MANUAL_BACKUP" // Job Manual/Sekali Jalan
	}

	// 4. Buat Struct Job (Model) yang akan dikirim ke Service Layer
	newJob := models.ScheduledJob{
		UserID:          userID,
		Name:            req.JobName,
		SourceType:      req.SourceType,
		OperationMode:   opMode,
		SourcePath:      req.SourcePath,
		RemoteName:      req.RemoteName,
		DestinationPath: req.DestinationPath,
		ScheduleCron:    req.ScheduleCron,
		IsEncrypted:     req.IsEncrypted,

		// KRITIS: DbUser dan DbPass disimpan di struct sebelum dikirim ke Service
		DbUser: req.DbUser,
		DbPass: req.DbPass, // Akan dienkripsi oleh Service sebelum disimpan di DB

		StatusQueue: "PENDING",
	}

	// 5. Panggil Service untuk Dispatch Job
	// Service akan menangani: (A) Enkripsi/Penyimpanan di DB, dan (B) Peluncuran Goroutine.
	if err := h.BackupSvc.CreateJobAndDispatch(&newJob); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Gagal mendispatch Job: %v", err.Error())})
	}

	// 6. Respon Sukses
	return c.JSON(http.StatusAccepted, map[string]string{"message": "Job berhasil diterima dan dikirim untuk eksekusi."})
}
