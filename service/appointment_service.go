package service

import (
	"fmt"
	"log"
	"pet-service/dto"
	"pet-service/middleware"
	"pet-service/scheduler"
	"pet-service/utils"
	"time"

	"gorm.io/gorm"
)

type appointmentService struct {
	db *gorm.DB
}

// NewAppointmentService creates a new appointment service instance
func NewAppointmentService(db *gorm.DB) IAppointmentService {
	return &appointmentService{
		db: db,
	}
}

func (s *appointmentService) RegisterAppointment(userInfo middleware.UserInfo, req dto.AppointmentRequest) (*dto.AppointmentResponse, error) {
	code := utils.GenerateTransactionCode()

	// Schedule email to be sent 10 seconds later
	sch := scheduler.GetScheduler()
	
	// Calculate when to send the email (10 seconds from now)
	sendTime := time.Now().Add(10 * time.Second)
	cronSpec := fmt.Sprintf("%d %d %d %d %d *", 
		sendTime.Second(), sendTime.Minute(), sendTime.Hour(), 
		sendTime.Day(), int(sendTime.Month()))

	_, err := sch.AddJob(cronSpec, func() {
		// This would be the email sending logic
		log.Printf("Sending appointment confirmation email to %s for appointment %s at %s",
			userInfo.Email, code, req.StartTime)
		// TODO: Implement actual email sending using SMTP or email service
	})

	if err != nil {
		log.Printf("Failed to schedule email: %v", err)
	}

	return &dto.AppointmentResponse{
		Code:      code,
		StartTime: req.StartTime,
	}, nil
}
