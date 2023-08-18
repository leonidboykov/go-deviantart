package deviantart

// OffsetParams params for offset-based pagination.
type OffsetParams struct {
	// The pagination offset.
	Offset uint32 `url:"limit,omitempty"`

	// The pagination limit.
	Limit uint32 `url:"limit,omitempty"`
}

type OffsetResponse[T any] struct {
	Results    []T    `json:"results"`
	HasMore    bool   `json:"has_more"`
	NextOffset uint32 `json:"next_offset,omitempty"`

	// This field is used in some endpoints with a query parameter.
	// TODO: Use separate struct and API method (?).
	EstimatedTotal uint32 `json:"estimated_total,omitempty"`
}

// CursorParams params for cursor-based pagination.
type CursorParams struct {
	Cursor string `url:"cursor,omitempty"`
}

type CursorResponse[T any] struct {
	Results    []T    `json:"results"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor"`
	PrevCursor string `json:"prev_cursor"`
}

// singleResponse represents response without pagination.
type singleResponse[T any] struct {
	Results []T `json:"results"`
}
