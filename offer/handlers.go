package offer

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guilhebl/go-offer/common/model"
	"github.com/guilhebl/go-offer/common/config"
	"github.com/guilhebl/go-offer/offer/walmart"
	"strings"
)

func Index(w http.ResponseWriter, r *http.Request) {

	// offer list
	offers := SearchOffers()

	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(offers); err != nil {
		panic(err)
	}
}

func Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string
	var err error
	if id = vars["id"]; err != nil {
		panic(err)
	}

	idType := r.FormValue("idType")
	if idType == "" {
		panic("invalid IdType")
	}

	source := r.FormValue("source")
	if source == "" {
		panic(err)
	}

	offerDetail := GetOfferDetail(id, idType, source)
	if offerDetail.Offer.Id != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(offerDetail); err != nil {
			panic(err)
		}
		return
	}

	// If we didn't find it, 404
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusNotFound)
	if err := json.NewEncoder(w).Encode(model.JsonErr{Code: http.StatusNotFound, Text: "Not Found"}); err != nil {
		panic(err)
	}
}

// Searches for offers from marketplace providers
func Search(w http.ResponseWriter, r *http.Request) {

	// decode request
	decoder := json.NewDecoder(r.Body)
	var req model.ListRequest
	err := decoder.Decode(&req)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	// build param map
	m := req.Map()
	country := m["country"]

	if country == "" {
		country = model.UnitedStates
	}

	// build empty response
	capacity := config.GetIntProperty("defaultOfferListCapacity")

	list := &model.OfferList{
		List:    make([]model.Offer, 10, capacity),
		Summary: model.Summary{Page: 1, PageCount: 1, TotalCount: 0},
	}

	// search providers
	providers := getProvidersByCountry(country)

	for i := 0; i < len(providers); i++ {
		mergeSearchResponse(list, search(providers[i], m))
	}

	// set response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(list); err != nil {
		panic(err)
	}
}

func mergeSearchResponse(list *model.OfferList, list2 *model.OfferList) {
	if list2 != nil && list2.TotalCount > 0 {
		list.List = append(list.List, list2.List...)
		list.TotalCount += list2.TotalCount
		list.PageCount += list2.PageCount
	}
}

func search(provider string, m map[string]string) *model.OfferList {
	switch provider {
		case model.Walmart: return walmart.SearchOffers(m)
	}
	return nil
}

func getProvidersByCountry(country string) []string {
	switch country {
	case model.Canada:
		return strings.Split(config.GetProperty("marketplaceProvidersCanada"), ",")
	default:
		return strings.Split(config.GetProperty("marketplaceProvidersCanada"), ",")
	}
}
