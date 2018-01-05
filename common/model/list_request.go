package model

import "errors"

type ListRequest struct {
	SearchColumns []NameValue `json:"searchColumns"`
	SortOrder     string      `json:"sortOrder"`
	Page          int         `json:"page"`
	RowsPerPage   int         `json:"rowsPerPage"`
}

func (r *ListRequest) Map() (map[string]string, error) {
	if !isValid(r) {
		return nil, errors.New("invalid request")
	}

	m := make(map[string]string)

	for _, p := range r.SearchColumns {
		m[p.Name] = p.Value
	}
	return m, nil
}

// Checks if request is valid
func isValid(r *ListRequest) bool {
	if r.Page <= 0 || r.RowsPerPage <= 0 {
		return false
	}

	if r.SortOrder != "asc" && r.SortOrder != "desc" {
		return false
	}
	return true
}
