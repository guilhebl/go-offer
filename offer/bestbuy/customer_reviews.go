package bestbuy

type CustomerReviews struct {
	AverageScore float32 `json:"averageScore", omitempty`
	Count        int     `json:"count", omitempty`
}
