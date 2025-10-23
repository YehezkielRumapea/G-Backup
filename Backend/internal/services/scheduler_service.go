package services

import (
	"fmt"
	"gbackup-system/backend/internal/models"
	"gbackup-system/backend/internal/repository"
	"time"

	"github.com/robfig/cron/v3"
)

type BackupService interface {
	StartNewJob(job models.ScheduledJob)
	CreateJobAndDispatch(job *models.ScheduledJob) error
}
type SchedulerService interface {
	StartDaemon()
	RunScheduledJob() error
	CalculateNextRun(schedule string, lastRun time.Time) time.Time
}

type SchedulerServiceImpl struct {
	JobRepo     repository.JobRepository
	BackupSync  BackupService
	intervalCek time.Duration
}

func NewSchedulerService(JRepo repository.JobRepository, Bsyc BackupService) SchedulerService {
	return &SchedulerServiceImpl{JobRepo: JRepo, BackupSync: Bsyc}
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

			s.BackupSync.StartNewJob(job)
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
