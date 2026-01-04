package dto

// Standard API Response structures following Golang best practices

// ErrorDetail represents a single validation error detail
type ErrorDetail struct {
	Field   string `json:"field" example:"email"`
	Message string `json:"message" example:"Email must be a valid email address"`
}

// ErrorResponse represents standard error response
type ErrorResponse struct {
	Code    string        `json:"code" example:"VALIDATION_ERROR"`
	Message string        `json:"message" example:"Validation failed"`
	Details []ErrorDetail `json:"details,omitempty"`
}

type MessageResponse struct {
	Message string `json:"message"`
}

type PaginationMeta struct {
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
}

type PaginationResponse struct {
	Data interface{}    `json:"data"`
	Meta PaginationMeta `json:"meta"`
}

// User DTOs
type UserRegisterRequest struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required,email"`
	Phone     string `json:"phone" binding:"required"`
	Gender    bool   `json:"gender"`
	Password  string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Expire       int64  `json:"expire"`
}

type UserResponse struct {
	ID          string   `json:"id"`
	FirstName   string   `json:"first_name"`
	LastName    string   `json:"last_name"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	IsAdmin     bool     `json:"is_admin"`
	Avatar      string   `json:"avatar,omitempty"`
	Roles       []string `json:"roles,omitempty"`
	Permissions []string `json:"permissions,omitempty"`
}

type ChangePasswordRequest struct {
	OldPassword   string `json:"old_password" binding:"required,min=6"`
	NewPassword   string `json:"new_password" binding:"required,min=6"`
	ReNewPassword string `json:"re_new_password" binding:"required,min=6"`
}

type CommentRequest struct {
	Content  string `json:"content" binding:"required"`
	ParentID string `json:"parent_id"`
}

type CommentResponse struct {
	ID        string `json:"id"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	ParentID  string `json:"parent_id"`
	UserID    string `json:"user_id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	AvatarURL string `json:"avatar_url"`
}

// Pet DTOs
type PetCreateRequest struct {
	Name        string `json:"name" binding:"required"`
	Gender      bool   `json:"gender"`
	DateOfBirth string `json:"date_of_birth" binding:"required"`
	DateOfDeath string `json:"date_of_death"`
	Breed       string `json:"breed"`
	Description string `json:"description"`
	Type        string `json:"type" binding:"required"`
}

type PetLifeEventRequest struct {
	PetID    string `json:"pet_id" binding:"required"`
	Title    string `json:"title" binding:"required"`
	Date     string `json:"date" binding:"required"`
	Location string `json:"location"`
	Story    string `json:"story"`
}

type PetResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Gender      bool   `json:"gender"`
	DateOfBirth string `json:"date_of_birth"`
	DateOfDeath string `json:"date_of_death"`
	Breed       string `json:"breed"`
	Description string `json:"description"`
	Type        string `json:"type"`
	AvtURL      string `json:"avt_url"`
}

type PetDetailResponse struct {
	ID          string               `json:"id"`
	Name        string               `json:"name"`
	Gender      bool                 `json:"gender"`
	DateOfBirth string               `json:"date_of_birth"`
	DateOfDeath string               `json:"date_of_death"`
	Breed       string               `json:"breed"`
	Description string               `json:"description"`
	Events      []PetLifeEventItem   `json:"events"`
	Medias      []MediaItem          `json:"medias"`
}

type PetLifeEventItem struct {
	ID       string `json:"id"`
	Title    string `json:"title"`
	Date     string `json:"date"`
	Location string `json:"location"`
	Story    string `json:"story"`
}

type MediaItem struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}

// Appointment DTOs
type AppointmentRequest struct {
	StartTime string `json:"start_time" binding:"required"`
	Message   string `json:"message"`
}

type AppointmentResponse struct {
	Code      string `json:"code"`
	StartTime string `json:"start_time"`
}

// Additional response DTOs for type safety
type PetLifeEventResponse struct {
	ID     string `json:"id"`
	PetID  string `json:"pet_id"`
	Title  string `json:"title"`
	Date   string `json:"date"`
}

type MediaResponse struct {
	ID  string `json:"id"`
	URL string `json:"url"`
}
