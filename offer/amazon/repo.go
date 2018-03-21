package amazon

import (
	"github.com/guilhebl/go-offer/common/config"
	"github.com/guilhebl/go-offer/common/model"
	"github.com/guilhebl/go-offer/common/util"
	"github.com/guilhebl/go-offer/offer/monitor"
	"github.com/guilhebl/go-worker-pool"
	"log"
	"math/rand"
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Creates Job for Searching offers from Amazon and returns a Channel with jobResults
func SearchOffers(m map[string]string) *job.Job {
	// create output channel
	out := job.NewJobResultChannel()

	// let's create a job with the payload
	task := NewSearchTask()
	job := job.NewJob(&task, m, out)
	return &job
}

// Searches for offers from amazon
func search(m map[string]string) *model.OfferList {
	// try to acquire lock from request Monitor
	if !monitor.IsServiceAvailable(model.Amazon) {
		log.Printf("Unable to acquire lock from Request Monitor")
		return nil
	}

	// format vendor specific params
	p := filterParams(m)
	page, err := strconv.Atoi(p[model.Page])

	if err != nil {
		log.Printf("format page error: %s", err.Error())
		return nil
	}

	region := config.GetProperty("amazonDefaultRegion")
	accessKeyId := config.GetProperty("amazonAccessKeyId")
	secretKey := config.GetProperty("amazonSecretKey")
	associateTag := config.GetProperty("amazonAssociateTag")
	timeout := time.Duration(config.GetIntProperty("marketplaceDefaultTimeout")) * time.Millisecond

	cfg := NewConfig(accessKeyId, secretKey, associateTag, region, true)
	client := NewClient(cfg)

	query := ItemSearchQuery{
		SearchIndex:    "All",
		Keywords:       p[model.Keywords],
		ItemPage:       p[model.Page],
		ResponseGroups: []string{"Images", "ItemAttributes", "Offers"},
	}
	response, err := client.ItemSearch(query, timeout)

	if err != nil {
		log.Printf("%s", err.Error())
		return nil
	}

	return buildSearchResponse(response, page)
}

// builds Offer list response mapping from vendor specific params
func buildSearchResponse(r *ItemSearchResponse, page int) *model.OfferList {
	items := r.Items
	total := len(items.Items)
	if total == 0 {
		return nil
	}
	totalPages := items.TotalPages

	list := buildSearchItemList(items.Items)
	o := model.NewOfferList(list, page, totalPages, total)
	return o
}

func buildSearchItemList(items []Item) []model.Offer {
	list := make([]model.Offer, 0)
	proxyRequired := config.IsProxyRequired(model.Amazon)

	for _, item := range items {
		o := buildOffer(&item, proxyRequired)
		list = append(list, *o)
	}

	return list
}

func buildOffer(item *Item, proxyRequired bool) *model.Offer {
	itemAttrs := item.ItemAttributes
	summary := item.OfferSummary

	imgUrl := ""
	if item.LargeImage != nil {
		imgUrl = item.LargeImage.URL
	}

	upc := ""
	if itemAttrs != nil && itemAttrs.UPC != "" {
		upc = itemAttrs.UPC
	}

	title := ""
	if itemAttrs != nil && itemAttrs.Title != "" {
		title = itemAttrs.Title
	}

	productGroup := ""
	if itemAttrs != nil && itemAttrs.ProductGroup != "" {
		productGroup = itemAttrs.ProductGroup
	}

	o := model.NewOffer(
		util.GenerateStringUUID(),
		item.ASIN,
		upc,
		title,
		model.Amazon,
		item.DetailPageURL,
		config.BuildImgUrlExternal(imgUrl, proxyRequired),
		config.BuildImgUrl("amazon-logo.png"),
		productGroup,
		buildPrice(summary.LowestNewPrice, summary.LowerUsedPrice),
		0.0,
		0,
	)

	return o
}

// builds amazon item price based on lowestNew or lowestUsed
func buildPrice(lowestNew, lowestUsed Price) float32 {
	if lowestNew.FormattedPrice != "" {
		return getFormattedPriceValue(lowestNew.FormattedPrice)
	} else if lowestUsed.FormattedPrice != "" {
		return getFormattedPriceValue(lowestUsed.FormattedPrice)
	}

	return 0.0
}

func getFormattedPriceValue(priceStr string) float32 {
	var re = regexp.MustCompile(`[^\d.]`)
	formatted := re.ReplaceAllString(priceStr, "")
	price, err := strconv.ParseFloat(formatted, 32)
	if err != nil {
		log.Printf("error on parsing price for item string: %s", priceStr)
		return 0.0
	}
	return float32(price)
}

// filters vendor specific params from generic offer model params
func filterParams(m map[string]string) map[string]string {
	p := make(map[string]string)

	// get search keyword phrase
	if m[model.Name] != "" {
		p[model.Keywords] = m[model.Name]
	} else {
		// amazon does not have a trending api so we need to fetch random query searches
		p[model.Keywords] = getRandomSearchQuery()
	}

	// get page - defaults to 1
	if m[model.Page] != "" {
		p[model.Page] = m[model.Page]
	} else {
		p[model.Page] = "1"
	}
	return p
}

// gets a random string inside an array of strings of queries
func getRandomSearchQuery() string {
	query := config.GetProperty("amazonDefaultSearchQuery")
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
		return "ASIN"
	case model.Upc:
		return "UPC"
	case model.Ean:
		return "EAN"
	case model.Isbn:
		return "ISBN"
	default:
		return ""
	}
}

// Search for a specific product detail either by Id or Upc
func GetOfferDetail(id string, idType string, country string) *model.OfferDetail {
	log.Printf("Get Detail: %s, %s, %s", id, idType, country)

	// try to acquire lock from request Monitor
	if !monitor.IsServiceAvailable(model.Amazon) {
		log.Printf("Unable to acquire lock from Request Monitor")
		return nil
	}

	if idType == model.Id || idType == model.Upc {
		idTypeVendor := getIdTypeVendor(idType)

		region := config.GetProperty("amazonDefaultRegion")
		accessKeyId := config.GetProperty("amazonAccessKeyId")
		secretKey := config.GetProperty("amazonSecretKey")
		associateTag := config.GetProperty("amazonAssociateTag")
		timeout := time.Duration(config.GetIntProperty("marketplaceDefaultTimeout")) * time.Millisecond

		cfg := NewConfig(accessKeyId, secretKey, associateTag, region, true)
		client := NewClient(cfg)

		indexSearch := ""
		if idTypeVendor != "ASIN" {
			indexSearch = "All"
		}

		query := ItemLookupQuery{
			SearchIndex:    indexSearch,
			IDType:         idTypeVendor,
			ItemIDs:        []string{id},
			ResponseGroups: []string{"Images", "ItemAttributes", "Offers"},
		}

		response, err := client.ItemLookup(query, timeout)
		if err != nil {
			log.Printf("error: %s", err)
			return nil
		}
		return buildProductDetailResponse(response)
	}

	return nil
}

func buildAttributes(a *ItemAttributes) map[string]string {
	attrs := make(map[string]string)

	if a == nil {
		return attrs
	}

	m := a.Model
	if m != "" {
		attrs[model.Model] = m
	}

	manufacturer := a.Manufacturer
	if manufacturer != "" {
		attrs[model.Manufacturer] = manufacturer
	}

	brand := a.Brand
	if brand != "" {
		attrs[model.Brand] = brand
	}

	publisher := a.Publisher
	if publisher != "" {
		attrs[model.Publisher] = publisher
	}

	return attrs
}

func buildProductDetailResponse(response *ItemLookupResponse) *model.OfferDetail {
	items := response.Items
	item := items.Item
	if item.ASIN == "" {
		return nil
	}

	o := buildOffer(&item, config.IsProxyRequired(model.Amazon))

	desc := ""
	if item.ItemAttributes != nil {
		desc = item.ItemAttributes.Feature
	}

	attrs := buildAttributes(item.ItemAttributes)
	detItems := make([]model.OfferDetailItem, 0)
	return model.NewOfferDetail(*o, desc, attrs, detItems)
}
