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

	source, err := h.s.GetSourceByID(id)
	source.UID = id

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(id, err, "source")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	c.JSON(http.StatusOK, source)
}

func (h *sourceHandler) CreateNewSource(c *gin.Context) {
	var newSource models.Source

	if err := c.ShouldBindJSON(&newSource); err != nil {
		c.AbortWithStatusJSON(error_utils.APIInvalidRequestBody{DetailErr: err}.GetAPIError())
		return
	}

	err := h.s.CreateNewSource(newSource)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(newSource.UID, err, "source")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	c.Status(http.StatusCreated)
}

func (h *sourceHandler) UpdateSource(c *gin.Context) {
	id := c.Param("id")
	var newSource models.Source

	if err := c.ShouldBindJSON(&newSource); err != nil {
		c.AbortWithStatusJSON(error_utils.APIInvalidRequestBody{DetailErr: err}.GetAPIError())
		return
	}

	newSource.UID = id
	err := h.s.CreateNewSource(newSource)

	if err != nil {
		apiErr := error_utils.CheckServiceErrors(newSource.UID, err, "source")
		c.AbortWithStatusJSON(apiErr.GetAPIError())
		return
	}

	c.Status(http.StatusCreated)
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
