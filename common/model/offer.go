package model

// represents an offer for a Product or service offered by a provider in a specific moment of Time
// Id represents the internal Id of this offer to uniquely address this offer within the system
// External Id represents the external Id used by the provider to address this entity in their system
type Offer struct {
	Id                string  `json:"id"`
	ExternalId        string  `json:"externalId"`
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

func NewOffer(id, externalId, upc, name, partyName, semanticName, mainImageUrl, partyImageUrl, productCategory string, price, rating float32, numReviews int) *Offer {
	o := &Offer{
		Id:                id,
		ExternalId:        externalId,
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

// represents an Offer with Rank information:
// num keywords - the number of keywords found out of a group of N keywords
// total - total matches of keywords found in this offer
// lowest index - the lowest index found
type OfferKeywordRank struct {
	Offer        Offer
	NumKeywords  int
	TotalMatches int
	LowestIndex  int
}

func NewOfferKeywordRank(o Offer, nk, tm, i int) *OfferKeywordRank {
	or := &OfferKeywordRank{
		Offer:        o,
		NumKeywords:  nk,
		TotalMatches: tm,
		LowestIndex:  i,
	}
	return or
}
