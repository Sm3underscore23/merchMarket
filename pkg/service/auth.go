package service

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"

	"github.com/Sm3underscore23/merchStore/internal/customerrors"
	"github.com/Sm3underscore23/merchStore/internal/models"
	"github.com/Sm3underscore23/merchStore/pkg/repository"
)

const tokenTTL = 12 * time.Hour

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	authRepo         repository.Authorization
	userProviderRepo repository.UserProvider
	authConfig       models.AuthConfig
}

func NewAuthService(
	authRepo repository.Authorization,
	userProviderRepo repository.UserProvider,
	authConfig models.AuthConfig,
) *AuthService {
	return &AuthService{
		authRepo:         authRepo,
		userProviderRepo: userProviderRepo,
		authConfig:       authConfig,
	}
}

func generatePasswordHash(password string, salt string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(salt)))
}

func (s *AuthService) AuthUser(username, password string) (string, error) {
	id, err := s.userProviderRepo.GetUserId(username)
	if errors.Is(err, customerrors.ErrUserNotFound) {
		id, err = s.userProviderRepo.CreateUser(
			username,
			generatePasswordHash(
				password,
				s.authConfig.Salt,
			))
		if err != nil {
			return "", err
		}
	} else if err != nil {
		return "", err
	}

	err = s.authRepo.CheckPassword_hash(id,
		generatePasswordHash(
			password,
			s.authConfig.Salt,
		))

	if err != nil {
		return "", err
	}

	token, err := s.GenerateToken(id)
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *AuthService) GenerateToken(id int) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		id,
	})

	signedKey := s.authConfig.SignedKey

	return token.SignedString([]byte(signedKey))
}

func (s *AuthService) ParseToken(accessToken string) (int, error) {
	signedKey := s.authConfig.SignedKey
	token, err := jwt.ParseWithClaims(accessToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid sighing method")
		}
		return []byte(signedKey), nil
	})

	if err != nil {
		return 0, err
	}

	claims, ok := token.Claims.(*tokenClaims)

	if !ok {
		return 0, errors.New("token claims are not type *tokenClaims")
	}

	return claims.UserId, nil
}
