package models

import (
	"context"
	"net/http"

	"gitlab.com/prixfixe/prixfixe/internal/v1/search"
)

const (
	// ValidIngredientsSearchIndexName is the name of the index used to search through valid ingredients.
	ValidIngredientsSearchIndexName search.IndexName = "valid_ingredients"
)

type (
	// ValidIngredient represents a valid ingredient.
	ValidIngredient struct {
		ID                 uint64  `json:"id"`
		Name               string  `json:"name"`
		Variant            string  `json:"variant"`
		Description        string  `json:"description"`
		Warning            string  `json:"warning"`
		ContainsEgg        bool    `json:"containsEgg"`
		ContainsDairy      bool    `json:"containsDairy"`
		ContainsPeanut     bool    `json:"containsPeanut"`
		ContainsTreeNut    bool    `json:"containsTreeNut"`
		ContainsSoy        bool    `json:"containsSoy"`
		ContainsWheat      bool    `json:"containsWheat"`
		ContainsShellfish  bool    `json:"containsShellfish"`
		ContainsSesame     bool    `json:"containsSesame"`
		ContainsFish       bool    `json:"containsFish"`
		ContainsGluten     bool    `json:"containsGluten"`
		AnimalFlesh        bool    `json:"animalFlesh"`
		AnimalDerived      bool    `json:"animalDerived"`
		MeasurableByVolume bool    `json:"measurableByVolume"`
		Icon               string  `json:"icon"`
		CreatedOn          uint64  `json:"createdOn"`
		LastUpdatedOn      *uint64 `json:"lastUpdatedOn"`
		ArchivedOn         *uint64 `json:"archivedOn"`
	}

	// ValidIngredientList represents a list of valid ingredients.
	ValidIngredientList struct {
		Pagination
		ValidIngredients []ValidIngredient `json:"validIngredients"`
	}

	// ValidIngredientCreationInput represents what a user could set as input for creating valid ingredients.
	ValidIngredientCreationInput struct {
		Name               string `json:"name"`
		Variant            string `json:"variant"`
		Description        string `json:"description"`
		Warning            string `json:"warning"`
		ContainsEgg        bool   `json:"containsEgg"`
		ContainsDairy      bool   `json:"containsDairy"`
		ContainsPeanut     bool   `json:"containsPeanut"`
		ContainsTreeNut    bool   `json:"containsTreeNut"`
		ContainsSoy        bool   `json:"containsSoy"`
		ContainsWheat      bool   `json:"containsWheat"`
		ContainsShellfish  bool   `json:"containsShellfish"`
		ContainsSesame     bool   `json:"containsSesame"`
		ContainsFish       bool   `json:"containsFish"`
		ContainsGluten     bool   `json:"containsGluten"`
		AnimalFlesh        bool   `json:"animalFlesh"`
		AnimalDerived      bool   `json:"animalDerived"`
		MeasurableByVolume bool   `json:"measurableByVolume"`
		Icon               string `json:"icon"`
	}

	// ValidIngredientUpdateInput represents what a user could set as input for updating valid ingredients.
	ValidIngredientUpdateInput struct {
		Name               string `json:"name"`
		Variant            string `json:"variant"`
		Description        string `json:"description"`
		Warning            string `json:"warning"`
		ContainsEgg        bool   `json:"containsEgg"`
		ContainsDairy      bool   `json:"containsDairy"`
		ContainsPeanut     bool   `json:"containsPeanut"`
		ContainsTreeNut    bool   `json:"containsTreeNut"`
		ContainsSoy        bool   `json:"containsSoy"`
		ContainsWheat      bool   `json:"containsWheat"`
		ContainsShellfish  bool   `json:"containsShellfish"`
		ContainsSesame     bool   `json:"containsSesame"`
		ContainsFish       bool   `json:"containsFish"`
		ContainsGluten     bool   `json:"containsGluten"`
		AnimalFlesh        bool   `json:"animalFlesh"`
		AnimalDerived      bool   `json:"animalDerived"`
		MeasurableByVolume bool   `json:"measurableByVolume"`
		Icon               string `json:"icon"`
	}

	// ValidIngredientDataManager describes a structure capable of storing valid ingredients permanently.
	ValidIngredientDataManager interface {
		ValidIngredientExists(ctx context.Context, validIngredientID uint64) (bool, error)
		GetValidIngredient(ctx context.Context, validIngredientID uint64) (*ValidIngredient, error)
		GetAllValidIngredientsCount(ctx context.Context) (uint64, error)
		GetAllValidIngredients(ctx context.Context, resultChannel chan []ValidIngredient) error
		GetValidIngredients(ctx context.Context, filter *QueryFilter) (*ValidIngredientList, error)
		GetValidIngredientsWithIDs(ctx context.Context, limit uint8, ids []uint64) ([]ValidIngredient, error)
		CreateValidIngredient(ctx context.Context, input *ValidIngredientCreationInput) (*ValidIngredient, error)
		UpdateValidIngredient(ctx context.Context, updated *ValidIngredient) error
		ArchiveValidIngredient(ctx context.Context, validIngredientID uint64) error
	}

	// ValidIngredientDataServer describes a structure capable of serving traffic related to valid ingredients.
	ValidIngredientDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		SearchHandler(res http.ResponseWriter, req *http.Request)
		ListHandler(res http.ResponseWriter, req *http.Request)
		CreateHandler(res http.ResponseWriter, req *http.Request)
		ExistenceHandler(res http.ResponseWriter, req *http.Request)
		ReadHandler(res http.ResponseWriter, req *http.Request)
		UpdateHandler(res http.ResponseWriter, req *http.Request)
		ArchiveHandler(res http.ResponseWriter, req *http.Request)
	}
)

// Update merges an ValidIngredientInput with a valid ingredient.
func (x *ValidIngredient) Update(input *ValidIngredientUpdateInput) {
	if input.Name != "" && input.Name != x.Name {
		x.Name = input.Name
	}

	if input.Variant != "" && input.Variant != x.Variant {
		x.Variant = input.Variant
	}

	if input.Description != "" && input.Description != x.Description {
		x.Description = input.Description
	}

	if input.Warning != "" && input.Warning != x.Warning {
		x.Warning = input.Warning
	}

	if input.ContainsEgg != x.ContainsEgg {
		x.ContainsEgg = input.ContainsEgg
	}

	if input.ContainsDairy != x.ContainsDairy {
		x.ContainsDairy = input.ContainsDairy
	}

	if input.ContainsPeanut != x.ContainsPeanut {
		x.ContainsPeanut = input.ContainsPeanut
	}

	if input.ContainsTreeNut != x.ContainsTreeNut {
		x.ContainsTreeNut = input.ContainsTreeNut
	}

	if input.ContainsSoy != x.ContainsSoy {
		x.ContainsSoy = input.ContainsSoy
	}

	if input.ContainsWheat != x.ContainsWheat {
		x.ContainsWheat = input.ContainsWheat
	}

	if input.ContainsShellfish != x.ContainsShellfish {
		x.ContainsShellfish = input.ContainsShellfish
	}

	if input.ContainsSesame != x.ContainsSesame {
		x.ContainsSesame = input.ContainsSesame
	}

	if input.ContainsFish != x.ContainsFish {
		x.ContainsFish = input.ContainsFish
	}

	if input.ContainsGluten != x.ContainsGluten {
		x.ContainsGluten = input.ContainsGluten
	}

	if input.AnimalFlesh != x.AnimalFlesh {
		x.AnimalFlesh = input.AnimalFlesh
	}

	if input.AnimalDerived != x.AnimalDerived {
		x.AnimalDerived = input.AnimalDerived
	}

	if input.MeasurableByVolume != x.MeasurableByVolume {
		x.MeasurableByVolume = input.MeasurableByVolume
	}

	if input.Icon != "" && input.Icon != x.Icon {
		x.Icon = input.Icon
	}
}

// ToUpdateInput creates a ValidIngredientUpdateInput struct for a valid ingredient.
func (x *ValidIngredient) ToUpdateInput() *ValidIngredientUpdateInput {
	return &ValidIngredientUpdateInput{
		Name:               x.Name,
		Variant:            x.Variant,
		Description:        x.Description,
		Warning:            x.Warning,
		ContainsEgg:        x.ContainsEgg,
		ContainsDairy:      x.ContainsDairy,
		ContainsPeanut:     x.ContainsPeanut,
		ContainsTreeNut:    x.ContainsTreeNut,
		ContainsSoy:        x.ContainsSoy,
		ContainsWheat:      x.ContainsWheat,
		ContainsShellfish:  x.ContainsShellfish,
		ContainsSesame:     x.ContainsSesame,
		ContainsFish:       x.ContainsFish,
		ContainsGluten:     x.ContainsGluten,
		AnimalFlesh:        x.AnimalFlesh,
		AnimalDerived:      x.AnimalDerived,
		MeasurableByVolume: x.MeasurableByVolume,
		Icon:               x.Icon,
	}
}
