package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID        string     `gorm:"type:varchar(36);primaryKey" json:"id"`
	CreatedBy string     `gorm:"type:varchar(36)" json:"created_by"`
	UpdatedBy string     `gorm:"type:varchar(36)" json:"updated_by"`
	CreatedAt time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time `json:"updated_at"`
	IsActive  bool       `gorm:"default:true" json:"is_active"`
}

func (base *BaseModel) BeforeCreate(tx *gorm.DB) error {
	if base.ID == "" {
		base.ID = uuid.New().String()
	}
	return nil
}

// User model
type User struct {
	BaseModel
	FirstName     string         `gorm:"type:varchar(105);not null" json:"first_name"`
	LastName      string         `gorm:"type:varchar(105);not null" json:"last_name"`
	Email         string         `gorm:"type:varchar(100)" json:"email"`
	Phone         string         `gorm:"type:varchar(12)" json:"phone"`
	Username      string         `gorm:"type:varchar(50);unique;index;comment:tên tài khoản" json:"username"`
	Gender        bool           `gorm:"default:true" json:"gender"`
	Password      string         `gorm:"type:varchar(100);not null" json:"-"`
	AvatarURL     string         `gorm:"type:varchar(255)" json:"avatar_url"`
	IsAdmin       bool           `gorm:"default:false" json:"is_admin"`
	Roles         []Role         `gorm:"many2many:user_roles" json:"roles,omitempty"`
	Pets          []Pet          `gorm:"foreignKey:UserID" json:"pets,omitempty"`
	LoginHistory  []LoginHistory `gorm:"foreignKey:UserID" json:"-"`
	Appointments  []Appointment  `gorm:"foreignKey:UserID" json:"appointments,omitempty"`
}

func (User) TableName() string {
	return "users"
}

// Role model
type Role struct {
	BaseModel
	Name        string       `gorm:"type:varchar(50);unique;index" json:"name"`
	Permissions []Permission `gorm:"many2many:role_permissions" json:"permissions,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

// Permission model
type Permission struct {
	BaseModel
	Name string `gorm:"type:varchar(50);unique;index" json:"name"`
}

func (Permission) TableName() string {
	return "permissions"
}

// RolePermission model
type RolePermission struct {
	BaseModel
	RoleID       string `gorm:"type:varchar(36)" json:"role_id"`
	PermissionID string `gorm:"type:varchar(36)" json:"permission_id"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

// UserRole model
type UserRole struct {
	BaseModel
	UserID string `gorm:"type:varchar(36)" json:"user_id"`
	RoleID string `gorm:"type:varchar(36)" json:"role_id"`
}

func (UserRole) TableName() string {
	return "user_roles"
}

// Pet model
type Pet struct {
	BaseModel
	Name         string         `gorm:"type:varchar(105);not null;comment:Tên" json:"name"`
	DateOfBirth  *time.Time     `json:"date_of_birth"`
	DateOfDeath  *time.Time     `json:"date_of_death"`
	Gender       bool           `gorm:"default:true" json:"gender"`
	Breed        string         `gorm:"type:varchar(50)" json:"breed"`
	Description  string         `gorm:"type:varchar(255)" json:"description"`
	AvtURL       string         `gorm:"type:varchar(255)" json:"avt_url"`
	Type         string         `gorm:"type:varchar(50)" json:"type"`
	UserID       string         `gorm:"type:varchar(36);not null" json:"user_id"`
	User         User           `gorm:"foreignKey:UserID" json:"user,omitempty"`
	Medias       []Media        `gorm:"foreignKey:PetID" json:"medias,omitempty"`
	LifeEvents   []PetLifeEvent `gorm:"foreignKey:PetID" json:"life_events,omitempty"`
	Comments     []Comment      `gorm:"foreignKey:PetID" json:"comments,omitempty"`
}

func (Pet) TableName() string {
	return "pets"
}

// Media model
type Media struct {
	BaseModel
	Type  string `gorm:"type:varchar(50)" json:"type"`
	Name  string `gorm:"type:varchar(105);not null" json:"name"`
	URL   string `gorm:"type:varchar(255)" json:"url"`
	PetID string `gorm:"type:varchar(36)" json:"pet_id"`
	Pet   Pet    `gorm:"foreignKey:PetID" json:"pet,omitempty"`
}

func (Media) TableName() string {
	return "medias"
}

// Comment model
type Comment struct {
	BaseModel
	Content  string `gorm:"type:text" json:"content"`
	PetID    string `gorm:"type:varchar(36)" json:"pet_id"`
	ParentID string `gorm:"type:varchar(36)" json:"parent_id"`
	Pet      Pet    `gorm:"foreignKey:PetID" json:"pet,omitempty"`
}

func (Comment) TableName() string {
	return "comments"
}

// PetLifeEvent model
type PetLifeEvent struct {
	BaseModel
	PetID    string    `gorm:"type:varchar(36);not null" json:"pet_id"`
	Title    string    `gorm:"type:varchar(100);not null" json:"title"`
	Date     time.Time `gorm:"not null" json:"date"`
	Location string    `gorm:"type:varchar(255)" json:"location"`
	Story    string    `gorm:"type:varchar(255)" json:"story"`
	Pet      Pet       `gorm:"foreignKey:PetID" json:"pet,omitempty"`
}

func (PetLifeEvent) TableName() string {
	return "pet_life_events"
}

// Service model
type Service struct {
	BaseModel
	Code        string `gorm:"type:varchar(50);not null" json:"code"`
	Name        string `gorm:"type:varchar(100);not null" json:"name"`
	Price       int    `gorm:"not null" json:"price"`
	Duration    string `gorm:"type:varchar(100)" json:"duration"`
	Description string `gorm:"type:varchar(255)" json:"description"`
}

func (Service) TableName() string {
	return "services"
}

// Appointment model
type Appointment struct {
	BaseModel
	Code               string              `gorm:"type:varchar(50);not null" json:"code"`
	Status             string              `gorm:"type:varchar(20);default:PENDING" json:"status"`
	StartTime          *time.Time          `json:"start_time"`
	TotalPrice         int                 `gorm:"not null;default:0" json:"total_price"`
	IsOnline           bool                `gorm:"default:true" json:"is_online"`
	UserID             string              `gorm:"type:varchar(36);not null" json:"user_id"`
	User               User                `gorm:"foreignKey:UserID" json:"user,omitempty"`
	AppointmentDetails []AppointmentDetail `gorm:"foreignKey:AppointmentID" json:"appointment_details,omitempty"`
	Payments           []Payment           `gorm:"foreignKey:AppointmentID" json:"payments,omitempty"`
}

func (Appointment) TableName() string {
	return "appointments"
}

// AppointmentDetail model
type AppointmentDetail struct {
	BaseModel
	AppointmentID string      `gorm:"type:varchar(36);not null" json:"appointment_id"`
	ServiceID     string      `gorm:"type:varchar(36);not null" json:"service_id"`
	StartTime     *time.Time  `json:"start_time"`
	EndTime       *time.Time  `json:"end_time"`
	Status        string      `gorm:"type:varchar(20);default:PENDING" json:"status"`
	DiscountCode  string      `gorm:"type:varchar(100)" json:"discount_code"`
	DiscountPrice int         `json:"discount_price"`
	UnitPrice     int         `gorm:"not null" json:"unit_price"`
	Price         int         `json:"price"`
	Quantity      int         `json:"quantity"`
	Appointment   Appointment `gorm:"foreignKey:AppointmentID" json:"appointment,omitempty"`
}

func (AppointmentDetail) TableName() string {
	return "appointment_details"
}

// Payment model
type Payment struct {
	BaseModel
	AppointmentID string      `gorm:"type:varchar(36);not null" json:"appointment_id"`
	PaymentDate   *time.Time  `json:"payment_date"`
	Amount        int         `gorm:"not null" json:"amount"`
	Method        string      `gorm:"type:varchar(100);not null" json:"method"`
	Status        string      `gorm:"type:varchar(20);not null" json:"status"`
	Appointment   Appointment `gorm:"foreignKey:AppointmentID" json:"appointment,omitempty"`
}

func (Payment) TableName() string {
	return "payments"
}

// LoginHistory model
type LoginHistory struct {
	BaseModel
	UserID       string `gorm:"type:varchar(36)" json:"user_id"`
	JTI          string `gorm:"type:varchar(36)" json:"jti"`
	RefreshToken string `gorm:"type:text" json:"refresh_token"`
	AccessToken  string `gorm:"type:text" json:"access_token"`
}

func (LoginHistory) TableName() string {
	return "login_history"
}

// TokenBlacklist model
type TokenBlacklist struct {
	BaseModel
	JTI string `gorm:"type:varchar(36)" json:"jti"`
}

func (TokenBlacklist) TableName() string {
	return "token_blacklist"
}
