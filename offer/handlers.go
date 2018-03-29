package offer

import (
	"encoding/json"
	"net/http"

	"fmt"
	"github.com/gorilla/mux"
	"github.com/guilhebl/go-offer/common/config"
	"github.com/guilhebl/go-offer/common/model"
)

// handles error conditions coming from service layer
func handleErr(errMsg string, w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	switch errMsg {
	case model.InvalidRequest:
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.InvalidRequest)
	default:
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.InternalError)
	}
}

// Searches with no keywords for Trending and Promotional Deals in each marketplace provider
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Index: %s", r.URL)

	// search with empty keyword
	req := model.NewEmptyListRequest(config.GetIntProperty("defaultRowsPerPage"))
	result, err := SearchOffers(req)
	if err != nil {
		handleErr(err.Error(), w)
		return
	}

	// set ok response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

// Searches for offers from marketplace providers using keyword and other filters
func Search(w http.ResponseWriter, r *http.Request) {

	defer r.Body.Close()

	// decode request
	decoder := json.NewDecoder(r.Body)
	var req model.ListRequest
	var err error
	if err = decoder.Decode(&req); err != nil {
		handleErr(err.Error(), w)
		return
	}

	result, err := SearchOffers(&req)
	if err != nil {
		handleErr(err.Error(), w)
		return
	}

	// set response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

// Reset Db
func ResetDatastore(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Reset Datastore: %s", r.URL)

	err := ResetDb()
	if err != nil {
		handleErr(err.Error(), w)
		return
	}

	// set ok response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode("OK"); err != nil {
		panic(err)
	}
}

// Searches for all offers in Db
func SearchDatastore(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Search Datastore: %s", r.URL)

	req := model.NewEmptyListRequest(config.GetIntProperty("defaultRowsPerPage"))
	result, err := SearchOffersDb(req)
	if err != nil {
		handleErr(err.Error(), w)
		return
	}

	// set ok response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

// Adds to Datastore a new offer
func AddOffer(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Add to Datastore: %s", r.URL)

	defer r.Body.Close()

	// decode request
	decoder := json.NewDecoder(r.Body)
	var req model.Offer
	var err error
	if err = decoder.Decode(&req); err != nil {
		fmt.Printf("ERROR ON DECODE %e", err.Error())
		handleErr(err.Error(), w)
		return
	}

	result, err := AddOfferDb(&req)
	if err != nil {
		handleErr(err.Error(), w)
		return
	}

	// set response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(result); err != nil {
		panic(err)
	}
}

// Get Offer Detail from marketplace provider and associated competitors
func Show(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var id string
	var err error
	if id = vars["id"]; err != nil {
		handleErr(model.InvalidRequest, w)
		return
	}

	idType := r.FormValue("idType")
	source := r.FormValue("source")
	country := r.FormValue("country")

	request := model.NewDetailRequest(id, idType, source, country)
	result, err := GetOfferDetail(request)
	if err != nil {
		handleErr(err.Error(), w)
		return
	}

	if result != nil && result.Offer.Id != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(*result); err != nil {
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
