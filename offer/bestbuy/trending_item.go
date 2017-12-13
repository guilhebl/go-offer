package bestbuy

type TrendingItem struct {
	CustomerReviews CustomerReviews `json:"customerReviews"`
	Descriptions    Descriptions    `json:"descriptions"`
	Images          ProductImages   `json:"images"`
	Names           ProductNames    `json:"names"`
	Links           ProductLinks    `json:"links"`
	Prices          ProductPrices   `json:"prices"`
	Rank            int             `json:"rank"`
	Sku             string          `json:"sku"`
}
