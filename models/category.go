package models

import (
	"encoding/json"
	"fmt"
	"os"
	"reflect"

	"github.com/chuxorg/chux-datastore/db"
	dbl "github.com/chuxorg/chux-datastore/logging"
	"github.com/chuxorg/chux-models/errors"
	"github.com/chuxorg/chux-models/logging"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Category struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	ProductID     primitive.ObjectID `bson:"product_id"`
	Name          string             `bson:"name"`
	Index         int                `bson:"index"`
	ParentID      primitive.ObjectID `bson:"parent_id"`
	isNew         bool               `bson:"isNew,omitempty" json:"isNew,omitempty"`
	isDeleted     bool               `bson:"isDeleted,omitempty" json:"isDeleted,omitempty"`
	isDirty       bool               `bson:"isDirty,omitempty" json:"isDirty,omitempty"`
	originalState *Category          `bson:"-" json:"-"`
	DateCreated   CustomTime         `bson:"dateCreated,omitempty" json:"dateCreated,omitempty"`
	DateModified  CustomTime         `bson:"dateModified,omitempty" json:"dateModified,omitempty"`
	Logger        *logging.Logger    `bson:"-" json:"-"`
}

// Creates a NewCategory with Options
func NewCategory(options ...func(*Category)) *Category {

	c := &Category{}

	for _, option := range options {
		option(c)
	}
	dbLogger := dbl.NewLogger(dbl.LogLevelDebug)
	mongoDB = db.New(
		db.WithURI(c.GetURI()),
		db.WithDatabaseName(c.GetDatabaseName()),
		db.WithCollectionName(c.GetCollectionName()),
		db.WithTimeout(30),
		db.WithLogger(*dbLogger),
	)

	c.isNew = true
	c.isDeleted = false
	c.isDirty = false

	return c
}

func NewCategoryWithLogger(logger logging.Logger) func(*Category) {
	return func(c *Category) {
		c.Logger = &logger
	}
}

// GetCollectionName returns the name of the collection
func (c *Category) GetCollectionName() string {
	c.Logger.Debug("GetCollectionName() called")
	return "categories"
}

// GetDatabaseName returns the name of the database
func (c *Category) GetDatabaseName() string {
	c.Logger.Debug("GetDatabaseName() called")
	return os.Getenv("MONGO_DATABASE")
}

func (c *Category) GetURI() string {
	logging := c.Logger
	c.Logger.Debug("GetURI() called")
	username := os.Getenv("MONGO_USER_NAME")
	password := os.Getenv("MONGO_PASSWORD")

	uri := os.Getenv("MONGO_URI")
	mongoURI := fmt.Sprintf(uri, username, password)
	masked := fmt.Sprintf(uri, "********", "********")
	logging.Info("GetURI() returning: %s", masked)

	return mongoURI
}

func (c *Category) GetID() primitive.ObjectID {
	logging := c.Logger
	logging.Debug("GetID() called")
	return c.ID
}

func (c *Category) SetID(id primitive.ObjectID) {
	logging := c.Logger
	logging.Debug("SetID() called")
	c.ID = id
}

// If the Model has changes, will return true
func (c *Category) IsDirty() bool {
	logging := c.Logger
	logging.Debug("IsDirty() called")

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
	logging.Info("IsDirty() returning: %t", c.isDirty)
	return c.isDirty
}

// When the Model is first created,
// the model is considered New. After the model is
// Saved or Loaded it is no longer New
func (c *Category) IsNew() bool {
	logging := c.Logger
	logging.Debug("IsNew() called")
	return c.isNew
}

func (c *Category) Exists() ([]db.IMongoDocument, error) {
	logging := c.Logger
	logging.Debug("Exists() called")
	docs, err := mongoDB.Query(c, "name", c.Name)
	if err != nil {
		return nil, errors.NewChuxModelsError("Category.Exists() Error querying database", err)
	}

	return docs, nil
}

// CompareCategories takes two Category Structs  and compares their fields to see if anything has changed.
// Returns a map containing the field names as keys and a tuple of the old and new values as the corresponding values.
func (c *Category) CompareCategories(oldCategory, newCategory Category) (map[string][2]interface{}, error) {
	logging := c.Logger
	logging.Debug("CompareCategories() called")
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
	logging := c.Logger
	logging.Debug("Save() called")
	if c.isNew {
		logging.Info("Save() Category isNew")
		// -- Set the date created to now
		c.DateCreated.Now()
		//-- Create a new document
		err := mongoDB.Upsert(c)
		if err != nil {
			logging.Error("Save() Error creating Category in MongoDB: %s", err.Error())
			return errors.NewChuxModelsError("Category.Save() Error creating Category in MongoDB", err)
		}

	} else if c.IsDirty() && !c.isDeleted {
		logging.Info("Save() Category isDirty")
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
		logging.Info("Category.Save() Category isDeleted and is not New")
		//--delete the document
		err := mongoDB.Delete(c, c.ID.Hex())
		if err != nil {
			logging.Error("Category.Save() Error deleting Category in MongoDB: %s", err.Error())
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
	logging.Info("Category.Save() returning successfully.")
	return nil
}

// Loads a Model from MongoDB by id
func (c *Category) Load(id string) (interface{}, error) {
	logging := c.Logger
	logging.Debug("Category.Load() called")
	retVal, err := mongoDB.GetByID(c, id)
	if err != nil {
		logging.Error("Category.Load() Error loading Category from MongoDB: %s", err.Error())
		return nil, errors.NewChuxModelsError("Category.Load() Error loading Category from MongoDB", err)
	}
	category, ok := retVal.(*Category)
	if !ok {
		logging.Error("Category.Load() Error casting retVal to *Category")
		return nil, fmt.Errorf("unable to cast retVal to *Category")
	}
	serialized, err := category.Serialize()
	if err != nil {
		logging.Error("Category.Load() Error serializing Category: %s", err.Error())
		return nil, fmt.Errorf("unable to set internal state")
	}
	c.SetState(serialized)
	c.isNew = false
	c.isDirty = false
	c.isDeleted = false
	logging.Info("Category.Load() returning successfully.")
	return retVal, nil
}

func (c *Category) Query(args ...interface{}) ([]db.IMongoDocument, error) {
	logging := c.Logger
	logging.Debug("Category.Query() called")
	results, err := mongoDB.Query(c, args...)
	if err != nil {
		logging.Error("Category.Query() Error occurred querying Categories: %s", err.Error())
		return nil, errors.NewChuxModelsError("Category.Query() Error occurred querying Categories", err)
	}
	logging.Info("Category.Query() returning successfully." + fmt.Sprintf("Found %d Categories", len(results)))
	return results, nil
}

// Marks a Model for deletion from the Data Store
// when Save() is called, the Model will be deleted
func (c *Category) Delete() error {
	logging := c.Logger
	logging.Debug("Category.Delete() called")
	c.isDeleted = true
	return nil
}

// Sets the internal state of the model.
func (c *Category) SetState(json string) error {
	logging := c.Logger
	logging.Debug("Category.SetState() called")
	// Store the current state as the original state
	original := &Category{}
	*original = *c
	c.originalState = original

	// Deserialize the new state
	logging.Debug("Category.SetState() calling c.Deserialize()")
	return c.Deserialize([]byte(json))
}

// Sets the internal state of the model of a new Category
// from a JSON String.
func (c *Category) Parse(json string) error {
	logging := c.Logger
	logging.Debug("Category.Parse() called")
	err := c.SetState(json)
	if err != nil {
		logging.Error("Category.Parse() Error setting state: %s", err.Error())
		return errors.NewChuxModelsError("Category.Parse() Error setting state", err)
	}
	c.isNew = true // this is a new model
	return nil
}

func (c *Category) Search(args ...interface{}) ([]interface{}, error) {
	return nil, nil
}

func (c *Category) Serialize() (string, error) {
	logging := c.Logger
	logging.Debug("Category.Serialize() called")
	bytes, err := json.Marshal(c)
	if err != nil {
		logging.Error("Category.Serialize() Error serializing Category: %s", err.Error())
		return "", errors.NewChuxModelsError("Category.Serialize() Error serializing Category", err)
	}
	return string(bytes), nil
}

func (c *Category) Deserialize(jsonData []byte) error {
	logging := c.Logger
	logging.Debug("Category.Deserialize() called")
	err := json.Unmarshal(jsonData, c)
	if err != nil {
		logging.Error("Category.Deserialize() Error deserializing Category: %s", err.Error())
		return errors.NewChuxModelsError("Category.Deserialize() Error deserializing Category", err)
	}
	return nil
}
