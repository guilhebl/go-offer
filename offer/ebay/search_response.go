package ebay

type SearchResponse struct {
	FindItemsByKeywordsResponse []SearchResponseResult `json:"findItemsByKeywordsResponse"`
}
