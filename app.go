package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"./data"
	"./repositories"
	"./rest"
)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

func (a *App) Initialize(connType, connectionString string) {
	var err error
	a.DB, err = sql.Open(connType, connectionString)
	if err != nil {
		log.Fatal(err)
	}

	a.Router = mux.NewRouter()
	a.initializeDB()
	a.initializeRoutes()
}

func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	uuid4Regex := "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}"

	productSpecificRoute := fmt.Sprintf("/product/{id:%s}", uuid4Regex)
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc(productSpecificRoute, a.getProduct).Methods("GET")
	a.Router.HandleFunc(productSpecificRoute, a.updateProduct).Methods("PUT")
	a.Router.HandleFunc(productSpecificRoute, a.deleteProduct).Methods("DELETE")
}

func (a *App) initializeDB() {
	_, err := a.DB.Exec(`CREATE TABLE IF NOT EXISTS products (
        id VARCHAR(36) NOT NULL,
        name TEXT NOT NULL,
        price NUMERIC(10,2) NOT NULL DEFAULT 0.00,
        CONSTRAINT products_pkey PRIMARY KEY (id)
    )`)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {
	count, page := getPagingFromRequest(r)

	products, err := repositories.GetProducts(a.DB, page, count)
	if err != nil {
		rest.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	total := repositories.GetProductCount(a.DB)
	l := rest.ListingJSONResponse("/products", page, total, count,
		rest.ProductsToEntries(products))
	rest.RespondWithJSON(w, http.StatusOK, l)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.RespondWithError(w, http.StatusBadRequest, "Unable to read request body")
	}
	defer r.Body.Close()

	p, err := data.ParseProductDataJSON(body)
	if err != nil {
		rest.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = repositories.CreateProduct(a.DB, p)
	if err != nil {
		rest.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rest.RespondWithJSON(w, http.StatusCreated, p)
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := data.ParseUUID(vars["id"])

	p, err := repositories.GetProduct(a.DB, id.String())
	if err != nil {
		switch err {
		case sql.ErrNoRows:
			rest.RespondWithError(w, http.StatusNotFound, "Product not found")
		default:
			rest.RespondWithError(w, http.StatusInternalServerError, "Error loading")
		}
		return
	}

	rest.RespondWithJSON(w, http.StatusOK, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := data.ParseUUID(vars["id"])

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rest.RespondWithError(w, http.StatusBadRequest, "Unable to read request body")
		return
	}
	defer r.Body.Close()

	p, err := data.ParseProductDataJSON(body)
	if err != nil {
		rest.RespondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}

	err = repositories.UpdateProduct(a.DB, id.String(), p)
	if err != nil {
		rest.RespondWithError(w, http.StatusInternalServerError, fmt.Sprintf(
			"Unable to save product '%s' with data %+v", id.String(), p))
		return
	}
	m, _ := repositories.GetProduct(a.DB, id.String())

	rest.RespondWithJSON(w, http.StatusOK, m)
}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := data.ParseUUID(vars["id"])

	_, err := repositories.GetProduct(a.DB, id.String())
	if err != nil {
		rest.RespondWithError(w, http.StatusNotFound, fmt.Sprintf(
			"Product '%s' not found", id.String()))
		return
	}

	err = repositories.DeleteProduct(a.DB, id.String())
	if err != nil {
		rest.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	rest.RespondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
}

func getURLQueryParam(r *http.Request, key string) string {
	return url.QueryEscape(r.URL.Query().Get(key))
}

func getPagingFromRequest(r *http.Request) (uint8, uint64) {
	count, _ := strconv.ParseUint(getURLQueryParam(r, "count"), 10, 8)
	page, _ := strconv.ParseUint(getURLQueryParam(r, "page"), 10, 64)

	if count < 1 {
		count = 10
	}
	if count > 250 {
		count = 250
	}
	c := uint8(count)

	return c, page
}
