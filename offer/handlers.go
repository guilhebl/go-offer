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
		if err := json.NewEncoder(w).Encode(model.InvalidRequest); err != nil {
			panic(err)
		}
		return
	default:
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(model.InternalError); err != nil {
			panic(err)
		}
		return
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
	// decode request
	decoder := json.NewDecoder(r.Body)
	var req model.ListRequest
	err := decoder.Decode(&req)
	if err != nil {
		panic(err)
	}
	defer r.Body.Close()

	result, err := SearchOffers(&req)
	if err != nil {
		handleErr(err.Error(), w)
	}

	// set response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
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
		panic(err)
		handleErr(model.InvalidRequest, w)
	}

	idType := r.FormValue("idType")
	source := r.FormValue("source")
	country := r.FormValue("country")

	if idType == "" || source == "" || id == "" {
		handleErr(model.InvalidRequest, w)
	}

	result, err := GetOfferDetail(id, idType, source, country)
	if err != nil {
		handleErr(err.Error(), w)
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
