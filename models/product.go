package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chuxorg/chux-datastore/db"
	"github.com/chuxorg/chux-models/errors"
	"github.com/chuxorg/chux-models/logging"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

var mongoDB *db.MongoDB

type GTIN struct {
	Type  string `bson:"type,omitempty" json:"type"`
	Value string `bson:"value,omitempty" json:"value"`
}

type Product struct {
	ID                   primitive.ObjectID   `bson:"_id,omitempty" json:"_id,omitempty"`
	URL                  string               `bson:"url" json:"url"`
	CanonicalURL         string               `bson:"canonicalUrl" json:"canonicalUrl"`
	CompanyName          string               `bson:"companyName" json:"companyName"`
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
	CategoryID           primitive.ObjectID   `bson:"categoryId" json:"categoryId"`
	IsCategorized        bool                 `bson:"isCategorized" json:"isCategorized"`
	ImagesProcessed      bool                 `bson:"imagesProcessed" json:"imagesProcessed"`
	FilesProcessed       bool                 `bson:"filesProcessed" json:"filesProcessed"`
	originalState        *Product             `bson:"-" json:"-"`
}

func NewProduct() *Product {

	logging.Debug("NewProduct was called")
	p := &Product{}

	mongoDB = db.New(
		db.WithURI(p.GetURI()),
		db.WithDatabaseName(p.GetDatabaseName()),
		db.WithCollectionName(p.GetCollectionName()),
		db.WithTimeout(30),
	)

	p.isNew = true
	p.isDeleted = false
	p.isDirty = false
	return p
}

func (p *Product) GetCollectionName() string {
	logging.Debug("Product.GetCollectionName() was called")
	return "products"
}

func (p *Product) GetDatabaseName() string {
	logging.Debug("Product.GetDatabaseName() was called")
	return os.Getenv("MONGO_DATABASE")
}

func (p *Product) GetURI() string {
	logging.Debug("Product.GetURI() was called")
	username := os.Getenv("MONGO_USER_NAME")
	password := os.Getenv("MONGO_PASSWORD")

	uri := os.Getenv("MONGO_URI")
	mongoURI := fmt.Sprintf(uri, username, password)
	masked := fmt.Sprintf(uri, "********", "********")
	logging.Info("Mongo URI: %s", masked)
	return mongoURI
}

func (p *Product) GetID() primitive.ObjectID {
	logging.Debug("Product.GetID() was called")
	return p.ID
}
func (p *Product) SetID(id primitive.ObjectID) {
	logging.Debug("Product.SetID() was called")
	p.ID = id
}

// If the Model has changes, will return true
func (p *Product) IsDirty() bool {
	logging.Debug("Product.IsDirty() was called")
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
	logging.Info("Product.IsDirty() isDirty: %t", p.isDirty)
	return p.isDirty
}

// When the Model is first created,
// the model is considered New. After the model is
// Saved or Loaded it is no longer New
func (p *Product) IsNew() bool {
	logging.Debug("Product.IsNew() was called")
	return p.isNew
}

// Saves the Model to a Data Store
func (p *Product) Save() error {

	logging.Debug("Product.Save() was called")
	if p.isNew {
		logging.Debug("Product.Save() Product is new")
		companyName, err := ExtractCompanyName(p.CanonicalURL)
		logging.Info("Product.Save() Extracted Company Name: %s", companyName)
		if err != nil {
			logging.Error("Product.Save() Error extracting Product.CompanyName: %s", err.Error())
			return errors.NewChuxModelsError("Product.Save() Error extracting Product.CompanyName", err)
		}
		p.CompanyName = companyName
		// -- Set the date created to now
		p.DateCreated.Now()
		// -- Set the category to uncategorized
		p.IsCategorized = false
		// Mark product as processed
		p.FilesProcessed = true
		p.ImagesProcessed = false
		//-- Upsert document
		err = mongoDB.Upsert(p)
		if err != nil {
			logging.Error("Product.Save() Error creating/updating Product in MongoDB: %s", err.Error())
			return errors.NewChuxModelsError("Product.Save() Error creating/updating Product in MongoDB", err)
		}

	} else if p.IsDirty() && !p.isDeleted {
		logging.Debug("Product.Save() Product is dirty")
		// Ensure the ID is a valid hex string representation of an ObjectID
		_, err := primitive.ObjectIDFromHex(p.ID.Hex())
		if err != nil {
			logging.Error("Product.Save() invalid ObjectID: %s", err.Error())
			return errors.NewChuxModelsError("Product.Save() invalid ObjectID", err)
		}
		// -- Set the date modified to now
		p.DateModified.Now()
		//--update this document
		err = mongoDB.Update(p, p.ID.Hex())
		if err != nil {
			logging.Error("Product.Save() Error updating Product in MongoDB: %s", err.Error())
			return errors.NewChuxModelsError("Product.Save() Error updating Product in MongoDB", err)
		}
	} else if p.isDeleted && !p.isNew {
		logging.Info("Product.Save() Product is deleted")
		//--delete the document
		err := mongoDB.Delete(p, p.ID.Hex())
		logging.Info("Product.Save() Product was deleted")
		if err != nil {
			logging.Error("Product.Save() Error deleting Product in MongoDB: %s", err.Error())
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
			logging.Error("Product.Save() Error serializing Product: %s", err.Error())
			return errors.NewChuxModelsError("Product.Save() Error serializing Product.", err)
		}
		p.SetState(serialized)
	}

	logging.Info("Product.Save() Product saved successfully")
	return nil
}

// Loads a Model from MongoDB by id
func (p *Product) Load(id string) (interface{}, error) {
	logging.Debug("Product.Load() Product was called")

	retVal, err := mongoDB.GetByID(p, id)
	if err != nil {
		logging.Error("Product.Load() Error loading Product from MongoDB: %s", err.Error())
		return nil, errors.NewChuxModelsError("Product.Load() Error loading Product from MongoDB", err)
	}
	product, ok := retVal.(*Product)
	if !ok {
		logging.Error("Product.Load() unable to cast retVal to *Product")
		return nil, errors.NewChuxModelsError("Product.Load() unable to cast retVal to *Product", err)
	}
	serialized, err := product.Serialize()
	if err != nil {
		logging.Error("Product.Load() Error serializing Product: %s", err.Error())
		return nil, errors.NewChuxModelsError("Product.Load() Error serializing Product", err)
	}
	p.SetState(serialized)
	p.isNew = false
	p.isDirty = false
	p.isDeleted = false
	logging.Info("Product.Load() Product loaded successfully")
	return retVal, nil
}

func (p *Product) Query(args ...interface{}) ([]db.IMongoDocument, error) {
	logging.Debug("Product.Query() was called")

	results, err := mongoDB.Query(p, args...)
	if err != nil {
		logging.Error("Product.Query() Error occurred querying Products: %s", err.Error())
		return nil, errors.NewChuxModelsError("Product.Query() Error occurred querying Products", err)
	}
	logging.Info("Product.Query() Products queried successfully")
	return results, nil
}

func (p *Product) GetAll() ([]db.IMongoDocument, error) {

	logging.Debug("Product.GetAll() was called")

	mongoDB := &db.MongoDB{}
	products, err := mongoDB.GetAll(p)
	logging.Info()
	if err != nil {
		logging.Error("Product.GetAll() Error occurred getting all Products: %s", err.Error())
		return nil, errors.NewChuxModelsError("Product.GetAll() Error occurred getting all Products", err)
	}

	logging.Info("Product.GetAll() Products retrieved successfully")
	return products, nil
}

// Marks a Model for deletion from the Data Store
// when Save() is called, the Model will be deleted
func (p *Product) Delete() error {
	logging.Debug("Product.Delete() was called")
	p.isDeleted = true
	return nil
}

// Sets the internal state of the model.
func (p *Product) SetState(json string) error {
	logging.Debug("Product.SetState() was called")
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
	logging.Debug("Product.Parse() was called")
	err := p.SetState(json)
	if err != nil {
		logging.Error("Product.Parse() error setting state")
		return errors.NewChuxModelsError("Product.Parse() Error setting state", err)
	}
	p.isNew = true // this is a new model
	return nil
}

func (p *Product) Search(args ...interface{}) ([]interface{}, error) {
	logging.Debug("Product.Search() was called")
	return nil, nil
}

func (p *Product) Serialize() (string, error) {
	logging.Debug("Product.Serialize() was called")
	bytes, err := json.Marshal(p)
	if err != nil {
		logging.Error("Product.Serialize() error ocurred ", err)
		return "", errors.NewChuxModelsError("Product.Serialize() error occured", err)
	}
	return string(bytes), nil
}

func (p *Product) Deserialize(jsonData []byte) error {
	logging.Debug("Product.Deserialize() was called")
	err := json.Unmarshal(jsonData, p)
	if err != nil {
		logging.Error("Product.Deserialize() error occurred ", err)
		return errors.NewChuxModelsError("Product.Deserialize() error occured", err)
	}
	return nil
}
