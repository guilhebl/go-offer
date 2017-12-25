package walmart

import (
	"encoding/json"
	"fmt"
	"github.com/guilhebl/go-offer/common/config"
	"github.com/guilhebl/go-offer/common/model"
	"github.com/guilhebl/go-offer/offer/monitor"
	"github.com/guilhebl/go-strutil"
	"github.com/guilhebl/go-worker-pool"
	"log"
	"net/http"
	"strconv"
	"time"
)

// Creates Job for Searching offers from Walmart and returns a Channel with jobResults
func SearchOffers(m map[string]string) *job.Job {
	// create output channel
	out := job.NewJobResultChannel()

	// let's create a job with the payload
	task := NewSearchTask()
	job := job.NewJob(&task, m, out)
	return &job
}

// Searches for offers from Walmart
func search(m map[string]string) *model.OfferList {
	// try to acquire lock from request Monitor
	if !monitor.IsServiceAvailable(model.Walmart) {
		log.Printf("Unable to acquire lock from Request Monitor")
		return nil
	}

	// format walmart specific params
	p := filterParams(m)

	endpoint := config.GetProperty("walmartEndpoint")
	isKeywordSearch := p[model.Query] != ""
	page, _ := strconv.ParseInt(p[model.Page], 10, 0)
	pageSize := int(config.GetIntProperty("walmartDefaultPageSize"))
	responseGroup := config.GetProperty("walmartSearchResponseGroup")
	apiKey := config.GetProperty("walmartApiKey")
	affiliateId := config.GetProperty("walmartAffiliateId")
	timeout := time.Duration(config.GetIntProperty("marketplaceDefaultTimeout")) * time.Millisecond

	var start = 1
	if page > 1 {
		start = int(page) - 1*int(pageSize) + 1
	}

	if isKeywordSearch {
		path := config.GetProperty("walmartProductSearchPath")
		url := fmt.Sprintf("%s/%s", endpoint, path)

		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/json")
		q := req.URL.Query()
		q.Add("format", "json")
		q.Add("responseGroup", responseGroup)
		q.Add("apiKey", apiKey)
		q.Add("lsPublisherId", affiliateId)
		q.Add("query", p["query"])
		q.Add("start", strconv.Itoa(start))
		req.URL.RawQuery = q.Encode()
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

		var entity SearchResponse

		if err := json.NewDecoder(resp.Body).Decode(&entity); err != nil {
			log.Println(err)
			return nil
		}
		return buildSearchResponse(&entity, pageSize)
	} else {
		// search trending items if no keyword provided
		path := config.GetProperty("walmartProductTrendingPath")
		url := fmt.Sprintf("%s/%s", endpoint, path)
		req, err := http.NewRequest("GET", url, nil)

		req.Header.Set("Accept", "application/json")
		q := req.URL.Query()
		q.Add("format", "json")
		q.Add("responseGroup", responseGroup)
		q.Add("apiKey", apiKey)
		q.Add("lsPublisherId", affiliateId)
		req.URL.RawQuery = q.Encode()
		url = fmt.Sprintf(req.URL.String())
		log.Printf("Walmart trending: %s", url)

		client := &http.Client{
			Timeout: timeout,
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Do: ", err)
			return nil
		}
		defer resp.Body.Close()

		var entity TrendingResponse

		if err := json.NewDecoder(resp.Body).Decode(&entity); err != nil {
			log.Println(err)
			return nil
		}
		return buildTrendingResponse(&entity, int(page), pageSize)
	}

	return nil
}

func buildTrendingResponse(r *TrendingResponse, page, pageSize int) *model.OfferList {
	list := buildSearchItemList(r.Items)
	l := len(list)
	o := model.NewOfferList(list, page, l/pageSize, l)
	return o
}

func buildSearchResponse(r *SearchResponse, pageSize int) *model.OfferList {
	list := buildSearchItemList(r.Items)
	o := model.NewOfferList(list, r.Start/pageSize+1, r.TotalResults/pageSize, r.TotalResults)
	return o
}

func buildSearchItemList(items []SearchItem) []model.Offer {
	list := make([]model.Offer, 0)
	proxyRequired := config.IsProxyRequired(model.Walmart)

	for _, item := range items {
		rate := 0.0

		// format Rating string
		if item.CustomerRating != "" {
			formattedRate, err := strconv.ParseFloat(item.CustomerRating, 32)
			if err != nil {
				log.Println(err)
				rate = formattedRate
			}
		}

		o := model.NewOffer(
			strconv.Itoa(item.ItemId),
			item.Upc,
			item.Name,
			model.Walmart,
			item.ProductTrackingUrl,
			config.BuildImgUrlExternal(item.LargeImage, proxyRequired),
			config.BuildImgUrl("walmart-logo.png"),
			item.CategoryPath,
			item.SalePrice,
			float32(rate),
			item.NumReviews,
		)

		list = append(list, *o)
	}

	return list
}

func filterParams(m map[string]string) map[string]string {
	p := make(map[string]string)

	// get search keyword phrase
	if m[model.Name] != "" {
		p[model.Query] = m[model.Name]
	}

	// get page - defaults to 1
	if m[model.Page] != "" {
		p[model.Page] = m[model.Page]
	} else {
		p[model.Page] = "1"
	}

	return p
}

// Creates Job for fetching Product Detail and returns a Channel with jobResult
func GetDetailJob(id, idType, country string) *job.Job {
	// convert to map for job to consume
	m := make(map[string]string)
	m["id"], m["idType"], m["country"] = id, idType, country

	// create output channel
	out := job.NewJobResultChannel()

	// let's create a job with the payload
	task := NewGetDetailTask()
	job := job.NewJob(&task, m, out)
	return &job
}

// Search for a specific product detail either by Id or Upc
func GetOfferDetail(id string, idType string, country string) *model.OfferDetail {
	log.Printf("Get Detail: %s, %s, %s", id, idType, country)

	// try to acquire lock from request Monitor
	if !monitor.IsServiceAvailable(model.Walmart) {
		log.Printf("Unable to acquire lock from Request Monitor")
		return nil
	}

	endpoint := config.GetProperty("walmartEndpoint")
	path := config.GetProperty("walmartProductDetailPath")
	apiKey := config.GetProperty("walmartApiKey")
	affiliateId := config.GetProperty("walmartAffiliateId")
	timeout := time.Duration(config.GetIntProperty("marketplaceDefaultTimeout")) * time.Millisecond

	if idType == model.Id {
		url := fmt.Sprintf("%s/%s/%s", endpoint, path, id)
		req, err := http.NewRequest("GET", url, nil)

		req.Header.Set("Accept", "application/json")
		q := req.URL.Query()
		q.Add("format", "json")
		q.Add("apiKey", apiKey)
		q.Add("lsPublisherId", affiliateId)
		req.URL.RawQuery = q.Encode()
		url = fmt.Sprintf(req.URL.String())
		log.Printf("Walmart get: %s", url)

		client := &http.Client{
			Timeout: timeout,
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Do: ", err)
			return nil
		}
		defer resp.Body.Close()

		var entity SearchItem

		if err := json.NewDecoder(resp.Body).Decode(&entity); err != nil {
			log.Println(err)
			return nil
		}
		return buildProductDetail(&entity)
	} else if idType == model.Upc {
		url := endpoint + "/" + path
		req, err := http.NewRequest("GET", url, nil)

		req.Header.Set("Accept", "application/json")
		q := req.URL.Query()
		q.Add(model.Upc, id)
		q.Add("format", "json")
		q.Add("apiKey", apiKey)
		q.Add("lsPublisherId", affiliateId)
		req.URL.RawQuery = q.Encode()
		url = fmt.Sprintf(req.URL.String())
		log.Printf("Walmart get by UPC: %s", url)

		client := &http.Client{
			Timeout: timeout,
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Do: ", err)
			return nil
		}
		defer resp.Body.Close()

		var entity BaseSearchResponse

		if err := json.NewDecoder(resp.Body).Decode(&entity); err != nil {
			log.Println(err)
			return nil
		}
		return buildProductDetailSearchResponse(&entity)
	}

	return nil
}

func buildProductDetail(item *SearchItem) *model.OfferDetail {
	proxyRequired := config.IsProxyRequired(model.Walmart)

	rate, err := strconv.ParseFloat(item.CustomerRating, 32)
	if err != nil {
		log.Printf("error on parsing rate for item string: %s", item.CustomerRating)
		rate = 0.0
	}

	o := model.NewOffer(
		strconv.Itoa(item.ItemId),
		item.Upc,
		item.Name,
		model.Walmart,
		item.ProductTrackingUrl,
		config.BuildImgUrlExternal(item.LargeImage, proxyRequired),
		config.BuildImgUrl("walmart-logo.png"),
		item.CategoryPath,
		item.SalePrice,
		float32(rate),
		item.NumReviews,
	)

	attrs := make(map[string]string)
	detItems := make([]model.OfferDetailItem, 0)

	det := model.NewOfferDetail(
		*o,
		strutil.FilterHtmlTags(item.LongDescription),
		attrs,
		detItems,
	)

	return det
}

func buildProductDetailSearchResponse(item *BaseSearchResponse) *model.OfferDetail {
	if item == nil || len(item.Items) == 0 {
		return nil
	}
	p := item.Items[0]

	return buildProductDetail(&p)
}
