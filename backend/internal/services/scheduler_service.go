package services

import (
	"fmt"
	"gbackup-system/backend/internal/repository"
	"time"

	"github.com/robfig/cron/v3"
)

type SchedulerService interface {
	StartDaemon()
	RunScheduledJob() error
	CalculateNextRun(schedule string, lastRun time.Time) time.Time
}

type SchedulerServiceImpl struct {
	JobRepo repository.JobRepository
}

func NewSchedulerService(JRepo repository.JobRepository) SchedulerService {
	return &SchedulerServiceImpl{JobRepo: JRepo}
}

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

func (s *SchedulerServiceImpl) RunScheduledJob() error {
	jobs, err := s.JobRepo.FindActiveJobs()
	if err != nil {
		return fmt.Errorf("Gagal mengambil job aktif dari DB: %w", err)
	}
	now := time.Now()

	for _, job := range jobs {
		lastRunTime := time.Time{}
		if job.LastRun != nil {
			lastRunTime = *job.LastRun
		}

		nextRun := s.CalculateNextRun(job.Schedule, lastRunTime)

		if nextRun.Before(now) {
		}
	}

	fmt.Println("Cek Job Terjadawal")
	return nil
}

func (s *SchedulerServiceImpl) StartDaemon() {
	go func() {
		fmt.Println("Schedular Daemon Aktif, Pengecekan Job tiap 5 menit ")

		for {
			err := s.RunScheduledJob()
			if err != nil {
				fmt.Printf("Daemon Error: %v\n", err)
			}
			time.Sleep(5 * time.Minute)
		}
	}()
}
