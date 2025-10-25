package controllers

import (
	"net/http"

	"gbackup-system/backend/internal/repository"
	"gbackup-system/backend/internal/services"

	"github.com/labstack/echo/v4"
)

type JobHandler struct {
	Schedulersvc services.SchedulerService // Next Run Logic
	JobRepo      repository.JobRepository
}

func NewJobHandler(Schedsvc services.SchedulerService, JobRepo repository.JobRepository) *JobHandler {
	return &JobHandler{
		Schedulersvc: Schedsvc,
		JobRepo:      JobRepo,
	}
}

func (h *JobHandler) GetScheduledJob(c echo.Context) error {
	jobsDto, err := h.Schedulersvc.GetScheduledJobsInfo()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "gagal mengambil data job: " + err.Error()})
	}

	return c.JSON(http.StatusOK, jobsDto)
}
