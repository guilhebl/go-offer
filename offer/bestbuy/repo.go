package bestbuy

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
	"strings"
	"time"
)

// Creates Job for Searching offers and returns a Channel with jobResults
func SearchOffers(m map[string]string) *job.Job {
	// create output channel
	out := job.NewJobResultChannel()

	// let's create a job with the payload
	task := NewSearchTask()
	job := job.NewJob(&task, m, out)
	return &job
}

// Searches for offers from BBY
func search(m map[string]string) *model.OfferList {

	// try to acquire lock from request Monitor
	if !monitor.IsServiceAvailable(model.BestBuy) {
		log.Printf("Unable to acquire lock from Request Monitor")
		return nil
	}

	// format vendor specific params
	p := filterParams(m)

	endpoint := config.GetProperty("bestbuyEndpoint")
	isKeywordSearch := p[model.Keywords] != ""
	page, _ := strconv.ParseInt(p[model.Page], 10, 0)
	pageSize := int(config.GetIntProperty("bestbuyDefaultPageSize"))
	apiKey := config.GetProperty("bestbuyApiKey")
	affiliateId := config.GetProperty("bestbuyLinkShareId")
	timeout := time.Duration(config.GetIntProperty("marketplaceDefaultTimeout")) * time.Millisecond

	var start = 1
	if page > 1 {
		start = int(page) - 1*int(pageSize) + 1
	}

	if isKeywordSearch {
		listFields := config.GetProperty("bestbuyListFields")
		path := config.GetProperty("bestbuyProductSearchPath")
		url := fmt.Sprintf("%s/%s%s", endpoint, path, p[model.Keywords])
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/json")
		q := req.URL.Query()
		q.Add("format", "json")
		q.Add("apiKey", apiKey)
		q.Add("LID", affiliateId)
		q.Add("show", listFields)
		q.Add("page", strconv.Itoa(start))
		q.Add("pageSize", strconv.Itoa(pageSize))
		req.URL.RawQuery = q.Encode()
		url = fmt.Sprintf(req.URL.String())
		log.Printf("BestBuy search: %s", url)

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
		return buildSearchResponse(&entity)
	} else {

		// search trending items if no keyword provided
		url := endpoint + "/" + config.GetProperty("bestbuyProductTrendingPath")
		req, err := http.NewRequest("GET", url, nil)

		req.Header.Set("Accept", "application/json")
		q := req.URL.Query()
		q.Add("format", "json")
		q.Add("apiKey", apiKey)
		q.Add("LID", affiliateId)
		req.URL.RawQuery = q.Encode()
		url = fmt.Sprintf(req.URL.String())
		log.Printf("BestBuy trending: %s", url)

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
		return buildTrendingResponse(&entity)
	}

	return nil
}

func buildTrendingResponse(r *TrendingResponse) *model.OfferList {
	list := buildTrendingItemList(r.Results)
	o := model.NewOfferList(list, 1, 1, r.Metadata.ResultSet.Count)
	return o
}

func buildTrendingItemList(items []TrendingItem) []model.Offer {
	list := make([]model.Offer, 0)
	proxyRequired := config.IsProxyRequired(model.BestBuy)

	for _, item := range items {
		o := model.NewOffer(
			item.Sku,
			"",
			item.Names.Title,
			model.BestBuy,
			item.Links.Web,
			config.BuildImgUrlExternal(item.Images.Standard, proxyRequired),
			config.BuildImgUrl("best-buy-logo.png"),
			model.SpecialOffer,
			item.Prices.Current,
			item.CustomerReviews.AverageScore,
			item.CustomerReviews.Count,
		)

		list = append(list, *o)
	}

	return list
}

func buildSearchResponse(r *SearchResponse) *model.OfferList {
	list := buildSearchItemList(r.Products)
	o := model.NewOfferList(list, r.CurrentPage, r.TotalPages, r.Total)
	return o
}

func buildSearchItemList(items []SearchItem) []model.Offer {
	list := make([]model.Offer, 0)
	proxyRequired := config.IsProxyRequired(model.BestBuy)

	for _, item := range items {
		o := buildOffer(&item, proxyRequired)
		list = append(list, o)
	}

	return list
}

func buildCategoryPath(c []CategoryPath) string {
	names := make([]string, 0)
	for _, cp := range c {
		names = append(names, cp.Name)
	}
	return strings.Join(names, "-")
}

func filterParams(m map[string]string) map[string]string {
	p := make(map[string]string)

	// get search keyword phrase
	if m[model.Name] != "" {
		p[model.Keywords] = buildSearchPath(m[model.Name])
	}

	// get page - defaults to 1
	if m[model.Page] != "" {
		p[model.Page] = m[model.Page]
	} else {
		p[model.Page] = "1"
	}

	return p
}

// Builds search path pattern for US best buy api
// sample: input 'deals of the day' : output -> (search=deals&search=of&search=the&search=day)
func buildSearchPath(str string) string {
	s := "("
	keywords := strings.Split(str, " ")
	for _, keyword := range keywords {
		if keyword != "" {
			s += fmt.Sprintf("search=%s&", keyword)
		}
	}
	// remove last &
	s = strutil.TrimSuffix(s, "&")
	s += ")"

	return s
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

// method to map generic offer model idType to vendor specific idType string such as upc to UPC
func filterIdType(t string) string {
	switch t {
	case model.Id:
		{
			return "productId"
		}
	case model.Upc:
		{
			return model.Upc
		}
	default:
		{
			return ""
		}
	}
}

// Search for a specific product detail either by Id or Upc
func GetOfferDetail(id string, idType string, country string) *model.OfferDetail {
	log.Printf("Get Detail: %s, %s, %s", id, idType, country)

	// try to acquire lock from request Monitor
	if !monitor.IsServiceAvailable(model.BestBuy) {
		log.Printf("Unable to acquire lock from Request Monitor")
		return nil
	}

	endpoint := config.GetProperty("bestbuyEndpoint")
	path := config.GetProperty("bestbuyProductSearchPath")
	apiKey := config.GetProperty("bestbuyApiKey")
	affiliateId := config.GetProperty("bestbuyLinkShareId")
	timeout := time.Duration(config.GetIntProperty("marketplaceDefaultTimeout")) * time.Millisecond
	var idTypeProvider string

	if idTypeProvider = filterIdType(idType); idTypeProvider == "" {
		return nil
	}

	url := fmt.Sprintf("%s/%s(%s=%s)", endpoint, path, idTypeProvider, id)
	req, err := http.NewRequest("GET", url, nil)

	req.Header.Set("Accept", "application/json")
	q := req.URL.Query()
	q.Add("format", "json")
	q.Add("apiKey", apiKey)
	q.Add("LID", affiliateId)
	req.URL.RawQuery = q.Encode()
	url = fmt.Sprintf(req.URL.String())
	log.Printf("BestBuy get: %s", url)

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
	return buildProductDetail(&entity)

	return nil
}

func buildOffer(item *SearchItem, proxyRequired bool) model.Offer {

	o := model.NewOffer(
		strconv.Itoa(item.ProductId),
		item.Upc,
		item.Name,
		model.BestBuy,
		item.Url,
		config.BuildImgUrlExternal(item.Image, proxyRequired),
		config.BuildImgUrl("best-buy-logo.png"),
		buildCategoryPath(item.CategoryPath),
		item.SalePrice,
		item.CustomerReviewAverage,
		item.CustomerReviewCount,
	)

	return *o
}

func buildProductDetail(r *SearchResponse) *model.OfferDetail {
	if len(r.Products) == 0 {
		return nil
	}

	item := r.Products[0]
	proxyRequired := config.IsProxyRequired(model.BestBuy)

	o := buildOffer(&item, proxyRequired)
	detItems := make([]model.OfferDetailItem, 0)
	det := model.NewOfferDetail(
		o,
		"",
		buildProductDetailAttributes(item.Manufacturer),
		detItems,
	)

	return det
}

func buildProductDetailAttributes(attr string) map[string]string {
	if attr == "" {
		return nil
	}
	attrs := make(map[string]string)
	attrs[model.Manufacturer] = attr
	return attrs
}
