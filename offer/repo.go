package offer

import (
	"encoding/json"
	"github.com/guilhebl/go-offer/common/config"
	"github.com/guilhebl/go-offer/common/model"
	"github.com/guilhebl/go-offer/offer/amazon"
	"github.com/guilhebl/go-offer/offer/bestbuy"
	"github.com/guilhebl/go-offer/offer/ebay"
	"github.com/guilhebl/go-offer/offer/walmart"
	"github.com/guilhebl/go-worker-pool"
	"github.com/guilhebl/xcrypto"
	"log"
	"strings"
)

// searches offers - tries to fetch 1st in cache if not found calls marketplace
func SearchOffers(r *model.ListRequest) (*model.OfferList, error) {
	json, _ := json.Marshal(&r)
	key := string(json)
	if key == "" {
		key = model.Trending
	}

	// validate and tranform request before querying marketplace
	m, err := r.Map()
	if err != nil {
		return nil, err
	}

	// search first in cache
	hash := xcrypto.GenerateSHA1(key)
	cacheEnabled := config.GetBoolProperty("cacheEnabled")

	var obj *model.OfferList
	if cacheEnabled {
		obj, err := GetInstance().RedisCache.GetOfferList(hash)
		if obj != nil && err == nil {
			return obj, nil
		}
	}

	// store valid output in cache
	obj = searchOffers(m)
	if cacheEnabled && obj != nil {
		err = GetInstance().RedisCache.SetOfferList(hash, obj)
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// Searches marketplace providers by keyword
func searchOffers(m map[string]string) *model.OfferList {
	log.Printf("Search: %v", m)

	country := m["country"]
	if country == "" {
		country = model.UnitedStates
	}

	// build empty response
	rowsPerPage := int(config.GetIntProperty("defaultRowsPerPage"))
	numProviders := config.CountMarketplaceProviderListSize()
	capacity := numProviders * rowsPerPage
	list := model.NewOfferList(make([]model.Offer, 0, capacity), 1, 1, 0)

	// search providers
	providers := getProvidersByCountry(country)

	// create a slice of jobResult outputs
	jobOutputs := make([]<-chan job.JobResult, 0)

	for i := 0; i < len(providers); i++ {
		job := search(providers[i], m)
		if job != nil {
			jobOutputs = append(jobOutputs, job.ReturnChannel)
			// Push each job onto the queue.
			GetInstance().JobQueue <- *job
		}
	}

	// Consume the merged output from all jobs
	out := job.Merge(jobOutputs...)
	for r := range out {
		if r.Error == nil {
			mergeSearchResponse(list, r.Value.(*model.OfferList))
		}
	}
	return list
}

func mergeSearchResponse(list *model.OfferList, list2 *model.OfferList) {
	if list2 != nil && list2.TotalCount > 0 {
		list.List = append(list.List, list2.List...)
		list.TotalCount += list2.TotalCount
		list.PageCount += list2.PageCount
	}
}

// searches create a new Job to search in a provider that returns a OfferList channel
func search(provider string, m map[string]string) *job.Job {
	switch provider {
	case model.Amazon:
		return amazon.SearchOffers(m)
	case model.Walmart:
		return walmart.SearchOffers(m)
	case model.BestBuy:
		return bestbuy.SearchOffers(m)
	case model.Ebay:
		return ebay.SearchOffers(m)
	}

	return nil
}

func getProvidersByCountry(country string) []string {
	switch country {
	case model.Canada:
		return strings.Split(config.GetProperty("marketplaceProvidersCanada"), ",")
	default:
		return strings.Split(config.GetProperty("marketplaceProviders"), ",")
	}
}

// Gets Product Detail from marketplace provider by Id and IdType, fetching competitors prices using UPC
func GetOfferDetail(id, idType, source, country string) (*model.OfferDetail, error) {
	det := getDetail(id, idType, source, country)

	// if product has Upc fetch competitors details in parallel using worker pool jobs
	if det != nil && det.Offer.Upc != "" {
		providers := getProvidersByCountry(country)

		// create a slice of jobResult outputs
		jobOutputs := make([]<-chan job.JobResult, 0)

		for i := 0; i < len(providers); i++ {
			if p := providers[i]; p != source {
				job := getDetailJob(det.Offer.Upc, model.Upc, providers[i], country)
				if job != nil {
					jobOutputs = append(jobOutputs, job.ReturnChannel)
					// Push each job onto the queue.
					GetInstance().JobQueue <- *job
				}
			}
		}

		// Consume the merged output from all jobs
		out := job.Merge(jobOutputs...)
		for r := range out {
			if r.Error == nil {
				// build detail item
				d := r.Value.(*model.OfferDetail)

				detItem := model.NewOfferDetailItem(
					d.Offer.PartyName,
					d.Offer.SemanticName,
					d.Offer.PartyImageFileUrl,
					d.Offer.Price,
					d.Offer.Rating,
					d.Offer.NumReviews)

				det.ProductDetailItems = append(det.ProductDetailItems, *detItem)
			}
		}
	}

	return det, nil
}

// creates a job to fetch a product detail from a given source using id and idType and country
func getDetailJob(id, idType, source, country string) *job.Job {
	log.Printf("getDetail Job: %s, %s, %s, %s", id, idType, source, country)

	switch source {
	case model.Amazon:
		return amazon.GetDetailJob(id, idType, country)
	case model.Walmart:
		return walmart.GetDetailJob(id, idType, country)
	case model.BestBuy:
		return bestbuy.GetDetailJob(id, idType, country)
	case model.Ebay:
		return ebay.GetDetailJob(id, idType, country)
	}

	return nil
}

// gets a product detail from a given source using id and idType and country
func getDetail(id, idType, source, country string) *model.OfferDetail {
	log.Printf("get: %s, %s, %s, %s", id, idType, source, country)

	switch source {
	case model.Amazon:
		return amazon.GetOfferDetail(id, idType, country)
	case model.Walmart:
		return walmart.GetOfferDetail(id, idType, country)
	case model.BestBuy:
		return bestbuy.GetOfferDetail(id, idType, country)
	case model.Ebay:
		return ebay.GetOfferDetail(id, idType, country)
	}
	return nil
}
