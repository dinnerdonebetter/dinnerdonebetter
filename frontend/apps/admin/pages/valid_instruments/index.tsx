import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Button, Grid, Pagination, Stack, Table, Text, TextInput } from '@mantine/core';
import { AxiosError } from 'axios';
import { formatRelative } from 'date-fns';
import router from 'next/router';
import { IconSearch } from '@tabler/icons';
import { useEffect, useState } from 'react';

import { EitherErrorOr, IAPIError, QueryFilter, QueryFilteredResult, ValidInstrument } from '@dinnerdonebetter/models';
import { ServerTiming, ServerTimingHeaderName } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { buildServerSideClientOrRedirect } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideTracer } from '../../src/tracer';
import { valueOrDefault } from '../../src/utils';

declare interface ValidInstrumentsPageProps {
  pageLoadValidInstruments: EitherErrorOr<QueryFilteredResult<ValidInstrument>>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidInstrumentsPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidInstrumentsPage.getServerSideProps');

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
  let props!: GetServerSidePropsResult<ValidInstrumentsPageProps>;

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  const fetchValidVesselTimer = timing.addEvent('fetch valid instruments');
  await apiClient
    .getValidInstruments(qf)
    .then((res: QueryFilteredResult<ValidInstrument>) => {
      span.addEvent('valid instruments retrieved');
      props = {
        props: {
          pageLoadValidInstruments: JSON.parse(JSON.stringify(res)),
        },
      };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchValidVesselTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function ValidInstrumentsPage(props: ValidInstrumentsPageProps) {
  let { pageLoadValidInstruments } = props;

  const ogValidInstruments = valueOrDefault(pageLoadValidInstruments, new QueryFilteredResult<ValidInstrument>());
  const [validInstrumentsError] = useState<IAPIError | undefined>(pageLoadValidInstruments.error);
  const [validInstruments, setValidInstruments] = useState<QueryFilteredResult<ValidInstrument>>(ogValidInstruments);
  const [search, setSearch] = useState('');

  useEffect(() => {
    const apiClient = buildLocalClient();

    if (search.trim().length < 1) {
      const qf = QueryFilter.deriveFromGetServerSidePropsContext({ search });
      apiClient
        .getValidInstruments(qf)
        .then((res: QueryFilteredResult<ValidInstrument>) => {
          setValidInstruments(res);
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    } else {
      apiClient
        .searchForValidInstruments(search)
        .then((res: QueryFilteredResult<ValidInstrument>) => {
          setValidInstruments(res);
        })
        .catch((err: AxiosError) => {
          console.error(err);
        });
    }
  }, [search]);

  useEffect(() => {
    const apiClient = buildLocalClient();

    const qf = QueryFilter.deriveFromPage();
    qf.page = validInstruments.page;

    apiClient
      .getValidInstruments(qf)
      .then((res: QueryFilteredResult<ValidInstrument>) => {
        setValidInstruments(res);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [validInstruments.page]);

  const formatDate = (x: string | undefined): string => {
    return x ? formatRelative(new Date(x), new Date()) : 'never';
  };

  const rows = (validInstruments.data || []).map((instrument) => (
    <tr
      key={instrument.id}
      onClick={() => router.push(`/valid_instruments/${instrument.id}`)}
      style={{ cursor: 'pointer' }}
    >
      <td>{instrument.id}</td>
      <td>{instrument.name}</td>
      <td>{instrument.slug}</td>
      <td>{formatDate(instrument.createdAt)}</td>
      <td>{formatDate(instrument.lastUpdatedAt)}</td>
    </tr>
  ));

  return (
    <AppLayout title="Valid Instruments">
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
                router.push('/valid_instruments/new');
              }}
            >
              Create New
            </Button>
          </Grid.Col>
        </Grid>

        {validInstrumentsError && <Text color="tomato"> {validInstrumentsError.message} </Text>}

        {!validInstrumentsError && validInstruments.data && (
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
              page={validInstruments.page}
              total={Math.ceil(validInstruments.totalCount / validInstruments.limit)}
              onChange={(value: number) => {
                setValidInstruments({ ...validInstruments, page: value });
              }}
            />
          </>
        )}
      </Stack>
    </AppLayout>
  );
}

export default ValidInstrumentsPage;
