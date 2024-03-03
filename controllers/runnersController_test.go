package controllers

import (
	"database/sql"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/fentezi/runnerBook/models"
	"github.com/fentezi/runnerBook/repositories"
	"github.com/fentezi/runnerBook/services"
	"github.com/gin-gonic/gin"
	"github.com/magiconair/properties/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func initTestRouter(dbHandler *sql.DB) *gin.Engine {
	runnersRepository := repositories.NewRunnersRepository(dbHandler)
	runnersService := services.NewRunnersService(runnersRepository, nil)
	runnersController := NewRunnersController(runnersService)
	router := gin.Default()
	router.GET("/runner", runnersController.GetRunnersBatch)
	router.GET("/runner/:id", runnersController.GetRunner)
	return router
}

func TestGetRunner(t *testing.T) {
	dbHandler, mock, _ := sqlmock.New()
	defer dbHandler.Close()
	columns := []string{"id", "first_name", "last_name", "age",
		"is_active", "country", "personal_best", "season_best"}
	query := `SELECT * FROM runners WHERE id = $1`
	rows := mock.NewRows(columns).AddRow("1", "John", "Smith", 30, true, "United States", "02:00:41", "02:13:13")
	mock.ExpectQuery(query).WithArgs("1").WillReturnRows(rows)
	router := initTestRouter(dbHandler)
	request, _ := http.NewRequest("GET", "/runner/1", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	t.Log(recorder.Body.String())
	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
	var runner []*models.Runner
	json.Unmarshal(recorder.Body.Bytes(), &runner)
	assert.Equal(t, 1, len(runner))
}

func TestGetRunnersResponse(t *testing.T) {
	dhHandler, mock, _ := sqlmock.New()
	defer dhHandler.Close()
	columns := []string{"id", "first_name", "last_name", "age",
		"is_active", "country", "personal_best", "season_best"}
	mock.ExpectQuery("SELECT *").WillReturnRows(
		sqlmock.NewRows(columns).
			AddRow("1", "John", "Smith", 30, true,
				"United States", "02:00:41", "02:13:13").
			AddRow("2", "Marjanna", "Komathic", 24, true,
				"Serbia", "01:18:28", "01:18:28"))
	router := initTestRouter(dhHandler)
	request, _ := http.NewRequest("GET", "/runner", nil)
	recorder := httptest.NewRecorder()
	router.ServeHTTP(recorder, request)
	assert.Equal(t, http.StatusOK, recorder.Result().StatusCode)
	var runners []*models.Runner
	json.Unmarshal(recorder.Body.Bytes(), &runners)
	assert.Equal(t, 2, len(runners))
}
