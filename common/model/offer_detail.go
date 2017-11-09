package model

type OfferDetail struct {
	Offer              Offer             `json:"offer"`
	Description        string            `json:"description"`
	Attributes         []NameValue       `json:"attributes"`
	ProductDetailItems []OfferDetailItem `json:"productDetailItems"`
}

func NewOfferDetail(o Offer, desc string, attrs map[string]string, items []OfferDetailItem) *OfferDetail {
	od := &OfferDetail{
		Offer:              o,
		Description:        desc,
		Attributes:         NewNameValues(attrs),
		ProductDetailItems: items,
	}
	return od
}
