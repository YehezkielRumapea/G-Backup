package handler

import (
	"gbackup-new/backend/internal/service" // Sesuaikan path module
	"net/http"

	"github.com/labstack/echo/v4"
)

type BrowserHandler struct {
	BrowserSvc service.BrowserService
}

func NewBrowserHandler(svc service.BrowserService) *BrowserHandler {
	return &BrowserHandler{BrowserSvc: svc}
}

// ListFiles: Endpoint GET /api/v1/browser/list
func (h *BrowserHandler) ListFiles(c echo.Context) error {
	remoteName := c.QueryParam("remote") // e.g., "gdrive_akun_utama"
	path := c.QueryParam("path")         // e.g., "/backups/db" (diawali / atau tidak)

	files, err := h.BrowserSvc.ListFiles(remoteName, path)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, files)
}
