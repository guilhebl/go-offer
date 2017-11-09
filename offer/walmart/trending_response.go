package walmart

import "time"

type TrendingResponse struct {
	Time  time.Time    `json:"time"`
	Items []SearchItem `json:"items"`
}
