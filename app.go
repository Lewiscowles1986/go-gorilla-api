package main

import (
	"context"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/pprof"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"

	"github.com/Lewiscowles1986/go-gorilla-api/data"
	"github.com/Lewiscowles1986/go-gorilla-api/repositories"
	"github.com/Lewiscowles1986/go-gorilla-api/rest"
)

// App - Structure for Global State
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize - Setup App resources
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

// Run - Main Loop
func (a *App) Run(addr string, wait time.Duration) {
	srv := &http.Server{
		Addr: addr,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      maxClientsMiddleware(a.Router, 5), // Pass our instance of gorilla/mux in.
	}

	// Run our server in a goroutine so that it doesn't block.
	go func() {
		log.Fatal(srv.ListenAndServe())
	}()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()

	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.

	srv.Shutdown(ctx)
	a.DB.Close()

	log.Println("shutting down")
	os.Exit(0)
}

func (a *App) initializeRoutes() {
	uuid4Regex := "[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[89abAB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}"

	productSpecificRoute := fmt.Sprintf("/product/{id:%s}", uuid4Regex)
	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/product", a.createProduct).Methods("POST")
	a.Router.HandleFunc(productSpecificRoute, a.getProduct).Methods("GET")
	a.Router.HandleFunc(productSpecificRoute, a.updateProduct).Methods("PUT")
	a.Router.HandleFunc(productSpecificRoute, a.deleteProduct).Methods("DELETE")

	a.Router.HandleFunc("/debug/pprof/", pprof.Index)
	a.Router.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	a.Router.HandleFunc("/debug/pprof/profile", pprof.Profile)
	a.Router.HandleFunc("/debug/pprof/symbol", pprof.Symbol)

	// Manually add support for paths linked to by index page at /debug/pprof/
	a.Router.Handle("/debug/pprof/goroutine", pprof.Handler("goroutine"))
	a.Router.Handle("/debug/pprof/heap", pprof.Handler("heap"))
	a.Router.Handle("/debug/pprof/threadcreate", pprof.Handler("threadcreate"))
	a.Router.Handle("/debug/pprof/block", pprof.Handler("block"))
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

	_, getErr := repositories.GetProduct(a.DB, id.String())
	if getErr != nil {
		rest.RespondWithError(w, http.StatusNotFound, fmt.Sprintf(
			"Product '%s' not found", id.String()))
		return
	}

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

func maxClientsMiddleware(handler http.Handler, maxClients uint16) http.Handler {
	log.Printf("Max Clients Middleware %d", maxClients)
	sema := make(chan struct{}, int(maxClients))
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		timeOut := time.NewTimer(time.Second * 1)
		select {
		case sema <- struct{}{}:
			handler.ServeHTTP(w, r)
			defer func() { <-sema }()
		case <-timeOut.C:
			rest.RespondWithError(w, http.StatusServiceUnavailable, "Too many requests")
		}
	})
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
	if page < 1 {
		page = 1
	}
	c := uint8(count)

	return c, page
}
