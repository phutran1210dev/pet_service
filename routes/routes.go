package routes

import (
"pet-service/container"
"pet-service/middleware"

"github.com/gin-gonic/gin"
"gorm.io/gorm"
)

func SetupRoutes(router *gin.Engine, db *gorm.DB) {
	// Initialize dependency injection container
	c := container.NewContainer(db)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		// Authentication routes (no auth required)
		auth := v1.Group("")
		{
			auth.POST("/login", c.Handlers.User.Login)
			auth.POST("/user", c.Handlers.User.Register)
		}

		// Protected user routes
		users := v1.Group("")
		users.Use(middleware.AuthMiddleware())
		{
			users.GET("/me", c.Handlers.User.GetMe)
			users.POST("/logout", c.Handlers.User.Logout)
			users.GET("/users", c.Handlers.User.GetUsers)
			users.PATCH("/users/change-password", c.Handlers.User.ChangePassword)
		}

		// Comment routes (protected)
		comments := v1.Group("")
		comments.Use(middleware.AuthMiddleware())
		{
			comments.POST("/post/:pet_id/comment", c.Handlers.User.CreateComment)
			comments.PATCH("/post/:pet_id/comment/:comment_id", c.Handlers.User.EditComment)
			comments.GET("/post/:pet_id/comments", c.Handlers.User.GetComments)
		}

		// Pet routes (protected)
		pets := v1.Group("")
		pets.Use(middleware.AuthMiddleware())
		{
			pets.POST("/pet", c.Handlers.Pet.CreatePet)
			pets.GET("/pets", c.Handlers.Pet.GetPets)
			pets.GET("/pet/:id", c.Handlers.Pet.GetPetDetail)
			pets.POST("/pet/life-event", c.Handlers.Pet.CreatePetLifeEvent)
			pets.POST("/pet/:pet_id/images", c.Handlers.Pet.UploadAvatar)
			pets.POST("/pet/:pet_id/gallery", c.Handlers.Pet.UploadGallery)
		}

		// Appointment routes (protected)
		appointments := v1.Group("")
		appointments.Use(middleware.AuthMiddleware())
		{
			appointments.POST("/appointment/register", c.Handlers.Appointment.RegisterAppointment)
		}
	}
}
