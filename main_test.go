package main

import (
	"github.com/guilhebl/go-offer/offer"
	"github.com/stretchr/testify/assert"
	"gopkg.in/jarcoal/httpmock.v1"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

var app offer.Module

func init() {
	runtime.LockOSThread()
}

func TestMain(m *testing.M) {
	setup()
	go func() {
		exitVal := m.Run()
		teardown()
		os.Exit(exitVal)
	}()

	log.Println("setting up test server...")
	run()
}

func setup() {
	log.Println("SETUP")
}

func teardown() {
	log.Println("TEARDOWN")
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func readFile(path string) []byte {
	absPath, _ := filepath.Abs("./" + path)
	dat, err := ioutil.ReadFile(absPath)
	check(err)
	return dat
}

const (
	WalmartTrendingUrl = "http://api.walmartlabs.com/v1/trends"
	WalmartSearchUrl   = "https://api.walmartlabs.com/search"
	BestBuyTrendingUrl = "https://api.bestbuy.com/beta/products/trendingViewed"
	BestBuySearchUrl   = "https://api.bestbuy.com/v1/products"
	EbaySearchUrl      = "http://svcs.ebay.com/services/search/FindingService/v1"
)

// returns the bytes of a corresponding mock API call for an external resource
func getJsonMockBytes(url string) []byte {
	switch url {
	case WalmartTrendingUrl:
		return readFile("offer/walmart/walmart_sample_trending_response.json")
	case WalmartSearchUrl:
		return readFile("offer/walmart/walmart_sample_search_response.json")
	case BestBuyTrendingUrl:
		return readFile("offer/bestbuy/bestbuy_sample_trending_response.json")
	case BestBuySearchUrl:
		return readFile("offer/bestbuy/bestbuy_sample_search_response.json")
	case EbaySearchUrl:
		return readFile("offer/ebay/ebay_sample_search_response.json")

	default:
		return nil
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	offer.GetInstance().Router.ServeHTTP(rr, req)
	return rr
}

func registerMockResponder(httpMethod, apiUrl string, status int) {
	log.Printf("Mocking Search: %s %d - %s", httpMethod, status, apiUrl)
	httpmock.RegisterResponder(httpMethod, apiUrl, httpmock.NewBytesResponder(status, getJsonMockBytes(apiUrl)))
}

// Tests basic Search (no keywords) that returns trending results from external APIs
func TestSearch(t *testing.T) {

	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponder("GET", WalmartTrendingUrl, 200)
	registerMockResponder("GET", BestBuyTrendingUrl, 200)
	registerMockResponder("GET", EbaySearchUrl, 200)

	// call our local server API
	endpoint := "http://localhost:8080/"
	req, _ := http.NewRequest("GET", endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()
	assert.True(t, strings.HasPrefix(body, `{"list":[{"`))

	walmartSnippet := `{"id":"348726849","upc":"816586026705","name":"Best Choice Products 6' Exercise Tri-Fold Gym Mat For Gymnastics, Aerobics, Yoga, Martial Arts - Pink","partyName":"walmart.com"`
	assert.True(t, strings.Contains(body, walmartSnippet))

	bestBuySnippet := `{"id":"5714687","upc":"","name":"Alienware - Aurora R6 Desktop - Intel Core i7 - 16GB Memory - NVIDIA GeForce GTX 1070 - 256GB Solid State Drive + 1TB Hard Drive - Silver","partyName":"bestbuy.com"`
	assert.True(t, strings.Contains(body, bestBuySnippet))

	ebaySnippet := `,{"id":"262954865748","upc":"","name":"The Elder Scrolls V: Skyrim - Greatest Hits PS3 [Brand New]","partyName":"ebay.com"`
	assert.True(t, strings.Contains(body, ebaySnippet))

	// get the amount of calls for the registered responders
	assertCallsMade(t, "GET", WalmartTrendingUrl, 1)
	assertCallsMade(t, "GET", BestBuyTrendingUrl, 1)
	assertCallsMade(t, "GET", EbaySearchUrl, 1)
}

func assertCallsMade(t *testing.T, httpMethod, url string, expected int) {
	info := httpmock.GetCallCountInfo()
	count := info[httpMethod+" "+url]
	assert.Equal(t, expected, count)
	log.Printf("Total External API Calls made to %s: %d", url, count)
}
