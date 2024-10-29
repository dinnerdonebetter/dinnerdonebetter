import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { useForm, zodResolver } from '@mantine/form';
import { TextInput, Button, Group, Container, Divider, Space, Select } from '@mantine/core';
import { useState } from 'react';
import { useRouter } from 'next/router';
import { z } from 'zod';

import { APIResponse, EitherErrorOr, IAPIError, ServiceSetting } from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { AppLayout } from '../../../src/layouts';
import { buildServerSideClientOrRedirect } from '../../../src/client';
import { serverSideTracer } from '../../../src/tracer';
import { errorOrDefault } from '../../../src/utils';

declare interface ServiceSettingPageProps {
  pageLoadServiceSetting: EitherErrorOr<ServiceSetting>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ServiceSettingPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ServiceSettingPage.getServerSideProps');

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

  const { settingID } = context.query;
  if (!settingID) {
    throw new Error('service setting ID is somehow missing!');
  }

  const fetchServiceSettingTimer = timing.addEvent('fetch service setting');
  const pageLoadServiceSettingPromise = apiClient
    .getServiceSetting(settingID.toString())
    .then((result: APIResponse<ServiceSetting>) => {
      span.addEvent('service setting retrieved');
      return result.data;
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred');
      return { error };
    })
    .finally(() => {
      fetchServiceSettingTimer.end();
    });

  const [pageLoadServiceSetting] = await Promise.all([pageLoadServiceSettingPromise]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return {
    props: {
      pageLoadServiceSetting: JSON.parse(JSON.stringify(pageLoadServiceSetting)),
    },
  };
};

const serviceSettingUpdateFormSchema = z.object({
  name: z.string().min(1, 'name is required').trim(),
});

function ServiceSettingPage(props: ServiceSettingPageProps) {
  const router = useRouter();

  const apiClient = buildLocalClient();
  const { pageLoadServiceSetting } = props;

  const ogServiceSetting: ServiceSetting = errorOrDefault(pageLoadServiceSetting, new ServiceSetting());
  const [serviceSettingError] = useState<IAPIError | undefined>(pageLoadServiceSetting.error);
  const [serviceSetting] = useState<ServiceSetting>(ogServiceSetting);

  const updateForm = useForm({
    initialValues: serviceSetting,
    validate: zodResolver(serviceSettingUpdateFormSchema),
  });

  return (
    <AppLayout title="Service Setting">
      <Container size="sm">
        {serviceSettingError && <div>{serviceSettingError.message}</div>}
        {!serviceSettingError && (
          <form>
            <TextInput label="Name" placeholder="thing" {...updateForm.getInputProps('name')} />
            <TextInput label="Description" placeholder="thing" {...updateForm.getInputProps('description')} />

            <Select
              label="Type"
              required
              onChange={async (item: string) => {
                updateForm.setFieldValue('type', item);
              }}
              {...updateForm.getInputProps('type')}
              data={['user', 'household', 'membership']}
            />

            <TextInput label="Default Value" placeholder="thing" {...updateForm.getInputProps('defaultValue')} />
            <TextInput label="Enumeration" placeholder="thing" {...updateForm.getInputProps('enumeration')} />

            <Group position="center">
              <Button
                type="submit"
                color="red"
                fullWidth
                onClick={() => {
                  if (confirm('Are you sure you want to delete this service setting?')) {
                    apiClient.archiveServiceSetting(serviceSetting.id).then(() => {
                      router.push('/service_settings');
                    });
                  }
                }}
              >
                Delete
              </Button>
            </Group>
          </form>
        )}

        <Space h="xl" />
        <Divider />
        <Space h="xl" />
      </Container>
    </AppLayout>
  );
}

export default ServiceSettingPage;
