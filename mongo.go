package main

import (
	"fmt"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"log"
	"errors"
)

type DataStore struct {
	session *mgo.Session
}
var ds DataStore

const (
	dbUrl = "localhost"
	dbName = "OwnLocal"
	collection = "businesses"
	linksUrl = "http://%s/businesses?page=%d&size=%d"
)

// Connects to MongoDB
func init() {
	session, err := mgo.Dial(dbUrl)
	if err != nil {
		panic(err)
	}
	defer session.Close()

	session.SetMode(mgo.Monotonic, true)

	ds = DataStore{session.Copy()}
}

// keeps session from getting garbage collected
func dataStore() *DataStore {
	return &DataStore{ds.session.Copy()}
}
// helper to access collection
func c() *mgo.Collection {
	return dataStore().session.DB(dbName).C(collection)
}
func a() *mgo.Collection {
	return dataStore().session.DB(dbName).C("Authorization")
}

// structure and routines for constructing links
type RequestAll struct {
	Url string
	Page int
	Size int
	numBusinesses int
}

func (r RequestAll) lastCount() int {
	if (r.numBusinesses % r.Size == 0) {
		return (r.numBusinesses / r.Size) + 1
	} else {
		return r.numBusinesses / r.Size
	}
}
func (r RequestAll) firstUrl() string {
	return fmt.Sprintf(linksUrl, r.Url, 1, r.Size)
}
func (r RequestAll) lastUrl() string {
	return fmt.Sprintf(linksUrl, r.Url, r.lastCount() + 1, r.numBusinesses - (r.lastCount() * r.Size))
}

func (r RequestAll) prevUrl() string {
	if (r.Page > 1) {
		return fmt.Sprintf(linksUrl, r.Url, r.Page - 1, r.Size)
	}
	return ""
}

func (r RequestAll) nextUrl() string {
	if (r.Page * r.Size < r.numBusinesses) {
		return fmt.Sprintf(linksUrl, r.Url, r.Page + 1, r.Size)
	}
	return ""
}

// fetches collection of businesses
func FindBusinesses(request RequestAll)(*BusinessesPage, error) {
	var err error
	request.numBusinesses, err = c().Count()
	if err != nil {
		log.Fatal(err)
	}
	meta := Meta{request.lastCount() + 1}

	links := Links{
		request.firstUrl(),
		request.lastUrl(),
		request.prevUrl(),
		request.nextUrl(),
	}

	gte := (request.Page * request.Size) - request.Size
	lt := request.Page * request.Size

	// additional bounds checking on the page size and start page
	if (gte > request.numBusinesses) {
		return nil, errors.New("Out of Range")
	}
	// check for overflow, to reduce number of returned businesses if needed
	if (lt > request.numBusinesses) {
		lt = request.numBusinesses
	}
	businesses := make([]Business, lt - gte) // allocates memory for businesses
	page := BusinessesPage{meta, businesses, links }
	// stores result and optionally the error
	err = c().Find(bson.M{"_id": bson.M{"$gte": gte, "$lt": lt}}).Sort("_id").All(&businesses)
	if err != nil {
		log.Fatal(err)
	}

	return &page, err
}

// fetched single business
func FindBusiness(id int) (*Business, error) {
	business := &Business{} // allocates memory for Business
	err := c().Find(bson.M{"_id": id}).One(business) // stores result and optionally the error

	return business, err
}

func TokenExist(id string) bool {
	return (a().Find(bson.M{"token": id}).Count() > 0)

}
