import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import { useForm, zodResolver } from '@mantine/form';
import {
  ActionIcon,
  Autocomplete,
  AutocompleteItem,
  Button,
  Center,
  Container,
  Divider,
  Grid,
  Group,
  Pagination,
  Space,
  Switch,
  Table,
  Text,
  TextInput,
  ThemeIcon,
  Title,
} from '@mantine/core';
import { AxiosError } from 'axios';
import { z } from 'zod';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { IconTrash } from '@tabler/icons';

import {
  APIResponse,
  EitherErrorOr,
  IAPIError,
  QueryFilteredResult,
  ValidInstrument,
  ValidInstrumentUpdateRequestInput,
  ValidPreparation,
  ValidPreparationInstrument,
  ValidPreparationInstrumentCreationRequestInput,
} from '@dinnerdonebetter/models';
import { ServerTiming, ServerTimingHeaderName } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { AppLayout } from '../../../src/layouts';
import { buildServerSideClientOrRedirect } from '../../../src/client';
import { serverSideTracer } from '../../../src/tracer';
import { inputSlug } from '../../../src/schemas';
import { valueOrDefault } from '../../../src/utils';

declare interface ValidInstrumentPageProps {
  pageLoadValidInstrument: EitherErrorOr<ValidInstrument>;
  pageLoadPreparationInstruments: EitherErrorOr<QueryFilteredResult<ValidPreparationInstrument>>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidInstrumentPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidInstrumentPage.getServerSideProps');

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

  const { validInstrumentID } = context.query;
  if (!validInstrumentID) {
    throw new Error('valid instrument ID is somehow missing!');
  }

  const fetchValidInstrumentTimer = timing.addEvent('fetch valid instrument');
  const pageLoadValidInstrumentPromise = apiClient
    .getValidInstrument(validInstrumentID.toString())
    .then((result: APIResponse<ValidInstrument>) => {
      span.addEvent('valid instrument retrieved');
      return result;
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchValidInstrumentTimer.end();
    });

  const fetchPreparationInstrumentsTimer = timing.addEvent('fetch valid preparation instruments');
  const pageLoadPreparationInstrumentsPromise = apiClient
    .getValidPreparationInstrumentsByInstrument(validInstrumentID.toString())
    .then((res: QueryFilteredResult<ValidPreparationInstrument>) => {
      return res;
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchPreparationInstrumentsTimer.end();
    });

  const [pageLoadValidInstrument, pageLoadPreparationInstruments] = await Promise.all([
    pageLoadValidInstrumentPromise,
    pageLoadPreparationInstrumentsPromise,
  ]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return {
    props: {
      pageLoadValidInstrument: JSON.parse(JSON.stringify(pageLoadValidInstrument)),
      pageLoadPreparationInstruments: JSON.parse(JSON.stringify(pageLoadPreparationInstruments)),
    },
  };
};

const validInstrumentUpdateFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
  pluralName: z.string().trim().min(1, 'plural name is required'),
  slug: inputSlug,
});

function ValidInstrumentPage(props: ValidInstrumentPageProps) {
  const router = useRouter();

  const apiClient = buildLocalClient();
  const { pageLoadValidInstrument, pageLoadPreparationInstruments } = props;

  const ogValidIngredient = valueOrDefault(pageLoadValidInstrument, new ValidInstrument());
  const [validInstrumentError] = useState<IAPIError | undefined>(pageLoadValidInstrument.error);
  const [validInstrument, setValidInstrument] = useState<ValidInstrument>(ogValidIngredient);
  const [originalValidInstrument, setOriginalValidInstrument] = useState<ValidInstrument>(ogValidIngredient);

  const [newPreparationForInstrumentInput, setNewPreparationForInstrumentInput] =
    useState<ValidPreparationInstrumentCreationRequestInput>(
      new ValidPreparationInstrumentCreationRequestInput({
        validInstrumentID: validInstrument.id,
      }),
    );
  const [preparationQuery, setPreparationQuery] = useState('');

  const ogPreparationsForInstrument = valueOrDefault(
    pageLoadPreparationInstruments,
    new QueryFilteredResult<ValidPreparationInstrument>(),
  );
  const [preparationsForInstrumentError] = useState<IAPIError | undefined>(pageLoadPreparationInstruments.error);
  const [preparationsForInstrument, setPreparationsForInstrument] =
    useState<QueryFilteredResult<ValidPreparationInstrument>>(ogPreparationsForInstrument);
  const [suggestedPreparations, setSuggestedPreparations] = useState<ValidPreparation[]>([]);

  useEffect(() => {
    if (preparationQuery.length < 2) {
      setSuggestedPreparations([]);
      return;
    }

    const apiClient = buildLocalClient();
    apiClient
      .searchForValidPreparations(preparationQuery)
      .then((res: QueryFilteredResult<ValidPreparation>) => {
        const newSuggestions = (res.data || []).filter((mu: ValidPreparation) => {
          return !(preparationsForInstrument.data || []).some((vimu: ValidPreparationInstrument) => {
            return vimu.instrument.id === mu.id;
          });
        });

        setSuggestedPreparations(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [preparationQuery, preparationsForInstrument.data]);

  const updateForm = useForm({
    initialValues: validInstrument,
    validate: zodResolver(validInstrumentUpdateFormSchema),
  });

  const dataHasChanged = (): boolean => {
    return (
      originalValidInstrument.description !== updateForm.values.description ||
      originalValidInstrument.iconPath !== updateForm.values.iconPath ||
      originalValidInstrument.name !== updateForm.values.name ||
      originalValidInstrument.pluralName !== updateForm.values.pluralName ||
      originalValidInstrument.slug !== updateForm.values.slug ||
      originalValidInstrument.displayInSummaryLists !== updateForm.values.displayInSummaryLists ||
      originalValidInstrument.includeInGeneratedInstructions !== updateForm.values.includeInGeneratedInstructions
    );
  };

  const submit = async () => {
    const validation = updateForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const submission = new ValidInstrumentUpdateRequestInput({
      name: updateForm.values.name,
      pluralName: updateForm.values.pluralName,
      description: updateForm.values.description,
      iconPath: updateForm.values.iconPath,
      slug: updateForm.values.slug,
      displayInSummaryLists: updateForm.values.displayInSummaryLists,
      includeInGeneratedInstructions: updateForm.values.includeInGeneratedInstructions,
    });

    const apiClient = buildLocalClient();

    await apiClient
      .updateValidInstrument(validInstrument.id, submission)
      .then((result: APIResponse<ValidInstrument>) => {
        updateForm.setValues(result.data);
        setValidInstrument(result.data);
        setOriginalValidInstrument(result.data);
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="Valid Instrument">
      <Container size="sm">
        {validInstrumentError && <Text color="tomato"> {validInstrumentError.message} </Text>}
        {!validInstrumentError && validInstrument.id !== '' && (
          <>
            <form onSubmit={updateForm.onSubmit(submit)}>
              <TextInput label="Name" placeholder="thing" {...updateForm.getInputProps('name')} />
              <TextInput label="Plural Name" placeholder="things" {...updateForm.getInputProps('pluralName')} />
              <TextInput label="Slug" placeholder="thing" {...updateForm.getInputProps('slug')} />
              <TextInput label="Description" placeholder="thing" {...updateForm.getInputProps('description')} />

              <Switch
                checked={updateForm.values.displayInSummaryLists}
                label="Display in summary lists"
                {...updateForm.getInputProps('displayInSummaryLists')}
              />

              <Switch
                checked={updateForm.values.includeInGeneratedInstructions}
                label="Include in generated instructions"
                {...updateForm.getInputProps('includeInGeneratedInstructions')}
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
                    if (confirm('Are you sure you want to delete this valid instrument?')) {
                      apiClient.archiveValidInstrument(validInstrument.id).then(() => {
                        router.push('/valid_instruments');
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

        {/*

        INGREDIENT MEASUREMENT UNITS

        */}

        <Space h="xl" />
        <Divider />
        <Space h="xl" />

        <form>
          <Center>
            <Title order={4}>Preparations</Title>
          </Center>

          {preparationsForInstrumentError && <Text color="tomato">{preparationsForInstrumentError.message}</Text>}
          {!preparationsForInstrumentError &&
            preparationsForInstrument.data &&
            (preparationsForInstrument.data || []).length !== 0 && (
              <>
                <Table mt="xl" withColumnBorders>
                  <thead>
                    <tr>
                      <th>Name</th>
                      <th>Notes</th>
                      <th>
                        <Center>
                          <ThemeIcon variant="outline" color="gray">
                            <IconTrash size="sm" color="gray" />
                          </ThemeIcon>
                        </Center>
                      </th>
                    </tr>
                  </thead>
                  <tbody>
                    {(preparationsForInstrument.data || []).map((preparationInstrument: ValidPreparationInstrument) => {
                      return (
                        <tr key={preparationInstrument.id}>
                          <td>
                            <Link href={`/valid_preparations/${preparationInstrument.id}`}>
                              {preparationInstrument.preparation.name}
                            </Link>
                          </td>
                          <td>
                            <Text>{preparationInstrument.notes}</Text>
                          </td>
                          <td>
                            <Center>
                              <ActionIcon
                                variant="outline"
                                aria-label="remove valid preparation measurement unit"
                                onClick={async () => {
                                  await apiClient
                                    .archiveValidPreparationInstrument(preparationInstrument.id)
                                    .then(() => {
                                      setPreparationsForInstrument({
                                        ...preparationsForInstrument,
                                        data: preparationsForInstrument.data.filter(
                                          (x: ValidPreparationInstrument) => x.id !== preparationInstrument.id,
                                        ),
                                      });
                                    })
                                    .catch((error: AxiosError) => {
                                      console.error(error);
                                    });
                                }}
                              >
                                <IconTrash size="md" color="tomato" />
                              </ActionIcon>
                            </Center>
                          </td>
                        </tr>
                      );
                    })}
                  </tbody>
                </Table>

                <Space h="xs" />

                <Pagination
                  disabled={
                    Math.ceil(preparationsForInstrument.totalCount / preparationsForInstrument.limit) <=
                    preparationsForInstrument.page
                  }
                  position="center"
                  page={preparationsForInstrument.page}
                  total={Math.ceil(preparationsForInstrument.totalCount / preparationsForInstrument.limit)}
                  onChange={(value: number) => {
                    setPreparationsForInstrument({ ...preparationsForInstrument, page: value });
                  }}
                />
              </>
            )}

          <Grid>
            <Grid.Col span="auto">
              <Autocomplete
                placeholder="mince"
                label="Preparation"
                value={preparationQuery}
                onChange={setPreparationQuery}
                onItemSubmit={async (item: AutocompleteItem) => {
                  const selectedPreparation = (suggestedPreparations || []).find(
                    (x: ValidPreparation) => x.name === item.value,
                  );

                  if (!selectedPreparation) {
                    console.error(`selectedPreparation not found for item ${item.value}}`);
                    return;
                  }

                  setNewPreparationForInstrumentInput({
                    ...newPreparationForInstrumentInput,
                    validPreparationID: selectedPreparation.id,
                  });
                }}
                data={(suggestedPreparations || []).map((x: ValidPreparation) => {
                  return { value: x.name, label: x.name };
                })}
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <TextInput
                label="Notes"
                value={newPreparationForInstrumentInput.notes}
                onChange={(event) =>
                  setNewPreparationForInstrumentInput({
                    ...newPreparationForInstrumentInput,
                    notes: event.target.value,
                  })
                }
              />
            </Grid.Col>
            <Grid.Col span={2}>
              <Button
                mt="xl"
                disabled={
                  newPreparationForInstrumentInput.validInstrumentID === '' ||
                  newPreparationForInstrumentInput.validPreparationID === ''
                }
                onClick={async () => {
                  await apiClient
                    .createValidPreparationInstrument(newPreparationForInstrumentInput)
                    .then((res: APIResponse<ValidPreparationInstrument>) => {
                      // the returned value doesn't have enough information to put it in the list, so we have to fetch it
                      apiClient
                        .getValidPreparationInstrument(res.data.id)
                        .then((res: APIResponse<ValidPreparationInstrument>) => {
                          setPreparationsForInstrument({
                            ...preparationsForInstrument,
                            data: [...(preparationsForInstrument.data || []), res.data],
                          });

                          setNewPreparationForInstrumentInput(
                            new ValidPreparationInstrumentCreationRequestInput({
                              validInstrumentID: validInstrument.id,
                              validPreparationID: '',
                              notes: '',
                            }),
                          );

                          setPreparationQuery('');
                        });
                    })
                    .catch((error) => {
                      console.error(error);
                    });
                }}
              >
                Save
              </Button>
            </Grid.Col>
          </Grid>
        </form>
      </Container>
    </AppLayout>
  );
}

export default ValidInstrumentPage;
