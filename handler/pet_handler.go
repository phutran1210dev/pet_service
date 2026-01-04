package handler

import (
	"io"
	"pet-service/dto"
	"pet-service/middleware"
	"pet-service/service"
	"pet-service/utils"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PetHandler struct {
	petService service.IPetService
	db         *gorm.DB
}

// NewPetHandler creates a new pet handler instance
func NewPetHandler(petService service.IPetService, db *gorm.DB) *PetHandler {
	return &PetHandler{
		petService: petService,
		db:         db,
	}
}

// CreatePet godoc
// @Summary      Create a new pet
// @Description  Create a new pet for the current user
// @Tags         Pets
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body dto.PetCreateRequest true "Pet data"
// @Success      200  {object}  dto.PetResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /pet [post]
func (h *PetHandler) CreatePet(c *gin.Context) {
	userInfo, exists := middleware.GetCurrentUser(c)
	if !exists {
		utils.UnauthorizedError(c, utils.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	var req dto.PetCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	resp, err := h.petService.CreatePet(userInfo, req)
	if err != nil {
		utils.InternalServerError(c, utils.ErrCodeInternalError, err.Error())
		return
	}

	utils.CreatedResponse(c, resp)
}

// GetPets godoc
// @Summary      Get all pets
// @Description  Get list of all pets with pagination and filters
// @Tags         Pets
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        page query int false "Page number" default(1)
// @Param        page_size query int false "Page size" default(10)
// @Param        search query string false "Search query"
// @Param        name query string false "Pet name filter"
// @Success      200  {object}  dto.PaginationResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /pets [get]
func (h *PetHandler) GetPets(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	search := c.Query("search")
	name := c.Query("name")

	resp, err := h.petService.GetPets(h.db, page, pageSize, search, name)
	if err != nil {
		utils.InternalServerError(c, utils.ErrCodeInternalError, err.Error())
		return
	}

	utils.SuccessResponse(c, resp)
}

// GetPetDetail godoc
// @Summary      Get pet detail
// @Description  Get detailed information about a specific pet
// @Tags         Pets
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        id path string true "Pet ID"
// @Success      200  {object}  dto.PetDetailResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /pet/{id} [get]
func (h *PetHandler) GetPetDetail(c *gin.Context) {
	petID := c.Param("id")

	resp, err := h.petService.GetPetDetail(petID)
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

// CreatePetLifeEvent godoc
// @Summary      Create pet life event
// @Description  Create a life event for a pet
// @Tags         Pets
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body dto.PetLifeEventRequest true "Life event data"
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /pet/life-event [post]
func (h *PetHandler) CreatePetLifeEvent(c *gin.Context) {
	userInfo, exists := middleware.GetCurrentUser(c)
	if !exists {
		utils.UnauthorizedError(c, utils.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	var req dto.PetLifeEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	resp, err := h.petService.CreatePetLifeEvent(userInfo, req)
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

// UploadAvatar godoc
// @Summary      Upload pet avatar
// @Description  Upload an avatar image for a pet
// @Tags         Pets
// @Accept       multipart/form-data
// @Produce      json
// @Security     Bearer
// @Param        pet_id path string true "Pet ID"
// @Param        file formData file true "Avatar image file"
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /pet/{pet_id}/images [post]
func (h *PetHandler) UploadAvatar(c *gin.Context) {
	petID := c.Param("pet_id")

	file, err := c.FormFile("file")
	if err != nil {
		utils.BadRequestError(c, utils.ErrCodeInvalidInput, "File is required")
		return
	}

	openedFile, err := file.Open()
	if err != nil {
		utils.BadRequestError(c, utils.ErrCodeInvalidInput, "Cannot open file")
		return
	}
	defer openedFile.Close()

	fileData, err := io.ReadAll(openedFile)
	if err != nil {
		utils.BadRequestError(c, utils.ErrCodeInvalidInput, "Cannot read file")
		return
	}

	resp, err := h.petService.UploadAvatar(petID, fileData, file.Header.Get("Content-Type"))
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

// UploadGallery godoc
// @Summary      Upload pet gallery images
// @Description  Upload multiple images to pet gallery
// @Tags         Pets
// @Accept       multipart/form-data
// @Produce      json
// @Security     Bearer
// @Param        pet_id path string true "Pet ID"
// @Param        files formData file true "Gallery image files" 
// @Success      200  {object}  dto.MessageResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /pet/{pet_id}/gallery [post]
func (h *PetHandler) UploadGallery(c *gin.Context) {
	petID := c.Param("pet_id")

	form, err := c.MultipartForm()
	if err != nil {
		utils.BadRequestError(c, utils.ErrCodeInvalidInput, "Files are required")
		return
	}

	files := form.File["files"]
	if len(files) == 0 {
		utils.BadRequestError(c, utils.ErrCodeInvalidInput, "At least one file is required")
		return
	}

	var fileReaders []io.Reader
	var fileNames []string
	var contentTypes []string

	for _, file := range files {
		openedFile, err := file.Open()
		if err != nil {
			continue
		}
		defer openedFile.Close()

		fileReaders = append(fileReaders, openedFile)
		fileNames = append(fileNames, file.Filename)
		contentTypes = append(contentTypes, file.Header.Get("Content-Type"))
	}

	resp, err := h.petService.UploadGallery(petID, fileReaders, fileNames, contentTypes)
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
