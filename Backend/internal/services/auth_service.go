package services

import (
	"fmt"
	"time"

	"gbackup-system/backend/internal/models"
	"gbackup-system/backend/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Auth Service Interace
type AuthService interface {
	Authenticate(req *LoginRequest) (string, error)
	RegisterAdmin(username, password string) error
	SetSecretKey(key string)
}

// Implemen AuthService
type AuthServiceImpl struct {
	userRepo     repository.UserRepository
	jwtSecretKey string
}

// Constructoe Dependency Injection
func NewAuthService(userRepo repository.UserRepository) AuthService {
	return &AuthServiceImpl{userRepo: userRepo}
}

// Method to set secret key
func (s *AuthServiceImpl) SetSecretKey(key string) {
	s.jwtSecretKey = key
}

// Main logic autentication
func (s *AuthServiceImpl) Authenticate(req *LoginRequest) (string, error) {
	// Cari user di DB
	user, err := s.userRepo.FindByUsername(req.Username)
	if err != nil {
		return "", fmt.Errorf("gagal mengambil user dari DB: %w", err)
	}

	if user == nil {
		return "", fmt.Errorf("user tidak ditemukan")
	}

	// Verifikasi password
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		return "", fmt.Errorf("password salah")
	}

	// Jwt Token
	claims := jwt.MapClaims{
		"User_ID":  user.ID,
		"Username": user.Username,
		"exp":      time.Now().Add(time.Hour * 1).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(s.jwtSecretKey))
	if err != nil {
		return "", fmt.Errorf("gagal membuat token: %w", err)
	}

	return signedToken, nil
}

// Logic Register

func (s *AuthServiceImpl) RegisterAdmin(username, password string) error {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return err
	}

	if user != nil {
		return fmt.Errorf("user admin sudah ada")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return fmt.Errorf("gagal hash password: %w", err)
	}

	newUser := &models.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}
	return s.userRepo.CreateUser(newUser)
}
