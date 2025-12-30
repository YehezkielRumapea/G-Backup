// internal/handler/remote_handler.go
package handler

import (
	"gbackup-new/backend/internal/service"
	"net/http"

	"github.com/labstack/echo/v4"
)

type RemoteHandler struct {
	remoteService service.AddRemoteService
}

func NewRemoteHandler(remoteService service.AddRemoteService) *RemoteHandler {
	return &RemoteHandler{
		remoteService: remoteService,
	}
}

// InitAuth handles POST /api/remote/init-auth
func (h *RemoteHandler) InitAuth(c echo.Context) error {
	var req struct {
		Name string `json:"name" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if req.Name == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Remote name is required",
		})
	}

	authURL, state, err := h.remoteService.InitAuth(req.Name)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"auth_url": authURL,
		"state":    state,
	})
}

// OAuthCallback handles GET /api/remote/oauth-callback
func (h *RemoteHandler) OAuthCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	errorParam := c.QueryParam("error")

	// Handle user cancellation
	if errorParam != "" {
		html := `<!DOCTYPE html>
<html>
<head><title>Authentication Cancelled</title></head>
<body style="font-family: Arial, sans-serif; text-align: center; padding: 50px;">
    <h2>❌ Authentication Cancelled</h2>
    <p>You can close this window.</p>
    <script>
        if (window.opener) {
            window.opener.postMessage({type: 'oauth-error', error: 'cancelled'}, '*');
        }
        setTimeout(() => window.close(), 2000);
    </script>
</body>
</html>`
		return c.HTML(http.StatusOK, html)
	}

	if code == "" || state == "" {
		return c.HTML(http.StatusBadRequest, "<h2>Error: Missing parameters</h2>")
	}

	html := `<!DOCTYPE html>
<html>
<head>
    <title>Authentication Success</title>
    <style>
        body { 
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', sans-serif;
            display: flex; align-items: center; justify-content: center;
            height: 100vh; margin: 0;
            background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }
        .container {
            background: white; padding: 40px; border-radius: 12px;
            box-shadow: 0 10px 40px rgba(0,0,0,0.2); text-align: center;
            max-width: 400px;
        }
        h2 { color: #4CAF50; margin: 0 0 10px 0; font-size: 24px; }
        p { color: #666; margin: 0; font-size: 14px; }
        .spinner {
            margin: 20px auto 0; width: 40px; height: 40px;
            border: 4px solid #f3f3f3; border-top: 4px solid #4CAF50;
            border-radius: 50%; animation: spin 1s linear infinite;
        }
        @keyframes spin {
            0% { transform: rotate(0deg); }
            100% { transform: rotate(360deg); }
        }
    </style>
</head>
<body>
    <div class="container">
        <h2>✓ Autentikasi Berhasil!</h2>
        <p>Sedang menyimpan konfigurasi...</p>
        <div class="spinner"></div>
    </div>
    <script>
        if (window.opener) {
            window.opener.postMessage({
                type: 'oauth-callback',
                code: '` + code + `',
                state: '` + state + `'
            }, '*');
        }
        setTimeout(() => window.close(), 3000);
    </script>
</body>
</html>`

	return c.HTML(http.StatusOK, html)
}

// FinalizeConfig handles POST /api/remote/finalize
func (h *RemoteHandler) FinalizeConfig(c echo.Context) error {
	var req struct {
		Code  string `json:"code" validate:"required"`
		State string `json:"state" validate:"required"`
	}

	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Invalid request body",
		})
	}

	if req.Code == "" || req.State == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Code and state are required",
		})
	}

	remoteName, err := h.remoteService.FinalizeAuth(req.Code, req.State)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":     true,
		"remote_name": remoteName,
		"message":     "Remote successfully configured",
	})
}

// DeleteRemote handles DELETE /api/remote/:name
func (h *RemoteHandler) DeleteRemote(c echo.Context) error {
	remoteName := c.Param("name")

	if remoteName == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "Remote name is required",
		})
	}

	if err := h.remoteService.DeleteRemote(remoteName); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "Remote deleted successfully from rclone config",
	})
}

// ListRemotes handles GET /api/remote/list (optional)
func (h *RemoteHandler) ListRemotes(c echo.Context) error {
	remotes, err := h.remoteService.ListRemotes()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"remotes": remotes,
	})
}
