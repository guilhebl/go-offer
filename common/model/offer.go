package model

// represents an offer for a Product or service
type Offer struct {
	Id                string  `json:"id"`
	Upc               string  `json:"upc"`
	Name              string  `json:"name"`
	PartyName         string  `json:"partyName"`
	SemanticName      string  `json:"semanticName"`
	MainImageFileUrl  string  `json:"mainImageFileUrl"`
	PartyImageFileUrl string  `json:"partyImageFileUrl"`
	ProductCategory   string  `json:"productCategory"`
	Price             float32 `json:"price"`
	Rating            float32 `json:"rating"`
	NumReviews        int     `json:"numReviews"`
}

func NewOffer(id, upc, name, partyName, semanticName, mainImageUrl, partyImageUrl, productCategory string, price, rating float32, numReviews int) *Offer {
	o := &Offer{
		Id:                id,
		Upc:               upc,
		Name:              name,
		PartyName:         partyName,
		SemanticName:      semanticName,
		MainImageFileUrl:  mainImageUrl,
		PartyImageFileUrl: partyImageUrl,
		ProductCategory:   productCategory,
		Price:             price,
		Rating:            rating,
		NumReviews:        numReviews,
	}
	return o
}
