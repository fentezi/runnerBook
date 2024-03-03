package services

import (
	"github.com/fentezi/runnerBook/models"
	"github.com/magiconair/properties/assert"
	"net/http"
	"testing"
)

func TestValidateRunner(t *testing.T) {
	tests := []struct {
		name   string
		runner *models.Runner
		want   *models.ResponseError
	}{
		{
			name: "Invalid_First_Name",
			runner: &models.Runner{
				LastName: "Swith",
				Age:      30,
				Country:  "United States",
			},
			want: &models.ResponseError{
				Message: "Invalid first name",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Invalid_Last_Name",
			runner: &models.Runner{
				FirstName: "John",
				Age:       30,
				Country:   "United States",
			},
			want: &models.ResponseError{
				Message: "Invalid last name",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Invalid_Age",
			runner: &models.Runner{
				FirstName: "John",
				LastName:  "Smith",
				Country:   "United States",
			},
			want: &models.ResponseError{
				Message: "Invalid age",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Invalid_Country",
			runner: &models.Runner{
				FirstName: "John",
				LastName:  "Smith",
				Age:       30,
			},
			want: &models.ResponseError{
				Message: "Invalid country",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name: "Valid Runner",
			runner: &models.Runner{
				FirstName: "John",
				LastName:  "Smith",
				Age:       30,
				Country:   "United States",
			},
			want: nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			responseErr := validateRunner(test.runner)
			assert.Equal(t, test.want, responseErr)
		})
	}
}

func TestValidateRunnerID(t *testing.T) {
	tests := []struct {
		name     string
		runnerID string
		want     *models.ResponseError
	}{
		{
			name:     "Invalid_RunnerID",
			runnerID: "",
			want: &models.ResponseError{
				Message: "Invalid runner ID",
				Status:  http.StatusBadRequest,
			},
		},
		{
			name:     "Valid RunnerID",
			runnerID: "1",
			want:     nil,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			responseErr := validateRunnerID(test.runnerID)
			assert.Equal(t, test.want, responseErr)
		})
	}
}
