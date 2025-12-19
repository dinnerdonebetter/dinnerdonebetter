package integration

import (
	"testing"

	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/pointer"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMealLists_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("should CRUD", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)

		createRes, err := userClient.CreateMealList(ctx, &mealplanninggrpc.CreateMealListRequest{
			Input: &mealplanninggrpc.MealListCreationRequestInput{
				Name:        t.Name(),
				Description: "desc",
			},
		})
		require.NoError(t, err)
		require.NotNil(t, createRes)

		listID := createRes.Created.Id

		listRes, err := userClient.GetMealLists(ctx, &mealplanninggrpc.GetMealListsRequest{})
		require.NoError(t, err)
		require.NotNil(t, listRes)
		assert.NotEmpty(t, listRes.Results)

		newName := t.Name() + "_updated"
		newDesc := "new desc"
		_, err = userClient.UpdateMealList(ctx, &mealplanninggrpc.UpdateMealListRequest{
			MealListId: listID,
			Input: &mealplanninggrpc.MealListUpdateRequestInput{
				Name:        &newName,
				Description: &newDesc,
			},
		})
		assert.NoError(t, err)

		_, err = userClient.ArchiveMealList(ctx, &mealplanninggrpc.ArchiveMealListRequest{MealListId: listID})
		assert.NoError(t, err)
	})
}

func TestMealListItems_CompleteLifecycle(T *testing.T) {
	T.Parallel()

	T.Run("should CRUD items", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)

		// Create a meal list
		createListRes, err := userClient.CreateMealList(ctx, &mealplanninggrpc.CreateMealListRequest{
			Input: &mealplanninggrpc.MealListCreationRequestInput{
				Name:        t.Name(),
				Description: "desc",
			},
		})
		require.NoError(t, err)
		listID := createListRes.Created.Id

		// Create a meal to reference
		createdMeal := createMealForTest(t, userClient, nil)

		// Create item
		createItemRes, err := userClient.CreateMealListItem(ctx, &mealplanninggrpc.CreateMealListItemRequest{
			Input: &mealplanninggrpc.MealListItemCreationRequestInput{
				BelongsToMealList: listID,
				MealId:            createdMeal.ID,
				Notes:             "notes",
			},
		})
		require.NoError(t, err)
		require.NotNil(t, createItemRes)

		// List lists and confirm item present
		listsRes, err := userClient.GetMealLists(ctx, &mealplanninggrpc.GetMealListsRequest{})
		require.NoError(t, err)
		require.NotNil(t, listsRes)
		require.NotEmpty(t, listsRes.Results)

		var itemID string
		for _, l := range listsRes.Results {
			if l.Id == listID && len(l.Items) > 0 {
				itemID = l.Items[0].Id
				assert.Equal(t, listID, l.Items[0].BelongsToMealList)
				break
			}
		}
		require.NotEmpty(t, itemID)

		// Update item
		newNotes := pointer.To("new notes")
		_, err = userClient.UpdateMealListItem(ctx, &mealplanninggrpc.UpdateMealListItemRequest{
			MealListItemId: itemID,
			Input: &mealplanninggrpc.MealListItemUpdateRequestInput{
				BelongsToMealList: &listID,
				MealId:            &createdMeal.ID,
				Notes:             newNotes,
			},
		})
		assert.NoError(t, err)

		// Archive item
		_, err = userClient.ArchiveMealListItem(ctx, &mealplanninggrpc.ArchiveMealListItemRequest{
			MealListItemId: itemID,
			MealListId:     listID,
		})
		assert.NoError(t, err)
	})
}
