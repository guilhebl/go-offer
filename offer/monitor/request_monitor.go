package monitor

import (
	"sync"
	"github.com/guilhebl/go-offer/common/config"
	"github.com/guilhebl/go-offer/common/model"
	"time"
)

// RequestMonitor is a singleton responsible for controlling the outbound calls to Marketplace providers
// controlling volume of calls being made to the external marketplace environment, making sure number
// of calls per second are within the limits and boundaries of each provider API.
type RequestMonitor struct {
	lastCalls map[string]int64
	waitIntervals map[string]int64
}

var instance *RequestMonitor
var once sync.Once

func GetInstance() *RequestMonitor {
	once.Do(func() {

		walmartWaitInterval := config.GetIntProperty("walmartRequestWaitIntervalMilis")
		eBayWaitInterval := config.GetIntProperty("eBayRequestWaitIntervalMilis")
		amazonWaitInterval := config.GetIntProperty("amazonRequestWaitIntervalMilis")
		bestBuyWaitInterval := config.GetIntProperty("bestbuyRequestWaitIntervalMilis")

		lc := map[string]int64{
			model.Walmart: 0,
			model.Ebay: 0,
			model.Amazon: 0,
			model.BestBuy: 0,
		}

		wi := map[string]int64{
			model.Walmart: walmartWaitInterval,
			model.Ebay: eBayWaitInterval,
			model.Amazon: amazonWaitInterval,
			model.BestBuy: bestBuyWaitInterval,
		}

		instance = &RequestMonitor{lc,wi	}
	})
	return instance
}

// Checks if waitIntervalMilis has passed since last Call.
// Difference in miliseconds (1*1000) = 1 second
//
func IsServiceAvailable(name string) bool {
	now := timeMillis()
	s := GetInstance()

	if now - s.lastCalls[name] >= s.waitIntervals[name] {
		s.lastCalls[name] = now
		return true
	}

	return false
}

// gets time in millis
func timeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}