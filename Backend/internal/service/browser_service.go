package service

import (
	"encoding/json"
	"fmt"
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
	// 1. Generate command
	fullPath := fmt.Sprintf("%s:%s", remoteName, path)
	// rclone lsjson [remote:path]
	args := []string{"rclone", "lsjson", fullPath, "--no-mimetype", "--no-modtime-accuracy"}

	// 2. Eksekusi (menggunakan Executor Anda)
	result := ExecuteCliJob(args)
	if !result.Success {
		return nil, fmt.Errorf("gagal list files: %s", result.ErrorMsg)
	}

	// 3. Parse JSON
	var files []RcloneFile
	if err := json.Unmarshal([]byte(result.Output), &files); err != nil {
		return nil, fmt.Errorf("gagal parsing output lsjson: %w", err)
	}

	return files, nil
}
