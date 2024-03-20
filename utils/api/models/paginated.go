package models

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	page_type_error = "page must be a number"
)

type PaginatedSearch[T any] struct {
	CurrentPage int `json:"page"`
	TotalData   int `json:"total"`
	PageSize    int `json:"size"`
	Data        []T `json:"data"`
}

type PaginatedTypeError struct{}

func (pte PaginatedTypeError) Error() string {
	return page_type_error
}

func (pte PaginatedTypeError) GetAPIError() (int, gin.H) {
	return http.StatusBadRequest, gin.H{
		"error": page_type_error,
	}
}
