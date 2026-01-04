package repository

import (
	"pet-service/models"

	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{DB: db}
}

func (r *UserRepository) CreateUser(user *models.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) GetRoleByName(name string) (*models.Role, error) {
	var role models.Role
	err := r.DB.Where("name = ? AND is_active = ?", name, true).First(&role).Error
	if err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *UserRepository) CreateUserRole(userRole *models.UserRole) error {
	return r.DB.Create(userRole).Error
}

func (r *UserRepository) GetUserByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("email = ? AND is_active = ?", email, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUserByID(id string) (*models.User, error) {
	var user models.User
	err := r.DB.Where("id = ? AND is_active = ?", id, true).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetUsers() ([]models.User, error) {
	var users []models.User
	err := r.DB.Where("is_active = ?", true).Find(&users).Error
	return users, err
}

func (r *UserRepository) UpdateUser(user *models.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) GetRolesPermissionsByUserID(userID string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.DB.Table("roles").
		Select("roles.name as role_name, roles.id as role_id, permissions.name as permission_name, permissions.id as permission_id").
		Joins("JOIN user_roles ON user_roles.role_id = roles.id").
		Joins("JOIN role_permissions ON role_permissions.role_id = roles.id").
		Joins("JOIN permissions ON permissions.id = role_permissions.permission_id").
		Where("user_roles.user_id = ?", userID).
		Scan(&results).Error

	return results, err
}

func (r *UserRepository) GetPermissionsByUserID(userID string) ([]string, error) {
	var permissions []string

	err := r.DB.Table("permissions").
		Select("permissions.name").
		Joins("JOIN role_permissions ON role_permissions.permission_id = permissions.id").
		Joins("JOIN user_roles ON user_roles.role_id = role_permissions.role_id").
		Where("user_roles.user_id = ?", userID).
		Pluck("name", &permissions).Error

	return permissions, err
}

// Login History
func (r *UserRepository) CreateLoginHistory(history *models.LoginHistory) error {
	return r.DB.Create(history).Error
}

func (r *UserRepository) GetLoginHistoryByJTI(jti string) (*models.LoginHistory, error) {
	var history models.LoginHistory
	err := r.DB.Where("jti = ? AND is_active = ?", jti, true).First(&history).Error
	if err != nil {
		return nil, err
	}
	return &history, nil
}

func (r *UserRepository) UpdateLoginHistory(history *models.LoginHistory) error {
	return r.DB.Save(history).Error
}

// Token Blacklist
func (r *UserRepository) CreateTokenBlacklist(token *models.TokenBlacklist) error {
	return r.DB.Create(token).Error
}

func (r *UserRepository) IsTokenBlacklisted(jti string) bool {
	var count int64
	r.DB.Model(&models.TokenBlacklist{}).Where("jti = ? AND is_active = ?", jti, true).Count(&count)
	return count > 0
}

// Comments
func (r *UserRepository) CreateComment(comment *models.Comment) error {
	return r.DB.Create(comment).Error
}

func (r *UserRepository) GetCommentByID(id string) (*models.Comment, error) {
	var comment models.Comment
	err := r.DB.Where("id = ? AND is_active = ?", id, true).First(&comment).Error
	if err != nil {
		return nil, err
	}
	return &comment, nil
}

func (r *UserRepository) UpdateComment(comment *models.Comment) error {
	return r.DB.Save(comment).Error
}

func (r *UserRepository) GetCommentsByPetID(petID string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.DB.Table("comments").
		Select("comments.id, comments.content, comments.created_at, comments.updated_at, comments.parent_id, users.last_name, users.first_name, users.id as user_id, users.avatar_url").
		Joins("JOIN users ON users.id = comments.created_by").
		Where("comments.pet_id = ?", petID).
		Scan(&results).Error

	return results, err
}
