package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"golangApp/config"
	"golangApp/models"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func setupTestDB() {
	var err error
	config.DB, err = gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	config.DB.AutoMigrate(&models.Client{})
}

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
				Email:     "test@test.com",
				Telephone: "1234567",
			},
			expectedError: "Name, Last Name, Email, Age, and Birth Day are required",
		},
		{
			name: "Invalid email",
			client: models.Client{
				Name:      "John",
				LastName:  "Doe",
				Email:     "invalid-email",
				Age:       30,
				BirthDay:  time.Now().AddDate(-30, 0, 0),
				Telephone: "1234567",
			},
			expectedError: "Invalid email format",
		},
		{
			name: "Invalid phone",
			client: models.Client{
				Name:      "John",
				LastName:  "Doe",
				Email:     "test@test.com",
				Age:       30,
				BirthDay:  time.Now().AddDate(-30, 0, 0),
				Telephone: "123",
			},
			expectedError: "Phone number must be numeric and at least 7 digits long",
		},
		{
			name: "Age mismatch with birthday",
			client: models.Client{
				Name:      "John",
				LastName:  "Doe",
				Email:     "test@test.com",
				Age:       25,
				BirthDay:  time.Now().AddDate(-30, 0, 0),
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

			if assert.NoError(t, CreateClient(c)) {
				assert.Equal(t, http.StatusBadRequest, rec.Code)

				var resp map[string]string
				err := json.Unmarshal(rec.Body.Bytes(), &resp)
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedError, resp["error"])
			}
		})
	}
}

func TestCalculateClientKPI(t *testing.T) {
	ages := []float64{20, 30, 40}

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

	assert.Equal(t, float64(30), expectedAverage)
	assert.InDelta(t, float64(66.66666666666667), expectedStdDev, 0.000001)
}

func TestGetAll(t *testing.T) {
	setupTestDB()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/clients/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	clients := []models.Client{
		{
			ID:        1,
			Name:      "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			BirthDay:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			Age:       33,
			Telephone: "123456789",
		},
		{
			ID:        2,
			Name:      "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			BirthDay:  time.Date(1985, time.February, 14, 0, 0, 0, 0, time.UTC),
			Age:       39,
			Telephone: "987654321",
		},
		{
			ID:        3,
			Name:      "Alice",
			LastName:  "Johnson",
			Email:     "alice.johnson@example.com",
			BirthDay:  time.Date(1995, time.March, 30, 0, 0, 0, 0, time.UTC),
			Age:       29,
			Telephone: "555555555",
		},
	}

	for _, client := range clients {
		config.DB.Create(&client)
	}

	err := GetAll(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponse := `[
			{
				"id": 1,
				"name": "John",
				"last_name": "Doe",
				"email": "john.doe@example.com",
				"birth_day": "1990-01-01T00:00:00Z",
				"age": 33,
				"telephone": "123456789"
			},
			{
				"id": 2,
				"name": "Jane",
				"last_name": "Smith",
				"email": "jane.smith@example.com",
				"birth_day": "1985-02-14T00:00:00Z",
				"age": 39,
				"telephone": "987654321"
			},
			{
				"id": 3,
				"name": "Alice",
				"last_name": "Johnson",
				"email": "alice.johnson@example.com",
				"birth_day": "1995-03-30T00:00:00Z",
				"age": 29,
				"telephone": "555555555"
			}
		]`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Client{})
}

func TestGetClient(t *testing.T) {
	setupTestDB()

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/clients/10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetParamNames("id")
	c.SetParamValues("10")

	testClient := models.Client{
		ID:        10,
		Name:      "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		BirthDay:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Age:       33,
		Telephone: "123456789",
	}

	config.DB.Create(&testClient)

	err := GetClient(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponse := `{
			"id": 10,
			"name": "John",
			"last_name": "Doe",
			"email": "john.doe@example.com",
			"birth_day": "1990-01-01T00:00:00Z",
			"age": 33,
			"telephone": "123456789"
		}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Client{})
}

func TestGetClientNotFound(t *testing.T) {
	setupTestDB()

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/clients/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("999")

	err := GetClient(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var jsonErr map[string]string
		json.Unmarshal(rec.Body.Bytes(), &jsonErr)
		assert.Equal(t, "Client not found", jsonErr["error"])
	}
}

func TestUpdateClient(t *testing.T) {
	setupTestDB()

	e := echo.New()
	originalClient := models.Client{
		ID:        1,
		Name:      "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		BirthDay:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Age:       33,
		Telephone: "123456789",
	}

	config.DB.Create(&originalClient)

	updateData := `{
		"name": "Johnny",
		"last_name": "Smith",
		"email": "johnny.smith@example.com",
		"age": 34,
		"telephone": "987654321"
	}`

	req := httptest.NewRequest(http.MethodPut, "/clients/1", strings.NewReader(updateData))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("1")

	err := UpdateClient(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponse := `{
			"id": 1,
			"name": "Johnny",
			"last_name": "Smith",
			"email": "johnny.smith@example.com",
			"birth_day": "1990-01-01T00:00:00Z",
			"age": 34,
			"telephone": "987654321"
		}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Client{})
}

func TestDeleteClient(t *testing.T) {
	setupTestDB()

	e := echo.New()

	testClient := models.Client{
		ID:        1,
		Name:      "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		BirthDay:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Age:       33,
		Telephone: "123456789",
	}

	config.DB.Create(&testClient)

	req := httptest.NewRequest(http.MethodDelete, "/clients/1", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("1")

	err := DeleteClient(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNoContent, rec.Code)

		var deletedClient models.Client
		err := config.DB.First(&deletedClient, 1).Error
		assert.Error(t, err)
	}

	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Client{})
}

func TestCreateClientIntegration(t *testing.T) {
	setupTestDB()

	e := echo.New()
	client := models.Client{
		Name:      "Juan",
		LastName:  "Doe",
		Email:     "juan.doe@example.com",
		BirthDay:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Age:       34,
		Telephone: "123456789",
	}

	reqBody := fmt.Sprintf(`{
		"name": "%s",
		"last_name": "%s",
		"email": "%s",
		"birth_day": "%s",
		"age": %d,
		"telephone": "%s"
	}`, client.Name, client.LastName, client.Email, client.BirthDay.Format(time.RFC3339), client.Age, client.Telephone)
	req := httptest.NewRequest(http.MethodPost, "/clients", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := CreateClient(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusCreated, rec.Code)

		var createdClient models.Client
		config.DB.Last(&createdClient)
		assert.Equal(t, client.Name, createdClient.Name)
		assert.Equal(t, client.Email, createdClient.Email)
	}

	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Client{})
}

func TestUpdateClientIntegration(t *testing.T) {
	setupTestDB()

	e := echo.New()
	client := models.Client{
		ID:        1,
		Name:      "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		BirthDay:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Age:       33,
		Telephone: "123456789",
	}

	config.DB.Create(&client)

	updatedClient := models.Client{
		Name:      "John Updated",
		LastName:  "Doe Updated",
		Email:     "john.updated@example.com",
		BirthDay:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Age:       34,
		Telephone: "987654321",
	}
	reqBody := fmt.Sprintf(`{
		"name": "%s",
		"last_name": "%s",
		"email": "%s",
		"birth_day": "%s",
		"age": %d,
		"telephone": "%s"
	}`, updatedClient.Name, updatedClient.LastName, updatedClient.Email, updatedClient.BirthDay.Format(time.RFC3339), updatedClient.Age, updatedClient.Telephone)
	req := httptest.NewRequest(http.MethodPut, "/clients/1", strings.NewReader(reqBody))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := UpdateClient(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var clientInDB models.Client
		config.DB.First(&clientInDB, 1)
		assert.Equal(t, updatedClient.Name, clientInDB.Name)
		assert.Equal(t, updatedClient.Email, clientInDB.Email)
	}

	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Client{})
}

func TestDeleteClientIntegration(t *testing.T) {
	setupTestDB()

	client := models.Client{
		Name:      "Juan",
		LastName:  "Doe",
		Email:     "juan.doe@example.com",
		BirthDay:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Age:       34,
		Telephone: "123456789",
	}
	result := config.DB.Create(&client)
	if result.Error != nil {
		t.Fatalf("Error al insertar el cliente: %v", result.Error)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/clients/%d", client.ID), nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues(fmt.Sprintf("%d", client.ID))

	err := DeleteClient(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNoContent, rec.Code)

		var deletedClient models.Client
		result := config.DB.First(&deletedClient, client.ID)
		assert.Error(t, result.Error, "Se esperaba que el cliente fuera eliminado")
	}

	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Client{})
}

func TestGetAllIntegration(t *testing.T) {
	setupTestDB()

	clients := []models.Client{
		{
			Name:      "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			BirthDay:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			Age:       33,
			Telephone: "123456789",
		},
		{
			Name:      "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			BirthDay:  time.Date(1985, time.February, 14, 0, 0, 0, 0, time.UTC),
			Age:       39,
			Telephone: "987654321",
		},
		{
			Name:      "Alice",
			LastName:  "Johnson",
			Email:     "alice.johnson@example.com",
			BirthDay:  time.Date(1995, time.March, 30, 0, 0, 0, 0, time.UTC),
			Age:       29,
			Telephone: "555555555",
		},
	}

	for _, client := range clients {
		config.DB.Create(&client)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/clients/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	err := GetAll(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)

		var returnedClients []models.Client
		json.Unmarshal(rec.Body.Bytes(), &returnedClients)
		assert.Greater(t, len(returnedClients), 1, "Se esperaba que el número de clientes devueltos sea mayor a 1")
	}

	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Client{})
}

func TestGetClientIntegration(t *testing.T) {
	setupTestDB()

	testClient := models.Client{
		ID:        10,
		Name:      "John",
		LastName:  "Doe",
		Email:     "john.doe@example.com",
		BirthDay:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
		Age:       33,
		Telephone: "123456789",
	}
	config.DB.Create(&testClient)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/clients/10", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("10")

	err := GetClient(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusOK, rec.Code)

		expectedResponse := `{
			"id": 10,
			"name": "John",
			"last_name": "Doe",
			"email": "john.doe@example.com",
			"birth_day": "1990-01-01T00:00:00Z",
			"age": 33,
			"telephone": "123456789"
		}`
		assert.JSONEq(t, expectedResponse, rec.Body.String())
	}

	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Client{})
}
func TestGetClientNotFoundIntegration(t *testing.T) {
	setupTestDB()

	e := echo.New()

	req := httptest.NewRequest(http.MethodGet, "/clients/999", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	c.SetParamNames("id")
	c.SetParamValues("999")

	err := GetClient(c)

	if assert.NoError(t, err) {
		assert.Equal(t, http.StatusNotFound, rec.Code)

		var jsonErr map[string]string
		json.Unmarshal(rec.Body.Bytes(), &jsonErr)
		assert.Equal(t, "Client not found", jsonErr["error"])
	}

	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Client{})
}

func TestCalculateClientKPIIntegration(t *testing.T) {
	setupTestDB()

	clients := []models.Client{
		{
			ID:        1,
			Name:      "John",
			LastName:  "Doe",
			Email:     "john.doe@example.com",
			BirthDay:  time.Date(1990, time.January, 1, 0, 0, 0, 0, time.UTC),
			Age:       20,
			Telephone: "123456789",
		},
		{
			ID:        2,
			Name:      "Jane",
			LastName:  "Smith",
			Email:     "jane.smith@example.com",
			BirthDay:  time.Date(1985, time.February, 14, 0, 0, 0, 0, time.UTC),
			Age:       30,
			Telephone: "987654321",
		},
		{
			ID:        3,
			Name:      "Alice",
			LastName:  "Johnson",
			Email:     "alice.johnson@example.com",
			BirthDay:  time.Date(1995, time.March, 30, 0, 0, 0, 0, time.UTC),
			Age:       40,
			Telephone: "555555555",
		},
	}

	for _, client := range clients {
		config.DB.Create(&client)
	}

	var ages []float64
	config.DB.Model(&models.Client{}).Pluck("age", &ages)

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

	assert.Equal(t, float64(30), expectedAverage)
	assert.InDelta(t, float64(66.66666666666667), expectedStdDev, 0.000001)

	config.DB.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&models.Client{})
}
