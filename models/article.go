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

// The Article struct represents an Article Document in MongoDB
type Article struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	URL              string             `bson:"url"`
	CompanyName      string             `bson:"companyName, omitempty"`
	Probability      float64            `bson:"probability"`
	Headline         string             `bson:"headline"`
	DatePublished    CustomTime         `bson:"datePublished"`
	DatePublishedRaw string             `bson:"datePublishedRaw"`
	DateCreated      CustomTime         `bson:"dateCreated"`
	DateModified     CustomTime         `bson:"dateModified"`
	DateModifiedRaw  string             `bson:"dateModifiedRaw"`
	Author           string             `bson:"author"`
	AuthorsList      []string           `bson:"authorsList"`
	InLanguage       string             `bson:"inLanguage"`
	Breadcrumbs      []Breadcrumb       `bson:"breadcrumbs"`
	MainImage        string             `bson:"mainImage"`
	Images           []string           `bson:"images"`
	Description      string             `bson:"description"`
	ArticleBody      string             `bson:"articleBody"`
	ArticleBodyHTML  string             `bson:"articleBodyHtml"`
	CanonicalURL     string             `bson:"canonicalUrl"`
	isNew            bool               `bson:"isNew"`
	isDeleted        bool               `bson:"isDeleted"`
	isDirty          bool               `bson:"isDirty"`
	originalState    *Article           `bson:"-"`
}

func NewArticle(opts ...func(interfaces.IModel)) *Article {
	a := &Article{}
	a.Apply(opts...)
	return a
}

func (a *Article) Apply(opts ...func(interfaces.IModel)) {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	_cfg = config.New()

	for _, opt := range opts {
		opt(a)
	}

	mongoDB = db.New(
		db.WithURI(a.GetURI()),
		db.WithDatabaseName(a.GetDatabaseName()),
		db.WithCollectionName(a.GetCollectionName()),
	)

	a.isNew = true
	a.isDeleted = false
	a.isDirty = false
}

func (a *Article) SetLoggingLevel(level string) {
	_cfg.Logging.Level = level
}
func (a *Article) SetBizObjConfig(config config.BizObjConfig) {
	_cfg = &config
}
func (a *Article) SetDataStoresConfig(config config.DataStoresConfig) {
	_cfg.DataStores = config
}

func (a *Article) GetCollectionName() string {
	return "articles"
}

func (a *Article) GetDatabaseName() string {
	return os.Getenv("MONGO_DATABASE")
}

func (a *Article) GetURI() string {
	return os.Getenv("MONGO_URI")
}

func (a *Article) GetID() primitive.ObjectID {
	return a.ID
}

// If the Model has changes, will return true
func (a *Article) IsDirty() bool {
	if a.originalState == nil {
		return false
	}

	originalBytes, err := a.originalState.Serialize()
	if err != nil {
		return false
	}

	currentBytes, err := a.Serialize()
	if err != nil {
		return false
	}

	a.isDirty = string(originalBytes) != string(currentBytes)
	return a.isDirty
}

// When the Model is first created,
// the model is considered New. After the model is
// Saved or Loaded it is no longer New
func (a *Article) IsNew() bool {
	return a.isNew
}

// Saves the Model to a Data Store
func (a *Article) Save() error {
	if a.isNew {
		//--Create a new document
		var err error
		a.CompanyName, err = ExtractCompanyName(a.CanonicalURL)
		if err != nil {
			return errors.NewChuxModelsError("Artilce.Save() error extracting company name", err)
		}
		// Set the DateCreated to the current time
		a.DateCreated.Now()
		err = mongoDB.Create(a)
		if err != nil {
			errors.NewChuxModelsError("Artilce.Save() error creating Article", err)
		}

	} else if a.IsDirty() && !a.isDeleted {
		// Ensure the ID is a valid hex string representation of an ObjectID
		_, err := primitive.ObjectIDFromHex(a.ID.Hex())
		if err != nil {
			msg := fmt.Sprintf("Artilce.Save() invalid ObjectID: %v", err)
			return errors.NewChuxModelsError(msg, err)
		}
		// Set the DateModified to the current time
		a.DateModified.Now()
		//--update this document
		err = mongoDB.Update(a, a.ID.Hex())
		if err != nil {
			return errors.NewChuxModelsError("Artilce.Save() error updating Article", err)
		}
	} else if a.isDeleted && !a.isNew {
		//--delete the document
		err := mongoDB.Delete(a, a.ID.Hex())
		if err != nil {
			return errors.NewChuxModelsError("Artilce.Save() error deleting Article", err)
		}
	}

	// If the Article has been deleted, then this is a new Article
	a.isNew = a.isDeleted
	// little confusing but use the IsDirty() func to set isDirty field on Article struct
	a.isDirty = a.IsDirty()
	a.isDeleted = false

	// serialized will help set the current state
	var serialized string
	var err error
	if a.isNew {
		serialized = ""
		a.originalState = nil
	} else {
		//--reset state
		serialized, err = a.Serialize()
		if err != nil {
			return errors.NewChuxModelsError("Artilce.Save() unable to set current state", err)
		}
		a.SetState(serialized)
	}

	return nil
}

// Loads a Model from MongoDB by id
func (a *Article) Load(id string) (interface{}, error) {
	retVal, err := mongoDB.GetByID(a, id)
	if err != nil {
		return nil, errors.NewChuxModelsError("Artilce.Load() error loading Article", err)
	}
	article, ok := retVal.(*Article)
	if !ok {
		return nil, errors.NewChuxModelsError("Artilce.Load() unable to cast retVal to *Article", err)
	}
	serialized, err := article.Serialize()
	if err != nil {
		return nil, errors.NewChuxModelsError("Artilce.Load() unable to serialize Article", err)
	}
	a.SetState(serialized)
	a.isNew = false
	a.isDirty = false
	a.isDeleted = false

	return retVal, nil
}

func (a *Article) Query(args ...interface{}) ([]db.IMongoDocument, error) {
	results, err := mongoDB.Query(a, args...)
	if err != nil {
		return nil, errors.NewChuxModelsError("Article.Query() Error occurred querying Articles", err)
	}

	return results, nil
}

// Marks a Model for deletion from the Data Store
// when Save() is called, the Model will be deleted
func (a *Article) Delete() error {
	a.isDeleted = true
	return nil
}

// Sets the internal state of the model.
func (a *Article) SetState(json string) error {
	// Store the current state as the original state
	original := &Article{}
	*original = *a
	a.originalState = original

	// Deserialize the new state
	return a.Deserialize([]byte(json))
}

// Sets the internal state of the model of a new Product
// from a JSON String.
func (a *Article) Parse(json string) error {
	err := a.SetState(json)
	a.isNew = true // this is a new model
	if err != nil {
		return errors.NewChuxModelsError("Artilce.Parse() unable to parse article", err)
	}
	return nil
}

func (a *Article) Search(args ...interface{}) ([]interface{}, error) {
	return nil, nil
}

func (a *Article) Serialize() (string, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return "", errors.NewChuxModelsError("Artilce.Serialize() unable to serialize Article", err)
	}
	return string(bytes), nil
}

func (a *Article) Deserialize(jsonData []byte) error {
	err := json.Unmarshal(jsonData, a)
	if err != nil {
		return errors.NewChuxModelsError("Artilce.Deserialize() unable to deserialize Article", err)
	}
	return nil
}
