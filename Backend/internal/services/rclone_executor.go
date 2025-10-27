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

	if len(commandArgs) == 0 {
		return RcloneResult{
			Success:  false,
			ErrorMsg: "Command arguments are empty"}
	}

	startTime := time.Now()
	cmdName := "rclone"
	args := append([]string{"-v"}, commandArgs...)
	if len(args) == 0 {
		return RcloneResult{Success: false, ErrorMsg: "Command arguments are empty"}
	}
	cmd := exec.Command(cmdName, args...)

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
