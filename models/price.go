package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type PriceHistory struct {
	Date  time.Time `bson:"priceDate"`
	Value float64   `bson:"priceValue"`
}

type Price struct {
	ID           primitive.ObjectID `bson:"_id,omitempty"`
	ProductID    primitive.ObjectID `bson:"productID"`
	CurrentPrice float64            `bson:"currentPrice"`
	LastPrice    float64            `bson:"lastPrice"`
	PriceHistory []PriceHistory     `bson:"history"`
	Date         time.Time          `bson:"date"`
	High         float64            `bson:"high"`
	Low          float64            `bson:"low"`
}
