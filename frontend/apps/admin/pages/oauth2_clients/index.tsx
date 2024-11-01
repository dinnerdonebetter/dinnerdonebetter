import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { Button, Grid, Pagination, Stack, Table } from '@mantine/core';
import { AxiosError } from 'axios';
import { formatRelative } from 'date-fns';
import router from 'next/router';
import { useEffect, useState } from 'react';

import { EitherErrorOr, IAPIError, OAuth2Client, QueryFilter, QueryFilteredResult } from '@dinnerdonebetter/models';
import { ServerTiming, ServerTimingHeaderName } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { buildServerSideClientOrRedirect } from '../../src/client';
import { AppLayout } from '../../src/layouts';
import { serverSideTracer } from '../../src/tracer';
import { valueOrDefault } from '../../src/utils';

declare interface OAuth2ClientsPageProps {
  pageLoadOAuth2Clients: EitherErrorOr<QueryFilteredResult<OAuth2Client>>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<OAuth2ClientsPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('OAuth2ClientsPage.getServerSideProps');

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
  let props!: GetServerSidePropsResult<OAuth2ClientsPageProps>;

  const qf = QueryFilter.deriveFromGetServerSidePropsContext(context.query);
  qf.attachToSpan(span);

  const fetchOAuth2ClientsTimer = timing.addEvent('fetch oauth2 clients');
  await apiClient
    .getOAuth2Clients(qf)
    .then((res: QueryFilteredResult<OAuth2Client>) => {
      span.addEvent('oauth2 clients retrieved');
      props = {
        props: {
          pageLoadOAuth2Clients: JSON.parse(JSON.stringify(res)),
        },
      };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchOAuth2ClientsTimer.end();
    });

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return props;
};

function OAuth2ClientsPage(props: OAuth2ClientsPageProps) {
  let { pageLoadOAuth2Clients } = props;

  const ogOAuth2Clients: QueryFilteredResult<OAuth2Client> = valueOrDefault(
    pageLoadOAuth2Clients,
    new QueryFilteredResult<OAuth2Client>(),
  );
  const [oauth2ClientsError] = useState<IAPIError | undefined>(pageLoadOAuth2Clients.error);
  const [oauth2Clients, setOAuth2Clients] = useState<QueryFilteredResult<OAuth2Client>>(ogOAuth2Clients);

  useEffect(() => {
    const apiClient = buildLocalClient();

    const qf = QueryFilter.deriveFromPage();
    qf.page = oauth2Clients.page;

    apiClient
      .getOAuth2Clients(qf)
      .then((res: QueryFilteredResult<OAuth2Client>) => {
        setOAuth2Clients(res || []);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [oauth2Clients.page]);

  const formatDate = (x: string | undefined): string => {
    return x ? formatRelative(new Date(x), new Date()) : 'never';
  };

  const rows = (oauth2Clients.data || []).map((oauth2Client) => (
    <tr
      key={oauth2Client.id}
      onClick={() => router.push(`/oauth2_clients/${oauth2Client.id}`)}
      style={{ cursor: 'pointer' }}
    >
      <td>{oauth2Client.name}</td>
      <td>{formatDate(oauth2Client.createdAt)}</td>
    </tr>
  ));

  return (
    <AppLayout title="OAuth2 Clients">
      <Stack>
        <Grid justify="space-between">
          <Grid.Col md="content" sm={12}>
            <Button
              onClick={() => {
                router.push('/oauth2_clients/new');
              }}
            >
              Create New
            </Button>
          </Grid.Col>
        </Grid>

        {oauth2ClientsError && <div>{oauth2ClientsError.message}</div>}
        {!oauth2ClientsError && oauth2Clients.data.length !== 0 && (
          <>
            <Table mt="xl" striped highlightOnHover withBorder withColumnBorders>
              <thead>
                <tr>
                  <th>Name</th>
                  <th>Created At</th>
                </tr>
              </thead>
              <tbody>{rows}</tbody>
            </Table>

            <Pagination
              position="center"
              page={oauth2Clients.page}
              total={Math.ceil(oauth2Clients.totalCount / oauth2Clients.limit)}
              onChange={(value: number) => {
                setOAuth2Clients({ ...oauth2Clients, page: value });
              }}
            />
          </>
        )}
      </Stack>
    </AppLayout>
  );
}

export default OAuth2ClientsPage;
