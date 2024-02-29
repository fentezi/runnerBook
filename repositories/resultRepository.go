package repositories

import (
	"database/sql"
	"github.com/fentezi/runnerBook/models"
	"net/http"
)

type ResultsRepository struct {
	dbHandler   *sql.DB
	transaction *sql.Tx
}

func NewResultsRepository(dbHandler *sql.DB) *ResultsRepository {
	return &ResultsRepository{
		dbHandler: dbHandler,
	}
}

func (rr ResultsRepository) CreateResult(result *models.Result) (*models.Result, *models.ResponseError) {
	query := `
		INSERT INTO results(runner_id, race_result, location,
		                    position, year)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id
    `
	rows, err := rr.dbHandler.Query(query, result.RunnerID, result.RaceResult,
		result.Location, result.Position, result.Year)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	var resultID string
	for rows.Next() {
		err = rows.Scan(&resultID)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return &models.Result{
		ID:         resultID,
		RunnerID:   result.RunnerID,
		RaceResult: result.RaceResult,
		Location:   result.Location,
		Position:   result.Position,
		Year:       result.Year,
	}, nil
}

func (rr ResultsRepository) DeleteResult(resultID string) (*models.Result, *models.ResponseError) {
	query := `
		DELETE FROM results
		WHERE id = $1
		RETURNING runner_id, race_result, year`
	rows, err := rr.transaction.Query(query, resultID)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	var runner_id, raceResult string
	var year int
	for rows.Next() {
		err = rows.Scan(&runner_id, &raceResult, &year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return &models.Result{
		ID:         resultID,
		RunnerID:   runner_id,
		RaceResult: raceResult,
		Year:       year,
	}, nil
}

func (rr ResultsRepository) GetAllRunnersResults(
	runnerID string) ([]*models.Result, *models.ResponseError) {
	query := `
    	SELECT id, race_result, location, position, year
		FROM results
    	WHERE runner_id = $1
    `
	rows, err := rr.dbHandler.Query(query, runnerID)
	if err != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	results := make([]*models.Result, 0)
	var id, raceResult, location string
	var position, year int
	for rows.Next() {
		err = rows.Scan(&id, &raceResult, &location, &position, &year)
		if err != nil {
			return nil, &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
			result := &models.Result{
				ID:         id,
				RunnerID:   runnerID,
				RaceResult: raceResult,
				Location:   location,
				Position:   position,
				Year:       year,
			}
			results = append(results, result)
		}
	}
	if rows.Err() != nil {
		return nil, &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return results, nil
}

func (rr ResultsRepository) GetPersonalBestResults(
	runnerID string) (string, *models.ResponseError) {
	query := `
		SELECT MIN(race_result)
		FROM results
		WHERE runner_id = $1
    `
	rows, err := rr.dbHandler.Query(query, runnerID)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	var raceResult string
	for rows.Next() {
		err = rows.Scan(&raceResult)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if rows.Err() != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return raceResult, nil
}

func (rr ResultsRepository) GetSeasonBestResults(
	runnerID string, year int) (string, *models.ResponseError) {
	query := `
    	SELECT MIN(race_result)
		FROM results
    	WHERE runner_id = $1 AND year = $2
    `
	rows, err := rr.dbHandler.Query(query, runnerID, year)
	if err != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	defer rows.Close()
	var raceResult string
	for rows.Next() {
		err = rows.Scan(&raceResult)
		if err != nil {
			return "", &models.ResponseError{
				Message: err.Error(),
				Status:  http.StatusInternalServerError,
			}
		}
	}
	if rows.Err() != nil {
		return "", &models.ResponseError{
			Message: err.Error(),
			Status:  http.StatusInternalServerError,
		}
	}
	return raceResult, nil
}
