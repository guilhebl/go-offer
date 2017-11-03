package model

type OfferDetail struct {
	Offer Offer `json:"offer"`
	Description string `json:"description"`
	Attributes []NameValue `json:"attributes"`
	ProductDetailItems []OfferDetailItem`json:"productDetailItems"`
}
