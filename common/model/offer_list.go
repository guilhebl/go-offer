package model

type OfferList struct {
	List []Offer `json:"list"`
	Summary struct {
		Page       int `json:"page"`
		PageCount  int `json:"pageCount"`
		TotalCount int `json:"totalCount"`
	} `json:"summary"`
}
