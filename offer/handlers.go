package offer

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guilhebl/go-offer/common/model"
)

// Searches for Trending and Promotional Deals in each marketplace provider
func Index(w http.ResponseWriter, r *http.Request) {
	// search with empty keyword
	m := make(map[string]string)
	result := SearchOffers(m)
	if result == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode("internal error"); err != nil {
			panic(err)
		}
		return
	}

	// set response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(result); err != nil {
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

	m := req.Map()
	result := SearchOffers(m)
	if result == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode("internal error"); err != nil {
			panic(err)
		}
		return
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
	}

	idType := r.FormValue("idType")
	if idType == "" {
		panic("invalid IdType")
	}

	source := r.FormValue("source")
	if source == "" {
		panic(err)
	}

	country := r.FormValue("country")
	if source == "" {
		country = model.UnitedStates
	}

	result := GetOfferDetail(id, idType, source, country)
	if result == nil {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode("internal error"); err != nil {
			panic(err)
		}
		return
	}

	if result.Offer.Id != "" {
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
