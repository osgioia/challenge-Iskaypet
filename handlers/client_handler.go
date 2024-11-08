package handlers

import (
	"math"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"golangApp/config"
	"golangApp/models"

	"github.com/labstack/echo/v4"
)

type ClientKPI struct {
	AverageAge           float64 `json:"average_age"`
	AgeStandardDeviation float64 `json:"age_standard_deviation"`
}

// GetAll obtiene todos los clientes
// @Summary Obtiene todos los clientes
// @Description Recupera una lista de todos los clientes en la base de datos
// @Tags Clientes
// @Produce json
// @Success 200 {array} models.Client "Lista de clientes"
// @Router /api/v1/clients [get]
func GetAll(c echo.Context) error {
	var clients []models.Client
	if err := config.DB.Find(&clients).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, clients)
}

// GetClient obtiene un cliente por ID
// @Summary Obtener cliente por ID
// @Description Recupera un cliente específico usando su ID
// @Tags Clientes
// @Param id path int true "ID del Cliente"
// @Produce json
// @Success 200 {object} models.Client "Detalles del cliente"
// @Router /api/v1/clients/{id} [get]
func GetClient(c echo.Context) error {
	id := c.Param("id")
	var client models.Client
	if err := config.DB.First(&client, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
	}
	return c.JSON(http.StatusOK, client)
}

func validEmail(email string) bool {
	const emailRegex = `^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,}$`
	re := regexp.MustCompile(emailRegex)
	return re.MatchString(email)
}

// CreateClient crea un nuevo cliente
// @Summary Crear cliente
// @Description Crea un nuevo cliente con los datos proporcionados
// @Tags Clientes
// @Accept json
// @Produce json
// @Param client body models.Client true "Información del Cliente"
// @Success 201 {object} models.Client "Cliente creado exitosamente"
// @Router /api/v1/clients [post]
func CreateClient(c echo.Context) error {
	var client models.Client
	if err := c.Bind(&client); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
	}

	if client.Name == "" || client.LastName == "" || client.Email == "" || client.Age == 0 || client.BirthDay.IsZero() {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Name, Last Name, Email, Age, and Birth Day are required"})
	}

	if !validEmail(client.Email) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid email format"})
	}

	if len(client.Telephone) < 7 || !isNumeric(client.Telephone) {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Phone number must be numeric and at least 7 digits long"})
	}

	currentYear := time.Now().Year()
	calculatedAge := currentYear - client.BirthDay.Year()
	if client.BirthDay.After(time.Now().AddDate(-client.Age, 0, 0)) {
		calculatedAge--
	}
	if client.Age != calculatedAge {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Age does not match birth date"})
	}

	if err := config.DB.Create(&client).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, client)
}

// UpdateClient actualiza un cliente
// @Summary Actualizar cliente
// @Description Actualiza la información de un cliente existente
// @Tags Clientes
// @Param id path int true "ID del Cliente"
// @Param client body models.Client true "Información actualizada del Cliente"
// @Produce json
// @Success 200 {object} models.Client "Cliente actualizado exitosamente"
// @Router /api/v1/clients/{id} [put]
func UpdateClient(c echo.Context) error {
	id := c.Param("id")
	var client models.Client
	if err := config.DB.First(&client, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
	}

	if err := c.Bind(&client); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "Invalid request data"})
	}

	if err := config.DB.Save(&client).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, client)
}

// DeleteClient elimina un cliente por ID
// @Summary Eliminar cliente
// @Description Elimina un cliente específico usando su ID
// @Tags Clientes
// @Param id path int true "ID del Cliente"
// @Success 204 "Cliente eliminado exitosamente"
// @Router /api/v1/clients/{id} [delete]
func DeleteClient(c echo.Context) error {
	id := c.Param("id")
	if err := config.DB.Delete(&models.Client{}, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "Client not found"})
	}
	return c.NoContent(http.StatusNoContent)
}

// GetClientKPI calcula el promedio de edad y la desviación estándar de edad de los clientes
// @Summary KPI de clientes
// @Description Calcula el promedio y la desviación estándar de edad de los clientes
// @Tags Clientes
// @Produce json
// @Success 200 {object} ClientKPI "KPI de clientes calculado"
// @Router /api/v1/clients/kpi [get]
func GetClientKPI(c echo.Context) error {
	var clients []struct {
		BirthDay time.Time
	}

	if err := config.DB.Model(&models.Client{}).Select("birth_day").Find(&clients).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"error": "Failed to retrieve clients"})
	}

	var ages []float64
	currentYear := time.Now().Year()
	for _, client := range clients {
		age := float64(currentYear - client.BirthDay.Year())
		ages = append(ages, age)
	}

	var sum float64
	for _, age := range ages {
		sum += age
	}
	averageAge := sum / float64(len(ages))

	var varianceSum float64
	for _, age := range ages {
		varianceSum += math.Pow(age-averageAge, 2)
	}
	ageStandardDeviation := math.Sqrt(varianceSum / float64(len(ages)))

	kpi := ClientKPI{
		AverageAge:           averageAge,
		AgeStandardDeviation: ageStandardDeviation,
	}
	return c.JSON(http.StatusOK, kpi)
}

func isNumeric(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}
