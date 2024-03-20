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

type sourceHandler struct {
	s services.SourceService
}

func NewSourceHandler(s services.SourceService) *sourceHandler {
	return &sourceHandler{
		s: s,
	}
}

func (h *sourceHandler) GetSourceByID(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(error_utils.CantGetUser{}.GetAPIError())
		return
	}

	source, err := h.s.GetSourceByID(id, userID.(string))
	sourceData := models.SourceRetrieve{}

	sourceData.ParseSource(source)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(id, err, "source")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	c.JSON(http.StatusOK, sourceData)
}

func (h *sourceHandler) GetSourcesByUser(c *gin.Context) {
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

	sourceSearchResult, err := h.s.GetSourcesByUser(userID.(string), page)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(userID.(string), err, "source")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	sourceRetrievedData := []models.SourceRetrieve{}

	for _, sourceData := range sourceSearchResult.Data {
		retrievedData := models.SourceRetrieve{}
		retrievedData.ParseSource(sourceData)
		sourceRetrievedData = append(sourceRetrievedData, retrievedData)
	}

	sourceSearchRetrieve := util_models.PaginatedSearch[models.SourceRetrieve]{
		CurrentPage: sourceSearchResult.CurrentPage,
		TotalData:   sourceSearchResult.TotalData,
		PageSize:    sourceSearchResult.TotalData,
		Data:        sourceRetrievedData,
	}

	c.JSON(http.StatusOK, sourceSearchRetrieve)
}

func (h *sourceHandler) GetMovementsBySourceAndDate(c *gin.Context) {
	id := c.Param("id")

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

	data, err := h.s.GetMovementsBySourceAndDate(id, userID.(string), page, dateFrom, dateTo)

	if err != nil {
		c.AbortWithStatusJSON(error_utils.CheckServiceErrors(id, err, "source").GetAPIError())
		return
	}

	c.JSON(http.StatusOK, data)
}

func (h *sourceHandler) CreateNewSource(c *gin.Context) {
	var newSource models.SourceCreate

	if err := c.ShouldBindJSON(&newSource); err != nil {
		c.AbortWithStatusJSON(error_utils.APIInvalidRequestBody{DetailErr: err}.GetAPIError())
		return
	}

	source := newSource.ParseSource()
	userID, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(error_utils.CantGetUser{}.GetAPIError())
		return
	}
	source.OwnerId = userID.(string)

	source, err := h.s.CreateNewSource(source)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(source.UID, err, "source")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	sourceData := models.SourceRetrieve{}
	sourceData.ParseSource(source)

	c.JSON(http.StatusCreated, sourceData)
}

func (h *sourceHandler) UpdateSource(c *gin.Context) {
	id := c.Param("id")
	var updatedSource models.SourceUpdate

	if err := c.ShouldBindJSON(&updatedSource); err != nil {
		c.AbortWithStatusJSON(error_utils.APIInvalidRequestBody{DetailErr: err}.GetAPIError())
		return
	}

	source := updatedSource.ParseSource()

	source.UID = id
	userID, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(error_utils.CantGetUser{}.GetAPIError())
		return
	}
	source.OwnerId = userID.(string)
	source, err := h.s.UpdateSource(source)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(id, err, "source")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	var sourceData models.SourceRetrieve

	sourceData.ParseSource(source)

	c.JSON(http.StatusCreated, sourceData)
}

func (h *sourceHandler) DeleteSource(c *gin.Context) {
	id := c.Param("id")
	userID, exists := c.Get("user")
	if !exists {
		c.AbortWithStatusJSON(error_utils.CantGetUser{}.GetAPIError())
		return
	}

	err := h.s.DeleteSource(id, userID.(string))

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(id, err, "source")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	c.Status(http.StatusOK)
}
