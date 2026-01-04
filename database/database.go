package database

import (
	"fmt"
	"log"
	"pet-service/config"
	"pet-service/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func ConnectDatabase() {
	cfg := config.AppConfig
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable TimeZone=%s",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.TimeZone)

	var err error
	var dbLogger logger.Interface
	if cfg.Debug {
		dbLogger = logger.Default.LogMode(logger.Info)
	} else {
		dbLogger = logger.Default.LogMode(logger.Silent)
	}

	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: dbLogger,
	})

	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connection established")

	// Auto migrate tables
	if err := DB.AutoMigrate(
		&models.User{},
		&models.Role{},
		&models.Permission{},
		&models.RolePermission{},
		&models.UserRole{},
		&models.Pet{},
		&models.Media{},
		&models.PetLifeEvent{},
		&models.Comment{},
		&models.Appointment{},
		&models.LoginHistory{},
		&models.TokenBlacklist{},
	); err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	log.Println("Database migration completed")

	// Seed initial data
	seedData()
}

func GetDB() *gorm.DB {
	return DB
}

func seedData() {
	// Check if roles already exist
	var count int64
	DB.Model(&models.Role{}).Count(&count)
	if count > 0 {
		log.Println("Data already seeded, skipping")
		return
	}

	log.Println("Seeding initial data...")

	// Seed roles
	roles := []models.Role{
		{BaseModel: models.BaseModel{ID: "DAA6B933DAF84FBB99477508A4DAC571"}, Name: "Admin"},
		{BaseModel: models.BaseModel{ID: "DAA6B933DAF84FBB99477508A4DAC572"}, Name: "Editor"},
		{BaseModel: models.BaseModel{ID: "DAA6B933DAF84FBB99477508A4DAC573"}, Name: "User"},
	}
	for i := range roles {
		roles[i].IsActive = true
		DB.Create(&roles[i])
	}

	// Seed permissions
	permissions := []models.Permission{
		{BaseModel: models.BaseModel{ID: "A3CBE2F6B7CD9A34D0FA23593D0E42FA"}, Name: "view_pet"},
		{BaseModel: models.BaseModel{ID: "B02AC61D6F5B8D83111F42A7C75C1D2A"}, Name: "add_pet"},
		{BaseModel: models.BaseModel{ID: "E0B57C9A19E0412399B72391FC5D50CC"}, Name: "edit_pet"},
		{BaseModel: models.BaseModel{ID: "F53A2B89FFB1B3C42B9E6814E5869338"}, Name: "delete_pet"},
		{BaseModel: models.BaseModel{ID: "6D02A9C9378D1D86F39B3BDA48DA9F8E"}, Name: "view_user"},
		{BaseModel: models.BaseModel{ID: "6D02A9C9378D1D86F39B3BDA48DA9F8F"}, Name: "add_user"},
		{BaseModel: models.BaseModel{ID: "6D02A9C9378D1D86F39B3BDA48DA9F8D"}, Name: "edit_user"},
		{BaseModel: models.BaseModel{ID: "6D02A9C9378D1D86F39B3BDA48DA9F8G"}, Name: "delete_user"},
	}
	for i := range permissions {
		permissions[i].IsActive = true
		DB.Create(&permissions[i])
	}

	// Seed role permissions
	rolePerms := []models.RolePermission{
		// Editor permissions
		{BaseModel: models.BaseModel{ID: "1D02A9C9378D1D86F39B3BDA48DA9F8G"}, RoleID: "DAA6B933DAF84FBB99477508A4DAC572", PermissionID: "A3CBE2F6B7CD9A34D0FA23593D0E42FA"},
		{BaseModel: models.BaseModel{ID: "2D02A9C9378D1D86F39B3BDA48DA9F8G"}, RoleID: "DAA6B933DAF84FBB99477508A4DAC572", PermissionID: "B02AC61D6F5B8D83111F42A7C75C1D2A"},
		{BaseModel: models.BaseModel{ID: "3D02A9C9378D1D86F39B3BDA48DA9F8G"}, RoleID: "DAA6B933DAF84FBB99477508A4DAC572", PermissionID: "E0B57C9A19E0412399B72391FC5D50CC"},
		{BaseModel: models.BaseModel{ID: "4D02A9C9378D1D86F39B3BDA48DA9F8G"}, RoleID: "DAA6B933DAF84FBB99477508A4DAC572", PermissionID: "F53A2B89FFB1B3C42B9E6814E5869338"},
		{BaseModel: models.BaseModel{ID: "5D02A9C9378D1D86F39B3BDA48DA9F8G"}, RoleID: "DAA6B933DAF84FBB99477508A4DAC572", PermissionID: "6D02A9C9378D1D86F39B3BDA48DA9F8E"},
		// User permissions
		{BaseModel: models.BaseModel{ID: "6D02A9C9378D1D86F39B3BDA48DA9F8G"}, RoleID: "DAA6B933DAF84FBB99477508A4DAC573", PermissionID: "6D02A9C9378D1D86F39B3BDA48DA9F8E"},
		{BaseModel: models.BaseModel{ID: "7D02A9C9378D1D86F39B3BDA48DA9F8G"}, RoleID: "DAA6B933DAF84FBB99477508A4DAC573", PermissionID: "6D02A9C9378D1D86F39B3BDA48DA9F8F"},
		{BaseModel: models.BaseModel{ID: "8D02A9C9378D1D86F39B3BDA48DA9F8G"}, RoleID: "DAA6B933DAF84FBB99477508A4DAC573", PermissionID: "6D02A9C9378D1D86F39B3BDA48DA9F8D"},
	}
	for i := range rolePerms {
		rolePerms[i].IsActive = true
		DB.Create(&rolePerms[i])
	}

	// Seed test users (password: 123456)
	users := []models.User{
		{
			BaseModel: models.BaseModel{ID: "9D02A9C9378D1D86F39B3BDA48DA9F8G"},
			FirstName: "Admin",
			LastName:  "User",
			Email:     "admin@example.com",
			Phone:     "0123456789",
			Username:  "admin",
			Gender:    true,
			Password:  "$2b$12$1hw8y1rphV0JGChNge5/Lu.YMLCdOl/CX/q9YwLfgXi4Z1.u3qJ.i",
			IsAdmin:   true,
		},
		{
			BaseModel: models.BaseModel{ID: "1102A9C9378D1D86F39B3BDA48DA9F8G"},
			FirstName: "Editor",
			LastName:  "User",
			Email:     "editor@example.com",
			Phone:     "0123456789",
			Username:  "editor",
			Gender:    true,
			Password:  "$2b$12$1hw8y1rphV0JGChNge5/Lu.YMLCdOl/CX/q9YwLfgXi4Z1.u3qJ.i",
			IsAdmin:   false,
		},
		{
			BaseModel: models.BaseModel{ID: "1002A9C9378D1D86F39B3BDA48DA9F8G"},
			FirstName: "User",
			LastName:  "Test",
			Email:     "user@example.com",
			Phone:     "0123456789",
			Username:  "user",
			Gender:    true,
			Password:  "$2b$12$1hw8y1rphV0JGChNge5/Lu.YMLCdOl/CX/q9YwLfgXi4Z1.u3qJ.i",
			IsAdmin:   false,
		},
	}
	for i := range users {
		users[i].IsActive = true
		DB.Create(&users[i])
	}

	// Seed user roles
	userRoles := []models.UserRole{
		{BaseModel: models.BaseModel{ID: "AD02A9C9378D1D86F39B3BDA48DA9F8G"}, UserID: "9D02A9C9378D1D86F39B3BDA48DA9F8G", RoleID: "DAA6B933DAF84FBB99477508A4DAC571"},
		{BaseModel: models.BaseModel{ID: "BD02A9C9378D1D86F39B3BDA48DA9F8G"}, UserID: "1102A9C9378D1D86F39B3BDA48DA9F8G", RoleID: "DAA6B933DAF84FBB99477508A4DAC572"},
		{BaseModel: models.BaseModel{ID: "CD02A9C9378D1D86F39B3BDA48DA9F8G"}, UserID: "1002A9C9378D1D86F39B3BDA48DA9F8G", RoleID: "DAA6B933DAF84FBB99477508A4DAC573"},
	}
	for i := range userRoles {
		userRoles[i].IsActive = true
		DB.Create(&userRoles[i])
	}

	log.Println("Initial data seeded successfully")
}
