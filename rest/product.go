package rest

import (
	"database/sql"

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

func ProductListingJSONResponse(db *sql.DB, page, total uint64, count uint8, products []data.Product) ProductListing {

	l := ProductListing{}
	l.Data = products
	l.Limit = count
	l.Page = page
	l.Count = uint8(len(products))
	l.Total = total
	AddHyperMediaLinks(&l.Links, "/products", page, total, count)

	return l
}
