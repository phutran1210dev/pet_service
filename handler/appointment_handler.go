package handler

import (
	"pet-service/dto"
	"pet-service/middleware"
	"pet-service/service"
	"pet-service/utils"

	"github.com/gin-gonic/gin"
)

type AppointmentHandler struct {
	appointmentService service.IAppointmentService
}

// NewAppointmentHandler creates a new appointment handler instance
func NewAppointmentHandler(appointmentService service.IAppointmentService) *AppointmentHandler {
	return &AppointmentHandler{
		appointmentService: appointmentService,
	}
}

// RegisterAppointment godoc
// @Summary      Register appointment
// @Description  Create a new appointment for a pet
// @Tags         Appointments
// @Accept       json
// @Produce      json
// @Security     Bearer
// @Param        request body dto.AppointmentRequest true "Appointment data"
// @Success      200  {object}  dto.AppointmentResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Failure      400  {object}  dto.ErrorResponse
// @Router       /appointment/register [post]
func (h *AppointmentHandler) RegisterAppointment(c *gin.Context) {
	userInfo, exists := middleware.GetCurrentUser(c)
	if !exists {
		utils.UnauthorizedError(c, utils.ErrCodeUnauthorized, "Unauthorized")
		return
	}

	var req dto.AppointmentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ValidationError(c, err)
		return
	}

	resp, err := h.appointmentService.RegisterAppointment(userInfo, req)
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
