package service

import (
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type RcloneResult struct {
	Success          bool
	Output           string
	ErrorMsg         string
	Duration         time.Duration
	TransferredBytes int64
}

func ExecuteCliJob(commandArgs []string) RcloneResult {
	// 1. Safety Check: Pastikan command tidak kosong agar tidak panic
	if len(commandArgs) == 0 {
		return RcloneResult{
			Success:  false,
			ErrorMsg: "Command arguments cannot be empty",
		}
	}

	startTime := time.Now()

	cmdName := commandArgs[0]
	args := commandArgs[1:]

	cmd := exec.Command(cmdName, args...)

	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)

	outputStr := strings.TrimSpace(string(output))

	result := RcloneResult{
		Duration: duration,
		Output:   outputStr,
	}

	// Parse bytes baik sukses maupun gagal (kadang rclone error tapi sempat transfer data)
	result.TransferredBytes = parseTransferredBytes(outputStr)

	if err != nil {
		result.Success = false
		result.ErrorMsg = fmt.Sprintf("Exit Error: %v. Output: %s", err, result.Output)
		return result
	}

	result.Success = true
	return result
}

func parseTransferredBytes(output string) int64 {
	// Regex untuk menangkap angka dan unit
	// Menangkap: "Transferred:   12.500 MiB"
	re := regexp.MustCompile(`Transferred:\s+([\d\.]+)\s*(\w+i?B)`)

	// 1. PERBAIKAN KRITIS: Gunakan FindAll untuk mengambil SEMUA kemunculan
	matches := re.FindAllStringSubmatch(output, -1)

	if len(matches) == 0 {
		return 0
	}

	// 2. Ambil match TERAKHIR (status final)
	lastMatch := matches[len(matches)-1]

	valueStr := lastMatch[1]
	unit := strings.ToUpper(lastMatch[2])

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil {
		return 0
	}

	// 3. Konversi Unit
	switch unit {
	case "B":
		return int64(value)
	case "KB", "KIB":
		return int64(value * 1024)
	case "MB", "MIB":
		return int64(value * 1024 * 1024)
	case "GB", "GIB":
		return int64(value * 1024 * 1024 * 1024)
	case "TB", "TIB":
		return int64(value * 1024 * 1024 * 1024 * 1024)
	default:
		return int64(value)
	}
}
