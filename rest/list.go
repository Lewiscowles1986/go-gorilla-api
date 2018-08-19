package rest

import (
    "database/sql"
    "encoding/json"
    "net/http"

    "../data"
    "../repositories"
)

type Link struct {
    Href    string          `json:"href"`
    Rel     string          `json:"rel"`
    Type    string          `json:"type"`
}

type ProductListing struct {
    Data    []data.Product  `json:"data"`
    Total   uint64          `json:"total"`
    Count   uint16          `json:"count"`
    Offset  uint64          `json:"offset"`
    Limit   uint8          `json:"limit"`
    Links   []Link          `json:"links"`
}

type Product struct {
    Object  data.Product    `json:"object"`
    Links   []Link          `json:"links"`
}

func ProductListingJSONResponse(db *sql.DB, start uint64, count uint8, products []data.Product) ProductListing {
    l := ProductListing{}
    l.Data = products
    l.Limit = count
    l.Offset = start
    l.Count = uint16(len(products))
    l.Total = repositories.GetProductCount(db)
    l.Links = append(l.Links, Link{Href:"/products", Rel:"first", Type:"GET"})

    return l
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
