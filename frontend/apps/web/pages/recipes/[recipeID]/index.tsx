import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Text } from '@mantine/core';
import { useState } from 'react';

import { APIResponse, EitherErrorOr, IAPIError, Recipe } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildServerSideClientOrRedirect } from '../../../src/client';
import { AppLayout } from '../../../src/layouts';
import { RecipeComponent } from '../../../src/components';
import { serverSideTracer } from '../../../src/tracer';
import { serverSideAnalytics } from '../../../src/analytics';
import { userSessionDetailsOrRedirect } from '../../../src/auth';
import { valueOrDefault } from '../../../src/utils';

declare interface RecipePageProps {
  recipe: EitherErrorOr<Recipe>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<RecipePageProps>> => {
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

  const { recipeID } = context.query;
  if (!recipeID) {
    throw new Error('recipe ID is somehow missing!');
  }

  const extractCookieTimer = timing.addEvent('extract cookie');
  const sessionDetails = userSessionDetailsOrRedirect(context.req.cookies);
  if (sessionDetails.redirect) {
    span.end();
    return { redirect: sessionDetails.redirect };
  }
  const userSessionData = sessionDetails.details;
  extractCookieTimer.end();

  if (userSessionData?.userID) {
    serverSideAnalytics.page(userSessionData.userID, 'RECIPE_PAGE', context, {
      recipeID,
      householdID: userSessionData.householdID,
    });
  }

  const fetchRecipeTimer = timing.addEvent('fetch recipe');
  let props!: GetServerSidePropsResult<RecipePageProps>;
  await apiClient
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

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function RecipePage(props: RecipePageProps) {
  const pageLoadRecipe = props.recipe;
  const ogRecipe = valueOrDefault(pageLoadRecipe, new Recipe());
  const [recipe] = useState<Recipe>(ogRecipe);
  const [recipeError] = useState<IAPIError | undefined>(pageLoadRecipe.error);

  return (
    <AppLayout title={recipe.name} titlePosition="left" userLoggedIn>
      {recipeError && <Text color="tomato">{recipeError.message}</Text>}

      {!recipeError && <RecipeComponent recipe={recipe} />}
    </AppLayout>
  );
}

export default RecipePage;
