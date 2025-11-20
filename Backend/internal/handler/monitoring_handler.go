package handler

import (
	"net/http"
	"time"

	"gbackup-new/backend/internal/repository"
	"gbackup-new/backend/internal/service" // Sesuaikan path module

	"github.com/labstack/echo/v4"
)

// RemoteStatusResponse: Struct DTO untuk format output JSON ke frontend
type RemoteStatusResponse struct {
	Name           string  `json:"remote_name"`
	Status         string  `json:"status_connect"`
	StorageUsedGB  float64 `json:"used_storage_gb"`
	StorageTotalGB float64 `json:"total_storage_gb"`
	FreeStorageGB  float64 `json:"free_storage_gb"`
	LastChecked    string  `json:"last_checked_at"`
	ActiveJobCount int64   `json:"active_job_count"` // Placeholder (Akan diisi oleh Service)
}

// MonitoringHandler struct menampung dependency
type MonitoringHandler struct {
	MonitoringSvc service.MonitoringService
	SchedulerSvc  service.SchedulerService
	LogRepo       repository.LogRepository
}

// NewMonitoringHandler adalah constructor (Factory)
func NewMonitoringHandler(mSvc service.MonitoringService, sSvc service.SchedulerService, lRepo repository.LogRepository) *MonitoringHandler {
	return &MonitoringHandler{
		MonitoringSvc: mSvc,
		SchedulerSvc:  sSvc,
		LogRepo:       lRepo,
	}
}

// ----------------------------------------------------
// FUNGSI IMPLEMENTASI HANDLER
// ----------------------------------------------------

// GetRemoteStatusList: Endpoint GET /api/v1/monitoring/remotes
// Mengambil data status remote dari database untuk UI (Tampilan Remote Monitoring).
func (h *MonitoringHandler) GetRemoteStatusList(c echo.Context) error {

	// 1. Panggil Service Layer
	remotes, err := h.MonitoringSvc.GetRemoteStatusList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil data monitoring: " + err.Error(),
		})
	}

	// 2. Format Data (Mapping dari Model ke DTO)
	var responseData []RemoteStatusResponse
	for _, r := range remotes {
		// (Service Layer Anda harus diperbarui untuk menghitung JobRuns)
		responseData = append(responseData, RemoteStatusResponse{
			Name:           r.RemoteName,
			Status:         r.StatusConnect,
			StorageUsedGB:  r.UsedStorageGB,
			FreeStorageGB:  r.FreeStorageGB,
			StorageTotalGB: r.TotalStorageGB,
			LastChecked:    r.LastCheckedAt.Format(time.RFC3339),
			ActiveJobCount: r.ActiveJobCount, // Placeholder
		})
	}

	return c.JSON(http.StatusOK, responseData)
}

// GetJobLogs: Endpoint GET /api/v1/monitoring/logs
// Mengambil seluruh riwayat log dari database untuk Audit Trail.
func (h *MonitoringHandler) GetJobLogs(c echo.Context) error {

	// 1. Panggil Service Layer
	// (Asumsi GetJobLogs ada di MonitoringSvc dan memanggil LogRepo.FindAllLogs)
	logs, err := h.MonitoringSvc.GetJobLogs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil riwayat log: " + err.Error(),
		})
	}

	// 2. Kirim Logs
	return c.JSON(http.StatusOK, logs)
}

func (h *MonitoringHandler) GetScheduledJobs(c echo.Context) error {
	jobsDTO, err := h.SchedulerSvc.GetScheduledJobsInfo()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "gagal mengambil daftar jobs" + err.Error()})
	}
	return c.JSON(http.StatusOK, jobsDTO)
}
