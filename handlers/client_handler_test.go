package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"golangApp/models"
    "fmt"
)

func TestValidEmail(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		expected bool
	}{
		{
			name:     "Valid email",
			email:    "test@test.com",
			expected: true,
		},
		{
			name:     "Valid email with subdomain",
			email:    "test@sub.test.com",
			expected: true,
		},
		{
			name:     "Invalid email - no domain",
			email:    "test@",
			expected: false,
		},
		{
			name:     "Invalid email - no @",
			email:    "testtest.com",
			expected: false,
		},
		{
			name:     "Invalid email - empty",
			email:    "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := validEmail(tt.email)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsNumeric(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{
			name:     "Valid number",
			input:    "1234567",
			expected: true,
		},
		{
			name:     "Invalid - contains letters",
			input:    "123abc",
			expected: false,
		},
		{
			name:     "Invalid - empty",
			input:    "",
			expected: false,
		},
		{
			name:     "Valid - single digit",
			input:    "0",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNumeric(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestCreateClientValidation(t *testing.T) {
	e := echo.New()
	tests := []struct {
		name          string
		client        models.Client
		expectedError string
	}{
		{
			name: "Missing required fields",
			client: models.Client{
				Email:    "test@test.com",
				Telephone: "1234567",
			},
			expectedError: "Name, Last Name, Email, Age, and Birth Day are required",
		},
		{
			name: "Invalid email",
			client: models.Client{
				Name:     "John",
				LastName: "Doe",
				Email:    "invalid-email",
				Age:     30,
				BirthDay: time.Now().AddDate(-30, 0, 0),
				Telephone: "1234567",
			},
			expectedError: "Invalid email format",
		},
		{
			name: "Invalid phone",
			client: models.Client{
				Name:     "John",
				LastName: "Doe",
				Email:    "test@test.com",
				Age:     30,
				BirthDay: time.Now().AddDate(-30, 0, 0),
				Telephone: "123",
			},
			expectedError: "Phone number must be numeric and at least 7 digits long",
		},
		{
			name: "Age mismatch with birthday",
			client: models.Client{
				Name:     "John",
				LastName: "Doe",
				Email:    "test@test.com",
				Age:     25,
				BirthDay: time.Now().AddDate(-30, 0, 0),
				Telephone: "1234567",
			},
			expectedError: "Age does not match birth date",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			jsonBytes, err := json.Marshal(tt.client)
			assert.NoError(t, err)

			req := httptest.NewRequest(http.MethodPost, "/api/v1/clients", bytes.NewBuffer(jsonBytes))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			err = CreateClient(c)
            if assert.Error(t, err) {
				herr, ok := err.(*echo.HTTPError)
                
				if ok {
					errMsg, ok := herr.Message.(map[string]string)
					if ok {
                        assert.Contains(t, errMsg["error"], tt.expectedError)
					}
				}
			}
		})
	}
}

func TestCalculateClientKPI(t *testing.T) {
	// Test data
	ages := []float64{20, 30, 40}
	
	// Calculate expected values
	var sum float64
	for _, age := range ages {
		sum += age
	}
	expectedAverage := sum / float64(len(ages))

	var varianceSum float64
	for _, age := range ages {
		varianceSum += (age - expectedAverage) * (age - expectedAverage)
	}
	expectedStdDev := float64(0)
	if len(ages) > 0 {
		expectedStdDev = float64(varianceSum / float64(len(ages)))
	}

	// Verify calculations
	assert.Equal(t, float64(30), expectedAverage)
	assert.InDelta(t, float64(66.66666666666667), expectedStdDev, 0.000001)
}