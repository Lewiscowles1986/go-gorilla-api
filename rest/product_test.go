package rest

import (
	"testing"

	"../data"
)

func TestProductListingJSONResponse(t *testing.T) {
	result := ProductListingJSONResponse(1, 100, 10, []data.Product{})
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
