package model

type ListRequest struct {
	SearchColumns []NameValue `json:"searchColumns"`
	SortOrder   string `json:"sortOrder"`
	Page        int    `json:"page"`
	RowsPerPage int    `json:"rowsPerPage"`
}

func (r *ListRequest) Map() map[string]string {
	m := make(map[string]string)

	for _, p := range r.SearchColumns {
		m[p.Name] = p.Value
	}
	return m
}
