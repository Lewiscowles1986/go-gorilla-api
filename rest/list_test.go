package rest

import (
	"testing"
)

func TestPagesCountTotalSingle(t *testing.T) {
	result := getPages(uint8(1), uint64(100))
	if result != uint64(100) {
		t.Errorf("Expected 100 got %d", result)
	}
}

func TestPagesCountTotalEvenForEven(t *testing.T) {
	result := getPages(uint8(2), uint64(50))
	if result != uint64(25) {
		t.Errorf("Expected 25 got %d", result)
	}
}

func TestPagesCountTotalOddForOdd(t *testing.T) {
	result := getPages(uint8(3), uint64(13))
	if result != uint64(5) {
		t.Errorf("Expected 5 got %d", result)
	}
}

func TestPagesCountTotalEvenForOdd(t *testing.T) {
	result := getPages(uint8(10), uint64(97))
	if result != uint64(10) {
		t.Errorf("Expected 10 got %d", result)
	}
}

func TestPagesCountTotalOddForEven(t *testing.T) {
	result := getPages(uint8(7), uint64(100))
	if result != uint64(15) {
		t.Errorf("Expected 15 got %d", result)
	}
}

func TestGetPageLinkOutputsExpected(t *testing.T) {
	expected := Link{Href: "/products?page=1&count=50", Rel: "something", Type: "GET"}
	result := GetPageLink("/products", uint8(50), uint64(1), "something", "GET")
	if result != expected {
		t.Errorf("Expected: %+v, Got: %+v", expected, result)
	}
}

func TestHyperMediaLinksFirstPage(t *testing.T) {
	links := []Link{}
	expected := []Link{
		Link{Href: "/products?page=1&count=10", Rel: "first", Type: "GET"},
		Link{Href: "/products?page=1&count=10", Rel: "current", Type: "GET"},
		Link{Href: "/products?page=2&count=10", Rel: "next", Type: "GET"},
		Link{Href: "/products?page=10&count=10", Rel: "last", Type: "GET"},
	}
	AddHyperMediaLinks(&links, "/products", 1, uint64(100), uint8(10))
	linksLen := len(links)
	if linksLen != len(expected) && links[0] != expected[0] &&
		links[1] != expected[1] && links[2] != expected[2] &&
		links[3] != expected[3] && links[linksLen-1] != expected[linksLen-1] {
		t.Errorf("Expected: %+v, Got: %+v", expected, links)
	}
}

func TestHyperMediaLinksNthPage(t *testing.T) {
	links := []Link{}
	expected := []Link{
		Link{Href: "/products?page=1&count=7", Rel: "first", Type: "GET"},
		Link{Href: "/products?page=4&count=7", Rel: "prev", Type: "GET"},
		Link{Href: "/products?page=5&count=7", Rel: "current", Type: "GET"},
		Link{Href: "/products?page=6&count=7", Rel: "next", Type: "GET"},
		Link{Href: "/products?page=15&count=7", Rel: "last", Type: "GET"},
	}
	AddHyperMediaLinks(&links, "/products", 5, uint64(100), uint8(7))
	linksLen := len(links)
	if linksLen != len(expected) && links[0] != expected[0] &&
		links[1] != expected[1] && links[2] != expected[2] &&
		links[3] != expected[3] && links[linksLen-1] != expected[linksLen-1] {
		t.Errorf("Expected: %+v, Got: %+v", expected, links)
	}
}

func TestHyperMediaLinksLastPage(t *testing.T) {
	links := []Link{}
	expected := []Link{
		Link{Href: "/products?page=1&count=10", Rel: "first", Type: "GET"},
		Link{Href: "/products?page=99&count=10", Rel: "prev", Type: "GET"},
		Link{Href: "/products?page=100&count=10", Rel: "current", Type: "GET"},
		Link{Href: "/products?page=100&count=10", Rel: "last", Type: "GET"},
	}
	AddHyperMediaLinks(&links, "/products", uint64(100), uint64(1000), uint8(10))
	linksLen := len(links)
	if linksLen != len(expected) && links[0] != expected[0] &&
		links[1] != expected[1] && links[2] != expected[2] &&
		links[3] != expected[3] && links[linksLen-1] != expected[linksLen-1] {
		t.Errorf("Expected: %+v, Got: %+v", expected, links)
	}
}
