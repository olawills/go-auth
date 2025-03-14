package main

import (
	"auth_api_with_Go/controllers"
	database "auth_api_with_Go/db"
	"auth_api_with_Go/middlewares"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)
// migrate -path . -database "postgresql://postgres@localhost:5432/auth_db?sslmode=disable" up 
func main() {
	// Initialize Database
	// Load configuration
	viper.SetConfigFile("config.json")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config: %v", err)
	}

	// Set up database
	dbConfig := database.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		User:     viper.GetString("db.user"),
		Password: viper.GetString("db.password"),
		DBName:   viper.GetString("db.name"),
		SSLMode:  viper.GetString("db.sslmode"),
	}

	db, err := database.NewConnection(dbConfig)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()

	// Initialize Router
	router := initRouter()
	// Start server
	port := viper.GetString("server.port")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s...", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func initRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/auth")
	{
		api.POST("/token", controllers.GenerateToken)
		api.POST("/user/register", controllers.RegisterUser)
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/ping", controllers.Ping)
		}
	}
	return router
}
