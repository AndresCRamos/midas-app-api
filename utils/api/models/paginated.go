package models

type PaginatedSearch[T any] struct {
	CurrentPage int
	TotalData   int
	PageSize    int
	Data        []T
}
