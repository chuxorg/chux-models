package models

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/csailer/chux-bizobj/config"
	"github.com/csailer/chux-mongo/db"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// The Article struct represents an Article Document in MongoDB
type Article struct {
	ID               primitive.ObjectID `bson:"_id,omitempty"`
	URL              string             `bson:"url"`
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

// Create the MongoDB from chux-mongo
var mongoDB = &db.MongoDB{}

// This builder func is to be used by apps that use chux-bizobj as a dependent
func NewArticle(config config.BizObjConfig) (*Article, error) {
	// Use the provided config
	mongoConfig = db.MongoConfig{
		CollectionName: "articles",
		DatabaseName:   config.MongoDB.Database,
		URI:            config.MongoDB.URI,
		Timeout:        config.MongoDB.Timeout,
	}
	var err error
	mongoDB, err = db.NewMongoDB(mongoConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create a new MongoDB: %v", err))
	}

	return &Article{
		isDirty:   false,
		isNew:     true,
		isDeleted: false,
	}, err
}

// This builder func is provided if the configuration is given
// Locally to chux-bizobj by using the yml files in the config package
// This func should not be used to build a Product if chux-bizobj
// is a dependent of another library or application. In these
// cases, use the NewArticle(config config.BizObjConfig) builder
func NewArticleWithDefaultConfig() (*Article, error) {

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	_cfg, err := config.LoadConfig(env)
	if err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return NewArticleWithCustomURI(_cfg.MongoDB.URI)
}

// This builder function was added to allow an adhoc URI to be issued to
// `Article` this is necessary for unit tests and could be useful in
// other edge use cases as well
func NewArticleWithCustomURI(customURI string) (*Article, error) {

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
		CollectionName: "articles",
		DatabaseName:   _cfg.MongoDB.Database,
		URI:            customURI,
		Timeout:        _cfg.MongoDB.Timeout,
	}

	mongoDB, err = db.NewMongoDB(mongoConfig)
	if err != nil {
		panic(fmt.Sprintf("failed to create a new MongoDB: %v", err))
	}

	return &Article{
		isDirty:   false,
		isNew:     true,
		isDeleted: false,
	}, err
}

func (a *Article) GetCollectionName() string {
	return mongoConfig.CollectionName
}

func (a *Article) GetDatabaseName() string {
	return mongoConfig.DatabaseName
}

func (a *Article) GetURI() string {
	return mongoConfig.URI
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
	// little confusing but use the IsDirty() func to set isDirty field on the Article struct
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
		return nil, fmt.Errorf("unable to cast retVal to *Product")
	}
	serialized, err := article.Serialize()
	if err != nil {
		return nil, fmt.Errorf("unable to set internal state")
	}
	a.SetState(serialized)
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
