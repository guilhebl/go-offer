package main

type OfferDetailItem struct {
	PartyName string    	 `json:"partyName"`
	SemanticName string 	 `json:"semanticName"`
	PartyImageFileUrl string `json:"partyImageFileUrl"`
	Price float64 			 `json:"price"`
	Rating float32 			 `json:"rating"`
	NumReviews int			 `json:"numReviews"`
}
