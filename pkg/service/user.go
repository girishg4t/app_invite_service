package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/girishg4t/app_invite_service/pkg/logging"
	"github.com/girishg4t/app_invite_service/pkg/model"
	"github.com/girishg4t/app_invite_service/pkg/repo"
	"go.uber.org/zap"
)

type UserService struct {
	repo   repo.IUser
	logger *zap.Logger
}

func NewUserService(repository repo.IUser) *UserService {
	log := logging.GetLogger().Named("user_service")
	return &UserService{
		repo:   repository,
		logger: log,
	}
}

// Login verify if the user is valid and send the token for further request
func (s UserService) Login(w http.ResponseWriter, req *http.Request) {
	s.logger.Info("Login request")
	var u model.User
	err := json.NewDecoder(req.Body).Decode(&u)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	s.logger.Debug("Found valid reqest", zap.Any("reqest", u))
	authuser, err := s.repo.GetUser(&u)
	if err != nil {
		err := errors.New("Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}

	if authuser.Password != u.Password {
		err = errors.New("Username or Password is incorrect")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	s.logger.Info("Password is valid")
	validToken, err := generateJWT(authuser)
	if err != nil {
		err = errors.New("Failed to generate token")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(err)
		return
	}
	s.logger.Debug("Found valid token", zap.String("Token", validToken))
	var token model.Token
	token.Username = authuser.Username
	token.Role = authuser.Role
	token.TokenString = validToken
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(token)
}

func generateJWT(u model.User) (string, error) {
	var mySigningKey = []byte(os.Getenv("ACCESS_SECRET"))
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["username"] = u.Username
	claims["role"] = u.Role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()
	tokenString, err := token.SignedString(mySigningKey)

	if err != nil {
		err = fmt.Errorf("Something Went Wrong: %s", err.Error())
		return "", err
	}
	return tokenString, nil
}
