package grpc

import (
	"context"
	"testing"

	"github.com/dinnerdonebetter/backend/internal/authentication/sessions"
	"github.com/dinnerdonebetter/backend/internal/domain/comments"
	commentsfakes "github.com/dinnerdonebetter/backend/internal/domain/comments/fakes"
	commentsmanagermock "github.com/dinnerdonebetter/backend/internal/domain/comments/manager/mock"
	mealplanningfakes "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/fakes"
	mockmanagers "github.com/dinnerdonebetter/backend/internal/domain/mealplanning/managers/mock"
	mealplanninggrpc "github.com/dinnerdonebetter/backend/internal/grpc/generated/services/mealplanning"
	"github.com/dinnerdonebetter/backend/internal/platform/reflection"
	"github.com/dinnerdonebetter/backend/internal/platform/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestServiceImpl_AddCommentToRecipe(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		recipeID := mealplanningfakes.BuildFakeID()
		userID := mealplanningfakes.BuildFakeID()
		content := "test comment"

		mcm := &commentsmanagermock.MockCommentsDataManager{}
		fakeComment := commentsfakes.BuildFakeComment()
		fakeComment.TargetType = comments.CommentTargetTypeRecipes
		fakeComment.ReferencedID = recipeID

		mcm.On(reflection.GetMethodName(mcm.CreateComment), testutils.ContextMatcher, mock.MatchedBy(func(in interface{}) bool {
			ci, ok := in.(*comments.CommentCreationRequestInput)
			return ok && ci != nil && ci.TargetType == comments.CommentTargetTypeRecipes && ci.ReferencedID == recipeID && ci.BelongsToUser == userID
		})).Return(fakeComment, nil)
		s.commentsManager = mcm

		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{UserID: userID},
			}, nil
		}

		res, err := s.AddCommentToRecipe(ctx, &mealplanninggrpc.AddCommentToRecipeRequest{
			RecipeId: recipeID,
			Content:  content,
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, fakeComment.ID, res.Comment.Id)
		assert.Equal(t, fakeComment.Content, res.Comment.Content)

		mock.AssertExpectationsForObjects(t, mcm)
	})
}

func TestServiceImpl_AddCommentToMeal(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mealID := mealplanningfakes.BuildFakeID()
		userID := mealplanningfakes.BuildFakeID()
		content := "test comment on meal"

		mcm := &commentsmanagermock.MockCommentsDataManager{}
		fakeComment := commentsfakes.BuildFakeComment()
		fakeComment.TargetType = comments.CommentTargetTypeMeals
		fakeComment.ReferencedID = mealID
		mcm.On(reflection.GetMethodName(mcm.CreateComment), testutils.ContextMatcher, mock.MatchedBy(func(in interface{}) bool {
			ci, ok := in.(*comments.CommentCreationRequestInput)
			return ok && ci != nil && ci.TargetType == comments.CommentTargetTypeMeals && ci.ReferencedID == mealID && ci.BelongsToUser == userID
		})).Return(fakeComment, nil)
		s.commentsManager = mcm

		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{UserID: userID},
			}, nil
		}

		res, err := s.AddCommentToMeal(ctx, &mealplanninggrpc.AddCommentToMealRequest{
			MealId:  mealID,
			Content: content,
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, fakeComment.ID, res.Comment.Id)

		mock.AssertExpectationsForObjects(t, mcm)
	})
}

func TestServiceImpl_AddCommentToMealPlan(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mealPlanID := mealplanningfakes.BuildFakeID()
		accountID := mealplanningfakes.BuildFakeID()
		userID := mealplanningfakes.BuildFakeID()
		content := "test comment on meal plan"

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On(reflection.GetMethodName(mmpm.ReadMealPlan), testutils.ContextMatcher, mealPlanID, accountID).Return(mealplanningfakes.BuildFakeMealPlan(), nil)
		s.mealPlanningManager = mmpm

		mcm := &commentsmanagermock.MockCommentsDataManager{}
		fakeComment := commentsfakes.BuildFakeComment()
		fakeComment.TargetType = comments.CommentTargetTypeMealPlans
		fakeComment.ReferencedID = mealPlanID
		mcm.On(reflection.GetMethodName(mcm.CreateComment), testutils.ContextMatcher, mock.Anything).Return(fakeComment, nil)
		s.commentsManager = mcm

		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester:       sessions.RequesterInfo{UserID: userID},
				ActiveAccountID: accountID,
			}, nil
		}

		res, err := s.AddCommentToMealPlan(ctx, &mealplanninggrpc.AddCommentToMealPlanRequest{
			MealPlanId: mealPlanID,
			Content:    content,
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, fakeComment.ID, res.Comment.Id)

		mock.AssertExpectationsForObjects(t, mmpm, mcm)
	})
}

func TestServiceImpl_GetCommentsForReference(T *testing.T) {
	T.Parallel()

	T.Run("recipes", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		recipeID := mealplanningfakes.BuildFakeID()
		expected := commentsfakes.BuildFakeCommentList(comments.CommentTargetTypeRecipes, recipeID)

		mcm := &commentsmanagermock.MockCommentsDataManager{}
		mcm.On(reflection.GetMethodName(mcm.GetCommentsForReference), testutils.ContextMatcher, comments.CommentTargetTypeRecipes, recipeID, testutils.QueryFilterMatcher).Return(expected, nil)
		s.commentsManager = mcm

		res, err := s.GetCommentsForReference(ctx, &mealplanninggrpc.GetCommentsForReferenceRequest{
			TargetType:   comments.CommentTargetTypeRecipes,
			ReferencedId: recipeID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Data, len(expected.Data))

		mock.AssertExpectationsForObjects(t, mcm)
	})

	T.Run("meals", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mealID := mealplanningfakes.BuildFakeID()
		expected := commentsfakes.BuildFakeCommentList(comments.CommentTargetTypeMeals, mealID)

		mcm := &commentsmanagermock.MockCommentsDataManager{}
		mcm.On(reflection.GetMethodName(mcm.GetCommentsForReference), testutils.ContextMatcher, comments.CommentTargetTypeMeals, mealID, testutils.QueryFilterMatcher).Return(expected, nil)
		s.commentsManager = mcm

		res, err := s.GetCommentsForReference(ctx, &mealplanninggrpc.GetCommentsForReferenceRequest{
			TargetType:   comments.CommentTargetTypeMeals,
			ReferencedId: mealID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Data, len(expected.Data))

		mock.AssertExpectationsForObjects(t, mcm)
	})

	T.Run("meal_plans", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForMealPlanningTest(t)

		mealPlanID := mealplanningfakes.BuildFakeID()
		accountID := mealplanningfakes.BuildFakeID()
		expected := commentsfakes.BuildFakeCommentList(comments.CommentTargetTypeMealPlans, mealPlanID)

		mmpm := &mockmanagers.MockMealPlanningManager{}
		mmpm.On(reflection.GetMethodName(mmpm.ReadMealPlan), testutils.ContextMatcher, mealPlanID, accountID).Return(mealplanningfakes.BuildFakeMealPlan(), nil)
		s.mealPlanningManager = mmpm

		mcm := &commentsmanagermock.MockCommentsDataManager{}
		mcm.On(reflection.GetMethodName(mcm.GetCommentsForReference), testutils.ContextMatcher, comments.CommentTargetTypeMealPlans, mealPlanID, testutils.QueryFilterMatcher).Return(expected, nil)
		s.commentsManager = mcm

		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{ActiveAccountID: accountID}, nil
		}

		res, err := s.GetCommentsForReference(ctx, &mealplanninggrpc.GetCommentsForReferenceRequest{
			TargetType:   comments.CommentTargetTypeMealPlans,
			ReferencedId: mealPlanID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Len(t, res.Data, len(expected.Data))

		mock.AssertExpectationsForObjects(t, mmpm, mcm)
	})

	T.Run("invalid_target_type", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		res, err := s.GetCommentsForReference(ctx, &mealplanninggrpc.GetCommentsForReferenceRequest{
			TargetType:   "invalid",
			ReferencedId: mealplanningfakes.BuildFakeID(),
		})
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.InvalidArgument, st.Code())
	})
}

func TestServiceImpl_UpdateComment(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		commentID := commentsfakes.BuildFakeID()
		userID := mealplanningfakes.BuildFakeID()
		newContent := "updated content"

		fakeComment := commentsfakes.BuildFakeComment()
		fakeComment.ID = commentID
		fakeComment.BelongsToUser = userID

		mcm := &commentsmanagermock.MockCommentsDataManager{}
		mcm.On(reflection.GetMethodName(mcm.GetComment), testutils.ContextMatcher, commentID).Return(fakeComment, nil).Twice()
		mcm.On(reflection.GetMethodName(mcm.UpdateComment), testutils.ContextMatcher, commentID, userID, mock.Anything).Return(nil)
		s.commentsManager = mcm

		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{UserID: userID},
			}, nil
		}

		res, err := s.UpdateComment(ctx, &mealplanninggrpc.UpdateCommentRequest{
			CommentId: commentID,
			Content:   newContent,
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)
		assert.Equal(t, commentID, res.Comment.Id)

		mock.AssertExpectationsForObjects(t, mcm)
	})

	T.Run("permission_denied_when_different_user", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		commentID := commentsfakes.BuildFakeID()
		ownerID := commentsfakes.BuildFakeID()
		requestingUserID := mealplanningfakes.BuildFakeID()

		fakeComment := commentsfakes.BuildFakeComment()
		fakeComment.ID = commentID
		fakeComment.BelongsToUser = ownerID

		mcm := &commentsmanagermock.MockCommentsDataManager{}
		mcm.On(reflection.GetMethodName(mcm.GetComment), testutils.ContextMatcher, commentID).Return(fakeComment, nil)
		s.commentsManager = mcm

		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{UserID: requestingUserID},
			}, nil
		}

		res, err := s.UpdateComment(ctx, &mealplanninggrpc.UpdateCommentRequest{
			CommentId: commentID,
			Content:   "updated",
		})
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.PermissionDenied, st.Code())

		mock.AssertExpectationsForObjects(t, mcm)
	})
}

func TestServiceImpl_ArchiveComment(T *testing.T) {
	T.Parallel()

	T.Run("standard", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		commentID := commentsfakes.BuildFakeID()
		userID := mealplanningfakes.BuildFakeID()

		fakeComment := commentsfakes.BuildFakeComment()
		fakeComment.ID = commentID
		fakeComment.BelongsToUser = userID

		mcm := &commentsmanagermock.MockCommentsDataManager{}
		mcm.On(reflection.GetMethodName(mcm.GetComment), testutils.ContextMatcher, commentID).Return(fakeComment, nil)
		mcm.On(reflection.GetMethodName(mcm.ArchiveComment), testutils.ContextMatcher, commentID).Return(nil)
		s.commentsManager = mcm

		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{UserID: userID},
			}, nil
		}

		res, err := s.ArchiveComment(ctx, &mealplanninggrpc.ArchiveCommentRequest{
			CommentId: commentID,
		})
		assert.NoError(t, err)
		assert.NotNil(t, res)

		mock.AssertExpectationsForObjects(t, mcm)
	})

	T.Run("permission_denied_when_different_user", func(t *testing.T) {
		t.Parallel()

		ctx := t.Context()
		s := buildServiceImplForRecipesTest(t)

		commentID := commentsfakes.BuildFakeID()
		ownerID := commentsfakes.BuildFakeID()
		requestingUserID := mealplanningfakes.BuildFakeID()

		fakeComment := commentsfakes.BuildFakeComment()
		fakeComment.ID = commentID
		fakeComment.BelongsToUser = ownerID

		mcm := &commentsmanagermock.MockCommentsDataManager{}
		mcm.On(reflection.GetMethodName(mcm.GetComment), testutils.ContextMatcher, commentID).Return(fakeComment, nil)
		s.commentsManager = mcm

		s.sessionContextDataFetcher = func(ctx context.Context) (*sessions.ContextData, error) {
			return &sessions.ContextData{
				Requester: sessions.RequesterInfo{UserID: requestingUserID},
			}, nil
		}

		res, err := s.ArchiveComment(ctx, &mealplanninggrpc.ArchiveCommentRequest{
			CommentId: commentID,
		})
		assert.Error(t, err)
		assert.Nil(t, res)
		st, ok := status.FromError(err)
		assert.True(t, ok)
		assert.Equal(t, codes.PermissionDenied, st.Code())

		mock.AssertExpectationsForObjects(t, mcm)
	})
}
