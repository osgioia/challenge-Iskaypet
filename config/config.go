package config

import (
	"golangApp/models"
	"log"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB
var isTestEnv bool

// DBInterface es la interfaz para abstraer el uso de la base de datos en las pruebas
type DBInterface interface {
	Create(value interface{}) *gorm.DB
	Find(out interface{}, where ...interface{}) *gorm.DB
	Model(value interface{}) *gorm.DB
}

// InitDB inicializa la conexión con la base de datos real
func InitDB() {
	// Determinar si estamos en un entorno de prueba o producción
	isTestEnv = false

	// Para producción (SQLite en este caso)
	dbPath := "./database/app.db"
	var err error
	DB, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to SQLite database:", err)
	}

	// Auto migrar tablas
	DB.AutoMigrate(&models.User{}, &models.Group{}, &models.Client{})

	log.Println("Connected to SQLite database successfully")

	// Sembrar datos iniciales solo en producción
	if !isTestEnv {
		seedData()
	}
}

// SetupTestDB configura una base de datos en memoria para las pruebas
func SetupTestDB() {
	// Establecer que estamos en entorno de pruebas
	isTestEnv = true

	// Configurar la base de datos en memoria para pruebas
	var err error
	DB, err = gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		panic("failed to connect to the database")
	}

	// Auto migrar tablas
	DB.AutoMigrate(&models.Client{}, &models.User{}, &models.Group{})
}

// seedData crea datos iniciales en la base de datos
func seedData() {
	seedGroupsAndUsers()
	seedClients()
}

func seedClients() {
	var clientCount int64
	DB.Model(&models.Client{}).Count(&clientCount)
	if clientCount == 0 {
		// Datos de clientes iniciales
		clients := []models.Client{
			{Name: "John", LastName: "Doe", Email: "johndoe@example.com", BirthDay: time.Date(1985, time.May, 15, 0, 0, 0, 0, time.UTC)},
			{Name: "Jane", LastName: "Smith", Email: "janesmith@example.com", BirthDay: time.Date(1990, time.June, 20, 0, 0, 0, 0, time.UTC)},
			{Name: "Alice", LastName: "Johnson", Email: "alicej@example.com", BirthDay: time.Date(1978, time.December, 5, 0, 0, 0, 0, time.UTC)},
		}

		for _, client := range clients {
			DB.Create(&client)
		}
		log.Println("Seeded initial clients")
	}
}

func seedGroupsAndUsers() {
	var groupCount int64
	DB.Model(&models.Group{}).Count(&groupCount)
	if groupCount == 0 {
		adminGroup := models.Group{Name: "Admin", Description: "Administrator group"}
		userGroup := models.Group{Name: "User", Description: "Regular user group"}
		DB.Create(&adminGroup)
		DB.Create(&userGroup)
		log.Println("Seeded initial groups")

		users := []models.User{
			{Username: "admin", Email: "admin@example.com", Password: hashPassword("admin"), FirstName: "Admin", LastName: "User", Groups: []models.Group{adminGroup}},
			{Username: "user", Email: "user@example.com", Password: hashPassword("user"), FirstName: "Regular", LastName: "User", Groups: []models.Group{userGroup}},
		}
		for _, user := range users {
			DB.Create(&user)
		}
		log.Println("Seeded initial users with group assignments")
	}
}

func hashPassword(password string) string {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(hashedPassword)
}
