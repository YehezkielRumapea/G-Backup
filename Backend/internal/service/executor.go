package service

import (
	"fmt"
	"os/exec" // Package inti untuk menjalankan command
	"strings"
	"time"
)

// RcloneResult menampung output dan status eksekusi CLI
// (Nama ini bisa diubah menjadi 'CliResult' jika Anda mau,
// karena ini menangani Rclone, Bash, dan Mysqldump)
type RcloneResult struct {
	Success          bool
	Output           string
	ErrorMsg         string
	Duration         time.Duration
	TransferredBytes int64
}

// ExecuteCliJob: Menjalankan command CLI APAPUN
// Menerima commandArgs sebagai slice (e.g., {"bash", "-c", "echo 'hello'"})
func ExecuteCliJob(commandArgs []string) RcloneResult {
	startTime := time.Now()

	// 1. Ambil command utama (e.g., "rclone" or "bash")
	cmdName := commandArgs[0]
	// 2. Ambil sisa argumen
	args := commandArgs[1:]

	// 3. Buat command
	cmd := exec.Command(cmdName, args...)

	// 4. Tangkap semua output (stdout dan stderr digabungkan)
	output, err := cmd.CombinedOutput()
	duration := time.Since(startTime)

	result := RcloneResult{
		Duration: duration,
		Output:   strings.TrimSpace(string(output)),
	}

	// 5. Error Handling Kritis (Cek Exit Code)
	if err != nil {
		// Jika command gagal (exit code != 0)
		result.Success = false
		result.ErrorMsg = fmt.Sprintf("Exit Error: %v. Output: %s", err, result.Output)
		return result
	}

	// Jika command sukses (exit code == 0)
	result.Success = true
	return result
}
