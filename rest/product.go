package rest

import (
	"../data"
)

type ProductListing struct {
	Data  []data.Product `json:"data"`
	Total uint64         `json:"total"`
	Count uint8          `json:"count"`
	Page  uint64         `json:"page"`
	Limit uint8          `json:"limit"`
	Links []Link         `json:"links"`
}

type Product struct {
	Object data.Product `json:"object"`
	Links  []Link       `json:"links"`
}

func ProductToEntry(product data.Product) Entry {
	return Entry{Object: product}
}

func ProductsToEntries(products []data.Product) []Entry {
	entries := []Entry{}
	for _, p := range products {
		e := ProductToEntry(p)
		entries = append(entries, e)
	}
	return entries
}
