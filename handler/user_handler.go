package handler

import (
	"pet-service/dto"
	"pet-service/middleware"
	"pet-service/service"
	"pet-service/utils"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.IUserService
}

// NewUserHandler creates a new user handler instance
func NewUserHandler(userService service.IUserService) *UserHandler {
	return &UserHandler{
		userService: userService,
	}
}

// Register godoc
// @Summary      Register a new user
// @Description  Create a new user account
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body dto.UserRegisterRequest true "User registration data"
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      409  {object}  dto.ErrorResponse
// @Router       /user [post]
func (h *UserHandler) Register(c *gin.Context) {
	var req dto.UserRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	resp, err := h.userService.Register(req)
	if err != nil {
		if err.Error() == utils.EmailTaken {
			utils.ConflictError(c, utils.ErrCodeEmailTaken, utils.EmailTaken)
		} else {
			utils.InternalServerError(c, utils.ErrCodeInternalError, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, resp)
}

// Login godoc
// @Summary      Login user
// @Description  Authenticate user and return JWT tokens
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Param        request body dto.LoginRequest true "Login credentials"
// @Success      200  {object}  dto.LoginResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /login [post]
func (h *UserHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	resp, err := h.userService.Login(req)
	if err != nil {
		switch err.Error() {
		case utils.UserIsNotExist:
			utils.BadRequestError(c, utils.ErrCodeUserNotFound, utils.UserIsNotExist)
		case utils.PasswordInvalid:
			utils.BadRequestError(c, utils.ErrCodeInvalidPassword, utils.PasswordInvalid)
		default:
			utils.InternalServerError(c, utils.ErrCodeInternalError, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, resp)
}

// GetMe godoc
// @Summary      Get current user
// @Description  Get current authenticated user information
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  dto.UserResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userInfo, exists := middleware.GetCurrentUser(c)
	if !exists {
		utils.UnauthorizedError(c, utils.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	resp, err := h.userService.GetMe(userInfo)
	if err != nil {
		if err.Error() == utils.UserIsNotExist {
			utils.NotFoundError(c, utils.ErrCodeUserNotFound, utils.UserIsNotExist)
		} else {
			utils.InternalServerError(c, utils.ErrCodeInternalError, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, resp)
}

// Logout godoc
// @Summary      Logout user
// @Description  Logout user and blacklist token
// @Tags         Authentication
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /logout [post]
func (h *UserHandler) Logout(c *gin.Context) {
	userInfo, exists := middleware.GetCurrentUser(c)
	if !exists {
		utils.UnauthorizedError(c, utils.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	resp, err := h.userService.Logout(userInfo)
	if err != nil {
		if err.Error() == utils.JTINotExist {
			utils.BadRequestError(c, utils.ErrCodeInvalidToken, utils.JTINotExist)
		} else {
			utils.InternalServerError(c, utils.ErrCodeInternalError, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, resp)
}

// GetUsers godoc
// @Summary      Get all users
// @Description  Get list of all users
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Success      200  {object}  []dto.UserResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /users [get]
func (h *UserHandler) GetUsers(c *gin.Context) {
	resp, err := h.userService.GetUsers()
	if err != nil {
		utils.InternalServerError(c, utils.ErrCodeInternalError, err.Error())
		return
	}

	utils.SuccessResponse(c, resp)
}

// ChangePassword godoc
// @Summary      Change user password
// @Description  Change current user password
// @Tags         Users
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body dto.ChangePasswordRequest true "Password change data"
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /users/change-password [patch]
func (h *UserHandler) ChangePassword(c *gin.Context) {
	userInfo, exists := middleware.GetCurrentUser(c)
	if !exists {
		utils.UnauthorizedError(c, utils.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	var req dto.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	resp, err := h.userService.ChangePassword(userInfo, req)
	if err != nil {
		if err.Error() == utils.PasswordInvalid {
			utils.BadRequestError(c, utils.ErrCodeInvalidPassword, utils.PasswordInvalid)
		} else {
			utils.BadRequestError(c, utils.ErrCodeInvalidInput, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, resp)
}

// CreateComment godoc
// @Summary      Create comment
// @Description  Create a comment on a pet post
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        pet_id path string true "Pet ID"
// @Param        request body dto.CommentRequest true "Comment data"
// @Success      200  {object}  dto.CommentResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /post/{pet_id}/comment [post]
func (h *UserHandler) CreateComment(c *gin.Context) {
	userInfo, exists := middleware.GetCurrentUser(c)
	if !exists {
		utils.UnauthorizedError(c, utils.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	petID := c.Param("pet_id")

	var req dto.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	resp, err := h.userService.CreateComment(userInfo, petID, req)
	if err != nil {
		if err.Error() == utils.PetIDNotExist {
			utils.NotFoundError(c, utils.ErrCodePetNotFound, utils.PetIDNotExist)
		} else {
			utils.InternalServerError(c, utils.ErrCodeInternalError, err.Error())
		}
		return
	}

	utils.CreatedResponse(c, resp)
}

// EditComment godoc
// @Summary      Edit comment
// @Description  Edit an existing comment
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        pet_id path string true "Pet ID"
// @Param        comment_id path string true "Comment ID"
// @Param        request body dto.CommentRequest true "Comment data"
// @Success      200  {object}  dto.CommentResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /post/{pet_id}/comment/{comment_id} [patch]
func (h *UserHandler) EditComment(c *gin.Context) {
	userInfo, exists := middleware.GetCurrentUser(c)
	if !exists {
		utils.UnauthorizedError(c, utils.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	petID := c.Param("pet_id")
	commentID := c.Param("comment_id")

	var req dto.CommentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	resp, err := h.userService.EditComment(userInfo, petID, commentID, req)
	if err != nil {
		if err.Error() == utils.PetIDNotExist {
			utils.NotFoundError(c, utils.ErrCodePetNotFound, utils.PetIDNotExist)
		} else if err.Error() == utils.PermissionDenied {
			utils.ForbiddenError(c, utils.ErrCodePermissionDenied, utils.PermissionDenied)
		} else {
			utils.InternalServerError(c, utils.ErrCodeInternalError, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, resp)
}

// GetComments godoc
// @Summary      Get comments
// @Description  Get all comments for a pet
// @Tags         Comments
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        pet_id path string true "Pet ID"
// @Success      200  {object}  []dto.CommentResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /post/{pet_id}/comments [get]
func (h *UserHandler) GetComments(c *gin.Context) {
	petID := c.Param("pet_id")

	resp, err := h.userService.GetCommentsByPetID(petID)
	if err != nil {
		if err.Error() == utils.PetIDNotExist {
			utils.NotFoundError(c, utils.ErrCodePetNotFound, utils.PetIDNotExist)
		} else {
			utils.InternalServerError(c, utils.ErrCodeInternalError, err.Error())
		}
		return
	}

	utils.SuccessResponse(c, resp)
}
