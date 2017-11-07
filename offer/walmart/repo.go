package walmart

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"github.com/guilhebl/offergo/common/model"
	"github.com/guilhebl/offergo/common/config"
	"github.com/guilhebl/offergo/offer/monitor"
	"strconv"
	"time"
)

// Searches for offers from Walmart
func SearchOffers(m map[string]string) *model.OfferList {

	// try to acquire lock from request Monitor
	if !monitor.IsServiceAvailable(model.Walmart) {
		log.Printf("Unable to acquire lock from Request Monitor")
		return nil
	}

	// format walmart specific params
	p := filterParams(m)

	endpoint := config.GetProperty("walmartEndpoint")
	isKeywordSearch := p["query"] != ""
	page, _ := strconv.ParseInt(p[model.Page], 10, 0)
	pageSize := config.GetIntProperty("walmartDefaultPageSize")
	responseGroup := config.GetProperty("walmartSearchResponseGroup")
	apiKey := config.GetProperty("walmartApiKey")
	affiliateId := config.GetProperty("walmartAffiliateId")
	timeout := time.Duration(config.GetIntProperty("marketplaceDefaultTimeout")) * time.Millisecond

	var start = 1
	if page > 1 {
		start = int(page) - 1 * int(pageSize) + 1
	}

	if isKeywordSearch {

		url := endpoint + "/" + config.GetProperty("walmartProductSearchPath")
		req, err := http.NewRequest("GET", url, nil)

		req.Header.Set("Accept", "application/json")
		q := req.URL.Query()
		q.Add("format", "json")
		q.Add("responseGroup", responseGroup)
		q.Add("apiKey", apiKey)
		q.Add("lsPublisherId", affiliateId)
		q.Add("query", p["query"])
		q.Add("start", strconv.Itoa(start))

		url = fmt.Sprintf(req.URL.String())
		log.Printf("Walmart search: %s", url)

		client := &http.Client{
			Timeout: timeout,
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Do: ", err)
			return nil
		}
		defer resp.Body.Close()
	}
	return nil
}

func filterParams(m map[string]string) map[string]string {
	p := make(map[string]string)

	// get search keyword phrase
	if m[model.Name] != "" {
		p["query"] = m[model.Name]
	}

	// get page - defaults to 1
	if m[model.Page] != "" {
		p[model.Page] = m[model.Page]
	} else {
		p[model.Page] = "1"
	}

	return p
}

func GetOfferDetail(id string, idType string, source string) model.OfferDetail {
	endpoint := config.GetEndpoint() + "/" + id + "?idType=" + idType + "&source=" + source
	url := fmt.Sprintf(endpoint)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	var entity model.OfferDetail

	log.Println("getDetail: %s", req)

	if err != nil {
		log.Fatal("NewRequest: ", err)
		return entity
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return entity
	}
	defer resp.Body.Close()

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&entity); err != nil {
		log.Println(err)
	}

	return entity
}
