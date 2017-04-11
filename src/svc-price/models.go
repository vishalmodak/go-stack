package main

import (
  "time"
)

type Price struct {
	ItemId			string
	LastUpdate	time.Time
	ListPrice `json:"listPrice"`
	OfferPrice `json:"offerPrice"`
}

type ListPrice struct {
	MinPrice		float32
	MaxPrice		float32
	Price 			float32
	PriceType		string
}

type OfferPrice struct {
	MinPrice		float32
	MaxPrice		float32
	Price 			float32
	PriceType		string
	StartDate		time.Time
	EndDate			time.Time
}
