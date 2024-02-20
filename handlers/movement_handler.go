package handlers

import (
	"net/http"

	"github.com/AndresCRamos/midas-app-api/services"
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

	movement, err := h.s.GetMovementByID(id, "")

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": "error"})
	}

	c.JSON(200, movement)
}
