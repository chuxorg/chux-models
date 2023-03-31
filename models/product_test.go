package models

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/benweissmann/memongo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func TestProductCRUD(t *testing.T) {
	// Create instance of in-mem mongo
	mongoServer, err := memongo.Start("4.0.5")
	require.NoError(t, err)
	defer mongoServer.Stop()

	// Set environment variable to use the in-memory MongoDB instance
	os.Setenv("MONGO_URI", mongoServer.URI())

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoServer.URI()))
	require.NoError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	err = client.Connect(ctx)
	require.NoError(t, err)

	defer func() {
		err = client.Disconnect(ctx)
		assert.NoError(t, err)
	}()

	product, err := NewProductWithCustomURI(mongoServer.URI())
	require.NoError(t, err)

	assert.NotNil(t, product, "NewProduct should return a non-nil Product")
	assert.True(t, product.IsNew(), "NewProduct should return a Product with isNew set to true")
	assert.False(t, product.IsDirty(), "NewProduct should return a Product with isDirty set to false")
	assert.False(t, product.isDeleted, "NewProduct should return a Product with isDeleted set to false")

	// Insert a new product
	product.Name = "Test Product"
	product.SKU = "12345"
	product.Brand = "Test Brand"

	err = product.Save()
	require.NoError(t, err)
	require.False(t, product.IsNew())
	require.False(t, product.IsDirty())

	// Load the product by ID
	loadedProductInterface, err := product.Load(product.ID.Hex())
	require.NoError(t, err)

	loadedProduct, ok := loadedProductInterface.(*Product)
	require.True(t, ok)

	assert.Equal(t, product.Name, loadedProduct.Name)
	assert.Equal(t, product.SKU, loadedProduct.SKU)
	assert.Equal(t, product.Brand, loadedProduct.Brand)

	// Update the product
	loadedProduct.Name = "Updated Test Product"

	err = loadedProduct.Save()
	require.NoError(t, err)
	require.False(t, product.isDeleted)
	require.False(t, product.IsNew())
	require.False(t, product.IsDirty())

	// Load the updated product
	updatedProductInterface, err := product.Load(loadedProduct.ID.Hex())
	require.NoError(t, err)
	require.False(t, product.isDeleted)
	require.False(t, product.IsNew())
	require.False(t, product.IsDirty())

	updatedProduct, ok := updatedProductInterface.(*Product)
	require.True(t, ok)

	assert.Equal(t, "Updated Test Product", updatedProduct.Name)

	// Delete the product
	err = updatedProduct.Delete()
	require.NoError(t, err)
	assert.True(t, updatedProduct.isDeleted)

	err = updatedProduct.Save()
	require.NoError(t, err)
	require.False(t, product.isDeleted)
	require.True(t, product.IsNew())
	require.False(t, product.IsDirty())
	// Try to load the deleted product
	deletedProductInterface, err := product.Load(updatedProduct.ID.Hex())
	require.Error(t, err)
	assert.Nil(t, deletedProductInterface)
}
