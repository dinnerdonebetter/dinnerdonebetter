import { AxiosError } from 'axios';
import { useForm, zodResolver } from '@mantine/form';
import {
  ActionIcon,
  Alert,
  Autocomplete,
  AutocompleteItem,
  Button,
  Center,
  Container,
  Divider,
  Grid,
  Group,
  NumberInput,
  Pagination,
  Space,
  Switch,
  Table,
  Text,
  TextInput,
  ThemeIcon,
  Title,
} from '@mantine/core';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import Link from 'next/link';
import { useRouter } from 'next/router';
import { useEffect, useState } from 'react';
import { IconTrash } from '@tabler/icons';
import { z } from 'zod';

import {
  APIResponse,
  EitherErrorOr,
  IAPIError,
  QueryFilteredResult,
  ValidIngredient,
  ValidIngredientMeasurementUnit,
  ValidIngredientMeasurementUnitCreationRequestInput,
  ValidMeasurementUnit,
  ValidMeasurementUnitConversion,
  ValidMeasurementUnitConversionCreationRequestInput,
  ValidMeasurementUnitUpdateRequestInput,
} from '@dinnerdonebetter/models';
import { ServerTiming, ServerTimingHeaderName } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { AppLayout } from '../../../src/layouts';
import { buildServerSideClientOrRedirect } from '../../../src/client';
import { serverSideTracer } from '../../../src/tracer';
import { inputSlug } from '../../../src/schemas';
import { valueOrDefault } from '../../../src/utils';

declare interface ValidMeasurementUnitPageProps {
  pageLoadValidMeasurementUnit: EitherErrorOr<ValidMeasurementUnit>;
  pageLoadIngredientsForMeasurementUnit: EitherErrorOr<QueryFilteredResult<ValidIngredientMeasurementUnit>>;
  pageLoadMeasurementUnitConversionsFromUnit: EitherErrorOr<ValidMeasurementUnitConversion[]>;
  pageLoadMeasurementUnitConversionsToUnit: EitherErrorOr<ValidMeasurementUnitConversion[]>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<ValidMeasurementUnitPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('ValidMeasurementUnitPage.getServerSideProps');

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

  const { validMeasurementUnitID } = context.query;
  if (!validMeasurementUnitID) {
    throw new Error('valid measurement unit ID is somehow missing!');
  }

  const fetchValidMeasurementUnitTimer = timing.addEvent('fetch ');
  const pageLoadValidMeasurementUnitPromise = apiClient
    .getValidMeasurementUnit(validMeasurementUnitID.toString())
    .then((result: APIResponse<ValidMeasurementUnit>) => {
      span.addEvent('valid measurement unit retrieved');
      return result.data;
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchValidMeasurementUnitTimer.end();
    });

  const fetchIngredientsForMeasurementUnitTimer = timing.addEvent('fetch valid ingredient measurement units');
  const pageLoadIngredientsForMeasurementUnitPromise = apiClient
    .getValidIngredientMeasurementUnitsByMeasurementUnit(validMeasurementUnitID.toString())
    .then((res: QueryFilteredResult<ValidIngredientMeasurementUnit>) => {
      span.addEvent('valid ingredient measurement units retrieved');
      return res;
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchIngredientsForMeasurementUnitTimer.end();
    });

  const fetchMeasurementUnitConversionsFromUnitTimer = timing.addEvent('fetch measurement unit conversions from unit');
  const pageLoadMeasurementUnitConversionsFromUnitPromise = apiClient
    .getValidMeasurementUnitConversionsFromUnit(validMeasurementUnitID.toString())
    .then((res: QueryFilteredResult<ValidMeasurementUnitConversion>) => {
      span.addEvent('valid ingredient measurement units retrieved');
      return res.data;
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchMeasurementUnitConversionsFromUnitTimer.end();
    });

  const fetchMeasurementUnitConversionsToUnitTimer = timing.addEvent('fetch measurement unit conversions to unit');
  const pageLoadMeasurementUnitConversionsToUnitPromise = apiClient
    .getValidMeasurementUnitConversionsToUnit(validMeasurementUnitID.toString())
    .then((res: QueryFilteredResult<ValidMeasurementUnitConversion>) => {
      span.addEvent('valid ingredient measurement units retrieved');
      return res.data;
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchMeasurementUnitConversionsToUnitTimer.end();
    });

  const [
    pageLoadValidMeasurementUnit,
    pageLoadIngredientsForMeasurementUnit,
    pageLoadMeasurementUnitConversionsFromUnit,
    pageLoadMeasurementUnitConversionsToUnit,
  ] = await Promise.all([
    pageLoadValidMeasurementUnitPromise,
    pageLoadIngredientsForMeasurementUnitPromise,
    pageLoadMeasurementUnitConversionsFromUnitPromise,
    pageLoadMeasurementUnitConversionsToUnitPromise,
  ]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  span.end();
  return {
    props: {
      pageLoadValidMeasurementUnit: JSON.parse(JSON.stringify(pageLoadValidMeasurementUnit)),
      pageLoadIngredientsForMeasurementUnit: JSON.parse(JSON.stringify(pageLoadIngredientsForMeasurementUnit)),
      pageLoadMeasurementUnitConversionsFromUnit: JSON.parse(
        JSON.stringify(pageLoadMeasurementUnitConversionsFromUnit),
      ),
      pageLoadMeasurementUnitConversionsToUnit: JSON.parse(JSON.stringify(pageLoadMeasurementUnitConversionsToUnit)),
    },
  };
};

const validMeasurementUnitUpdateFormSchema = z.object({
  name: z.string().trim().min(1, 'name is required'),
  pluralName: z.string().trim().min(1, 'plural name is required'),
  slug: inputSlug,
});

function ValidMeasurementUnitPage(props: ValidMeasurementUnitPageProps) {
  const router = useRouter();

  const apiClient = buildLocalClient();
  const {
    pageLoadValidMeasurementUnit,
    pageLoadIngredientsForMeasurementUnit,
    pageLoadMeasurementUnitConversionsFromUnit,
    pageLoadMeasurementUnitConversionsToUnit,
  } = props;

  const ogValidMeasurementUnit = valueOrDefault(pageLoadValidMeasurementUnit, new ValidMeasurementUnit());
  const [validMeasurementUnit, setValidMeasurementUnit] = useState<ValidMeasurementUnit>(ogValidMeasurementUnit);
  const [originalValidMeasurementUnit, setOriginalValidMeasurementUnit] =
    useState<ValidMeasurementUnit>(ogValidMeasurementUnit);
  const [validMeasurementUnitError] = useState<IAPIError | undefined>(pageLoadValidMeasurementUnit.error);

  const [newIngredientForMeasurementUnitInput, setNewIngredientForMeasurementUnitInput] =
    useState<ValidIngredientMeasurementUnitCreationRequestInput>(
      new ValidIngredientMeasurementUnitCreationRequestInput({
        validMeasurementUnitID: validMeasurementUnit.id,
        allowableQuantity: { min: 1 },
      }),
    );
  const [ingredientQuery, setIngredientQuery] = useState('');

  const ogValidIngredient = valueOrDefault(
    pageLoadIngredientsForMeasurementUnit,
    new QueryFilteredResult<ValidIngredientMeasurementUnit>(),
  );
  const [ingredientsForMeasurementUnitError] = useState<IAPIError | undefined>(
    pageLoadIngredientsForMeasurementUnit.error,
  );
  const [ingredientsForMeasurementUnit, setIngredientsForMeasurementUnit] =
    useState<QueryFilteredResult<ValidIngredientMeasurementUnit>>(ogValidIngredient);
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
          return !(ingredientsForMeasurementUnit.data || []).some((vimu: ValidIngredientMeasurementUnit) => {
            return vimu.ingredient.id === mu.id;
          });
        });

        setSuggestedIngredients(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [ingredientQuery, ingredientsForMeasurementUnit.data]);

  const [newMeasurementUnitConversionFromMeasurementUnit, setNewMeasurementUnitConversionFromMeasurementUnit] =
    useState<ValidMeasurementUnitConversionCreationRequestInput>(
      new ValidMeasurementUnitConversionCreationRequestInput({
        from: validMeasurementUnit.id,
        modifier: 1,
      }),
    );
  const [conversionFromUnitQuery, setConversionFromUnitQuery] = useState('');

  const ogMeasurementUnitConversionsFrom = valueOrDefault(
    pageLoadMeasurementUnitConversionsFromUnit,
    new Array<ValidMeasurementUnitConversion>(),
  );
  const [measurementUnitsToConvertFromError] = useState<IAPIError | undefined>(
    pageLoadMeasurementUnitConversionsFromUnit.error,
  );
  const [measurementUnitsToConvertFrom, setMeasurementUnitsToConvertFrom] = useState<ValidMeasurementUnitConversion[]>(
    ogMeasurementUnitConversionsFrom,
  );
  const [suggestedMeasurementUnitsToConvertFrom, setSuggestedMeasurementUnitsToConvertFrom] = useState<
    ValidMeasurementUnit[]
  >([]);
  const [conversionFromOnlyIngredientQuery, setConversionFromOnlyIngredientQuery] = useState('');
  const [suggestedIngredientsToRestrictConversionFrom, setSuggestedIngredientsToRestrictConversionFrom] = useState<
    ValidIngredient[]
  >([]);

  useEffect(() => {
    if (conversionFromUnitQuery.length <= 2) {
      setSuggestedMeasurementUnitsToConvertFrom([]);
      return;
    }

    const apiClient = buildLocalClient();
    apiClient
      .searchForValidMeasurementUnits(conversionFromUnitQuery)
      .then((res: QueryFilteredResult<ValidMeasurementUnit>) => {
        const newSuggestions = (res.data || []).filter((mu: ValidMeasurementUnit) => {
          return mu.id != validMeasurementUnit.id;
        });

        setSuggestedMeasurementUnitsToConvertFrom(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [conversionFromUnitQuery, validMeasurementUnit.id]);

  const [newMeasurementUnitConversionToMeasurementUnit, setNewMeasurementUnitConversionToMeasurementUnit] =
    useState<ValidMeasurementUnitConversionCreationRequestInput>(
      new ValidMeasurementUnitConversionCreationRequestInput({
        to: validMeasurementUnit.id,
        modifier: 1,
      }),
    );
  const [conversionToUnitQuery, setConversionToUnitQuery] = useState('');

  const ogMeasurementUnitConversionsTo = valueOrDefault(
    pageLoadMeasurementUnitConversionsToUnit,
    new Array<ValidMeasurementUnitConversion>(),
  );
  const [measurementUnitsToConvertToError] = useState<IAPIError | undefined>(
    pageLoadMeasurementUnitConversionsToUnit.error,
  );
  const [measurementUnitsToConvertTo, setMeasurementUnitsToConvertTo] =
    useState<ValidMeasurementUnitConversion[]>(ogMeasurementUnitConversionsTo);
  const [suggestedMeasurementUnitsToConvertTo, setSuggestedMeasurementUnitsToConvertTo] = useState<
    ValidMeasurementUnit[]
  >([]);
  const [conversionToOnlyIngredientQuery, setConversionToOnlyIngredientQuery] = useState('');
  const [suggestedIngredientsToRestrictConversionTo, setSuggestedIngredientsToRestrictConversionTo] = useState<
    ValidIngredient[]
  >([]);

  useEffect(() => {
    if (conversionToUnitQuery.length <= 2) {
      setSuggestedMeasurementUnitsToConvertTo([]);
      return;
    }

    const apiClient = buildLocalClient();
    apiClient
      .searchForValidMeasurementUnits(conversionToUnitQuery)
      .then((res: QueryFilteredResult<ValidMeasurementUnit>) => {
        const newSuggestions = (res.data || []).filter((mu: ValidMeasurementUnit) => {
          return mu.id != validMeasurementUnit.id;
        });

        setSuggestedMeasurementUnitsToConvertTo(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [conversionToUnitQuery, validMeasurementUnit.id]);

  useEffect(() => {
    if (conversionFromOnlyIngredientQuery.length <= 2) {
      setSuggestedIngredientsToRestrictConversionFrom([]);
      return;
    }

    const apiClient = buildLocalClient();
    apiClient
      .searchForValidIngredients(conversionFromOnlyIngredientQuery)
      .then((res: QueryFilteredResult<ValidIngredient>) => {
        const newSuggestions = (res.data || []).filter((mu: ValidIngredient) => {
          return !(ingredientsForMeasurementUnit.data || []).some((vimu: ValidIngredientMeasurementUnit) => {
            return vimu.ingredient.id === mu.id;
          });
        });

        setSuggestedIngredientsToRestrictConversionFrom(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [conversionFromOnlyIngredientQuery, ingredientsForMeasurementUnit.data]);

  useEffect(() => {
    if (conversionToOnlyIngredientQuery.length <= 2) {
      setSuggestedIngredientsToRestrictConversionTo([]);
      return;
    }

    const apiClient = buildLocalClient();
    apiClient
      .searchForValidIngredients(conversionToOnlyIngredientQuery)
      .then((res: QueryFilteredResult<ValidIngredient>) => {
        const newSuggestions = (res.data || []).filter((mu: ValidIngredient) => {
          return !(ingredientsForMeasurementUnit.data || []).some((vimu: ValidIngredientMeasurementUnit) => {
            return vimu.ingredient.id === mu.id;
          });
        });

        setSuggestedIngredientsToRestrictConversionTo(newSuggestions);
      })
      .catch((err: AxiosError) => {
        console.error(err);
      });
  }, [conversionToOnlyIngredientQuery, ingredientsForMeasurementUnit.data]);

  const updateForm = useForm({
    initialValues: validMeasurementUnit,
    validate: zodResolver(validMeasurementUnitUpdateFormSchema),
  });

  const dataHasChanged = (): boolean => {
    return (
      originalValidMeasurementUnit.name !== updateForm.values.name ||
      originalValidMeasurementUnit.description !== updateForm.values.description ||
      originalValidMeasurementUnit.pluralName !== updateForm.values.pluralName ||
      originalValidMeasurementUnit.universal !== updateForm.values.universal ||
      originalValidMeasurementUnit.volumetric !== updateForm.values.volumetric ||
      originalValidMeasurementUnit.metric !== updateForm.values.metric ||
      originalValidMeasurementUnit.imperial !== updateForm.values.imperial ||
      originalValidMeasurementUnit.slug !== updateForm.values.slug
    );
  };

  const submit = async () => {
    const validation = updateForm.validate();
    if (validation.hasErrors) {
      console.error(validation.errors);
      return;
    }

    const submission = new ValidMeasurementUnitUpdateRequestInput({
      name: updateForm.values.name,
      description: updateForm.values.description,
      universal: updateForm.values.universal,
      metric: updateForm.values.metric,
      volumetric: updateForm.values.volumetric,
      imperial: updateForm.values.imperial,
      pluralName: updateForm.values.pluralName,
      slug: updateForm.values.slug,
    });

    if (updateForm.values.imperial && updateForm.values.metric) {
      updateForm.setErrors({
        metric: 'Cannot be imperial and metric at the same time',
        imperial: 'Cannot be imperial and metric at the same time',
      });
      return;
    }

    const apiClient = buildLocalClient();

    await apiClient
      .updateValidMeasurementUnit(validMeasurementUnit.id, submission)
      .then((result: APIResponse<ValidMeasurementUnit>) => {
        if (result) {
          updateForm.setValues(result.data);
          setValidMeasurementUnit(result.data);
          setOriginalValidMeasurementUnit(result.data);
        }
      })
      .catch((err) => {
        console.error(err);
      });
  };

  return (
    <AppLayout title="Valid Measurement Unit">
      <Container size="sm">
        {validMeasurementUnitError && <Alert color="tomato">{validMeasurementUnitError.message}</Alert>}

        {!validMeasurementUnitError && (
          <form onSubmit={updateForm.onSubmit(submit)}>
            <TextInput label="Name" placeholder="thing" {...updateForm.getInputProps('name')} />
            <TextInput label="Plural Name" placeholder="things" {...updateForm.getInputProps('pluralName')} />
            <TextInput label="Slug" placeholder="thing" {...updateForm.getInputProps('slug')} />
            <TextInput label="Description" placeholder="thing" {...updateForm.getInputProps('description')} />

            <Switch
              checked={updateForm.values.volumetric}
              label="Volumetric"
              {...updateForm.getInputProps('volumetric')}
            />
            <Switch
              checked={updateForm.values.universal}
              label="Universal"
              {...updateForm.getInputProps('universal')}
            />
            <Switch
              checked={updateForm.values.metric}
              disabled={updateForm.values.imperial}
              label="Metric"
              {...updateForm.getInputProps('metric')}
            />
            <Switch
              checked={updateForm.values.imperial}
              disabled={updateForm.values.metric}
              label="Imperial"
              {...updateForm.getInputProps('imperial')}
            />

            <Group position="center">
              <Button type="submit" mt="sm" fullWidth disabled={!dataHasChanged()}>
                Submit
              </Button>
              <Button
                type="submit"
                color="tomato"
                fullWidth
                onClick={() => {
                  if (confirm('Are you sure you want to delete this valid measurement unit?')) {
                    apiClient.archiveValidMeasurementUnit(validMeasurementUnit.id).then(() => {
                      router.push('/valid_measurement_units');
                    });
                  }
                }}
              >
                Delete
              </Button>
            </Group>
          </form>
        )}

        {!validMeasurementUnit.universal && (
          <>
            {/*

            INGREDIENTS

            */}

            <Space h="xl" />
            <Divider />
            <Space h="xl" />

            <form>
              <Center>
                <Title order={4}>Ingredients</Title>
              </Center>
              {ingredientsForMeasurementUnitError && (
                <Text color="tomato">
                  Error fetching ingredients for measurement unit: {ingredientsForMeasurementUnitError.message}
                </Text>
              )}

              {!ingredientsForMeasurementUnitError &&
                ingredientsForMeasurementUnit.data &&
                (ingredientsForMeasurementUnit.data || []).length !== 0 && (
                  <>
                    <Table mt="xl" withColumnBorders>
                      <thead>
                        <tr>
                          <th>Name</th>
                          <th>Min Qty</th>
                          <th>Max Qty</th>
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
                        {(ingredientsForMeasurementUnit.data || []).map(
                          (validIngredientMeasurementUnit: ValidIngredientMeasurementUnit) => {
                            return (
                              <tr key={validIngredientMeasurementUnit.id}>
                                <td>
                                  <Link href={`/valid_ingredients/${validIngredientMeasurementUnit.ingredient.id}`}>
                                    {validIngredientMeasurementUnit.ingredient.name}
                                  </Link>
                                </td>
                                <td>
                                  <Text>{validIngredientMeasurementUnit.allowableQuantity.min}</Text>
                                </td>
                                <td>
                                  <Text>{validIngredientMeasurementUnit.allowableQuantity.max}</Text>
                                </td>
                                <td>
                                  <Text>{validIngredientMeasurementUnit.notes}</Text>
                                </td>
                                <td>
                                  <Center>
                                    <ActionIcon
                                      variant="outline"
                                      aria-label="remove valid ingredient measurement unit"
                                      onClick={async () => {
                                        await apiClient
                                          .archiveValidIngredientMeasurementUnit(validIngredientMeasurementUnit.id)
                                          .then(() => {
                                            setIngredientsForMeasurementUnit({
                                              ...ingredientsForMeasurementUnit,
                                              data: ingredientsForMeasurementUnit.data.filter(
                                                (x: ValidIngredientMeasurementUnit) =>
                                                  x.id !== validIngredientMeasurementUnit.id,
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
                        Math.ceil(ingredientsForMeasurementUnit.totalCount / ingredientsForMeasurementUnit.limit) <=
                        ingredientsForMeasurementUnit.page
                      }
                      position="center"
                      page={ingredientsForMeasurementUnit.page}
                      total={Math.ceil(ingredientsForMeasurementUnit.totalCount / ingredientsForMeasurementUnit.limit)}
                      onChange={(value: number) => {
                        setIngredientsForMeasurementUnit({ ...ingredientsForMeasurementUnit, page: value });
                      }}
                    />
                  </>
                )}

              <Grid>
                <Grid.Col span="auto">
                  <Autocomplete
                    placeholder="garlic"
                    label="Ingredient"
                    value={ingredientQuery}
                    onChange={setIngredientQuery}
                    onItemSubmit={async (item: AutocompleteItem) => {
                      const selectedIngredient = suggestedIngredients.find(
                        (x: ValidIngredient) => x.name === item.value,
                      );

                      if (!selectedIngredient) {
                        console.error(`selectedIngredient not found for item ${item.value}}`);
                        return;
                      }

                      setNewIngredientForMeasurementUnitInput({
                        ...newIngredientForMeasurementUnitInput,
                        validIngredientID: selectedIngredient.id,
                      });
                    }}
                    data={suggestedIngredients.map((x: ValidIngredient) => {
                      return { value: x.name, label: x.name };
                    })}
                  />
                </Grid.Col>
                <Grid.Col span={2}>
                  <NumberInput
                    value={newIngredientForMeasurementUnitInput.allowableQuantity.min}
                    label="Min. Qty"
                    onChange={(value: number) =>
                      setNewIngredientForMeasurementUnitInput({
                        ...newIngredientForMeasurementUnitInput,
                        allowableQuantity: {
                          ...newIngredientForMeasurementUnitInput.allowableQuantity,
                          min: value,
                        },
                      })
                    }
                  />
                </Grid.Col>
                <Grid.Col span={2}>
                  <NumberInput
                    value={newIngredientForMeasurementUnitInput.allowableQuantity.max}
                    label="Max. Qty"
                    onChange={(value: number) =>
                      setNewIngredientForMeasurementUnitInput({
                        ...newIngredientForMeasurementUnitInput,
                        allowableQuantity: {
                          ...newIngredientForMeasurementUnitInput.allowableQuantity,
                          max: value,
                        },
                      })
                    }
                  />
                </Grid.Col>
                <Grid.Col span="auto">
                  <TextInput
                    label="Notes"
                    value={newIngredientForMeasurementUnitInput.notes}
                    onChange={(event) =>
                      setNewIngredientForMeasurementUnitInput({
                        ...newIngredientForMeasurementUnitInput,
                        notes: event.target.value,
                      })
                    }
                  />
                </Grid.Col>
                <Grid.Col span={2}>
                  <Button
                    mt="xl"
                    disabled={
                      newIngredientForMeasurementUnitInput.validMeasurementUnitID === '' ||
                      newIngredientForMeasurementUnitInput.validIngredientID === ''
                    }
                    onClick={async () => {
                      await apiClient
                        .createValidIngredientMeasurementUnit(newIngredientForMeasurementUnitInput)
                        .then((res: APIResponse<ValidIngredientMeasurementUnit>) => {
                          // the returned value doesn't have enough information to put it in the list, so we have to fetch it
                          apiClient
                            .getValidIngredientMeasurementUnit(res.data.id)
                            .then((res: APIResponse<ValidIngredientMeasurementUnit>) => {
                              setIngredientsForMeasurementUnit({
                                ...ingredientsForMeasurementUnit,
                                data: [...(ingredientsForMeasurementUnit.data || []), res.data],
                              });

                              setNewIngredientForMeasurementUnitInput(
                                new ValidIngredientMeasurementUnitCreationRequestInput({
                                  validMeasurementUnitID: validMeasurementUnit.id,
                                  allowableQuantity: { min: 1 },
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
          </>
        )}

        {/*

        CONVERSIONS FROM THIS MEASUREMENT UNIT

        */}

        <Space h="xl" />
        <Divider />
        <Space h="xl" />

        <form>
          <Center>
            <Title order={4}>Conversions From</Title>
          </Center>

          {measurementUnitsToConvertFromError && (
            <Text color="tomato">Error fetching conversion units: {measurementUnitsToConvertFromError.message}</Text>
          )}

          {!measurementUnitsToConvertFromError &&
            measurementUnitsToConvertFrom &&
            (measurementUnitsToConvertFrom || []).length !== 0 && (
              <>
                <Table mt="xl" withColumnBorders>
                  <thead>
                    <tr>
                      <th>From</th>
                      <th>To</th>
                      <th>Modifier</th>
                      <th>Only For Ingredient</th>
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
                    {(measurementUnitsToConvertFrom || []).map(
                      (validMeasurementUnitConversion: ValidMeasurementUnitConversion) => {
                        return (
                          <tr key={validMeasurementUnitConversion.id}>
                            <td>
                              <Link href={`/valid_measurement_units/${validMeasurementUnitConversion.from.id}`}>
                                {validMeasurementUnitConversion.from.pluralName}
                              </Link>
                            </td>
                            <td>
                              <Link href={`/valid_measurement_units/${validMeasurementUnitConversion.to.id}`}>
                                {validMeasurementUnitConversion.to.pluralName}
                              </Link>
                            </td>
                            <td>
                              <Text>{validMeasurementUnitConversion.modifier}</Text>
                            </td>
                            <td>
                              {(validMeasurementUnitConversion.onlyForIngredient && (
                                <Link
                                  href={`/valid_ingredients/${validMeasurementUnitConversion.onlyForIngredient.id}`}
                                >
                                  {validMeasurementUnitConversion.to.pluralName}
                                </Link>
                              )) || <Text> - </Text>}
                            </td>
                            <td>
                              <Text>{validMeasurementUnitConversion.notes}</Text>
                            </td>
                            <td>
                              <Center>
                                <ActionIcon
                                  variant="outline"
                                  aria-label="remove valid ingredient measurement unit"
                                  onClick={async () => {
                                    await apiClient
                                      .archiveValidMeasurementUnitConversion(validMeasurementUnitConversion.id)
                                      .then(() => {
                                        setMeasurementUnitsToConvertFrom([
                                          ...measurementUnitsToConvertFrom.filter(
                                            (x: ValidMeasurementUnitConversion) =>
                                              x.id !== validMeasurementUnitConversion.id,
                                          ),
                                        ]);
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

                {/*
              <Pagination
                disabled={
                  Math.ceil(measurementUnitsToConvertFrom.totalCount / measurementUnitsToConvertFrom.limit) <=
                  measurementUnitsToConvertFrom.page
                }
                position="center"
                page={measurementUnitsToConvertFrom.page}
                total={Math.ceil(measurementUnitsToConvertFrom.totalCount / measurementUnitsToConvertFrom.limit)}
                onChange={(value: number) => {
                  setMeasurementUnitsToConvertFrom({ ...measurementUnitsToConvertFrom, page: value });
                }}
              />
              */}
              </>
            )}

          <Grid>
            <Grid.Col span="auto">
              <Autocomplete
                placeholder="grams"
                label="Measurement Unit"
                value={conversionFromUnitQuery}
                onChange={setConversionFromUnitQuery}
                onItemSubmit={async (item: AutocompleteItem) => {
                  const selectedMeasurementUnit = suggestedMeasurementUnitsToConvertFrom.find(
                    (x: ValidMeasurementUnit) => x.name === item.value,
                  );

                  if (!selectedMeasurementUnit) {
                    console.error(`selectedMeasurementUnit not found for item ${item.value}}`);
                    return;
                  }

                  setNewMeasurementUnitConversionFromMeasurementUnit({
                    ...newMeasurementUnitConversionFromMeasurementUnit,
                    to: selectedMeasurementUnit.id,
                  });
                }}
                data={suggestedMeasurementUnitsToConvertFrom.map((x: ValidMeasurementUnit) => {
                  return { value: x.name, label: x.name };
                })}
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <NumberInput
                label="Modifier"
                step={0.00001}
                precision={5}
                value={newMeasurementUnitConversionFromMeasurementUnit.modifier}
                onChange={(value: number) =>
                  setNewMeasurementUnitConversionFromMeasurementUnit({
                    ...newMeasurementUnitConversionFromMeasurementUnit,
                    modifier: value,
                  })
                }
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <Autocomplete
                placeholder="garlic"
                label="Only For Ingredient"
                value={conversionFromOnlyIngredientQuery}
                onChange={setConversionFromOnlyIngredientQuery}
                onItemSubmit={async (item: AutocompleteItem) => {
                  const selectedValidIngredient = suggestedIngredientsToRestrictConversionFrom.find(
                    (x: ValidIngredient) => x.name === item.value,
                  );

                  if (!selectedValidIngredient) {
                    console.error(`selectedValidIngredient not found for item ${item.value}}`);
                    return;
                  }

                  setNewMeasurementUnitConversionFromMeasurementUnit({
                    ...newMeasurementUnitConversionFromMeasurementUnit,
                    onlyForIngredient: selectedValidIngredient.id,
                  });
                }}
                data={suggestedIngredientsToRestrictConversionFrom.map((x: ValidIngredient) => {
                  return { value: x.name, label: x.name };
                })}
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <TextInput
                label="Notes"
                value={newMeasurementUnitConversionFromMeasurementUnit.notes}
                onChange={(event) =>
                  setNewMeasurementUnitConversionFromMeasurementUnit({
                    ...newMeasurementUnitConversionFromMeasurementUnit,
                    notes: event.target.value,
                  })
                }
              />
            </Grid.Col>
            <Grid.Col span={2}>
              <Button
                mt="xl"
                disabled={
                  newMeasurementUnitConversionFromMeasurementUnit.from === '' ||
                  newMeasurementUnitConversionFromMeasurementUnit.to === ''
                }
                onClick={async () => {
                  await apiClient
                    .createValidMeasurementUnitConversion(newMeasurementUnitConversionFromMeasurementUnit)
                    .then((res: APIResponse<ValidMeasurementUnitConversion>) => {
                      // the returned value doesn't have enough information to put it in the list, so we have to fetch it
                      apiClient
                        .getValidMeasurementUnitConversion(res.data.id)
                        .then((res: APIResponse<ValidMeasurementUnitConversion>) => {
                          setMeasurementUnitsToConvertFrom([...(measurementUnitsToConvertFrom || []), res.data]);

                          setNewMeasurementUnitConversionFromMeasurementUnit(
                            new ValidMeasurementUnitConversionCreationRequestInput({
                              from: validMeasurementUnit.id,
                              to: '',
                              notes: '',
                              modifier: 1,
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

        {/*

        CONVERSIONS FROM THIS MEASUREMENT UNIT

        */}

        <Space h="xl" />
        <Divider />
        <Space h="xl" />

        <form>
          <Center>
            <Title order={4}>Conversions To</Title>
          </Center>

          {measurementUnitsToConvertToError && (
            <Text color="tomato">Error fetching conversion units: {measurementUnitsToConvertToError.message}</Text>
          )}

          {!measurementUnitsToConvertToError &&
            measurementUnitsToConvertTo &&
            (measurementUnitsToConvertTo || []).length !== 0 && (
              <>
                <Table mt="xl" withColumnBorders>
                  <thead>
                    <tr>
                      <th>From</th>
                      <th>To</th>
                      <th>Modifier</th>
                      <th>Only for Ingredient</th>
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
                    {(measurementUnitsToConvertTo || []).map(
                      (validMeasurementUnitConversion: ValidMeasurementUnitConversion) => {
                        return (
                          <tr key={validMeasurementUnitConversion.id}>
                            <td>
                              <Link href={`/valid_measurement_units/${validMeasurementUnitConversion.from.id}`}>
                                {validMeasurementUnitConversion.from.pluralName}
                              </Link>
                            </td>
                            <td>
                              <Link href={`/valid_measurement_units/${validMeasurementUnitConversion.to.id}`}>
                                {validMeasurementUnitConversion.to.pluralName}
                              </Link>
                            </td>
                            <td>
                              <Text>{validMeasurementUnitConversion.modifier}</Text>
                            </td>
                            <td>
                              {(validMeasurementUnitConversion.onlyForIngredient && (
                                <Link
                                  href={`/valid_ingredients/${validMeasurementUnitConversion.onlyForIngredient.id}`}
                                >
                                  {validMeasurementUnitConversion.to.pluralName}
                                </Link>
                              )) || <Text> - </Text>}
                            </td>
                            <td>
                              <Text>{validMeasurementUnitConversion.notes}</Text>
                            </td>
                            <td>
                              <Center>
                                <ActionIcon
                                  variant="outline"
                                  aria-label="remove valid ingredient measurement unit"
                                  onClick={async () => {
                                    await apiClient
                                      .archiveValidMeasurementUnitConversion(validMeasurementUnitConversion.id)
                                      .then(() => {
                                        setMeasurementUnitsToConvertTo([
                                          ...measurementUnitsToConvertTo.filter(
                                            (x: ValidMeasurementUnitConversion) =>
                                              x.id !== validMeasurementUnitConversion.id,
                                          ),
                                        ]);
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

                {/*
              <Pagination
                disabled={
                  Math.ceil(measurementUnitsToConvertTo.totalCount / measurementUnitsToConvertTo.limit) <=
                  measurementUnitsToConvertTo.page
                }
                position="center"
                page={measurementUnitsToConvertTo.page}
                total={Math.ceil(measurementUnitsToConvertTo.totalCount / measurementUnitsToConvertTo.limit)}
                onChange={(value: number) => {
                  setMeasurementUnitsToConvertTo({ ...measurementUnitsToConvertTo, page: value });
                }}
              />
              */}
              </>
            )}

          <Grid>
            <Grid.Col span="auto">
              <Autocomplete
                placeholder="grams"
                label="Measurement Unit"
                value={conversionToUnitQuery}
                onChange={setConversionToUnitQuery}
                onItemSubmit={async (item: AutocompleteItem) => {
                  const selectedMeasurementUnit = suggestedMeasurementUnitsToConvertTo.find(
                    (x: ValidMeasurementUnit) => x.name === item.value,
                  );

                  if (!selectedMeasurementUnit) {
                    console.error(`selectedMeasurementUnit not found for item ${item.value}}`);
                    return;
                  }

                  setNewMeasurementUnitConversionToMeasurementUnit({
                    ...newMeasurementUnitConversionToMeasurementUnit,
                    from: selectedMeasurementUnit.id,
                  });
                }}
                data={suggestedMeasurementUnitsToConvertTo.map((x: ValidMeasurementUnit) => {
                  return { value: x.name, label: x.name };
                })}
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <NumberInput
                label="Modifier"
                step={0.00001}
                precision={5}
                value={newMeasurementUnitConversionToMeasurementUnit.modifier}
                onChange={(value: number) =>
                  setNewMeasurementUnitConversionToMeasurementUnit({
                    ...newMeasurementUnitConversionToMeasurementUnit,
                    modifier: value,
                  })
                }
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <Autocomplete
                placeholder="garlic"
                label="Only For Ingredient"
                value={conversionToOnlyIngredientQuery}
                onChange={setConversionToOnlyIngredientQuery}
                onItemSubmit={async (item: AutocompleteItem) => {
                  const selectedValidIngredient = suggestedIngredientsToRestrictConversionTo.find(
                    (x: ValidIngredient) => x.name === item.value,
                  );

                  if (!selectedValidIngredient) {
                    console.error(`selectedValidIngredient not found for item ${item.value}}`);
                    return;
                  }

                  setNewMeasurementUnitConversionToMeasurementUnit({
                    ...newMeasurementUnitConversionToMeasurementUnit,
                    onlyForIngredient: selectedValidIngredient.id,
                  });
                }}
                data={suggestedIngredientsToRestrictConversionTo.map((x: ValidIngredient) => {
                  return { value: x.name, label: x.name };
                })}
              />
            </Grid.Col>
            <Grid.Col span="auto">
              <TextInput
                label="Notes"
                value={newMeasurementUnitConversionToMeasurementUnit.notes}
                onChange={(event) =>
                  setNewMeasurementUnitConversionToMeasurementUnit({
                    ...newMeasurementUnitConversionToMeasurementUnit,
                    notes: event.target.value,
                  })
                }
              />
            </Grid.Col>
            <Grid.Col span={2}>
              <Button
                mt="xl"
                disabled={
                  newMeasurementUnitConversionToMeasurementUnit.from === '' ||
                  newMeasurementUnitConversionToMeasurementUnit.to === ''
                }
                onClick={async () => {
                  await apiClient
                    .createValidMeasurementUnitConversion(newMeasurementUnitConversionToMeasurementUnit)
                    .then((res: APIResponse<ValidMeasurementUnitConversion>) => {
                      // the returned value doesn't have enough information to put it in the list, so we have to fetch it
                      apiClient
                        .getValidMeasurementUnitConversion(res.data.id)
                        .then((res: APIResponse<ValidMeasurementUnitConversion>) => {
                          setMeasurementUnitsToConvertTo([...(measurementUnitsToConvertTo || []), res.data]);

                          setNewMeasurementUnitConversionToMeasurementUnit(
                            new ValidMeasurementUnitConversionCreationRequestInput({
                              to: validMeasurementUnit.id,
                              from: '',
                              notes: '',
                              modifier: 1,
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

export default ValidMeasurementUnitPage;
