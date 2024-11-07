package handlers

import (
	"golangApp/config"
	"golangApp/models"
	"net/http"
	"time"

	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
)

// HandleLogin autentica al usuario
// @Summary Autenticar usuario
// @Description Autentica un usuario y crea una sesión
// @Tags Autenticación
// @Accept x-www-form-urlencoded
// @Produce json
// @Param username formData string true "Nombre de usuario"
// @Param password formData string true "Contraseña"
// @Success 200 {string} string "Inicio de sesión exitoso"
// @Router /login [post]
func HandleLogin(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	// Buscar el usuario en la base de datos
	var user models.User
	if err := config.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Invalid username or password",
		})
	}

	// Verificar la contraseña
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"error": "Invalid username or password",
		})
	}

	// Actualizar el tiempo del último inicio de sesión
	user.LastLogin = time.Now()
	if err := config.DB.Save(&user).Error; err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to update last login time",
		})
	}

	// Configurar la sesión
	sess, err := session.Get("session", c)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to get session",
		})
	}

	sess.Values["username"] = user.Username
	sess.Values["userID"] = user.ID
	if err := sess.Save(c.Request(), c.Response()); err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{
			"message": "Failed to save session",
		})
	}

	// Devolver el hash de la contraseña en la respuesta
	return c.JSON(http.StatusOK, echo.Map{
		"hash": user.Password,
	})
}
