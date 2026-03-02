package dto

type Pagination struct {
	CurrentPage     int64   `json:"current_page,omitempty" extensions:"x-order=0"`
	CurrentElements int64   `json:"current_elements,omitempty" extensions:"x-order=1"`
	TotalPages      int64   `json:"total_pages,omitempty" extensions:"x-order=2"`
	TotalElements   int64   `json:"total_elements,omitempty" extensions:"x-order=3"`
	SortBy          string  `json:"sort_by,omitempty" extensions:"x-order=4"`
	SortDir         string  `json:"sort_dir,omitempty" extensions:"x-order=5"`
	CursorStart     *string `json:"cursor_start,omitempty" extensions:"x-order=6"`
	CursorEnd       *string `json:"cursor_end,omitempty" extensions:"x-order=7"`

	Page  string `json:"page,omitempty"`
	Limit string `json:"limit,omitempty"`
	Total int64  `json:"total,omitempty"`
}
