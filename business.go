package main

import (
	"time"
)

// Paginated data structure
type BusinessesPage struct {
	Meta       Meta       `json:"meta"`
	Businesses []Business `json:"businesses"`
	Links      Links      `json:"links"`
}

// Paginated prefix structure
type Meta struct {
	TotalPages int `json:"total-pages"`
}

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

// Paginated links
type Links struct {
	First string `json:"first"`
	Last  string `json:"last"`
	Prev  string `json:"prev"`
	Next  string `json:"next"`
}
