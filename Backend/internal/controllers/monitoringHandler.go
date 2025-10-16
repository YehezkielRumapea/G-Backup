package controllers

import (
	"gbackup-system/backend/internal/services"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

type RemoteStatusResponse struct {
	Name           string  `json:"RemoteName"`
	Status         string  `json:"statusConnect"`
	StorageUsedGb  float64 `jsonn:"usedStorageGB"`
	TotalStorageGB float64 `jsonn:"totalStorageGB"`
	LastChecked    string  `json:"lasrCheckedAt"`
	JobRuns        int     `json:"jobRuns"`
}

type MonitoringHandler struct {
	MonitoringSvc services.MonitoringService
}

func NewMonitoringHandler(Svc services.MonitoringService) *MonitoringHandler {
	return &MonitoringHandler{MonitoringSvc: Svc}
}

func (h *MonitoringHandler) GetRemoteStatusList(c echo.Context) error {
	// Memanggil service layer untuk mengamboil data monit
	remotes, err := h.MonitoringSvc.GetRemoteStatusList()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "gagal mengambil data monit:" + err.Error(),
		})
	}

	var responseData []RemoteStatusResponse

	for _, r := range remotes {
		responseData = append(responseData, RemoteStatusResponse{
			Name:           r.RemoteName,
			Status:         r.StatusConnect,
			StorageUsedGb:  r.UsedStorageGB,
			TotalStorageGB: r.TotalStorageGB,
			LastChecked:    r.LastCheckedAt.Format(time.Kitchen),
			JobRuns:        0, // Placeholder for job runs,
		})
	}

	// JSON Response
	return c.JSON(http.StatusOK, responseData)
}

func (h *MonitoringHandler) GetJobLogs(c echo.Context)
