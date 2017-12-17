package ebay

import "time"

type SearchItem struct {
	ItemID          []string `json:"itemId"`
	Title           []string `json:"title"`
	GlobalID        []string `json:"globalId"`
	PrimaryCategory []struct {
		CategoryID   []string `json:"categoryId"`
		CategoryName []string `json:"categoryName"`
	} `json:"primaryCategory"`
	GalleryURL    []string `json:"galleryURL"`
	ViewItemURL   []string `json:"viewItemURL"`
	PaymentMethod []string `json:"paymentMethod"`
	AutoPay       []string `json:"autoPay"`
	PostalCode    []string `json:"postalCode"`
	Location      []string `json:"location"`
	Country       []string `json:"country"`
	ShippingInfo  []struct {
		ShippingServiceCost []struct {
			CurrencyID string `json:"@currencyId"`
			Value      string `json:"__value__"`
		} `json:"shippingServiceCost"`
		ShippingType            []string `json:"shippingType"`
		ShipToLocations         []string `json:"shipToLocations"`
		ExpeditedShipping       []string `json:"expeditedShipping"`
		OneDayShippingAvailable []string `json:"oneDayShippingAvailable"`
		HandlingTime            []string `json:"handlingTime"`
	} `json:"shippingInfo"`
	SellingStatus []struct {
		CurrentPrice []struct {
			CurrencyID string `json:"@currencyId"`
			Value      string `json:"__value__"`
		} `json:"currentPrice"`
		ConvertedCurrentPrice []struct {
			CurrencyID string `json:"@currencyId"`
			Value      string `json:"__value__"`
		} `json:"convertedCurrentPrice"`
		SellingState []string `json:"sellingState"`
		TimeLeft     []string `json:"timeLeft"`
	} `json:"sellingStatus"`
	ListingInfo []struct {
		BestOfferEnabled  []string    `json:"bestOfferEnabled"`
		BuyItNowAvailable []string    `json:"buyItNowAvailable"`
		StartTime         []time.Time `json:"startTime"`
		EndTime           []time.Time `json:"endTime"`
		ListingType       []string    `json:"listingType"`
		Gift              []string    `json:"gift"`
		WatchCount        []string    `json:"watchCount"`
	} `json:"listingInfo"`
	ReturnsAccepted []string `json:"returnsAccepted"`
	Condition       []struct {
		ConditionID          []string `json:"conditionId"`
		ConditionDisplayName []string `json:"conditionDisplayName"`
	} `json:"condition"`
	IsMultiVariationListing []string `json:"isMultiVariationListing"`
	PictureURLLarge         []string `json:"pictureURLLarge,omitempty"`
	TopRatedListing         []string `json:"topRatedListing"`
	ProductID               []struct {
		Type  string `json:"@type"`
		Value string `json:"__value__"`
	} `json:"productId,omitempty"`
	CharityID []string `json:"charityId,omitempty"`
}
