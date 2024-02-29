package services

import (
	"github.com/fentezi/runnerBook/models"
	"github.com/fentezi/runnerBook/repositories"
	"net/http"
	"strconv"
	"time"
)

type RunnersService struct {
	runnersRepository *repositories.RunnersRepository
	resultsRepository *repositories.ResultsRepository
}

func NewRunnersService(
	runnersRepository *repositories.RunnersRepository,
	resultsRepository *repositories.ResultsRepository) *RunnersService {
	return &RunnersService{
		runnersRepository: runnersRepository,
		resultsRepository: resultsRepository,
	}
}

func (rs RunnersService) CreateRunner(
	runner *models.Runner) (*models.Runner, *models.ResponseError) {
	responseErr := validateRunner(runner)
	if responseErr != nil {
		return nil, responseErr
	}
	return rs.runnersRepository.CreateRunner(runner)
}

func (rs RunnersService) UpdateRunner(
	runner *models.Runner) *models.ResponseError {
	responseErr := validateRunnerID(runner.ID)
	if responseErr != nil {
		return responseErr
	}
	responseErr = validateRunner(runner)
	if responseErr != nil {
		return responseErr
	}
	return rs.runnersRepository.UpdateRunner(runner)
}

func (rs RunnersService) DeleteRunner(
	runnerID string) *models.ResponseError {
	responseErr := validateRunnerID(runnerID)
	if responseErr != nil {
		return responseErr
	}
	return rs.runnersRepository.DeleteRunner(runnerID)
}

func (rs RunnersService) GetRunner(
	runnerID string) (*models.Runner, *models.ResponseError) {
	responseErr := validateRunnerID(runnerID)
	if responseErr != nil {
		return nil, responseErr
	}
	runner, responseErr := rs.runnersRepository.GetRunner(runnerID)
	if responseErr != nil {
		return nil, responseErr
	}
	results, responseErr := rs.resultsRepository.GetAllRunnersResults(runnerID)
	if responseErr != nil {
		return nil, responseErr
	}
	runner.Results = results
	return runner, nil
}

func (rs RunnersService) GetRunnersBatch(
	country, year string) ([]*models.Runner, *models.ResponseError) {
	if country != "" && year != "" {
		return nil, &models.ResponseError{
			Message: "Only one parameter can be passed",
			Status:  http.StatusBadRequest,
		}
	}
	if country != "" {
		return rs.runnersRepository.GetRunnersByCountry(country)
	}
	if year != "" {
		intYear, err := strconv.Atoi(year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}
		currentYear := time.Now().Year()
		if intYear < 0 || intYear > currentYear {
			return nil, &models.ResponseError{
				Message: "Invalid year",
				Status:  http.StatusBadRequest,
			}
		}
		return rs.runnersRepository.GetRunnersByYear(intYear)
	}
	return rs.runnersRepository.GetAllRunners()
}

func validateRunner(runner *models.Runner) *models.ResponseError {
	if runner.FirstName == "" {
		return &models.ResponseError{
			Message: "Invalid first name",
			Status:  http.StatusBadRequest,
		}
	}
	if runner.LastName == "" {
		return &models.ResponseError{
			Message: "Invalid last name",
			Status:  http.StatusBadRequest,
		}
	}
	if runner.Age <= 16 || runner.Age > 125 {
		return &models.ResponseError{
			Message: "Invalid age",
			Status:  http.StatusBadRequest,
		}
	}
	if runner.Country == "" {
		return &models.ResponseError{
			Message: "Invalid country",
			Status:  http.StatusBadRequest,
		}
	}
	return nil
}

func validateRunnerID(runnerID string) *models.ResponseError {
	if runnerID == "" {
		return &models.ResponseError{
			Message: "Invalid runner ID",
			Status:  http.StatusBadRequest,
		}
	}
	return nil
}
