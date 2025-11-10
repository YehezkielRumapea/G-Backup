package handler

import (
	"fmt"
	"net/http"
	"strconv" // Diperlukan untuk parsing ID dari URL

	"gbackup-new/backend/internal/service" // Sesuaikan path module

	"github.com/labstack/echo/v4"
)

// JobHandler struct menampung dependency (Service Layer)
type JobHandler struct {
	SchedulerSvc service.SchedulerService
	BackupSvc    service.BackupService // Dibutuhkan untuk memicu Job Manual
}

// NewJobHandler adalah constructor (Factory)
func NewJobHandler(schedulerSvc service.SchedulerService, backupSvc service.BackupService) *JobHandler {
	return &JobHandler{
		SchedulerSvc: schedulerSvc,
		BackupSvc:    backupSvc,
	}
}

// ----------------------------------------------------
// FUNGSI IMPLEMENTASI HANDLER
// ----------------------------------------------------

// GetScheduledJobs: Endpoint GET /api/v1/jobs/scheduled
// Mengambil data Job Monitoring yang sudah diformat dari Service Layer.
func (h *JobHandler) GetScheduledJobs(c echo.Context) error {

	// 1. Panggil method Service Layer yang sudah melakukan semua perhitungan
	// Method ini (GetScheduledJobsInfo) sudah menghitung NextRun dan memformat data.
	jobsDTO, err := h.SchedulerSvc.GetScheduledJobsInfo()

	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Gagal mengambil daftar Job terformat: " + err.Error()})
	}

	// 2. Kirim DTO (dari package service) yang sudah diformat langsung ke frontend
	return c.JSON(http.StatusOK, jobsDTO)
}

// GetJobScript: Endpoint GET /api/v1/jobs/script/:id (Fitur Pratinjau Script)
func (h *JobHandler) GetJobScript(c echo.Context) error {
	// 1. Ambil JobID dari URL parameter
	jobIDStr := c.Param("id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Job ID tidak valid"})
	}

	// 2. Panggil Service untuk men-generate script
	script, err := h.SchedulerSvc.GetGeneratedScript(uint(jobID))
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	}

	// 3. Kembalikan script sebagai JSON
	return c.JSON(http.StatusOK, map[string]string{"script_preview": script})
}

// TriggerManualJob: Endpoint POST /api/v1/jobs/trigger/:id
// Memicu Job yang sudah ada di DB (Tombol "Run Now")
func (h *JobHandler) TriggerManualJob(c echo.Context) error {
	// 1. Ambil JobID dari URL parameter
	jobIDStr := c.Param("id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Job ID tidak valid"})
	}

	// 2. Panggil Backup Service untuk memicu Job (Service akan mengambil dari DB)
	if err := h.BackupSvc.TriggerManualJob(uint(jobID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Gagal memicu Job: %v", err.Error())})
	}

	return c.JSON(http.StatusAccepted, map[string]string{"message": "Job berhasil dipicu."})
}

func (h *JobHandler) GetManualJob(c echo.Context) error {
	jobsDTO, err := h.SchedulerSvc.GetManualJob()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": "Gagal mengambil daftar Job Manual: " + err.Error(),
		})
	}
	return c.JSON(http.StatusOK, jobsDTO)
}

func (h *JobHandler) DeleteJob(c echo.Context) error {
	jobIDStr := c.Param("id")
	jobID, err := strconv.ParseUint(jobIDStr, 10, 64)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Job ID tidak valid"})
	}
	if err := h.BackupSvc.DeleteJob(uint(jobID)); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("Gagal menghapus Job: %v", err.Error())})
	}
	return c.JSON(http.StatusAccepted, map[string]string{"message": "Job berhasil dihapus."})
}
