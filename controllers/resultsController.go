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

type ResultsController struct {
	resultsService *services.ResultsService
	usersService   *services.UsersService
}

func NewResultsController(resultsService *services.ResultsService, usersService *services.UsersService) *ResultsController {
	return &ResultsController{
		resultsService: resultsService,
		usersService:   usersService,
	}
}

func (rh ResultsController) CreateResult(c *gin.Context) {
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
			"Error while reading create result request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	var result models.Result
	err = json.Unmarshal(body, &result)
	if err != nil {
		log.Println(
			"Error while unmarshaling "+
				"creates result request body", err)
		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	response, responseErr := rh.resultsService.CreateResult(&result)
	if responseErr != nil {
		c.JSON(responseErr.Status, responseErr)
		return
	}
	c.JSON(http.StatusOK, response)
}

func (rh *ResultsController) DeleteResult(c *gin.Context) {
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
	resultID := c.Param("id")
	responseErr = rh.resultsService.DeleteResult(resultID)
	if responseErr != nil {
		c.JSON(responseErr.Status, responseErr)
		return
	}
	c.Status(http.StatusNoContent)
}
