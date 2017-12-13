package bestbuy

type SearchResponse struct {
	From         int          `json:"from"`
	To           int          `json:"to"`
	Total        int          `json:"total"`
	CurrentPage  int          `json:"currentPage"`
	TotalPages   int          `json:"totalPages"`
	QueryTime    string       `json:"queryTime"`
	TotalTime    string       `json:"totalTime"`
	Partial      bool         `json:"partial"`
	CanonicalUrl string       `json:"canonicalUrl"`
	Products     []SearchItem `json:"products"`
}
