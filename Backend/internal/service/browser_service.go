package service

import (
	"encoding/json"
	"fmt"
	"strings"
	// Sesuaikan path module
)

// RcloneFile: Struct untuk parsing output 'rclone lsjson'
type RcloneFile struct {
	Name    string `json:"Name"`
	IsDir   bool   `json:"IsDir"`
	Size    int64  `json:"Size"`
	ModTime string `json:"ModTime"` // Bisa time.Time jika parsing diperlukan
}

// BrowserService interface
type BrowserService interface {
	ListFiles(remoteName string, path string) ([]RcloneFile, error)
}

type browserServiceImpl struct {
	// (Tambahkan LogRepo jika perlu logging error)
}

func NewBrowserService() BrowserService {
	return &browserServiceImpl{}
}

// ListFiles: Menjalankan 'rclone lsjson'
func (s *browserServiceImpl) ListFiles(remoteName string, path string) ([]RcloneFile, error) {

	// Pastikan path selalu string kosong jika inputnya "/"
	CleanPath := strings.TrimPrefix(path, "/")

	// 1. Generate command
	fullPath := fmt.Sprintf("%s:%s", remoteName, CleanPath)

	// ✅ Tambahkan --max-depth 1 untuk navigasi folder interaktif
	args := []string{
		"rclone",
		"lsjson",
		fullPath,
		// "--no-mimetype",
		// "--max-depth", "1", // Limit ke folder saat ini
	}

	// 2. Eksekusi
	result := ExecuteCliJob(args)
	if !result.Success {
		// Gabungkan stderr/stdout dari Rclone untuk debugging yang lebih baik
		detailedError := fmt.Sprintf("Rclone failed for path '%s'. Detail: %s", fullPath, result.ErrorMsg)
		return nil, fmt.Errorf(detailedError)
	}

	// 3. Parse JSON
	var files []RcloneFile
	if err := json.Unmarshal([]byte(result.Output), &files); err != nil {
		return nil, fmt.Errorf("gagal parsing output lsjson: %w", err)
	}

	return files, nil
}
