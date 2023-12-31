package handlers

import (
	"errors"
	"log"
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

	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"message": err})
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *userHandler) CreateNewUser(c *gin.Context) {
	var newUser models.User

	if err := c.BindJSON(&newUser); err != nil {
		c.AbortWithStatusJSON(error_utils.InvalidRequestBody{}.GetAPIError())
		return
	}

	err := h.s.CreateNewUser(newUser)

	if err != nil {
		apiErr := checkServiceErrors(newUser, err)
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	c.Status(http.StatusCreated)
}

func checkServiceErrors(user models.User, err error) error_utils.APIError {
	log.Print(err)

	alreadyExists := &error_utils.FirestoreAlreadyExistsError{}
	unauthorized := &error_utils.FirebaseUnauthorizedError{}
	notFound := &error_utils.FirestoreNotFoundError{}

	if errors.As(err, unauthorized) {
		return error_utils.APIUnauthorized{}
	}
	if errors.As(err, alreadyExists) {
		return error_utils.UserDuplicated{UserID: user.UID}
	}
	if errors.As(err, notFound) {
		return error_utils.APIUnauthorized{}
	}
	return error_utils.APIUnknown{}
}
