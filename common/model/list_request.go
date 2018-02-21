package model

import (
	"regexp"
	"strconv"
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

// builds a map out of request params
func (r *ListRequest) Map() map[string]string {
	m := make(map[string]string)

	// put search columns
	for _, p := range r.SearchColumns {
		m[p.Name] = p.Value
	}

	// put other fields
	m[SortBy] = r.SortBy
	m[SortOrder] = r.SortOrder
	m[Page] = strconv.Itoa(r.Page)
	m[RowsPerPage] = strconv.Itoa(r.RowsPerPage)

	return m
}

// Checks if request is valid
func (r *ListRequest) IsValid() bool {

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
