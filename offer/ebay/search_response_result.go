package ebay

import "time"

type SearchResponseResult struct {
	Ack          []string    `json:"ack"`
	Version      []string    `json:"version"`
	Timestamp    []time.Time `json:"timestamp"`
	SearchResult []struct {
		Count string       `json:"@count"`
		Item  []SearchItem `json:"item"`
	} `json:"searchResult"`
	PaginationOutput []struct {
		PageNumber     []string `json:"pageNumber"`
		EntriesPerPage []string `json:"entriesPerPage"`
		TotalPages     []string `json:"totalPages"`
		TotalEntries   []string `json:"totalEntries"`
	} `json:"paginationOutput"`
	ItemSearchURL []string `json:"itemSearchURL"`
}
