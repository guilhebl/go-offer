package offer

import (
	"encoding/json"
	"errors"
	"github.com/guilhebl/go-offer/common/config"
	"github.com/guilhebl/go-offer/common/model"
	"github.com/guilhebl/go-offer/offer/amazon"
	"github.com/guilhebl/go-offer/offer/bestbuy"
	"github.com/guilhebl/go-offer/offer/ebay"
	"github.com/guilhebl/go-offer/offer/walmart"
	"github.com/guilhebl/go-worker-pool"
	"github.com/guilhebl/xcrypto"
	"log"
	"math/rand"
	"sort"
	"strings"
	"github.com/guilhebl/go-offer/common/db"
)

// searches offers - tries to fetch 1st in cache if not found calls marketplace
func SearchOffers(r *model.ListRequest) (*model.OfferList, error) {
	// validates request before querying marketplace
	if !r.IsValid() {
		return nil, errors.New(model.InvalidRequest)
	}

	// transform request
	jsonReq, _ := json.Marshal(&r)
	key := string(jsonReq)
	if key == "" {
		key = model.Trending
	}

	// search first in cache
	hash := xcrypto.GenerateSHA1(key)
	cacheEnabled := isCacheEnabled()

	var obj *model.OfferList
	if cacheEnabled {
		if data, _ := GetInstance().RedisCache.Get(hash); data != "" {
			if err := json.Unmarshal([]byte(data), &obj); err != nil {
				return nil, err
			}
			return obj, nil
		}
	}

	// if not found in cache search and store valid output in cache
	obj = searchOffers(r.Map())
	if cacheEnabled && obj != nil {
		data, _ := json.Marshal(&obj)
		err := GetInstance().RedisCache.Set(hash, string(data))
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
}

// resets Db
func ResetDb() error {
	return db.Reset()
}

// searches offers in Db
func SearchOffersDb(r *model.ListRequest) (*model.OfferList, error) {
	// validates request before querying marketplace
	if !r.IsValid() {
		return nil, errors.New(model.InvalidRequest)
	}

	list := model.NewOfferList(make([]model.Offer, 0, 40), 1, 1, 0)
	var err error
	list.List, err = db.GetOffers()
	return list, err
}

// add offer to Db - returns offer with new id
func AddOfferDb(r *model.Offer) (*model.Offer, error) {
	// validates request before querying marketplace
	if r.Name == "" {
		return nil, errors.New(model.InvalidRequest)
	}

	return db.InsertOffer(r)
}

// Searches marketplace providers by keyword
func searchOffers(m map[string]string) *model.OfferList {
	log.Printf("Search: %v", m)

	country := m[model.Country]
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

	// sort list
	sortList(list, country, m[model.Name], m[model.SortBy], m[model.SortOrder] == "asc")

	return list
}

// sorts offer list by field
func sortList(list *model.OfferList, country, keyword, sortBy string, asc bool) {
	switch sortBy {

	case model.Id:
		sort.Slice(list.List, func(i, j int) bool {
			ret := list.List[i].Id < list.List[j].Id
			if asc {
				return ret
			}
			return !ret
		})
	case model.Name:
		sort.Slice(list.List, func(i, j int) bool {
			ret := list.List[i].Name < list.List[j].Name
			if asc {
				return ret
			}
			return !ret
		})
	case model.Price:
		sort.Slice(list.List, func(i, j int) bool {
			ret := list.List[i].Price < list.List[j].Price
			if asc {
				return ret
			}
			return !ret
		})
	case model.Rating:
		sort.Slice(list.List, func(i, j int) bool {
			ret := list.List[i].Rating < list.List[j].Rating
			if asc {
				return ret
			}
			return !ret
		})
	case model.NumReviews:
		sort.Slice(list.List, func(i, j int) bool {
			ret := list.List[i].NumReviews < list.List[j].NumReviews
			if asc {
				return ret
			}
			return !ret
		})
	default:
		if keyword != "" {
			sortByBestResults(list, keyword)
		} else {
			sortGroupedByProvider(list, country)
		}
	}
}

// Groups the results in buckets with each provider appearing in the first group of results
func sortGroupedByProvider(list *model.OfferList, country string) {
	m := make(map[string][]model.Offer)

	// split providers into buckets
	providers := getProvidersByCountry(country)
	for _, p := range providers {
		m[p] = make([]model.Offer, 0)
	}

	// fill each providers queue
	for _, o := range list.List {
		m[o.PartyName] = append(m[o.PartyName], o)
	}

	// shuffle each queue
	for _, p := range providers {
		src := m[p]
		dest := make([]model.Offer, len(src))
		perm := rand.Perm(len(src))
		for i, v := range perm {
			dest[v] = src[i]
		}
		m[p] = dest
	}

	// create output list
	listSorted := make([]model.Offer, 0)

	// keep filling groups of items from each provider queue until empty
	for i := 0; i < len(list.List); i++ {

		// fetch a random index of provider slice
		idx := rand.Intn(len(providers))
		p := providers[idx]

		// pick first element from queue if not empty
		if len(m[p]) > 0 {
			offerList := m[p]
			o := offerList[0]

			// append to output list
			listSorted = append(listSorted, o)

			// remove element from queue
			m[p] = append(offerList[:0], offerList[1:]...)
		}

		// remove provider from next round
		providers = append(providers[:idx], providers[idx+1:]...)

		// reset list of providers if all removed
		if len(providers) == 0 {
			providers = getProvidersByCountry(country)
		}
	}

	list.List = listSorted
}

// ranks by keyword distance and sorts list based on ranking found
func sortByBestResults(list *model.OfferList, keyword string) {

	// sort source list
	dest := make([]model.Offer, len(list.List))
	perm := rand.Perm(len(list.List))
	for i, v := range perm {
		dest[v] = list.List[i]
	}
	list.List = dest

	// filter keywords
	keywords := strings.Split(keyword, " ")

	// filter all blank or empty spaces
	words := make([]string, 0)
	for _, w := range keywords {
		t := strings.TrimSpace(w)
		if t != "" {
			words = append(words, t)
		}
	}

	if len(words) == 0 {
		return
	}

	// with remaining words search for distance ranking for all offers
	rankings := make([]model.OfferKeywordRank, 0)

	for i := 0; i < len(list.List); i++ {
		offer := list.List[i]
		numKeywords, totalMatches, lowestIndex := calculateKeywordsRank(words, &offer)
		offerRank := model.NewOfferKeywordRank(offer, numKeywords, totalMatches, lowestIndex)
		rankings = append(rankings, *offerRank)
	}

	// sort list by ranking
	sort.Slice(rankings, func(i, j int) bool {

		// try by num unique keywords
		if rankings[i].NumKeywords != rankings[j].NumKeywords {
			return rankings[i].NumKeywords > rankings[j].NumKeywords
		}

		// if same try by Index
		if rankings[i].LowestIndex != rankings[j].LowestIndex {
			return rankings[i].LowestIndex < rankings[j].LowestIndex
		}

		// try by total matches
		return rankings[i].TotalMatches > rankings[j].TotalMatches
	})

	// copy rankings offers to new list
	listSorted := make([]model.Offer, 0)

	for i := 0; i < len(rankings); i++ {
		listSorted = append(listSorted, rankings[i].Offer)
	}

	list.List = listSorted
}

// calculates the ranking based on keywords from a slice that are found in the title of this offer, num keywords found,
// lowest index found for any of the keywords
func calculateKeywordsRank(keywords []string, o *model.Offer) (int, int, int) {

	numKeywords, totalMatches, i := 0, 0, -1

	// to keep track of the unique found already
	found := make(map[string]bool)
	text := strings.ToLower(o.Name)

	for _, w := range keywords {
		word := strings.ToLower(w)
		j := strings.Index(text, word)

		if j != -1 {
			if !found[word] {
				numKeywords += 1
				found[word] = true
			}
			if i == -1 || i > j {
				i = j
			}
			totalMatches += 1
		}
	}

	return numKeywords, totalMatches, i
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
func GetOfferDetail(r *model.DetailRequest) (*model.OfferDetail, error) {
	// validate and transform request before querying marketplace
	if !r.IsValid() {
		return nil, errors.New(model.InvalidRequest)
	}
	jsonReq, _ := json.Marshal(&r)
	key := string(jsonReq)

	// search first in cache
	hash := xcrypto.GenerateSHA1(key)
	cacheEnabled := isCacheEnabled()

	var obj *model.OfferDetail
	if cacheEnabled {
		if data, _ := GetInstance().RedisCache.Get(hash); data != "" {
			if err := json.Unmarshal([]byte(data), &obj); err != nil {
				return nil, err
			}
			return obj, nil
		}
	}

	// store valid output in cache
	obj = getDetail(r.Id, r.IdType, r.Source, r.Country)

	// if product has Upc fetch competitors details in parallel using worker pool jobs
	if obj != nil && obj.Offer.Upc != "" {
		providers := getProvidersByCountry(r.Country)

		// create a slice of jobResult outputs
		jobOutputs := make([]<-chan job.JobResult, 0)

		for i := 0; i < len(providers); i++ {
			if p := providers[i]; p != r.Source {
				job := getDetailJob(obj.Offer.Upc, model.Upc, providers[i], r.Country)
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

				obj.ProductDetailItems = append(obj.ProductDetailItems, *detItem)
			}
		}
	}

	// store in cache if possible
	if cacheEnabled && obj != nil {
		data, err := json.Marshal(&obj)
		if err != nil {
			return nil, err
		}

		err = GetInstance().RedisCache.Set(hash, string(data))
		if err != nil {
			return nil, err
		}
	}

	return obj, nil
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

func isCacheEnabled() bool {
	return config.GetBoolProperty("cacheEnabled")
}
