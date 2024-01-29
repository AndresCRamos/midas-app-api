package handlers

import (
	"net/http"

	"github.com/AndresCRamos/midas-app-api/models"
	"github.com/AndresCRamos/midas-app-api/services"
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

	err := h.s.DeleteSource(id)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(id, err, "source")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	c.Status(http.StatusOK)
}
