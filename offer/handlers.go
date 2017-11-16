package offer

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/guilhebl/go-offer/common/model"
	"github.com/guilhebl/go-worker-pool"
)

// Searches for Trending and Promotional Deals in each marketplace provider
func Index(w http.ResponseWriter, r *http.Request) {
	// search with empty keyword
	m := make(map[string]string)

	// let's create a job with the payload
	ret := make(chan interface{})
	defer close(ret)
	work := worker.NewJob(JobTypeSearch, m, ret)

	// Push the work onto the queue.
	module := GetInstance()
	module.JobQueue <- work

	// wait for response from Job
	resp := <-ret

	var list *model.OfferList

	switch resp.(type) {
	case model.OfferList:
		{
			list = resp.(*model.OfferList)
		}
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	// set response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(list); err != nil {
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

	// let's create a job with the payload
	ret := make(chan interface{})
	defer close(ret)
	m := req.Map()
	work := worker.NewJob(JobTypeSearch, m, ret)

	// Push the work onto the queue.
	module := GetInstance()
	module.JobQueue <- work

	// wait for response from Job
	resp := <-ret

	var list *model.OfferList

	switch resp.(type) {
	case model.OfferList:
		{
			list = resp.(*model.OfferList)
		}
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	// set response
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(list); err != nil {
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

	// let's create a job with the payload
	ret := make(chan interface{})
	defer close(ret)
	work := worker.NewJob(JobTypeGet, buildGetParams(id, idType, source, country), ret)

	// Push the work onto the queue.
	module := GetInstance()
	module.JobQueue <- work

	// wait for response from Job
	resp := <-ret

	var detail *model.OfferDetail

	switch resp.(type) {
	case model.OfferDetail:
		{
			detail = resp.(*model.OfferDetail)
		}
	default:
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusInternalServerError)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			panic(err)
		}
		return
	}

	if detail.Offer.Id != "" {
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(*detail); err != nil {
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

func buildGetParams(id, idType, source, country string) map[string]string {
	m := make(map[string]string, 0)
	m["id"] = id
	m["idType"] = idType
	m["source"] = source
	m["country"] = country
	return m
}
