package monitor

import (
	"github.com/guilhebl/go-offer/common/config"
	"github.com/guilhebl/go-offer/common/model"
	"sync"
	"time"
)

// RequestMonitor is a singleton responsible for controlling the outbound calls to Marketplace providers
// controlling volume of calls being made to the external marketplace environment, making sure number
// of calls per second are within the limits and boundaries of each provider API.
type RequestMonitor struct {
	lastCalls     sync.Map
	waitIntervals map[string]int64
}

// checks if this api is available after waiting a certain "treshold" this avoids flooding this external resource with
// calls, certain external APIs have quotas as max X requests per second.
func (r *RequestMonitor) isServiceAvailable(name string, timestamp int64) bool {
	s := GetInstance()
	lastCall, ok := instance.lastCalls.Load(name)

	if ok && timestamp-lastCall.(int64) >= s.waitIntervals[name] {
		instance.lastCalls.Store(name, timestamp)
		return true
	}

	return false
}

var instance *RequestMonitor
var once sync.Once

func GetInstance() *RequestMonitor {
	once.Do(func() {

		walmartWaitInterval := config.GetIntProperty("walmartRequestWaitIntervalMilis")
		eBayWaitInterval := config.GetIntProperty("eBayRequestWaitIntervalMilis")
		amazonWaitInterval := config.GetIntProperty("amazonRequestWaitIntervalMilis")
		bestBuyWaitInterval := config.GetIntProperty("bestbuyRequestWaitIntervalMilis")

		wi := map[string]int64{
			model.Walmart: walmartWaitInterval,
			model.Ebay:    eBayWaitInterval,
			model.Amazon:  amazonWaitInterval,
			model.BestBuy: bestBuyWaitInterval,
		}

		instance = &RequestMonitor{sync.Map{}, wi}

		// store initial last calls timestamps
		instance.lastCalls.Store(model.Walmart, int64(0))
		instance.lastCalls.Store(model.Ebay, int64(0))
		instance.lastCalls.Store(model.Amazon, int64(0))
		instance.lastCalls.Store(model.BestBuy, int64(0))
	})
	return instance
}

// Checks if waitIntervalMilis has passed since last Call.
// Difference in miliseconds (1*1000) = 1 second
func IsServiceAvailable(name string) bool {
	now := timeMillis()
	return instance.isServiceAvailable(name, now)
}

// gets time in millis
func timeMillis() int64 {
	return time.Now().UnixNano() / 1000000
}
