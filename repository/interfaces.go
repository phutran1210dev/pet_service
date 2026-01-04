package repository

import (
	"pet-service/models"

	"gorm.io/gorm"
)

// UserRepository defines the interface for user data access operations
type IUserRepository interface {
	// User operations
	CreateUser(user *models.User) error
	GetUserByEmail(email string) (*models.User, error)
	GetUserByID(id string) (*models.User, error)
	GetUsers() ([]models.User, error)
	UpdateUser(user *models.User) error

	// Role operations
	GetRoleByName(name string) (*models.Role, error)
	CreateUserRole(userRole *models.UserRole) error
	GetRolesPermissionsByUserID(userID string) ([]map[string]interface{}, error)
	GetPermissionsByUserID(userID string) ([]string, error)

	// Login history operations
	CreateLoginHistory(history *models.LoginHistory) error
	GetLoginHistoryByJTI(jti string) (*models.LoginHistory, error)
	UpdateLoginHistory(history *models.LoginHistory) error

	// Token blacklist operations
	CreateTokenBlacklist(token *models.TokenBlacklist) error
	IsTokenBlacklisted(jti string) bool

	// Comment operations
	CreateComment(comment *models.Comment) error
	GetCommentByID(id string) (*models.Comment, error)
	UpdateComment(comment *models.Comment) error
	GetCommentsByPetID(petID string) ([]map[string]interface{}, error)
}

// PetRepository defines the interface for pet data access operations
type IPetRepository interface {
	// Pet operations
	CreatePet(pet *models.Pet) error
	GetPetByID(id string) (*models.Pet, error)
	GetPets(query *gorm.DB) *gorm.DB
	UpdatePet(pet *models.Pet) error
	GetPetDetail(petID string) ([]map[string]interface{}, error)

	// Pet life event operations
	CreateLifeEvent(event *models.PetLifeEvent) error

	// Media operations
	CreateMediaBatch(medias []models.Media) error
}
