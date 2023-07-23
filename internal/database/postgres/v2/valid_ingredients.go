package v2

import (
	"context"
	"database/sql"
	"time"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/doug-martin/goqu/v9"
	"github.com/jinzhu/copier"
)

const (
	validIngredientsTableName = "valid_ingredients"
)

type (
	// ValidIngredient represents a valid ingredient.
	ValidIngredient struct {
		_ struct{}

		CreatedAt                               time.Time  `db:"created_at"                                   goqu:"skipinsert"`
		LastUpdatedAt                           *time.Time `db:"last_updated_at"                              goqu:"skipinsert"`
		ArchivedAt                              *time.Time `db:"archived_at"                                  goqu:"skipinsert"`
		MaximumIdealStorageTemperatureInCelsius *float32   `db:"maximum_ideal_storage_temperature_in_celsius"`
		MinimumIdealStorageTemperatureInCelsius *float32   `db:"minimum_ideal_storage_temperature_in_celsius"`
		IconPath                                string     `db:"icon_path"`
		Warning                                 string     `db:"warning"`
		PluralName                              string     `db:"plural_name"`
		StorageInstructions                     string     `db:"storage_instructions"`
		Name                                    string     `db:"name"`
		ID                                      string     `db:"id"`
		Description                             string     `db:"description"`
		Slug                                    string     `db:"slug"`
		ShoppingSuggestions                     string     `db:"shopping_suggestions"`
		ContainsShellfish                       bool       `db:"contains_shellfish"`
		IsMeasuredVolumetrically                bool       `db:"volumetric"`
		IsLiquid                                bool       `db:"is_liquid"`
		ContainsPeanut                          bool       `db:"contains_peanut"`
		ContainsTreeNut                         bool       `db:"contains_tree_nut"`
		ContainsEgg                             bool       `db:"contains_egg"`
		ContainsWheat                           bool       `db:"contains_wheat"`
		ContainsSoy                             bool       `db:"contains_soy"`
		AnimalDerived                           bool       `db:"animal_derived"`
		RestrictToPreparations                  bool       `db:"restrict_to_preparations"`
		ContainsSesame                          bool       `db:"contains_sesame"`
		ContainsFish                            bool       `db:"contains_fish"`
		ContainsGluten                          bool       `db:"contains_gluten"`
		ContainsDairy                           bool       `db:"contains_dairy"`
		ContainsAlcohol                         bool       `db:"contains_alcohol"`
		AnimalFlesh                             bool       `db:"animal_flesh"`
		IsStarch                                bool       `db:"is_starch"`
		IsProtein                               bool       `db:"is_protein"`
		IsGrain                                 bool       `db:"is_grain"`
		IsFruit                                 bool       `db:"is_fruit"`
		IsSalt                                  bool       `db:"is_salt"`
		IsFat                                   bool       `db:"is_fat"`
		IsAcid                                  bool       `db:"is_acid"`
		IsHeat                                  bool       `db:"is_heat"`
	}
)

// CreateValidIngredient gets a valid ingredient from the database.
func (c *DatabaseClient) CreateValidIngredient(ctx context.Context, input *ValidIngredient) (*types.ValidIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	q := c.xdb.Insert(validIngredientsTableName).Rows(
		input,
	)

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return nil, observability.PrepareError(err, span, "creating valid ingredient")
	}

	var output types.ValidIngredient
	if err := copier.Copy(&output, input); err != nil {
		return nil, observability.PrepareError(err, span, "copying input to output")
	}

	return &output, nil
}

// GetValidIngredient gets a valid ingredient from the database.
func (c *DatabaseClient) GetValidIngredient(ctx context.Context, validIngredientID string) (*types.ValidIngredient, error) {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	x := &ValidIngredient{}
	q := c.xdb.From(validIngredientsTableName).Where(goqu.Ex{
		idColumn:         validIngredientID,
		archivedAtColumn: nil,
	})

	found, err := q.ScanStructContext(ctx, x)
	if err != nil {
		return nil, err
	} else if !found {
		return nil, sql.ErrNoRows
	}

	var output types.ValidIngredient
	if err = copier.Copy(&output, x); err != nil {
		return nil, observability.PrepareError(err, span, "copying input to output")
	}

	return &output, nil
}

// UpdateValidIngredient gets a valid ingredient from the database.
func (c *DatabaseClient) UpdateValidIngredient(ctx context.Context, input *types.ValidIngredient) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	var updateInput ValidIngredient
	if err := copier.Copy(&updateInput, input); err != nil {
		return observability.PrepareError(err, span, "copying input to output")
	}

	q := c.xdb.Update(validIngredientsTableName).Set(
		updateInput,
	)

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return observability.PrepareError(err, span, "updating valid ingredient")
	}

	return nil
}

// ArchiveValidIngredient gets a valid ingredient from the database.
func (c *DatabaseClient) ArchiveValidIngredient(ctx context.Context, validIngredientID string) error {
	ctx, span := c.tracer.StartSpan(ctx)
	defer span.End()

	q := c.xdb.Update(validIngredientsTableName).
		Set(goqu.Record{archivedAtColumn: goqu.L("NOW()")}).
		Where(goqu.Ex{idColumn: validIngredientID})

	if _, err := q.Executor().ExecContext(ctx); err != nil {
		return observability.PrepareError(err, span, "archiving valid ingredient")
	}

	return nil
}
