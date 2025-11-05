package handler

import (
	"fmt"
	"gbackup-new/backend/internal/models"
	"gbackup-new/backend/internal/service"
	"net/http"
	"time" // Diperlukan untuk auto-generate JobName

	"github.com/labstack/echo/v4"
)

// RestoreRequestDTO (SANGAT DISEDERHANAKAN)
// Hanya menerima input minimal yang diperlukan untuk Restore.
type RestoreRequestDTO struct {
	SourcePath      string `json:"source_path" validate:"required"` // Path file di Cloud
	OperationMode   string `json:"operation_mode" validate:"required"`
	RemoteName      string `json:"remote_name" validate:"required"`
	DestinationPath string `json:"destination_path" validate:"required"` // Path Lokal Target
}

type RestoreHandler struct {
	BackupSvc service.BackupService
}

func NewRestoreHandler(svc service.BackupService) *RestoreHandler {
	return &RestoreHandler{BackupSvc: svc}
}

// TriggerRestore: Endpoint POST /api/v1/jobs/restore
func (h *RestoreHandler) TriggerRestore(c echo.Context) error {
	var req RestoreRequestDTO

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format input JSON tidak valid."})
	}

	userID := uint(1) // Placeholder Admin

	// 2. Buat Struct Job (Model) secara otomatis
	// (Backend men-generate JobName dan menyetel semua default Restore)
	newJob := models.ScheduledJob{
		UserID: userID,
		// Auto-generate nama Job (karena user tidak input)
		JobName: fmt.Sprintf("Manual Restore - %d", time.Now().Unix()),

		OperationMode:   "RESTORE", // Flag Kritis
		RcloneMode:      "copy",    // Restore selalu "copy"
		SourcePath:      req.SourcePath,
		RemoteName:      req.RemoteName,
		DestinationPath: req.DestinationPath,

		ScheduleCron: "", // Job Manual
		PreScript:    "", // Tidak ada Pre-Script
		PostScript:   "", // Tidak ada Post-Script

		IsActive:    true,
		StatusQueue: "PENDING",
	}

	// 3. Panggil Service untuk Dispatch Job
	if err := h.BackupSvc.CreateJobAndDispatch(&newJob); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Gagal mendispatch Restore Job: %v", err.Error())})
	}

	return c.JSON(http.StatusAccepted, map[string]string{"message": "Restore Job diterima dan dimulai."})
}
