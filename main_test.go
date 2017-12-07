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

func buildWalmartStub() *httptest.Server {
	var resp []byte
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("calling Server Handler FUNC STUB ")

		switch r.Host {
		case "api.walmartlabs.com":
			resp = readFile("offer/walmart/walmart_sample_search_response")
		default:
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		w.Write(resp)
	}))
}

func getJsonMockBytes(url string) []byte {
	switch url {
	case "http://api.walmartlabs.com/v1/trends":
		return readFile("offer/walmart/walmart_sample_trending_response.json")
	case "https://api.walmartlabs.com/search":
		return readFile("offer/walmart/walmart_sample_search_response.json")
	default:
		return nil
	}
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	offer.GetInstance().Router.ServeHTTP(rr, req)
	return rr
}

func TestSearch(t *testing.T) {

	// register mock for external API endpoint
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpMethod := "GET"
	apiUrl := "http://api.walmartlabs.com/v1/trends"
	log.Printf("Mocking Search: %s", apiUrl)
	httpmock.RegisterResponder(httpMethod, apiUrl, httpmock.NewBytesResponder(200, getJsonMockBytes(apiUrl)))

	// call our local server API
	endpoint := "http://localhost:8080/"
	req, _ := http.NewRequest(httpMethod, endpoint, nil)
	response := executeRequest(req)
	assert.Equal(t, 200, response.Code)

	// verify responses
	jsonSnippet := `{"list":[{"id":"55760264","upc":"065857174434","name":"Better Homes and Gardens Leighton Twin-Over-Full Bunk Bed, Multiple Colors","partyName":"walmart.com"`
	body := response.Body.String()
	assert.True(t, strings.HasPrefix(body, jsonSnippet))

	// get the amount of calls for the registered responder
	info := httpmock.GetCallCountInfo()
	count := info[httpMethod+" "+apiUrl]
	assert.Equal(t, 1, count)
	log.Printf("Total External API Calls made: %d", count)
}
