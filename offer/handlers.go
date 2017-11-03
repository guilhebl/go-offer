package offer

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guilhebl/offergo/common/model"
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
