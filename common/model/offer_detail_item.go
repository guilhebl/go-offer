package model

type OfferDetailItem struct {
	PartyName         string  `json:"partyName"`
	SemanticName      string  `json:"semanticName"`
	PartyImageFileUrl string  `json:"partyImageFileUrl"`
	Price             float32 `json:"price"`
	Rating            float32 `json:"rating"`
	NumReviews        int     `json:"numReviews"`
}

func NewOfferDetailItem(partyName, semanticName, partyImageFileUrl string, price, rating float32, numReviews int) *OfferDetailItem {
	od := &OfferDetailItem{
		PartyName:         partyName,
		SemanticName:      semanticName,
		PartyImageFileUrl: partyImageFileUrl,
		Price:             price,
		Rating:            rating,
		NumReviews:        numReviews,
	}
	return od
}
