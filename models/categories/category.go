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
	ParentID  primitive.ObjectID `bson:"parent_id"`
	isNew                bool                        `bson:"isNew,omitempty" json:"isNew,omitempty"`
	isDeleted            bool                        `bson:"isDeleted,omitempty" json:"isDeleted,omitempty"`
	isDirty              bool                        `bson:"isDirty,omitempty" json:"isDirty,omitempty"`
	originalState        *Category                    `bson:"-" json:"-"`
	DateCreated          models.CustomTime           `bson:"dateCreated,omitempty" json:"dateCreated,omitempty"`
	DateModified         models.CustomTime           `bson:"dateModified,omitempty" json:"dateModified,omitempty"`
}

var _cfg *config.BizObjConfig
var mongoDB *db.MongoDB

func New(options ...func(*Category)) *Category {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	_cfg = config.New()
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
	return func(product *Category) {
		_cfg.Logging.Level = level
	}
}

func WithBizObjConfig(config config.BizObjConfig) func(*Category) {
	return func(product *Category) {
		_cfg = &config
	}
}

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
	if err != nil {
		return false
	}

	currentBytes, err := c.Serialize()
	if err != nil {
		return false
	}

	c.isDirty = string(originalBytes) != string(currentBytes)
	return c.isDirty
}

// When the Model is first created,
// the model is considered New. After the model is
// Saved or Loaded it is no longer New
func (c *Category) IsNew() bool {
	return c.isNew
}

func (c *Category) Exists() ([]db.IMongoDocument, error) {

	docs, err := mongoDB.Query(c, "name", c.Name)
	if err != nil {
		return nil, errors.NewChuxModelsError("Category.Exists() Error querying database", err)
	}

	return docs, nil
}

// CompareCategories takes two Category Structs  and compares their fields to see if anything has changed.
// Returns a map containing the field names as keys and a tuple of the old and new values as the corresponding values.
func CompareCategories(oldCategory, newCategory Category) (map[string][2]interface{}, error) {
	changes := make(map[string][2]interface{})

	v1 := reflect.ValueOf(oldCategory)
	v2 := reflect.ValueOf(newCategory)

	// Loop through the fields of the Category struct
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
			}
			if len(changes) > 0 {
				fmt.Println(changes)
			}
		}
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
			}
		} else {
			if len(changes) > 0 {
				// The product could exist in Mongo yet was changed on the website on the last crawl
				// If it did, the product will be updated with the new data
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
		}
	}

	// If the Product has been deleted, then this is a new Product
	c.isNew = c.isDeleted
	// little confusing but use the IsDirty() func to set isDirty field on Product struct
	c.isDirty = c.IsDirty()
	c.isDeleted = false

	// serialized will help set the current state
	var serialized string
	var err error
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
	}

	return nil
}



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

	return retVal, nil
}

// Marks a Model for deletion from the Data Store
// when Save() is called, the Model will be deleted
func (c *Category) Delete() error {
	c.isDeleted = true
	return nil
}

// Sets the internal state of the model.
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
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (c *Category) Deserialize(jsonData []byte) error {
	err := json.Unmarshal(jsonData, c)
	if err != nil {
		return err
	}
	return nil
}
