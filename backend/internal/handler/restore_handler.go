package handler

import (
	"fmt"
	"net/http"

	"gbackup-new/backend/internal/models" // Sesuaikan path module
	"gbackup-new/backend/internal/service"

	"github.com/labstack/echo/v4"
)

// RestoreRequestDTO (Data Transfer Object) untuk menerima input Restore Job Manual
type RestoreRequestDTO struct {
	JobName         string `json:"job_name" validate:"required"`
	SourcePath      string `json:"source_path" validate:"required"` // Path file di Cloud
	RemoteName      string `json:"remote_name" validate:"required"`
	DestinationPath string `json:"destination_path" validate:"required"` // Path Lokal Target
	RcloneMode      string `json:"rclone_mode"`                          // "COPY" (Default)

	// Script Mentah (Opsional untuk Restore)
	PreScript  string `json:"pre_script"`
	PostScript string `json:"post_script"`
}

type RestoreHandler struct {
	BackupSvc service.BackupService
}

// NewRestoreHandler adalah constructor (Factory)
func NewRestoreHandler(svc service.BackupService) *RestoreHandler {
	return &RestoreHandler{BackupSvc: svc}
}

// ----------------------------------------------------
// FUNGSI IMPLEMENTASI HANDLER
// ----------------------------------------------------

// TriggerRestore: Endpoint POST /api/v1/jobs/restore
func (h *RestoreHandler) TriggerRestore(c echo.Context) error {
	var req RestoreRequestDTO

	// 1. Binding dan Validasi Input
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format input JSON tidak valid."})
	}

	// 2. Dapatkan User ID (dari Middleware JWT)
	userID := uint(1) // Placeholder untuk admin tunggal

	// 3. Buat Struct Job (Model) dengan Mode RESTORE (Job Manual)
	newJob := models.ScheduledJob{
		UserID:  userID,
		JobName: req.JobName,

		// FLAG KRITIS: Memberitahu Service untuk membalik flow
		OperationMode: "RESTORE",

		RcloneMode:      "COPY",         // Restore selalu "copy"
		SourcePath:      req.SourcePath, // Path di Cloud
		RemoteName:      req.RemoteName,
		DestinationPath: req.DestinationPath, // Path di Lokal

		// Script (jika user ingin melakukan sesuatu setelah restore)
		PreScript:  req.PreScript,
		PostScript: req.PostScript,

		ScheduleCron: "", // Job Manual (NULL)
		IsActive:     true,
		StatusQueue:  "PENDING",
	}

	// 4. Panggil Service untuk Dispatch Job
	// Service akan: (A) Menyimpan ke DB, (B) Meluncurkan Goroutine
	if err := h.BackupSvc.CreateJobAndDispatch(&newJob); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Gagal mendispatch Restore Job: %v", err.Error())})
	}

	// 5. Respon Sukses (202 Accepted)
	return c.JSON(http.StatusAccepted, map[string]string{"message": "Restore Job diterima dan dimulai."})
}
