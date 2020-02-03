package models

import (
	"context"
	"net/http"
)

type (
	// Ingredient represents an ingredient
	Ingredient struct {
		ID                uint64  `json:"id"`
		Name              string  `json:"name"`
		Variant           string  `json:"variant"`
		Description       string  `json:"description"`
		Warning           string  `json:"warning"`
		ContainsEgg       bool    `json:"contains_egg"`
		ContainsDairy     bool    `json:"contains_dairy"`
		ContainsPeanut    bool    `json:"contains_peanut"`
		ContainsTreeNut   bool    `json:"contains_tree_nut"`
		ContainsSoy       bool    `json:"contains_soy"`
		ContainsWheat     bool    `json:"contains_wheat"`
		ContainsShellfish bool    `json:"contains_shellfish"`
		ContainsSesame    bool    `json:"contains_sesame"`
		ContainsFish      bool    `json:"contains_fish"`
		ContainsGluten    bool    `json:"contains_gluten"`
		AnimalFlesh       bool    `json:"animal_flesh"`
		AnimalDerived     bool    `json:"animal_derived"`
		ConsideredStaple  bool    `json:"considered_staple"`
		Icon              string  `json:"icon"`
		CreatedOn         uint64  `json:"created_on"`
		UpdatedOn         *uint64 `json:"updated_on"`
		ArchivedOn        *uint64 `json:"archived_on"`
		BelongsTo         uint64  `json:"belongs_to"`
	}

	// IngredientList represents a list of ingredients
	IngredientList struct {
		Pagination
		Ingredients []Ingredient `json:"ingredients"`
	}

	// IngredientCreationInput represents what a user could set as input for creating ingredients
	IngredientCreationInput struct {
		Name              string `json:"name"`
		Variant           string `json:"variant"`
		Description       string `json:"description"`
		Warning           string `json:"warning"`
		ContainsEgg       bool   `json:"contains_egg"`
		ContainsDairy     bool   `json:"contains_dairy"`
		ContainsPeanut    bool   `json:"contains_peanut"`
		ContainsTreeNut   bool   `json:"contains_tree_nut"`
		ContainsSoy       bool   `json:"contains_soy"`
		ContainsWheat     bool   `json:"contains_wheat"`
		ContainsShellfish bool   `json:"contains_shellfish"`
		ContainsSesame    bool   `json:"contains_sesame"`
		ContainsFish      bool   `json:"contains_fish"`
		ContainsGluten    bool   `json:"contains_gluten"`
		AnimalFlesh       bool   `json:"animal_flesh"`
		AnimalDerived     bool   `json:"animal_derived"`
		ConsideredStaple  bool   `json:"considered_staple"`
		Icon              string `json:"icon"`
		BelongsTo         uint64 `json:"-"`
	}

	// IngredientUpdateInput represents what a user could set as input for updating ingredients
	IngredientUpdateInput struct {
		Name              string `json:"name"`
		Variant           string `json:"variant"`
		Description       string `json:"description"`
		Warning           string `json:"warning"`
		ContainsEgg       bool   `json:"contains_egg"`
		ContainsDairy     bool   `json:"contains_dairy"`
		ContainsPeanut    bool   `json:"contains_peanut"`
		ContainsTreeNut   bool   `json:"contains_tree_nut"`
		ContainsSoy       bool   `json:"contains_soy"`
		ContainsWheat     bool   `json:"contains_wheat"`
		ContainsShellfish bool   `json:"contains_shellfish"`
		ContainsSesame    bool   `json:"contains_sesame"`
		ContainsFish      bool   `json:"contains_fish"`
		ContainsGluten    bool   `json:"contains_gluten"`
		AnimalFlesh       bool   `json:"animal_flesh"`
		AnimalDerived     bool   `json:"animal_derived"`
		ConsideredStaple  bool   `json:"considered_staple"`
		Icon              string `json:"icon"`
		BelongsTo         uint64 `json:"-"`
	}

	// IngredientDataManager describes a structure capable of storing ingredients permanently
	IngredientDataManager interface {
		GetIngredient(ctx context.Context, ingredientID, userID uint64) (*Ingredient, error)
		GetIngredientCount(ctx context.Context, filter *QueryFilter, userID uint64) (uint64, error)
		GetAllIngredientsCount(ctx context.Context) (uint64, error)
		GetIngredients(ctx context.Context, filter *QueryFilter, userID uint64) (*IngredientList, error)
		GetAllIngredientsForUser(ctx context.Context, userID uint64) ([]Ingredient, error)
		CreateIngredient(ctx context.Context, input *IngredientCreationInput) (*Ingredient, error)
		UpdateIngredient(ctx context.Context, updated *Ingredient) error
		ArchiveIngredient(ctx context.Context, id, userID uint64) error
	}

	// IngredientDataServer describes a structure capable of serving traffic related to ingredients
	IngredientDataServer interface {
		CreationInputMiddleware(next http.Handler) http.Handler
		UpdateInputMiddleware(next http.Handler) http.Handler

		ListHandler() http.HandlerFunc
		CreateHandler() http.HandlerFunc
		ReadHandler() http.HandlerFunc
		UpdateHandler() http.HandlerFunc
		ArchiveHandler() http.HandlerFunc
	}
)

// Update merges an IngredientInput with an ingredient
func (x *Ingredient) Update(input *IngredientUpdateInput) {
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

	if input.ConsideredStaple != x.ConsideredStaple {
		x.ConsideredStaple = input.ConsideredStaple
	}

	if input.Icon != "" && input.Icon != x.Icon {
		x.Icon = input.Icon
	}
}

// ToInput creates a IngredientUpdateInput struct for an ingredient
func (x *Ingredient) ToInput() *IngredientUpdateInput {
	return &IngredientUpdateInput{
		Name:              x.Name,
		Variant:           x.Variant,
		Description:       x.Description,
		Warning:           x.Warning,
		ContainsEgg:       x.ContainsEgg,
		ContainsDairy:     x.ContainsDairy,
		ContainsPeanut:    x.ContainsPeanut,
		ContainsTreeNut:   x.ContainsTreeNut,
		ContainsSoy:       x.ContainsSoy,
		ContainsWheat:     x.ContainsWheat,
		ContainsShellfish: x.ContainsShellfish,
		ContainsSesame:    x.ContainsSesame,
		ContainsFish:      x.ContainsFish,
		ContainsGluten:    x.ContainsGluten,
		AnimalFlesh:       x.AnimalFlesh,
		AnimalDerived:     x.AnimalDerived,
		ConsideredStaple:  x.ConsideredStaple,
		Icon:              x.Icon,
	}
}
