import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Button, Grid, Pagination, Stack, Table, TextInput, Text } from '@mantine/core';
import { AxiosError } from 'axios';
import { formatRelative } from 'date-fns';
import router from 'next/router';
import { IconSearch } from '@tabler/icons';
import { useState, useEffect } from 'react';

import {
  EitherErrorOr,
  IAPIError,
  QueryFilter,
  QueryFilteredResult,
  ValidIngredientGroup,
} from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { buildServerSideClientOrRedirect } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideTracer } from '../../src/tracer';
import { errorOrDefault } from '../../src/utils';

declare interface ValidIngredientGroupsPageProps {
  pageLoadValidIngredientGroups: EitherErrorOr<QueryFilteredResult<ValidIngredientGroup>>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidIngredientGroupsPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidIngredientGroupsPage.getServerSideProps');

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
  let props!: GetServerSidePropsResult<ValidIngredientGroupsPageProps>;

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  const fetchValidIngredientGroupsTimer = timing.addEvent('fetch valid ingredient groups');
  await apiClient
    .getValidIngredientGroups(qf)
    .then((res: QueryFilteredResult<ValidIngredientGroup>) => {
      span.addEvent('valid ingredient groups retrieved');
      props = {
        props: {
          pageLoadValidIngredientGroups: JSON.parse(JSON.stringify(res)),
        },
      };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred');
      return { error };
    })
    .finally(() => {
      fetchValidIngredientGroupsTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function ValidIngredientGroupsPage(props: ValidIngredientGroupsPageProps) {
  let { pageLoadValidIngredientGroups } = props;

  const ogValidIngredientGroups = errorOrDefault(
    pageLoadValidIngredientGroups,
    new QueryFilteredResult<ValidIngredientGroup>(),
  );
  const [validIngredientGroupsError] = useState<IAPIError | undefined>(pageLoadValidIngredientGroups.error);

  const [validIngredientGroups, setValidIngredientGroups] =
    useState<QueryFilteredResult<ValidIngredientGroup>>(ogValidIngredientGroups);
  const [search, setSearch] = useState('');

  useEffect(() => {
    const apiClient = buildLocalClient();

    if (search.trim().length < 1) {
      const qf = QueryFilter.deriveFromGetServerSidePropsContext({ search });
      apiClient
        .getValidIngredientGroups(qf)
        .then((res: QueryFilteredResult<ValidIngredientGroup>) => {
          setValidIngredientGroups(res);
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    } else {
      apiClient
        .searchForValidIngredientGroups(search)
        .then((res: QueryFilteredResult<ValidIngredientGroup>) => {
          setValidIngredientGroups(res);
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    }
  }, [search]);

  useEffect(() => {
    const apiClient = buildLocalClient();

    const qf = QueryFilter.deriveFromPage();
    qf.page = validIngredientGroups.page;

    apiClient
      .getValidIngredientGroups(qf)
      .then((res: QueryFilteredResult<ValidIngredientGroup>) => {
        setValidIngredientGroups(res);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [validIngredientGroups.page]);

  const formatDate = (x: string | undefined): string => {
    return x ? formatRelative(new Date(x), new Date()) : 'never';
  };

  const rows = (validIngredientGroups.data || []).map((ingredient) => (
    <tr
      key={ingredient.id}
      onClick={() => router.push(`/valid_ingredient_groups/${ingredient.id}`)}
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
    <AppLayout title="valid ingredient groups">
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
                router.push('/valid_ingredient_groups/new');
              }}
            >
              Create New
            </Button>
          </Grid.Col>
        </Grid>

        {validIngredientGroupsError && <Text color="tomato"> {validIngredientGroupsError.message} </Text>}

        {!validIngredientGroupsError && validIngredientGroups.data.length > 0 && (
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
              page={validIngredientGroups.page}
              total={Math.ceil(validIngredientGroups.totalCount / validIngredientGroups.limit)}
              onChange={(value: number) => {
                setValidIngredientGroups({ ...validIngredientGroups, page: value });
              }}
            />
          </>
        )}
      </Stack>
    </AppLayout>
  );
}

export default ValidIngredientGroupsPage;
