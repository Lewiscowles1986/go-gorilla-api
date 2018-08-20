package rest

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Link struct {
	Href string `json:"href"`
	Rel  string `json:"rel"`
	Type string `json:"type"`
}

func AddHyperMediaLinks(links *[]Link, basePath string, page, total uint64, count uint8) {
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

func GetPageLink(basePath string, count uint8, page uint64, name, method string) Link {
	url := basePath + "?" + BuildListingQuery(page, count)
	return Link{Href: url, Rel: name, Type: method}
}

func getPages(count uint8, total uint64) uint64 {
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

func RespondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
