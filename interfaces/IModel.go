package interfaces

import (
	"github.com/chuxorg/chux-datastore/db"
	"github.com/chuxorg/chux-models/config"
)

// An Interface for Models that interact with a data store
type IModel interface {
	// If the Model has changes, will return true
	IsDirty() bool
	// When the Model is first created,
	// the model is considered New. After the model is
	// Saved or Loaded it is no longer New
	IsNew() bool
	// Saves the Model to a Data Store
	Save() error
	// Loads a Model from the Data Store
	Load(id string) (interface{}, error)
	// Loads a Model from the Data Store based on a query
	Query(args ...interface{}) ([]db.IMongoDocument, error)
	// Searches for items in the data store
	Search(args ...interface{}) ([]interface{}, error)
	// Deletes a Model from the Data Store
	Delete() error
	// Sets the internal state of the model.
	SetState(json string) error
	// Applies variations of constructor functions to the model
	Apply(opts ...func(IModel))
	SetLoggingLevel(level string)
	SetBizObjConfig(config config.BizObjConfig)
	SetDataStoresConfig(config config.DataStoresConfig)
}
