import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Button, Grid, Pagination, Stack, Table, Text, TextInput } from '@mantine/core';
import { AxiosError } from 'axios';
import { formatRelative } from 'date-fns';
import router from 'next/router';
import { IconSearch } from '@tabler/icons';
import { useState, useEffect } from 'react';

import { EitherErrorOr, IAPIError, QueryFilter, QueryFilteredResult, ValidIngredient } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { buildServerSideClientOrRedirect } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideTracer } from '../../src/tracer';
import { errorOrDefault } from '../../src/utils';

declare interface ValidIngredientsPageProps {
  pageLoadValidIngredients: EitherErrorOr<QueryFilteredResult<ValidIngredient>>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidIngredientsPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidIngredientsPage.getServerSideProps');

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

  // TODO: parse context.query as QueryFilter.
  let props!: GetServerSidePropsResult<ValidIngredientsPageProps>;

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  console.log('getting valid ingredients from getServerSideProps');
  const fetchValidIngredientsTimer = timing.addEvent('fetch valid ingredients');
  await apiClient
    .getValidIngredients(qf)
    .then((res: QueryFilteredResult<ValidIngredient>) => {
      span.addEvent('valid ingredients retrieved');
      props = {
        props: {
          pageLoadValidIngredients: JSON.parse(JSON.stringify(res)),
        },
      };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred');
      return { error };
    })
    .finally(() => {
      fetchValidIngredientsTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function ValidIngredientsPage(props: ValidIngredientsPageProps) {
  let { pageLoadValidIngredients } = props;

  const ogValidIngredient = errorOrDefault(pageLoadValidIngredients, new QueryFilteredResult<ValidIngredient>());
  const [validIngredientsError] = useState<IAPIError | undefined>(pageLoadValidIngredients.error);
  const [validIngredients, setValidIngredients] = useState<QueryFilteredResult<ValidIngredient>>(ogValidIngredient);
  const [search, setSearch] = useState('');
  const [currentPage, setCurrentPage] = useState(1);

  useEffect(() => {
    const qf = QueryFilter.deriveFromGetServerSidePropsContext({ search });
    const apiClient = buildLocalClient();

    if (search.trim().length < 1) {
      console.log('getting valid ingredients from search useEffect');
      apiClient
        .getValidIngredients(qf)
        .then((res: QueryFilteredResult<ValidIngredient>) => {
          setValidIngredients(res);
        })
        .catch((err: AxiosError) => {
          console.error('getting valid ingredients: ', err);
        });
    } else {
      console.log('searching for valid ingredients from search useEffect');
      apiClient
        .searchForValidIngredients(search, qf)
        .then((res: QueryFilteredResult<ValidIngredient>) => {
          setValidIngredients(res);
        })
        .catch((err: AxiosError) => {
          console.error('searching for valid ingredients: ', err);
        });
    }
  }, [search]);

  useEffect(() => {
    const apiClient = buildLocalClient();

    const qf = QueryFilter.deriveFromPage();
    qf.page = currentPage;

    console.log('getting valid ingredients from page useEffect');
    apiClient
      .getValidIngredients(qf)
      .then((res: QueryFilteredResult<ValidIngredient>) => {
        setValidIngredients(res);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [currentPage]);

  const formatDate = (x: string | undefined): string => {
    return x ? formatRelative(new Date(x), new Date()) : 'never';
  };

  const rows = (validIngredients.data || []).map((ingredient) => (
    <tr
      key={ingredient.id}
      onClick={() => router.push(`/valid_ingredients/${ingredient.id}`)}
      style={{ cursor: 'pointer' }}
    >
      <td>{ingredient.id}</td>
      <td>{ingredient.name}</td>
      <td>{ingredient.slug}</td>
      <td>{formatDate(ingredient.createdAt)}</td>
      <td>{formatDate(ingredient.lastUpdatedAt)}</td>
    </tr>
  ));

  return (
    <AppLayout title="Valid Ingredients">
      <Stack>
        <Grid justify="space-between">
          <Grid.Col md="auto" sm={12}>
            <TextInput
              placeholder="Search..."
              icon={<IconSearch size={14} />}
              onChange={(event) => setSearch(event.target.value || '')}
            />
          </Grid.Col>
          <Grid.Col md="content" sm={12}>
            <Button
              onClick={() => {
                router.push('/valid_ingredients/new');
              }}
            >
              Create New
            </Button>
          </Grid.Col>
        </Grid>

        {validIngredientsError && <Text color="tomato">{validIngredientsError.message}</Text>}
        {!validIngredientsError && validIngredients.data && (
          <>
            <Table mt="xl" striped highlightOnHover withBorder withColumnBorders>
              <thead>
                <tr>
                  <th>ID</th>
                  <th>Name</th>
                  <th>Slug</th>
                  <th>Created At</th>
                  <th>Last Updated At</th>
                </tr>
              </thead>
              <tbody>{rows}</tbody>
            </Table>

            <Pagination
              disabled={search.trim().length > 0}
              position="center"
              page={validIngredients.page}
              total={Math.ceil(validIngredients.totalCount / validIngredients.limit)}
              onChange={(value: number) => {
                setCurrentPage(value);
              }}
            />
          </>
        )}
      </Stack>
    </AppLayout>
  );
}

export default ValidIngredientsPage;
