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

// TestMain builds an instance of the application in test mode to run e2e tests. Every external calls are intercepted through stubs
// that returns test JSON data instead of really calling external services.
// This way it's possible to build and run offline Functional Tests on top of the actual running stack and test multiple end-to-end scenarios.
func TestMain(m *testing.M) {
	setup()
	go func() {
		exitVal := m.Run()
		teardown()
		os.Exit(exitVal)
	}()

	log.Println("setting up test server...")
	run("test")
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
	WalmartTrendingUrl       = "http://api.walmartlabs.com/v1/trends"
	WalmartSearchUrl         = "http://api.walmartlabs.com/v1/search"
	WalmartGetDetailUrl      = "http://api.walmartlabs.com/v1/items/53966162"
	WalmartGetDetailByUpcUrl = "http://api.walmartlabs.com/v1/items"
	BestBuyTrendingUrl       = "https://api.bestbuy.com/beta/products/trendingViewed"
	BestBuySearchUrl         = "https://api.bestbuy.com/v1/products(search=skyrim)"
	BestBuyGetDetailUrl      = "https://api.bestbuy.com/v1/products(productId=5529006)"
	BestBuyGetDetailByUpcUrl = "https://api.bestbuy.com/v1/products(upc=065857174434)"
	EbaySearchUrl            = "http://svcs.ebay.com/services/search/FindingService/v1"
	EbayGetDetailUrl         = "http://svcs.ebay.com/services/search/FindingService/v1"
	AmazonSearchUrl          = "https://webservices.amazon.com/onca/xml"
	AmazonGetDetailUrl       = "https://webservices.amazon.com/onca/xml"
)

// returns the bytes of a corresponding mock API call for an external resource for the 'Trending' API CALL
func getJsonBytesTrendingMock(url string) []byte {
	switch url {
	case WalmartTrendingUrl:
		return readFile("offer/walmart/testdata/walmart_sample_trending_response.json")
	case BestBuyTrendingUrl:
		return readFile("offer/bestbuy/testdata/bestbuy_sample_trending_response.json")
	case EbaySearchUrl:
		return readFile("offer/ebay/testdata/ebay_sample_trending_response.json")
	case AmazonSearchUrl:
		return readFile("offer/amazon/testdata/amazon_sample_trending_response.xml")

	default:
		return nil
	}
}

// returns the bytes of a corresponding mock API call for an external resource for the 'Search' API CALL
func getJsonBytesSearchMock(url string) []byte {
	switch url {
	case WalmartSearchUrl:
		return readFile("offer/walmart/testdata/walmart_sample_search_response.json")
	case BestBuySearchUrl:
		return readFile("offer/bestbuy/testdata/bestbuy_sample_search_response.json")
	case EbaySearchUrl:
		return readFile("offer/ebay/testdata/ebay_sample_search_response.json")
	case AmazonSearchUrl:
		return readFile("offer/amazon/testdata/amazon_sample_search_response.xml")

	default:
		return nil
	}
}

// returns the bytes of a corresponding mock API call for an external resource for the 'Search No results' API CALL
func getJsonBytesSearchNoResultsMock(url string) []byte {
	switch url {
	case WalmartSearchUrl:
		return readFile("offer/walmart/testdata/walmart_search_no_results.json")
	case BestBuySearchUrl:
		return readFile("offer/bestbuy/testdata/bestbuy_search_no_results.json")
	case EbaySearchUrl:
		return readFile("offer/ebay/testdata/ebay_search_no_results.json")
	case AmazonSearchUrl:
		return readFile("offer/amazon/testdata/amazon_search_no_results.xml")

	default:
		return nil
	}
}

// returns the bytes of a corresponding mock API call for an external resource for the 'GetDetail' API CALL
func getJsonBytesGetDetailByIdMock(url string) []byte {
	switch url {
	case WalmartGetDetailUrl:
		return readFile("offer/walmart/testdata/walmart_sample_get_detail_by_id_response.json")
	case BestBuyGetDetailUrl:
		return readFile("offer/bestbuy/testdata/bestbuy_sample_get_detail_by_id_response.json")
	case EbayGetDetailUrl:
		return readFile("offer/ebay/testdata/ebay_sample_get_detail_by_id_response.json")
	case AmazonGetDetailUrl:
		return readFile("offer/amazon/testdata/amazon_sample_get_detail_by_id_response.xml")

	default:
		return nil
	}
}

// returns the bytes of a corresponding mock API call for an external resource for the 'GetDetail' API CALL By UPC
func getJsonBytesGetDetailByUpcMock(url string) []byte {
	switch url {
	case WalmartGetDetailByUpcUrl:
		return readFile("offer/walmart/testdata/walmart_sample_get_detail_by_upc_response.json")
	case BestBuyGetDetailByUpcUrl:
		return readFile("offer/bestbuy/testdata/bestbuy_sample_get_detail_by_upc_response.json")
	case EbayGetDetailUrl:
		return readFile("offer/ebay/testdata/ebay_find_by_upc.json")
	case AmazonGetDetailUrl:
		return readFile("offer/amazon/testdata/amazon_sample_get_detail_by_upc_response.xml")

	default:
		return nil
	}
}

// returns the bytes of a corresponding mock API call for an external resource for the 'GetDetail' API CALL By UPC With no Results
func getJsonBytesGetDetailByUpcNoResultsMock(url string) []byte {
	switch url {
	case WalmartGetDetailByUpcUrl:
		return readFile("offer/walmart/testdata/walmart_get_by_upc_not_found.json")
	case BestBuyGetDetailByUpcUrl:
		return readFile("offer/bestbuy/testdata/best_buy_get_by_upc_prod_detail_not_found.json")
	case EbayGetDetailUrl:
		return readFile("offer/ebay/testdata/ebay_find_by_upc_no_result.json")
	case AmazonGetDetailUrl:
		return readFile("offer/amazon/testdata/amazon_get_product_detail_by_upc_not_found.xml")

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

// Registers Mock endpoint responders for Search based API calls
func registerMockResponderSearch(httpMethod, apiUrl, apiType string, status int) {
	log.Printf("Mocking Search: %s %d - %s", httpMethod, status, apiUrl)

	switch apiType {
	case model.Trending:
		httpmock.RegisterResponder(httpMethod, apiUrl, httpmock.NewBytesResponder(status, getJsonBytesTrendingMock(apiUrl)))
		break
	case model.Search:
		httpmock.RegisterResponder(httpMethod, apiUrl, httpmock.NewBytesResponder(status, getJsonBytesSearchMock(apiUrl)))
		break
	case model.NoResults:
		httpmock.RegisterResponder(httpMethod, apiUrl, httpmock.NewBytesResponder(status, getJsonBytesSearchNoResultsMock(apiUrl)))
		break
	}
}

// Registers Mock endpoint responders for Get Detail based API calls
func registerMockResponderGetDetail(httpMethod, apiUrl, apiType string, status int) {
	log.Printf("Mocking GetDetail: %s %d - %s", httpMethod, status, apiUrl)

	switch apiType {
	case model.Id:
		httpmock.RegisterResponder(httpMethod, apiUrl, httpmock.NewBytesResponder(status, getJsonBytesGetDetailByIdMock(apiUrl)))
		break
	case model.Upc:
		httpmock.RegisterResponder(httpMethod, apiUrl, httpmock.NewBytesResponder(status, getJsonBytesGetDetailByUpcMock(apiUrl)))
		break
	case model.NoResults:
		httpmock.RegisterResponder(httpMethod, apiUrl, httpmock.NewBytesResponder(status, getJsonBytesGetDetailByUpcNoResultsMock(apiUrl)))
		break
	}
}

// Tests basic Search (no keywords) that returns trending results from external APIs
func TestSearch(t *testing.T) {

	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderSearch(http.MethodGet, WalmartTrendingUrl, model.Trending, 200)
	registerMockResponderSearch(http.MethodGet, BestBuyTrendingUrl, model.Trending, 200)
	registerMockResponderSearch(http.MethodGet, EbaySearchUrl, model.Trending, 200)
	registerMockResponderSearch(http.MethodGet, AmazonSearchUrl, model.Trending, 200)

	// call our local server API
	endpoint := "http://localhost:8080/"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()

	assert.True(t, strings.HasPrefix(body, `{"list":[{"`))

	walmartSnippet := `"externalId":"348726849","upc":"816586026705","name":"Best Choice Products 6' Exercise Tri-Fold Gym Mat For Gymnastics, Aerobics, Yoga, Martial Arts - Pink","partyName":"walmart.com"`
	assert.True(t, strings.Contains(body, walmartSnippet))

	ebaySnippet := `"externalId":"282629961650","upc":"","name":"Reverb Cross Men s Running Shoes","partyName":"ebay.com"`
	assert.True(t, strings.Contains(body, ebaySnippet))

	amazonSnippet := `"externalId":"B0743W4Y75","upc":"701649356113","name":"Bluetooth Smart Watch with Camera, Aosmart B23 Smart Watch for Android Smartphones (White)","partyName":"amazon.com"`
	assert.True(t, strings.Contains(body, amazonSnippet))

	bestBuySnippet := `"externalId":"5714687","upc":"","name":"Alienware - Aurora R6 Desktop - Intel Core i7 - 16GB Memory - NVIDIA GeForce GTX 1070 - 256GB Solid State Drive + 1TB Hard Drive - Silver","partyName":"bestbuy.com"`
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
	registerMockResponderSearch(http.MethodGet, WalmartSearchUrl, model.Search, 200)
	registerMockResponderSearch(http.MethodGet, BestBuySearchUrl, model.Search, 200)
	registerMockResponderSearch(http.MethodGet, EbaySearchUrl, model.Search, 200)
	registerMockResponderSearch(http.MethodGet, AmazonSearchUrl, model.Search, 200)

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

	walmartSnippet := `"externalId":"53966162","upc":"093155171244","name":"Skyrim Special Edition (Xbox One)","partyName":"walmart.com",`
	assert.True(t, strings.Contains(body, walmartSnippet))

	bestBuySnippet := `"externalId":"5626200","upc":"600603210488","name":"The Elder Scrolls V: Skyrim Special Edition Best Buy Exclusive Dragonborn Bundle - Xbox One","partyName":"bestbuy.com"`
	assert.True(t, strings.Contains(body, bestBuySnippet))

	ebaySnippet := `"externalId":"223482818","upc":"","name":"Elder Scrolls V: Skyrim - Special Edition With Bonus Steelbook Case PS4 ","partyName":"ebay.com"`
	assert.True(t, strings.Contains(body, ebaySnippet))

	amazonSnippet := `"externalId":"B01GW8XJVU","upc":"093155171251","name":"The Elder Scrolls V: Skyrim - Special Edition - PlayStation 4","partyName":"amazon.com"`
	assert.True(t, strings.Contains(body, amazonSnippet))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, WalmartSearchUrl, 1)
	assertCallsMade(t, http.MethodGet, BestBuySearchUrl, 1)
	assertCallsMade(t, http.MethodGet, EbaySearchUrl, 1)
	assertCallsMade(t, http.MethodGet, AmazonSearchUrl, 1)
}

// Tests Search with keywords that returns results from external APIs - SORT by name DESC
func TestSearchWithKeywordsSortByNameDesc(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderSearch(http.MethodGet, WalmartSearchUrl, model.Search, 200)
	registerMockResponderSearch(http.MethodGet, BestBuySearchUrl, model.Search, 200)
	registerMockResponderSearch(http.MethodGet, EbaySearchUrl, model.Search, 200)
	registerMockResponderSearch(http.MethodGet, AmazonSearchUrl, model.Search, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers"
	var jsonRequest = []byte(`{"searchColumns":[{"name":"name","value":"skyrim"}],"sortBy":"name", "sortOrder":"desc","page":1,"rowsPerPage":10}`)

	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()

	assert.True(t, strings.HasPrefix(body, `{"list":[{"`))
	assert.True(t, strings.Contains(body, `"externalId":"0","upc":"","name":"Skyrim VR Digital - PlayStation 4`))

	walmartSnippet := `"externalId":"53966162","upc":"093155171244","name":"Skyrim Special Edition (Xbox One)","partyName":"walmart.com",`
	assert.True(t, strings.Contains(body, walmartSnippet))

	bestBuySnippet := `"externalId":"5626200","upc":"600603210488","name":"The Elder Scrolls V: Skyrim Special Edition Best Buy Exclusive Dragonborn Bundle - Xbox One","partyName":"bestbuy.com"`
	assert.True(t, strings.Contains(body, bestBuySnippet))

	ebaySnippet := `"externalId":"223482818","upc":"","name":"Elder Scrolls V: Skyrim - Special Edition With Bonus Steelbook Case PS4 ","partyName":"ebay.com"`
	assert.True(t, strings.Contains(body, ebaySnippet))

	amazonSnippet := `"externalId":"B01GW8XJVU","upc":"093155171251","name":"The Elder Scrolls V: Skyrim - Special Edition - PlayStation 4","partyName":"amazon.com"`
	assert.True(t, strings.Contains(body, amazonSnippet))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, WalmartSearchUrl, 1)
	assertCallsMade(t, http.MethodGet, BestBuySearchUrl, 1)
	assertCallsMade(t, http.MethodGet, EbaySearchUrl, 1)
	assertCallsMade(t, http.MethodGet, AmazonSearchUrl, 1)
}

// Tests Search SORT by price ASC
func TestSearchWithKeywordsSortByPriceAsc(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderSearch(http.MethodGet, WalmartSearchUrl, model.Search, 200)
	registerMockResponderSearch(http.MethodGet, BestBuySearchUrl, model.Search, 200)
	registerMockResponderSearch(http.MethodGet, EbaySearchUrl, model.Search, 200)
	registerMockResponderSearch(http.MethodGet, AmazonSearchUrl, model.Search, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers"
	var jsonRequest = []byte(`{"searchColumns":[{"name":"name","value":"skyrim"}],"sortBy":"price", "sortOrder":"asc","page":1,"rowsPerPage":10}`)

	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()

	assert.True(t, strings.HasPrefix(body, `{"list":[{"`))
	assert.True(t, strings.Contains(body, `"externalId":"0","upc":"","name":"Skyrim VR Digital - PlayStation 4`))

	walmartSnippet := `"externalId":"53966162","upc":"093155171244","name":"Skyrim Special Edition (Xbox One)","partyName":"walmart.com",`
	assert.True(t, strings.Contains(body, walmartSnippet))

	bestBuySnippet := `"externalId":"5626200","upc":"600603210488","name":"The Elder Scrolls V: Skyrim Special Edition Best Buy Exclusive Dragonborn Bundle - Xbox One","partyName":"bestbuy.com"`
	assert.True(t, strings.Contains(body, bestBuySnippet))

	ebaySnippet := `"externalId":"223482818","upc":"","name":"Elder Scrolls V: Skyrim - Special Edition With Bonus Steelbook Case PS4 ","partyName":"ebay.com"`
	assert.True(t, strings.Contains(body, ebaySnippet))

	amazonSnippet := `"externalId":"B01GW8XJVU","upc":"093155171251","name":"The Elder Scrolls V: Skyrim - Special Edition - PlayStation 4","partyName":"amazon.com"`
	assert.True(t, strings.Contains(body, amazonSnippet))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, WalmartSearchUrl, 1)
	assertCallsMade(t, http.MethodGet, BestBuySearchUrl, 1)
	assertCallsMade(t, http.MethodGet, EbaySearchUrl, 1)
	assertCallsMade(t, http.MethodGet, AmazonSearchUrl, 1)
}

// tests Reset and List from Datastore - Cassandra must be running
func TestSearchDatastore(t *testing.T) {
	// 1st call Reset to reset database and re-create keystore
	endpoint1 := "http://localhost:8080/reset"
	req1, _ := http.NewRequest(http.MethodGet, endpoint1, nil)
	response1 := executeRequest(req1)
	assert.Equal(t, 200, response1.Code)

	// 2nd call our local server API to fetch list
	endpoint := "http://localhost:8080/offerlist"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)
	body := response.Body.String()

	assert.True(t, strings.HasPrefix(body, `{"list":[{"`))
	assert.True(t, strings.Contains(body, `"externalId":"1","upc":"upc12345678","name":"offer 1"`))
	assert.True(t, strings.Contains(body, `"externalId":"2","upc":"upc22345678","name":"offer 2"`))
	assert.True(t, strings.Contains(body, `"externalId":"3","upc":"upc32345678","name":"offer 3"`))
	assert.True(t, strings.Contains(body, `"externalId":"4","upc":"upc42345678","name":"offer 4"`))
}

// tests Reset and Add to Datastore - Cassandra must be running
func TestResetAddOfferDatastore(t *testing.T) {
	// 1st call Reset to reset database and re-create keystore
	endpoint1 := "http://localhost:8080/reset"
	req1, _ := http.NewRequest(http.MethodGet, endpoint1, nil)
	response1 := executeRequest(req1)
	assert.Equal(t, 200, response1.Code)

	// call our local server API to add
	endpoint := "http://localhost:8080/offerlist"
	var jsonRequest = []byte(`{"externalId":"1", "upc":"upc999","name":"test record","partyName":"amazon.com","semanticName":"http:/item01","mainImageFileUrl":"http:/item01.jpg","partyImageFileUrl":"amazon-logo.jpg","productCategory":"laptops","price":500,"rating":3.88,"numReviews":120}`)

	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	assert.Equal(t, 201, response.Code)

	// verify responses
	body := response.Body.String()
	assert.True(t, strings.Contains(body, `"externalId":"1","upc":"upc999","name":"test record","partyName":"amazon.com","semanticName":"http:/item01"`))
}

// Tests Search with keywords invalid expects Bad Request 400
func testSearchWithKeywordsInvalidRequest(t *testing.T, json []byte) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderSearch(http.MethodGet, WalmartSearchUrl, model.Search, 200)
	registerMockResponderSearch(http.MethodGet, BestBuySearchUrl, model.Search, 200)
	registerMockResponderSearch(http.MethodGet, EbaySearchUrl, model.Search, 200)
	registerMockResponderSearch(http.MethodGet, AmazonSearchUrl, model.Search, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers"

	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(json))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	assert.Equal(t, 400, response.Code)

	// verify responses
	body := response.Body.String()
	assert.True(t, strings.Contains(body, model.InvalidRequest))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, WalmartSearchUrl, 0)
	assertCallsMade(t, http.MethodGet, BestBuySearchUrl, 0)
	assertCallsMade(t, http.MethodGet, EbaySearchUrl, 0)
	assertCallsMade(t, http.MethodGet, AmazonSearchUrl, 0)
}

// Tests Search with keywords invalid expects Bad Request 400 - sort order
func TestSearchWithKeywordsInvalidSortOrder(t *testing.T) {
	var jsonRequest = []byte(`{"searchColumns":[{"name":"name","value":"skyrim"}],"sortBy":"name", "sortOrder":"abc","page":1,"rowsPerPage":10}`)
	testSearchWithKeywordsInvalidRequest(t, jsonRequest)
}

// Tests Search with keywords invalid expects Bad Request 400 - page
func TestSearchWithKeywordsInvalidPage(t *testing.T) {
	var jsonRequest = []byte(`{"searchColumns":[{"name":"name","value":"skyrim"}],"sortBy":"name","sortOrder":"asc","page":-1,"rowsPerPage":10}`)
	testSearchWithKeywordsInvalidRequest(t, jsonRequest)
}

// Tests Search with keywords invalid expects Bad Request 400 - rowsPerPage
func TestSearchWithKeywordsInvalidRowsPerPage(t *testing.T) {
	var jsonRequest = []byte(`{"searchColumns":[{"name":"name","value":"skyrim"}],"sortBy":"name","sortOrder":"asc","page":1,"rowsPerPage":-10}`)
	testSearchWithKeywordsInvalidRequest(t, jsonRequest)
}

// Tests Search No results
func TestSearchNoResults(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderSearch(http.MethodGet, WalmartSearchUrl, model.NoResults, 200)
	registerMockResponderSearch(http.MethodGet, BestBuySearchUrl, model.NoResults, 200)
	registerMockResponderSearch(http.MethodGet, EbaySearchUrl, model.NoResults, 200)
	registerMockResponderSearch(http.MethodGet, AmazonSearchUrl, model.NoResults, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers"
	var jsonRequest = []byte(`{"searchColumns":[{"name":"name","value":"skyrim"}],"sortOrder":"asc","page":1,"rowsPerPage":10}`)

	req, _ := http.NewRequest(http.MethodPost, endpoint, bytes.NewBuffer(jsonRequest))
	req.Header.Set("Content-Type", "application/json")
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()
	assert.True(t, strings.HasPrefix(body, `{"list":[],"summary":{"page":1,"pageCount":1,"totalCount":0}`))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, WalmartSearchUrl, 1)
	assertCallsMade(t, http.MethodGet, BestBuySearchUrl, 1)
	assertCallsMade(t, http.MethodGet, EbaySearchUrl, 1)
	assertCallsMade(t, http.MethodGet, AmazonSearchUrl, 1)
}

// Tests GetDetail By Id - walmart
func TestGetDetailByIdWalmart(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderGetDetail(http.MethodGet, WalmartGetDetailUrl, model.Id, 200)
	registerMockResponderGetDetail(http.MethodGet, BestBuyGetDetailByUpcUrl, model.Upc, 200)
	registerMockResponderGetDetail(http.MethodGet, EbayGetDetailUrl, model.Upc, 200)
	registerMockResponderGetDetail(http.MethodGet, AmazonGetDetailUrl, model.Upc, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers/53966162?idType=id&source=walmart.com"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()

	//assert.True(t, strings.HasPrefix(body, `{"offer":{"id":"55760264","upc":"065857174434",`))

	walmartSnippet := `"name":"Better Homes and Gardens Leighton Twin-Over-Full Bunk Bed, Multiple Colors","partyName":"walmart.com",`
	assert.True(t, strings.Contains(body, walmartSnippet))

	bestBuySnippet := `{"partyName":"bestbuy.com","semanticName":"https://api.bestbuy.com/click/-/5529006/pdp"`
	assert.True(t, strings.Contains(body, bestBuySnippet))

	ebaySnippet := `{"partyName":"ebay.com","semanticName":"http://www.ebay.com/itm/New-Laptop-Toshiba-Satellite-L355-S7907-17-Intel-Pentium-Dual-core-T3400-2-16Gh`
	assert.True(t, strings.Contains(body, ebaySnippet))

	amazonSnippet := `{"partyName":"amazon.com","semanticName":"https://www.amazon.com/Elder-Scrolls-Skyrim-strategy-bundle-Playstation`
	assert.True(t, strings.Contains(body, amazonSnippet))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, WalmartGetDetailUrl, 1)
	assertCallsMade(t, http.MethodGet, BestBuyGetDetailByUpcUrl, 1)
	assertCallsMade(t, http.MethodGet, EbayGetDetailUrl, 1)
	assertCallsMade(t, http.MethodGet, AmazonGetDetailUrl, 1)
}

// Tests GetDetail By Upc Not Found Walmart
func TestGetDetailByUpcNotFoundWalmart(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderGetDetail(http.MethodGet, WalmartGetDetailByUpcUrl, model.NoResults, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers/12345678?idType=upc&source=walmart.com"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 404, response.Code)

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, WalmartGetDetailByUpcUrl, 1)
}

// Tests GetDetail By Id - No Competitors search by UPC detail items found - walmart
func TestGetDetailByIdWalmartNoDetailItems(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderGetDetail(http.MethodGet, WalmartGetDetailUrl, model.Id, 200)
	registerMockResponderGetDetail(http.MethodGet, BestBuyGetDetailByUpcUrl, model.NoResults, 200)
	registerMockResponderGetDetail(http.MethodGet, EbayGetDetailUrl, model.NoResults, 200)
	registerMockResponderGetDetail(http.MethodGet, AmazonGetDetailUrl, model.NoResults, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers/53966162?idType=id&source=walmart.com"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()

	assert.True(t, strings.HasPrefix(body, `{"offer":{"id":"`))
	assert.True(t, strings.Contains(body, `"externalId":"55760264","upc":"065857174434","name":"Better Homes`))
	assert.True(t, strings.Contains(body, `"productDetailItems":[]`))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, WalmartGetDetailUrl, 1)
	assertCallsMade(t, http.MethodGet, BestBuyGetDetailByUpcUrl, 1)
	assertCallsMade(t, http.MethodGet, EbayGetDetailUrl, 1)
	assertCallsMade(t, http.MethodGet, AmazonGetDetailUrl, 1)
}

// Tests GetDetail By Id BestBuy
func TestGetDetailByIdBestBuy(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderGetDetail(http.MethodGet, BestBuyGetDetailUrl, model.Id, 200)
	registerMockResponderGetDetail(http.MethodGet, WalmartGetDetailByUpcUrl, model.Upc, 200)
	registerMockResponderGetDetail(http.MethodGet, EbayGetDetailUrl, model.Upc, 200)
	registerMockResponderGetDetail(http.MethodGet, AmazonGetDetailUrl, model.Upc, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers/5529006?idType=id&source=bestbuy.com"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()

	assert.True(t, strings.HasPrefix(body, `{"offer":{"id":"`))
	assert.True(t, strings.Contains(body, `"externalId":"5529006","upc":"849803052423","name":"Funko - Elder Scrolls V`))

	walmartSnippet := `{"partyName":"walmart.com","semanticName":"http://linksynergy.walmart.com/fs-bin/click?id=12345678`
	assert.True(t, strings.Contains(body, walmartSnippet))

	bestBuySnippet := `"partyName":"bestbuy.com","semanticName":"https://api.bestbuy.com/click`
	assert.True(t, strings.Contains(body, bestBuySnippet))

	ebaySnippet := `{"partyName":"ebay.com","semanticName":"http://www.ebay.com/itm`
	assert.True(t, strings.Contains(body, ebaySnippet))

	amazonSnippet := `{"partyName":"amazon.com","semanticName":"https://www.amazon.com/Elder-Scrolls-Skyrim-strategy-bundle-Playstation`
	assert.True(t, strings.Contains(body, amazonSnippet))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, BestBuyGetDetailUrl, 1)
	assertCallsMade(t, http.MethodGet, WalmartGetDetailByUpcUrl, 1)
	assertCallsMade(t, http.MethodGet, EbayGetDetailUrl, 1)
	assertCallsMade(t, http.MethodGet, AmazonGetDetailUrl, 1)
}

// Tests GetDetail By Upc Not Found - Best Buy
func TestGetDetailByUpcNotFoundBestBuy(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderGetDetail(http.MethodGet, BestBuyGetDetailByUpcUrl, model.NoResults, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers/065857174434?idType=upc&source=bestbuy.com"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 404, response.Code)

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, BestBuyGetDetailByUpcUrl, 1)
}

// Tests GetDetail By Id - No Competitors search by UPC detail items found - BestBuy
func TestGetDetailByIdBestBuyNoDetailItems(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderGetDetail(http.MethodGet, BestBuyGetDetailUrl, model.Id, 200)
	registerMockResponderGetDetail(http.MethodGet, WalmartGetDetailByUpcUrl, model.NoResults, 200)
	registerMockResponderGetDetail(http.MethodGet, EbayGetDetailUrl, model.NoResults, 200)
	registerMockResponderGetDetail(http.MethodGet, AmazonGetDetailUrl, model.NoResults, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers/5529006?idType=id&source=bestbuy.com"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()

	//assert.True(t, strings.HasPrefix(body, `{"offer":{"id":"5529006","upc":"849803052423","name":"Funko - Elder Scrolls V: Skyrim`))
	assert.True(t, strings.Contains(body, `"productDetailItems":[]`))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, BestBuyGetDetailUrl, 1)
	assertCallsMade(t, http.MethodGet, WalmartGetDetailByUpcUrl, 1)
	assertCallsMade(t, http.MethodGet, EbayGetDetailUrl, 1)
	assertCallsMade(t, http.MethodGet, AmazonGetDetailUrl, 1)
}

// Tests GetDetail By Id Ebay returns no UPC so product detail items is empty (not fetching others competitors prices)
func TestGetDetailByIdEbay(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderGetDetail(http.MethodGet, EbayGetDetailUrl, model.Id, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers/62923188?idType=id&source=ebay.com"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()

	assert.True(t, strings.HasPrefix(body, `{"offer":{"id":"`))
	assert.True(t, strings.Contains(body, `"externalId":"62923188","upc":"","name":"Harry Potter and the Order of the Phoenix-(DVD, Widescreen`))

	ebaySnippet := `partyName":"ebay.com","semanticName":"http://www.ebay.com/itm/Harry-Potter-and-Order-Phoenix-DVD-Widescreen-Edition-BRAND-NEW`
	assert.True(t, strings.Contains(body, ebaySnippet))

	assert.True(t, strings.Contains(body, `"price":5.62,"rating":0,"numReviews":0},"description":"","attributes":[],"productDetailItems":[]}`))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, EbayGetDetailUrl, 1)
}

// Tests GetDetail By Upc Not Found - Ebay
func TestGetDetailByUpcNotFoundEbay(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderGetDetail(http.MethodGet, EbayGetDetailUrl, model.NoResults, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers/123456789?idType=upc&source=ebay.com"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 404, response.Code)

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, EbayGetDetailUrl, 1)
}

// Tests GetDetail By Id Amazon
func TestGetDetailByIdAmazon(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderGetDetail(http.MethodGet, AmazonGetDetailUrl, model.Id, 200)
	registerMockResponderGetDetail(http.MethodGet, BestBuyGetDetailByUpcUrl, model.Upc, 200)
	registerMockResponderGetDetail(http.MethodGet, WalmartGetDetailByUpcUrl, model.Upc, 200)
	registerMockResponderGetDetail(http.MethodGet, EbayGetDetailUrl, model.Upc, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers/5529006?idType=id&source=amazon.com"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()

	assert.True(t, strings.HasPrefix(body, `{"offer":{"id":"`))
	assert.True(t, strings.Contains(body, `"externalId":"B01GW8XJVU","upc":"065857174434","name":"The Elder Scrolls V: Skyrim - Special Edition - PlayStation 4"`))

	desc := `"description":"All-new features include remastered art and effects, volumetric god rays, dynamic depth of field, screen-space reflections, and more."`
	assert.True(t, strings.Contains(body, desc))
	assert.True(t, strings.Contains(body, `"name":"manufacturer","value":"Bethesda"`))

	walmartSnippet := `{"partyName":"walmart.com","semanticName":"http://linksynergy.walmart.com/fs-bin/click?id=12345678`
	assert.True(t, strings.Contains(body, walmartSnippet))

	bestBuySnippet := `{"partyName":"bestbuy.com","semanticName":"https://api.bestbuy.com/click/-/5529006/pdp"`
	assert.True(t, strings.Contains(body, bestBuySnippet))

	ebaySnippet := `{"partyName":"ebay.com","semanticName":"http://www.ebay.com/itm`
	assert.True(t, strings.Contains(body, ebaySnippet))

	amazonSnippet := `"partyName":"amazon.com","semanticName":"https://www.amazon.com/Elder-Scrolls-Skyrim-Special-PlayStation-4`
	assert.True(t, strings.Contains(body, amazonSnippet))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, AmazonGetDetailUrl, 1)
	assertCallsMade(t, http.MethodGet, BestBuyGetDetailByUpcUrl, 1)
	assertCallsMade(t, http.MethodGet, WalmartGetDetailByUpcUrl, 1)
	assertCallsMade(t, http.MethodGet, EbayGetDetailUrl, 1)
}

// Tests GetDetail By Upc Not Found - Amazon
func TestGetDetailByUpcNotFoundAmazon(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderGetDetail(http.MethodGet, AmazonGetDetailUrl, model.NoResults, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers/123456789?idType=upc&source=amazon.com"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 404, response.Code)

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, AmazonGetDetailUrl, 1)
}

// Tests GetDetail By Id - No Competitors search by UPC detail items found - Amazon
func TestGetDetailByIdAmazonNoDetailItems(t *testing.T) {
	// register mock for external API endpoints
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	// External Vendor Apis
	registerMockResponderGetDetail(http.MethodGet, AmazonGetDetailUrl, model.Id, 200)
	registerMockResponderGetDetail(http.MethodGet, BestBuyGetDetailByUpcUrl, model.NoResults, 200)
	registerMockResponderGetDetail(http.MethodGet, WalmartGetDetailByUpcUrl, model.NoResults, 200)
	registerMockResponderGetDetail(http.MethodGet, EbayGetDetailUrl, model.NoResults, 200)

	// call our local server API
	endpoint := "http://localhost:8080/offers/123456789?idType=id&source=amazon.com"
	req, _ := http.NewRequest(http.MethodGet, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	body := response.Body.String()

	assert.True(t, strings.HasPrefix(body, `{"offer":{"id":"`))
	assert.True(t, strings.Contains(body, `"externalId":"B01GW8XJVU","upc":"065857174434","name":"The Elder Scrolls V: Skyrim - Special Edition`))
	assert.True(t, strings.Contains(body, `"productDetailItems":[]`))

	// get the amount of calls for the registered responders
	assertCallsMade(t, http.MethodGet, AmazonGetDetailUrl, 1)
	assertCallsMade(t, http.MethodGet, BestBuyGetDetailByUpcUrl, 1)
	assertCallsMade(t, http.MethodGet, WalmartGetDetailByUpcUrl, 1)
	assertCallsMade(t, http.MethodGet, EbayGetDetailUrl, 1)
}
