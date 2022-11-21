package rest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Listing struct {
	Data  []Entry `json:"data"`
	Total uint64  `json:"total"`
	Count uint8   `json:"count"`
	Page  uint64  `json:"page"`
	Limit uint8   `json:"limit"`
	Links []Link  `json:"links"`
}

type Entry struct {
	Object interface{} `json:"object"`
	Links  []Link      `json:"links"`
}

type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
	Type string `json:"type"`
}

func AddHyperMediaLinks(links *[]Link, basePath string, page, total uint64, count uint8) {
	if page < 1 {
		page = 1
	}
	pages := getPages(count, total)
	*links = append(*links, GetPageLink(basePath, count, 1, "first", "GET"))
	if page > 1 {
		prevPg := page - 1
		*links = append(*links, GetPageLink(basePath, count, prevPg, "prev", "GET"))
	}
	*links = append(*links, GetPageLink(basePath, count, page, "current", "GET"))
	if page < pages {
		nextPg := page + 1
		*links = append(*links, GetPageLink(basePath, count, nextPg, "next", "GET"))
	}
	*links = append(*links, GetPageLink(basePath, count, pages, "last", "GET"))
}

func ListingJSONResponse(basePath string, page, total uint64, count uint8, entries []Entry) Listing {
	if page < 1 {
		page = 1
	}
	l := Listing{}
	l.Data = entries
	l.Limit = count
	l.Page = page
	l.Count = uint8(len(entries))
	l.Total = total
	AddHyperMediaLinks(&l.Links, basePath, page, total, count)

	return l
}

func GetPageLink(basePath string, count uint8, page uint64, name, method string) Link {
	url := basePath + "?" + BuildListingQuery(page, count)
	return Link{Href: url, Rel: name, Type: method}
}

func getPages(count uint8, total uint64) uint64 {
	if total < 1 || count < 1 {
		return 1
	}
	pages := total / uint64(count)
	carry := total % uint64(count)
	if carry != 0 {
		pages++
	}
	return pages
}

func BuildListingQuery(page uint64, count uint8) string {
	return fmt.Sprintf("page=%d&count=%d", page, count)
}

func RespondWithError(w http.ResponseWriter, code int, message string) {
	RespondWithJSON(w, code, map[string]string{"error": message})
}

func JSONMarshal(t interface{}) ([]byte, error) {
    buffer := &bytes.Buffer{}
    encoder := json.NewEncoder(buffer)
    encoder.SetEscapeHTML(false)
    err := encoder.Encode(t)
    return buffer.Bytes(), err
}

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := JSONMarshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	w.Write(response)
}
