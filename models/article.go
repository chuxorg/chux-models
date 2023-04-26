package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/chuxorg/chux-datastore/db"
	dbl "github.com/chuxorg/chux-datastore/logging"
	"github.com/chuxorg/chux-models/errors"
	"github.com/chuxorg/chux-models/logging"

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
	FilesProcessed   bool               `bson:"filesProcessed" json:"filesProcessed"`
	ImagesProcessed  bool               `bson:"imagesProcessed" json:"imagesProcessed"`
	originalState    *Article           `bson:"-"`
	Logger           *logging.Logger    `bson:"-"`
}

func NewArticle(options ...func(*Article)) *Article {

	a := &Article{}

	for _, option := range options {
		option(a)
	}
	dbLogger := dbl.NewLogger(dbl.LogLevelDebug)
	mongoDB = db.New(
		db.WithURI(a.GetURI()),
		db.WithDatabaseName(a.GetDatabaseName()),
		db.WithCollectionName(a.GetCollectionName()),
		db.WithTimeout(30),
		db.WithLogger(*dbLogger),
	)

	a.isNew = true
	a.isDeleted = false
	a.isDirty = false
	return a
}

func NewArticleWithLogger(logger logging.Logger) func(*Article) {
	return func(a *Article) {
		a.Logger = &logger
	}
}

func (a *Article) GetCollectionName() string {
	a.Logger.Debug("Article.GetCollectionName() called")
	return "articles"
}

func (a *Article) GetDatabaseName() string {
	a.Logger.Debug("Article.GetDatabaseName() called")
	return os.Getenv("MONGO_DATABASE")
}

func (a *Article) GetURI() string {

	logging := a.Logger

	logging.Debug("Article.GetURI() called")
	username := os.Getenv("MONGO_USER_NAME")
	password := os.Getenv("MONGO_PASSWORD")

	uri := os.Getenv("MONGO_URI")
	mongoURI := fmt.Sprintf(uri, username, password)
	masked := fmt.Sprintf(uri, "********", "********")
	logging.Info("Article.GetURI() returning: %s", masked)
	return mongoURI
}

func (a *Article) GetID() primitive.ObjectID {
	a.Logger.Debug("Article.GetID() called")
	return a.ID
}

// If the Model has changes, will return true
func (a *Article) IsDirty() bool {
	a.Logger.Debug("Article.IsDirty() called")
	if a.originalState == nil {
		a.Logger.Info("Article.IsDirty() original state is nil")
		return false
	}

	originalBytes, err := a.originalState.Serialize()
	if err != nil {
		a.Logger.Error("Article.IsDirty() error serializing original state", err)
		return false
	}

	currentBytes, err := a.Serialize()
	if err != nil {
		a.Logger.Info("Article.IsDirty() Could not Serialize current state.", err)
		return false
	}

	a.isDirty = string(originalBytes) != string(currentBytes)
	a.Logger.Debug("Article.IsDirty() returning isDirty=%t", a.isDirty)
	return a.isDirty
}

// When the Model is first created,
// the model is considered New. After the model is
// Saved or Loaded it is no longer New
func (a *Article) IsNew() bool {
	a.Logger.Debug("Article.IsNew() called")
	return a.isNew
}

func (a *Article) SetID(id primitive.ObjectID) {
	a.Logger.Debug("Article.SetID() called")
	a.ID = id
}

// Saves the Model to a Data Store
func (a *Article) Save() error {
	logging := a.Logger
	a.Logger.Debug("Article.Save() called")
	if a.isNew {
		a.Logger.Debug("Article.Save() is new")
		//--Create a new document
		var err error
		a.CompanyName, err = ExtractCompanyName(a.CanonicalURL)
		if err != nil {
			logging.Error("Article.Save() error extracting company name", err)
			return errors.NewChuxModelsError("Article.Save() error extracting company name", err)
		}
		// Set the DateCreated to the current time
		a.DateCreated.Now()
		a.FilesProcessed = true
		err = mongoDB.Upsert(a)
		if err != nil {
			errors.NewChuxModelsError("Article.Save() error creating Article", err)
		}

		logging.Info("Article.Save() Successfully created new Article")

	} else if a.IsDirty() && !a.isDeleted {
		logging.Info("Article.Save() is dirty and not isDeleted")
		// Ensure the ID is a valid hex string representation of an ObjectID
		_, err := primitive.ObjectIDFromHex(a.ID.Hex())
		if err != nil {
			msg := fmt.Sprintf("Article.Save() invalid ObjectID: %v", err)
			logging.Error(msg)
			return errors.NewChuxModelsError(msg, err)
		}
		// Set the DateModified to the current time
		a.DateModified.Now()
		//--update this document
		err = mongoDB.Update(a, a.ID.Hex())
		if err != nil {
			logging.Error("Article.Save() error updating Article", err)
			return errors.NewChuxModelsError("Article.Save() error updating Article", err)
		}
		logging.Info("Article.Save() Successfully updated Article")
	} else if a.isDeleted && !a.isNew {
		logging.Info("Article.Save() isDeleted and not isNew")
		//--delete the document
		err := mongoDB.Delete(a, a.ID.Hex())
		if err != nil {
			logging.Error("Article.Save() error deleting Article", err)
			return errors.NewChuxModelsError("Article.Save() error deleting Article", err)
		}
		logging.Info("Article.Save() Successfully deleted Article")
	}

	// If the Article has been deleted, then this is a new Article
	a.isNew = a.isDeleted
	// little confusing but use the IsDirty() func to set isDirty field on Article struct
	a.isDirty = a.IsDirty()
	a.isDeleted = false

	// serialized will help set the current state
	var serialized string
	var err error
	logging.Info("Article.Save() Setting current state to original state after MongoDB operation.")
	if a.isNew {
		logging.Debug("Article.Save() is new so emptying original state")
		serialized = ""
		a.originalState = nil
	} else {
		logging.Info("Article.Save() is not new so setting original state to current state")
		//--reset state
		serialized, err = a.Serialize()
		if err != nil {
			logging.Error("Article.Save() unable to set current state", err)
			return errors.NewChuxModelsError("Article.Save() unable to set current state", err)
		}
		logging.Info("Article.Save() Setting internal state")
		err = a.SetState(serialized)
		if err != nil {
			logging.Error("Article.Save() unable to set current state", err)
			return errors.NewChuxModelsError("Article.Save() unable to set current state", err)
		}
	}

	return nil
}

// Loads a Model from MongoDB by id
func (a *Article) Load(id string) (interface{}, error) {
	logging := a.Logger
	logging.Debug("Article.Load() called")
	retVal, err := mongoDB.GetByID(a, id)
	if err != nil {
		logging.Error("Article.Load() error loading Article", err)
		return nil, errors.NewChuxModelsError("Article.Load() error loading Article", err)
	}
	article, ok := retVal.(*Article)
	if !ok {
		logging.Error("Article.Load() unable to cast retVal to *Article", err)
		return nil, errors.NewChuxModelsError("Article.Load() unable to cast retVal to *Article", err)
	}
	serialized, err := article.Serialize()
	if err != nil {
		logging.Error("Article.Load() unable to serialize Article", err)
		return nil, errors.NewChuxModelsError("Article.Load() unable to serialize Article", err)
	}
	logging.Info("Article.Load() Setting internal state")
	a.SetState(serialized)
	a.isNew = false
	a.isDirty = false
	a.isDeleted = false

	return retVal, nil
}

func (a *Article) Query(args ...interface{}) ([]db.IMongoDocument, error) {
	logging := a.Logger
	logging.Debug("Article.Query() called")
	results, err := mongoDB.Query(a, args...)
	if err != nil {
		logging.Error("Article.Query() Error occurred querying Articles", err)
		return nil, errors.NewChuxModelsError("Article.Query() Error occurred querying Articles", err)
	}

	return results, nil
}

// Marks a Model for deletion from the Data Store
// when Save() is called, the Model will be deleted
func (a *Article) Delete() error {
	a.Logger.Debug("Article.Delete() called")
	a.isDeleted = true
	return nil
}

// Sets the internal state of the model.
func (a *Article) SetState(json string) error {
	logging := a.Logger
	logging.Debug("Article.SetState() called")
	// Store the current state as the original state
	original := &Article{}
	*original = *a
	a.originalState = original

	// Deserialize the new state
	logging.Info("Article.SetState() deserializing new state and returning")
	return a.Deserialize([]byte(json))
}

// Sets the internal state of the model of a new Product
// from a JSON String.
func (a *Article) Parse(json string) error {
	logging := a.Logger

	logging.Debug("Article.Parse() called")
	err := a.SetState(json)
	a.isNew = true // this is a new model
	if err != nil {
		logging.Error("Article.Parse() unable to parse article", err)
		return errors.NewChuxModelsError("Article.Parse() unable to parse article", err)
	}
	return nil
}

func (a *Article) Search(args ...interface{}) ([]interface{}, error) {
	a.Logger.Debug("Article.Search() called")
	return nil, nil
}

func (a *Article) Serialize() (string, error) {
	logging := a.Logger
	logging.Debug("Article.Serialize() called")
	bytes, err := json.Marshal(a)
	if err != nil {
		logging.Error("Article.Serialize() unable to serialize Article", err)
		return "", errors.NewChuxModelsError("Article.Serialize() unable to serialize Article", err)
	}
	return string(bytes), nil
}

func (a *Article) Deserialize(jsonData []byte) error {
	logging := a.Logger
	logging.Debug("Article.Deserialize() called")
	err := json.Unmarshal(jsonData, a)
	if err != nil {
		logging.Error("Article.Deserialize() unable to deserialize Article", err)
		return errors.NewChuxModelsError("Article.Deserialize() unable to deserialize Article", err)
	}
	return nil
}
