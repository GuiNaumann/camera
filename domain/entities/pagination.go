package entities

type PaginatedList struct {
	Items []interface{} `json:"items"`

	Page int64 `json:"page"`

	TotalCount int64 `json:"totalCount"`

	RequestedCount int64 `json:"requestedCount"`
}

type PaginatedListUpdated[T any] struct {
	Items []T `json:"items"`

	Page int64 `json:"page"`

	TotalCount int64 `json:"totalCount"`

	RequestedCount int64 `json:"requestedCount"`
}
