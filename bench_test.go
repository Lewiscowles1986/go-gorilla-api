package main

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"github.com/Lewiscowles1986/go-gorilla-api/data"
	"github.com/Lewiscowles1986/go-gorilla-api/repositories"
)

var response *http.Response

func BenchmarkListingRequests(b *testing.B) {
	request, _ := http.NewRequest("GET", "/products", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response = executeRequest(request).Result()
	}
}
func BenchmarkCreateRequests(b *testing.B) {
	request, _ := http.NewRequest("POST", "/product", strings.NewReader(`{"name":"something","price":99.99}`))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response = executeRequest(request).Result()
	}
}
func BenchmarkGetRecordRequests(b *testing.B) {

	p := data.CreateProduct("something", 99.99)
	repositories.CreateProduct(a.DB, p)

	request, _ := http.NewRequest("GET", fmt.Sprintf("/product/%s", p.GetID()), nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response = executeRequest(request).Result()
	}
}
