// internal/service/remote_service.go
package service

import (
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os/exec"
	"regexp"
	"strings"
	"sync"
	"time"
)

type AddRemoteService interface {
	InitAuth(remoteName string) (authURL string, state string, err error)
	FinalizeAuth(code, state string) (remoteName string, err error)
	DeleteRemote(remoteName string) error
	ListRemotes() ([]string, error)
	VerifyRemote(remoteName string) error
}

type addremoteServiceImpl struct {
	redirectURI string
	sessions    sync.Map
}

type authSession struct {
	RemoteName string
	Timestamp  time.Time
}

func NewAddRemoteService(redirectURI string) AddRemoteService {
	service := &addremoteServiceImpl{
		redirectURI: redirectURI,
	}

	go service.cleanupExpiredSessions()

	return service
}

// ============================================================
// PUBLIC METHODS
// ============================================================

func (s *addremoteServiceImpl) InitAuth(remoteName string) (string, string, error) {
	if err := s.validateRemoteName(remoteName); err != nil {
		return "", "", err
	}

	if s.remoteExistsInRclone(remoteName) {
		return "", "", errors.New("remote name already exists in rclone config")
	}

	state := generateStateToken()

	s.sessions.Store(state, authSession{
		RemoteName: remoteName,
		Timestamp:  time.Now(),
	})

	authURL := s.buildGoogleAuthURL(state)

	return authURL, state, nil
}

func (s *addremoteServiceImpl) FinalizeAuth(code, state string) (string, error) {
	sessionValue, exists := s.sessions.Load(state)
	if !exists {
		return "", errors.New("invalid or expired session")
	}

	session := sessionValue.(authSession)
	s.sessions.Delete(state)

	if time.Since(session.Timestamp) > 15*time.Minute {
		return "", errors.New("session expired")
	}

	token, err := s.exchangeCodeForToken(code)
	if err != nil {
		return "", fmt.Errorf("failed to exchange token: %w", err)
	}

	if err := s.createRcloneConfig(session.RemoteName, token); err != nil {
		return "", fmt.Errorf("failed to create rclone config: %w", err)
	}

	if err := s.verifyRemoteConnection(session.RemoteName); err != nil {
		s.deleteRcloneConfig(session.RemoteName)
		return "", fmt.Errorf("failed to verify connection: %w", err)
	}

	return session.RemoteName, nil
}

func (s *addremoteServiceImpl) DeleteRemote(remoteName string) error {
	if !s.remoteExistsInRclone(remoteName) {
		return errors.New("remote not found in rclone config")
	}

	return s.deleteRcloneConfig(remoteName)
}

func (s *addremoteServiceImpl) ListRemotes() ([]string, error) {
	cmd := exec.Command("rclone", "listremotes")

	output, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("failed to list remotes: %w", err)
	}

	var remotes []string
	for _, line := range strings.Split(string(output), "\n") {
		line = strings.TrimSpace(line)
		if line != "" {
			remotes = append(remotes, strings.TrimSuffix(line, ":"))
		}
	}

	return remotes, nil
}

func (s *addremoteServiceImpl) VerifyRemote(remoteName string) error {
	return s.verifyRemoteConnection(remoteName)
}

// ============================================================
// PRIVATE HELPERS
// ============================================================

func (s *addremoteServiceImpl) buildGoogleAuthURL(state string) string {
	params := url.Values{}
	params.Set("client_id", "202264815644.apps.googleusercontent.com")
	params.Set("redirect_uri", s.redirectURI)
	params.Set("response_type", "code")
	params.Set("scope", "https://www.googleapis.com/auth/drive")
	params.Set("access_type", "offline")
	params.Set("state", state)
	params.Set("prompt", "consent")

	return "https://accounts.google.com/o/oauth2/auth?" + params.Encode()
}

func (s *addremoteServiceImpl) exchangeCodeForToken(code string) (string, error) {
	tokenURL := "https://oauth2.googleapis.com/token"

	data := url.Values{}
	data.Set("code", code)
	data.Set("client_id", "202264815644.apps.googleusercontent.com")
	data.Set("client_secret", "X4Z3ca8xfWDb1Voo-F9a7ZxJ")
	data.Set("redirect_uri", s.redirectURI)
	data.Set("grant_type", "authorization_code")

	resp, err := http.PostForm(tokenURL, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("token exchange failed: %s", string(body))
	}

	var tokenResponse struct {
		AccessToken  string `json:"access_token"`
		RefreshToken string `json:"refresh_token"`
		TokenType    string `json:"token_type"`
		ExpiresIn    int    `json:"expires_in"`
	}

	if err := json.Unmarshal(body, &tokenResponse); err != nil {
		return "", err
	}

	tokenData := map[string]interface{}{
		"access_token":  tokenResponse.AccessToken,
		"token_type":    tokenResponse.TokenType,
		"refresh_token": tokenResponse.RefreshToken,
		"expiry":        time.Now().Add(time.Duration(tokenResponse.ExpiresIn) * time.Second).Format(time.RFC3339),
	}

	tokenJSON, err := json.Marshal(tokenData)
	if err != nil {
		return "", err
	}

	return string(tokenJSON), nil
}

func (s *addremoteServiceImpl) createRcloneConfig(remoteName, token string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "rclone", "config", "create",
		remoteName,
		"drive",
		"token", token,
		"scope", "drive",
	)

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("rclone config create failed: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

func (s *addremoteServiceImpl) verifyRemoteConnection(remoteName string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "rclone", "lsd", remoteName+":", "--max-depth", "1")

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("verification failed: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

func (s *addremoteServiceImpl) deleteRcloneConfig(remoteName string) error {
	cmd := exec.Command("rclone", "config", "delete", remoteName)

	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("delete failed: %w, stderr: %s", err, stderr.String())
	}

	return nil
}

func (s *addremoteServiceImpl) remoteExistsInRclone(remoteName string) bool {
	cmd := exec.Command("rclone", "listremotes")

	output, err := cmd.Output()
	if err != nil {
		return false
	}

	remotes := string(output)
	expectedRemote := remoteName + ":"

	for _, line := range strings.Split(remotes, "\n") {
		if strings.TrimSpace(line) == expectedRemote {
			return true
		}
	}

	return false
}

func (s *addremoteServiceImpl) validateRemoteName(name string) error {
	if name == "" {
		return errors.New("remote name cannot be empty")
	}
	if len(name) > 64 {
		return errors.New("remote name too long (max 64 characters)")
	}
	matched, _ := regexp.MatchString("^[a-z0-9-]+$", name)
	if !matched {
		return errors.New("remote name can only contain lowercase letters, numbers, and dashes")
	}
	return nil
}

func (s *addremoteServiceImpl) cleanupExpiredSessions() {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now()
		s.sessions.Range(func(key, value interface{}) bool {
			session := value.(authSession)
			if now.Sub(session.Timestamp) > 15*time.Minute {
				s.sessions.Delete(key)
			}
			return true
		})
	}
}

func generateStateToken() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
