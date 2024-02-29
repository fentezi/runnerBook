package server

import (
	"database/sql"
	"github.com/fentezi/runnerBook/controllers"
	"github.com/fentezi/runnerBook/repositories"
	"github.com/fentezi/runnerBook/services"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"log"
)

type HttpServer struct {
	config            *viper.Viper
	router            *gin.Engine
	runnersController *controllers.RunnersController
	resultController  *controllers.ResultsController
}

func InitHttpServer(config *viper.Viper,
	dbHandler *sql.DB) HttpServer {
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	resultRepository := repositories.NewResultsRepository(dbHandler)
	runnersService := services.NewRunnersService(
		runnersRepository, resultRepository)
	resultsService := services.NewResultsService(
		resultRepository, runnersRepository)
	runnersController := controllers.NewRunnersController(runnersService)
	resultsController := controllers.NewResultsController(resultsService)
	router := gin.Default()
	router.POST("/runner", runnersController.CreateRunner)
	router.PUT("/runner", runnersController.UpdateRunner)
	router.DELETE("/runner/:id", runnersController.DeleteRunner)
	router.GET("/runner/:id", runnersController.GetRunner)
	router.GET("/runner", runnersController.GetRunnersBatch)
	router.POST("/result", resultsController.CreateResult)
	router.DELETE("/result/:id", resultsController.DeleteResult)
	return HttpServer{
		config:            config,
		router:            router,
		runnersController: runnersController,
		resultController:  resultsController,
	}
}

func (h *HttpServer) Start() {
	err := h.router.Run(h.config.GetString("http.server_address"))
	if err != nil {
		log.Fatalf("Error while statring HTTP server: %v", err)
	}
}
