import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { TextInput, Button, Group, Container, Divider, Text } from '@mantine/core';
import { useState } from 'react';
import { useRouter } from 'next/router';

import { APIResponse, EitherErrorOr, IAPIError, OAuth2Client } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { AppLayout } from '../../../src/layouts';
import { buildServerSideClientOrRedirect } from '../../../src/client';
import { serverSideTracer } from '../../../src/tracer';
import { errorOrDefault } from '../../../src/utils';

declare interface OAuth2ClientPageProps {
  pageLoadOAuth2Client: EitherErrorOr<OAuth2Client>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<OAuth2ClientPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('OAuth2ClientPage.getServerSideProps');

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

  const { oauth2ClientID } = context.query;
  if (!oauth2ClientID) {
    throw new Error('oauth2 client ID is somehow missing!');
  }

  const fetchOAuth2ClientTimer = timing.addEvent('fetch OAuth2 client');
  const pageLoadOAuth2ClientPromise = apiClient
    .getOAuth2Client(oauth2ClientID.toString())
    .then((result: APIResponse<OAuth2Client>) => {
      span.addEvent('oauth2 client retrieved');
      return result.data;
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred');
      return { error };
    })
    .finally(() => {
      fetchOAuth2ClientTimer.end();
    });

  const [pageLoadOAuth2Client] = await Promise.all([pageLoadOAuth2ClientPromise]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return {
    props: {
      pageLoadOAuth2Client: JSON.parse(JSON.stringify(pageLoadOAuth2Client)),
    },
  };
};

function OAuth2ClientPage(props: OAuth2ClientPageProps) {
  const router = useRouter();

  const apiClient = buildLocalClient();
  const { pageLoadOAuth2Client } = props;

  const ogOAuth2Client: OAuth2Client = errorOrDefault(pageLoadOAuth2Client, new OAuth2Client());

  const [oauth2ClientError, _setOAuth2ClientError] = useState<IAPIError | undefined>(pageLoadOAuth2Client.error);
  const [oauth2Client, _setOAuth2Client] = useState<OAuth2Client>(ogOAuth2Client);

  return (
    <AppLayout title="Valid Preparation">
      <Container size="sm">
        {oauth2ClientError && <Text color="tomato"> {oauth2ClientError.message} </Text>}

        {!oauth2ClientError && oauth2Client.id !== '' && (
          <>
            <TextInput label="Name" value={oauth2Client.name} onChange={() => {}} />
            <TextInput label="Description" value={oauth2Client.description} onChange={() => {}} />
            <TextInput label="Client ID" value={oauth2Client.clientID} onChange={() => {}} />
            <TextInput label="Client Secret" value={oauth2Client.clientSecret} onChange={() => {}} />

            <Divider my="xl" />

            <Group position="center">
              <Button
                type="submit"
                color="tomato"
                fullWidth
                onClick={() => {
                  if (confirm('Are you sure you want to delete this OAuth2 client?')) {
                    apiClient.archiveOAuth2Client(oauth2Client.id).then(() => {
                      router.push('/oauth_clients');
                    });
                  }
                }}
              >
                Delete
              </Button>
            </Group>
          </>
        )}
      </Container>
    </AppLayout>
  );
}

export default OAuth2ClientPage;
