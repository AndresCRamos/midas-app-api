package handlers

import (
	"net/http"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/services"
	error_utils "github.com/AndresCRamos/midas-app-api/utils/errors"
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
	user.UID = id

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(id, err, "user")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) CreateNewUser(c *gin.Context) {
	var newUser models.User

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.AbortWithStatusJSON(error_utils.APIInvalidRequestBody{DetailErr: err}.GetAPIError())
		return
	}

	err := h.s.CreateNewUser(newUser)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(newUser.UID, err, "user")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	c.Status(http.StatusCreated)
}
