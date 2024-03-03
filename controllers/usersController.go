package controllers

import (
	"github.com/fentezi/runnerBook/services"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

type UsersController struct {
	usersService *services.UsersService
}

func NewUsersController(usersService *services.UsersService) *UsersController {
	return &UsersController{usersService: usersService}
}

func (uc UsersController) Login(c *gin.Context) {
	username, password, ok := c.Request.BasicAuth()
	if !ok {
		log.Println("Error while reading credentials")
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	accessToken, responseErr := uc.usersService.Login(username, password)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	c.JSON(http.StatusOK, accessToken)
}

func (uc UsersController) Logout(c *gin.Context) {
	accessToken := c.Request.Header.Get("Token")
	responseErr := uc.usersService.Logout(accessToken)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	c.Status(http.StatusNoContent)
}
