package main

import (
	"net/http"
	"strconv"
	"github.com/gorilla/mux"
	"fmt"
)

func Auth(w http.ResponseWriter, r *http.Request) bool {
	//Get the authorization header.
	_, password, ok := r.BasicAuth()

	if (ok) {
		if (len(password == 0) || !TokenExist(password)) {
			sendJsonError(w, http.StatusForbidden)
			return false
		}
	} else {
		w.Header().Set("WWW-Authenticate", "Basic realm=\"user\"")
		http.Error(w, http.StatusText(401), 401)
		return false
	}
	return true
}

// Send basic usage statement (html encoding?)
func Index(w http.ResponseWriter, r *http.Request) {
	if (!Auth(w, r)) {
		return
	}
	fmt.Fprintln(w, "Welcome to the REST Businesses API!")
	fmt.Fprintln(w, "Usage:")
	fmt.Fprintln(w, `GET /businesses?page={number}&size={number}
Fetch list of businesses with pagination.`)
	fmt.Fprintln(w, "")
	fmt.Fprintln(w, `GET /business/{id}
Fetch a specific business as a JSON object`)
}

// sends JSON message of all/paginated businesses
func BusinessesAll(w http.ResponseWriter, r *http.Request) {
	if (!Auth(w, r)) {
		return
	}
	var page, size int
	var err error
	// get and bounds check the page number
	if page, err = strconv.Atoi(r.FormValue("page")); err != nil || page < 1 {
		page = 1
	}
	// get and bounds check the page size
	if size, err = strconv.Atoi(r.FormValue("size")); err != nil || size < 1 || size > 1000 {
		size = 50
	}

	// Fetch the corresponding businesses from MongoDB
	var businessesPage *BusinessesPage
	if businessesPage, err = FindBusinesses(RequestAll{Url:r.Host, Page:page, Size:size}); err == nil {
		// send thee data
		sendJsonData(w, businessesPage)
		return
	}

	// If we didn't find it, 404 Not Found
	sendJsonError(w, http.StatusNotFound)
}

// sends JSON message of a specific business
func BusinessById(w http.ResponseWriter, r *http.Request) {
        if (!Auth(w, r)) {
		return
	}

	var id int
	var err error

	// capture the business id
	if id, err = strconv.Atoi(mux.Vars(r)["id"]); err != nil {
		// Bad Request 400
		sendJsonError(w, http.StatusBadRequest)
		return
	}

	// find the business object in MongoDB
	var business *Business
	if business, err = FindBusiness(id); err == nil {
		// Send the data
		sendJsonData(w, business)
		return
	}

	// If we didn't find it, 404 Not Found
	sendJsonError(w, http.StatusNotFound)
}
