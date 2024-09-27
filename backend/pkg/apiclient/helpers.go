package apiclient

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"strconv"

	"github.com/dinnerdonebetter/backend/internal/observability"
	"github.com/dinnerdonebetter/backend/pkg/apiclient/generated"
	"github.com/dinnerdonebetter/backend/pkg/types"

	"github.com/jinzhu/copier"
)

// errorFromResponse returns library errors according to a response's status code.
func errorFromResponse(res *http.Response) error {
	if res == nil {
		return ErrNilResponse
	}

	switch res.StatusCode {
	case http.StatusNotFound:
		return ErrNotFound
	case http.StatusBadRequest:
		return ErrInvalidRequestInput
	case http.StatusUnauthorized, http.StatusForbidden:
		return ErrUnauthorized
	case http.StatusInternalServerError:
		return ErrInternalServerError
	default:
		return nil
	}
}

// argIsNotPointer checks an argument and returns whether it is a pointer.
func argIsNotPointer(i any) (bool, error) {
	if i == nil || reflect.TypeOf(i).Kind() != reflect.Ptr {
		return true, ErrArgumentIsNotPointer
	}

	return false, nil
}

// argIsNotNil checks an argument and returns whether it is nil.
func argIsNotNil(i any) (bool, error) {
	if i == nil {
		return true, ErrNilInputProvided
	}

	return false, nil
}

// argIsNotPointerOrNil does what it says on the tin. This function is primarily useful for detecting
// if a destination value is valid before decoding an HTTP response, for instance.
func argIsNotPointerOrNil(i any) error {
	if nn, err := argIsNotNil(i); nn || err != nil {
		return err
	}

	if np, err := argIsNotPointer(i); np || err != nil {
		return err
	}

	return nil
}

// unmarshalBody takes an HTTP response and JSON decodes its body into a destination value. The error returned here
// should only ever be received in testing, and should never be encountered by an end-user.
func (c *Client) unmarshalBody(ctx context.Context, res *http.Response, dest any) error {
	_, span := c.tracer.StartSpan(ctx)
	defer span.End()

	logger := c.logger.WithResponse(res)

	if err := argIsNotPointerOrNil(dest); err != nil {
		return observability.PrepareAndLogError(err, logger, span, "nil marshal target")
	}

	bodyBytes, err := io.ReadAll(res.Body)
	if err != nil {
		return observability.PrepareAndLogError(err, logger, span, "unmarshalling error response")
	}

	if res.StatusCode >= http.StatusBadRequest {
		return fmt.Errorf("unexpected status: %d", res.StatusCode)
	}

	if err = c.encoder.Unmarshal(ctx, bodyBytes, &dest); err != nil {
		logger = logger.WithValue("raw_body", string(bodyBytes))
		return observability.PrepareAndLogError(err, logger, span, "unmarshalling response body")
	}

	return nil
}

func (c *Client) queryFilterCleaner(ctx context.Context, req *http.Request) error {
	_, span := c.tracer.StartSpan(ctx)
	defer span.End()

	oldQuery := req.URL.Query()

	newQuery := url.Values{}
	for key, values := range oldQuery {
		switch key {
		case types.QueryKeyLimit,
			types.QueryKeyPage:
			for _, value := range values {
				if number, err := strconv.ParseUint(value, 10, 64); err == nil && number > 0 {
					newQuery.Set(key, strconv.Itoa(int(number)))
				}
			}
		case types.QueryKeySearch,
			types.QueryKeySearchWithDatabase,
			types.QueryKeyCreatedBefore,
			types.QueryKeyCreatedAfter,
			types.QueryKeyUpdatedBefore,
			types.QueryKeyUpdatedAfter,
			types.QueryKeyIncludeArchived,
			types.QueryKeySortBy:
			for _, value := range values {
				if value != "" {
					newQuery.Set(key, value)
				}
			}
		}
	}

	req.URL.RawQuery = newQuery.Encode()
	return nil
}

func (c *Client) copyType(to, from any) {
	if err := copier.Copy(to, from); err != nil {
		panic(err)
	}

	// the above doesn't _always_ work how we'd like it, so we handle special cases here:
	switch t := to.(type) {
	case *generated.UpdateRecipeStepProductJSONRequestBody:
		if f, ok := from.(*types.RecipeStepProduct); ok {
			if f.MeasurementUnit != nil {
				t.MeasurementUnitID = &f.MeasurementUnit.ID
			}
		}
	case *generated.CreateRecipeJSONRequestBody:
		if f, ok := from.(*types.RecipeCreationRequestInput); ok {
			for i, step := range f.Steps {
				if step != nil {
					for j, cc := range step.CompletionConditions {
						if cc != nil {
							// fucking slice pointers in generated code, man.
							bodySteps := *t.Steps
							bodyCCs := bodySteps[i].CompletionConditions
							(*bodyCCs)[j].IngredientState = &cc.IngredientStateID

							for k, ingredientID := range cc.Ingredients {
								(*(*bodyCCs)[j].Ingredients)[k] = int(ingredientID)
							}
						}
					}
				}
			}
		}
	case *generated.CreateRecipeStepJSONRequestBody:
		if f, ok := from.(*types.RecipeStepCreationRequestInput); ok {
			for j, cc := range f.CompletionConditions {
				if cc != nil {
					// fucking slice pointers in generated code, man.
					bodyCCs := t.CompletionConditions
					(*bodyCCs)[j].IngredientState = &cc.IngredientStateID

					for k, ingredientID := range cc.Ingredients {
						(*(*bodyCCs)[j].Ingredients)[k] = int(ingredientID)
					}
				}
			}
		}
	case *generated.RecipeStepInstrumentUpdateRequestInput:
		if f, ok := from.(*types.RecipeStepInstrument); ok {
			if f.Instrument != nil {
				t.InstrumentID = &f.Instrument.ID
			}
		}
	case *generated.RecipeStepIngredientUpdateRequestInput:
		if f, ok := from.(*types.RecipeStepIngredient); ok {
			t.MeasurementUnitID = &f.MeasurementUnit.ID
			if f.Ingredient != nil {
				t.IngredientID = &f.Ingredient.ID
			}
		}
	case *generated.UpdateRecipeStepCompletionConditionJSONRequestBody:
		if f, ok := from.(*types.RecipeStepCompletionCondition); ok {
			t.IngredientState = &f.IngredientState.ID
		}
	case *generated.UpdateMealPlanOptionJSONRequestBody:
		if f, ok := from.(*types.MealPlanOption); ok {
			t.MealID = &f.Meal.ID
		}
	case *generated.UpdateHouseholdInstrumentOwnershipJSONRequestBody:
		if f, ok := from.(*types.HouseholdInstrumentOwnership); ok {
			t.ValidInstrumentID = &f.Instrument.ID
		}
	}
}
