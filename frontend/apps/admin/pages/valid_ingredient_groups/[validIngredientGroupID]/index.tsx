import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { useForm, zodResolver } from '@mantine/form';
import { TextInput, Button, Group, Container, Text } from '@mantine/core';
import { z } from 'zod';
import { useState } from 'react';
import { useRouter } from 'next/router';

import {
  APIResponse,
  EitherErrorOr,
  IAPIError,
  ValidIngredientGroup,
  ValidIngredientGroupUpdateRequestInput,
} from '@dinnerdonebetter/models';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { AppLayout } from '../../../src/layouts';
import { buildServerSideClientOrRedirect } from '../../../src/client';
import { serverSideTracer } from '../../../src/tracer';
import { inputSlug } from '../../../src/schemas';
import { errorOrDefault } from '../../../src/utils';

declare interface ValidIngredientGroupPageProps {
  pageLoadValidIngredientGroup: EitherErrorOr<ValidIngredientGroup>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidIngredientGroupPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidIngredientGroupPage.getServerSideProps');

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

  const { validIngredientGroupID } = context.query;
  if (!validIngredientGroupID) {
    throw new Error('valid ingredient group ID is somehow missing!');
  }

  const fetchValidIngredientGroupsTimer = timing.addEvent('fetch valid ingredient groups');
  const pageLoadValidIngredientGroupPromise = apiClient
    .getValidIngredientGroup(validIngredientGroupID.toString())
    .then((result: APIResponse<ValidIngredientGroup>) => {
      span.addEvent('valid ingredient group retrieved');
      return result.data;
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred');
      return { error };
    })
    .finally(() => {
      fetchValidIngredientGroupsTimer.end();
    });

  const [pageLoadValidIngredientGroup] = await Promise.all([pageLoadValidIngredientGroupPromise]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return {
    props: {
      pageLoadValidIngredientGroup: JSON.parse(JSON.stringify(pageLoadValidIngredientGroup)),
    },
  };
};

const validIngredientGroupUpdateFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
  pluralName: z.string().trim().min(1, 'plural name is required'),
  slug: inputSlug,
});

function ValidIngredientGroupPage(props: ValidIngredientGroupPageProps) {
  const router = useRouter();

  const { pageLoadValidIngredientGroup } = props;

  const apiClient = buildLocalClient();

  const ogValidIngredientGroup: ValidIngredientGroup = errorOrDefault(
    pageLoadValidIngredientGroup,
    new ValidIngredientGroup(),
  );
  const [validIngredientGroupError] = useState<IAPIError | undefined>(pageLoadValidIngredientGroup.error);

  const [validIngredientGroup, setValidIngredientGroup] = useState<ValidIngredientGroup>(ogValidIngredientGroup);
  const [originalValidIngredientGroup, setOriginalValidIngredientGroup] =
    useState<ValidIngredientGroup>(ogValidIngredientGroup);

  const updateForm = useForm({
    initialValues: validIngredientGroup,
    validate: zodResolver(validIngredientGroupUpdateFormSchema),
  });

  const dataHasChanged = (): boolean => {
    return (
      originalValidIngredientGroup.name !== updateForm.values.name ||
      originalValidIngredientGroup.description !== updateForm.values.description ||
      originalValidIngredientGroup.slug !== updateForm.values.slug
    );
  };

  const submit = async () => {
    const validation = updateForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const submission = new ValidIngredientGroupUpdateRequestInput({
      name: updateForm.values.name,
      description: updateForm.values.description,
      slug: updateForm.values.slug,
    });

    await apiClient
      .updateValidIngredientGroup(validIngredientGroup.id, submission)
      .then((result: APIResponse<ValidIngredientGroup>) => {
        updateForm.setValues(result.data);
        setValidIngredientGroup(result.data);
        setOriginalValidIngredientGroup(result.data);
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="valid ingredient group">
      <Container size="sm">
        {validIngredientGroupError && <Text color="tomato"> {validIngredientGroupError.message} </Text>}
        {!validIngredientGroupError && validIngredientGroup.id !== '' && (
          <>
            <form onSubmit={updateForm.onSubmit(submit)}>
              <TextInput label="Name" placeholder="thing" {...updateForm.getInputProps('name')} />
              <TextInput label="Slug" placeholder="thing" {...updateForm.getInputProps('slug')} />
              <TextInput
                label="Description"
                placeholder="stuff about things"
                {...updateForm.getInputProps('description')}
              />

              <Group position="center">
                <Button type="submit" mt="sm" fullWidth disabled={!dataHasChanged()}>
                  Submit
                </Button>
                <Button
                  type="submit"
                  color="red"
                  fullWidth
                  onClick={() => {
                    if (confirm('Are you sure you want to delete this valid ingredient group?')) {
                      apiClient.archiveValidIngredientGroup(validIngredientGroup.id).then(() => {
                        router.push('/valid_ingredient_groups');
                      });
                    }
                  }}
                >
                  Delete
                </Button>
              </Group>
            </form>
          </>
        )}
      </Container>
    </AppLayout>
  );
}

export default ValidIngredientGroupPage;
