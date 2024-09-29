import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Button, Grid, Pagination, Stack, Table, TextInput } from '@mantine/core';
import { AxiosError } from 'axios';
import { formatRelative } from 'date-fns';
import router from 'next/router';
import { IconSearch } from '@tabler/icons';
import { useState, useEffect } from 'react';

import { QueryFilter, QueryFilteredResult, ValidIngredient } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { buildServerSideClient } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideTracer } from '../../src/tracer';

declare interface ValidIngredientsPageProps {
  pageLoadValidIngredients: QueryFilteredResult<ValidIngredient>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidIngredientsPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidIngredientsPage.getServerSideProps');
  const apiClient = buildServerSideClient(context).withSpan(span);

  // TODO: parse context.query as QueryFilter.
  let props!: GetServerSidePropsResult<ValidIngredientsPageProps>;

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  const fetchValidIngredientsTimer = timing.addEvent('fetch valid ingredients');
  await apiClient
    .getValidIngredients(qf)
    .then((res: QueryFilteredResult<ValidIngredient>) => {
      span.addEvent('valid ingredients retrieved');
      props = { props: { pageLoadValidIngredients: res } };
    })
    .catch((error: AxiosError) => {
      console.error(`getting valid ingredients`, error.status);
      span.addEvent('error occurred');
      if (error.status === '401') {
        props = {
          redirect: {
            destination: `/login?dest=${encodeURIComponent(context.resolvedUrl)}`,
            permanent: false,
          },
        };
      }
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

  const [validIngredients, setValidIngredients] =
    useState<QueryFilteredResult<ValidIngredient>>(pageLoadValidIngredients);
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
      apiClient
        .searchForValidIngredients(search, qf)
        .then((res: QueryFilteredResult<ValidIngredient>) => {
          setValidIngredients({
            ...QueryFilter.Default(),
            data: res.data || [],
            filteredCount: (res.data || []).length,
            totalCount: (res.data || []).length,
          });
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
      </Stack>
    </AppLayout>
  );
}

export default ValidIngredientsPage;
