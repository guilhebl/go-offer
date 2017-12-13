package bestbuy

type TrendingResponse struct {
	Metadata Metadata       `json:"metadata"`
	Results  []TrendingItem `json:"results"`
}
