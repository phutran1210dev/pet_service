package repository

import (
	"pet-service/models"

	"gorm.io/gorm"
)

type PetRepository struct {
	DB *gorm.DB
}

func NewPetRepository(db *gorm.DB) *PetRepository {
	return &PetRepository{DB: db}
}

func (r *PetRepository) CreatePet(pet *models.Pet) error {
	return r.DB.Create(pet).Error
}

func (r *PetRepository) GetPetByID(id string) (*models.Pet, error) {
	var pet models.Pet
	err := r.DB.Where("id = ? AND is_active = ?", id, true).First(&pet).Error
	if err != nil {
		return nil, err
	}
	return &pet, nil
}

func (r *PetRepository) GetPets(query *gorm.DB) *gorm.DB {
	return query.Where("is_active = ?", true)
}

func (r *PetRepository) UpdatePet(pet *models.Pet) error {
	return r.DB.Save(pet).Error
}

// Pet Life Events
func (r *PetRepository) CreateLifeEvent(event *models.PetLifeEvent) error {
	return r.DB.Create(event).Error
}

func (r *PetRepository) GetPetDetail(petID string) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	err := r.DB.Table("pets").
		Select(`pets.id as pet_id, pets.name as pet_name, pets.gender as pet_gender, 
				pets.breed as pet_breed, pets.description as pet_description, 
				pets.date_of_birth as pet_date_of_birth, pets.date_of_death as pet_date_of_death,
				pet_life_events.id as event_id, pet_life_events.title as event_title, 
				pet_life_events.date as event_date, pet_life_events.location as event_location, 
				pet_life_events.story as event_story,
				medias.id as media_id, medias.url as media_url`).
		Joins("LEFT JOIN pet_life_events ON pet_life_events.pet_id = pets.id AND pet_life_events.is_active = true").
		Joins("LEFT JOIN medias ON medias.pet_id = pets.id AND medias.is_active = true").
		Where("pets.id = ? AND pets.is_active = ?", petID, true).
		Scan(&results).Error

	return results, err
}

// Media
func (r *PetRepository) CreateMediaBatch(medias []models.Media) error {
	return r.DB.Create(&medias).Error
}
