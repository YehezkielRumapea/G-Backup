package services

import (
	"fmt"
	"os/exec"
	"strings"
	"time"
)

type RcloneResult struct {
	Success  bool
	Output   string
	ErrorMsg string
	Duration time.Duration
}

func ExecuteRcloneJob(commandArgs []string) RcloneResult {
	startTime := time.Now()

	cmd := exec.Command("rclone", commandArgs...)

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)

	result := RcloneResult{
		Duration: duration,
		Output:   strings.TrimSpace(string(output)),
	}

	if err != nil {
		result.Success = false

		result.ErrorMsg = fmt.Sprintf("Exit Error: %v. Output: %s", err, result.Output)
		return result
	}

	result.Success = true
	return result
}
