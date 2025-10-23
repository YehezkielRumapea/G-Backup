package controllers

import (
	"net/http"
	"time"

	"gbackup-system/backend/internal/repository"
	"gbackup-system/backend/internal/services"

	"github.com/labstack/echo/v4"
)

// JobResponse == Struct unruk response ke frontend

type JobResponse struct {
	ID          uint       `json:"id"`
	Name        string     `json:"name"`
	RemoteName  string     `json:"remoteName"`
	SourceType  string     `json:"SourceType"`
	RcloneMode  string     `json:"rcloneMode"`
	LastRun     *time.Time `json:"lastRun"`
	StatusQueue string     `json:"statusQueue"`
	NextRun     string     `json:"nextRun"`
}

type JobHandler struct {
	Schedulersvc services.SchedulerService // Next Run Logic
	JobRepo      repository.JobRepository  // Untuk FindActiveJobs
}

func NewJobHandler(Schedsvc services.SchedulerService, JobRepo repository.JobRepository) *JobHandler {
	return &JobHandler{
		Schedulersvc: Schedsvc,
		JobRepo:      JobRepo,
	}
}

func (h *JobHandler) GetScheduledJob(c echo.Context) error {
	// ambil semua job yang active dari repository
	jobs, err := h.JobRepo.FindActiveJobs()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil Job yang aktif: " + err.Error(),
		})
	}

	var responseData []JobResponse

	for _, job := range jobs {
		LastRun := time.Time{}
		if job.LastRun != nil {
			LastRun = *job.LastRun
		}

		// Logic Next Run
		nextRun := h.Schedulersvc.CalculateNextRun(job.ScheduleCron, LastRun)

		responseData = append(responseData, JobResponse{
			ID:          job.ID,
			Name:        job.Name,
			RemoteName:  job.RemoteName,
			SourceType:  job.SourceType,
			RcloneMode:  job.RcloneMode,
			LastRun:     &LastRun,
			StatusQueue: job.StatusQueue,
			NextRun:     nextRun.Format(time.RFC3339), // FormatWaktu untuk FrontEnd\
		})
	}
	return c.JSON(http.StatusOK, responseData)
}
