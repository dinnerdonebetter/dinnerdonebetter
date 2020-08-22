package bleve

import (
	"context"
	"os"
	"testing"

	"gitlab.com/prixfixe/prixfixe/internal/v1/search"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gitlab.com/verygoodsoftwarenotvirus/logging/v1/noop"
)

type (
	exampleType struct {
		ID   uint64 `json:"id"`
		Name string `json:"name"`
	}

	exampleTypeWithStringID struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
)

func TestNewBleveIndexManager(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		exampleIndexPath := search.IndexPath("constructor_test_happy_path.bleve")

		_, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})

	T.Run("invalid path", func(t *testing.T) {
		exampleIndexPath := search.IndexPath("")

		_, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.Error(t, err)
	})

	T.Run("invalid name", func(t *testing.T) {
		exampleIndexPath := search.IndexPath("constructor_test_invalid_name.bleve")

		_, err := NewBleveIndexManager(exampleIndexPath, "invalid", noop.ProvideNoopLogger())
		assert.Error(t, err)
	})
}

func TestBleveIndexManager_Index(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "index_test"
		exampleIndexPath := search.IndexPath("index_test_obligatory.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleType{
			ID:   123,
			Name: exampleQuery,
		}
		assert.NoError(t, im.Index(ctx, x.ID, x))

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})
}

func TestBleveIndexManager_Search(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "search_test"
		exampleIndexPath := search.IndexPath("search_test_obligatory.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleType{
			ID:   123,
			Name: exampleQuery,
		}
		assert.NoError(t, im.Index(ctx, x.ID, x))

		results, err := im.Search(ctx, x.Name)
		assert.NotEmpty(t, results)
		assert.NoError(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})

	T.Run("with empty index and search", func(t *testing.T) {
		ctx := context.Background()

		exampleIndexPath := search.IndexPath("search_test_empty_index.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		results, err := im.Search(ctx, "")
		assert.Empty(t, results)
		assert.NoError(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})

	T.Run("with closed index", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "search_test"
		exampleIndexPath := search.IndexPath("search_test_closed_index.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleType{
			ID:   123,
			Name: exampleQuery,
		}
		assert.NoError(t, im.Index(ctx, x.ID, x))

		assert.NoError(t, im.(*bleveIndexManager).index.Close())

		results, err := im.Search(ctx, x.Name)
		assert.Empty(t, results)
		assert.Error(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})

	T.Run("with invalid ID", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "search_test"
		exampleIndexPath := search.IndexPath("search_test_invalid_id.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleTypeWithStringID{
			ID:   "whatever",
			Name: exampleQuery,
		}
		assert.NoError(t, im.(*bleveIndexManager).index.Index(x.ID, x))

		results, err := im.Search(ctx, x.Name)
		assert.Empty(t, results)
		assert.Error(t, err)

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})
}

func TestBleveIndexManager_Delete(T *testing.T) {
	T.Parallel()

	T.Run("obligatory", func(t *testing.T) {
		ctx := context.Background()

		const exampleQuery = "delete_test"
		exampleIndexPath := search.IndexPath("delete_test.bleve")

		im, err := NewBleveIndexManager(exampleIndexPath, testingSearchIndexName, noop.ProvideNoopLogger())
		assert.NoError(t, err)
		require.NotNil(t, im)

		x := &exampleType{
			ID:   123,
			Name: exampleQuery,
		}
		assert.NoError(t, im.Index(ctx, x.ID, x))

		assert.NoError(t, im.Delete(ctx, x.ID))

		assert.NoError(t, os.RemoveAll(string(exampleIndexPath)))
	})
}
