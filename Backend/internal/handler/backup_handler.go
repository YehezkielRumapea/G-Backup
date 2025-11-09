package handler

import (
	"fmt"
	"net/http"

	"gbackup-new/backend/internal/models" // Sesuaikan path module
	"gbackup-new/backend/internal/service"

	"github.com/labstack/echo/v4"
)

// BackupRequestDTO (Data Transfer Object) untuk menerima input Job Baru
// Ini harus sesuai dengan skema "Script Runner" baru.
type BackupRequestDTO struct {
	JobName         string `json:"job_name" validate:"required"`
	SourcePath      string `json:"source_path" validate:"required"` // Path yang akan digunakan Rclone
	RemoteName      string `json:"remote_name" validate:"required"`
	DestinationPath string `json:"destination_path" validate:"required"`
	ScheduleCron    string `json:"schedule_cron"` // Kosong jika Manual

	// Script Mentah dari User
	PreScript  string `json:"pre_script"`
	PostScript string `json:"post_script"`
}

type BackupHandler struct {
	BackupSvc service.BackupService
}

// NewBackupHandler adalah constructor (Factory)
func NewBackupHandler(svc service.BackupService) *BackupHandler {
	return &BackupHandler{BackupSvc: svc}
}

// ----------------------------------------------------
// FUNGSI IMPLEMENTASI HANDLER
// ----------------------------------------------------

// CreateNewJob: Endpoint POST /api/v1/jobs/new
func (h *BackupHandler) CreateNewJob(c echo.Context) error {
	var req BackupRequestDTO

	// 1. Binding dan Validasi Input
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format input JSON tidak valid."})
	}

	// 2. Dapatkan User ID (dari Middleware JWT)
	// (Di implementasi nyata, ambil dari context: userID := c.Get("userID").(uint))
	userID := uint(1) // Placeholder untuk admin tunggal

	rcloneMode := "COPY"
	// 3. Buat Struct Job (Model) dari DTO
	// (Kita set OperationMode default ke BACKUP)
	newJob := models.ScheduledJob{
		UserID:          userID,
		JobName:         req.JobName,
		RcloneMode:      rcloneMode,
		SourcePath:      req.SourcePath,
		RemoteName:      req.RemoteName,
		DestinationPath: req.DestinationPath,
		PreScript:       req.PreScript,
		PostScript:      req.PostScript,
		ScheduleCron:    req.ScheduleCron, // Jika "" -> Job Manual (Template)
		IsActive:        true,
		StatusQueue:     "PENDING",
	}

	// 4. Panggil Service untuk Dispatch Job
	// Service akan: (A) Menyimpan ke DB, (B) Meluncurkan Goroutine (jika Manual)
	if err := h.BackupSvc.CreateJobAndDispatch(&newJob); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Gagal mendispatch Job: %v", err.Error())})
	}

	// 5. Respon Sukses (202 Accepted karena Job mungkin diproses di background)
	return c.JSON(http.StatusAccepted, map[string]string{"message": "Job berhasil diterima dan dikirim untuk eksekusi."})
}
