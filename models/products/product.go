package products

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/chuxorg/chux-datastore/db"
	"github.com/chuxorg/chux-models/config"
	"github.com/chuxorg/chux-models/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type GTIN struct {
	Type  string `bson:"type,omitempty" json:"type"`
	Value string `bson:"value,omitempty" json:"value"`
}

type Product struct {
	ID                   primitive.ObjectID          `bson:"_id,omitempty" json:"_id,omitempty"`
	URL                  string                      `bson:"url" json:"url"`
	CanonicalURL         string                      `bson:"canonicalUrl" json:"canonicalUrl"`
	Probability          float64                     `bson:"probability" json:"probability"`
	Name                 string                      `bson:"name" json:"name"`
	Offers               []models.Offer              `bson:"offers" json:"offers"`
	SKU                  string                      `bson:"sku" json:"sku"`
	MPN                  string                      `bson:"mpn,omitempty" json:"mpn,omitempty"`
	Brand                string                      `bson:"brand,omitempty" json:"brand,omitempty"`
	Breadcrumbs          []models.Breadcrumb         `bson:"breadcrumbs" json:"breadcrumbs"`
	MainImage            string                      `bson:"mainImage" json:"mainImage"`
	Images               []string                    `bson:"images" json:"images"`
	Description          string                      `bson:"description" json:"description"`
	DescriptionHTML      string                      `bson:"descriptionHtml" json:"descriptionHtml"`
	AdditionalProperties []models.AdditionalProperty `bson:"additionalProperty" json:"additionalProperty"`
	AggregateRating      models.AggregateRating      `bson:"aggregateRating" json:"aggregateRating"`
	GTINs                []GTIN                      `bson:"gtins,omitempty" json:"gtin,omitempty"`
	Color                string                      `bson:"color,omitempty" json:"color,omitempty"`
	Style                string                      `bson:"style,omitempty" json:"style,omitempty"`
	DateCreated          models.CustomTime           `bson:"dateCreated,omitempty" json:"dateCreated,omitempty"`
	DateModified         models.CustomTime           `bson:"dateModified,omitempty" json:"dateModified,omitempty"`
	isNew                bool                        `bson:"isNew,omitempty" json:"isNew,omitempty"`
	isDeleted            bool                        `bson:"isDeleted,omitempty" json:"isDeleted,omitempty"`
	isDirty              bool                        `bson:"isDirty,omitempty" json:"isDirty,omitempty"`
	originalState        *Product                    `bson:"-" json:"-"`
}

var _cfg *config.BizObjConfig
var mongoDB *db.MongoDB

func New(options ...func(*Product)) *Product {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	_cfg = config.New()
	product := &Product{}
	for _, option := range options {
		option(product)
	}

	mongoDB = db.New(
		db.WithURI(product.GetURI()),
		db.WithDatabaseName(product.GetDatabaseName()),
		db.WithCollectionName(product.GetCollectionName()),
		db.WithTimeout(float64(_cfg.DataStores.DataStoreMap["mongo"].Timeout)),
	)

	product.isNew = true
	product.isDeleted = false
	product.isDirty = false
	return product
}

func WithLoggingLevel(level string) func(*Product) {
	return func(product *Product) {
		_cfg.Logging.Level = level
	}
}

func WithBizObjConfig(config config.BizObjConfig) func(*Product) {
	return func(product *Product) {
		_cfg = &config
	}
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
		return nil, err
	}

	return docs, nil
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

// Saves the Model to a Data Store
func (p *Product) Save() error {

	if p.isNew {
		var exists bool
		changes := make(map[string][2]interface{})
		products, err := p.Exists()
		if err != nil {
			return err // return the error
		}

		if len(products) > 0 {
			product := products[0].(*Product)
			var err error
			changes, err = CompareProducts(*product, *p)
			if err != nil {
				return err
			}
			if len(changes) > 0 {
				fmt.Println(changes)
			}
		}
		exists = len(products) > 0 && len(changes) > 0
		if !exists {
			//-- This check that the product does not exist in the database.
			//-- it is required because new products that come in from a parse
			//-- run will bring in dupes, which we don't want to save.

			//-- Create a new document
			err := mongoDB.Create(p)
			if err != nil {
				return err
			}
		} else {
			p.isNew = false
			p.isDirty = true
			//-- this will cause an update to the document
			return p.Save()
		}

	} else if p.IsDirty() && !p.isDeleted {
		// Ensure the ID is a valid hex string representation of an ObjectID
		_, err := primitive.ObjectIDFromHex(p.ID.Hex())
		if err != nil {
			return fmt.Errorf("invalid ObjectID: %v", err)
		}
		//--update this document
		err = mongoDB.Update(p, p.ID.Hex())
		if err != nil {
			return err
		}
	} else if p.isDeleted && !p.isNew {
		//--delete the document
		err := mongoDB.Delete(p, p.ID.Hex())
		if err != nil {
			return err
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
			return fmt.Errorf("unable to set internal state")
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
