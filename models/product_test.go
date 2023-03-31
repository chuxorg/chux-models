package models

import (
	"os"
	"testing"

	"github.com/chuxorg/chux-models/config"
	"github.com/stretchr/testify/assert"
)

// TestNew tests the New function with different options.
func TestNew(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	product := New()

	assert.NotNil(t, product)
	assert.Equal(t, "testdb", product.GetDatabaseName())
	assert.Equal(t, "testcollection", product.GetCollectionName())
	assert.Equal(t, "mongodb://localhost:27017", product.GetURI())

	productWithLoggingLevel := New(WithLoggingLevel("debug"))

	assert.NotNil(t, productWithLoggingLevel)
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

	productWithCustomConfig := New(WithBizObjConfig(customConfig))

	assert.NotNil(t, productWithCustomConfig)
	assert.Equal(t, "customdb", productWithCustomConfig.GetDatabaseName())
	assert.Equal(t, "customcollection", productWithCustomConfig.GetCollectionName())
	assert.Equal(t, "mongodb://localhost:27017", productWithCustomConfig.GetURI())
}

// TestWithLoggingLevel tests the WithLoggingLevel function.
func TestWithLoggingLevel(t *testing.T) {
	product := &Product{}
	withLoggingLevel := WithLoggingLevel("error")
	withLoggingLevel(product)

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

	product := &Product{}
	withBizObjConfig := WithBizObjConfig(customConfig)
	withBizObjConfig(product)

	assert.Equal(t, "warning", _cfg.Logging.Level)
	assert.Equal(t, "customdb", _cfg.DataStores.DataStoreMap["mongo"].DatabaseName)
	assert.Equal(t, "customcollection", _cfg.DataStores.DataStoreMap["mongo"].CollectionName)
	assert.Equal(t, "mongodb://localhost:27017", _cfg.DataStores.DataStoreMap["mongo"].URI)
}
