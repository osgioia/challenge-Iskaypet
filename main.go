package main

import (
	"golangApp/config"
	"golangApp/handlers"
	"golangApp/middlewares"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/swaggo/echo-swagger"
	_ "golangApp/docs"
)

// @title Swagger Example API
// @version 1.0
// @description This is a sample server for using Swagger with Echo.
// @host localhost:8080
// @BasePath /
func main() {
	config.InitDB()

	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Setup session middleware
	e.Use(session.Middleware(sessions.NewCookieStore([]byte("secret-key"))))

	// Routes that don't require authentication
	e.POST("/login", handlers.HandleLogin)

	// Group of routes that require authentication
	auth := e.Group("/api/v1")

	// Apply the Basic Auth Middleware only to specific routes
	auth.Use(middleware.BasicAuth(middlewares.BasicAuthMiddleware))

	// Define routes without authentication
	auth.GET("/clients/:id", handlers.GetClient)        // Obtiene un cliente por ID
	auth.GET("/clients", handlers.GetAll)               // Obtiene todos los clientes
	auth.GET("/clients/kpi", handlers.GetClientKPI)     // KPI de clientes
	auth.POST("/clients", handlers.CreateClient)        // Crea un nuevo cliente
	auth.PUT("/clients/:id", handlers.UpdateClient)     // Actualiza un cliente por ID
	auth.DELETE("/clients/:id", handlers.DeleteClient)  // Elimina un cliente por ID
	auth.GET("/users/:id", handlers.GetUser)
	auth.GET("/users", handlers.GetAllUsers)
	auth.GET("/groups/:id", handlers.GetGroup)
	auth.GET("/groups", handlers.GetAllGroups)
	auth.POST("/users", handlers.CreateUser)
	auth.POST("/groups", handlers.CreateGroup)
	auth.PUT("/users/:id", handlers.UpdateUser)
	auth.PUT("/users/:id/enable", handlers.EnableUser)
	auth.PUT("/users/:id/disable", handlers.DisableUser)
	auth.DELETE("/users/:id", handlers.DeleteUser)
	auth.POST("/users/:id/groups/:group_id", handlers.AssignGroup)
	auth.DELETE("/users/:id/groups/:group_id", handlers.RemoveAssignGroup)
	auth.DELETE("/groups/:group_id", handlers.RemoveGroup)
	auth.PUT("/users/:id/reset_password", handlers.ResetPassword)

	e.GET("/swagger/*", echoSwagger.WrapHandler)
	// Start server
	e.Logger.Fatal(e.Start(":8080"))
}