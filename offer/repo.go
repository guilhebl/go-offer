package offer

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/guilhebl/go-offer/common/config"
	"github.com/guilhebl/go-offer/common/model"
	"github.com/guilhebl/go-offer/offer/walmart"
	"strings"
)

// Searches marketplace providers by keyword
func SearchOffers(m map[string]string) *model.OfferList {

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

	return list
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
	case model.Walmart:
		return walmart.SearchOffers(m)
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
func GetOfferDetail(id, idType, source, country string) *model.OfferDetail {
	det := getDetail(id, idType, source, country)

	if det != nil && det.Offer.Upc != "" {
		providers := getProvidersByCountry(country)
		for i := 0; i < len(providers); i++ {
			if p := providers[i]; p != source {
				d := getDetail(det.Offer.Upc, model.Upc, providers[i], country)

				// build detail item
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

	return det
}

func getDetail(id, idType, source, country string) *model.OfferDetail {
	switch source {
	case model.Walmart:
		return walmart.GetOfferDetail(id, idType, country)
	}
	return nil
}
