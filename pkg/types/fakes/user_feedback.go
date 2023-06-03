package fakes

import (
	"github.com/dinnerdonebetter/backend/pkg/types"
	"github.com/dinnerdonebetter/backend/pkg/types/converters"

	fake "github.com/brianvoe/gofakeit/v5"
)

// BuildFakeUserFeedback builds a faked valid ingredient.
func BuildFakeUserFeedback() *types.UserFeedback {
	var x *types.UserFeedback
	fake.Struct(&x)
	return x
}

// BuildFakeUserFeedbackList builds a faked UserFeedbackList.
func BuildFakeUserFeedbackList() *types.QueryFilteredResult[types.UserFeedback] {
	var examples []*types.UserFeedback
	for i := 0; i < exampleQuantity; i++ {
		examples = append(examples, BuildFakeUserFeedback())
	}

	return &types.QueryFilteredResult[types.UserFeedback]{
		Pagination: types.Pagination{
			Page:          1,
			Limit:         20,
			FilteredCount: exampleQuantity / 2,
			TotalCount:    exampleQuantity,
		},
		Data: examples,
	}
}

// BuildFakeUserFeedbackCreationRequestInput builds a faked UserFeedbackCreationRequestInput.
func BuildFakeUserFeedbackCreationRequestInput() *types.UserFeedbackCreationRequestInput {
	validIngredient := BuildFakeUserFeedback()
	return converters.ConvertUserFeedbackToUserFeedbackCreationRequestInput(validIngredient)
}
