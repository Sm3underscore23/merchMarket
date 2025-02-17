package service

import (
	"crypto/sha1"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/pkg/repository"
)

const (
	salt      = "jernfn32esck4334ejrkf54"
	signedKey = "ejfnw67qd732vbewx38"
	tokenTTL  = 12 * time.Hour
)

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) GetUser(username, password string) (int, error) {
	id, err := s.repo.GetUser(username, generatePasswordHash(password))
	if errors.Is(err, customerrors.ErrUserNotFound) {
		if err := s.CreateUser(username, password); err != nil {
			return 0, err
		}
	} else if errors.Is(err, customerrors.ErrWrongPasswod) {
		return 0, err
	}
	return id, nil
}

func (s *AuthService) CreateUser(username, password string) error {
	password_hash := generatePasswordHash(password)
	return s.repo.CreateUser(username, password_hash)
}

func (s *AuthService) GenerateToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})

	return token.SignedString([]byte(signedKey))
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))

	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}
