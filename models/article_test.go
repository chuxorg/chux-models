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

// Testing the CRUD operations of the Article struct
func TestArticleCRUD(t *testing.T) {
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

	article, err := NewArticleWithCustomURI(mongoServer.URI())
	require.NoError(t, err)

	assert.NotNil(t, article, "NewArticle should return a non-nil article")
	assert.True(t, article.IsNew(), "NewArticle should return a article with isNew set to true")
	assert.False(t, article.IsDirty(), "NewArticle should return a article with isDirty set to false")
	assert.False(t, article.isDeleted, "NewArticle should return a article with isDeleted set to false")

	// Insert a new article
	article.Headline = "Test article"
	article.Author = "12345"
	article.Description = "Test Description"

	err = article.Save()
	require.NoError(t, err)
	require.False(t, article.IsNew())
	require.False(t, article.IsDirty())

	// Load the article by ID
	loadedArticleInterface, err := article.Load(article.ID.Hex())
	require.NoError(t, err)

	loadedArticle, ok := loadedArticleInterface.(*Article)
	require.True(t, ok)

	assert.Equal(t, article.Headline, loadedArticle.Headline)
	assert.Equal(t, article.Author, loadedArticle.Author)
	assert.Equal(t, article.Description, loadedArticle.Description)

	// Update the article
	loadedArticle.Headline = "Updated Test article"

	err = loadedArticle.Save()
	require.NoError(t, err)
	require.False(t, article.isDeleted)
	require.False(t, article.IsNew())
	require.False(t, article.IsDirty())

	// Load the updated article
	updatedArticleInterface, err := article.Load(loadedArticle.ID.Hex())
	require.NoError(t, err)
	require.False(t, article.isDeleted)
	require.False(t, article.IsNew())
	require.False(t, article.IsDirty())

	updatedArticle, ok := updatedArticleInterface.(*Article)
	require.True(t, ok)

	assert.Equal(t, "Updated Test article", updatedArticle.Headline)

	// Delete the article
	err = updatedArticle.Delete()
	require.NoError(t, err)
	assert.True(t, updatedArticle.isDeleted)

	err = updatedArticle.Save()
	require.NoError(t, err)
	require.False(t, article.isDeleted)
	require.True(t, article.IsNew())
	require.False(t, article.IsDirty())
	// Try to load the deleted article
	deletedArticleInterface, err := article.Load(updatedArticle.ID.Hex())
	require.Error(t, err)
	assert.Nil(t, deletedArticleInterface)
}
