import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { Button, Center, Container, List, Text } from '@mantine/core';
import { useState } from 'react';

import { Meal, QueryFilteredResult, QueryFilter, EitherErrorOr, IAPIError, MealPlan } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';

import { serverSideTracer } from '../../src/tracer';
import { buildServerSideClientOrRedirect } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideAnalytics } from '../../src/analytics';
import { userSessionDetailsOrRedirect } from '../../src/auth';
import { valueOrDefault } from '../../src/utils';

declare interface MealsPageProps {
  meals: EitherErrorOr<Meal[]>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<MealsPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('MealsPage.getServerSideProps');

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  let props!: GetServerSidePropsResult<MealsPageProps>;

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
    serverSideAnalytics.page(userSessionData.userID, 'MEALS_PAGE', context, {
      query: context.query,
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

  const fetchMealsTimer = timing.addEvent('fetch meals');
  await apiClient
    .getMeals(qf)
    .then((result: QueryFilteredResult<Meal>) => {
      span.addEvent('meals retrieved');
      return { data: result.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchMealsTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function MealsPage(props: MealsPageProps) {
  const router = useRouter();

  const pageLoadMeals = props.meals;
  const ogMeals = valueOrDefault(pageLoadMeals, new Array<MealPlan>());
  const [mealsError] = useState<IAPIError | undefined>(pageLoadMeals.error);
  const [meals] = useState<Meal[]>(ogMeals);

  const mealItems = (meals || []).map((meal: Meal) => (
    <List.Item key={meal.id}>
      <Link href={`/meals/${meal.id}`}>{meal.name}</Link>
    </List.Item>
  ));

  return (
    <AppLayout title="Meals" userLoggedIn>
      <Container size="xs">
        <Center>
          <Button
            my="lg"
            onClick={() => {
              router.push('/meals/new');
            }}
          >
            New Meal
          </Button>
        </Center>

        {mealsError && <Text color="tomato">{mealsError.message}</Text>}

        {!mealsError && <List>{mealItems}</List>}
      </Container>
    </AppLayout>
  );
}

export default MealsPage;
