import { format, formatDuration, subSeconds, intervalToDuration } from 'date-fns';
import { GetServerSideProps, GetServerSidePropsContext, GetServerSidePropsResult } from 'next';
import {
  Title,
  Grid,
  Center,
  Button,
  Divider,
  Card,
  Stack,
  ActionIcon,
  Indicator,
  Text,
  List,
  Box,
  Table,
  NumberInput,
  Tooltip,
  Badge,
  Avatar,
} from '@mantine/core';
import Link from 'next/link';
import router from 'next/router';
import { Reducer, useReducer, useState } from 'react';
import { IconCheck, IconCircleX, IconThumbUp, IconTrash } from '@tabler/icons';

import {
  APIResponse,
  EitherErrorOr,
  Household,
  HouseholdUserMembershipWithUser,
  IAPIError,
  MealComponent,
  MealPlan,
  MealPlanEvent,
  MealPlanGroceryListItem,
  MealPlanGroceryListItemUpdateRequestInput,
  MealPlanOption,
  MealPlanOptionVote,
  MealPlanTask,
  MealPlanTaskStatusChangeRequestInput,
  QueryFilteredResult,
  Recipe,
  RecipePrepTaskStep,
  RecipeStep,
  RecipeStepIngredient,
} from '@dinnerdonebetter/models';
import { getEarliestEvent, getLatestEvent } from '@dinnerdonebetter/utils';
import { ServerTimingHeaderName, ServerTiming } from '@dinnerdonebetter/server-timing';
import { buildLocalClient } from '@dinnerdonebetter/api-client';

import { buildServerSideClientOrRedirect } from '../../../src/client';
import { AppLayout } from '../../../src/layouts';
import { serverSideTracer } from '../../../src/tracer';
import { serverSideAnalytics } from '../../../src/analytics';
import { userSessionDetailsOrRedirect } from '../../../src/auth';
import { valueOrDefault } from '../../../src/utils';

declare interface MealPlanPageProps {
  mealPlan: EitherErrorOr<MealPlan>;
  userID: string;
  household: EitherErrorOr<Household>;
  groceryList: EitherErrorOr<MealPlanGroceryListItem[]>;
  tasks: EitherErrorOr<MealPlanTask[]>;
}

export const getServerSideProps: GetServerSideProps = async (
  context: GetServerSidePropsContext,
): Promise<GetServerSidePropsResult<MealPlanPageProps>> => {
  const timing = new ServerTiming();
  const span = serverSideTracer.startSpan('MealPlanPage.getServerSideProps');

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

  const { mealPlanID: mealPlanIDParam } = context.query;
  if (!mealPlanIDParam) {
    throw new Error('meal plan ID is somehow missing!');
  }

  const mealPlanID = mealPlanIDParam.toString();

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
    serverSideAnalytics.page(userSessionData.userID, 'MEAL_PLAN_PAGE', context, {
      mealPlanID,
      householdID: userSessionData.householdID,
    });
    analyticsTimer.end();
  } else {
    console.log(`no user session data found for ${context.req.url}`);
  }

  const fetchMealPlanTimer = timing.addEvent('fetch meal plan');
  const mealPlanPromise = apiClient
    .getMealPlan(mealPlanID)
    .then((result: APIResponse<MealPlan>) => {
      span.addEvent(`meal plan retrieved`);
      return { data: result.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchMealPlanTimer.end();
    });

  const fetchHouseholdTimer = timing.addEvent('fetch household');
  const householdPromise = apiClient
    .getActiveHousehold()
    .then((result: APIResponse<Household>) => {
      span.addEvent(`household retrieved`);
      return { data: result.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchHouseholdTimer.end();
    });

  const fetchMealPlanTasksTimer = timing.addEvent('fetch meal plan tasks');
  const tasksPromise = apiClient
    .getMealPlanTasks(mealPlanID)
    .then((result: QueryFilteredResult<MealPlanTask>) => {
      span.addEvent('meal plan grocery list items retrieved');
      return { data: result.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchMealPlanTasksTimer.end();
    });

  const fetchMealPlanGroceryListItemsTimer = timing.addEvent('fetch meal plan grocery list items');
  const groceryListPromise = apiClient
    .getMealPlanGroceryListItemsForMealPlan(mealPlanID)
    .then((result: QueryFilteredResult<MealPlanGroceryListItem>) => {
      span.addEvent('meal plan grocery list items retrieved');
      return { data: result.data };
    })
    .catch((error: IAPIError) => {
      span.addEvent('error occurred', { error: error.message });
      return { error };
    })
    .finally(() => {
      fetchMealPlanGroceryListItemsTimer.end();
    });

  const retrievedData = await Promise.all([mealPlanPromise, householdPromise, groceryListPromise, tasksPromise]);

  context.res.setHeader(ServerTimingHeaderName, timing.headerValue());

  const [mealPlan, household, groceryList, tasks] = retrievedData;

  span.end();

  return {
    props: {
      mealPlan: mealPlan!,
      household: household!,
      userID: userSessionData?.userID || '',
      tasks: tasks,
      groceryList: groceryList || [],
    },
  };
};

const dateFormat = "h aa 'on' iiii',' M/d";

/* BEGIN Meal Plan Creation Reducer */

type mealPlanPageAction =
  | { type: 'MOVE_OPTION'; eventIndex: number; optionIndex: number; direction: 'up' | 'down' }
  | { type: 'ADD_VOTES_TO_MEAL_PLAN'; eventIndex: number; votes: MealPlanOptionVote[] }
  | { type: 'UPDATE_MEAL_PLAN_GROCERY_LIST_ITEM'; newItem: MealPlanGroceryListItem }
  | { type: 'UPDATE_MEAL_PLAN_TASK'; newTask: MealPlanTask };

export class MealPlanPageState {
  mealPlan: MealPlan = new MealPlan();
  groceryList: MealPlanGroceryListItem[] = [];
  tasks: MealPlanTask[] = [];

  constructor(mealPlan: MealPlan, groceryList: MealPlanGroceryListItem[], tasks: MealPlanTask[]) {
    this.mealPlan = mealPlan;
    this.groceryList = groceryList;
    this.tasks = tasks;
  }
}

const useMealPlanReducer: Reducer<MealPlanPageState, mealPlanPageAction> = (
  state: MealPlanPageState,
  action: mealPlanPageAction,
): MealPlanPageState => {
  switch (action.type) {
    case 'MOVE_OPTION':
      return {
        ...state,
        mealPlan: {
          ...state.mealPlan,
          events: (state.mealPlan.events || []).map((event: MealPlanEvent, eventIndex: number) => {
            if (
              (action.optionIndex === 0 && action.direction === 'up') ||
              (action.optionIndex === event.options.length - 1 && action.direction === 'down')
            ) {
              return event;
            }

            const options = [...event.options];
            [
              options[action.direction === 'up' ? action.optionIndex - 1 : action.optionIndex + 1],
              options[action.optionIndex],
            ] = [
              options[action.optionIndex],
              options[action.direction === 'up' ? action.optionIndex - 1 : action.optionIndex + 1],
            ];

            return eventIndex !== action.eventIndex
              ? event
              : {
                  ...event,
                  options: options,
                };
          }),
        },
      };

    case 'ADD_VOTES_TO_MEAL_PLAN':
      return {
        ...state,
        mealPlan: {
          ...state.mealPlan,
          events: (state.mealPlan.events || []).map((event: MealPlanEvent, eventIndex: number) => {
            return eventIndex !== action.eventIndex
              ? event
              : new MealPlanEvent({
                  ...event,
                  options: event.options.map((option: MealPlanOption) => {
                    const votes = (action.votes || []).filter(
                      (vote: MealPlanOptionVote) => vote.belongsToMealPlanOption === option.id,
                    );
                    return new MealPlanOption({
                      ...option,
                      votes: votes,
                    });
                  }),
                });
          }),
        },
      };

    case 'UPDATE_MEAL_PLAN_TASK':
      return {
        ...state,
        tasks: state.tasks.map((task: MealPlanTask) => {
          return task.id === action.newTask.id ? action.newTask : task;
        }),
      };

    case 'UPDATE_MEAL_PLAN_GROCERY_LIST_ITEM':
      return {
        ...state,
        groceryList: state.groceryList.map((item: MealPlanGroceryListItem) => {
          return item.id === action.newItem.id ? action.newItem : item;
        }),
      };

    default:
      console.error(`Unhandled action type`);
      return state;
  }
};

/* END Meal Plan Creation Reducer */

const getMissingVotersForMealPlanEvent = (
  mealPlanEvent: MealPlanEvent,
  household: Household,
  userID: string,
): Array<string> => {
  const missingVotes: Set<string> = new Set<string>();

  mealPlanEvent.options.forEach((option: MealPlanOption) => {
    household.members.forEach((member: HouseholdUserMembershipWithUser) => {
      if (
        (option.votes || []).find((vote: MealPlanOptionVote) => vote.byUser === member.belongsToUser?.id) === undefined
      ) {
        missingVotes.add(member.belongsToUser?.id !== userID ? member.belongsToUser?.username || 'UNKNOWN' : 'you');
      }
    });
  });

  return Array.from(missingVotes.values());
};

const optionWasChosen = (option: MealPlanOption) => option.chosen;
const userVotedForMealPlanOption = (userID: string) => (vote: MealPlanOptionVote) => vote.byUser === userID;

const getUnvotedMealPlanEvents = (mealPlan: MealPlan, userID: string): Array<MealPlanEvent> => {
  return (mealPlan.events || []).filter((event: MealPlanEvent) => {
    return (
      event.options.find(
        (option: MealPlanOption) => (option.votes || []).find(userVotedForMealPlanOption(userID)) === undefined,
      ) !== undefined
    );
  });
};

const getChosenMealPlanEvents = (mealPlan: MealPlan): Array<MealPlanEvent> => {
  return mealPlan.events.filter((event: MealPlanEvent) => {
    return (event.options || []).find(optionWasChosen) !== undefined;
  });
};

const getMealPlanTasksForRecipe = (tasks: MealPlanTask[], recipeID: string): Array<MealPlanTask> => {
  return tasks.filter((task: MealPlanTask) => {
    return task.recipePrepTask.belongsToRecipe === recipeID;
  });
};

const getRecipesForMealPlanOptions = (options: MealPlanOption[]): Array<Recipe> => {
  return options
    .map((opt: MealPlanOption) => opt.meal.components.map((component: MealComponent) => component.recipe))
    .flat();
};

const findRecipeInMealPlan = (mealPlan: MealPlan, recipeID: string): Recipe | undefined => {
  let recipeToReturn = undefined;

  mealPlan.events.forEach((event: MealPlanEvent) => {
    event.options.forEach((option: MealPlanOption) => {
      option.meal.components.forEach((component: MealComponent) => {
        if (component.recipe.id === recipeID) {
          recipeToReturn = component.recipe;
        }
      });
    });
  });

  return recipeToReturn;
};

const getUserFromHouseholdByID = (
  household: Household,
  userID: string,
): HouseholdUserMembershipWithUser | undefined => {
  return household.members.find((member: HouseholdUserMembershipWithUser) => member.belongsToUser?.id === userID);
};

function MealPlanPage(props: MealPlanPageProps) {
  const apiClient = buildLocalClient();

  const pageLoadMealPlan = props.mealPlan;
  const userID = props.userID;
  const pageLoadHousehold = props.household;
  const pageLoadGroceryList = props.groceryList;
  const pageLoadTasks = props.tasks;

  const mealPlan = valueOrDefault(pageLoadMealPlan, new MealPlan());
  const [mealPlanError] = useState<IAPIError | undefined>(pageLoadMealPlan.error);

  const household = valueOrDefault(pageLoadHousehold, new MealPlan());
  const [householdError] = useState<IAPIError | undefined>(pageLoadHousehold.error);

  const groceryList = valueOrDefault(pageLoadGroceryList, new MealPlan());
  const [groceryListError] = useState<IAPIError | undefined>(pageLoadGroceryList.error);

  const tasks = valueOrDefault(pageLoadTasks, new MealPlan());
  const [tasksError] = useState<IAPIError | undefined>(pageLoadTasks.error);

  const [pageState, dispatchPageEvent] = useReducer(
    useMealPlanReducer,
    new MealPlanPageState(mealPlan, groceryList, tasks),
  );

  return (
    <AppLayout title="Meal Plan" containerSize="xl" userLoggedIn>
      <Center>
        <Stack>
          <Center>
            <Title order={3} p={5}>
              {`${format(new Date(getEarliestEvent(pageState.mealPlan).startsAt), dateFormat)} - ${format(
                new Date(getLatestEvent(pageState.mealPlan).startsAt),
                dateFormat,
              )}`}
            </Title>
          </Center>

          {mealPlan.status === 'finalized' && <Divider label="voted for" labelPosition="center" />}
          {mealPlan.status === 'awaiting_votes' && <Divider label="awaiting votes" labelPosition="center" />}
          {pageState.mealPlan.events.filter(
            (event: MealPlanEvent) => event.options.filter((option: MealPlanOption) => !option.chosen).length === 0,
          ).length > 0 && <Divider my="xl" label="decided" labelPosition="center" />}

          {mealPlan.status === 'awaiting_votes' && getUnvotedMealPlanEvents(pageState.mealPlan, userID).length > 0 && (
            <Button onClick={() => router.push(`/meal_plans/${mealPlan.id}/ballot`)}>Vote</Button>
          )}

          {mealPlan.events.map((event: MealPlanEvent, eventIndex: number) => {
            return (
              <Card shadow="xs" radius="md" withBorder my="xl" key={eventIndex}>
                <Grid justify="center" align="center">
                  <Title order={4}>{format(new Date(event.startsAt), 'iiii, M/d/yy @ h aa')}</Title>
                </Grid>
                {event.options
                  .sort((a: MealPlanOption, b: MealPlanOption) => (a.chosen ? -1 : b.chosen ? 1 : 0))
                  .map((option: MealPlanOption, optionIndex: number) => {
                    return (
                      <Grid key={optionIndex}>
                        <Grid.Col span="auto">
                          <Indicator
                            position="top-start"
                            offset={2}
                            label={
                              (option.votes || []).find(
                                (vote: MealPlanOptionVote) => vote.byUser === userID && vote.rank === 0,
                              ) !== undefined
                                ? '⭐'
                                : ''
                            }
                            color="none"
                          >
                            <Card shadow="xs" radius="md" withBorder mt="xs">
                              <Grid grow justify="space-between">
                                <Grid.Col span="auto">
                                  <Link key={option.meal.id} href={`/meals/${option.meal.id}`}>
                                    {option.meal.name}
                                  </Link>
                                </Grid.Col>
                                <Grid.Col span="auto">
                                  {!householdError && (
                                    <Box sx={{ float: 'right' }}>
                                      {(option.votes || []).map((vote: MealPlanOptionVote) => {
                                        const userWhoVoted = getUserFromHouseholdByID(
                                          household,
                                          vote.byUser,
                                        )?.belongsToUser;
                                        return (
                                          <Tooltip
                                            label={`${
                                              userWhoVoted?.firstName || userWhoVoted?.username || 'UNKNOWN'
                                            } ranked this choice #${vote.rank + 1}`}
                                            withArrow
                                            withinPortal
                                          >
                                            <Indicator
                                              mr="xs"
                                              offset={4}
                                              inline
                                              color={
                                                (vote.rank === 0 && 'yellow') ||
                                                (vote.rank === 1 && 'gray') ||
                                                (vote.rank === 2 && '#CD7F32') ||
                                                'blue'
                                              }
                                            >
                                              <Avatar
                                                radius="xl"
                                                src={userWhoVoted?.avatar || null}
                                                alt={`${
                                                  userWhoVoted?.firstName || userWhoVoted?.username || 'UNKNOWN'
                                                }'s avatar`}
                                              />
                                            </Indicator>
                                          </Tooltip>
                                        );
                                      })}
                                    </Box>
                                  )}
                                </Grid.Col>
                              </Grid>
                            </Card>
                          </Indicator>
                        </Grid.Col>
                      </Grid>
                    );
                  })}

                {householdError && <Text color="tomato">{householdError.message}</Text>}

                {!householdError && getMissingVotersForMealPlanEvent(event, household, userID).length > 0 && (
                  <Grid justify="center" align="center">
                    <Grid.Col span="auto">
                      <small>{`(awaiting votes from ${new Intl.ListFormat('en').format(
                        getMissingVotersForMealPlanEvent(event, household, userID),
                      )})`}</small>
                    </Grid.Col>
                  </Grid>
                )}
              </Card>
            );
          })}

          {mealPlanError && <Text color="tomato">{mealPlanError.message}</Text>}

          <Grid>
            <Grid.Col span={12} md={7}>
              {!mealPlanError && (getChosenMealPlanEvents(pageState.mealPlan) || []).length > 0 && (
                <>
                  <Divider label="tasks" labelPosition="center" />

                  {tasksError && <Text color="red">{tasksError.message}</Text>}

                  {!tasksError && !mealPlanError && (
                    <List icon={<></>}>
                      {(getChosenMealPlanEvents(pageState.mealPlan) || []).map(
                        (event: MealPlanEvent, eventIndex: number) => {
                          return (
                            <div key={eventIndex}>
                              <List.Item>
                                For{' '}
                                <Link
                                  href={`/meals/${
                                    (event.options || []).find((opt: MealPlanOption) => opt.chosen)!.meal.id
                                  }`}
                                >
                                  &nbsp;{event.mealName}&nbsp;
                                </Link>{' '}
                                at {format(new Date(event.startsAt), "h aa 'on' M/d/yy")}:&nbsp;
                              </List.Item>
                              <List icon={<></>} withPadding>
                                {getRecipesForMealPlanOptions(
                                  (event.options || []).filter((opt: MealPlanOption) => opt.chosen),
                                ).map((recipe: Recipe, recipeIndex: number) => {
                                  return (
                                    <div key={recipeIndex}>
                                      <List.Item>
                                        {'For'}&nbsp;
                                        <Link href={`/meal_plans/${mealPlan.id}/recipe/${recipe.id}`}>
                                          {recipe.name}
                                        </Link>
                                        :&nbsp;
                                      </List.Item>

                                      <List icon={<></>} withPadding>
                                        {getMealPlanTasksForRecipe(pageState.tasks, recipe.id).map(
                                          (mealPlanTask: MealPlanTask, taskIndex: number) => {
                                            return (
                                              <Box key={taskIndex}>
                                                <List.Item>
                                                  <Grid grow>
                                                    <Grid.Col span="content">
                                                      <Tooltip label="Mark as done">
                                                        <ActionIcon
                                                          disabled={!['unfinished'].includes(mealPlanTask.status)}
                                                          onClick={() => {
                                                            apiClient
                                                              .updateMealPlanTaskStatus(
                                                                mealPlan.id,
                                                                mealPlanTask.id,
                                                                new MealPlanTaskStatusChangeRequestInput({
                                                                  status: 'finished',
                                                                }),
                                                              )
                                                              .then((res: APIResponse<MealPlanTask>) => {
                                                                dispatchPageEvent({
                                                                  type: 'UPDATE_MEAL_PLAN_TASK',
                                                                  newTask: res.data,
                                                                });
                                                              });
                                                          }}
                                                        >
                                                          <IconCheck />
                                                        </ActionIcon>
                                                      </Tooltip>
                                                    </Grid.Col>
                                                    <Grid.Col span="content">
                                                      <Tooltip label="Cancel">
                                                        <ActionIcon
                                                          disabled={!['unfinished'].includes(mealPlanTask.status)}
                                                          onClick={() => {
                                                            apiClient
                                                              .updateMealPlanTaskStatus(
                                                                mealPlan.id,
                                                                mealPlanTask.id,
                                                                new MealPlanTaskStatusChangeRequestInput({
                                                                  status: 'canceled',
                                                                }),
                                                              )
                                                              .then((res: APIResponse<MealPlanTask>) => {
                                                                dispatchPageEvent({
                                                                  type: 'UPDATE_MEAL_PLAN_TASK',
                                                                  newTask: res.data,
                                                                });
                                                              });
                                                          }}
                                                        >
                                                          <IconCircleX />
                                                        </ActionIcon>
                                                      </Tooltip>
                                                    </Grid.Col>
                                                    <Grid.Col span="content">
                                                      <Text
                                                        strikethrough={['ignored', 'finished'].includes(
                                                          mealPlanTask.status,
                                                        )}
                                                      >
                                                        {`Between ${formatDuration(
                                                          intervalToDuration({
                                                            start: subSeconds(
                                                              new Date(event.startsAt),
                                                              mealPlanTask.recipePrepTask
                                                                .timeBufferBeforeRecipeInSeconds.max || 0,
                                                            ),
                                                            end: new Date(event.startsAt),
                                                          }),
                                                        )} before and `}
                                                        {mealPlanTask.recipePrepTask.timeBufferBeforeRecipeInSeconds
                                                          .min === 0
                                                          ? `time of ${event.mealName} prep, ${mealPlanTask.recipePrepTask.name}`
                                                          : format(
                                                              subSeconds(
                                                                new Date(event.startsAt),
                                                                mealPlanTask.recipePrepTask
                                                                  .timeBufferBeforeRecipeInSeconds.min,
                                                              ),
                                                              "h aa 'on' M/d/yy",
                                                            )}
                                                        <Text size="xs">
                                                          (store {mealPlanTask.recipePrepTask.storageType})
                                                          <Badge ml="xs" size="sm" color="orange">
                                                            Optional
                                                          </Badge>
                                                        </Text>
                                                      </Text>
                                                    </Grid.Col>
                                                  </Grid>
                                                </List.Item>
                                                {!['ignored', 'finished'].includes(mealPlanTask.status) && (
                                                  <List icon={<></>} withPadding mx="lg" mb="lg">
                                                    {mealPlanTask.recipePrepTask.recipeSteps.map(
                                                      (prepTaskStep: RecipePrepTaskStep, prepTaskStepIndex: number) => {
                                                        const relevantRecipe = findRecipeInMealPlan(
                                                          pageState.mealPlan,
                                                          mealPlanTask.recipePrepTask.belongsToRecipe,
                                                        )!;
                                                        const relevantRecipeStep = relevantRecipe.steps.find(
                                                          (step: RecipeStep) =>
                                                            step.id === prepTaskStep.belongsToRecipeStep,
                                                        )!;

                                                        return (
                                                          <List.Item key={prepTaskStepIndex} my="-sm">
                                                            <Text
                                                              strikethrough={['ignored', 'finished'].includes(
                                                                mealPlanTask.status,
                                                              )}
                                                            >
                                                              Step #
                                                              {relevantRecipe.steps.indexOf(relevantRecipeStep) + 1} (
                                                              {relevantRecipeStep.preparation.name}{' '}
                                                              {new Intl.ListFormat('en').format(
                                                                relevantRecipeStep.ingredients.map(
                                                                  (ingredient: RecipeStepIngredient) =>
                                                                    ingredient.ingredient?.pluralName ||
                                                                    ingredient.name,
                                                                ),
                                                              )}
                                                              )
                                                            </Text>
                                                          </List.Item>
                                                        );
                                                      },
                                                    )}
                                                  </List>
                                                )}
                                              </Box>
                                            );
                                          },
                                        )}
                                      </List>
                                    </div>
                                  );
                                })}
                              </List>
                            </div>
                          );
                        },
                      )}
                    </List>
                  )}
                </>
              )}
            </Grid.Col>
            <Grid.Col span={12} md={5}>
              {groceryListError && <Text color="tomato">{groceryListError.message}</Text>}

              {!groceryListError && (pageState.groceryList || []).length > 0 && (
                <>
                  <Divider label="grocery list" labelPosition="center" />

                  <Table mt="xl">
                    <thead>
                      <tr>
                        <th>Ingredient</th>
                        <th>Quantity</th>
                        <th colSpan={3}>Actions</th>
                      </tr>
                    </thead>
                    <tbody>
                      {(pageState.groceryList || []).map(
                        (groceryListItem: MealPlanGroceryListItem, groceryListItemIndex: number) => {
                          return (
                            <tr key={groceryListItemIndex}>
                              <td>
                                <Text strikethrough={['already owned', 'acquired'].includes(groceryListItem.status)}>
                                  {groceryListItem.quantityNeeded.min === 1
                                    ? groceryListItem.ingredient.name
                                    : groceryListItem.ingredient.pluralName}
                                </Text>
                              </td>
                              <td>
                                <Grid>
                                  <Grid.Col span={12} md={6}>
                                    {(!['already owned', 'acquired'].includes(groceryListItem.status) && (
                                      <NumberInput hideControls value={groceryListItem.quantityNeeded.min} />
                                    )) || (
                                      <Text strikethrough size="sm" mt="xs">
                                        {groceryListItem.quantityNeeded.min}
                                      </Text>
                                    )}
                                  </Grid.Col>
                                  <Grid.Col span={12} md={6} mt="xs">
                                    <Text
                                      strikethrough={['already owned', 'acquired'].includes(groceryListItem.status)}
                                    >
                                      {groceryListItem.quantityNeeded.min === 1
                                        ? groceryListItem.measurementUnit.name
                                        : groceryListItem.measurementUnit.pluralName}
                                    </Text>
                                  </Grid.Col>
                                </Grid>
                              </td>
                              <td>
                                {!['already owned', 'acquired'].includes(groceryListItem.status) && (
                                  <Tooltip label="Got it!">
                                    <ActionIcon
                                      disabled={['already owned', 'acquired'].includes(groceryListItem.status)}
                                      onClick={() => {
                                        apiClient
                                          .updateMealPlanGroceryListItem(
                                            mealPlan.id,
                                            groceryListItem.id,
                                            new MealPlanGroceryListItemUpdateRequestInput({ status: 'acquired' }),
                                          )
                                          .then((res: APIResponse<MealPlanGroceryListItem>) => {
                                            dispatchPageEvent({
                                              type: 'UPDATE_MEAL_PLAN_GROCERY_LIST_ITEM',
                                              newItem: res.data,
                                            });
                                          });
                                      }}
                                    >
                                      <IconCheck />
                                    </ActionIcon>
                                  </Tooltip>
                                )}
                              </td>
                              <td>
                                {!['already owned', 'acquired'].includes(groceryListItem.status) && (
                                  <Tooltip label="Had it">
                                    <ActionIcon
                                      onClick={() => {
                                        apiClient
                                          .updateMealPlanGroceryListItem(
                                            mealPlan.id,
                                            groceryListItem.id,
                                            new MealPlanGroceryListItemUpdateRequestInput({ status: 'already owned' }),
                                          )
                                          .then((res: APIResponse<MealPlanGroceryListItem>) => {
                                            dispatchPageEvent({
                                              type: 'UPDATE_MEAL_PLAN_GROCERY_LIST_ITEM',
                                              newItem: res.data,
                                            });
                                          });
                                      }}
                                    >
                                      <IconThumbUp />
                                    </ActionIcon>
                                  </Tooltip>
                                )}
                              </td>
                              <td>
                                {!['already owned', 'acquired'].includes(groceryListItem.status) && (
                                  <Tooltip label="Don't need it">
                                    <ActionIcon
                                      onClick={() => {
                                        if (confirm('Are you sure you want to delete this item?')) {
                                          apiClient.archiveMealPlanGroceryListItem(mealPlan.id, groceryListItem.id);
                                        }
                                      }}
                                    >
                                      <IconTrash />
                                    </ActionIcon>
                                  </Tooltip>
                                )}
                              </td>
                            </tr>
                          );
                        },
                      )}
                    </tbody>
                  </Table>
                </>
              )}
            </Grid.Col>
          </Grid>
        </Stack>
      </Center>
    </AppLayout>
  );
}

export default MealPlanPage;
