package textsearch

import (
	"context"
)

const (
	// IndexTypeRecipes represents the recipes index.
	IndexTypeRecipes = "recipes"
	// IndexTypeMeals represents the meals index.
	IndexTypeMeals = "meals"
	// IndexTypeValidIngredients represents the valid_ingredients index.
	IndexTypeValidIngredients = "valid_ingredients"
	// IndexTypeValidInstruments represents the valid_instruments index.
	IndexTypeValidInstruments = "valid_instruments"
	// IndexTypeValidMeasurementUnits represents the valid_measurement_units index.
	IndexTypeValidMeasurementUnits = "valid_measurement_units"
	// IndexTypeValidPreparations represents the  valid_preparations index.
	IndexTypeValidPreparations = "valid_preparations"
	// IndexTypeValidIngredientStates represents the valid_ingredient_states index.
	IndexTypeValidIngredientStates = "valid_ingredient_states"
	// IndexTypeValidVessels represents the vessels index.
	IndexTypeValidVessels = "vessels"
	// IndexTypeUsers represents the users index.
	IndexTypeUsers = "users"
)

type ( // IndexSearcher is our wrapper interface for querying a text search index.
	IndexSearcher[T any] interface {
		Search(ctx context.Context, query string) (ids []*T, err error)
	}

	// IndexManager is our wrapper interface for a text search index.
	IndexManager interface {
		Index(ctx context.Context, id string, value any) error
		Delete(ctx context.Context, id string) (err error)
		Wipe(ctx context.Context) error
	}

	// Index is our wrapper interface for a text search index.
	Index[T any] interface {
		IndexSearcher[T]
		IndexManager
	}
)
