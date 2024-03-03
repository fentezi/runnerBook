package services

import (
	"encoding/base64"
	"github.com/fentezi/runnerBook/models"
	"github.com/fentezi/runnerBook/repositories"
	"golang.org/x/crypto/bcrypt"
	"net/http"
)

type UsersService struct {
	usersRepository *repositories.UsersRepository
}

func NewUsersService(usersRepository *repositories.UsersRepository) *UsersService {
	return &UsersService{usersRepository: usersRepository}
}

func (uc UsersService) Login(username, password string) (string, *models.ResponseError) {
	if username == "" || password == "" {
		return "", &models.ResponseError{
			Message: "Invalid username or password",
			Status:  http.StatusBadRequest,
		}
	}
	id, responseErr := uc.usersRepository.LoginUser(
		username, password)
	if responseErr != nil {
		return "", responseErr
	}
	if id == "" {
		return "", &models.ResponseError{
			Message: "Login failed",
			Status:  http.StatusUnauthorized,
		}
	}
	accessToken, responseErr := generateAccessToken(username)
	if responseErr != nil {
		return "", responseErr
	}
	uc.usersRepository.SetAccessToken(accessToken, id)
	return accessToken, nil
}

func (uc UsersService) Logout(accessToken string) *models.ResponseError {
	if accessToken == "" {
		return &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}
	return uc.usersRepository.RemoveAccessToken(accessToken)
}

func (uc UsersService) AuthorizeUser(accessToken string, expectedRoles []string) (bool, *models.ResponseError) {
	if accessToken == "" {
		return false, &models.ResponseError{
			Message: "Invalid access token",
			Status:  http.StatusBadRequest,
		}
	}
	role, responseErr := uc.usersRepository.GetUserRole(accessToken)
	if responseErr != nil {
		return false, responseErr
	}
	if role == "" {
		return false, &models.ResponseError{
			Message: "Failed to authorize user",
			Status:  http.StatusUnauthorized,
		}
	}
	for _, expectedRole := range expectedRoles {
		if expectedRole == role {
			return true, nil
		}
	}
	return false, nil
}

func generateAccessToken(username string) (string, *models.ResponseError) {
	hash, err := bcrypt.GenerateFromPassword([]byte(username), bcrypt.DefaultCost)
	if err != nil {
		return "", &models.ResponseError{
			Message: "Failed to generate token",
			Status:  http.StatusInternalServerError,
		}
	}
	return base64.StdEncoding.EncodeToString(hash), nil
}
