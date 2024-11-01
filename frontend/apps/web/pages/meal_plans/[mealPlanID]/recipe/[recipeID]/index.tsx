import { Text } from '@mantine/core';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { useState } from 'react';

import {
  APIResponse,
  EitherErrorOr,
  Household,
  IAPIError,
  MealPlan,
  MealPlanGroceryListItem,
  MealPlanTask,
  QueryFilteredResult,
  Recipe,
} from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildServerSideClientOrRedirect } from '../../../../../src/client';
import { AppLayout } from '../../../../../src/layouts';
import { RecipeComponent } from '../../../../../src/components';
import { serverSideTracer } from '../../../../../src/tracer';
import { serverSideAnalytics } from '../../../../../src/analytics';
import { userSessionDetailsOrRedirect } from '../../../../../src/auth';
import { valueOrDefault } from '../../../../../src/utils';

declare interface MealPlanRecipePageProps {
  recipe: EitherErrorOr<Recipe>;
  mealPlan: EitherErrorOr<MealPlan>;
  userID: string;
  household: EitherErrorOr<Household>;
  groceryList: EitherErrorOr<MealPlanGroceryListItem[]>;
  tasks: EitherErrorOr<MealPlanTask[]>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<MealPlanRecipePageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('RecipePage.getServerSideProps');

  const clientOrRedirect = buildServerSideClientOrRedirect(context);
  if (clientOrRedirect.redirect) {
    span.end();
    return { redirect: clientOrRedirect.redirect };
  }

  if (!clientOrRedirect.client) {
    // this should never occur if the above state is false
    throw new Error('no client returned');
  }
  const apiClient = clientOrRedirect.client.withSpan(span);

  const { mealPlanID: mealPlanIDParam, recipeID: recipeIDParam } = context.query;
  if (!mealPlanIDParam) {
    throw new Error('meal plan ID is somehow missing!');
  }
  const mealPlanID = mealPlanIDParam.toString();

  if (!recipeIDParam) {
    throw new Error('recipe ID is somehow missing!');
  }
  const recipeID = recipeIDParam.toString();

  const extractCookieTimer = timing.addEvent('extract cookie');
  const sessionDetails = userSessionDetailsOrRedirect(context.req.cookies);
  if (sessionDetails.redirect) {
    span.end();
    return { redirect: sessionDetails.redirect };
  }
  const userSessionData = sessionDetails.details;
  extractCookieTimer.end();

  if (userSessionData?.userID) {
    const analyticsTimer = timing.addEvent('analytics');
    serverSideAnalytics.page(userSessionData.userID, 'RECIPE_PAGE', context, {
      recipeID,
      householdID: userSessionData.householdID,
    });
    analyticsTimer.end();
  }

  const fetchMealPlanTimer = timing.addEvent('fetch meal plan');
  const mealPlanPromise = apiClient
    .getMealPlan(mealPlanID)
    .then((result: APIResponse<MealPlan>) => {
      span.addEvent(`meal plan retrieved`);
      return { data: result.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchMealPlanTimer.end();
    });

  const fetchHouseholdTimer = timing.addEvent('fetch household');
  const householdPromise = apiClient
    .getActiveHousehold()
    .then((result: APIResponse<Household>) => {
      span.addEvent(`household retrieved`);
      return { data: result.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchHouseholdTimer.end();
    });

  const fetchMealPlanTasksTimer = timing.addEvent('fetch meal plan tasks');
  const tasksPromise = apiClient
    .getMealPlanTasks(mealPlanID)
    .then((result: QueryFilteredResult<MealPlanTask>) => {
      span.addEvent('meal plan grocery list items retrieved');
      return { data: result.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchMealPlanTasksTimer.end();
    });

  const fetchMealPlanGroceryListItemsTimer = timing.addEvent('fetch meal plan grocery list items');
  const groceryListPromise = apiClient
    .getMealPlanGroceryListItemsForMealPlan(mealPlanID)
    .then((result: QueryFilteredResult<MealPlanGroceryListItem>) => {
      span.addEvent('meal plan grocery list items retrieved');
      return { data: result.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchMealPlanGroceryListItemsTimer.end();
    });

  const fetchRecipeTimer = timing.addEvent('fetch recipe');
  const recipePromise = apiClient
    .getRecipe(recipeID.toString())
    .then((result: APIResponse<Recipe>) => {
      span.addEvent(`recipe retrieved`);
      return { data: result.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchRecipeTimer.end();
    });

  const retrievedData = await Promise.all([
    mealPlanPromise,
    householdPromise,
    groceryListPromise,
    tasksPromise,
    recipePromise,
  ]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  const [mealPlan, household, groceryList, tasks, recipe] = retrievedData;

  span.end();
  return {
    props: {
      recipe: recipe,
      mealPlan: mealPlan!,
      household: household!,
      userID: userSessionData?.userID || '',
      tasks: tasks,
      groceryList: groceryList || [],
    },
  };
};

export default function MealPlanRecipePage(props: MealPlanRecipePageProps) {
  const pageLoadRecipe = props.recipe;

  const ogRecipe = valueOrDefault(pageLoadRecipe, new Recipe());
  const [recipe] = useState<Recipe>(ogRecipe);
  const [recipeError] = useState<IAPIError | undefined>(pageLoadRecipe.error);

  return (
    <AppLayout title={recipe.name} userLoggedIn>
      {recipeError && <Text color="tomato">Error loading recipe: {recipeError.message}</Text>}
      {!recipeError && <RecipeComponent recipe={recipe} />}
    </AppLayout>
  );
}
