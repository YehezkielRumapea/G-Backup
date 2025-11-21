package handler

import (
	"net/http"
	"time"

	"gbackup-new/backend/internal/repository"
	"gbackup-new/backend/internal/service"

	"github.com/labstack/echo/v4"
)

type MonitoringHandler struct {
	MonitoringSvc service.MonitoringService
	SchedulerSvc  service.SchedulerService
	LogRepo       repository.LogRepository
}

func NewMonitoringHandler(mSvc service.MonitoringService, sSvc service.SchedulerService, lRepo repository.LogRepository) *MonitoringHandler {
	return &MonitoringHandler{
		MonitoringSvc: mSvc,
		SchedulerSvc:  sSvc,
		LogRepo:       lRepo,
	}
}

func (h *MonitoringHandler) SyncRemotes(c echo.Context) error {
	err := h.MonitoringSvc.SyncRemotesWithRclone()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal sync remote: " + err.Error(),
		})
	}

	remotes, err := h.MonitoringSvc.GetRemoteStatusList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil remote setelah sync: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "Remote berhasil disinkronkan",
		"count":   len(remotes),
		"remotes": remotes,
	})
}

func (h *MonitoringHandler) GetRemoteStatusList(c echo.Context) error {
	remotes, err := h.MonitoringSvc.GetRemoteStatusList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil data monitoring: " + err.Error(),
		})
	}

	var responseData []map[string]interface{}
	for _, r := range remotes {
		responseData = append(responseData, map[string]interface{}{
			"remote_name":      r.RemoteName,
			"status_connect":   r.StatusConnect,
			"used_storage_gb":  r.UsedStorageGB,
			"free_storage_gb":  r.FreeStorageGB,
			"total_storage_gb": r.TotalStorageGB,
			"last_checked_at":  r.LastCheckedAt.Format(time.RFC3339),
			"active_job_count": r.ActiveJobCount,
			"system_message":   r.SystemMessage,
		})
	}

	return c.JSON(http.StatusOK, responseData)
}

func (h *MonitoringHandler) GetJobLogs(c echo.Context) error {
	logs, err := h.MonitoringSvc.GetJobLogs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil riwayat log: " + err.Error(),
		})
	}

	return c.JSON(http.StatusOK, logs)
}

func (h *MonitoringHandler) GetScheduledJobs(c echo.Context) error {
	jobsDTO, err := h.SchedulerSvc.GetScheduledJobsInfo()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "gagal mengambil daftar jobs: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, jobsDTO)
}
