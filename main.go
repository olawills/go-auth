package main

import (
	"auth_api_with_Go/controllers"
	database "auth_api_with_Go/db"
	"auth_api_with_Go/middlewares"
	"auth_api_with_Go/model"
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func main() {

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
	}

	database.NewConnection(dbConfig)
	database.Migrate()

	// Initialize Router
	router := initRouter()
	router.Run(":8080")

}

func initRouter() *gin.Engine {
	router := gin.Default()
	api := router.Group("/api/v1")
	{
		api.POST("/token", controllers.GenerateToken)
		api.POST("/auth/register", controllers.RegisterUser)
		api.POST("/auth/login", func(c *gin.Context) {
			var login model.Login
			if err := c.ShouldBindJSON(&login); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
			controllers.LoginUser(&login, c)
		})
		secured := api.Group("/secured").Use(middlewares.Auth())
		{
			secured.GET("/ping", controllers.Ping)
		}
	}
	return router
}
