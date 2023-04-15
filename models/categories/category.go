package categories

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/chuxorg/chux-datastore/db"
	"github.com/chuxorg/chux-models/config"
	"github.com/chuxorg/chux-models/errors"
	"github.com/chuxorg/chux-models/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID        primitive.ObjectID `bson:"_id,omitempty"`
	ProductID primitive.ObjectID `bson:"product_id"`
	Name      string             `bson:"name"`
	Index     int                `bson:"index"`
<<<<<<< HEAD
=======
	DateCreated          models.CustomTime           `bson:"dateCreated,omitempty" json:"dateCreated,omitempty"`
	DateModified         models.CustomTime           `bson:"dateModified,omitempty" json:"dateModified,omitempty"`
	isNew                bool                        `bson:"isNew"`
	isDeleted            bool                        `bson:"isDeleted"`
	isDirty              bool                        `bson:"isDirty"`
	originalState        *Category                   `bson:"-"`
	ParentCategoryID     primitive.ObjectID          `bson:"parentCategory,omitempty" json:"parentCategory,omitempty"`
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
}

var _cfg *config.BizObjConfig
var mongoDB *db.MongoDB

func New(options ...func(*Category)) *Category {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	_cfg = config.New()
<<<<<<< HEAD
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
=======
	category := &Category{}
	for _, option := range options {
		option(category)
	}

	mongoDB = db.New(
		db.WithURI(category.GetURI()),
		db.WithDatabaseName(category.GetDatabaseName()),
		db.WithCollectionName(category.GetCollectionName()),
		db.WithTimeout(float64(_cfg.DataStores.DataStoreMap["mongo"].Timeout)),
	)

	category.isNew = true
	category.isDeleted = false
	category.isDirty = false
	return category
}

func WithLoggingLevel(level string) func(*Category) {
	return func(category *Category) {
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
		_cfg.Logging.Level = level
	}
}

<<<<<<< HEAD
func WithBizObjConfig(config config.BizObjConfig) func(*Product) {
	return func(product *Product) {
=======
func WithBizObjConfig(config config.BizObjConfig) func(*Category) {
	return func(product *Category) {
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
		_cfg = &config
	}
}

<<<<<<< HEAD
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
=======
func (c *Category) GetCollectionName() string {
	return "categories"
}

func (c *Category) GetDatabaseName() string {
	return _cfg.DataStores.DataStoreMap["mongo"].DatabaseName
}

func (c *Category) GetURI() string {
	return _cfg.DataStores.DataStoreMap["mongo"].URI
}

func (c *Category) GetID() primitive.ObjectID {
	return c.ID
}

// If the Model has changes, will return true
func (c *Category) IsDirty() bool {
	if c.originalState == nil {
		return false
	}

	originalBytes, err := c.originalState.Serialize()
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
	if err != nil {
		return false
	}

<<<<<<< HEAD
	currentBytes, err := p.Serialize()
=======
	currentBytes, err := c.Serialize()
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
	if err != nil {
		return false
	}

<<<<<<< HEAD
	p.isDirty = string(originalBytes) != string(currentBytes)
	return p.isDirty
=======
	c.isDirty = string(originalBytes) != string(currentBytes)
	return c.isDirty
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
}

// When the Model is first created,
// the model is considered New. After the model is
// Saved or Loaded it is no longer New
<<<<<<< HEAD
func (p *Product) IsNew() bool {
	return p.isNew
}

func (p *Product) Exists() ([]db.IMongoDocument, error) {

	docs, err := mongoDB.Query(p, "canonicalUrl", p.CanonicalURL, "description", p.Description)
	if err != nil {
		return nil, errors.NewChuxModelsError("Product.Exists() Error querying database", err)
=======
func (c *Category) IsNew() bool {
	return c.isNew
}

func (c *Category) Exists() ([]db.IMongoDocument, error) {

	docs, err := mongoDB.Query(c, "name", c.Name)
	if err != nil {
		return nil, errors.NewChuxModelsError("Category.Exists() Error querying database", err)
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
	}

	return docs, nil
}

<<<<<<< HEAD
// CompareProducts takes two Product structs and compares their fields to see if anything has changed.
// Returns a map containing the field names as keys and a tuple of the old and new values as the corresponding values.
func CompareProducts(oldProduct, newProduct Product) (map[string][2]interface{}, error) {
	changes := make(map[string][2]interface{})

	v1 := reflect.ValueOf(oldProduct)
	v2 := reflect.ValueOf(newProduct)

	// Loop through the fields of the Product struct
=======
// CompareCategories takes two Category Structs  and compares their fields to see if anything has changed.
// Returns a map containing the field names as keys and a tuple of the old and new values as the corresponding values.
func CompareCategories(oldCategory, newCategory Category) (map[string][2]interface{}, error) {
	changes := make(map[string][2]interface{})

	v1 := reflect.ValueOf(oldCategory)
	v2 := reflect.ValueOf(newCategory)

	// Loop through the fields of the Category struct
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
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
<<<<<<< HEAD
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
=======
func (c *Category) Save() error {

	if c.isNew {
		var exists bool
		changes := make(map[string][2]interface{})
		categories, err := c.Exists()
		if err != nil {
			return errors.NewChuxModelsError("Category.Save() Error checking if category exists", err)
		}

		if len(categories) > 0 {
			category := categories[0].(*Category)
			var err error
			c.ID = category.ID
			//-- The c.Exists() call above will return a slice of categories if there are any
			//-- that match the name. We need to compare the
			//-- incoming category with the existing category in Mongo to see if there are any
			//-- changes. If there are, we need to update the existing category in Mongo with
			//-- the new values.
			changes, err = CompareCategories(*category, *c)
			if err != nil {
				return errors.NewChuxModelsError("Category.Save() Error comparing Categories", err)
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
			}
			if len(changes) > 0 {
				fmt.Println(changes)
			}
		}
<<<<<<< HEAD
		exists = len(products) > 0
		if !exists {
			//-- This checks that the product does not exist in the database.
			//-- it is required because new products that come in from a parse
			//-- run will bring in dupes, which I don't want to save.

			// -- Set the company name
			p.CompanyName, err = models.ExtractCompanyName(p.CanonicalURL)
			if err != nil {				
				return errors.NewChuxModelsError("Product.Save() Error extracting Product.CompanyName", err)
			}
			// -- Set the date created to now
			p.DateCreated.Now()
			//-- Create a new document
			err := mongoDB.Create(p)
			if err != nil {
				return errors.NewChuxModelsError("Product.Save() Error creating Product in MongoDB", err)
=======
		exists = len(categories) > 0
		if !exists {
			//-- This checks that the category does not exist in the database.
			//-- it is required because new categories that come in from a parse
			//-- run will bring in dupes, which I don't want to save.

			// -- Set the date created to now
			c.DateCreated.Now()
			//-- Create a new document
			err := mongoDB.Create(c)
			if err != nil {
				return errors.NewChuxModelsError("Category.Save() Error creating Category in MongoDB", err)
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
			}
		} else {
			if len(changes) > 0 {
				// The product could exist in Mongo yet was changed on the website on the last crawl
				// If it did, the product will be updated with the new data
<<<<<<< HEAD
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
=======
				c.isNew = false
				c.isDirty = true
				//-- this will cause an update to the document
				return c.Save()
			}
		}

	} else if c.IsDirty() && !c.isDeleted {
		// Ensure the ID is a valid hex string representation of an ObjectID
		_, err := primitive.ObjectIDFromHex(c.ID.Hex())
		if err != nil {
			return errors.NewChuxModelsError("Category.Save() invalid ObjectID", err)
		}
		// -- Set the date modified to now
		c.DateModified.Now()
		//--update this document
		err = mongoDB.Update(c, c.ID.Hex())
		if err != nil {
			return errors.NewChuxModelsError("Category.Save() Error updating the Category in MongoDB", err)
		}
	} else if c.isDeleted && !c.isNew {
		//--delete the document
		err := mongoDB.Delete(c, c.ID.Hex())
		if err != nil {
			return errors.NewChuxModelsError("Category.Save() Error deleting Product in MongoDB", err)
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
		}
	}

	// If the Product has been deleted, then this is a new Product
<<<<<<< HEAD
	p.isNew = p.isDeleted
	// little confusing but use the IsDirty() func to set isDirty field on Product struct
	p.isDirty = p.IsDirty()
	p.isDeleted = false
=======
	c.isNew = c.isDeleted
	// little confusing but use the IsDirty() func to set isDirty field on Product struct
	c.isDirty = c.IsDirty()
	c.isDeleted = false
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e

	// serialized will help set the current state
	var serialized string
	var err error
<<<<<<< HEAD
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
=======
	if c.isNew {
		serialized = ""
		c.originalState = nil
	} else {
		//--reset state
		serialized, err = c.Serialize()
		if err != nil {
			return errors.NewChuxModelsError("Category.Save() Error serializing Product.", err)
		}
		c.SetState(serialized)
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
	}

	return nil
}

<<<<<<< HEAD
func ExtractCompanyName(s string) {
	panic("unimplemented")
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
=======


// Loads a Model from MongoDB by id
func (c *Category) Load(id string) (interface{}, error) {
	
	retVal, err := mongoDB.GetByID(c, id)
	if err != nil {
		return nil, errors.NewChuxModelsError("Category.Load() Error loading Category from MongoDB", err)
	}
	category, ok := retVal.(*Category)
	if !ok {
		return nil, fmt.Errorf("unable to cast retVal to *Category")
	}
	serialized, err := category.Serialize()
	if err != nil {
		return nil, fmt.Errorf("unable to set internal state")
	}
	c.SetState(serialized)
	c.isNew = false
	c.isDirty = false
	c.isDeleted = false
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e

	return retVal, nil
}

// Marks a Model for deletion from the Data Store
// when Save() is called, the Model will be deleted
<<<<<<< HEAD
func (p *Product) Delete() error {
	p.isDeleted = true
=======
func (c *Category) Delete() error {
	c.isDeleted = true
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
	return nil
}

// Sets the internal state of the model.
<<<<<<< HEAD
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
=======
func (c *Category) SetState(json string) error {
	// Store the current state as the original state
	original := &Category{}
	*original = *c
	c.originalState = original

	// Deserialize the new state
	return c.Deserialize([]byte(json))
}

// Sets the internal state of the model of a new Category
// from a JSON String.
func (c *Category) Parse(json string) error {
	err := c.SetState(json)
	c.isNew = true // this is a new model
	return err
}

func (c *Category) Search(args ...interface{}) ([]interface{}, error) {
	return nil, nil
}

func (c *Category) Serialize() (string, error) {
	bytes, err := json.Marshal(c)
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

<<<<<<< HEAD
func (p *Product) Deserialize(jsonData []byte) error {
	err := json.Unmarshal(jsonData, p)
=======
func (c *Category) Deserialize(jsonData []byte) error {
	err := json.Unmarshal(jsonData, c)
>>>>>>> 625cc0e9edd33f30d7e59f65a87b1848deb7403e
	if err != nil {
		return err
	}
	return nil
}
