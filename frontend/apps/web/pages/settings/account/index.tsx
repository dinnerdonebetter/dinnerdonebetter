import { AxiosError } from 'axios';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import {
  ActionIcon,
  Alert,
  Avatar,
  Box,
  Button,
  Center,
  Container,
  Divider,
  Grid,
  List,
  Paper,
  Select,
  SimpleGrid,
  Space,
  Stack,
  Text,
  Textarea,
  TextInput,
  Title,
  Tooltip,
} from '@mantine/core';
import { useForm, zodResolver } from '@mantine/form';
import Link from 'next/link';
import { useState } from 'react';
import { IconAlertCircle, IconInfoCircle } from '@tabler/icons';
import { z } from 'zod';

import {
  APIResponse,
  EitherErrorOr,
  Account,
  AccountInvitation,
  AccountInvitationCreationRequestInput,
  AccountUpdateRequestInput,
  AccountUserMembershipWithUser,
  IAPIError,
  QueryFilteredResult,
  ServiceSettingConfiguration,
  User,
} from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { buildServerSideClientOrRedirect } from '../../../src/client';
import { AppLayout } from '../../../src/layouts';
import { serverSideTracer } from '../../../src/tracer';
import { serverSideAnalytics } from '../../../src/analytics';
import { userSessionDetailsOrRedirect } from '../../../src/auth';
import { valueOrDefault } from '../../../src/utils';

declare interface AccountSettingsPageProps {
  account: EitherErrorOr<Account>;
  user: EitherErrorOr<User>;
  invitations: EitherErrorOr<QueryFilteredResult<AccountInvitation>>;
  accountSettings: EitherErrorOr<QueryFilteredResult<ServiceSettingConfiguration>>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<AccountSettingsPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('AccountSettingsPage.getServerSideProps');

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
    serverSideAnalytics.page(userSessionData.userID, 'HOUSEHOLD_SETTINGS_PAGE', context, {
      accountID: userSessionData.accountID,
    });
    analyticsTimer.end();
  }

  const fetchUserTimer = timing.addEvent('fetch user');
  const userPromise = apiClient
    .getSelf()
    .then((result: APIResponse<User>) => {
      span.addEvent('user retrieved');
      return { data: result.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchUserTimer.end();
    });

  const fetchAccountTimer = timing.addEvent('fetch account');
  const accountPromise = apiClient
    .getActiveAccount()
    .then((result: APIResponse<Account>) => {
      span.addEvent('account retrieved');
      return { data: result.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchAccountTimer.end();
    });

  const fetchInvitationsTimer = timing.addEvent('fetch received invitations');
  const invitationsPromise = apiClient
    .getSentAccountInvitations()
    .then((result: QueryFilteredResult<AccountInvitation>) => {
      span.addEvent('invitations retrieved');
      return { data: result };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchInvitationsTimer.end();
    });

  const fetchSettingConfigurationsForAccountTimer = timing.addEvent('fetch configured settings for account');
  const rawAccountSettingsPromise = apiClient
    .getServiceSettingConfigurationsForAccount()
    .then((result: QueryFilteredResult<ServiceSettingConfiguration>) => {
      span.addEvent('service settings retrieved');
      return { data: result };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchSettingConfigurationsForAccountTimer.end();
    });

  const retrievedData = await Promise.all([
    userPromise,
    accountPromise,
    invitationsPromise,
    rawAccountSettingsPromise,
  ]);

  const [user, account, invitations, rawAccountSettings] = retrievedData;

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return {
    props: {
      user,
      account,
      invitations: invitations,
      accountSettings: rawAccountSettings,
    },
  };
};

const inviteFormSchema = z.object({
  emailAddress: z.string().trim().email({ message: 'Invalid email' }),
  toName: z.string().trim().optional(),
  note: z.string().trim().optional(),
});

const accountUpdateSchema = z.object({
  name: z.string().trim().optional(),
  contactPhone: z.string().trim().optional(),
  addressLine1: z.string().trim().optional(),
  addressLine2: z.string().trim().optional(),
  city: z.string().trim().optional(),
  state: z.string().trim().optional(),
  zipCode: z.string().trim().regex(/\d{5}/, 'token must be 6 digits').optional(),
});

const allStates = [
  'AL',
  'AK',
  'AZ',
  'AR',
  'CA',
  'CO',
  'CT',
  'DE',
  'DC',
  'FL',
  'GA',
  'HI',
  'ID',
  'IL',
  'IN',
  'IA',
  'KS',
  'KY',
  'LA',
  'ME',
  'MD',
  'MA',
  'MI',
  'MN',
  'MS',
  'MO',
  'MT',
  'NE',
  'NV',
  'NH',
  'NJ',
  'NM',
  'NY',
  'NC',
  'ND',
  'OH',
  'OK',
  'OR',
  'PA',
  'RI',
  'SC',
  'SD',
  'TN',
  'TX',
  'UT',
  'VT',
  'VA',
  'WA',
  'WV',
  'WI',
  'WY',
];

const exampleNames: { first: string; last: string }[] = [
  { first: 'Julia', last: 'Child' },
  { first: 'Ina', last: 'Garten' },
  { first: 'Lidia', last: 'Bastianich' },
];

export default function AccountSettingsPage(props: AccountSettingsPageProps): JSX.Element {
  const ogAccount = valueOrDefault(props.account, new Account());
  const [accountError] = useState<IAPIError | undefined>(props.account.error);
  const [account, setAccount] = useState<Account>(ogAccount);

  const ogUser = valueOrDefault(props.user, new User());
  const [userError] = useState<IAPIError | undefined>(props.user.error);
  const [user] = useState<User>(ogUser);

  const ogInvitations = valueOrDefault(props.invitations, new QueryFilteredResult<AccountInvitation>());
  const [invitationsError] = useState<IAPIError | undefined>(props.invitations.error);
  const [invitations] = useState<QueryFilteredResult<AccountInvitation>>(ogInvitations);

  const [invitationSubmissionError, setInvitationSubmissionError] = useState('');
  const [userIsAccountAdmin] = useState(
    user.emailAddressVerifiedAt &&
      account.members.find((x: AccountUserMembershipWithUser) => x.belongsToUser?.id === user.id)?.accountRole ===
        'account_admin',
  );

  const outboundPendingInvites = (invitations.data || []).map((invite: AccountInvitation) => {
    return (
      <List.Item key={invite.id}>
        {invite.toEmail} - {invite.status}
      </List.Item>
    );
  });

  const inviteForm = useForm({
    initialValues: {
      emailAddress: '',
      note: '',
      toName: '',
    },
    validate: zodResolver(inviteFormSchema),
  });

  const accountUpdateForm = useForm({
    initialValues: {
      name: account.name,
      contactPhone: account.contactPhone,
      addressLine1: account.addressLine1,
      addressLine2: account.addressLine2,
      city: account.city,
      state: account.state,
      zipCode: account.zipCode,
      country: 'USA',
    },
    validate: zodResolver(accountUpdateSchema),
  });

  const submitInvite = async () => {
    setInvitationSubmissionError('');
    const validation = inviteForm.validate();
    if (validation.hasErrors) {
      return;
    }

    const accountInvitationInput = new AccountInvitationCreationRequestInput({
      toEmail: inviteForm.values.emailAddress,
      note: inviteForm.values.note,
    });

    const apiClient = buildLocalClient();

    await apiClient
      .createAccountInvitation(account.id, accountInvitationInput)
      .then(() => {
        inviteForm.reset();
      })
      .catch((err: AxiosError<IAPIError>) => {
        setInvitationSubmissionError(err?.response?.data.message || 'unknown error occurred');
      });
  };

  const accountDataHasChanged = (): boolean => {
    return (
      account.name !== accountUpdateForm.values.name ||
      account.contactPhone !== accountUpdateForm.values.contactPhone ||
      account.addressLine1 !== accountUpdateForm.values.addressLine1 ||
      account.addressLine2 !== accountUpdateForm.values.addressLine2 ||
      account.city !== accountUpdateForm.values.city ||
      account.state !== accountUpdateForm.values.state ||
      account.zipCode !== accountUpdateForm.values.zipCode ||
      account.country !== accountUpdateForm.values.country
    );
  };

  const updateAccount = async () => {
    const validation = accountUpdateForm.validate();
    if (validation.hasErrors) {
      return;
    }

    const updateInput = new AccountUpdateRequestInput({
      name: accountUpdateForm.values.name,
      contactPhone: accountUpdateForm.values.contactPhone,
      addressLine1: accountUpdateForm.values.addressLine1,
      addressLine2: accountUpdateForm.values.addressLine2,
      city: accountUpdateForm.values.city,
      state: accountUpdateForm.values.state,
      zipCode: accountUpdateForm.values.zipCode,
      country: accountUpdateForm.values.country,
    });

    const apiClient = buildLocalClient();

    await apiClient
      .updateAccount(account.id, updateInput)
      .then(() => {
        setAccount(account);
        accountUpdateForm.reset();
      })
      .catch((err: AxiosError<IAPIError>) => {
        setInvitationSubmissionError(err?.response?.data.message || 'unknown error occurred');
      });
  };

  return (
    <AppLayout title="Account Settings" userLoggedIn>
      <Container size="xs">
        <Center>
          <Title order={2}>Account Settings</Title>
        </Center>

        {accountError && <Text color="tomato">{accountError.message}</Text>}

        {!accountError && (account.members || []).length > 0 && (
          <>
            <Divider my="lg" label="Members" labelPosition="center" />
            <SimpleGrid cols={1}>
              {(account.members || []).map((member: AccountUserMembershipWithUser) => {
                return (
                  <Paper withBorder style={{ width: '100%' }} key={member.id} p="md">
                    <Grid gutter="xl">
                      <Grid.Col span={1}>
                        {member.belongsToUser?.avatar && (
                          <Avatar radius={100} component="a" src={member.belongsToUser.avatar} alt="it's me" />
                        )}

                        {!member.belongsToUser?.avatar && <Avatar radius={100} src={null} alt="no image here" />}
                      </Grid.Col>
                      <Grid.Col span="auto" px="xl" mt={7}>
                        {(member.belongsToUser?.id === user.id && (
                          <Link href="/settings/user">
                            {member.belongsToUser?.firstName ?? member.belongsToUser?.username}
                          </Link>
                        )) || <Text>{member.belongsToUser?.firstName ?? member.belongsToUser?.username}</Text>}
                      </Grid.Col>
                      <Grid.Col span={4} offset={3}>
                        <Grid gutter="xs">
                          <Grid.Col span={10} mr="-xs">
                            <Select
                              disabled={!userIsAccountAdmin}
                              value={member.accountRole === 'account_admin' ? 'Admin' : 'Member'}
                              data={['Admin', 'Member']}
                              onChange={async (role: string) => {
                                if (member.accountRole === 'account_admin' && role === 'Member') {
                                  if (confirm("Are you sure you want to remove this user's admin privileges?")) {
                                    // TODO: update account membership
                                  }
                                } else if (member.accountRole !== 'account_admin' && role === 'Admin') {
                                  if (confirm('Are you sure you want to grant this user admin privileges?')) {
                                    // TODO: update account membership
                                  }
                                }
                              }}
                            />
                          </Grid.Col>
                          <Grid.Col span={2} ml={3} mt={4}>
                            <Tooltip
                              label={
                                member.accountRole === 'account_admin'
                                  ? `Admins are capable of inviting new members, creating meal plans, and generally managing the account.`
                                  : `Members are capable of participating in meal planning, but can't do things like invite new members or propose meal plans.`
                              }
                            >
                              <ActionIcon>
                                <IconInfoCircle size={20} />
                              </ActionIcon>
                            </Tooltip>
                          </Grid.Col>
                        </Grid>
                      </Grid.Col>
                    </Grid>
                  </Paper>
                );
              })}
            </SimpleGrid>
          </>
        )}

        <Divider my="lg" label="Information" labelPosition="center" />
        {!accountError && (
          <Box my="xl">
            <form
              onSubmit={(e) => {
                e.preventDefault();
                updateAccount();
              }}
            >
              <Stack>
                <Alert icon={<IconAlertCircle size={16} />} color="blue">
                  We don&apos;t require you to fill this info out to use the service, but future experiments involving
                  features like grocery delivery may require this information.
                </Alert>

                <TextInput
                  label="Name"
                  placeholder=""
                  disabled={!userIsAccountAdmin}
                  {...accountUpdateForm.getInputProps('name')}
                />
                <Grid>
                  <Grid.Col span={7}>
                    <TextInput
                      label="Address Line 1"
                      placeholder=""
                      disabled={!userIsAccountAdmin}
                      {...accountUpdateForm.getInputProps('addressLine1')}
                    />
                  </Grid.Col>
                  <Grid.Col span={5}>
                    <TextInput
                      label="Address Line 2"
                      placeholder=""
                      disabled={!userIsAccountAdmin}
                      {...accountUpdateForm.getInputProps('addressLine2')}
                    />
                  </Grid.Col>
                </Grid>
                <Grid>
                  <Grid.Col span={7}>
                    <TextInput
                      label="City"
                      placeholder=""
                      disabled={!userIsAccountAdmin}
                      {...accountUpdateForm.getInputProps('city')}
                    />
                  </Grid.Col>
                  <Grid.Col span={2}>
                    <Select
                      label="State"
                      searchable
                      placeholder=""
                      disabled={!userIsAccountAdmin}
                      value={accountUpdateForm.getInputProps('state').value}
                      data={allStates.map((state) => {
                        return { value: state, label: state };
                      })}
                      onChange={(e) => {
                        accountUpdateForm.getInputProps('state').onChange(e);
                      }}
                    />
                  </Grid.Col>
                  <Grid.Col span={3}>
                    <TextInput
                      label="Zip Code"
                      placeholder=""
                      disabled={!userIsAccountAdmin}
                      {...accountUpdateForm.getInputProps('zipCode')}
                    />
                  </Grid.Col>
                </Grid>
                <Button
                  type="submit"
                  disabled={
                    !accountUpdateForm.isValid() || !user.emailAddressVerifiedAt || !accountDataHasChanged()
                  }
                  fullWidth
                >
                  Update
                </Button>
              </Stack>
            </form>
          </Box>
        )}

        {invitationsError && <Text color="tomato">{invitationsError.message}</Text>}

        {!invitationsError && outboundPendingInvites.length > 0 && (
          <>
            <Divider my="lg" label="Awaiting Invites" labelPosition="center" />
            <List>{outboundPendingInvites}</List>
          </>
        )}

        {!invitationsError && invitationSubmissionError && (
          <>
            <Space h="md" />
            <Alert title="Oh no!" color="tomato">
              {invitationSubmissionError}
            </Alert>
          </>
        )}

        <Divider my="lg" label="Send Invite" labelPosition="center" />

        {userError && <Text color="tomato">{userError.message}</Text>}

        {!userError && (
          <form
            onSubmit={(e) => {
              e.preventDefault();
              submitInvite();
            }}
          >
            <Grid>
              <Grid.Col md={12} lg="auto">
                <TextInput
                  label="Email Address"
                  disabled={!user.emailAddressVerifiedAt}
                  placeholder="neato_person@fake.email"
                  {...inviteForm.getInputProps('emailAddress')}
                />
              </Grid.Col>
              <Grid.Col md={12} lg="auto">
                <TextInput
                  label="Name"
                  placeholder={exampleNames[Math.floor(Math.random() * exampleNames.length)].first}
                  disabled={!user.emailAddressVerifiedAt}
                  {...inviteForm.getInputProps('toName')}
                />
              </Grid.Col>
            </Grid>
            <Grid>
              <Grid.Col md={12} lg="auto">
                <Textarea
                  label="Note"
                  disabled={!user.emailAddressVerifiedAt}
                  placeholder="Join my account on Dinner Done Better!"
                  {...inviteForm.getInputProps('note')}
                />
              </Grid.Col>
            </Grid>
            <Grid>
              <Grid.Col md={12} lg={12}>
                <Button type="submit" disabled={!inviteForm.isValid() || !user.emailAddressVerifiedAt} fullWidth>
                  Send Invite
                </Button>
              </Grid.Col>
            </Grid>
          </form>
        )}
      </Container>
    </AppLayout>
  );
}
