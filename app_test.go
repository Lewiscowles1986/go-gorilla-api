package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	uuid "github.com/satori/go.uuid"

	"github.com/Lewiscowles1986/go-gorilla-api/data"
	"github.com/Lewiscowles1986/go-gorilla-api/repositories"
	"github.com/Lewiscowles1986/go-gorilla-api/rest"
	"github.com/Lewiscowles1986/go-gorilla-api/settings"
)

var a App

func clearTable() {
	a.DB.Exec("DELETE FROM products")
	a.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

func TestMain(m *testing.M) {
	a = App{}

	dbType := settings.Getenv("TEST_DB_TYPE", "sqlite3")
	dbConnStr := settings.GetDBConnStr(
		dbType,
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		settings.Getenv("TEST_DB_NAME", ":memory:"))
	a.Initialize(dbType, dbConnStr)

	fmt.Printf("Testing with settings\n%s: %s\n\n", dbType, dbConnStr)

	code := m.Run()

	clearTable()

	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products", nil)
	response := executeRequest(req)

	total := repositories.GetProductCount(a.DB)
	products := make([]data.Product, 0)
	blankListing := rest.ListingJSONResponse("/products", 0, total, 10,
		rest.ProductsToEntries(products))

	expected, _ := data.JSONMarshal(blankListing)

	checkResponseCode(t, http.StatusOK, response.Code)

	if result := response.Body.String(); result != string(expected) {
		t.Errorf("Expected %s. Got %s", string(expected), result)
	}
}

func TestEmptyTableLargeCount(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/products?count=255", nil)
	response := executeRequest(req)

	total := repositories.GetProductCount(a.DB)
	products := make([]data.Product, 0)
	blankListing := rest.ListingJSONResponse("/products", 0, total, 250,
		rest.ProductsToEntries(products))

	expected, _ := data.JSONMarshal(blankListing)

	checkResponseCode(t, http.StatusOK, response.Code)

	if result := response.Body.String(); result != string(expected) {
		t.Errorf("Expected %s. Got %s", string(expected), result)
	}
}

func TestCreateProduct(t *testing.T) {
	clearTable()

	payload := []byte(`{"name":"test product","price":11.22}`)

	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}
}

func TestCreateProductFailsWithGarbage(t *testing.T) {
	clearTable()

	payload := []byte(`"garbage"`)

	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)
}

func TestGetProduct(t *testing.T) {
	clearTable()

	p := data.CreateProduct("test", 9.99)
	repositories.CreateProduct(a.DB, p)

	req, _ := http.NewRequest("GET", fmt.Sprintf("/product/%s", p.GetID()), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)
}

func TestGetNonExistentProduct(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", fmt.Sprintf("/product/%s",
		uuid.Must(uuid.NewV4(), nil).String()), nil)
	response := executeRequest(req)

	checkResponseCode(t, http.StatusNotFound, response.Code)

	var m map[string]string
	json.Unmarshal(response.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestUpdateProduct(t *testing.T) {
	clearTable()

	p := data.CreateProduct("cheap trash", .99)
	err := repositories.CreateProduct(a.DB, p)
	if err != nil {
		t.Errorf("Unable to save initial model to database")
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("/product/%s", p.GetID()), nil)
	response := executeRequest(req)

	originalProduct, _ := data.ParseProductDataJSON(response.Body.Bytes())

	payload := []byte(`{"name":"test product - updated name","price":11.22}`)

	req, _ = http.NewRequest(
		"PUT", fmt.Sprintf("/product/%s", p.GetID()), bytes.NewBuffer(payload))
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	modifiedProduct, _ := data.ParseProductDataJSON(response.Body.Bytes())

	if modifiedProduct.GetID() != originalProduct.GetID() {
		t.Errorf("Expected the id to remain the same (%v). Got %v",
			originalProduct.GetID(), modifiedProduct.GetID())
	}

	if originalProduct.GetName() == modifiedProduct.GetName() {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'",
			originalProduct.GetName(), modifiedProduct.GetName(),
			modifiedProduct.GetName())
	}

	if originalProduct.GetPrice() == modifiedProduct.GetPrice() {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'",
			originalProduct.GetPrice(), modifiedProduct.GetPrice(),
			modifiedProduct.GetPrice())
	}
}

func TestUpdateProductWithGarbageFails(t *testing.T) {
	clearTable()

	p := data.CreateProduct("cheap trash", .99)
	err := repositories.CreateProduct(a.DB, p)
	if err != nil {
		t.Errorf("Unable to save initial model to database")
	}

	payload := []byte(`gibberish`)

	req, _ := http.NewRequest(
		"PUT", fmt.Sprintf("/product/%s", p.GetID()), bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, http.StatusBadRequest, response.Code)

}

func TestDeleteProductNonExist(t *testing.T) {
	clearTable()

	id := "49458d94-3347-4c3b-a12f-91f2b33fa3ad"
	req, _ := http.NewRequest("DELETE", fmt.Sprintf("/product/%s", id), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func TestDeleteProductExist(t *testing.T) {
	clearTable()

	p := data.CreateProduct("something we're ashamed of", 500000.00)
	err := repositories.CreateProduct(a.DB, p)
	if err != nil {
		t.Errorf("Unable to save initial model to database")
	}

	req, _ := http.NewRequest("GET", fmt.Sprintf("/product/%s", p.GetID()), nil)
	response := executeRequest(req)
	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", fmt.Sprintf("/product/%s", p.GetID()), nil)
	response = executeRequest(req)

	checkResponseCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", fmt.Sprintf("/product/%s", p.GetID()), nil)
	response = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, response.Code)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)

	return rr
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}
