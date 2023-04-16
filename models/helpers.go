package models

import (
	"fmt"
	"net/url"
	"reflect"
	"strings"

	"github.com/chuxorg/chux-datastore/db"
	"github.com/chuxorg/chux-models/config"
	"github.com/chuxorg/chux-models/errors"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

// config for models package
var _cfg *config.BizObjConfig
var mongoDB *db.MongoDB

func ExtractCompanyName(urlString string) (string, error) {
	parsedURL, err := url.Parse(urlString)
	if err != nil {
		return "", errors.NewChuxModelsError("ExtractCompanyName() Unable to parse the url", err)
	}

	host := parsedURL.Host
	// Split the host into parts
	parts := strings.Split(host, ".")

	// If there are at least two parts (subdomain(s) and domain)
	if len(parts) >= 2 {
		// Return the second last part, which is the domain without the extension
		return parts[len(parts)-2], nil
	}
	msg := fmt.Sprintf("Could not extract company name from url: %s", urlString)
	return "", errors.NewChuxModelsError(msg, nil)
}

// Categorizes all products which are not already categorized
func Categorize(cfg config.BizObjConfig) error {
	// - Get all products that are not categorized
	prd := NewProduct(
		NewProductWithBizObjConfig(*_cfg),
	)
	products, err := mongoDB.Query(prd, "isCatagorized", false)

	if err != nil {
		return errors.NewChuxModelsError("Product.GetUncategorized() Error querying database", err)
	}

	for _, product := range products {
		// -- Iterate over the product's breadcrumbs and create categories
		createdCategories := make([]*Category, len(product.(*Product).Breadcrumbs))

		for index, breadcrumb := range product.(*Product).Breadcrumbs {
			// -- Create a category document
			category := NewCategory(
				WithBizObjConfig(*_cfg),
			)
			category.ProductID = product.(*Product).ID
			category.Name = breadcrumb.Name
			category.Index = index
			category.ParentID = primitive.NewObjectID()

			err := category.Save()
			if err != nil {
				return errors.NewChuxModelsError("Product.Categorize() Error saving category", err)
			}

			createdCategories[index] = category
		}

		/*
			After all categories are created for a product, iterate over the created categories and set the ParentID accordingly.
			The ParentID of the first category in the list (index 0) will remain nil.
			This will help with the tree structure of the categories.
		*/
		for index, category := range createdCategories {
			if index > 0 {
				category.ParentID = createdCategories[index-1].ID
				err := category.Save()
				if err != nil {
					return errors.NewChuxModelsError("Product.Categorize() Error updating category ParentID", err)
				}
			}
		}
	}
	return nil
}

// CompareProducts takes two Product structs and compares their fields to see if anything has changed.
// Returns a map containing the field names as keys and a tuple of the old and new values as the corresponding values.
func CompareProducts(oldProduct, newProduct Product) (map[string][2]interface{}, error) {
	changes := make(map[string][2]interface{})

	v1 := reflect.ValueOf(oldProduct)
	v2 := reflect.ValueOf(newProduct)

	// Loop through the fields of the Product struct
	for i := 0; i < v1.NumField(); i++ {
		field1 := v1.Field(i)
		field2 := v2.Field(i)

		// Ignore unexported fields
		if field1.CanInterface() && field2.CanInterface() {
			// Compare field values
			if !reflect.DeepEqual(field1.Interface(), field2.Interface()) {
				fieldName := v1.Type().Field(i).Name
				changes[fieldName] = [2]interface{}{field1.Interface(), field2.Interface()}
			}
		}
	}

	return changes, nil
}
