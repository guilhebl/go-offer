package model

type Offer struct {
	Id        string    	 `json:"id"`
	Upc 	  string    	 `json:"upc"`
	Name      string    	 `json:"name"`
	PartyName string    	 `json:"partyName"`
	SemanticName string 	 `json:"semanticName"`
	MainImageFileUrl string  `json:"mainImageFileUrl"`
	PartyImageFileUrl string `json:"partyImageFileUrl"`
	Price float64 			 `json:"price"`
	ProductCategory string 	 `json:"productCategory"`
	Rating float32 			 `json:"rating"`
	NumReviews int			 `json:"numReviews"`
}

type Offers []Offer