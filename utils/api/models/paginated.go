package models

type PaginatedSearch[T any] struct {
	CurrentPage int `json:"page"`
	TotalData   int `json:"total"`
	PageSize    int `json:"size"`
	Data        []T `json:"data"`
}
