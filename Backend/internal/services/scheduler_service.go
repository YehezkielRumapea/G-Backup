package services

import (
	"fmt"
	"gbackup-system/backend/internal/models"
	"gbackup-system/backend/internal/repository"
	"time"

	"github.com/robfig/cron/v3"
)

type JobMonitoring struct {
	ID           uint   `json:"id"`
	JobName      string `json:"job_name"`
	Type         string `json:"type"`
	GdriveTarget string `json:"gdrive_target"`
	Mode         string `json:"mode"`
	LastRun      string `json:"last_run"`
	Status       string `json:"status"`
	NextRun      string `json:"next_run"`
}

type BackupService interface {
	StartNewJob(job models.ScheduledJob)
	CreateJobAndDispatch(job *models.ScheduledJob) error
}
type SchedulerService interface {
	StartDaemon()
	RunScheduledJob() error
	CalculateNextRun(schedule string, lastRun time.Time) time.Time
	GetScheduledJobsInfo() ([]JobMonitoring, error)
}

type SchedulerServiceImpl struct {
	JobRepo     repository.JobRepository
	BackupSvc   BackupService
	intervalCek time.Duration
}

func NewSchedulerService(JRepo repository.JobRepository, Bsvc BackupService) SchedulerService {
	return &SchedulerServiceImpl{
		JobRepo:     JRepo,
		BackupSvc:   Bsvc,
		intervalCek: 5 * time.Minute,
	}
}

// Main Logic
// Menghitung waktu run job selanjutnya
func (s *SchedulerServiceImpl) CalculateNextRun(schedule string, lastRun time.Time) time.Time {
	parser := cron.NewParser(
		cron.Minute | cron.Hour | cron.Dom | cron.Month | cron.Dow,
	)

	sched, err := parser.Parse(schedule)

	if err != nil {
		fmt.Printf("Error, gagal parsing cron expression: %v\n", err)
		return time.Time{}
	}

	baseTime := lastRun
	if lastRun.IsZero() {
		baseTime = time.Now()
	}
	return sched.Next(baseTime)
}

// Me-Run Job yang terjadwal
func (s *SchedulerServiceImpl) RunScheduledJob() error {
	jobs, err := s.JobRepo.FindActiveJobs()
	if err != nil {
		return fmt.Errorf("gagal mengambil job aktif dari DB: %w", err)
	}
	now := time.Now()

	for _, job := range jobs {
		lastRunTime := time.Time{}
		if job.LastRun != nil {
			lastRunTime = *job.LastRun
		}

		nextRun := s.CalculateNextRun(job.ScheduleCron, lastRunTime)

		if nextRun.Before(now) {

			lockTime := time.Now() // Job locking

			err := s.JobRepo.UpdateLastRunStatus(job.ID, lockTime, "RUNNING")

			if err != nil {
				fmt.Printf("Job %d (%s) Sudah berjalan : %v\n",
					job.ID,
					job.Name,
					err)
				continue
			}

			fmt.Printf("[Scheduler] Dispatching Job %d (%s) - Next Run: %s\n",
				job.ID,
				job.Name,
				nextRun.Format(time.RFC3339))

			s.BackupSvc.StartNewJob(job)
		}
	}
	return nil
}

// Goroutine
func (s *SchedulerServiceImpl) StartDaemon() {
	if s.intervalCek == 0 {
		s.intervalCek = 5 * time.Minute
	}
	go func() {
		fmt.Println("Schedular Daemon Aktif, Pengecekan Job tiap 5 menit ")

		for {
			err := s.RunScheduledJob()
			if err != nil {
				fmt.Printf("Daemon Error: %v\n", err)
			}
			time.Sleep(s.intervalCek)
		}
	}()
}

func (s *SchedulerServiceImpl) GetScheduledJobsInfo() ([]JobMonitoring, error) {
	// Ambil semua Job yang aktif
	jobs, err := s.JobRepo.FindActiveJobs()
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil job aktif dari DB: %w", err)
	}

	var Output []JobMonitoring

	//  Loop dan hitung format data
	for _, job := range jobs {
		lastRunTime := time.Time{}
		lastRunStr := ""

		if job.LastRun != nil {
			lastRunTime = *job.LastRun
			lastRunStr = job.LastRun.Format("02-01-2006 15:04")
		}

		//Next Run
		nextRunTime := s.CalculateNextRun(job.ScheduleCron, lastRunTime)

		// Pembuatan Format Outptu
		mode := "Auto"
		JobTypeFormatted := fmt.Sprintf("%s: %s", job.SourceType, job.SourcePath)

		// format Output
		Output = append(Output, JobMonitoring{
			ID:           job.ID,
			JobName:      job.Name,
			Type:         JobTypeFormatted,
			GdriveTarget: job.RemoteName,
			Mode:         mode,
			LastRun:      lastRunStr,
			Status:       job.StatusQueue,
			NextRun:      nextRunTime.Format("02-01-2006 15:04"),
		})
	}
	return Output, nil
}
