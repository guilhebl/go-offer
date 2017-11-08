package model

// represents a List of offers response with a summary
type OfferList struct {
	List []Offer `json:"list"`
	Summary `json:"summary"`
}

type Summary struct {
	Page       int `json:"page"`
	PageCount  int `json:"pageCount"`
	TotalCount int `json:"totalCount"`
}

func NewOfferList(list []Offer, page int, pageCount int, total int) *OfferList {
	o := &OfferList{
		List: list,
		Summary: Summary{Page: page, PageCount: pageCount, TotalCount: total},
	}
	return o
}