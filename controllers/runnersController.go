package controllers

import (
	"encoding/json"
	"github.com/fentezi/runnerBook/models"
	"github.com/fentezi/runnerBook/services"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"net/http"
)

const ROLE_ADMIN = "admin"
const ROLE_RUNNER = "runner"

type RunnersController struct {
	runnersService *services.RunnersService
	usersService   *services.UsersService
}

func NewRunnersController(runnersService *services.RunnersService, usersService *services.UsersService) *RunnersController {
	return &RunnersController{
		runnersService: runnersService,
		usersService:   usersService,
	}
}

func (rh RunnersController) CreateRunner(c *gin.Context) {
	accessToken := c.Request.Header.Get("Token")
	auth, responseErr := rh.usersService.AuthorizeUser(
		accessToken, []string{ROLE_ADMIN})
	if responseErr != nil {
		c.JSON(responseErr.Status, responseErr)
		return
	}
	if !auth {
		c.Status(http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(
			"Error while reading create runner request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println(
			"Error while unmarshaling "+
				"create runner request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response, responseErr := rh.runnersService.CreateRunner(&runner)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	c.JSON(http.StatusOK, response)
}
func (rh RunnersController) UpdateRunner(c *gin.Context) {
	accessToken := c.Request.Header.Get("Token")
	auth, responseErr := rh.usersService.AuthorizeUser(
		accessToken, []string{ROLE_ADMIN})
	if responseErr != nil {
		c.JSON(responseErr.Status, responseErr)
		return
	}
	if !auth {
		c.Status(http.StatusUnauthorized)
		return
	}
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		log.Println(
			"Error while reading update runner request body", err)
		c.AbortWithStatusJSON(http.StatusInternalServerError, err)
		return
	}
	var runner models.Runner
	err = json.Unmarshal(body, &runner)
	if err != nil {
		log.Println(
			"Error while unmarshaling "+
				"create runner request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	responseErr = rh.runnersService.UpdateRunner(&runner)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	c.Status(http.StatusNoContent)
}
func (rh RunnersController) DeleteRunner(c *gin.Context) {
	accessToken := c.Request.Header.Get("Token")
	auth, responseErr := rh.usersService.AuthorizeUser(
		accessToken, []string{ROLE_ADMIN})
	if responseErr != nil {
		c.JSON(responseErr.Status, responseErr)
		return
	}
	if !auth {
		c.Status(http.StatusUnauthorized)
		return
	}
	runnerID := c.Param("id")
	responseErr = rh.runnersService.DeleteRunner(runnerID)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	c.Status(http.StatusNoContent)
}
func (rh RunnersController) GetRunner(c *gin.Context) {
	accessToken := c.Request.Header.Get("Token")
	auth, responseErr := rh.usersService.AuthorizeUser(
		accessToken, []string{ROLE_ADMIN, ROLE_RUNNER})
	if responseErr != nil {
		c.JSON(responseErr.Status, responseErr)
		return
	}
	if !auth {
		c.Status(http.StatusUnauthorized)
		return
	}
	runnerID := c.Param("id")
	response, responseErr := rh.runnersService.GetRunner(runnerID)
	if responseErr != nil {
		c.JSON(responseErr.Status, responseErr)
		return
	}
	c.JSON(http.StatusOK, response)
}
func (rh RunnersController) GetRunnersBatch(c *gin.Context) {
	accessToken := c.Request.Header.Get("Token")
	auth, responseErr := rh.usersService.AuthorizeUser(
		accessToken, []string{ROLE_ADMIN, ROLE_RUNNER})
	if responseErr != nil {
		c.JSON(responseErr.Status, responseErr)
		return
	}
	if !auth {
		c.Status(http.StatusUnauthorized)
		return
	}
	params := c.Request.URL.Query()
	country := params.Get("country")
	year := params.Get("year")
	response, responseErr := rh.runnersService.GetRunnersBatch(country, year)
	if responseErr != nil {
		c.JSON(responseErr.Status, responseErr)
		return
	}
	c.JSON(http.StatusOK, response)
}
