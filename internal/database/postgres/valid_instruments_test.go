package postgres

import (
	"context"
	"database/sql"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/prixfixeco/api_server/internal/database/postgres/generated"
	"github.com/prixfixeco/api_server/pkg/types"
	"github.com/prixfixeco/api_server/pkg/types/fakes"
	testutils "github.com/prixfixeco/api_server/tests/utils"
)

func TestQuerier_ValidInstrumentExists(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		c, db := buildTestClient(t)

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"ValidInstrumentExists",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return(true, nil)
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.ValidInstrumentExists(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)
		assert.True(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockGeneratedQuerier)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, _ := buildTestClient(t)

		actual, err := c.ValidInstrumentExists(ctx, "")
		assert.Error(t, err)
		assert.False(t, actual)
	})

	T.Run("with sql.ErrNoRows", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		c, db := buildTestClient(t)

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"ValidInstrumentExists",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return(false, sql.ErrNoRows)
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.ValidInstrumentExists(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockGeneratedQuerier)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		c, db := buildTestClient(t)

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"ValidInstrumentExists",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return(false, errors.New("blah"))
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.ValidInstrumentExists(ctx, exampleValidInstrument.ID)
		assert.Error(t, err)
		assert.False(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockGeneratedQuerier)
	})
}

func TestQuerier_GetValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		querierResponse := &generated.GetValidInstrumentRow{
			ID:          exampleValidInstrument.ID,
			Name:        exampleValidInstrument.Name,
			PluralName:  exampleValidInstrument.PluralName,
			Description: exampleValidInstrument.Description,
			IconPath:    exampleValidInstrument.IconPath,
			CreatedOn:   int64(exampleValidInstrument.CreatedOn),
		}

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return(querierResponse, nil)
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.GetValidInstrument(ctx, exampleValidInstrument.ID)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrument, actual)

		mock.AssertExpectationsForObjects(t, db, mockGeneratedQuerier)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.GetValidInstrument(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"GetValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return((*generated.GetValidInstrumentRow)(nil), errors.New("blah"))
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.GetValidInstrument(ctx, exampleValidInstrument.ID)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockGeneratedQuerier)
	})
}

func TestQuerier_GetRandomValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		querierResponse := &generated.GetRandomValidInstrumentRow{
			ID:          exampleValidInstrument.ID,
			Name:        exampleValidInstrument.Name,
			PluralName:  exampleValidInstrument.PluralName,
			Description: exampleValidInstrument.Description,
			IconPath:    exampleValidInstrument.IconPath,
			CreatedOn:   int64(exampleValidInstrument.CreatedOn),
		}

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"GetRandomValidInstrument",
			testutils.ContextMatcher,
		).Return(querierResponse, nil)
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.GetRandomValidInstrument(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrument, actual)

		mock.AssertExpectationsForObjects(t, db, mockGeneratedQuerier)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"GetRandomValidInstrument",
			testutils.ContextMatcher,
		).Return((*generated.GetRandomValidInstrumentRow)(nil), errors.New("blah"))
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.GetRandomValidInstrument(ctx)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, db, mockGeneratedQuerier)
	})
}

func TestQuerier_SearchForValidInstruments(T *testing.T) {
	T.Parallel()

	exampleQuery := "blah"

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstruments := fakes.BuildFakeValidInstrumentList()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		results := []*generated.SearchForValidInstrumentsRow{}
		for _, instrument := range exampleValidInstruments.ValidInstruments {
			results = append(results, &generated.SearchForValidInstrumentsRow{
				ID:          instrument.ID,
				Name:        instrument.Name,
				PluralName:  instrument.PluralName,
				Description: instrument.Description,
				IconPath:    instrument.IconPath,
				CreatedOn:   int64(instrument.CreatedOn),
			})
		}

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"SearchForValidInstruments",
			testutils.ContextMatcher,
			wrapQueryForILIKE(exampleQuery),
		).Return(results, nil)
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.SearchForValidInstruments(ctx, exampleQuery)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstruments.ValidInstruments, actual)

		mock.AssertExpectationsForObjects(t, mockGeneratedQuerier)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.SearchForValidInstruments(ctx, "")
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"SearchForValidInstruments",
			testutils.ContextMatcher,
			wrapQueryForILIKE(exampleQuery),
		).Return([]*generated.SearchForValidInstrumentsRow(nil), errors.New("blah"))
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.SearchForValidInstruments(ctx, exampleQuery)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, mockGeneratedQuerier)
	})
}

func TestQuerier_GetTotalValidInstrumentCount(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		exampleCount := uint64(123)

		c, _ := buildTestClient(t)

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"GetTotalValidInstrumentCount",
			testutils.ContextMatcher,
		).Return(int64(exampleCount), nil)
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.GetTotalValidInstrumentCount(ctx)
		assert.NoError(t, err)
		assert.Equal(t, exampleCount, actual)

		mock.AssertExpectationsForObjects(t, mockGeneratedQuerier)
	})

	T.Run("error executing query", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()

		c, db := buildTestClient(t)

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"GetTotalValidInstrumentCount",
			testutils.ContextMatcher,
		).Return(int64(0), errors.New("blah"))
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.GetTotalValidInstrumentCount(ctx)
		assert.Error(t, err)
		assert.Zero(t, actual)

		mock.AssertExpectationsForObjects(t, db)
	})
}

func TestQuerier_GetValidInstruments(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()
		exampleValidInstrumentList := fakes.BuildFakeValidInstrumentList()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		results := []*generated.GetValidInstrumentsRow{}
		for _, instrument := range exampleValidInstrumentList.ValidInstruments {
			results = append(results, &generated.GetValidInstrumentsRow{
				ID:            instrument.ID,
				Name:          instrument.Name,
				PluralName:    instrument.PluralName,
				Description:   instrument.Description,
				IconPath:      instrument.IconPath,
				CreatedOn:     int64(instrument.CreatedOn),
				FilteredCount: int64(exampleValidInstrumentList.FilteredCount),
				TotalCount:    int64(exampleValidInstrumentList.TotalCount),
			})
		}

		args := &generated.GetValidInstrumentsParams{
			CreatedAfter:  nullInt64ForUint64Field(filter.CreatedAfter),
			CreatedBefore: nullInt64ForUint64Field(filter.CreatedBefore),
			UpdatedAfter:  nullInt64ForUint64Field(filter.UpdatedAfter),
			UpdatedBefore: nullInt64ForUint64Field(filter.UpdatedBefore),
			Limit:         nullInt32ForUint8Field(filter.Limit),
		}

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			args,
		).Return(results, nil)
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.GetValidInstruments(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, mockGeneratedQuerier)
	})

	T.Run("with nil filter", func(t *testing.T) {
		t.Parallel()

		filter := (*types.QueryFilter)(nil)
		exampleValidInstrumentList := fakes.BuildFakeValidInstrumentList()
		exampleValidInstrumentList.Page = 0
		exampleValidInstrumentList.Limit = 0

		ctx := context.Background()
		c, _ := buildTestClient(t)

		f := types.DefaultQueryFilter()

		results := []*generated.GetValidInstrumentsRow{}
		for _, instrument := range exampleValidInstrumentList.ValidInstruments {
			results = append(results, &generated.GetValidInstrumentsRow{
				ID:            instrument.ID,
				Name:          instrument.Name,
				PluralName:    instrument.PluralName,
				Description:   instrument.Description,
				IconPath:      instrument.IconPath,
				CreatedOn:     int64(instrument.CreatedOn),
				FilteredCount: int64(exampleValidInstrumentList.FilteredCount),
				TotalCount:    int64(exampleValidInstrumentList.TotalCount),
			})
		}

		args := &generated.GetValidInstrumentsParams{
			CreatedAfter:  nullInt64ForUint64Field(f.CreatedAfter),
			CreatedBefore: nullInt64ForUint64Field(f.CreatedBefore),
			UpdatedAfter:  nullInt64ForUint64Field(f.UpdatedAfter),
			UpdatedBefore: nullInt64ForUint64Field(f.UpdatedBefore),
			Limit:         nullInt32ForUint8Field(f.Limit),
		}

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			args,
		).Return(results, nil)
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.GetValidInstruments(ctx, filter)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrumentList, actual)

		mock.AssertExpectationsForObjects(t, mockGeneratedQuerier)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		filter := types.DefaultQueryFilter()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		args := &generated.GetValidInstrumentsParams{
			CreatedAfter:  nullInt64ForUint64Field(filter.CreatedAfter),
			CreatedBefore: nullInt64ForUint64Field(filter.CreatedBefore),
			UpdatedAfter:  nullInt64ForUint64Field(filter.UpdatedAfter),
			UpdatedBefore: nullInt64ForUint64Field(filter.UpdatedBefore),
			Limit:         nullInt32ForUint8Field(filter.Limit),
		}

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"GetValidInstruments",
			testutils.ContextMatcher,
			args,
		).Return([]*generated.GetValidInstrumentsRow(nil), errors.New("blah"))
		c.generatedQuerier = mockGeneratedQuerier

		actual, err := c.GetValidInstruments(ctx, filter)
		assert.Error(t, err)
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, mockGeneratedQuerier)
	})
}

func TestQuerier_CreateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		exampleValidInstrument.ID = "1"
		exampleInput := fakes.BuildFakeValidInstrumentDatabaseCreationInputFromValidInstrument(exampleValidInstrument)

		ctx := context.Background()
		c, _ := buildTestClient(t)

		args := &generated.CreateValidInstrumentParams{
			ID:          exampleInput.ID,
			Name:        exampleInput.Name,
			PluralName:  exampleInput.PluralName,
			Description: exampleInput.Description,
			IconPath:    exampleInput.IconPath,
		}

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"CreateValidInstrument",
			testutils.ContextMatcher,
			args,
		).Return(nil)
		c.generatedQuerier = mockGeneratedQuerier

		c.timeFunc = func() uint64 {
			return exampleValidInstrument.CreatedOn
		}

		actual, err := c.CreateValidInstrument(ctx, exampleInput)
		assert.NoError(t, err)
		assert.Equal(t, exampleValidInstrument, actual)

		mock.AssertExpectationsForObjects(t, mockGeneratedQuerier)
	})

	T.Run("with invalid input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		actual, err := c.CreateValidInstrument(ctx, nil)
		assert.Error(t, err)
		assert.Nil(t, actual)
	})

	T.Run("with error executing query", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New(t.Name())
		exampleValidInstrument := fakes.BuildFakeValidInstrument()
		exampleInput := fakes.BuildFakeValidInstrumentDatabaseCreationInputFromValidInstrument(exampleValidInstrument)

		ctx := context.Background()
		c, _ := buildTestClient(t)

		args := &generated.CreateValidInstrumentParams{
			ID:          exampleInput.ID,
			Name:        exampleInput.Name,
			PluralName:  exampleInput.PluralName,
			Description: exampleInput.Description,
			IconPath:    exampleInput.IconPath,
		}

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"CreateValidInstrument",
			testutils.ContextMatcher,
			args,
		).Return(expectedErr)
		c.generatedQuerier = mockGeneratedQuerier

		c.timeFunc = func() uint64 {
			return exampleValidInstrument.CreatedOn
		}

		actual, err := c.CreateValidInstrument(ctx, exampleInput)
		assert.Error(t, err)
		assert.True(t, errors.Is(err, expectedErr))
		assert.Nil(t, actual)

		mock.AssertExpectationsForObjects(t, mockGeneratedQuerier)
	})
}

func TestQuerier_UpdateValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := &generated.UpdateValidInstrumentParams{
			ID:          exampleValidInstrument.ID,
			Name:        exampleValidInstrument.Name,
			PluralName:  exampleValidInstrument.PluralName,
			Description: exampleValidInstrument.Description,
			IconPath:    exampleValidInstrument.IconPath,
		}

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"UpdateValidInstrument",
			testutils.ContextMatcher,
			args,
		).Return(nil)
		c.generatedQuerier = mockGeneratedQuerier

		assert.NoError(t, c.UpdateValidInstrument(ctx, exampleValidInstrument))

		mock.AssertExpectationsForObjects(t, db, mockGeneratedQuerier)
	})

	T.Run("with nil input", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.UpdateValidInstrument(ctx, nil))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		args := &generated.UpdateValidInstrumentParams{
			ID:          exampleValidInstrument.ID,
			Name:        exampleValidInstrument.Name,
			PluralName:  exampleValidInstrument.PluralName,
			Description: exampleValidInstrument.Description,
			IconPath:    exampleValidInstrument.IconPath,
		}

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"UpdateValidInstrument",
			testutils.ContextMatcher,
			args,
		).Return(errors.New("blah"))
		c.generatedQuerier = mockGeneratedQuerier

		assert.Error(t, c.UpdateValidInstrument(ctx, exampleValidInstrument))

		mock.AssertExpectationsForObjects(t, db, mockGeneratedQuerier)
	})
}

func TestQuerier_ArchiveValidInstrument(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"ArchiveValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return(nil)
		c.generatedQuerier = mockGeneratedQuerier

		assert.NoError(t, c.ArchiveValidInstrument(ctx, exampleValidInstrument.ID))

		mock.AssertExpectationsForObjects(t, db, mockGeneratedQuerier)
	})

	T.Run("with invalid valid instrument ID", func(t *testing.T) {
		t.Parallel()

		ctx := context.Background()
		c, _ := buildTestClient(t)

		assert.Error(t, c.ArchiveValidInstrument(ctx, ""))
	})

	T.Run("with error writing to database", func(t *testing.T) {
		t.Parallel()

		exampleValidInstrument := fakes.BuildFakeValidInstrument()

		ctx := context.Background()
		c, db := buildTestClient(t)

		mockGeneratedQuerier := &mockQuerier{}
		mockGeneratedQuerier.On(
			"ArchiveValidInstrument",
			testutils.ContextMatcher,
			exampleValidInstrument.ID,
		).Return(errors.New("blah"))
		c.generatedQuerier = mockGeneratedQuerier

		assert.Error(t, c.ArchiveValidInstrument(ctx, exampleValidInstrument.ID))

		mock.AssertExpectationsForObjects(t, db, mockGeneratedQuerier)
	})
}
