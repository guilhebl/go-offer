package model

import (
	"errors"
	"regexp"
)

type ListRequest struct {
	SearchColumns []NameValue `json:"searchColumns"`
	SortBy        string      `json:"sortBy"`
	SortOrder     string      `json:"sortOrder"`
	Page          int         `json:"page"`
	RowsPerPage   int         `json:"rowsPerPage"`
}

func NewListRequest(searchColumns []NameValue, sortBy, sortOrder string, page, rowsPerPage int) *ListRequest {
	o := &ListRequest{
		SearchColumns: searchColumns,
		SortBy:        sortBy,
		SortOrder:     sortOrder,
		Page:          page,
		RowsPerPage:   rowsPerPage,
	}
	return o
}
func NewEmptyListRequest(rowsPerPage int) *ListRequest {
	o := &ListRequest{
		SearchColumns: make([]NameValue, 0),
		SortBy:        "",
		SortOrder:     "",
		Page:          1,
		RowsPerPage:   rowsPerPage,
	}
	return o
}

func (r *ListRequest) Map() (map[string]string, error) {
	if !isValid(r) {
		return nil, errors.New(InvalidRequest)
	}

	m := make(map[string]string)

	for _, p := range r.SearchColumns {
		m[p.Name] = p.Value
	}
	return m, nil
}

// Checks if request is valid
func isValid(r *ListRequest) bool {

	// check valid page
	if r.Page <= 0 || r.RowsPerPage <= 0 {
		return false
	}

	// check valid sort
	if r.SortBy != "" {
		if match, _ := regexp.MatchString("asc|desc", r.SortOrder); !match {
			return false
		}
	}

	return true
}
