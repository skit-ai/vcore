package handlers

import (
	"fmt"
	"strings"
)

// PaginatedList represents a paginated list of data items.
type PaginatedList struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	TotalPages int         `json:"total_pages"`
	TotalItems int         `json:"total_items"`
	Items      interface{} `json:"items"`
	First      string      `json:"first,omitempty"`
	Prev       string      `json:"prev,omitempty"`
	Next       string      `json:"next,omitempty"`
	Last       string      `json:"last,omitempty"`
}

// NewPaginatedList creates a new Paginated instance.
// The page parameter is 1-based and refers to the current page index/number.
// The perPage parameter refers to the number of items on each page.
// And the total parameter specifies the total number of data items.
// If total is less than 0, it means total is unknown.
func NewPaginatedList(page, pageSize, totalItems int) *PaginatedList {
	if pageSize < 1 {
		pageSize = 10
	}
	totalPages := -1
	if totalItems >= 0 {
		totalPages = (totalItems + pageSize - 1) / pageSize
		if page > totalPages {
			page = totalPages
		}
	}
	if page < 1 {
		page = 1
	}

	return &PaginatedList{
		Page:       page,
		PageSize:   pageSize,
		TotalItems: totalItems,
		TotalPages: totalPages,
	}
}

// Offset returns the OFFSET value that can be used in a SQL statement.
func (p *PaginatedList) Offset() int {
	return (p.Page - 1) * p.PageSize
}

// Limit returns the LIMIT value that can be used in a SQL statement.
func (p *PaginatedList) Limit() int {
	return p.PageSize
}

// Build page links to first, prev, next and last page
func (p *PaginatedList) BuildPageLinks(baseUrl string, defaultPageSize int) {
	links := p.buildLinks(baseUrl, defaultPageSize)
	if links[0] != "" {
		p.First = links[0]
		p.Prev = links[1]
	}
	if links[2] != "" {
		p.Next = links[2]
		if links[3] != "" {
			p.Last = links[3]
		}
	}
}

// buildLinks returns the first, prev, next, and last links corresponding to the pagination.
// A link could be an empty string if it is not needed.
// For example, if the pagination is at the first page, then both first and prev links
// will be empty.
func (p *PaginatedList) buildLinks(baseUrl string, defaultPageSize int) [4]string {
	var links [4]string
	pageCount := p.TotalPages
	page := p.Page
	if pageCount >= 0 && page > pageCount {
		page = pageCount
	}
	if strings.Contains(baseUrl, "?") {
		baseUrl += "&"
	} else {
		baseUrl += "?"
	}
	if page > 1 {
		links[0] = fmt.Sprintf("%vpage=%v", baseUrl, 1)
		links[1] = fmt.Sprintf("%vpage=%v", baseUrl, page-1)
	}
	if pageCount >= 0 && page < pageCount {
		links[2] = fmt.Sprintf("%vpage=%v", baseUrl, page+1)
		links[3] = fmt.Sprintf("%vpage=%v", baseUrl, pageCount)
	} else if pageCount < 0 {
		links[2] = fmt.Sprintf("%vpage=%v", baseUrl, page+1)
	}
	if pageSize := p.PageSize; pageSize != defaultPageSize {
		for i := 0; i < 4; i++ {
			if links[i] != "" {
				links[i] += fmt.Sprintf("&page_size=%v", pageSize)
			}
		}
	}

	return links
}
