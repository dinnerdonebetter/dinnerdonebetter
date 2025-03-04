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
  Select,
  Space,
  Table,
  Text,
  TextInput,
  ThemeIcon,
  Title,
} from '@mantine/core';
import { useRouter } from 'next/router';
import { AxiosError } from 'axios';
import { useEffect, useState } from 'react';
import Link from 'next/link';
import { IconTrash } from '@tabler/icons';
import { z } from 'zod';

import {
  APIResponse,
  EitherErrorOr,
  IAPIError,
  QueryFilteredResult,
  ValidIngredient,
  ValidIngredientState,
  ValidIngredientStateIngredient,
  ValidIngredientStateIngredientCreationRequestInput,
  ValidIngredientStateUpdateRequestInput,
} from '@dinnerdonebetter/models';
import { ServerTiming, ServerTimingHeaderName } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { AppLayout } from '../../../src/layouts';
import { buildServerSideClientOrRedirect } from '../../../src/client';
import { serverSideTracer } from '../../../src/tracer';
import { inputSlug } from '../../../src/schemas';
import { valueOrDefault } from '../../../src/utils';

declare interface ValidIngredientStatePageProps {
  pageLoadValidIngredientState: EitherErrorOr<ValidIngredientState>;
  pageLoadValidIngredientStates: EitherErrorOr<QueryFilteredResult<ValidIngredientStateIngredient>>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidIngredientStatePageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidIngredientStatePage.getServerSideProps');

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

  const { validIngredientStateID } = context.query;
  if (!validIngredientStateID) {
    throw new Error('valid ingredient state ID is somehow missing!');
  }

  const fetchValidIngredientStateTimer = timing.addEvent('fetch valid ingredient state');
  const pageLoadValidIngredientStatePromise = apiClient
    .getValidIngredientState(validIngredientStateID.toString())
    .then((result: APIResponse<ValidIngredientState>) => {
      span.addEvent('valid ingredient state retrieved');
      return result;
    })
    .finally(() => {
      fetchValidIngredientStateTimer.end();
    });

  const fetchValidIngredientStatesTimer = timing.addEvent('fetch valid ingredient states for ingredient');
  const pageLoadValidIngredientStatesPromise = apiClient
    .getValidIngredientStateIngredientsByIngredientState(validIngredientStateID.toString())
    .then((res: QueryFilteredResult<ValidIngredientStateIngredient>) => {
      span.addEvent('valid ingredient states retrieved');
      return res;
    })
    .finally(() => {
      fetchValidIngredientStatesTimer.end();
    });

  const [pageLoadValidIngredientState, pageLoadValidIngredientStates] = await Promise.all([
    pageLoadValidIngredientStatePromise,
    pageLoadValidIngredientStatesPromise,
  ]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return {
    props: {
      pageLoadValidIngredientState: JSON.parse(JSON.stringify(pageLoadValidIngredientState)),
      pageLoadValidIngredientStates: JSON.parse(JSON.stringify(pageLoadValidIngredientStates)),
    },
  };
};

const validIngredientStateUpdateFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
  pastTense: z.string().trim().min(1, 'past tense is required'),
  slug: inputSlug,
  attributeType: z.enum([
    'texture',
    'consistency',
    'temperature',
    'color',
    'appearance',
    'odor',
    'taste',
    'sound',
    'other',
  ]),
});

function ValidIngredientStatePage(props: ValidIngredientStatePageProps) {
  const router = useRouter();

  const apiClient = buildLocalClient();
  const { pageLoadValidIngredientState, pageLoadValidIngredientStates } = props;

  const ogValidIngredientState: ValidIngredientState = valueOrDefault(
    pageLoadValidIngredientState,
    new ValidIngredientState(),
  );
  const [validIngredientStateError] = useState<IAPIError | undefined>(pageLoadValidIngredientState.error);
  const [validIngredientState, setValidIngredientState] = useState<ValidIngredientState>(ogValidIngredientState);
  const [originalValidIngredientState, setOriginalValidIngredientState] =
    useState<ValidIngredientState>(ogValidIngredientState);

  const [newIngredientForIngredientStateInput, setNewIngredientForIngredientStateInput] =
    useState<ValidIngredientStateIngredientCreationRequestInput>(
      new ValidIngredientStateIngredientCreationRequestInput({
        validIngredientStateID: validIngredientState.id,
      }),
    );
  const [ingredientQuery, setIngredientQuery] = useState('');

  const ogIngredientsForIngredientState: QueryFilteredResult<ValidIngredientStateIngredient> = valueOrDefault(
    pageLoadValidIngredientStates,
    new QueryFilteredResult<ValidIngredientStateIngredient>(),
  );
  const [ingredientsForIngredientStateError] = useState<IAPIError | undefined>(pageLoadValidIngredientStates.error);
  const [ingredientsForIngredientState, setIngredientsForIngredientState] = useState<
    QueryFilteredResult<ValidIngredientStateIngredient>
  >(ogIngredientsForIngredientState);

  const [suggestedIngredients, setSuggestedIngredients] = useState<ValidIngredient[]>([]);

  useEffect(() => {
    if (ingredientQuery.length <= 2) {
      setSuggestedIngredients([]);
      return;
    }

    const apiClient = buildLocalClient();
    apiClient
      .searchForValidIngredients(ingredientQuery)
      .then((res: QueryFilteredResult<ValidIngredient>) => {
        const newSuggestions = (res.data || []).filter((mu: ValidIngredient) => {
          return !(ingredientsForIngredientState.data || []).some((vimu: ValidIngredientStateIngredient) => {
            return vimu.ingredient.id === mu.id;
          });
        });

        setSuggestedIngredients(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [ingredientQuery, ingredientsForIngredientState.data]);

  const updateForm = useForm({
    initialValues: validIngredientState,
    validate: zodResolver(validIngredientStateUpdateFormSchema),
  });

  const dataHasChanged = (): boolean => {
    return (
      originalValidIngredientState.name !== updateForm.values.name ||
      originalValidIngredientState.description !== updateForm.values.description ||
      originalValidIngredientState.pastTense !== updateForm.values.pastTense ||
      originalValidIngredientState.slug !== updateForm.values.slug ||
      originalValidIngredientState.attributeType !== updateForm.values.attributeType
    );
  };

  const submit = async () => {
    const validation = updateForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const submission = new ValidIngredientStateUpdateRequestInput({
      name: updateForm.values.name,
      description: updateForm.values.description,
      pastTense: updateForm.values.pastTense,
      slug: updateForm.values.slug,
      attributeType: updateForm.values.attributeType,
    });

    const apiClient = buildLocalClient();

    await apiClient
      .updateValidIngredientState(validIngredientState.id, submission)
      .then((result: APIResponse<ValidIngredientState>) => {
        if (result) {
          updateForm.setValues(result.data);
          setValidIngredientState(result.data);
          setOriginalValidIngredientState(result.data);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="Valid Ingredient State">
      <Container size="sm">
        {validIngredientStateError && <Text color="tomato">{validIngredientStateError.message}</Text>}

        {!validIngredientStateError && validIngredientState.id !== '' && (
          <>
            <form onSubmit={updateForm.onSubmit(submit)}>
              <TextInput label="Name" placeholder="thing" {...updateForm.getInputProps('name')} />
              <TextInput label="Past Tense" placeholder="things" {...updateForm.getInputProps('pastTense')} />
              <TextInput label="Slug" placeholder="thing" {...updateForm.getInputProps('slug')} />
              <TextInput label="Description" placeholder="thing" {...updateForm.getInputProps('description')} />

              <Select
                label="Component Type"
                placeholder="Type"
                value={updateForm.values.attributeType}
                data={[
                  { value: 'texture', label: 'texture' },
                  { value: 'consistency', label: 'consistency' },
                  { value: 'temperature', label: 'temperature' },
                  { value: 'color', label: 'color' },
                  { value: 'appearance', label: 'appearance' },
                  { value: 'odor', label: 'odor' },
                  { value: 'taste', label: 'taste' },
                  { value: 'sound', label: 'sound' },
                  { value: 'other', label: 'other' },
                ]}
                {...updateForm.getInputProps('attributeType')}
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
                    if (confirm('Are you sure you want to delete this valid ingredient state?')) {
                      apiClient.archiveValidIngredientState(validIngredientState.id).then(() => {
                        router.push('/valid_ingredient_states');
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
            <Title order={4}>Ingredients</Title>
          </Center>

          {ingredientsForIngredientStateError && (
            <Text color="tomato">{ingredientsForIngredientStateError.message}</Text>
          )}

          {ingredientsForIngredientState.data && (ingredientsForIngredientState.data || []).length !== 0 && (
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
                  {(ingredientsForIngredientState.data || []).map(
                    (ingredientStateIngredient: ValidIngredientStateIngredient) => {
                      return (
                        <tr key={ingredientStateIngredient.id}>
                          <td>
                            <Link href={`/valid_ingredients/${ingredientStateIngredient.ingredient.id}`}>
                              {ingredientStateIngredient.ingredient.name}
                            </Link>
                          </td>
                          <td>
                            <Text>{ingredientStateIngredient.notes}</Text>
                          </td>
                          <td>
                            <Center>
                              <ActionIcon
                                variant="outline"
                                aria-label="remove valid ingredient measurement unit"
                                onClick={async () => {
                                  await apiClient
                                    .archiveValidIngredientStateIngredient(ingredientStateIngredient.id)
                                    .then(() => {
                                      setIngredientsForIngredientState({
                                        ...ingredientsForIngredientState,
                                        data: ingredientsForIngredientState.data.filter(
                                          (x: ValidIngredientStateIngredient) => x.id !== ingredientStateIngredient.id,
                                        ),
                                      });
                                    })
                                    .catch((error) => {
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
                    },
                  )}
                </tbody>
              </Table>

              <Space h="xs" />

              <Pagination
                disabled={
                  Math.ceil(ingredientsForIngredientState.totalCount / ingredientsForIngredientState.limit) <=
                  ingredientsForIngredientState.page
                }
                position="center"
                page={ingredientsForIngredientState.page}
                total={Math.ceil(ingredientsForIngredientState.totalCount / ingredientsForIngredientState.limit)}
                onChange={(value: number) => {
                  setIngredientsForIngredientState({ ...ingredientsForIngredientState, page: value });
                }}
              />
            </>
          )}

          <Grid>
            <Grid.Col span="auto">
              <Autocomplete
                placeholder="grams"
                label="Ingredient"
                value={ingredientQuery}
                onChange={setIngredientQuery}
                onItemSubmit={async (item: AutocompleteItem) => {
                  const selectedIngredient = suggestedIngredients.find((x: ValidIngredient) => x.name === item.value);

                  if (!selectedIngredient) {
                    console.error(`selectedIngredient not found for item ${item.value}}`);
                    return;
                  }

                  setNewIngredientForIngredientStateInput({
                    ...newIngredientForIngredientStateInput,
                    validIngredientID: selectedIngredient.id,
                  });
                }}
                data={suggestedIngredients.map((x: ValidIngredient) => {
                  return { value: x.name, label: x.name };
                })}
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <TextInput
                label="Notes"
                value={newIngredientForIngredientStateInput.notes}
                onChange={(event) =>
                  setNewIngredientForIngredientStateInput({
                    ...newIngredientForIngredientStateInput,
                    notes: event.target.value,
                  })
                }
              />
            </Grid.Col>
            <Grid.Col span={2}>
              <Button
                mt="xl"
                disabled={
                  newIngredientForIngredientStateInput.validIngredientID === '' ||
                  newIngredientForIngredientStateInput.validIngredientID === ''
                }
                onClick={async () => {
                  await apiClient
                    .createValidIngredientStateIngredient(newIngredientForIngredientStateInput)
                    .then((res: APIResponse<ValidIngredientStateIngredient>) => {
                      // the returned value doesn't have enough information to put it in the list, so we have to fetch it
                      apiClient
                        .getValidIngredientStateIngredient(res.data.id)
                        .then((res: APIResponse<ValidIngredientStateIngredient>) => {
                          const returnedValue = res;

                          setIngredientsForIngredientState({
                            ...ingredientsForIngredientState,
                            data: [...(ingredientsForIngredientState.data || []), returnedValue.data],
                          });

                          setNewIngredientForIngredientStateInput(
                            new ValidIngredientStateIngredientCreationRequestInput({
                              validIngredientStateID: validIngredientState.id,
                              validIngredientID: '',
                              notes: '',
                            }),
                          );

                          setIngredientQuery('');
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

export default ValidIngredientStatePage;
