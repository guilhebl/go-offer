package main

type ListRequest struct {
	SearchColumns []NameValue `json:"searchColumns"`
	SortOrder   string `json:"sortOrder"`
	Page        int    `json:"page"`
	RowsPerPage int    `json:"rowsPerPage"`
}