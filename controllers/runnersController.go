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

type RunnersController struct {
	runnersService *services.RunnersService
}

func NewRunnersController(runnersService *services.RunnersService) *RunnersController {
	return &RunnersController{
		runnersService: runnersService,
	}
}

func (rh RunnersController) CreateRunner(c *gin.Context) {
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
	responseErr := rh.runnersService.UpdateRunner(&runner)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	c.Status(http.StatusNoContent)
}
func (rh RunnersController) DeleteRunner(c *gin.Context) {
	runnerID := c.Param("id")
	responseErr := rh.runnersService.DeleteRunner(runnerID)
	if responseErr != nil {
		c.AbortWithStatusJSON(responseErr.Status, responseErr)
		return
	}
	c.Status(http.StatusNoContent)
}
func (rh RunnersController) GetRunner(c *gin.Context) {
	runnerID := c.Param("id")
	response, responseErr := rh.runnersService.GetRunner(runnerID)
	if responseErr != nil {
		c.JSON(responseErr.Status, responseErr)
		return
	}
	c.JSON(http.StatusOK, response)
}
func (rh RunnersController) GetRunnersBatch(c *gin.Context) {
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
