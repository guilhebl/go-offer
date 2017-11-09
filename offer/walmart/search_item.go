package walmart

type SearchItem struct {
	ItemId             int     `json:"itemId"`
	ParentItemId       int     `json:"parentItemId"`
	Upc                string  `json:"upc"`
	Name               string  `json:"name"`
	SalePrice          float32 `json:"salePrice", omitempty`
	CategoryPath       string  `json:"categoryPath"`
	ShortDescription   string  `json:"shortDescription", omitempty`
	LongDescription    string  `json:"longDescription", omitempty`
	ThumbnailImage     string  `json:"thumbnailImage", omitempty`
	MediumImage        string  `json:"mediumImage", omitempty`
	LargeImage         string  `json:"largeImage", omitempty`
	ProductTrackingUrl string  `json:"productTrackingUrl"`
	ModelNumber        string  `json:"modelNumber", omitempty`
	ProductUrl         string  `json:"productUrl"`
	CustomerRating     string  `json:"CustomerRating", omitempty`
	NumReviews         int     `json:"numReviews", omitempty`
}
