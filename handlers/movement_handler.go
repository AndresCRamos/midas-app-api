package handlers

import (
	"net/http"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/services"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
	"github.com/gin-gonic/gin"
)

type movementHandler struct {
	s services.MovementService
}

func NewMovementHandler(s services.MovementService) *movementHandler {
	return &movementHandler{
		s: s,
	}
}

func (h *movementHandler) GetMovementByID(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(error_utils.CantGetUser{}.GetAPIError())
		return
	}

	movement, err := h.s.GetMovementByID(id, userID.(string))
	movementData := models.MovementRetrieve{}

	movementData.ParseMovement(movement)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(id, err, "movement")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	c.JSON(http.StatusOK, movementData)
}
