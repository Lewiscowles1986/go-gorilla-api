package main_test

import (
	"net/http"
	"testing"
)

var response *http.Response

func BenchmarkListingRequests(b *testing.B) {
	request, _ := http.NewRequest("GET", "/products", nil)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		response = executeRequest(request).Result()
	}
}
