package handlers

import (
	"net/http"
	"strconv"
	"time"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/services"
	util_models "github.com/AndresCRamos/midas-app-api/utils/api/models"
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

func (h *movementHandler) GetMovementsByUserAndDate(c *gin.Context) {
	userID, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(error_utils.CantGetUser{}.GetAPIError())
		return
	}

	var page int
	var err error

	pageStr, exists := c.GetQuery("page")
	if !exists {
		page = 1
	} else if page, err = strconv.Atoi(pageStr); err != nil {
		c.AbortWithStatusJSON(util_models.PaginatedTypeError{}.GetAPIError())
		return
	}

	dateFromStr, exists := c.GetQuery("date_from")
	var dateFrom time.Time
	if !exists {
		dateFrom = time.Now().UTC().Add(-30 * 24 * time.Hour)
	} else if dateFrom, err = time.Parse(time.DateOnly, dateFromStr); err != nil {
		c.AbortWithStatusJSON(error_utils.APIBadDateFormat{DateString: dateFromStr, DateField: "date_from", Format: time.DateOnly}.GetAPIError())
		return
	}

	dateToStr, exists := c.GetQuery("date_to")
	var dateTo time.Time
	if !exists {
		dateTo = time.Now().UTC()
	} else if dateTo, err = time.Parse(time.DateOnly, dateToStr); err != nil {
		c.AbortWithStatusJSON(error_utils.APIBadDateFormat{DateString: dateToStr, DateField: "date_to", Format: time.DateOnly}.GetAPIError())
		return
	}

	movementSearchResult, err := h.s.GetMovementsByUserAndDate(userID.(string), page, dateFrom, dateTo)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(userID.(string), err, "movement")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	movementRetrievedData := []models.MovementRetrieve{}

	for _, movementData := range movementSearchResult.Data {
		retrievedData := models.MovementRetrieve{}
		retrievedData.ParseMovement(movementData)
		movementRetrievedData = append(movementRetrievedData, retrievedData)
	}

	movementSearchRetrieve := util_models.PaginatedSearch[models.MovementRetrieve]{
		CurrentPage: movementSearchResult.CurrentPage,
		TotalData:   movementSearchResult.TotalData,
		PageSize:    movementSearchResult.TotalData,
		Data:        movementRetrievedData,
	}

	c.JSON(http.StatusOK, movementSearchRetrieve)
}

func (h *movementHandler) CreateNewMovement(c *gin.Context) {
	var newMovement models.MovementCreate

	if err := c.ShouldBindJSON(&newMovement); err != nil {
		c.AbortWithStatusJSON(error_utils.APIInvalidRequestBody{DetailErr: err}.GetAPIError())
		return
	}

	movement := newMovement.ParseMovement()
	userID, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(error_utils.CantGetUser{}.GetAPIError())
		return
	}
	movement.OwnerId = userID.(string)

	movement, err := h.s.CreateNewMovement(movement)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(movement.UID, err, "movement")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	movementData := models.MovementRetrieve{}
	movementData.ParseMovement(movement)

	c.JSON(http.StatusCreated, movementData)
}

func (h *movementHandler) UpdateMovement(c *gin.Context) {
	id := c.Param("id")
	var updatedMovement models.MovementUpdate

	if err := c.ShouldBindJSON(&updatedMovement); err != nil {
		c.AbortWithStatusJSON(error_utils.APIInvalidRequestBody{DetailErr: err}.GetAPIError())
		return
	}

	movement := updatedMovement.ParseMovement()

	movement.UID = id
	userID, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(error_utils.CantGetUser{}.GetAPIError())
		return
	}
	movement.OwnerId = userID.(string)
	movement, err := h.s.UpdateMovement(movement)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(id, err, "movement")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	var movementData models.MovementRetrieve

	movementData.ParseMovement(movement)

	c.JSON(http.StatusCreated, movementData)
}
