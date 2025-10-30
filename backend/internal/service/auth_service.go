package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"gbackup-new/backend/internal/models"
	"gbackup-new/backend/internal/repository"

	jwt "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

var jwtSecretKey string

// init() dijalankan saat package service di-load
func init() {
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY") // Diambil dari .env
	if jwtSecretKey == "" {
		log.Fatal("FATAL: JWT_SECRET_KEY environment variable tidak diatur.")
	}
}

// Struct untuk payload login
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// AuthService interface (Kontrak)
type AuthService interface {
	Authenticate(req *LoginRequest) (string, error)
	RegisterAdmin(username, password string) error
}

// authServiceImpl: Struct implementasi
type authServiceImpl struct {
	UserRepo repository.UserRepository
}

// NewAuthService: Constructor untuk Dependency Injection
func NewAuthService(uRepo repository.UserRepository) AuthService {
	return &authServiceImpl{UserRepo: uRepo}
}

// ----------------------------------------------------
// LOGIKA UTAMA OTENTIKASI
// ----------------------------------------------------

// Authenticate: Memverifikasi user dan menghasilkan JWT.
func (s *authServiceImpl) Authenticate(req *LoginRequest) (string, error) {
	// 1. Panggil Repository (yang baru kita buat)
	user, err := s.UserRepo.FindByUsername(req.Username)
	if err != nil {
		return "", fmt.Errorf("error database saat mencari user: %w", err)
	}
	if user == nil {
		return "", fmt.Errorf("kredensial tidak valid (user tidak ditemukan)")
	}

	// 2. Verifikasi Password (Bcrypt)
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", fmt.Errorf("kredensial tidak valid (password salah)")
	}

	// 3. Generate JWT
	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(time.Hour * 72).Unix(), // Berlaku 72 jam
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("gagal menandatangani token JWT: %w", err)
	}

	return signedToken, nil
}

// RegisterAdmin: Untuk inisialisasi admin (Model Akun Tertutup)
func (s *authServiceImpl) RegisterAdmin(username, password string) error {
	user, err := s.UserRepo.FindByUsername(username)
	if err != nil {
		return err
	}
	if user != nil {
		return fmt.Errorf("user admin sudah terdaftar")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal hash password: %w", err)
	}

	newUser := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        fmt.Sprintf("%s@internal.com", username), // Email default
	}
	return s.UserRepo.CreateUser(newUser)
}
