package service

import (
	"errors"
	"pet-service/dto"
	"pet-service/middleware"
	"pet-service/models"
	"pet-service/repository"
	"pet-service/utils"
	"strings"
	"time"
)

type userService struct {
	userRepo repository.IUserRepository
}

// NewUserService creates a new user service instance
func NewUserService(userRepo repository.IUserRepository) IUserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) Register(req dto.UserRegisterRequest) (*dto.UserResponse, error) {
	// Check if email exists
	existingUser, _ := s.userRepo.GetUserByEmail(req.Email)
	if existingUser != nil {
		return nil, errors.New(utils.EmailTaken)
	}

	// Generate username from email (part before @)
	username := strings.Split(req.Email, "@")[0]

	// Hash password
	hashedPassword, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user
	user := &models.User{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		Email:     req.Email,
		Phone:     req.Phone,
		Username:  username,
		Gender:    req.Gender,
		Password:  hashedPassword,
		IsAdmin:   false,
	}

	if err := s.userRepo.CreateUser(user); err != nil {
		return nil, err
	}

	// Assign default User role
	role, err := s.userRepo.GetRoleByName(utils.RoleUser)
	if err == nil && role != nil {
		userRole := &models.UserRole{
			UserID: user.ID,
			RoleID: role.ID,
		}
		userRole.CreatedBy = user.ID
		_ = s.userRepo.CreateUserRole(userRole)
	}

	// Load roles and permissions for the newly created user
	rolesPermissions, _ := s.userRepo.GetRolesPermissionsByUserID(user.ID)
	var roles []string
	var permissions []string

	for _, rp := range rolesPermissions {
		if roleName, ok := rp["role_name"].(string); ok {
			// Add role if not already in list
			found := false
			for _, r := range roles {
				if r == roleName {
					found = true
					break
				}
			}
			if !found {
				roles = append(roles, roleName)
			}
		}
		if permName, ok := rp["permission_name"].(string); ok {
			// Add permission if not already in list
			found := false
			for _, p := range permissions {
				if p == permName {
					found = true
					break
				}
			}
			if !found {
				permissions = append(permissions, permName)
			}
		}
	}

	return &dto.UserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Phone:       user.Phone,
		IsAdmin:     user.IsAdmin,
		Avatar:      user.AvatarURL,
		Roles:       roles,
		Permissions: permissions,
	}, nil
}

func (s *userService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Get user by email
	user, err := s.userRepo.GetUserByEmail(req.Email)
	if err != nil {
		return nil, errors.New(utils.UserIsNotExist)
	}

	// Verify password
	if !utils.VerifyPassword(req.Password, user.Password) {
		return nil, errors.New(utils.PasswordInvalid)
	}

	// Generate JTI
	jti := utils.GenerateUUID()

	// Generate tokens
	accessToken, refreshToken, expire, err := middleware.GenerateToken(
		user.ID, user.Username, user.FirstName, user.LastName, user.Email, jti, user.IsAdmin,
	)
	if err != nil {
		return nil, errors.New(utils.LoginError)
	}

	// Save login history
	history := &models.LoginHistory{
		UserID:       user.ID,
		JTI:          jti,
		RefreshToken: refreshToken,
		AccessToken:  accessToken,
	}
	history.CreatedBy = user.ID

	if err := s.userRepo.CreateLoginHistory(history); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		Expire:       expire,
	}, nil
}

func (s *userService) GetMe(userInfo middleware.UserInfo) (*dto.UserResponse, error) {
	user, err := s.userRepo.GetUserByID(userInfo.UserID)
	if err != nil {
		return nil, err
	}

	rolesPerms, _ := s.userRepo.GetRolesPermissionsByUserID(userInfo.UserID)

	// Extract unique roles and permissions
	rolesMap := make(map[string]bool)
	permsMap := make(map[string]bool)

	for _, rp := range rolesPerms {
		if roleName, ok := rp["role_name"].(string); ok {
			rolesMap[roleName] = true
		}
		if permName, ok := rp["permission_name"].(string); ok {
			permsMap[permName] = true
		}
	}

	roles := make([]string, 0, len(rolesMap))
	for role := range rolesMap {
		roles = append(roles, role)
	}

	perms := make([]string, 0, len(permsMap))
	for perm := range permsMap {
		perms = append(perms, perm)
	}

	return &dto.UserResponse{
		ID:          user.ID,
		FirstName:   user.FirstName,
		LastName:    user.LastName,
		Email:       user.Email,
		Phone:       user.Phone,
		Avatar:      user.AvatarURL,
		IsAdmin:     user.IsAdmin,
		Roles:       roles,
		Permissions: perms,
	}, nil
}

func (s *userService) Logout(userInfo middleware.UserInfo) (*dto.MessageResponse, error) {
	// Get login history by JTI
	history, err := s.userRepo.GetLoginHistoryByJTI(userInfo.JTI)
	if err != nil {
		return nil, errors.New(utils.JTINotExist)
	}

	// Mark history as inactive
	history.IsActive = false
	if err := s.userRepo.UpdateLoginHistory(history); err != nil {
		return nil, err
	}

	// Add token to blacklist
	blacklist := &models.TokenBlacklist{
		JTI: userInfo.JTI,
	}
	blacklist.CreatedBy = userInfo.UserID

	if err := s.userRepo.CreateTokenBlacklist(blacklist); err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		Message: "Logout successfully",
	}, nil
}

func (s *userService) GetUsers() ([]dto.UserResponse, error) {
	users, err := s.userRepo.GetUsers()
	if err != nil {
		return nil, err
	}

	var userResponses []dto.UserResponse
	for _, user := range users {
		userResponses = append(userResponses, dto.UserResponse{
			ID:        user.ID,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Email:     user.Email,
			Phone:     user.Phone,
			IsAdmin:   user.IsAdmin,
		})
	}

	return userResponses, nil
}

func (s *userService) ChangePassword(userInfo middleware.UserInfo, req dto.ChangePasswordRequest) (*dto.MessageResponse, error) {
	// Validate passwords match
	if req.NewPassword != req.ReNewPassword {
		return nil, errors.New("re_new_password does not match new_password")
	}

	if req.OldPassword == req.NewPassword {
		return nil, errors.New("new password must be different from old password")
	}

	// Get user
	user, err := s.userRepo.GetUserByID(userInfo.UserID)
	if err != nil {
		return nil, err
	}

	// Verify old password
	if !utils.VerifyPassword(req.OldPassword, user.Password) {
		return nil, errors.New(utils.PasswordInvalid)
	}

	// Hash new password
	hashedPassword, err := utils.HashPassword(req.NewPassword)
	if err != nil {
		return nil, err
	}

	// Update password
	user.Password = hashedPassword
	now := time.Now()
	user.UpdatedAt = &now
	user.UpdatedBy = userInfo.UserID

	if err := s.userRepo.UpdateUser(user); err != nil {
		return nil, err
	}

	return &dto.MessageResponse{
		Message: "Password changed successfully",
	}, nil
}

// Comment methods
func (s *userService) CreateComment(userInfo middleware.UserInfo, petID string, req dto.CommentRequest) (*dto.CommentResponse, error) {
	comment := &models.Comment{
		Content:  req.Content,
		PetID:    petID,
		ParentID: req.ParentID,
	}
	comment.CreatedBy = userInfo.UserID

	if err := s.userRepo.CreateComment(comment); err != nil {
		return nil, err
	}

	return &dto.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		CreatedAt: comment.CreatedAt.Format("2006-01-02 15:04:05"),
		UserID:    userInfo.UserID,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
	}, nil
}

func (s *userService) EditComment(userInfo middleware.UserInfo, petID, commentID string, req dto.CommentRequest) (*dto.CommentResponse, error) {
	comment, err := s.userRepo.GetCommentByID(commentID)
	if err != nil {
		return nil, errors.New("Comment not found")
	}

	// Update comment
	comment.Content = req.Content
	now := time.Now()
	comment.UpdatedAt = &now
	comment.UpdatedBy = userInfo.UserID

	if err := s.userRepo.UpdateComment(comment); err != nil {
		return nil, err
	}

	return &dto.CommentResponse{
		ID:        comment.ID,
		Content:   comment.Content,
		UpdatedAt: comment.UpdatedAt.Format("2006-01-02 15:04:05"),
		UserID:    userInfo.UserID,
		FirstName: userInfo.FirstName,
		LastName:  userInfo.LastName,
	}, nil
}

func (s *userService) GetCommentsByPetID(petID string) ([]dto.CommentResponse, error) {
	results, err := s.userRepo.GetCommentsByPetID(petID)
	if err != nil {
		return nil, err
	}

	var comments []dto.CommentResponse
	for _, r := range results {
		comment := dto.CommentResponse{
			ID:      r["id"].(string),
			Content: r["content"].(string),
			UserID:  r["user_id"].(string),
		}
		if firstName, ok := r["first_name"].(string); ok {
			comment.FirstName = firstName
		}
		if lastName, ok := r["last_name"].(string); ok {
			comment.LastName = lastName
		}
		if avatarURL, ok := r["avatar_url"].(string); ok {
			comment.AvatarURL = avatarURL
		}
		if parentID, ok := r["parent_id"].(string); ok {
			comment.ParentID = parentID
		}
		if createdAt, ok := r["created_at"].(time.Time); ok {
			comment.CreatedAt = createdAt.Format("2006-01-02 15:04:05")
		}
		if updatedAt, ok := r["updated_at"].(time.Time); ok {
			comment.UpdatedAt = updatedAt.Format("2006-01-02 15:04:05")
		}
		comments = append(comments, comment)
	}

	return comments, nil
}
