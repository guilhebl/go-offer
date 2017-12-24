package main

import (
	"bytes"
	"github.com/guilhebl/go-offer/common/model"
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

// TestMain acts as a Functional Test, it starts the server and actually is a running instance of the whole app, whenever
// the server triggers external calls these are captured and mock data is returned instead, this way all functional tests can be done
// without any external dependencies (offline mode)
func TestMain(m *testing.M) {
	setup()
	go func() {
		exitVal := m.Run()
		teardown()
		os.Exit(exitVal)
	}()

	log.Println("setting up test server...")
	run(model.Test)
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
	WalmartSearchUrl   = "http://api.walmartlabs.com/v1/search"
	BestBuyTrendingUrl = "https://api.bestbuy.com/beta/products/trendingViewed"
	BestBuySearchUrl   = "https://api.bestbuy.com/v1/products(search=skyrim)"
	EbaySearchUrl      = "http://svcs.ebay.com/services/search/FindingService/v1"
	AmazonSearchUrl    = "https://webservices.amazon.com/onca/xml"
)

// returns the bytes of a corresponding mock API call for an external resource for the 'Trending' API CALL
func getJsonBytesTrendingMock(url string) []byte {
	switch url {
	case WalmartTrendingUrl:
		return readFile("offer/walmart/walmart_sample_trending_response.json")
	case BestBuyTrendingUrl:
		return readFile("offer/bestbuy/bestbuy_sample_trending_response.json")
	case EbaySearchUrl:
		return readFile("offer/ebay/ebay_sample_trending_response.json")
	case AmazonSearchUrl:
		return readFile("offer/amazon/amazon_sample_trending_response.xml")

	default:
		return nil
	}
}

// returns the bytes of a corresponding mock API call for an external resource for the 'Search' API CALL
func getJsonBytesSearchMock(url string) []byte {
	switch url {
	case WalmartSearchUrl:
		return readFile("offer/walmart/walmart_sample_search_response.json")
	case BestBuySearchUrl:
		return readFile("offer/bestbuy/bestbuy_sample_search_response.json")
	case EbaySearchUrl:
		return readFile("offer/ebay/ebay_sample_search_response.json")
	case AmazonSearchUrl:
		return readFile("offer/amazon/amazon_sample_search_response.xml")

	default:
		return nil
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	offer.GetInstance().Router.ServeHTTP(rr, req)
	return rr
}

func assertCallsMade(t *testing.T, httpMethod, url string, expected int) {
	info := httpmock.GetCallCountInfo()
	count := info[httpMethod+" "+url]
	assert.Equal(t, expected, count)
	log.Printf("Total External API Calls made to %s: %d", url, count)
}

func registerMockResponder(httpMethod, apiUrl, apiType string, status int) {
	log.Printf("Mocking Search: %s %d - %s", httpMethod, status, apiUrl)

	switch apiType {
	case model.Trending:
		httpmock.RegisterResponder(httpMethod, apiUrl, httpmock.NewBytesResponder(status, getJsonBytesTrendingMock(apiUrl)))
	case model.Search:
		httpmock.RegisterResponder(httpMethod, apiUrl, httpmock.NewBytesResponder(status, getJsonBytesSearchMock(apiUrl)))
	}
}

// Tests basic Search (no keywords) that returns trending results from external APIs
func TestSearch(t *testing.T) {

	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponder(http.MethodGet, WalmartTrendingUrl, model.Trending, 200)
	registerMockResponder(http.MethodGet, BestBuyTrendingUrl, model.Trending, 200)
	registerMockResponder(http.MethodGet, EbaySearchUrl, model.Trending, 200)
	registerMockResponder(http.MethodGet, AmazonSearchUrl, model.Trending, 200)

	// call our local server API
	endpoint := "http://localhost:8080/"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()

	assert.True(t, strings.HasPrefix(body, `{"list":[{"`))

	walmartSnippet := `{"id":"348726849","upc":"816586026705","name":"Best Choice Products 6' Exercise Tri-Fold Gym Mat For Gymnastics, Aerobics, Yoga, Martial Arts - Pink","partyName":"walmart.com"`
	assert.True(t, strings.Contains(body, walmartSnippet))

	ebaySnippet := `{"id":"282629961650","upc":"","name":"Reverb Cross Men s Running Shoes","partyName":"ebay.com"`
	assert.True(t, strings.Contains(body, ebaySnippet))

	amazonSnippet := `{"id":"B0743W4Y75","upc":"701649356113","name":"Bluetooth Smart Watch with Camera, Aosmart B23 Smart Watch for Android Smartphones (White)","partyName":"amazon.com"`
	assert.True(t, strings.Contains(body, amazonSnippet))

	bestBuySnippet := `{"id":"5714687","upc":"","name":"Alienware - Aurora R6 Desktop - Intel Core i7 - 16GB Memory - NVIDIA GeForce GTX 1070 - 256GB Solid State Drive + 1TB Hard Drive - Silver","partyName":"bestbuy.com"`
	assert.True(t, strings.Contains(body, bestBuySnippet))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, WalmartTrendingUrl, 1)
	assertCallsMade(t, http.MethodGet, BestBuyTrendingUrl, 1)
	assertCallsMade(t, http.MethodGet, EbaySearchUrl, 1)
	assertCallsMade(t, http.MethodGet, AmazonSearchUrl, 1)
}

// Tests Search with keywords that returns results from external APIs
func TestSearchWithKeywords(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponder(http.MethodGet, WalmartSearchUrl, model.Search, 200)
	registerMockResponder(http.MethodGet, BestBuySearchUrl, model.Search, 200)
	registerMockResponder(http.MethodGet, EbaySearchUrl, model.Search, 200)
	registerMockResponder(http.MethodGet, AmazonSearchUrl, model.Search, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers"
	var jsonRequest = []byte(`{"searchColumns":[{"name":"name","value":"skyrim"}],"sortOrder":"asc","page":1,"rowsPerPage":10}`)

	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()

	assert.True(t, strings.HasPrefix(body, `{"list":[{"`))

	walmartSnippet := `{"id":"53966162","upc":"093155171244","name":"Skyrim Special Edition (Xbox One)","partyName":"walmart.com",`
	assert.True(t, strings.Contains(body, walmartSnippet))

	bestBuySnippet := `{"id":"5626200","upc":"600603210488","name":"The Elder Scrolls V: Skyrim Special Edition Best Buy Exclusive Dragonborn Bundle - Xbox One","partyName":"bestbuy.com"`
	assert.True(t, strings.Contains(body, bestBuySnippet))

	ebaySnippet := `{"id":"223482818","upc":"","name":"Elder Scrolls V: Skyrim - Special Edition With Bonus Steelbook Case PS4 ","partyName":"ebay.com"`
	assert.True(t, strings.Contains(body, ebaySnippet))

	amazonSnippet := `{"id":"B01GW8XJVU","upc":"093155171251","name":"The Elder Scrolls V: Skyrim - Special Edition - PlayStation 4","partyName":"amazon.com"`
	assert.True(t, strings.Contains(body, amazonSnippet))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, WalmartSearchUrl, 1)
	assertCallsMade(t, http.MethodGet, BestBuySearchUrl, 1)
	assertCallsMade(t, http.MethodGet, EbaySearchUrl, 1)
	assertCallsMade(t, http.MethodGet, AmazonSearchUrl, 1)
}
