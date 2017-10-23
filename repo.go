package main

import (
	"fmt"
	"encoding/json"
	"log"
	"net/http"
	"bytes"
)

func SearchOffers() OfferList {
	endpoint := GetEndpoint()
	url := fmt.Sprintf(endpoint)

	// Build the request
	l := ListRequest{}
	jsonValue, _ := json.Marshal(l)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonValue))
	var offerList OfferList

	log.Println("search: %s", req)

	if err != nil {
		log.Fatal("NewRequest: ", err)
		return offerList
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return offerList
	}
	defer resp.Body.Close()

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&offerList); err != nil {
		log.Println(err)
	}

	return offerList
}

func GetOfferDetail(id string, idType string, source string) OfferDetail {
	endpoint := GetEndpoint() + "/" + id + "?idType=" + idType + "&source=" + source
	url := fmt.Sprintf(endpoint)

	// Build the request
	req, err := http.NewRequest("GET", url, nil)
	var entity OfferDetail

	log.Println("getDetail: %s", req)

	if err != nil {
		log.Fatal("NewRequest: ", err)
		return entity
	}

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		log.Fatal("Do: ", err)
		return entity
	}
	defer resp.Body.Close()

	// Use json.Decode for reading streams of JSON data
	if err := json.NewDecoder(resp.Body).Decode(&entity); err != nil {
		log.Println(err)
	}

	return entity
}
