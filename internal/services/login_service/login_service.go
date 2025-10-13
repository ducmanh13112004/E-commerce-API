package loginservice

import (
	"ecom_promotion_v2/internal/repositories"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type LoginService interface {
	Login(username, funcName, password string) (string, error)
}

type loginservice struct {
	repo      *repositories.Repositories
	jwtSecret string
}

func NewLoginService(repo *repositories.Repositories, jwtSecret string) LoginService {
	return &loginservice{repo: repo, jwtSecret: jwtSecret}
}

func (s *loginservice) Login(username, funcName, password string) (string, error) {
	user, err := s.repo.LoginApi.GetByUserName(funcName, username)
	if err != nil {
		return "", fmt.Errorf("user not found")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return "", fmt.Errorf("invalid password")
	}

	claims := jwt.MapClaims{
		"userId":   user.ID,
		"username": user.Username,
		"fullname": user.Fullname,
		"email":    user.Email,
		"exp":      time.Now().Add(time.Hour * 2).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", fmt.Errorf("failed to generate token: %v", err)
	}

	return signedToken, nil
}
