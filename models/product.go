package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chuxorg/chux-models/config"
	"github.com/chuxorg/chux-datastore/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

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

type Product struct {
	ID                   primitive.ObjectID   `bson:"_id,omitempty"`
	URL                  string               `bson:"url"`
	CanonicalURL         string               `bson:"canonicalUrl"`
	Probability          float64              `bson:"probability"`
	Name                 string               `bson:"name"`
	Offers               []Offer              `bson:"offers"`
	SKU                  string               `bson:"sku"`
	MPN                  string               `bson:"mpn"`
	Brand                string               `bson:"brand"`
	Breadcrumbs          []Breadcrumb         `bson:"breadcrumbs"`
	MainImage            string               `bson:"mainImage"`
	Images               []string             `bson:"images"`
	Description          string               `bson:"description"`
	DescriptionHTML      string               `bson:"descriptionHtml"`
	AdditionalProperties []AdditionalProperty `bson:"additionalProperty"`
	AggregateRating      AggregateRating      `bson:"aggregateRating"`
	Color                string               `bson:"color"`
	Style                string               `bson:"style"`
	DateCreated          CustomTime           `bson:"dateCreated"`
	DateModified         CustomTime           `bson:"dateModified"`
	isNew                bool                 `bson:"isNew"`
	isDeleted            bool                 `bson:"isDeleted"`
	isDirty              bool                 `bson:"isDirty"`
	originalState        *Product             `bson:"-"`
}

var mongoConfig db.MongoConfig
var _cfg *config.BizObjConfig

func New(options ...func(*Product)) *Product {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	product := &Product{}
	for _, option := range options {
		option(product)
	}
	return product
}

func WithMongoConfig(config db.MongoConfig) func(*Product) {
	return func(product *Product) {
		mongoConfig = config
	}
}

func WithLoggingLevel(level string) func(*Product) {
	return func(product *Product) {
		_cfg.Logging.Level = level
	}
}

func WithBizObjConfig(config config.BizObjConfig) func(*Product) {
	return func(product *Product) {
		_cfg = &config
		mongoConfig = db.MongoConfig{
			CollectionName: "products",
			DatabaseName:   config.MongoDB.Database,
			URI:            config.MongoDB.URI,
			Timeout:        config.MongoDB.Timeout,
	}
}

// This builder func is to be used by apps that use chux-bizobj as a dependent
func NewProduct(config config.BizObjConfig) (*Product, error) {
	// Use the provided config
	mongoConfig = db.MongoConfig{
		CollectionName: "products",
		DatabaseName:   config.MongoDB.Database,
		URI:            config.MongoDB.URI,
		Timeout:        config.MongoDB.Timeout,
	}
	var err error
	mongoDB, err = db.NewMongoDB(mongoConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create a new MongoDB: %v", err))
	}

	return &Product{
		isDirty:   false,
		isNew:     true,
		isDeleted: false,
	}, err
}

// This builder func is provided if the configuration is given
// Locally to chux-bizobj by using the yml files in the config package
// This func should not be used to build a Product if chux-bizobj
// is a dependent of another library or application. In these
// cases, use the NewProduct(config config.BizObjConfig) builder
func NewProductWithDefaultConfig() (*Product, error) {

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	_cfg, err := config.LoadConfig(env)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return NewProductWithCustomURI(_cfg.MongoDB.URI)
}

// This builder function was added to allow an adhoc URI to be issued to
// `Product` this is necessary for unit tests and could be useful in
// other edge use cases as well
func NewProductWithCustomURI(customURI string) (*Product, error) {

	var err error

	if _cfg == nil {
		env := os.Getenv("APP_ENV")
		if env == "" {
			env = "development"
		}

		_cfg, err = config.LoadConfig(env)
		if err != nil {
			panic(fmt.Sprintf("failed to load config: %v", err))
		}
	}

	// Use the provided customURI instead of the one from the config
	mongoConfig = db.MongoConfig{
		CollectionName: "products",
		DatabaseName:   _cfg.MongoDB.Database,
		URI:            customURI,
		Timeout:        _cfg.MongoDB.Timeout,
	}

	mongoDB, err = db.NewMongoDB(mongoConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create a new MongoDB: %v", err))
	}

	return &Product{
		isDirty:   false,
		isNew:     true,
		isDeleted: false,
	}, err
}

func (p *Product) GetCollectionName() string {
	return mongoConfig.CollectionName
}

func (p *Product) GetDatabaseName() string {
	return mongoConfig.DatabaseName
}

func (p *Product) GetURI() string {
	return mongoConfig.URI
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

// Saves the Model to a Data Store
func (p *Product) Save() error {
	if p.isNew {
		//--Create a new document
		err := mongoDB.Create(p)
		if err != nil {
			return err
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
