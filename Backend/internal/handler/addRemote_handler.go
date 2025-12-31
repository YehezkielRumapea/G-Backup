// internal/handler/remote_handler.go
package handler

import (
	"gbackup-new/backend/internal/service"
	"net/http"
	"strings"

	"github.com/labstack/echo/v4"
)

type RemoteHandler struct {
	remoteService service.AddRemoteService
}

func NewAddRemoteHandler(remoteService service.AddRemoteService) *RemoteHandler {
	return &RemoteHandler{remoteService: remoteService}
}

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

	baseURL := getBaseURL(c)

	authURL, state, err := h.remoteService.InitAuth(req.Name, baseURL)
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

func getBaseURL(c echo.Context) string {
	scheme := c.Request().Header.Get("X-Forwarded-Proto")
	if scheme == "" {
		if c.Request().TLS != nil {
			scheme = "https"
		} else {
			scheme = "http"
		}
	}

	host := c.Request().Header.Get("X-Forwarded-Host")
	if host == "" {
		host = c.Request().Host
	}

	if scheme == "http" && strings.HasSuffix(host, ":80") {
		host = strings.TrimSuffix(host, ":80")
	}
	if scheme == "https" && strings.HasSuffix(host, ":443") {
		host = strings.TrimSuffix(host, ":443")
	}

	return scheme + "://" + host
}

func (h *RemoteHandler) OAuthCallback(c echo.Context) error {
	code := c.QueryParam("code")
	state := c.QueryParam("state")
	errorParam := c.QueryParam("error")

	if errorParam != "" {
		return c.HTML(http.StatusOK, `<!DOCTYPE html>
<html>
<head><title>Cancelled</title></head>
<body style="font-family:Arial;text-align:center;padding:50px;background:#f5f5f5;">
<div style="background:white;padding:40px;border-radius:8px;max-width:400px;margin:0 auto;box-shadow:0 2px 10px rgba(0,0,0,0.1);">
<h2 style="color:#f44336;">❌ Authentication Cancelled</h2>
<p style="color:#666;">You can close this window.</p>
</div>
</body></html>`)
	}

	if code == "" || state == "" {
		return c.HTML(http.StatusBadRequest, `<!DOCTYPE html>
<html>
<head><title>Error</title></head>
<body style="font-family:Arial;text-align:center;padding:50px;background:#f5f5f5;">
<div style="background:white;padding:40px;border-radius:8px;max-width:400px;margin:0 auto;box-shadow:0 2px 10px rgba(0,0,0,0.1);">
<h2 style="color:#f44336;">❌ Error: Missing Parameters</h2>
<p style="color:#666;">Please try again.</p>
</div>
</body></html>`)
	}

	return c.HTML(http.StatusOK, `<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>Authentication Success</title>
<style>
*{margin:0;padding:0;box-sizing:border-box;}
body{font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',sans-serif;background:linear-gradient(135deg,#667eea 0%,#764ba2 100%);min-height:100vh;display:flex;align-items:center;justify-content:center;padding:20px;}
.container{background:white;border-radius:16px;box-shadow:0 20px 60px rgba(0,0,0,0.3);padding:40px;max-width:500px;width:100%;}
.icon{text-align:center;font-size:64px;margin-bottom:20px;}
h1{color:#4CAF50;font-size:28px;text-align:center;margin-bottom:10px;}
.msg{text-align:center;color:#666;font-size:14px;margin-bottom:20px;}
.auto{text-align:center;color:#4CAF50;font-size:13px;padding:15px;background:#e8f5e9;border-radius:8px;margin:20px 0;}
.spinner{border:3px solid #f3f3f3;border-top:3px solid #4CAF50;border-radius:50%;width:40px;height:40px;animation:spin 1s linear infinite;margin:20px auto;}
@keyframes spin{0%{transform:rotate(0deg);}100%{transform:rotate(360deg);}}
.btn{background:#666;color:white;border:none;padding:12px 24px;border-radius:6px;cursor:pointer;font-size:14px;font-weight:600;width:100%;margin-top:20px;}
.btn:hover{background:#444;}
</style>
</head>
<body>
<div class="container">
<div class="icon">✅</div>
<h1>Autentikasi Berhasil!</h1>
<p class="msg">Google Drive berhasil diotorisasi</p>
<div class="auto">⚡ Proses finalisasi sedang berjalan...</div>
<div class="spinner"></div>
<button class="btn" onclick="window.close()" style="display:none;" id="closeBtn">Close Window</button>
</div>
<script>
fetch('/api/v1/remote/finalize', {
    method: 'POST',
    headers: {'Content-Type': 'application/json'},
    body: JSON.stringify({code: '`+code+`', state: '`+state+`'})
})
.then(r => r.json())
.then(d => {
    const autoDiv = document.querySelector('.auto');
    const spinner = document.querySelector('.spinner');
    const closeBtn = document.getElementById('closeBtn');
    
    if(d.success) {
        autoDiv.innerHTML = '✅ Remote "' + d.remote_name + '" berhasil ditambahkan!';
        autoDiv.style.background = '#e8f5e9';
        autoDiv.style.color = '#4CAF50';
        spinner.style.display = 'none';
        closeBtn.style.display = 'block';
        setTimeout(() => window.close(), 3000);
    } else {
        autoDiv.innerHTML = '❌ Gagal: ' + (d.error || 'Unknown error');
        autoDiv.style.background = '#ffebee';
        autoDiv.style.color = '#c62828';
        spinner.style.display = 'none';
        closeBtn.style.display = 'block';
    }
})
.catch(e => {
    const autoDiv = document.querySelector('.auto');
    const spinner = document.querySelector('.spinner');
    const closeBtn = document.getElementById('closeBtn');
    
    autoDiv.innerHTML = '❌ Error: ' + e.message;
    autoDiv.style.background = '#ffebee';
    autoDiv.style.color = '#c62828';
    spinner.style.display = 'none';
    closeBtn.style.display = 'block';
});
</script>
</body>
</html>`)
}

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
		return c.JSON(http.StatusBadRequest, map[string]interface{}{
			"success": false,
			"error":   "Code and state are required",
		})
	}

	remoteName, err := h.remoteService.FinalizeAuth(req.Code, req.State)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		})
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"success":     true,
		"remote_name": remoteName,
		"message":     "Remote successfully configured",
	})
}

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
		"message": "Remote deleted successfully",
	})
}

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
