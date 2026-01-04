package service

import (
	"io"
	"pet-service/dto"
	"pet-service/middleware"

	"gorm.io/gorm"
)

// IUserService defines the interface for user business logic operations
type IUserService interface {
	Register(req dto.UserRegisterRequest) (*dto.UserResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	GetMe(userInfo middleware.UserInfo) (*dto.UserResponse, error)
	Logout(userInfo middleware.UserInfo) (*dto.MessageResponse, error)
	GetUsers() ([]dto.UserResponse, error)
	ChangePassword(userInfo middleware.UserInfo, req dto.ChangePasswordRequest) (*dto.MessageResponse, error)
	
	// Comment operations
	CreateComment(userInfo middleware.UserInfo, petID string, req dto.CommentRequest) (*dto.CommentResponse, error)
	EditComment(userInfo middleware.UserInfo, petID, commentID string, req dto.CommentRequest) (*dto.CommentResponse, error)
	GetCommentsByPetID(petID string) ([]dto.CommentResponse, error)
}

// IPetService defines the interface for pet business logic operations
type IPetService interface {
	CreatePet(userInfo middleware.UserInfo, req dto.PetCreateRequest) (*dto.PetResponse, error)
	GetPets(db *gorm.DB, page, pageSize int, search, name string) (*dto.PaginationResponse, error)
	GetPetDetail(petID string) (*dto.PetDetailResponse, error)
	CreatePetLifeEvent(userInfo middleware.UserInfo, req dto.PetLifeEventRequest) (*dto.PetLifeEventResponse, error)
	UploadAvatar(petID string, fileData []byte, contentType string) (*dto.MediaResponse, error)
	UploadGallery(petID string, files []io.Reader, fileNames []string, contentTypes []string) ([]dto.MediaResponse, error)
}

// IAppointmentService defines the interface for appointment business logic operations
type IAppointmentService interface {
	RegisterAppointment(userInfo middleware.UserInfo, req dto.AppointmentRequest) (*dto.AppointmentResponse, error)
}
