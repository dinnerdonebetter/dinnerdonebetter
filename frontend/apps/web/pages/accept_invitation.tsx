import { Grid, Button, Stack, Space } from '@mantine/core';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { useRouter } from 'next/router';

import { HouseholdInvitationUpdateRequestInput } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';
import { buildBrowserSideClient } from '@dinnerdonebetter/api-client';

import { AppLayout } from '../src/layouts';
import { serverSideTracer } from '../src/tracer';
import { userSessionDetailsOrRedirect } from '../src/auth';

declare interface AcceptInvitationPageProps {
  invitationToken: string;
  invitationID: string;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<AcceptInvitationPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('RegistrationPage.getServerSideProps');

  const invitationToken = context.query['t']?.toString() || '';
  const invitationID = context.query['i']?.toString() || '';

  let props: GetServerSidePropsResult<AcceptInvitationPageProps> = {
    props: {
      invitationID: invitationID,
      invitationToken: invitationToken,
    },
  };

  const extractCookieTimer = timing.addEvent('extract cookie');
  const sessionDetails = userSessionDetailsOrRedirect(context.req.cookies);
  if (sessionDetails.redirect) {
    span.end();
    return { redirect: sessionDetails.redirect };
  }
  const userSessionData = sessionDetails.details;
  extractCookieTimer.end();

  if (!userSessionData?.userID) {
    return {
      redirect: {
        destination: `/register?i=${invitationID}&t=${invitationToken}`,
        permanent: false,
      },
    };
  }

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();

  return props;
};

function AcceptInvitationPage(props: AcceptInvitationPageProps) {
  const { invitationID, invitationToken } = props;
  const router = useRouter();
  const apiClient = buildBrowserSideClient();

  const acceptInvite = async () => {
    const proceed = confirm('Are you sure you want to accept this invite?');
    if (proceed) {
      await apiClient
        .acceptHouseholdInvitation(
          invitationID,
          new HouseholdInvitationUpdateRequestInput({
            token: invitationToken,
          }),
        )
        .then(() => {
          router.push('/');
        })
        .catch(() => {
          window.location.href = `/register?i=${invitationID}&t=${invitationToken}`;
        });
    }
  };

  const rejectInvite = async () => {
    const proceed = confirm('Are you sure you want to reject this invite?');
    if (proceed) {
      await apiClient
        .rejectHouseholdInvitation(
          invitationID,
          new HouseholdInvitationUpdateRequestInput({
            token: invitationToken,
          }),
        )
        .finally(() => {
          router.push('/');
        });
    }
  };

  return (
    <AppLayout title="Accept Invitation" userLoggedIn={false}>
      {' '}
      {/* TODO: this is actually unknown, not false */}
      <Grid mt="xl">
        <Grid.Col span={4}>
          <Space h="xl" />
        </Grid.Col>
        <Grid.Col span="auto">
          <Stack>
            <Button onClick={acceptInvite}>Accept Invite</Button>
            <Button onClick={rejectInvite}>Reject Invite</Button>
          </Stack>
        </Grid.Col>
        <Grid.Col span={4}>
          <Space h="xl" />
        </Grid.Col>
      </Grid>
    </AppLayout>
  );
}

export default AcceptInvitationPage;
