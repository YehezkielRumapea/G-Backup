package controllers

import (
	"fmt"
	"gbackup-system/backend/internal/models"
	"gbackup-system/backend/internal/services"
	"net/http"

	"github.com/labstack/echo/v4"
)

// struct import job manual
type RestoreRequest struct {
	JobName     string `json:"job_name" validate:"required"`
	RemoteName  string `json:"remote_name" validate:"required"`
	Sourcepath  string `json:"source_path" validate:"required"`
	RestorePath string `json:"restore_path" validate:"required"`
	SourceType  string `json:"job_type" validate:"required"`

	// Credential
	Dbuser string `json:"db_user"`
	Dbpass string `json:"db_pass"`
}

type RestoreHandler struct {
	Backupsvc services.BackupService
}

func NewRestoreHandler(svc services.BackupService) *RestoreHandler {
	return &RestoreHandler{Backupsvc: svc}
}

// TriggerRestore && Endpoint API/v1/restore
func (h *RestoreHandler) TriggerRestore(c echo.Context) error {
	var req RestoreRequest

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Format Json tidak valid"})
	}

	// Ambil User ID
	userID := uint(1)

	// Struct Job dengan model Restore
	newjob := models.ScheduledJob{
		UserID:     userID,
		Name:       req.JobName,
		SourceType: req.SourceType,
		RcloneMode: "RESTORE",

		// SourcePath
		SourcePath: req.Sourcepath,
		// DestinationPath
		DestinationPath: req.RestorePath,

		RemoteName:   req.RemoteName,
		ScheduleCron: "", // Manual Job
		DbUser:       req.Dbuser,
		DbPass:       req.Dbpass,
		IsEncrypted:  false,
		StatusQueue:  "PENDING",
	}
	if err := h.Backupsvc.CreateJobAndDispatch(&newjob); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Gagal me-restore job: %v", err.Error())})
	}
	return c.JSON(http.StatusAccepted, map[string]string{"message": "Job restore berhasil di-trigger"})
}
