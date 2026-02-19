package integration

import (
	"testing"

	"github.com/dinnerdonebetter/backend/internal/domain/comments"
	commentsgrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/comments"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/pkg/client"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func commentsServiceAddCommentToRecipe(t *testing.T, recipeID string, c client.Client, content string) *commentsgrpc.Comment {
	t.Helper()
	ctx := t.Context()

	if content == "" {
		content = "test comment via CommentsService"
	}

	res, err := c.CommentsService().AddCommentToRecipe(ctx, &commentsgrpc.AddCommentToRecipeRequest{
		RecipeId: recipeID,
		Input:    &commentsgrpc.CommentCreationRequestInput{Content: content},
	})
	require.NoError(t, err)
	require.NotNil(t, res)
	require.NotNil(t, res.Comment)

	return res.Comment
}

func TestCommentsService_AddCommentToRecipe(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)
		_, testClient := createUserAndClientForTest(t)

		comment := commentsServiceAddCommentToRecipe(t, createdRecipe.ID, testClient, "via CommentsService")

		assert.NotEmpty(t, comment.Id)
		assert.Equal(t, comments.CommentTargetTypeRecipes, comment.TargetType)
		assert.Equal(t, createdRecipe.ID, comment.ReferencedId)
		assert.Equal(t, "via CommentsService", comment.Content)

		listRes, err := testClient.CommentsService().GetCommentsForReference(ctx, &commentsgrpc.GetCommentsForReferenceRequest{
			TargetType:   comments.CommentTargetTypeRecipes,
			ReferencedId: createdRecipe.ID,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(listRes.Data), 1)

		_, _ = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: createdRecipe.ID})
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)
		c := buildUnauthenticatedGRPCClientForTest(t)

		res, err := c.CommentsService().AddCommentToRecipe(ctx, &commentsgrpc.AddCommentToRecipeRequest{
			RecipeId: createdRecipe.ID,
			Input:    &commentsgrpc.CommentCreationRequestInput{Content: "test"},
		})
		assert.Error(t, err)
		assert.Nil(t, res)

		_, _ = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: createdRecipe.ID})
	})
}

func TestCommentsService_AddCommentToMeal(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		createdMeal := createMealForTest(t, userClient, nil)

		res, err := userClient.CommentsService().AddCommentToMeal(ctx, &commentsgrpc.AddCommentToMealRequest{
			MealId: createdMeal.ID,
			Input:  &commentsgrpc.CommentCreationRequestInput{Content: "comment via CommentsService"},
		})
		require.NoError(t, err)
		require.NotNil(t, res.Comment)
		assert.Equal(t, comments.CommentTargetTypeMeals, res.Comment.TargetType)
		assert.Equal(t, createdMeal.ID, res.Comment.ReferencedId)

		listRes, err := userClient.CommentsService().GetCommentsForReference(ctx, &commentsgrpc.GetCommentsForReferenceRequest{
			TargetType:   comments.CommentTargetTypeMeals,
			ReferencedId: createdMeal.ID,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(listRes.Data), 1)

		_, _ = userClient.ArchiveMeal(ctx, &mealplanninggrpc.ArchiveMealRequest{MealId: createdMeal.ID})
	})
}

func TestCommentsService_AddCommentToMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		createdMealPlan := createMealPlanForTest(t, userClient, nil)

		res, err := userClient.CommentsService().AddCommentToMealPlan(ctx, &commentsgrpc.AddCommentToMealPlanRequest{
			MealPlanId: createdMealPlan.ID,
			Input:      &commentsgrpc.CommentCreationRequestInput{Content: "comment via CommentsService"},
		})
		require.NoError(t, err)
		require.NotNil(t, res.Comment)
		assert.Equal(t, comments.CommentTargetTypeMealPlans, res.Comment.TargetType)
		assert.Equal(t, createdMealPlan.ID, res.Comment.ReferencedId)

		listRes, err := userClient.CommentsService().GetCommentsForReference(ctx, &commentsgrpc.GetCommentsForReferenceRequest{
			TargetType:   comments.CommentTargetTypeMealPlans,
			ReferencedId: createdMealPlan.ID,
		})
		require.NoError(t, err)
		assert.GreaterOrEqual(t, len(listRes.Data), 1)

		_, _ = userClient.ArchiveMealPlan(ctx, &mealplanninggrpc.ArchiveMealPlanRequest{MealPlanId: createdMealPlan.ID})
	})
}

func TestCommentsService_CreateComment(T *testing.T) {
	T.Parallel()

	T.Run("recipe target", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)
		_, testClient := createUserAndClientForTest(t)

		res, err := testClient.CommentsService().CreateComment(ctx, &commentsgrpc.CreateCommentRequest{
			Input: &commentsgrpc.CommentCreationRequestInput{
				Content:      "created via CreateComment",
				TargetType:   comments.CommentTargetTypeRecipes,
				ReferencedId: createdRecipe.ID,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res.Comment)
		assert.Equal(t, comments.CommentTargetTypeRecipes, res.Comment.TargetType)
		assert.Equal(t, createdRecipe.ID, res.Comment.ReferencedId)
		assert.Equal(t, "created via CreateComment", res.Comment.Content)

		_, _ = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: createdRecipe.ID})
	})

	T.Run("meal plan target", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, userClient := createUserAndClientForTest(t)
		createdMealPlan := createMealPlanForTest(t, userClient, nil)

		res, err := userClient.CommentsService().CreateComment(ctx, &commentsgrpc.CreateCommentRequest{
			Input: &commentsgrpc.CommentCreationRequestInput{
				Content:      "CreateComment on meal plan",
				TargetType:   comments.CommentTargetTypeMealPlans,
				ReferencedId: createdMealPlan.ID,
			},
		})
		require.NoError(t, err)
		require.NotNil(t, res.Comment)
		assert.Equal(t, comments.CommentTargetTypeMealPlans, res.Comment.TargetType)
		assert.Equal(t, createdMealPlan.ID, res.Comment.ReferencedId)

		_, _ = userClient.ArchiveMealPlan(ctx, &mealplanninggrpc.ArchiveMealPlanRequest{MealPlanId: createdMealPlan.ID})
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)
		c := buildUnauthenticatedGRPCClientForTest(t)

		res, err := c.CommentsService().CreateComment(ctx, &commentsgrpc.CreateCommentRequest{
			Input: &commentsgrpc.CommentCreationRequestInput{
				Content:      "test",
				TargetType:   comments.CommentTargetTypeRecipes,
				ReferencedId: createdRecipe.ID,
			},
		})
		assert.Error(t, err)
		assert.Nil(t, res)

		_, _ = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: createdRecipe.ID})
	})
}

func TestCommentsService_GetCommentsForReference(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)
		_, testClient := createUserAndClientForTest(t)
		_ = commentsServiceAddCommentToRecipe(t, createdRecipe.ID, testClient, "")

		listRes, err := testClient.CommentsService().GetCommentsForReference(ctx, &commentsgrpc.GetCommentsForReferenceRequest{
			TargetType:   comments.CommentTargetTypeRecipes,
			ReferencedId: createdRecipe.ID,
		})
		require.NoError(t, err)
		require.NotNil(t, listRes)
		assert.GreaterOrEqual(t, len(listRes.Data), 1)

		_, _ = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: createdRecipe.ID})
	})

	T.Run("requires auth", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)
		_, testClient := createUserAndClientForTest(t)
		_ = commentsServiceAddCommentToRecipe(t, createdRecipe.ID, testClient, "")

		c := buildUnauthenticatedGRPCClientForTest(t)
		listRes, err := c.CommentsService().GetCommentsForReference(ctx, &commentsgrpc.GetCommentsForReferenceRequest{
			TargetType:   comments.CommentTargetTypeRecipes,
			ReferencedId: createdRecipe.ID,
		})
		assert.Error(t, err)
		assert.Nil(t, listRes)

		_, _ = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: createdRecipe.ID})
	})
}

func TestCommentsService_UpdateComment(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)
		_, testClient := createUserAndClientForTest(t)
		createdComment := commentsServiceAddCommentToRecipe(t, createdRecipe.ID, testClient, "original")

		_, err := testClient.CommentsService().UpdateComment(ctx, &commentsgrpc.UpdateCommentRequest{
			CommentId: createdComment.Id,
			Input:     &commentsgrpc.CommentUpdateRequestInput{Content: "updated via CommentsService"},
		})
		require.NoError(t, err)

		listRes, err := testClient.CommentsService().GetCommentsForReference(ctx, &commentsgrpc.GetCommentsForReferenceRequest{
			TargetType:   comments.CommentTargetTypeRecipes,
			ReferencedId: createdRecipe.ID,
		})
		require.NoError(t, err)
		for _, c := range listRes.Data {
			if c.Id == createdComment.Id {
				assert.Equal(t, "updated via CommentsService", c.Content)
				break
			}
		}

		_, _ = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: createdRecipe.ID})
	})
}

func TestCommentsService_ArchiveComment(T *testing.T) {
	T.Parallel()

	T.Run("happy path", func(t *testing.T) {
		t.Parallel()
		ctx := t.Context()

		_, _, createdRecipe := createRecipeForTest(t, nil)
		_, testClient := createUserAndClientForTest(t)
		createdComment := commentsServiceAddCommentToRecipe(t, createdRecipe.ID, testClient, "")

		_, err := testClient.CommentsService().ArchiveComment(ctx, &commentsgrpc.ArchiveCommentRequest{
			CommentId: createdComment.Id,
		})
		require.NoError(t, err)

		listRes, err := testClient.CommentsService().GetCommentsForReference(ctx, &commentsgrpc.GetCommentsForReferenceRequest{
			TargetType:   comments.CommentTargetTypeRecipes,
			ReferencedId: createdRecipe.ID,
		})
		require.NoError(t, err)
		for _, c := range listRes.Data {
			assert.NotEqual(t, createdComment.Id, c.Id, "archived comment should not appear")
		}

		_, _ = adminClient.ArchiveRecipe(ctx, &mealplanninggrpc.ArchiveRecipeRequest{RecipeId: createdRecipe.ID})
	})
}
