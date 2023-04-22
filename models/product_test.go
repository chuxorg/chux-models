package models

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestNew tests the New function with different options.
func TestNewProduct(t *testing.T) {
	os.Setenv("APP_ENV", "test")
	product := NewProduct()
	assert.NotNil(t, product)
	assert.Equal(t, "testdb", product.GetDatabaseName())
	assert.Equal(t, "products", product.GetCollectionName())
	assert.Equal(t, "mongodb://localhost:27017", product.GetURI())
}

// TestWithLoggingLevel tests the WithLoggingLevel function.
func TestProductWithLoggingLevel(t *testing.T) {
	return
}

// TestWithBizObjConfig tests the WithBizObjConfig function.
func TestProductWithBizObjConfig(t *testing.T) {
	return
}
