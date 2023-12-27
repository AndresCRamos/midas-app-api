package handlers

import (
	"net/http"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/services"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	s services.UserService
}

func NewUserHandler(s services.UserService) *userHandler {
	return &userHandler{
		s: s,
	}
}

func (h *userHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")

	user, err := h.s.GetUserByID(id)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
	}

	c.JSON(200, user)
}

func (h *userHandler) CreateNewUser(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	err := h.s.CreateNewUser(newUser)

	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": err})
	}

	c.Status(http.StatusCreated)
}
