package rest

import (
	"testing"

	"github.com/Lewiscowles1986/go-gorilla-api/data"
)

func TestProductListingJSONResponse(t *testing.T) {
	result := ListingJSONResponse("/products", 1, 100, 10,
		ProductsToEntries([]data.Product{}))
	expected := ProductListing{
		Data:  []data.Product{},
		Total: 100,
		Count: 0,
		Page:  1,
		Limit: 10,
		Links: []Link{
			{Href: "/products?page=1&count=10", Rel: "first", Type: "GET"},
			{Href: "/products?page=1&count=10", Rel: "current", Type: "GET"},
			{Href: "/products?page=2&count=10", Rel: "next", Type: "GET"},
			{Href: "/products?page=10&count=10", Rel: "last", Type: "GET"}}}
	if len(result.Data) != len(expected.Data) {
		t.Fatalf("Expected %d Got %d", len(expected.Data), len(result.Data))
	}
	if len(result.Links) != len(expected.Links) {
		t.Fatalf("Expected %d Got %d", len(expected.Links), len(result.Links))
	}
}

func TestProductToEntryEmpty(t *testing.T) {
	products := []data.Product{}

	results := ProductsToEntries(products)
	if len(results) != len(products) {
		t.Fatalf("Both input and output should be zero-length.\n\t%+v", results)
	}
}

func TestProductToEntrySingle(t *testing.T) {
	products := []data.Product{data.CreateProduct("test", 9.99)}

	results := ProductsToEntries(products)
	if len(results) != len(products) {
		t.Fatalf("Both input and output should be equal length (1).\n\t%+v", results)
	}
}

func TestProductToEntryMultiple(t *testing.T) {
	products := []data.Product{
		data.CreateProduct("test1", 9.99),
		data.CreateProduct("test2", 99.99)}

	results := ProductsToEntries(products)
	if len(results) != len(products) {
		t.Fatalf("Both input and output should be equal length (2).\n\t%+v", results)
	}
}

func TestProductToEntry(t *testing.T) {
	entry := ProductToEntry(data.CreateProduct("test", 9.99))
	if false {
		t.Fatalf("If it compiles we shouldn't reach this. %+v", entry)
	}
}
