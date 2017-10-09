package main

import (
	"fmt"
	"strconv"
)

var currentId int

var offers Offers

// Give us some seed data
func init() {
	RepoCreateOffer(
		Offer{
		Name: "Toshiba Laptop",
		Upc: "upc1",
		SemanticName: "prod1",
		PartyName: "amazon.com",
		PartyImageFileUrl: "amazon.jpg",
		MainImageFileUrl: "/img/sample_001.jpg",
		ProductCategory: "Computer, Laptops",
		NumReviews: 10,
		Price: 785.99,
		Rating: 4.5,
	})

	RepoCreateOffer(
		Offer{
			Name: "HP Laptop",
			Upc: "upc2",
			SemanticName: "prod2",
			PartyName: "bestbuy.com",
			PartyImageFileUrl: "bestbuy.jpg",
			MainImageFileUrl: "/img/sample_002.jpg",
			ProductCategory: "Computer, Laptops",
			NumReviews: 100,
			Price: 615.40,
			Rating: 3.9,
		})

}

func RepoFindOffer(id string) Offer {
	for _, t := range offers {
		if t.Id == id {
			return t
		}
	}
	// return empty Offer if not found
	return Offer{}
}

//this is bad, I don't think it passes race condtions
func RepoCreateOffer(t Offer) Offer {
	currentId += 1
	t.Id = strconv.Itoa(currentId)
	offers = append(offers, t)
	return t
}

func RepoDestroyOffer(id string) error {
	for i, t := range offers {
		if t.Id == id {
			offers = append(offers[:i], offers[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("Could not find Offer with id of %d to delete", id)
}