// handlers/user_handler.go
package handlers

import (
	"golangApp/config"
	"golangApp/models"
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// GetUser obtiene un usuario por ID
// @Summary Obtener usuario por ID
// @Description Recupera un usuario específico usando su ID
// @Tags Usuarios
// @Param id path int true "ID del Usuario"
// @Produce json
// @Success 200 {object} models.User "Detalles del usuario"
// @Router /api/v1/users/{id} [get]
func GetUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	// Mengambil user beserta grup yang diassign menggunakan Preload
	if err := config.DB.Preload("Groups").First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// GetAllUsers obtiene todos los usuarios
// @Summary Obtener todos los usuarios
// @Description Recupera una lista de todos los usuarios registrados
// @Tags Usuarios
// @Produce json
// @Success 200 {array} models.User "Lista de usuarios"
// @Router /api/v1/users [get]
func GetAllUsers(c echo.Context) error {
	var users []models.User

	// Mengambil semua user beserta grup yang diassign menggunakan Preload
	if err := config.DB.Preload("Groups").Find(&users).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to retrieve users",
		})
	}

	return c.JSON(http.StatusOK, users)
}

// CreateUser crea un nuevo usuario
// @Summary Crear usuario
// @Description Crea un nuevo usuario con los datos proporcionados
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param user body models.User true "Información del Usuario"
// @Success 201 {object} models.User "Usuario creado exitosamente"
// @Router /api/v1/users [post]
func CreateUser(c echo.Context) error {
	var user models.User

	// Bind input ke struct user
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid input",
		})
	}

	// Hash password sebelum disimpan
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to hash password",
		})
	}
	user.Password = string(hashedPassword)

	// Menyimpan user baru ke database dengan GORM
	if err := config.DB.Create(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to create user",
		})
	}

	return c.JSON(http.StatusCreated, user)
}

// UpdateUser actualiza un usuario por ID
// @Summary Actualizar usuario
// @Description Actualiza los datos de un usuario específico
// @Tags Usuarios
// @Accept json
// @Produce json
// @Param id path int true "ID del Usuario"
// @Param user body models.User true "Información del Usuario"
// @Success 200 {object} models.User "Usuario actualizado exitosamente"
// @Router /api/v1/users/{id} [put]
func UpdateUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	// Mencari user berdasarkan ID dengan GORM
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
		})
	}

	// Bind input ke struct user (hanya field yang ingin diupdate)
	if err := c.Bind(&user); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid input",
		})
	}

	// Mengupdate user di database dengan GORM
	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to update user",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// EnableUser habilita un usuario
// @Summary Habilitar usuario
// @Description Cambia el estado de un usuario a habilitado
// @Tags Usuarios
// @Param id path int true "ID del Usuario"
// @Produce json
// @Success 200 {object} models.User "Usuario habilitado"
// @Router /api/v1/users/{id}/enable [put]
func EnableUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	// Mencari user berdasarkan ID
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
		})
	}

	// Mengubah is_enabled menjadi true
	user.IsEnabled = true
	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to enable user",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// DisableUser deshabilita un usuario
// @Summary Deshabilitar usuario
// @Description Cambia el estado de un usuario a deshabilitado
// @Tags Usuarios
// @Param id path int true "ID del Usuario"
// @Produce json
// @Success 200 {object} models.User "Usuario deshabilitado"
// @Router /api/v1/users/{id}/disable [put]
func DisableUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	// Mencari user berdasarkan ID
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
		})
	}

	// Mengubah is_enabled menjadi false
	user.IsEnabled = false
	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to disable user",
		})
	}

	return c.JSON(http.StatusOK, user)
}

// DeleteUser elimina un usuario
// @Summary Eliminar usuario
// @Description Elimina un usuario específico usando su ID
// @Tags Usuarios
// @Param id path int true "ID del Usuario"
// @Success 204 "Usuario eliminado exitosamente"
// @Router /api/v1/users/{id} [delete]
func DeleteUser(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	// Mencari user berdasarkan ID
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
		})
	}

	// Menghapus user dari database
	if err := config.DB.Delete(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to delete user",
		})
	}

	return c.NoContent(http.StatusNoContent)
}

// ResetPassword restablece la contraseña de un usuario
// @Summary Restablecer contraseña
// @Description Cambia la contraseña de un usuario a una nueva
// @Tags Usuarios
// @Param id path int true "ID del Usuario"
// @Param new_password body string true "Nueva contraseña"
// @Success 200 "Contraseña restablecida exitosamente"
// @Router /api/v1/users/{id}/reset_password [put]
func ResetPassword(c echo.Context) error {
	id := c.Param("id")
	var user models.User

	// Mencari user berdasarkan ID
	if err := config.DB.First(&user, id).Error; err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{
			"message": "User not found",
		})
	}

	// Mendapatkan password baru dari request body
	type ResetPasswordRequest struct {
		NewPassword string `json:"new_password"`
	}

	var req ResetPasswordRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"message": "Invalid input",
		})
	}

	// Hash password baru
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to hash password",
		})
	}

	// Mengupdate password user di database
	user.Password = string(hashedPassword)
	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to reset password",
		})
	}

	return c.NoContent(http.StatusOK)
}