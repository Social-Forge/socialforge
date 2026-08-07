package helpers

import (
	"strconv"

	"github.com/gofiber/fiber/v3"
)

// Default + bounds for list pagination.
const (
	DefaultPage    = 1
	DefaultPerPage = 20
	MaxPerPage     = 100
)

// PageMeta is the standard pagination envelope returned in ApiResponse.meta.
type PageMeta struct {
	Page       int   `json:"page"`
	PerPage    int   `json:"per_page"`
	Total      int64 `json:"total"`
	TotalPages int   `json:"total_pages"`
	HasMore    bool  `json:"has_more"`
}

// PageParams carries a parsed page/limit/offset for a list query.
type PageParams struct {
	Page   int
	Limit  int
	Offset int
}

// ParsePageParams reads `page` and `per_page` (alias `limit`) from the query,
// clamping to sane bounds.
func ParsePageParams(c fiber.Ctx) PageParams {
	page := atoiDefault(c.Query("page"), DefaultPage)
	if page < 1 {
		page = 1
	}
	perPage := c.Query("per_page")
	if perPage == "" {
		perPage = c.Query("limit")
	}
	limit := atoiDefault(perPage, DefaultPerPage)
	if limit < 1 {
		limit = DefaultPerPage
	}
	if limit > MaxPerPage {
		limit = MaxPerPage
	}
	return PageParams{Page: page, Limit: limit, Offset: (page - 1) * limit}
}

// NewPageMeta builds the pagination metadata from the params + total row count.
func NewPageMeta(p PageParams, total int64) PageMeta {
	totalPages := 0
	if p.Limit > 0 {
		totalPages = int((total + int64(p.Limit) - 1) / int64(p.Limit))
	}
	return PageMeta{
		Page:       p.Page,
		PerPage:    p.Limit,
		Total:      total,
		TotalPages: totalPages,
		HasMore:    p.Page < totalPages,
	}
}

// RespondWithMeta is like Respond but also attaches pagination (or any) meta.
func RespondWithMeta(c fiber.Ctx, status int, message string, payload interface{}, meta interface{}) error {
	success := status >= 200 && status < 300
	response := ApiResponse{
		Status:  status,
		Success: success,
		Message: message,
		Meta:    meta,
	}
	if success {
		response.Data = payload
	} else {
		response.Error = payload
	}
	return c.Status(status).JSON(response)
}

func atoiDefault(s string, def int) int {
	if s == "" {
		return def
	}
	n, err := strconv.Atoi(s)
	if err != nil {
		return def
	}
	return n
}
