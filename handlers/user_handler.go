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
	userID, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(error_utils.CantGetUser{}.GetAPIError())
		return
	}

	userIDStr := userID.(string)

	user, err := h.s.GetUserByID(userIDStr)
	user.UID = userIDStr

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(userIDStr, err, "user")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) CreateNewUser(c *gin.Context) {
	userID, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(error_utils.CantGetUser{}.GetAPIError())
		return
	}
	var newUser models.UserCreate

	if err := c.ShouldBindJSON(&newUser); err != nil {
		c.AbortWithStatusJSON(error_utils.APIInvalidRequestBody{DetailErr: err}.GetAPIError())
		return
	}

	user := newUser.ParseUser()

	user.UID = userID.(string)

	err := h.s.CreateNewUser(user)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(userID.(string), err, "user")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	c.Status(http.StatusCreated)
}
