package mealplanning

import "errors"

var (
	// ErrDuplicateMeal is returned when creating a meal that already exists (same name and components for the creator).
	ErrDuplicateMeal = errors.New("meal with same name and components already exists")
	// ErrDuplicateMealInList is returned when adding a meal to a list that already contains it.
	ErrDuplicateMealInList = errors.New("meal already exists in list")
	// ErrDuplicateMealPlanOption is returned when adding a meal as an option to an event that already has it.
	ErrDuplicateMealPlanOption = errors.New("meal already exists as option for this event")

	// ErrNoMatchingMeal is a sentinel returned when FindMealWithSameComponents finds no duplicate.
	// It is not an error; callers should treat it as "no match found" and proceed.
	ErrNoMatchingMeal = errors.New("no meal with matching components found")
)
