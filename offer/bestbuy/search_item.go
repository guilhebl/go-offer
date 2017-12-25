package bestbuy

type SearchItem struct {
	ProductId             int            `json:"productId"`
	Upc                   string         `json:"upc", omitempty`
	Sku                   int64          `json:"sku", omitempty`
	Name                  string         `json:"name"`
	SalePrice             float32        `json:"salePrice"`
	ReleaseDate           string         `json:"releaseDate", omitempty`
	Url                   string         `json:"url", omitempty`
	Image                 string         `json:"image", omitempty`
	ThumbnailImage        string         `json:"thumbnailImage", omitempty`
	Manufacturer          string         `json:"manufacturer", omitempty`
	Department            string         `json:"department", omitempty`
	CustomerReviewAverage float32        `json:"customerReviewAverage", omitempty`
	CustomerReviewCount   int            `json:"customerReviewCount", omitempty`
	CategoryPath          []CategoryPath `json:"categoryPath"`
}
