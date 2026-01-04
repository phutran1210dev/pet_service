package container

import (
	"pet-service/handler"
	"pet-service/repository"
	"pet-service/service"

	"gorm.io/gorm"
)

// Container holds all application dependencies
type Container struct {
	DB       *gorm.DB
	Repos    *Repositories
	Services *Services
	Handlers *Handlers
}

// Repositories holds all repository instances
type Repositories struct {
	User repository.IUserRepository
	Pet  repository.IPetRepository
}

// Services holds all service instances
type Services struct {
	User        service.IUserService
	Pet         service.IPetService
	Appointment service.IAppointmentService
}

// Handlers holds all handler instances
type Handlers struct {
	User        *handler.UserHandler
	Pet         *handler.PetHandler
	Appointment *handler.AppointmentHandler
}

// NewContainer creates and wires up all dependencies
func NewContainer(db *gorm.DB) *Container {
	// Initialize repositories
	repos := &Repositories{
		User: repository.NewUserRepository(db),
		Pet:  repository.NewPetRepository(db),
	}

	// Initialize services with repository interfaces
	services := &Services{
		User:        service.NewUserService(repos.User),
		Pet:         service.NewPetService(repos.Pet),
		Appointment: service.NewAppointmentService(db),
	}

	// Initialize handlers with service interfaces
	handlers := &Handlers{
		User:        handler.NewUserHandler(services.User),
		Pet:         handler.NewPetHandler(services.Pet, db),
		Appointment: handler.NewAppointmentHandler(services.Appointment),
	}

	return &Container{
		DB:       db,
		Repos:    repos,
		Services: services,
		Handlers: handlers,
	}
}
