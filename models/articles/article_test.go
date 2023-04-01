package articles

import (
	"os"
	"testing"

	"github.com/chuxorg/chux-models/config"
	"github.com/stretchr/testify/assert"
)

// TestNew tests the New function with different options.
func TestNew(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	article := New()
	assert.NotNil(t, article)
	assert.Equal(t, "testdb", article.GetDatabaseName())
	assert.Equal(t, "articles", article.GetCollectionName())
	assert.Equal(t, "mongodb://localhost:27017", article.GetURI())

	articleWithLoggingLevel := New(WithLoggingLevel("debug"))

	assert.NotNil(t, articleWithLoggingLevel)
	assert.Equal(t, "debug", _cfg.Logging.Level)

	customConfig := config.BizObjConfig{
		Logging: struct {
			Level string `mapstructure:"level"`
		}{
			Level: "info",
		},
		DataStores: struct {
			DataStoreMap map[string]config.DataStoreConfig `mapstructure:"dataStore"`
		}{
			DataStoreMap: map[string]config.DataStoreConfig{
				"mongo": {
					Target:         "mongo",
					URI:            "mongodb://localhost:27017",
					Timeout:        10,
					DatabaseName:   "customdb",
					CollectionName: "customcollection",
				},
			},
		},
	}

	articleWithCustomConfig := New(WithBizObjConfig(customConfig))

	assert.NotNil(t, articleWithCustomConfig)
	assert.Equal(t, "customdb", articleWithCustomConfig.GetDatabaseName())
	// Should return articles instead of customdb
	assert.Equal(t, "articles", articleWithCustomConfig.GetCollectionName())
	assert.Equal(t, "mongodb://localhost:27017", articleWithCustomConfig.GetURI())
}

// TestWithLoggingLevel tests the WithLoggingLevel function.
func TestWithLoggingLevel(t *testing.T) {
	article := &Article{}
	withLoggingLevel := WithLoggingLevel("error")
	withLoggingLevel(article)

	assert.Equal(t, "error", _cfg.Logging.Level)
}

// TestWithBizObjConfig tests the WithBizObjConfig function.
func TestWithBizObjConfig(t *testing.T) {
	customConfig := config.BizObjConfig{
		Logging: struct {
			Level string `mapstructure:"level"`
		}{
			Level: "warning",
		},
		DataStores: struct {
			DataStoreMap map[string]config.DataStoreConfig `mapstructure:"dataStore"`
		}{
			DataStoreMap: map[string]config.DataStoreConfig{
				"mongo": {
					Target:         "mongo",
					URI:            "mongodb://localhost:27017",
					Timeout:        10,
					DatabaseName:   "customdb",
					CollectionName: "customcollection",
				},
			},
		},
	}

	/*
	    article := &Article{}: This line creates a new Article struct and assigns its address to the article variable. 
		The & symbol is used to get the address of the newly created struct.

    	withBizObjConfig := WithBizObjConfig(customConfig): This line calls the WithBizObjConfig function with a custom configuration 
		(assumed to be of type config.BizObjConfig). 
		The function returns a closure (a function with access to the variables from its parent scope) that takes 
		an *Article as an argument. The closure is assigned to the withBizObjConfig variable.

    	withBizObjConfig(article): This line calls the closure stored in the withBizObjConfig variable, passing in the product variable 
		(which is a pointer to an Article struct). This closure sets the _cfg global variable to the custom configuration passed 
		to the WithBizObjConfig function.
	*/

	article := &Article{}
	withBizObjConfig := WithBizObjConfig(customConfig)
	withBizObjConfig(article)

	assert.Equal(t, "warning", _cfg.Logging.Level)
	assert.Equal(t, "customdb", _cfg.DataStores.DataStoreMap["mongo"].DatabaseName)
	assert.Equal(t, "customcollection", _cfg.DataStores.DataStoreMap["mongo"].CollectionName)
	assert.Equal(t, "mongodb://localhost:27017", _cfg.DataStores.DataStoreMap["mongo"].URI)
}
