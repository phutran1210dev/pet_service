package main

import (
	"log"
	"pet-service/config"
	"pet-service/database"
	_ "pet-service/docs" // Swagger docs
	"pet-service/middleware"
	"pet-service/routes"
	"pet-service/scheduler"
	"pet-service/storage"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title           Pet Service API
// @version         1.0
// @description     Pet Service API with JWT authentication
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.email  support@petservice.com

// @license.name  MIT
// @license.url   https://opensource.org/licenses/MIT

// @host      localhost:8001
// @BasePath  /api/v1

// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	// Load configuration
	config.LoadConfig()

	// Connect to database
	database.ConnectDatabase()

	// Initialize MinIO
	if err := storage.InitMinio(); err != nil {
		log.Fatalf("Failed to initialize MinIO: %v", err)
	}

	// Initialize scheduler
	scheduler.InitScheduler()

	// Setup Gin
	if !config.AppConfig.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.Default()

	// Middleware
	router.Use(middleware.LoggingMiddleware())
	router.Use(func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	// Setup routes
	routes.SetupRoutes(router, database.GetDB())

	// Swagger documentation
	log.Printf("Debug mode: %v", config.AppConfig.Debug)
	if config.AppConfig.Debug {
		log.Println("Registering Swagger routes...")
		router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		router.GET("/docs", func(c *gin.Context) {
			c.Redirect(301, "/swagger/index.html")
		})
		router.GET("/", func(c *gin.Context) {
			c.Redirect(301, "/swagger/index.html")
		})
	}

	// Health check
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
			"app":    config.AppConfig.ProjectName,
		})
	})

	// Start server
	port := ":" + config.AppConfig.ServerPort
	log.Printf("Server starting on port %s", port)
	if err := router.Run(port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
