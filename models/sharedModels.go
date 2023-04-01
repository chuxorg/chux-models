package models

type Offer struct {
	Price        string `bson:"price"`
	Currency     string `bson:"currency"`
	Availability string `bson:"availability"`
}

type Breadcrumb struct {
	Name string `bson:"name"`
	Link string `bson:"link"`
}

type AdditionalProperty struct {
	Name  string `bson:"name"`
	Value string `bson:"value"`
}

type AggregateRating struct {
	RatingValue float64 `bson:"ratingValue"`
	BestRating  float64 `bson:"bestRating"`
	ReviewCount int     `bson:"reviewCount"`
}
