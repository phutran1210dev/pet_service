package service

import (
	"errors"
	"io"
	"math"
	"pet-service/dto"
	"pet-service/middleware"
	"pet-service/models"
	"pet-service/repository"
	"pet-service/storage"
	"pet-service/utils"
	"time"

	"gorm.io/gorm"
)

type petService struct {
	petRepo repository.IPetRepository
}

// NewPetService creates a new pet service instance
func NewPetService(petRepo repository.IPetRepository) IPetService {
	return &petService{
		petRepo: petRepo,
	}
}

func (s *petService) CreatePet(userInfo middleware.UserInfo, req dto.PetCreateRequest) (*dto.PetResponse, error) {
	dateOfBirth, _ := utils.ParseDateTime(req.DateOfBirth)
	var dateOfDeath *time.Time
	if req.DateOfDeath != "" {
		dateOfDeath, _ = utils.ParseDateTime(req.DateOfDeath)
	}

	pet := &models.Pet{
		Name:        req.Name,
		Gender:      req.Gender,
		DateOfBirth: dateOfBirth,
		DateOfDeath: dateOfDeath,
		Breed:       req.Breed,
		Description: req.Description,
		Type:        req.Type,
		UserID:      userInfo.UserID,
	}
	pet.CreatedBy = userInfo.UserID

	if err := s.petRepo.CreatePet(pet); err != nil {
		return nil, err
	}

	return &dto.PetResponse{
		ID:   pet.ID,
		Name: pet.Name,
		Gender: pet.Gender,
		DateOfBirth: pet.DateOfBirth.Format("2006-01-02"),
		Breed: pet.Breed,
		Description: pet.Description,
		Type: pet.Type,
	}, nil
}

func (s *petService) GetPets(db *gorm.DB, page, pageSize int, search, name string) (*dto.PaginationResponse, error) {
	query := db.Model(&models.Pet{}).Where("is_active = ?", true)

	// Apply search filters
	if search != "" {
		query = query.Where("name ILIKE ? OR breed ILIKE ? OR type ILIKE ?",
			"%"+search+"%", "%"+search+"%", "%"+search+"%")
	}

	if name != "" {
		query = query.Where("name ILIKE ?", "%"+name+"%")
	}

	// Count total
	var totalItem int64
	query.Count(&totalItem)

	// Calculate pagination
	totalPage := int64(math.Ceil(float64(totalItem) / float64(pageSize)))
	offset := (page - 1) * pageSize

	// Get data with user preload
	var pets []models.Pet
	query.Preload("User").Order("date_of_birth DESC, name ASC").Limit(pageSize).Offset(offset).Find(&pets)

	return &dto.PaginationResponse{
		Data: pets,
		Meta: dto.PaginationMeta{
			TotalItems: totalItem,
			TotalPages: totalPage,
			Page:       page,
			PageSize:   pageSize,
		},
	}, nil
}

func (s *petService) GetPetDetail(petID string) (*dto.PetDetailResponse, error) {
	results, err := s.petRepo.GetPetDetail(petID)
	if err != nil || len(results) == 0 {
		return nil, errors.New(utils.PetIDNotExist)
	}

	// Build response
	response := dto.PetDetailResponse{
		Events: []dto.PetLifeEventItem{},
		Medias: []dto.MediaItem{},
	}

	eventIDs := make(map[string]bool)
	mediaIDs := make(map[string]bool)

	for _, r := range results {
		// Set pet info (only once)
		if response.ID == "" {
			response.ID = r["pet_id"].(string)
			response.Name = r["pet_name"].(string)
			response.Gender = r["pet_gender"].(bool)
			if breed, ok := r["pet_breed"].(string); ok {
				response.Breed = breed
			}
			if desc, ok := r["pet_description"].(string); ok {
				response.Description = desc
			}
			if dob, ok := r["pet_date_of_birth"].(time.Time); ok {
				response.DateOfBirth = dob.Format("2006-01-02")
			}
			if dod, ok := r["pet_date_of_death"].(time.Time); ok {
				response.DateOfDeath = dod.Format("2006-01-02")
			}
		}

		// Add events
		if eventID, ok := r["event_id"].(string); ok && eventID != "" && !eventIDs[eventID] {
			event := dto.PetLifeEventItem{
				ID:    eventID,
				Title: r["event_title"].(string),
			}
			if date, ok := r["event_date"].(time.Time); ok {
				event.Date = date.Format("2006-01-02")
			}
			if loc, ok := r["event_location"].(string); ok {
				event.Location = loc
			}
			if story, ok := r["event_story"].(string); ok {
				event.Story = story
			}
			response.Events = append(response.Events, event)
			eventIDs[eventID] = true
		}

		// Add medias
		if mediaID, ok := r["media_id"].(string); ok && mediaID != "" && !mediaIDs[mediaID] {
			media := dto.MediaItem{
				ID:  mediaID,
				URL: r["media_url"].(string),
			}
			response.Medias = append(response.Medias, media)
			mediaIDs[mediaID] = true
		}
	}

	return &response, nil
}

func (s *petService) CreatePetLifeEvent(userInfo middleware.UserInfo, req dto.PetLifeEventRequest) (*dto.PetLifeEventResponse, error) {
	date, _ := utils.ParseDateTime(req.Date)

	event := &models.PetLifeEvent{
		PetID:    req.PetID,
		Title:    req.Title,
		Date:     *date,
		Location: req.Location,
		Story:    req.Story,
	}
	event.CreatedBy = userInfo.UserID

	if err := s.petRepo.CreateLifeEvent(event); err != nil {
		return nil, err
	}

	return &dto.PetLifeEventResponse{
		ID:     event.ID,
		PetID:  event.PetID,
		Title:  event.Title,
		Date:   event.Date.Format("2006-01-02"),
	}, nil
}

func (s *petService) UploadAvatar(petID string, fileData []byte, contentType string) (*dto.MediaResponse, error) {
	pet, err := s.petRepo.GetPetByID(petID)
	if err != nil {
		return nil, errors.New(utils.PetIDNotExist)
	}

	imageID := utils.GenerateUUID()
	objectName := "pet/" + petID + "/" + imageID

	minioClient := storage.GetMinioClient()
	url, err := minioClient.UploadFile(objectName, fileData, contentType)
	if err != nil {
		return nil, err
	}

	pet.AvtURL = url
	if err := s.petRepo.UpdatePet(pet); err != nil {
		return nil, err
	}

	return &dto.MediaResponse{
		ID:  imageID,
		URL: url,
	}, nil
}

func (s *petService) UploadGallery(petID string, files []io.Reader, fileNames []string, contentTypes []string) ([]dto.MediaResponse, error) {
	_, err := s.petRepo.GetPetByID(petID)
	if err != nil {
		return nil, errors.New(utils.PetIDNotExist)
	}

	minioClient := storage.GetMinioClient()
	var medias []models.Media

	for i, file := range files {
		mediaID := utils.GenerateUUID()
		objectName := "pet/" + petID + "/" + mediaID

		// Read file data
		fileData, err := io.ReadAll(file)
		if err != nil {
			continue
		}

		url, err := minioClient.UploadFile(objectName, fileData, contentTypes[i])
		if err != nil {
			continue
		}

		media := models.Media{
			Name:  fileNames[i],
			URL:   url,
			PetID: petID,
		}
		media.ID = mediaID
		medias = append(medias, media)
	}

	if err := s.petRepo.CreateMediaBatch(medias); err != nil {
		return nil, err
	}

	var response []dto.MediaResponse
	for _, media := range medias {
		response = append(response, dto.MediaResponse{
			ID:  media.ID,
			URL: media.URL,
		})
	}

	return response, nil
}
