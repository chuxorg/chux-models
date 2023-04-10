package articles

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/chuxorg/chux-datastore/db"
	"github.com/chuxorg/chux-models/config"
	"github.com/chuxorg/chux-models/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// The Article struct represents an Article Document in MongoDB
type Article struct {
	ID               primitive.ObjectID  `bson:"_id,omitempty"`
	URL              string              `bson:"url"`
	Probability      float64             `bson:"probability"`
	Headline         string              `bson:"headline"`
	DatePublished    models.CustomTime   `bson:"datePublished"`
	DatePublishedRaw string              `bson:"datePublishedRaw"`
	DateCreated      models.CustomTime   `bson:"dateCreated"`
	DateModified     models.CustomTime   `bson:"dateModified"`
	DateModifiedRaw  string              `bson:"dateModifiedRaw"`
	Author           string              `bson:"author"`
	AuthorsList      []string            `bson:"authorsList"`
	InLanguage       string              `bson:"inLanguage"`
	Breadcrumbs      []models.Breadcrumb `bson:"breadcrumbs"`
	MainImage        string              `bson:"mainImage"`
	Images           []string            `bson:"images"`
	Description      string              `bson:"description"`
	ArticleBody      string              `bson:"articleBody"`
	ArticleBodyHTML  string              `bson:"articleBodyHtml"`
	CanonicalURL     string              `bson:"canonicalUrl"`
	isNew            bool                `bson:"isNew"`
	isDeleted        bool                `bson:"isDeleted"`
	isDirty          bool                `bson:"isDirty"`
	originalState    *Article            `bson:"-"`
}

var _cfg *config.BizObjConfig
var mongoDB *db.MongoDB

func New(options ...func(*Article)) *Article {
	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}
	var err error
	_cfg, err = config.LoadConfig(env)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	article := &Article{}
	for _, option := range options {
		option(article)
	}

	mongoDB = db.New(
		db.WithURI(article.GetURI()),
		db.WithDatabaseName(article.GetDatabaseName()),
		db.WithCollectionName(article.GetCollectionName()),
	)

	article.isNew = true
	article.isDeleted = false
	article.isDirty = false
	return article
}

func WithLoggingLevel(level string) func(*Article) {
	return func(article *Article) {
		_cfg.Logging.Level = level
	}
}

func WithBizObjConfig(config config.BizObjConfig) func(*Article) {
	return func(article *Article) {
		_cfg = &config
	}
}

func (a *Article) GetCollectionName() string {
	return "articles"
}

func (a *Article) GetDatabaseName() string {
	return _cfg.DataStores.DataStoreMap["mongo"].DatabaseName
}

func (a *Article) GetURI() string {
	return _cfg.DataStores.DataStoreMap["mongo"].URI
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
		err := mongoDB.Create(a)
		if err != nil {
			return err
		}

	} else if a.IsDirty() && !a.isDeleted {
		// Ensure the ID is a valid hex string representation of an ObjectID
		_, err := primitive.ObjectIDFromHex(a.ID.Hex())
		if err != nil {
			return fmt.Errorf("invalid ObjectID: %v", err)
		}
		//--update this document
		err = mongoDB.Update(a, a.ID.Hex())
		if err != nil {
			return err
		}
	} else if a.isDeleted && !a.isNew {
		//--delete the document
		err := mongoDB.Delete(a, a.ID.Hex())
		if err != nil {
			return err
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
			return fmt.Errorf("unable to set internal state")
		}
		a.SetState(serialized)
	}

	return nil
}

// Loads a Model from MongoDB by id
func (a *Article) Load(id string) (interface{}, error) {
	retVal, err := mongoDB.GetByID(a, id)
	if err != nil {
		return nil, err
	}
	article, ok := retVal.(*Article)
	if !ok {
		return nil, fmt.Errorf("unable to cast retVal to *Article")
	}
	serialized, err := article.Serialize()
	if err != nil {
		return nil, fmt.Errorf("unable to set internal state")
	}
	a.SetState(serialized)
	a.isNew = false
	a.isDirty = false
	a.isDeleted = false

	return retVal, nil
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
	return err
}

func (a *Article) Search(args ...interface{}) ([]interface{}, error) {
	return nil, nil
}

func (a *Article) Serialize() (string, error) {
	bytes, err := json.Marshal(a)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (a *Article) Deserialize(jsonData []byte) error {
	err := json.Unmarshal(jsonData, a)
	if err != nil {
		return err
	}
	return nil
}
