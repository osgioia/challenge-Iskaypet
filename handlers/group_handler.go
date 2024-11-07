// handlers/group_handler.go
package handlers

import (
	"golangApp/config"
	"golangApp/models"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

// GetGroup obtiene un grupo por ID
// @Summary Obtiene un grupo
// @Description Recupera un grupo específico basado en el ID proporcionado
// @Tags Group
// @Param id path int true "ID del grupo"
// @Success 200 {object} models.Group
// @Router /groups/{id} [get]
func GetGroup(c echo.Context) error {
	id := c.Param("id")
	var group models.Group

	// Mencari group berdasarkan ID
	if err := config.DB.First(&group, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Group not found",
		})
	}

	return c.JSON(http.StatusOK, group)
}

// GetAllGroups obtiene todos los grupos
// @Summary Obtiene todos los grupos
// @Description Recupera todos los grupos de la base de datos
// @Tags Group
// @Success 200 {array} models.Group
// @Router /groups [get]
func GetAllGroups(c echo.Context) error {
	var groups []models.Group

	// Mengambil semua grup
	if err := config.DB.Find(&groups).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to retrieve groups",
		})
	}

	return c.JSON(http.StatusOK, groups)
}

// CreateGroup crea un nuevo grupo
// @Summary Crea un grupo
// @Description Crea un nuevo grupo en la base de datos
// @Tags Group
// @Param group body models.Group true "Datos del grupo"
// @Success 201 {object} models.Group
// @Router /groups [post]
func CreateGroup(c echo.Context) error {
	var group models.Group

	// Bind input JSON to the group struct
	if err := c.Bind(&group); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid input",
		})
	}

	// Validate the required fields
	if group.Name == "" {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Group name is required",
		})
	}

	// Create the new group in the database
	if err := config.DB.Create(&group).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to create group",
		})
	}

	return c.JSON(http.StatusCreated, group)
}

// AssignGroup asigna un grupo a un usuario
// @Summary Asigna un grupo a un usuario
// @Description Asigna un grupo existente a un usuario basado en el ID del usuario y del grupo
// @Tags Group
// @Param id path int true "ID del usuario"
// @Param group_id path int true "ID del grupo"
// @Success 200 "Group assigned"
// @Router /users/{id}/groups/{group_id} [put]
func AssignGroup(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid user ID",
		})
	}

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid group ID",
		})
	}

	var user models.User
	var group models.Group

	// Mencari user berdasarkan ID
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
		})
	}

	// Mencari group berdasarkan ID
	if err := config.DB.First(&group, groupID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Group not found",
		})
	}

	// Menetapkan group ke user
	if err := config.DB.Model(&user).Association("Groups").Append(&group); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to assign group",
		})
	}

	return c.NoContent(http.StatusOK)
}

// RemoveAssignGroup elimina la asignación de un grupo de un usuario
// @Summary Elimina la asignación de un grupo a un usuario
// @Description Elimina la relación de un grupo asignado a un usuario basado en sus IDs
// @Tags Group
// @Param id path int true "ID del usuario"
// @Param group_id path int true "ID del grupo"
// @Success 200 "Group unassigned"
// @Router /users/{id}/groups/{group_id} [delete]
func RemoveAssignGroup(c echo.Context) error {
	userID, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid user ID",
		})
	}

	groupID, err := strconv.Atoi(c.Param("group_id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid group ID",
		})
	}

	var user models.User
	var group models.Group

	// Cek apakah user dan group ada di database
	if err := config.DB.First(&user, userID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
		})
	}

	if err := config.DB.First(&group, groupID).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Group not found",
		})
	}

	// Menghapus relasi group dari user
	if err := config.DB.Model(&user).Association("Groups").Delete(&group); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to remove group assignment",
		})
	}

	return c.NoContent(http.StatusOK)
}

// RemoveGroup elimina un grupo
// @Summary Elimina un grupo
// @Description Elimina un grupo de la base de datos basado en el ID proporcionado
// @Tags Group
// @Param group_id path int true "ID del grupo"
// @Success 204 "Group deleted"
// @Router /groups/{group_id} [delete]
func RemoveGroup(c echo.Context) error {
	id := c.Param("group_id")
	var group models.Group

	// Mencari group berdasarkan ID
	if err := config.DB.First(&group, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "Group not found",
		})
	}

	// Menghapus grup dari database
	if err := config.DB.Delete(&group).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to delete group",
		})
	}

	return c.NoContent(http.StatusNoContent)
}