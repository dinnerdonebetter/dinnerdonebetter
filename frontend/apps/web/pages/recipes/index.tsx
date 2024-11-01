import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Container, Grid, List, Space, Text, TextInput } from '@mantine/core';
import Link from 'next/link';
import { useState } from 'react';
import { IconSearch } from '@tabler/icons';

import { QueryFilter, Recipe, QueryFilteredResult, IAPIError, EitherErrorOr } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { buildServerSideClientOrRedirect } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideTracer } from '../../src/tracer';
import { userSessionDetailsOrRedirect } from '../../src/auth';
import { serverSideAnalytics } from '../../src/analytics';
import { valueOrDefault } from '../../src/utils';

declare interface RecipesPageProps {
  recipes: EitherErrorOr<QueryFilteredResult<Recipe>>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<RecipesPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('RecipesPage.getServerSideProps');

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  let props!: GetServerSidePropsResult<RecipesPageProps>;

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
    serverSideAnalytics.page(userSessionData.userID, 'RECIPES_PAGE', context, {
      householdID: userSessionData.householdID,
    });
    analyticsTimer.end();
  } else {
    return {
      redirect: {
        destination: `/login?dest=${encodeURIComponent(context.resolvedUrl)}`,
        permanent: false,
      },
    };
  }

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

  const fetchRecipesTimer = timing.addEvent('fetch recipes');
  await apiClient
    .getRecipes(qf)
    .then((res: QueryFilteredResult<Recipe>) => {
      span.addEvent('recipes retrieved');
      return { data: res.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchRecipesTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function RecipesPage(props: RecipesPageProps) {
  const pageLoadRecipes = props.recipes;

  const ogRecipes = valueOrDefault(pageLoadRecipes, new QueryFilteredResult<Recipe>());
  const [recipesError] = useState<IAPIError | undefined>(pageLoadRecipes.error);
  const [recipes] = useState<Recipe[]>(ogRecipes.data);

  const recipeItems = (recipes || []).map((recipe: Recipe) => (
    <List.Item key={recipe.id}>
      <Link href={`/recipes/${recipe.id}`}>{recipe.name}</Link>
    </List.Item>
  ));

  return (
    <AppLayout title="Recipes" userLoggedIn>
      <Container size="xs">
        <Grid justify="space-between">
          <Grid.Col md="auto" sm={12}>
            <TextInput
              placeholder="Search..."
              disabled
              icon={<IconSearch size={14} />}
              // onChange={(event) => setSearch(event.target.value || '')}
            />
          </Grid.Col>
        </Grid>

        <Space my="xl" />

        {recipesError && <Text color="tomato">{recipesError.message}</Text>}

        {!recipesError && <List>{recipeItems}</List>}
      </Container>
    </AppLayout>
  );
}

export default RecipesPage;
