package ebay

import (
	"encoding/json"
	"fmt"
	"github.com/guilhebl/go-offer/common/config"
	"github.com/guilhebl/go-offer/common/model"
	"github.com/guilhebl/go-offer/offer/monitor"
	"github.com/guilhebl/go-worker-pool"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Creates Job for Searching offers from Ebay and returns a Channel with jobResults
func SearchOffers(m map[string]string) *job.Job {
	// create output channel
	out := job.NewJobResultChannel()

	// let's create a job with the payload
	task := NewSearchTask()
	job := job.NewJob(&task, m, out)
	return &job
}

// Searches for offers from ebay
func search(m map[string]string) *model.OfferList {
	// try to acquire lock from request Monitor
	if !monitor.IsServiceAvailable(model.Ebay) {
		log.Printf("Unable to acquire lock from Request Monitor")
		return nil
	}

	// format vendor specific params
	p := filterParams(m)

	endpoint := config.GetProperty("eBayEndpoint")
	path := config.GetProperty("eBayProductSearchPath")
	pageSize := config.GetProperty("eBayDefaultPageSize")
	securityAppName := config.GetProperty("eBaySecurityAppName")
	defaultDataFormat := config.GetProperty("eBayDefaultDataFormat")
	affiliateNetworkId := config.GetProperty("eBayAffiliateNetworkId")
	affiliateTrackingId := config.GetProperty("eBayAffiliateTrackingId")
	affiliateCustomId := config.GetProperty("ebayAffiliateCustomId")
	timeout := time.Duration(config.GetIntProperty("marketplaceDefaultTimeout")) * time.Millisecond

	url := fmt.Sprintf("%s/%s", endpoint, path)

	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("Accept", "application/json")
	q := req.URL.Query()
	q.Add("OPERATION-NAME", "findItemsByKeywords")
	q.Add("SERVICE-VERSION", "1.0.0")
	q.Add("SECURITY-APPNAME", securityAppName)
	q.Add("GLOBAL-ID", getGlobalId(p[model.Country]))
	q.Add("RESPONSE-DATA-FORMAT", defaultDataFormat)
	q.Add("affiliate.networkId", affiliateNetworkId)
	q.Add("affiliate.trackingId", affiliateTrackingId)
	q.Add("affiliate.customId", affiliateCustomId)
	q.Add("outputSelector", "PictureURLLarge") // add large picture to standard result
	q.Add("paginationInput.pageNumber", p[model.Page])
	q.Add("paginationInput.entriesPerPage", pageSize)
	q.Add(model.Keywords, p[model.Keywords])

	req.URL.RawQuery = q.Encode()
	url = fmt.Sprintf(req.URL.String())

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
}

// get Ebay global market Id
func getGlobalId(country string) string {
	switch country {

	//Canada
	case model.Canada:
		{
			return "EBAY-ENCA"
		}
	default:
		{
			return "EBAY-US"
		}
	}
}

// builds Offer list response mapping from vendor specific params
func buildSearchResponse(r *SearchResponse) *model.OfferList {
	if len(r.FindItemsByKeywordsResponse) == 0 || len(r.FindItemsByKeywordsResponse[0].PaginationOutput) == 0 {
		return nil
	}

	head := r.FindItemsByKeywordsResponse[0]
	pg := head.PaginationOutput[0]

	page, err := strconv.Atoi(pg.PageNumber[0])
	if err != nil {
		log.Println(err)
		return nil
	}

	totalPages, err := strconv.Atoi(pg.TotalPages[0])
	if err != nil {
		log.Println(err)
		return nil
	}

	total, err := strconv.Atoi(pg.TotalEntries[0])
	if err != nil {
		log.Println(err)
		return nil
	}

	list := buildSearchItemList(head.SearchResult[0].Item)
	o := model.NewOfferList(list, page, totalPages, total)
	return o
}

func buildSearchItemList(items []SearchItem) []model.Offer {
	list := make([]model.Offer, 0)
	proxyRequired := config.IsProxyRequired(model.Ebay)

	for _, item := range items {
		o := buildOffer(&item, proxyRequired)
		list = append(list, *o)
	}

	return list
}

func buildOffer(item *SearchItem, proxyRequired bool) *model.Offer {
	price, err := strconv.ParseFloat(item.SellingStatus[0].ConvertedCurrentPrice[0].Value, 32)
	if err != nil {
		log.Printf("error on parsing price for item string: %v", item)
		price = 0.0
	}

	id := ""
	if len(item.ProductID) > 0 {
		id = item.ProductID[0].Value
	} else if len(item.ItemID) > 0 {
		id = item.ItemID[0]
	}

	url := ""
	if len(item.ViewItemURL) > 0 {
		url = item.ViewItemURL[0]
	}

	imgUrl := ""
	if len(item.PictureURLLarge) > 0 {
		imgUrl = item.PictureURLLarge[0]
	}

	o := model.NewOffer(
		id,
		"",
		strings.Join(item.Title, ""),
		model.Ebay,
		url,
		config.BuildImgUrlExternal(imgUrl, proxyRequired),
		config.BuildImgUrl("ebay-logo.png"),
		strings.Join(item.PrimaryCategory[0].CategoryName, ""),
		float32(price),
		0.0,
		0,
	)
	return o
}

// filters vendor specific params from generic offer model params
func filterParams(m map[string]string) map[string]string {
	p := make(map[string]string)

	// get search keyword phrase
	if m[model.Name] != "" {
		p[model.Keywords] = m[model.Name]
	} else {
		// ebay does not have a trending api so we need to fetch random query searches
		p[model.Keywords] = getRandomSearchQuery()
	}

	// get page - defaults to 1
	if m[model.Page] != "" {
		p[model.Page] = m[model.Page]
	} else {
		p[model.Page] = "1"
	}

	if m[model.Country] != "" {
		p[model.Country] = m[model.Country]
	} else {
		p[model.Country] = model.UnitedStates
	}

	return p
}

// gets a random string inside an array of strings of queries
func getRandomSearchQuery() string {
	query := config.GetProperty("eBayDefaultSearchQuery")
	keywords := strings.Split(query, ",")
	i := rand.Intn(len(keywords))
	return keywords[i]
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

// filters Id Type for this vendor
func getIdTypeVendor(idType string) string {
	switch idType {
	case model.Id:
		return "ReferenceID"
	case model.Upc:
		return "UPC"
	default:
		return ""
	}
}

// Search for a specific product detail either by Id or Upc
func GetOfferDetail(id string, idType string, country string) *model.OfferDetail {
	log.Printf("Get Detail: %s, %s, %s", id, idType, country)

	// try to acquire lock from request Monitor
	if !monitor.IsServiceAvailable(model.Ebay) {
		log.Printf("Unable to acquire lock from Request Monitor")
		return nil
	}

	if idType == model.Id || idType == model.Upc {
		idTypeVendor := getIdTypeVendor(idType)
		endpoint := config.GetProperty("eBayEndpoint")
		path := config.GetProperty("eBayProductSearchPath")
		securityAppName := config.GetProperty("eBaySecurityAppName")
		defaultDataFormat := config.GetProperty("eBayDefaultDataFormat")
		affiliateNetworkId := config.GetProperty("eBayAffiliateNetworkId")
		affiliateTrackingId := config.GetProperty("eBayAffiliateTrackingId")
		affiliateCustomId := config.GetProperty("ebayAffiliateCustomId")
		timeout := time.Duration(config.GetIntProperty("marketplaceDefaultTimeout")) * time.Millisecond

		url := fmt.Sprintf("%s/%s", endpoint, path)
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("Accept", "application/json")
		q := req.URL.Query()

		q.Add("OPERATION-NAME", "findItemsByProduct")
		q.Add("SERVICE-VERSION", "1.0.0")
		q.Add("SECURITY-APPNAME", securityAppName)
		q.Add("GLOBAL-ID", getGlobalId(country))
		q.Add("RESPONSE-DATA-FORMAT", defaultDataFormat)
		q.Add("affiliate.networkId", affiliateNetworkId)
		q.Add("affiliate.trackingId", affiliateTrackingId)
		q.Add("affiliate.customId", affiliateCustomId)
		q.Add("outputSelector", "PictureURLLarge") // add large picture to standard result
		q.Add("productId.@type", idTypeVendor)
		q.Add("productId", id)
		q.Add("paginationInput.entriesPerPage", "1")

		req.URL.RawQuery = q.Encode()
		url = fmt.Sprintf(req.URL.String())
		log.Printf("Ebay get: %s", url)

		client := &http.Client{
			Timeout: timeout,
		}

		resp, err := client.Do(req)
		if err != nil {
			log.Fatal("Do: ", err)
			return nil
		}
		defer resp.Body.Close()

		var entity ProductDetailResponse

		if err := json.NewDecoder(resp.Body).Decode(&entity); err != nil {
			log.Println(err)
			return nil
		}
		return buildProductDetailResponse(&entity)
	}

	return nil
}

func buildProductDetail(item *SearchItem) *model.OfferDetail {
	proxyRequired := config.IsProxyRequired(model.Ebay)
	o := buildOffer(item, proxyRequired)

	attrs := make(map[string]string)
	detItems := make([]model.OfferDetailItem, 0)

	det := model.NewOfferDetail(
		*o,
		"",
		attrs,
		detItems,
	)

	return det
}

func buildProductDetailResponse(item *ProductDetailResponse) *model.OfferDetail {
	if item == nil || len(item.FindItemsByProductResponse) == 0 || len(item.FindItemsByProductResponse[0].SearchResult) == 0 {
		return nil
	}
	p := item.FindItemsByProductResponse[0].SearchResult[0].Item[0]
	return buildProductDetail(&p)
}
