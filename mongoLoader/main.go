package main

import (
	"fmt"
	"github.com/gocarina/gocsv"
	"gopkg.in/mgo.v2"
	"os"
	"time"
)

// helper for input string date format
type DateTime struct {
	time.Time
}

// fundamental business struct
type Business struct {
	Id        int      `json:"id" bson:"_id" csv:"id"`
	Uuid      string   `json:"uuid" bson:"uuid" csv:"uuid"`
	Name      string   `json:"name" bson:"name" csv:"name"`
	Address   string   `json:"address" bson:"address" csv:"address"`
	Address2  string   `json:"address2" bson:"address2" csv:"address2"`
	City      string   `json:"city" bson:"city" csv:"city"`
	State     string   `json:"state" bson:"state" csv:"state"`
	Zip       string   `json:"zip" bson:"zip" csv:"zip"`
	Country   string   `json:"country" bson:"country" csv:"country"`
	Phone     string   `json:"phone" bson:"phone" csv:"phone"`
	Website   string   `json:"website" bson:"website" csv:"website"`
	CreatedAt DateTime `json:"created_at" bson:"created_at" csv:"created_at"`
}

// loads the given CSV file into a MongoDB instance
func main() {
	businessFile, err := os.Open("engineering_project_businesses.csv")
	if err != nil {
		panic(err)
	}
	defer businessFile.Close()

	businesses := []*Business{}

	if err := gocsv.UnmarshalFile(businessFile, &businesses); err != nil { // Load businesses from file
		panic(err)
	}

	session, err := mgo.Dial("localhost")
	if err != nil {
		panic(err)
	}
	defer session.Close()

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	c := session.DB("OwnLocal").C("businesses")

	cnt := 0
	for _, business := range businesses {
		c.Insert(business)
		fmt.Println("Inserting...", business.Name)
		cnt++
	}
	fmt.Println("Inserted %d records", cnt)
}

// Convert CSV string to time type
func (date *DateTime) UnmarshalCSV(csv string) (err error) {
	date.Time, err = time.Parse("2006-01-02 15:04:05", csv)
	if err != nil {
		return err
	}
	return nil
}
