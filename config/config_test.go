package config_test

import (
	"os"
	"testing"

	"github.com/chuxorg/chux-models/config"
	"github.com/chuxorg/chux-models/models/products"
	"github.com/stretchr/testify/assert"
)

// TestNew tests the New function with different options.
func TestNew(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	_cfg, err := config.LoadConfig("test")
	assert.Nil(t, err)
	product := products.New()
	assert.NotNil(t, product)
	assert.Equal(t, "testdb", product.GetDatabaseName())
	assert.Equal(t, "products", product.GetCollectionName())
	assert.Equal(t, "mongodb://localhost:27017", product.GetURI())

	productWithLoggingLevel := products.New(
		products.WithLoggingLevel("debug"),
	)

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

	productWithCustomConfig := products.New(
		products.WithBizObjConfig(customConfig),
	)

	assert.NotNil(t, productWithCustomConfig)
	assert.Equal(t, "customdb", productWithCustomConfig.GetDatabaseName())
	assert.Equal(t, "products", productWithCustomConfig.GetCollectionName())
	assert.Equal(t, "mongodb://localhost:27017", productWithCustomConfig.GetURI())
}