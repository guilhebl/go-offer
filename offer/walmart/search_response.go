package walmart

type SearchResponse struct {
	Query string `json:"query"`
	Sort string `json:"sort"`
	ResponseGroup string `json:"responseGroup"`
	TotalResults int `json:"totalResults"`
	Start int `json:"start"`
	NumItems int `json:"numItems"`
	Items []SearchItem `json:"items"`
}
