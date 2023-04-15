package models

import (
	"context"
	"fmt"
	"net/url"
	"strings"

	"github.com/chuxorg/chux-models/errors"
	"github.com/chuxorg/chux-models/models/categories"
	"github.com/chuxorg/chux-models/models/products"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// ...

func createCategoriesFromProducts(ctx context.Context, client *mongo.Client) error {
	// Get the products and categories collections
	productsColl := client.Database("chux-cprs").Collection("products")
	categoriesColl := client.Database("chux-cprs").Collection("categories")

	// Ensure the "name" field in the categories collection is unique and indexed
	nameIndexModel := mongo.IndexModel{
		Keys:    bson.M{"name": 1},
		Options: options.Index().SetUnique(true).SetCollation(&options.Collation{Locale: "en", Strength: 1}),
	}
	_, err := categoriesColl.Indexes().CreateOne(ctx, nameIndexModel)
	if err != nil {
		return fmt.Errorf("failed to create index: %v", err)
	}

	// Find uncategorized products
	filter := bson.M{"isCategorized": false}
	cursor, err := productsColl.Find(ctx, filter)
	if err != nil {
		return fmt.Errorf("failed to find uncategorized products: %v", err)
	}
	defer cursor.Close(ctx)

	// Iterate over the uncategorized products
	for cursor.Next(ctx) {
		
		var product products.Product
		err = cursor.Decode(&product)
		if err != nil {
			return fmt.Errorf("failed to decode product: %v", err)
		}

		// Iterate over the product's breadcrumbs and create categories
		for index, breadcrumb := range product.Breadcrumbs {
			// Create a category document
			category := categories.Category{
				ProductID: product.ID,
				Name:      strings.ToLower(breadcrumb.Name),
				Index:     index,
			}

			// Insert the category document into the categories collection
			_, err := categoriesColl.InsertOne(ctx, category)
			if err != nil {
				return fmt.Errorf("failed to insert category: %v", err)
			}
		}

		// Update the product's isCategorized field
		updateFilter := bson.M{"_id": product.ID}
		update := bson.M{"$set": bson.M{"isCategorized": true}}
		_, err = productsColl.UpdateOne(ctx, updateFilter, update)
		if err != nil {
			return fmt.Errorf("failed to update product: %v", err)
		}
	}

	return nil
}


func ExtractCompanyName(urlStr string) (string, error) {
	parsedURL, err := url.Parse(urlStr)
	if err != nil {
		return "", err
	}

	host := parsedURL.Host
	// Split the host into parts
	parts := strings.Split(host, ".")

	// If there are at least two parts (subdomain(s) and domain)
	if len(parts) >= 2 {
		// Return the second last part, which is the domain without the extension
		return parts[len(parts)-2], nil
	}
	msg := fmt.Sprintf("Could not extract company name from url: %s", urlStr)
	return "", errors.NewChuxModelsError(msg, nil)
}
