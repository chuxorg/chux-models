package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chuxorg/chux-datastore/db"
	"github.com/chuxorg/chux-models/config"
	"github.com/chuxorg/chux-models/errors"
	"github.com/chuxorg/chux-models/interfaces"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GTIN struct {
	Type  string `bson:"type,omitempty" json:"type"`
	Value string `bson:"value,omitempty" json:"value"`
}

type Product struct {
	ID                   primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	URL                  string               `bson:"url" json:"url"`
	CanonicalURL         string               `bson:"canonicalUrl" json:"canonicalUrl"`
	CompanyName          string               `bson:"companyName, omitempty" json:"companyName"`
	Probability          float64              `bson:"probability" json:"probability"`
	Name                 string               `bson:"name" json:"name"`
	Offers               []Offer              `bson:"offers" json:"offers"`
	SKU                  string               `bson:"sku" json:"sku"`
	MPN                  string               `bson:"mpn,omitempty" json:"mpn,omitempty"`
	Brand                string               `bson:"brand,omitempty" json:"brand,omitempty"`
	Breadcrumbs          []Breadcrumb         `bson:"breadcrumbs" json:"breadcrumbs"`
	MainImage            string               `bson:"mainImage" json:"mainImage"`
	Images               []string             `bson:"images" json:"images"`
	Description          string               `bson:"description" json:"description"`
	DescriptionHTML      string               `bson:"descriptionHtml" json:"descriptionHtml"`
	AdditionalProperties []AdditionalProperty `bson:"additionalProperty" json:"additionalProperty"`
	AggregateRating      AggregateRating      `bson:"aggregateRating" json:"aggregateRating"`
	GTINs                []GTIN               `bson:"gtins,omitempty" json:"gtin,omitempty"`
	Color                string               `bson:"color,omitempty" json:"color,omitempty"`
	Style                string               `bson:"style,omitempty" json:"style,omitempty"`
	DateCreated          CustomTime           `bson:"dateCreated,omitempty" json:"dateCreated,omitempty"`
	DateModified         CustomTime           `bson:"dateModified,omitempty" json:"dateModified,omitempty"`
	isNew                bool                 `bson:"isNew,omitempty" json:"isNew,omitempty"`
	isDeleted            bool                 `bson:"isDeleted,omitempty" json:"isDeleted,omitempty"`
	isDirty              bool                 `bson:"isDirty,omitempty" json:"isDirty,omitempty"`
	IsCategorized        bool                 `bson:"isCategorized" json:"isCategorized"`
	originalState        *Product             `bson:"-" json:"-"`
}

func NewProduct(opts ...func(interfaces.IModel)) *Product {
	p := &Product{}
	p.Apply(opts...)
	return p
}

func (p *Product) Apply(opts ...func(interfaces.IModel)) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	_cfg = config.New()
	
	for _, opt := range opts {
		opt(p)
	}

	mongoDB = db.New(
		db.WithURI(p.GetURI()),
		db.WithDatabaseName(p.GetDatabaseName()),
		db.WithCollectionName(p.GetCollectionName()),
	)

	p.isNew = true
	p.isDeleted = false
	p.isDirty = false
}

func (p *Product) SetLoggingLevel(level string) {
	_cfg.Logging.Level = level
}
func (p *Product) SetBizObjConfig(config config.BizObjConfig) {
	_cfg = &config
}
func (p *Product) SetDataStoresConfig(config config.DataStoresConfig) {
	_cfg.DataStores = config
}


func (p *Product) GetCollectionName() string {
	return "products"
}

func (p *Product) GetDatabaseName() string {
	return _cfg.DataStores.DataStoreMap["mongo"].DatabaseName
}

func (p *Product) GetURI() string {
	return _cfg.DataStores.DataStoreMap["mongo"].URI
}

func (p *Product) GetID() primitive.ObjectID {
	return p.ID
}

// If the Model has changes, will return true
func (p *Product) IsDirty() bool {
	if p.originalState == nil {
		return false
	}

	originalBytes, err := p.originalState.Serialize()
	if err != nil {
		return false
	}

	currentBytes, err := p.Serialize()
	if err != nil {
		return false
	}

	p.isDirty = string(originalBytes) != string(currentBytes)
	return p.isDirty
}

// When the Model is first created,
// the model is considered New. After the model is
// Saved or Loaded it is no longer New
func (p *Product) IsNew() bool {
	return p.isNew
}

func (p *Product) Exists() ([]db.IMongoDocument, error) {

	docs, err := mongoDB.Query(p, "canonicalUrl", p.CanonicalURL, "description", p.Description)
	if err != nil {
		return nil, errors.NewChuxModelsError("Product.Exists() Error querying database", err)
	}

	return docs, nil
}

// Saves the Model to a Data Store
func (p *Product) Save() error {

	if p.isNew {
		var exists bool
		changes := make(map[string][2]interface{})
		products, err := p.Exists()
		if err != nil {
			return errors.NewChuxModelsError("Product.Save() Error checking if product exists", err)
		}

		if len(products) > 0 {
			product := products[0].(*Product)
			var err error
			p.ID = product.ID
			//-- The p.Exists() call above will return a slice of products if there are any
			//-- that match the canonicalUrl and description. We need to compare the
			//-- incoming product with the existing product to see if there are any
			//-- changes. If there are, we need to update the existing product with
			changes, err = CompareProducts(*product, *p)
			if err != nil {
				return errors.NewChuxModelsError("Product.Save() Error comparing Products", err)
			}
			if len(changes) > 0 {
				fmt.Println(changes)
			}
		}
		exists = len(products) > 0
		if !exists {
			//-- This checks that the product does not exist in the database.
			//-- it is required because new products that come in from a parse
			//-- run will bring in dupes, which I don't want to save.

			// -- Set the company name
			p.CompanyName, err = ExtractCompanyName(p.CanonicalURL)
			if err != nil {
				return errors.NewChuxModelsError("Product.Save() Error extracting Product.CompanyName", err)
			}
			// -- Set the date created to now
			p.DateCreated.Now()
			// -- Set the category to uncategorized
			p.IsCategorized = false
			//-- Create a new document
			err := mongoDB.Create(p)
			if err != nil {
				return errors.NewChuxModelsError("Product.Save() Error creating Product in MongoDB", err)
			}
		} else {
			if len(changes) > 0 {
				// The product could exist in Mongo yet was changed on the website on the last crawl
				// If it did, the product will be updated with the new data
				p.isNew = false
				p.isDirty = true
				//-- this will cause an update to the document
				//TODO: Set and Save Price History, Inventory History, and other data
				return p.Save()
			}
		}

	} else if p.IsDirty() && !p.isDeleted {
		// Ensure the ID is a valid hex string representation of an ObjectID
		_, err := primitive.ObjectIDFromHex(p.ID.Hex())
		if err != nil {
			return errors.NewChuxModelsError("Product.Save() invalid ObjectID", err)
		}
		// -- Set the date modified to now
		p.DateModified.Now()
		//--update this document
		err = mongoDB.Update(p, p.ID.Hex())
		if err != nil {
			return errors.NewChuxModelsError("Product.Save() Error updating Product in MongoDB", err)
		}
	} else if p.isDeleted && !p.isNew {
		//--delete the document
		err := mongoDB.Delete(p, p.ID.Hex())
		if err != nil {
			return errors.NewChuxModelsError("Product.Save() Error deleting Product in MongoDB", err)
		}
	}

	// If the Product has been deleted, then this is a new Product
	p.isNew = p.isDeleted
	// little confusing but use the IsDirty() func to set isDirty field on Product struct
	p.isDirty = p.IsDirty()
	p.isDeleted = false

	// serialized will help set the current state
	var serialized string
	var err error
	if p.isNew {
		serialized = ""
		p.originalState = nil
	} else {
		//--reset state
		serialized, err = p.Serialize()
		if err != nil {
			return errors.NewChuxModelsError("Product.Save() Error serializing Product.", err)
		}
		p.SetState(serialized)
	}

	return nil
}

// Loads a Model from MongoDB by id
func (p *Product) Load(id string) (interface{}, error) {
	retVal, err := mongoDB.GetByID(p, id)
	if err != nil {
		return nil, err
	}
	product, ok := retVal.(*Product)
	if !ok {
		return nil, fmt.Errorf("unable to cast retVal to *Product")
	}
	serialized, err := product.Serialize()
	if err != nil {
		return nil, fmt.Errorf("unable to set internal state")
	}
	p.SetState(serialized)
	p.isNew = false
	p.isDirty = false
	p.isDeleted = false

	return retVal, nil
}

func (p *Product) GetAll() ([]db.IMongoDocument, error) {
	mongoDB := &db.MongoDB{}
	products, err := mongoDB.GetAll(p)
	if err != nil {
		return nil, err
	}
	return products, nil
}

// Marks a Model for deletion from the Data Store
// when Save() is called, the Model will be deleted
func (p *Product) Delete() error {
	p.isDeleted = true
	return nil
}

// Sets the internal state of the model.
func (p *Product) SetState(json string) error {
	// Store the current state as the original state
	original := &Product{}
	*original = *p
	p.originalState = original

	// Deserialize the new state
	return p.Deserialize([]byte(json))
}

// Sets the internal state of the model of a new Product
// from a JSON String.
func (p *Product) Parse(json string) error {
	err := p.SetState(json)
	p.isNew = true // this is a new model
	return err
}

func (p *Product) Search(args ...interface{}) ([]interface{}, error) {
	return nil, nil
}

func (p *Product) Serialize() (string, error) {
	bytes, err := json.Marshal(p)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (p *Product) Deserialize(jsonData []byte) error {
	err := json.Unmarshal(jsonData, p)
	if err != nil {
		return err
	}
	return nil
}
