// GENERATED CODE, DO NOT EDIT MANUALLY

import axios, {
  AxiosInstance,
  AxiosError,
  AxiosRequestConfig,
  AxiosResponse,
  HeadersDefaults,
  InternalAxiosRequestConfig,
} from 'axios';
import { Span } from '@opentelemetry/api';

import { buildServerSideLogger, LoggerType } from '@dinnerdonebetter/logger';
import {
  APIResponse,
  AuditLogEntry,
  AvatarUpdateInput,
  CreateMealPlanTasksRequest,
  CreateMealPlanTasksResponse,
  EmailAddressVerificationRequestInput,
  FinalizeMealPlansRequest,
  FinalizeMealPlansResponse,
  Household,
  HouseholdCreationRequestInput,
  HouseholdInstrumentOwnership,
  HouseholdInstrumentOwnershipCreationRequestInput,
  HouseholdInstrumentOwnershipUpdateRequestInput,
  HouseholdInvitation,
  HouseholdInvitationCreationRequestInput,
  HouseholdInvitationUpdateRequestInput,
  HouseholdOwnershipTransferInput,
  HouseholdUpdateRequestInput,
  HouseholdUserMembership,
  InitializeMealPlanGroceryListRequest,
  InitializeMealPlanGroceryListResponse,
  JWTResponse,
  Meal,
  MealCreationRequestInput,
  MealPlan,
  MealPlanCreationRequestInput,
  MealPlanEvent,
  MealPlanEventCreationRequestInput,
  MealPlanEventUpdateRequestInput,
  MealPlanGroceryListItem,
  MealPlanGroceryListItemCreationRequestInput,
  MealPlanGroceryListItemUpdateRequestInput,
  MealPlanOption,
  MealPlanOptionCreationRequestInput,
  MealPlanOptionUpdateRequestInput,
  MealPlanOptionVote,
  MealPlanOptionVoteCreationRequestInput,
  MealPlanOptionVoteUpdateRequestInput,
  MealPlanTask,
  MealPlanTaskCreationRequestInput,
  MealPlanTaskStatusChangeRequestInput,
  MealPlanUpdateRequestInput,
  ModifyUserPermissionsInput,
  OAuth2Client,
  OAuth2ClientCreationRequestInput,
  OAuth2ClientCreationResponse,
  PasswordResetResponse,
  PasswordResetToken,
  PasswordResetTokenCreationRequestInput,
  PasswordResetTokenRedemptionRequestInput,
  PasswordUpdateInput,
  QueryFilter,
  QueryFilteredResult,
  Recipe,
  RecipeCreationRequestInput,
  RecipePrepTask,
  RecipePrepTaskCreationRequestInput,
  RecipePrepTaskStep,
  RecipePrepTaskUpdateRequestInput,
  RecipeRating,
  RecipeRatingCreationRequestInput,
  RecipeRatingUpdateRequestInput,
  RecipeStep,
  RecipeStepCompletionCondition,
  RecipeStepCompletionConditionForExistingRecipeCreationRequestInput,
  RecipeStepCompletionConditionUpdateRequestInput,
  RecipeStepCreationRequestInput,
  RecipeStepIngredient,
  RecipeStepIngredientCreationRequestInput,
  RecipeStepIngredientUpdateRequestInput,
  RecipeStepInstrument,
  RecipeStepInstrumentCreationRequestInput,
  RecipeStepInstrumentUpdateRequestInput,
  RecipeStepProduct,
  RecipeStepProductCreationRequestInput,
  RecipeStepProductUpdateRequestInput,
  RecipeStepUpdateRequestInput,
  RecipeStepVessel,
  RecipeStepVesselCreationRequestInput,
  RecipeStepVesselUpdateRequestInput,
  RecipeUpdateRequestInput,
  ServiceSetting,
  ServiceSettingConfiguration,
  ServiceSettingConfigurationCreationRequestInput,
  ServiceSettingConfigurationUpdateRequestInput,
  ServiceSettingCreationRequestInput,
  TOTPSecretRefreshInput,
  TOTPSecretRefreshResponse,
  TOTPSecretVerificationInput,
  User,
  UserAccountStatusUpdateInput,
  UserCreationResponse,
  UserDetailsUpdateRequestInput,
  UserEmailAddressUpdateInput,
  UserIngredientPreference,
  UserIngredientPreferenceCreationRequestInput,
  UserIngredientPreferenceUpdateRequestInput,
  UserLoginInput,
  UserNotification,
  UserNotificationCreationRequestInput,
  UserNotificationUpdateRequestInput,
  UserPermissionsRequestInput,
  UserPermissionsResponse,
  UserRegistrationInput,
  UserStatusResponse,
  UsernameReminderRequestInput,
  UsernameUpdateInput,
  ValidIngredient,
  ValidIngredientCreationRequestInput,
  ValidIngredientGroup,
  ValidIngredientGroupCreationRequestInput,
  ValidIngredientGroupUpdateRequestInput,
  ValidIngredientMeasurementUnit,
  ValidIngredientMeasurementUnitCreationRequestInput,
  ValidIngredientMeasurementUnitUpdateRequestInput,
  ValidIngredientPreparation,
  ValidIngredientPreparationCreationRequestInput,
  ValidIngredientPreparationUpdateRequestInput,
  ValidIngredientState,
  ValidIngredientStateCreationRequestInput,
  ValidIngredientStateIngredient,
  ValidIngredientStateIngredientCreationRequestInput,
  ValidIngredientStateIngredientUpdateRequestInput,
  ValidIngredientStateUpdateRequestInput,
  ValidIngredientUpdateRequestInput,
  ValidInstrument,
  ValidInstrumentCreationRequestInput,
  ValidInstrumentUpdateRequestInput,
  ValidMeasurementUnit,
  ValidMeasurementUnitConversion,
  ValidMeasurementUnitConversionCreationRequestInput,
  ValidMeasurementUnitConversionUpdateRequestInput,
  ValidMeasurementUnitCreationRequestInput,
  ValidMeasurementUnitUpdateRequestInput,
  ValidPreparation,
  ValidPreparationCreationRequestInput,
  ValidPreparationInstrument,
  ValidPreparationInstrumentCreationRequestInput,
  ValidPreparationInstrumentUpdateRequestInput,
  ValidPreparationUpdateRequestInput,
  ValidPreparationVessel,
  ValidPreparationVesselCreationRequestInput,
  ValidPreparationVesselUpdateRequestInput,
  ValidVessel,
  ValidVesselCreationRequestInput,
  ValidVesselUpdateRequestInput,
  Webhook,
  WebhookCreationRequestInput,
  WebhookTriggerEvent,
  WebhookTriggerEventCreationRequestInput,
} from '@dinnerdonebetter/models';

function _curlFromAxiosConfig(config: InternalAxiosRequestConfig): string {
  const method = (config?.method || 'UNKNOWN').toUpperCase();
  const url = config.url;
  const headers = config.headers || {};
  const data = config.data;

  ['get', 'delete', 'head', 'post', 'put', 'patch'].forEach((method) => {
    delete headers[method];
  });

  // iterate through headers["common"], and add each key's value to headers
  const headerDefault = headers as unknown as HeadersDefaults;
  for (const key in headerDefault['common']) {
    if (headerDefault['common'].hasOwnProperty(key)) {
      headers[key] = headerDefault['common'][key];
    }
  }
  delete headers['common'];

  let curlCommand = `curl -X ${method} "${config?.baseURL || 'MISSING_BASE_URL'}${url}"`;

  for (const key in headers) {
    if (headers.hasOwnProperty(key)) {
      curlCommand += ` -H "${key}: ${headers[key]}"`;
    }
  }

  if (data) {
    curlCommand += ` -d '${JSON.stringify(data)}'`;
  }

  return curlCommand;
}

export class DinnerDoneBetterAPIClient {
  baseURL: string;
  client: AxiosInstance;
  requestInterceptorID: number;
  responseInterceptorID: number;
  logger: LoggerType = buildServerSideLogger('api_client');

  constructor(clientName: string = 'DDB-Service-Client', baseURL: string = '', oauth2Token?: string) {
    this.baseURL = baseURL;

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'X-Request-Source': 'webapp',
      'X-Service-Client': clientName,
    };

    if (oauth2Token) {
      headers['Authorization'] = `Bearer ${oauth2Token}`;
    }

    this.client = axios.create({
      baseURL,
      timeout: 10000,
      withCredentials: true,
      crossDomain: true,
      headers,
    } as AxiosRequestConfig);

    this.requestInterceptorID = this.client.interceptors.request.use(
      (request: InternalAxiosRequestConfig) => {
        // this.logger.debug(`Request: ${request.method} ${request.baseURL}${request.url}`);
        console.log(`${_curlFromAxiosConfig(request)}`);

        return request;
      },
      (error) => {
        // Do whatever you want with the response error here
        // But, be SURE to return the rejected promise, so the caller still has
        // the option of additional specialized handling at the call-site:
        return Promise.reject(error);
      },
    );

    this.responseInterceptorID = this.client.interceptors.response.use(
      (response: AxiosResponse) => {
        this.logger.debug(
          `Response: ${response.status} ${response.config.method} ${response.config.url}`,
          // response.data,
        );

        // console.log(`${response.status} ${_curlFromAxiosConfig(response.config)}`);

        return response;
      },
      (error) => {
        return Promise.reject(error);
      },
    );
  }

  withSpan(span: Span): DinnerDoneBetterAPIClient {
    const spanContext = span.spanContext();
    const spanLogDetails = { spanID: spanContext.spanId, traceID: spanContext.traceId };

    this.client.interceptors.request.eject(this.requestInterceptorID);
    this.requestInterceptorID = this.client.interceptors.request.use(
      (request: InternalAxiosRequestConfig) => {
        this.logger.debug(`Request: ${request.method} ${request.url}`, spanLogDetails);

        // console.log(_curlFromAxiosConfig(request));

        if (spanContext.traceId) {
          request.headers.set('traceparent', spanContext.traceId);
        }

        return request;
      },
      (error) => {
        return Promise.reject(error);
      },
    );

    return this;
  }

  // eslint-disable-next-line no-unused-vars
  configureRouterRejectionInterceptor(redirectCallback: (_: Location) => void) {
    this.client.interceptors.response.eject(this.responseInterceptorID);
    this.responseInterceptorID = this.client.interceptors.response.use(
      (response: AxiosResponse) => {
        this.logger.debug(
          `Response: ${response.status} ${response.config.method} ${response.config.url}${response.config.method === 'POST' || response.config.method === 'PUT' ? ` ${JSON.stringify(response.config.data)}` : ''}`,
        );

        return response;
      },
      (error: AxiosError) => {
        console.debug(`Request failed: ${error.response?.status}`);
        if (error.response?.status === 401) {
          redirectCallback(window.location);
        }

        return Promise.reject(error);
      },
    );
  }

  async acceptHouseholdInvitation(
    householdInvitationID: string,
    input: HouseholdInvitationUpdateRequestInput,
  ): Promise<APIResponse<HouseholdInvitation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<HouseholdInvitation>>(
        `/api/v1/household_invitations/${householdInvitationID}/accept`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async adminLoginForJWT(input: UserLoginInput): Promise<APIResponse<JWTResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<JWTResponse>>(`/users/login/jwt/admin`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async adminUpdateUserStatus(input: UserAccountStatusUpdateInput): Promise<APIResponse<UserStatusResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<UserStatusResponse>>(`/api/v1/admin/users/status`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveHousehold(householdID: string): Promise<APIResponse<Household>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<Household>>(`/api/v1/households/${householdID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveHouseholdInstrumentOwnership(
    householdInstrumentOwnershipID: string,
  ): Promise<APIResponse<HouseholdInstrumentOwnership>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<HouseholdInstrumentOwnership>>(
        `/api/v1/households/instruments/${householdInstrumentOwnershipID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveMeal(mealID: string): Promise<APIResponse<Meal>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<Meal>>(`/api/v1/meals/${mealID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveMealPlan(mealPlanID: string): Promise<APIResponse<MealPlan>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<MealPlan>>(`/api/v1/meal_plans/${mealPlanID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveMealPlanEvent(mealPlanID: string, mealPlanEventID: string): Promise<APIResponse<MealPlanEvent>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<MealPlanEvent>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveMealPlanGroceryListItem(
    mealPlanID: string,
    mealPlanGroceryListItemID: string,
  ): Promise<APIResponse<MealPlanGroceryListItem>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<MealPlanGroceryListItem>>(
        `/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveMealPlanOption(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
  ): Promise<APIResponse<MealPlanOption>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<MealPlanOption>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveMealPlanOptionVote(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
    mealPlanOptionVoteID: string,
  ): Promise<APIResponse<MealPlanOptionVote>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<MealPlanOptionVote>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveOAuth2Client(oauth2ClientID: string): Promise<APIResponse<OAuth2Client>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<OAuth2Client>>(
        `/api/v1/oauth2_clients/${oauth2ClientID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveRecipe(recipeID: string): Promise<APIResponse<Recipe>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<Recipe>>(`/api/v1/recipes/${recipeID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveRecipePrepTask(recipeID: string, recipePrepTaskID: string): Promise<APIResponse<RecipePrepTask>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<RecipePrepTask>>(
        `/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveRecipeRating(recipeID: string, recipeRatingID: string): Promise<APIResponse<RecipeRating>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<RecipeRating>>(
        `/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveRecipeStep(recipeID: string, recipeStepID: string): Promise<APIResponse<RecipeStep>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<RecipeStep>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveRecipeStepCompletionCondition(
    recipeID: string,
    recipeStepID: string,
    recipeStepCompletionConditionID: string,
  ): Promise<APIResponse<RecipeStepCompletionCondition>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<RecipeStepCompletionCondition>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveRecipeStepIngredient(
    recipeID: string,
    recipeStepID: string,
    recipeStepIngredientID: string,
  ): Promise<APIResponse<RecipeStepIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<RecipeStepIngredient>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveRecipeStepInstrument(
    recipeID: string,
    recipeStepID: string,
    recipeStepInstrumentID: string,
  ): Promise<APIResponse<RecipeStepInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<RecipeStepInstrument>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveRecipeStepProduct(
    recipeID: string,
    recipeStepID: string,
    recipeStepProductID: string,
  ): Promise<APIResponse<RecipeStepProduct>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<RecipeStepProduct>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveRecipeStepVessel(
    recipeID: string,
    recipeStepID: string,
    recipeStepVesselID: string,
  ): Promise<APIResponse<RecipeStepVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<RecipeStepVessel>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveServiceSetting(serviceSettingID: string): Promise<APIResponse<ServiceSetting>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ServiceSetting>>(
        `/api/v1/settings/${serviceSettingID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveServiceSettingConfiguration(
    serviceSettingConfigurationID: string,
  ): Promise<APIResponse<ServiceSettingConfiguration>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ServiceSettingConfiguration>>(
        `/api/v1/settings/configurations/${serviceSettingConfigurationID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveUser(userID: string): Promise<APIResponse<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<User>>(`/api/v1/users/${userID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveUserIngredientPreference(
    userIngredientPreferenceID: string,
  ): Promise<APIResponse<UserIngredientPreference>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<UserIngredientPreference>>(
        `/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveUserMembership(householdID: string, userID: string): Promise<APIResponse<HouseholdUserMembership>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<HouseholdUserMembership>>(
        `/api/v1/households/${householdID}/members/${userID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidIngredient(validIngredientID: string): Promise<APIResponse<ValidIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidIngredient>>(
        `/api/v1/valid_ingredients/${validIngredientID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidIngredientGroup(validIngredientGroupID: string): Promise<APIResponse<ValidIngredientGroup>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidIngredientGroup>>(
        `/api/v1/valid_ingredient_groups/${validIngredientGroupID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidIngredientMeasurementUnit(
    validIngredientMeasurementUnitID: string,
  ): Promise<APIResponse<ValidIngredientMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidIngredientMeasurementUnit>>(
        `/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidIngredientPreparation(
    validIngredientPreparationID: string,
  ): Promise<APIResponse<ValidIngredientPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidIngredientPreparation>>(
        `/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidIngredientState(validIngredientStateID: string): Promise<APIResponse<ValidIngredientState>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidIngredientState>>(
        `/api/v1/valid_ingredient_states/${validIngredientStateID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidIngredientStateIngredient(
    validIngredientStateIngredientID: string,
  ): Promise<APIResponse<ValidIngredientStateIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidIngredientStateIngredient>>(
        `/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidInstrument(validInstrumentID: string): Promise<APIResponse<ValidInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidInstrument>>(
        `/api/v1/valid_instruments/${validInstrumentID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidMeasurementUnit(validMeasurementUnitID: string): Promise<APIResponse<ValidMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidMeasurementUnit>>(
        `/api/v1/valid_measurement_units/${validMeasurementUnitID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidMeasurementUnitConversion(
    validMeasurementUnitConversionID: string,
  ): Promise<APIResponse<ValidMeasurementUnitConversion>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidMeasurementUnitConversion>>(
        `/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidPreparation(validPreparationID: string): Promise<APIResponse<ValidPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidPreparation>>(
        `/api/v1/valid_preparations/${validPreparationID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidPreparationInstrument(
    validPreparationVesselID: string,
  ): Promise<APIResponse<ValidPreparationInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidPreparationInstrument>>(
        `/api/v1/valid_preparation_instruments/${validPreparationVesselID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidPreparationVessel(validPreparationVesselID: string): Promise<APIResponse<ValidPreparationVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidPreparationVessel>>(
        `/api/v1/valid_preparation_vessels/${validPreparationVesselID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveValidVessel(validVesselID: string): Promise<APIResponse<ValidVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<ValidVessel>>(`/api/v1/valid_vessels/${validVesselID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveWebhook(webhookID: string): Promise<APIResponse<Webhook>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<Webhook>>(`/api/v1/webhooks/${webhookID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async archiveWebhookTriggerEvent(
    webhookID: string,
    webhookTriggerEventID: string,
  ): Promise<APIResponse<WebhookTriggerEvent>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.delete<APIResponse<WebhookTriggerEvent>>(
        `/api/v1/webhooks/${webhookID}/trigger_events/${webhookTriggerEventID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async cancelHouseholdInvitation(
    householdInvitationID: string,
    input: HouseholdInvitationUpdateRequestInput,
  ): Promise<APIResponse<HouseholdInvitation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<HouseholdInvitation>>(
        `/api/v1/household_invitations/${householdInvitationID}/cancel`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async checkPermissions(input: UserPermissionsRequestInput): Promise<APIResponse<UserPermissionsResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<UserPermissionsResponse>>(
        `/api/v1/users/permissions/check`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async cloneRecipe(recipeID: string): Promise<APIResponse<Recipe>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<Recipe>>(`/api/v1/recipes/${recipeID}/clone`);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createHousehold(input: HouseholdCreationRequestInput): Promise<APIResponse<Household>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<Household>>(`/api/v1/households`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createHouseholdInstrumentOwnership(
    input: HouseholdInstrumentOwnershipCreationRequestInput,
  ): Promise<APIResponse<HouseholdInstrumentOwnership>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<HouseholdInstrumentOwnership>>(
        `/api/v1/households/instruments`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createHouseholdInvitation(
    householdID: string,
    input: HouseholdInvitationCreationRequestInput,
  ): Promise<APIResponse<HouseholdInvitation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<HouseholdInvitation>>(
        `/api/v1/households/${householdID}/invitations`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createMeal(input: MealCreationRequestInput): Promise<APIResponse<Meal>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<Meal>>(`/api/v1/meals`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createMealPlan(input: MealPlanCreationRequestInput): Promise<APIResponse<MealPlan>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<MealPlan>>(`/api/v1/meal_plans`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createMealPlanEvent(
    mealPlanID: string,
    input: MealPlanEventCreationRequestInput,
  ): Promise<APIResponse<MealPlanEvent>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<MealPlanEvent>>(
        `/api/v1/meal_plans/${mealPlanID}/events`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createMealPlanGroceryListItem(
    mealPlanID: string,
    input: MealPlanGroceryListItemCreationRequestInput,
  ): Promise<APIResponse<MealPlanGroceryListItem>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<MealPlanGroceryListItem>>(
        `/api/v1/meal_plans/${mealPlanID}/grocery_list_items`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createMealPlanOption(
    mealPlanID: string,
    mealPlanEventID: string,
    input: MealPlanOptionCreationRequestInput,
  ): Promise<APIResponse<MealPlanOption>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<MealPlanOption>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createMealPlanOptionVote(
    mealPlanID: string,
    mealPlanEventID: string,
    input: MealPlanOptionVoteCreationRequestInput,
  ): Promise<APIResponse<MealPlanOptionVote>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<MealPlanOptionVote>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/vote`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createMealPlanTask(
    mealPlanID: string,
    input: MealPlanTaskCreationRequestInput,
  ): Promise<APIResponse<MealPlanTask>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<MealPlanTask>>(
        `/api/v1/meal_plans/${mealPlanID}/tasks`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createOAuth2Client(
    input: OAuth2ClientCreationRequestInput,
  ): Promise<APIResponse<OAuth2ClientCreationResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<OAuth2ClientCreationResponse>>(
        `/api/v1/oauth2_clients`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createRecipe(input: RecipeCreationRequestInput): Promise<APIResponse<Recipe>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<Recipe>>(`/api/v1/recipes`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createRecipePrepTask(
    recipeID: string,
    input: RecipePrepTaskCreationRequestInput,
  ): Promise<APIResponse<RecipePrepTask>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<RecipePrepTask>>(
        `/api/v1/recipes/${recipeID}/prep_tasks`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createRecipeRating(
    recipeID: string,
    input: RecipeRatingCreationRequestInput,
  ): Promise<APIResponse<RecipeRating>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<RecipeRating>>(`/api/v1/recipes/${recipeID}/ratings`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createRecipeStep(recipeID: string, input: RecipeStepCreationRequestInput): Promise<APIResponse<RecipeStep>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<RecipeStep>>(`/api/v1/recipes/${recipeID}/steps`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createRecipeStepCompletionCondition(
    recipeID: string,
    recipeStepID: string,
    input: RecipeStepCompletionConditionForExistingRecipeCreationRequestInput,
  ): Promise<APIResponse<RecipeStepCompletionCondition>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<RecipeStepCompletionCondition>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createRecipeStepIngredient(
    recipeID: string,
    recipeStepID: string,
    input: RecipeStepIngredientCreationRequestInput,
  ): Promise<APIResponse<RecipeStepIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<RecipeStepIngredient>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createRecipeStepInstrument(
    recipeID: string,
    recipeStepID: string,
    input: RecipeStepInstrumentCreationRequestInput,
  ): Promise<APIResponse<RecipeStepInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<RecipeStepInstrument>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createRecipeStepProduct(
    recipeID: string,
    recipeStepID: string,
    input: RecipeStepProductCreationRequestInput,
  ): Promise<APIResponse<RecipeStepProduct>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<RecipeStepProduct>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createRecipeStepVessel(
    recipeID: string,
    recipeStepID: string,
    input: RecipeStepVesselCreationRequestInput,
  ): Promise<APIResponse<RecipeStepVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<RecipeStepVessel>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createServiceSetting(input: ServiceSettingCreationRequestInput): Promise<APIResponse<ServiceSetting>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ServiceSetting>>(`/api/v1/settings`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createServiceSettingConfiguration(
    input: ServiceSettingConfigurationCreationRequestInput,
  ): Promise<APIResponse<ServiceSettingConfiguration>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ServiceSettingConfiguration>>(
        `/api/v1/settings/configurations`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createUser(input: UserRegistrationInput): Promise<APIResponse<UserCreationResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<UserCreationResponse>>(`/users`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createUserIngredientPreference(
    input: UserIngredientPreferenceCreationRequestInput,
  ): Promise<APIResponse<UserIngredientPreference>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<UserIngredientPreference>>(
        `/api/v1/user_ingredient_preferences`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createUserNotification(input: UserNotificationCreationRequestInput): Promise<APIResponse<UserNotification>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<UserNotification>>(`/api/v1/user_notifications`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidIngredient(input: ValidIngredientCreationRequestInput): Promise<APIResponse<ValidIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidIngredient>>(`/api/v1/valid_ingredients`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidIngredientGroup(
    input: ValidIngredientGroupCreationRequestInput,
  ): Promise<APIResponse<ValidIngredientGroup>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidIngredientGroup>>(
        `/api/v1/valid_ingredient_groups`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidIngredientMeasurementUnit(
    input: ValidIngredientMeasurementUnitCreationRequestInput,
  ): Promise<APIResponse<ValidIngredientMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidIngredientMeasurementUnit>>(
        `/api/v1/valid_ingredient_measurement_units`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidIngredientPreparation(
    input: ValidIngredientPreparationCreationRequestInput,
  ): Promise<APIResponse<ValidIngredientPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidIngredientPreparation>>(
        `/api/v1/valid_ingredient_preparations`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidIngredientState(
    input: ValidIngredientStateCreationRequestInput,
  ): Promise<APIResponse<ValidIngredientState>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidIngredientState>>(
        `/api/v1/valid_ingredient_states`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidIngredientStateIngredient(
    input: ValidIngredientStateIngredientCreationRequestInput,
  ): Promise<APIResponse<ValidIngredientStateIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidIngredientStateIngredient>>(
        `/api/v1/valid_ingredient_state_ingredients`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidInstrument(input: ValidInstrumentCreationRequestInput): Promise<APIResponse<ValidInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidInstrument>>(`/api/v1/valid_instruments`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidMeasurementUnit(
    input: ValidMeasurementUnitCreationRequestInput,
  ): Promise<APIResponse<ValidMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidMeasurementUnit>>(
        `/api/v1/valid_measurement_units`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidMeasurementUnitConversion(
    input: ValidMeasurementUnitConversionCreationRequestInput,
  ): Promise<APIResponse<ValidMeasurementUnitConversion>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidMeasurementUnitConversion>>(
        `/api/v1/valid_measurement_conversions`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidPreparation(input: ValidPreparationCreationRequestInput): Promise<APIResponse<ValidPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidPreparation>>(`/api/v1/valid_preparations`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidPreparationInstrument(
    input: ValidPreparationInstrumentCreationRequestInput,
  ): Promise<APIResponse<ValidPreparationInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidPreparationInstrument>>(
        `/api/v1/valid_preparation_instruments`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidPreparationVessel(
    input: ValidPreparationVesselCreationRequestInput,
  ): Promise<APIResponse<ValidPreparationVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidPreparationVessel>>(
        `/api/v1/valid_preparation_vessels`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createValidVessel(input: ValidVesselCreationRequestInput): Promise<APIResponse<ValidVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<ValidVessel>>(`/api/v1/valid_vessels`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createWebhook(input: WebhookCreationRequestInput): Promise<APIResponse<Webhook>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<Webhook>>(`/api/v1/webhooks`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async createWebhookTriggerEvent(
    webhookID: string,
    input: WebhookTriggerEventCreationRequestInput,
  ): Promise<APIResponse<WebhookTriggerEvent>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<WebhookTriggerEvent>>(
        `/api/v1/webhooks/${webhookID}/trigger_events`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async finalizeMealPlan(mealPlanID: string): Promise<APIResponse<FinalizeMealPlansResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<FinalizeMealPlansResponse>>(
        `/api/v1/meal_plans/${mealPlanID}/finalize`,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getActiveHousehold(): Promise<APIResponse<Household>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Household>>(`/api/v1/households/current`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getAuditLogEntriesForHousehold(): Promise<APIResponse<AuditLogEntry>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<AuditLogEntry>>(`/api/v1/audit_log_entries/for_household`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getAuditLogEntriesForUser(): Promise<APIResponse<AuditLogEntry>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<AuditLogEntry>>(`/api/v1/audit_log_entries/for_user`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getAuditLogEntryByID(auditLogEntryID: string): Promise<APIResponse<AuditLogEntry>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<AuditLogEntry>>(
        `/api/v1/audit_log_entries/${auditLogEntryID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getAuthStatus(): Promise<APIResponse<UserStatusResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<UserStatusResponse>>(`/auth/status`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getHousehold(householdID: string): Promise<APIResponse<Household>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Household>>(`/api/v1/households/${householdID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getHouseholdInstrumentOwnership(
    householdInstrumentOwnershipID: string,
  ): Promise<APIResponse<HouseholdInstrumentOwnership>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<HouseholdInstrumentOwnership>>(
        `/api/v1/households/instruments/${householdInstrumentOwnershipID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getHouseholdInstrumentOwnerships(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<HouseholdInstrumentOwnership>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<HouseholdInstrumentOwnership>>>(
        `/api/v1/households/instruments`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<HouseholdInstrumentOwnership>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getHouseholdInvitation(householdInvitationID: string): Promise<APIResponse<HouseholdInvitation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<HouseholdInvitation>>(
        `/api/v1/household_invitations/${householdInvitationID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getHouseholdInvitationByID(
    householdID: string,
    householdInvitationID: string,
  ): Promise<APIResponse<HouseholdInvitation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<HouseholdInvitation>>(
        `/api/v1/households/${householdID}/invitations/${householdInvitationID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getHouseholds(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Household>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<Household>>>(`/api/v1/households`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<Household>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getMeal(mealID: string): Promise<APIResponse<Meal>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Meal>>(`/api/v1/meals/${mealID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getMealPlan(mealPlanID: string): Promise<APIResponse<MealPlan>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<MealPlan>>(`/api/v1/meal_plans/${mealPlanID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getMealPlanEvent(mealPlanID: string, mealPlanEventID: string): Promise<APIResponse<MealPlanEvent>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<MealPlanEvent>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getMealPlanEvents(
    filter: QueryFilter = QueryFilter.Default(),
    mealPlanID: string,
  ): Promise<QueryFilteredResult<MealPlanEvent>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<MealPlanEvent>>>(
        `/api/v1/meal_plans/${mealPlanID}/events`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<MealPlanEvent>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getMealPlanGroceryListItem(
    mealPlanID: string,
    mealPlanGroceryListItemID: string,
  ): Promise<APIResponse<MealPlanGroceryListItem>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<MealPlanGroceryListItem>>(
        `/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getMealPlanGroceryListItemsForMealPlan(
    filter: QueryFilter = QueryFilter.Default(),
    mealPlanID: string,
  ): Promise<QueryFilteredResult<MealPlanGroceryListItem>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<MealPlanGroceryListItem>>>(
        `/api/v1/meal_plans/${mealPlanID}/grocery_list_items`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<MealPlanGroceryListItem>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getMealPlanOption(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
  ): Promise<APIResponse<MealPlanOption>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<MealPlanOption>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getMealPlanOptionVote(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
    mealPlanOptionVoteID: string,
  ): Promise<APIResponse<MealPlanOptionVote>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<MealPlanOptionVote>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getMealPlanOptionVotes(
    filter: QueryFilter = QueryFilter.Default(),
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
  ): Promise<QueryFilteredResult<MealPlanOptionVote>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<MealPlanOptionVote>>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<MealPlanOptionVote>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getMealPlanOptions(
    filter: QueryFilter = QueryFilter.Default(),
    mealPlanID: string,
    mealPlanEventID: string,
  ): Promise<QueryFilteredResult<MealPlanOption>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<MealPlanOption>>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<MealPlanOption>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getMealPlanTask(mealPlanID: string, mealPlanTaskID: string): Promise<APIResponse<MealPlanTask>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<MealPlanTask>>(
        `/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getMealPlanTasks(
    filter: QueryFilter = QueryFilter.Default(),
    mealPlanID: string,
  ): Promise<QueryFilteredResult<MealPlanTask>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<MealPlanTask>>>(
        `/api/v1/meal_plans/${mealPlanID}/tasks`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<MealPlanTask>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getMealPlans(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<MealPlan>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<MealPlan>>>(`/api/v1/meal_plans`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<MealPlan>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getMeals(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Meal>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<Meal>>>(`/api/v1/meals`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<Meal>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getMermaidDiagramForRecipe(recipeID: string): Promise<APIResponse<string>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<string>>(`/api/v1/recipes/${recipeID}/mermaid`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getOAuth2Client(oauth2ClientID: string): Promise<APIResponse<OAuth2Client>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<OAuth2Client>>(`/api/v1/oauth2_clients/${oauth2ClientID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getOAuth2Clients(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<OAuth2Client>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<OAuth2Client>>>(`/api/v1/oauth2_clients`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<OAuth2Client>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getRandomValidIngredient(): Promise<APIResponse<ValidIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidIngredient>>(`/api/v1/valid_ingredients/random`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRandomValidInstrument(): Promise<APIResponse<ValidInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidInstrument>>(`/api/v1/valid_instruments/random`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRandomValidPreparation(): Promise<APIResponse<ValidPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidPreparation>>(`/api/v1/valid_preparations/random`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRandomValidVessel(): Promise<APIResponse<ValidVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidVessel>>(`/api/v1/valid_vessels/random`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getReceivedHouseholdInvitations(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<HouseholdInvitation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<HouseholdInvitation>>>(
        `/api/v1/household_invitations/received`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<HouseholdInvitation>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getRecipe(recipeID: string): Promise<APIResponse<Recipe>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Recipe>>(`/api/v1/recipes/${recipeID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRecipeMealPlanTasks(recipeID: string): Promise<APIResponse<RecipePrepTaskStep>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<RecipePrepTaskStep>>(
        `/api/v1/recipes/${recipeID}/prep_steps`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRecipePrepTask(recipeID: string, recipePrepTaskID: string): Promise<APIResponse<RecipePrepTask>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<RecipePrepTask>>(
        `/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRecipePrepTasks(
    filter: QueryFilter = QueryFilter.Default(),
    recipeID: string,
  ): Promise<QueryFilteredResult<RecipePrepTask>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<RecipePrepTask>>>(
        `/api/v1/recipes/${recipeID}/prep_tasks`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<RecipePrepTask>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getRecipeRating(recipeID: string, recipeRatingID: string): Promise<APIResponse<RecipeRating>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<RecipeRating>>(
        `/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRecipeRatings(
    filter: QueryFilter = QueryFilter.Default(),
    recipeID: string,
  ): Promise<QueryFilteredResult<RecipeRating>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<RecipeRating>>>(`/api/v1/recipes/${recipeID}/ratings`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<RecipeRating>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getRecipeStep(recipeID: string, recipeStepID: string): Promise<APIResponse<RecipeStep>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<RecipeStep>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRecipeStepCompletionCondition(
    recipeID: string,
    recipeStepID: string,
    recipeStepCompletionConditionID: string,
  ): Promise<APIResponse<RecipeStepCompletionCondition>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<RecipeStepCompletionCondition>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRecipeStepCompletionConditions(
    filter: QueryFilter = QueryFilter.Default(),
    recipeID: string,
    recipeStepID: string,
  ): Promise<QueryFilteredResult<RecipeStepCompletionCondition>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<RecipeStepCompletionCondition>>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<RecipeStepCompletionCondition>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getRecipeStepIngredient(
    recipeID: string,
    recipeStepID: string,
    recipeStepIngredientID: string,
  ): Promise<APIResponse<RecipeStepIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<RecipeStepIngredient>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRecipeStepIngredients(
    filter: QueryFilter = QueryFilter.Default(),
    recipeID: string,
    recipeStepID: string,
  ): Promise<QueryFilteredResult<RecipeStepIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<RecipeStepIngredient>>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<RecipeStepIngredient>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getRecipeStepInstrument(
    recipeID: string,
    recipeStepID: string,
    recipeStepInstrumentID: string,
  ): Promise<APIResponse<RecipeStepInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<RecipeStepInstrument>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRecipeStepInstruments(
    filter: QueryFilter = QueryFilter.Default(),
    recipeID: string,
    recipeStepID: string,
  ): Promise<QueryFilteredResult<RecipeStepInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<RecipeStepInstrument>>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<RecipeStepInstrument>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getRecipeStepProduct(
    recipeID: string,
    recipeStepID: string,
    recipeStepProductID: string,
  ): Promise<APIResponse<RecipeStepProduct>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<RecipeStepProduct>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRecipeStepProducts(
    filter: QueryFilter = QueryFilter.Default(),
    recipeID: string,
    recipeStepID: string,
  ): Promise<QueryFilteredResult<RecipeStepProduct>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<RecipeStepProduct>>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<RecipeStepProduct>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getRecipeStepVessel(
    recipeID: string,
    recipeStepID: string,
    recipeStepVesselID: string,
  ): Promise<APIResponse<RecipeStepVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<RecipeStepVessel>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getRecipeStepVessels(
    filter: QueryFilter = QueryFilter.Default(),
    recipeID: string,
    recipeStepID: string,
  ): Promise<QueryFilteredResult<RecipeStepVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<RecipeStepVessel>>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<RecipeStepVessel>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getRecipeSteps(
    filter: QueryFilter = QueryFilter.Default(),
    recipeID: string,
  ): Promise<QueryFilteredResult<RecipeStep>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<RecipeStep>>>(`/api/v1/recipes/${recipeID}/steps`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<RecipeStep>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getRecipes(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Recipe>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<Recipe>>>(`/api/v1/recipes`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<Recipe>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getSelf(): Promise<APIResponse<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<User>>(`/api/v1/users/self`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getSentHouseholdInvitations(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<HouseholdInvitation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<HouseholdInvitation>>>(
        `/api/v1/household_invitations/sent`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<HouseholdInvitation>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getServiceSetting(serviceSettingID: string): Promise<APIResponse<ServiceSetting>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ServiceSetting>>(`/api/v1/settings/${serviceSettingID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getServiceSettingConfigurationByName(
    filter: QueryFilter = QueryFilter.Default(),
    serviceSettingConfigurationName: string,
  ): Promise<QueryFilteredResult<ServiceSettingConfiguration>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ServiceSettingConfiguration>>>(
        `/api/v1/settings/configurations/user/${serviceSettingConfigurationName}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ServiceSettingConfiguration>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getServiceSettingConfigurationsForHousehold(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ServiceSettingConfiguration>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ServiceSettingConfiguration>>>(
        `/api/v1/settings/configurations/household`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ServiceSettingConfiguration>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getServiceSettingConfigurationsForUser(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ServiceSettingConfiguration>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ServiceSettingConfiguration>>>(
        `/api/v1/settings/configurations/user`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ServiceSettingConfiguration>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getServiceSettings(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<ServiceSetting>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ServiceSetting>>>(`/api/v1/settings`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ServiceSetting>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getUser(userID: string): Promise<APIResponse<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<User>>(`/api/v1/users/${userID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getUserIngredientPreferences(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<UserIngredientPreference>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<UserIngredientPreference>>>(
        `/api/v1/user_ingredient_preferences`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<UserIngredientPreference>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getUserNotification(userNotificationID: string): Promise<APIResponse<UserNotification>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<UserNotification>>(
        `/api/v1/user_notifications/${userNotificationID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getUserNotifications(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<UserNotification>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<UserNotification>>>(`/api/v1/user_notifications`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<UserNotification>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getUsers(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<User>>>(`/api/v1/users`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<User>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredient(validIngredientID: string): Promise<APIResponse<ValidIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidIngredient>>(
        `/api/v1/valid_ingredients/${validIngredientID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidIngredientGroup(validIngredientGroupID: string): Promise<APIResponse<ValidIngredientGroup>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidIngredientGroup>>(
        `/api/v1/valid_ingredient_groups/${validIngredientGroupID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidIngredientGroups(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientGroup>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientGroup>>>(
        `/api/v1/valid_ingredient_groups`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientGroup>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredientMeasurementUnit(
    validIngredientMeasurementUnitID: string,
  ): Promise<APIResponse<ValidIngredientMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidIngredientMeasurementUnit>>(
        `/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidIngredientMeasurementUnits(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientMeasurementUnit>>>(
        `/api/v1/valid_ingredient_measurement_units`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientMeasurementUnit>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredientMeasurementUnitsByIngredient(
    filter: QueryFilter = QueryFilter.Default(),
    validIngredientID: string,
  ): Promise<QueryFilteredResult<ValidIngredientMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientMeasurementUnit>>>(
        `/api/v1/valid_ingredient_measurement_units/by_ingredient/${validIngredientID}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientMeasurementUnit>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredientMeasurementUnitsByMeasurementUnit(
    filter: QueryFilter = QueryFilter.Default(),
    validMeasurementUnitID: string,
  ): Promise<QueryFilteredResult<ValidIngredientMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientMeasurementUnit>>>(
        `/api/v1/valid_ingredient_measurement_units/by_measurement_unit/${validMeasurementUnitID}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientMeasurementUnit>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredientPreparation(
    validIngredientPreparationID: string,
  ): Promise<APIResponse<ValidIngredientPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidIngredientPreparation>>(
        `/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidIngredientPreparations(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientPreparation>>>(
        `/api/v1/valid_ingredient_preparations`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientPreparation>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredientPreparationsByIngredient(
    filter: QueryFilter = QueryFilter.Default(),
    validIngredientID: string,
  ): Promise<QueryFilteredResult<ValidIngredientPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientPreparation>>>(
        `/api/v1/valid_ingredient_preparations/by_ingredient/${validIngredientID}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientPreparation>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredientPreparationsByPreparation(
    filter: QueryFilter = QueryFilter.Default(),
    validPreparationID: string,
  ): Promise<QueryFilteredResult<ValidIngredientPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientPreparation>>>(
        `/api/v1/valid_ingredient_preparations/by_preparation/${validPreparationID}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientPreparation>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredientState(validIngredientStateID: string): Promise<APIResponse<ValidIngredientState>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidIngredientState>>(
        `/api/v1/valid_ingredient_states/${validIngredientStateID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidIngredientStateIngredient(
    validIngredientStateIngredientID: string,
  ): Promise<APIResponse<ValidIngredientStateIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidIngredientStateIngredient>>(
        `/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidIngredientStateIngredients(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientStateIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientStateIngredient>>>(
        `/api/v1/valid_ingredient_state_ingredients`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientStateIngredient>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredientStateIngredientsByIngredient(
    filter: QueryFilter = QueryFilter.Default(),
    validIngredientID: string,
  ): Promise<QueryFilteredResult<ValidIngredientStateIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientStateIngredient>>>(
        `/api/v1/valid_ingredient_state_ingredients/by_ingredient/${validIngredientID}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientStateIngredient>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredientStateIngredientsByIngredientState(
    filter: QueryFilter = QueryFilter.Default(),
    validIngredientStateID: string,
  ): Promise<QueryFilteredResult<ValidIngredientStateIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientStateIngredient>>>(
        `/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/${validIngredientStateID}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientStateIngredient>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredientStates(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientState>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientState>>>(
        `/api/v1/valid_ingredient_states`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientState>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredients(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredient>>>(`/api/v1/valid_ingredients`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredient>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidIngredientsByPreparation(
    filter: QueryFilter = QueryFilter.Default(),
    q: string,
    validPreparationID: string,
  ): Promise<QueryFilteredResult<ValidIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredient>>>(
        `/api/v1/valid_ingredients/by_preparation/${validPreparationID}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredient>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidInstrument(validInstrumentID: string): Promise<APIResponse<ValidInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidInstrument>>(
        `/api/v1/valid_instruments/${validInstrumentID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidInstruments(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidInstrument>>>(`/api/v1/valid_instruments`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidInstrument>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidMeasurementUnit(validMeasurementUnitID: string): Promise<APIResponse<ValidMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidMeasurementUnit>>(
        `/api/v1/valid_measurement_units/${validMeasurementUnitID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidMeasurementUnitConversion(
    validMeasurementUnitConversionID: string,
  ): Promise<APIResponse<ValidMeasurementUnitConversion>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidMeasurementUnitConversion>>(
        `/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidMeasurementUnitConversionsFromUnit(
    validMeasurementUnitID: string,
  ): Promise<APIResponse<ValidMeasurementUnitConversion>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidMeasurementUnitConversion>>(
        `/api/v1/valid_measurement_conversions/from_unit/${validMeasurementUnitID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidMeasurementUnits(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidMeasurementUnit>>>(
        `/api/v1/valid_measurement_units`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidMeasurementUnit>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidPreparation(validPreparationID: string): Promise<APIResponse<ValidPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidPreparation>>(
        `/api/v1/valid_preparations/${validPreparationID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidPreparationInstrument(
    validPreparationVesselID: string,
  ): Promise<APIResponse<ValidPreparationInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidPreparationInstrument>>(
        `/api/v1/valid_preparation_instruments/${validPreparationVesselID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidPreparationInstruments(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidPreparationInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidPreparationInstrument>>>(
        `/api/v1/valid_preparation_instruments`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidPreparationInstrument>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidPreparationInstrumentsByInstrument(
    filter: QueryFilter = QueryFilter.Default(),
    validInstrumentID: string,
  ): Promise<QueryFilteredResult<ValidPreparationInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidPreparationInstrument>>>(
        `/api/v1/valid_preparation_instruments/by_instrument/${validInstrumentID}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidPreparationInstrument>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidPreparationInstrumentsByPreparation(
    filter: QueryFilter = QueryFilter.Default(),
    validPreparationID: string,
  ): Promise<QueryFilteredResult<ValidPreparationInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidPreparationInstrument>>>(
        `/api/v1/valid_preparation_instruments/by_preparation/${validPreparationID}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidPreparationInstrument>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidPreparationVessel(validPreparationVesselID: string): Promise<APIResponse<ValidPreparationVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidPreparationVessel>>(
        `/api/v1/valid_preparation_vessels/${validPreparationVesselID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidPreparationVessels(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidPreparationVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidPreparationVessel>>>(
        `/api/v1/valid_preparation_vessels`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidPreparationVessel>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidPreparationVesselsByPreparation(
    filter: QueryFilter = QueryFilter.Default(),
    validPreparationID: string,
  ): Promise<QueryFilteredResult<ValidPreparationVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidPreparationVessel>>>(
        `/api/v1/valid_preparation_vessels/by_preparation/${validPreparationID}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidPreparationVessel>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidPreparationVesselsByVessel(
    filter: QueryFilter = QueryFilter.Default(),
    ValidVesselID: string,
  ): Promise<QueryFilteredResult<ValidPreparationVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidPreparationVessel>>>(
        `/api/v1/valid_preparation_vessels/by_vessel/${ValidVesselID}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidPreparationVessel>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidPreparations(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidPreparation>>>(`/api/v1/valid_preparations`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidPreparation>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getValidVessel(validVesselID: string): Promise<APIResponse<ValidVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidVessel>>(`/api/v1/valid_vessels/${validVesselID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getValidVessels(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<ValidVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidVessel>>>(`/api/v1/valid_vessels`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidVessel>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async getWebhook(webhookID: string): Promise<APIResponse<Webhook>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Webhook>>(`/api/v1/webhooks/${webhookID}`, {});

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async getWebhooks(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Webhook>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<Webhook>>>(`/api/v1/webhooks`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<Webhook>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async loginForJWT(input: UserLoginInput): Promise<APIResponse<JWTResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<JWTResponse>>(`/users/login/jwt`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async pOST_households_householdID_invite(
    householdID: string,
    input: HouseholdInvitationCreationRequestInput,
  ): Promise<APIResponse<HouseholdInvitation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<HouseholdInvitation>>(
        `/api/v1/households/${householdID}/invite`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async redeemPasswordResetToken(input: PasswordResetTokenRedemptionRequestInput): Promise<APIResponse<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<User>>(`/users/password/reset/redeem`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async refreshTOTPSecret(input: TOTPSecretRefreshInput): Promise<APIResponse<TOTPSecretRefreshResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<TOTPSecretRefreshResponse>>(
        `/api/v1/users/totp_secret/new`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async rejectHouseholdInvitation(
    householdInvitationID: string,
    input: HouseholdInvitationUpdateRequestInput,
  ): Promise<APIResponse<HouseholdInvitation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<HouseholdInvitation>>(
        `/api/v1/household_invitations/${householdInvitationID}/reject`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async requestPasswordResetToken(
    input: PasswordResetTokenCreationRequestInput,
  ): Promise<APIResponse<PasswordResetToken>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<PasswordResetToken>>(`/users/password/reset`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async requestUsernameReminder(input: UsernameReminderRequestInput): Promise<APIResponse<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<User>>(`/users/username/reminder`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async runFinalizeMealPlanWorker(input: FinalizeMealPlansRequest): Promise<APIResponse<FinalizeMealPlansResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<FinalizeMealPlansResponse>>(
        `/api/v1/workers/finalize_meal_plans`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async runMealPlanGroceryListInitializerWorker(
    input: InitializeMealPlanGroceryListRequest,
  ): Promise<APIResponse<InitializeMealPlanGroceryListResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<InitializeMealPlanGroceryListResponse>>(
        `/api/v1/workers/meal_plan_grocery_list_init`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async runMealPlanTaskCreatorWorker(
    input: CreateMealPlanTasksRequest,
  ): Promise<APIResponse<CreateMealPlanTasksResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<CreateMealPlanTasksResponse>>(
        `/api/v1/workers/meal_plan_tasks`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async searchForMeals(filter: QueryFilter = QueryFilter.Default(), q: string): Promise<QueryFilteredResult<Meal>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<Meal>>>(`/api/v1/meals/search`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<Meal>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async searchForRecipes(filter: QueryFilter = QueryFilter.Default(), q: string): Promise<QueryFilteredResult<Recipe>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<Recipe>>>(`/api/v1/recipes/search`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<Recipe>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async searchForServiceSettings(
    filter: QueryFilter = QueryFilter.Default(),
    q: string,
  ): Promise<QueryFilteredResult<ServiceSetting>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ServiceSetting>>>(`/api/v1/settings/search`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ServiceSetting>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async searchForUsers(filter: QueryFilter = QueryFilter.Default(), q: string): Promise<QueryFilteredResult<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<User>>>(`/api/v1/users/search`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<User>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async searchForValidIngredientGroups(
    filter: QueryFilter = QueryFilter.Default(),
    q: string,
  ): Promise<QueryFilteredResult<ValidIngredientGroup>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientGroup>>>(
        `/api/v1/valid_ingredient_groups/search`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientGroup>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async searchForValidIngredientStates(
    filter: QueryFilter = QueryFilter.Default(),
    q: string,
  ): Promise<QueryFilteredResult<ValidIngredientState>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredientState>>>(
        `/api/v1/valid_ingredient_states/search`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredientState>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async searchForValidIngredients(
    filter: QueryFilter = QueryFilter.Default(),
    q: string,
  ): Promise<QueryFilteredResult<ValidIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidIngredient>>>(`/api/v1/valid_ingredients/search`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidIngredient>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async searchForValidInstruments(
    filter: QueryFilter = QueryFilter.Default(),
    q: string,
  ): Promise<QueryFilteredResult<ValidInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidInstrument>>>(`/api/v1/valid_instruments/search`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidInstrument>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async searchForValidMeasurementUnits(
    filter: QueryFilter = QueryFilter.Default(),
    q: string,
  ): Promise<QueryFilteredResult<ValidMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidMeasurementUnit>>>(
        `/api/v1/valid_measurement_units/search`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidMeasurementUnit>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async searchForValidPreparations(
    filter: QueryFilter = QueryFilter.Default(),
    q: string,
  ): Promise<QueryFilteredResult<ValidPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidPreparation>>>(
        `/api/v1/valid_preparations/search`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidPreparation>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async searchForValidVessels(
    filter: QueryFilter = QueryFilter.Default(),
    q: string,
  ): Promise<QueryFilteredResult<ValidVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidVessel>>>(`/api/v1/valid_vessels/search`, {
        params: filter.asRecord(),
      });

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidVessel>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async searchValidMeasurementUnitsByIngredient(
    filter: QueryFilter = QueryFilter.Default(),
    q: string,
    validIngredientID: string,
  ): Promise<QueryFilteredResult<ValidMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<Array<ValidMeasurementUnit>>>(
        `/api/v1/valid_measurement_units/by_ingredient/${validIngredientID}`,
        {
          params: filter.asRecord(),
        },
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      const result = new QueryFilteredResult<ValidMeasurementUnit>({
        data: response.data.data,
        totalCount: response.data.pagination?.totalCount,
        page: response.data.pagination?.page,
        limit: response.data.pagination?.limit,
      });

      resolve(result);
    });
  }

  async setDefaultHousehold(householdID: string): Promise<APIResponse<Household>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<Household>>(`/api/v1/households/${householdID}/default`);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async transferHouseholdOwnership(
    householdID: string,
    input: HouseholdOwnershipTransferInput,
  ): Promise<APIResponse<Household>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<Household>>(
        `/api/v1/households/${householdID}/transfer`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateHousehold(householdID: string, input: HouseholdUpdateRequestInput): Promise<APIResponse<Household>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<Household>>(`/api/v1/households/${householdID}`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateHouseholdInstrumentOwnership(
    householdInstrumentOwnershipID: string,
    input: HouseholdInstrumentOwnershipUpdateRequestInput,
  ): Promise<APIResponse<HouseholdInstrumentOwnership>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<HouseholdInstrumentOwnership>>(
        `/api/v1/households/instruments/${householdInstrumentOwnershipID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateHouseholdMemberPermissions(
    householdID: string,
    userID: string,
    input: ModifyUserPermissionsInput,
  ): Promise<APIResponse<UserPermissionsResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.patch<APIResponse<UserPermissionsResponse>>(
        `/api/v1/households/${householdID}/members/${userID}/permissions`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateMealPlan(mealPlanID: string, input: MealPlanUpdateRequestInput): Promise<APIResponse<MealPlan>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<MealPlan>>(`/api/v1/meal_plans/${mealPlanID}`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateMealPlanEvent(
    mealPlanID: string,
    mealPlanEventID: string,
    input: MealPlanEventUpdateRequestInput,
  ): Promise<APIResponse<MealPlanEvent>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<MealPlanEvent>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateMealPlanGroceryListItem(
    mealPlanID: string,
    mealPlanGroceryListItemID: string,
    input: MealPlanGroceryListItemUpdateRequestInput,
  ): Promise<APIResponse<MealPlanGroceryListItem>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<MealPlanGroceryListItem>>(
        `/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateMealPlanOption(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
    input: MealPlanOptionUpdateRequestInput,
  ): Promise<APIResponse<MealPlanOption>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<MealPlanOption>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateMealPlanOptionVote(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
    mealPlanOptionVoteID: string,
    input: MealPlanOptionVoteUpdateRequestInput,
  ): Promise<APIResponse<MealPlanOptionVote>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<MealPlanOptionVote>>(
        `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateMealPlanTaskStatus(
    mealPlanID: string,
    mealPlanTaskID: string,
    input: MealPlanTaskStatusChangeRequestInput,
  ): Promise<APIResponse<MealPlanTask>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.patch<APIResponse<MealPlanTask>>(
        `/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updatePassword(input: PasswordUpdateInput): Promise<APIResponse<PasswordResetResponse>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<PasswordResetResponse>>(`/api/v1/users/password/new`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateRecipe(recipeID: string, input: RecipeUpdateRequestInput): Promise<APIResponse<Recipe>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<Recipe>>(`/api/v1/recipes/${recipeID}`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateRecipePrepTask(
    recipeID: string,
    recipePrepTaskID: string,
    input: RecipePrepTaskUpdateRequestInput,
  ): Promise<APIResponse<RecipePrepTask>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<RecipePrepTask>>(
        `/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateRecipeRating(
    recipeID: string,
    recipeRatingID: string,
    input: RecipeRatingUpdateRequestInput,
  ): Promise<APIResponse<RecipeRating>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<RecipeRating>>(
        `/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateRecipeStep(
    recipeID: string,
    recipeStepID: string,
    input: RecipeStepUpdateRequestInput,
  ): Promise<APIResponse<RecipeStep>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<RecipeStep>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateRecipeStepCompletionCondition(
    recipeID: string,
    recipeStepID: string,
    recipeStepCompletionConditionID: string,
    input: RecipeStepCompletionConditionUpdateRequestInput,
  ): Promise<APIResponse<RecipeStepCompletionCondition>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<RecipeStepCompletionCondition>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateRecipeStepIngredient(
    recipeID: string,
    recipeStepID: string,
    recipeStepIngredientID: string,
    input: RecipeStepIngredientUpdateRequestInput,
  ): Promise<APIResponse<RecipeStepIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<RecipeStepIngredient>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateRecipeStepInstrument(
    recipeID: string,
    recipeStepID: string,
    recipeStepInstrumentID: string,
    input: RecipeStepInstrumentUpdateRequestInput,
  ): Promise<APIResponse<RecipeStepInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<RecipeStepInstrument>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateRecipeStepProduct(
    recipeID: string,
    recipeStepID: string,
    recipeStepProductID: string,
    input: RecipeStepProductUpdateRequestInput,
  ): Promise<APIResponse<RecipeStepProduct>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<RecipeStepProduct>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateRecipeStepVessel(
    recipeID: string,
    recipeStepID: string,
    recipeStepVesselID: string,
    input: RecipeStepVesselUpdateRequestInput,
  ): Promise<APIResponse<RecipeStepVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<RecipeStepVessel>>(
        `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateServiceSettingConfiguration(
    serviceSettingConfigurationID: string,
    input: ServiceSettingConfigurationUpdateRequestInput,
  ): Promise<APIResponse<ServiceSettingConfiguration>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ServiceSettingConfiguration>>(
        `/api/v1/settings/configurations/${serviceSettingConfigurationID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateUserDetails(input: UserDetailsUpdateRequestInput): Promise<APIResponse<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<User>>(`/api/v1/users/details`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateUserEmailAddress(input: UserEmailAddressUpdateInput): Promise<APIResponse<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<User>>(`/api/v1/users/email_address`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateUserIngredientPreference(
    userIngredientPreferenceID: string,
    input: UserIngredientPreferenceUpdateRequestInput,
  ): Promise<APIResponse<UserIngredientPreference>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<UserIngredientPreference>>(
        `/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateUserNotification(
    userNotificationID: string,
    input: UserNotificationUpdateRequestInput,
  ): Promise<APIResponse<UserNotification>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.patch<APIResponse<UserNotification>>(
        `/api/v1/user_notifications/${userNotificationID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateUserUsername(input: UsernameUpdateInput): Promise<APIResponse<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<User>>(`/api/v1/users/username`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidIngredient(
    validIngredientID: string,
    input: ValidIngredientUpdateRequestInput,
  ): Promise<APIResponse<ValidIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidIngredient>>(
        `/api/v1/valid_ingredients/${validIngredientID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidIngredientGroup(
    validIngredientGroupID: string,
    input: ValidIngredientGroupUpdateRequestInput,
  ): Promise<APIResponse<ValidIngredientGroup>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidIngredientGroup>>(
        `/api/v1/valid_ingredient_groups/${validIngredientGroupID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidIngredientMeasurementUnit(
    validIngredientMeasurementUnitID: string,
    input: ValidIngredientMeasurementUnitUpdateRequestInput,
  ): Promise<APIResponse<ValidIngredientMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidIngredientMeasurementUnit>>(
        `/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidIngredientPreparation(
    validIngredientPreparationID: string,
    input: ValidIngredientPreparationUpdateRequestInput,
  ): Promise<APIResponse<ValidIngredientPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidIngredientPreparation>>(
        `/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidIngredientState(
    validIngredientStateID: string,
    input: ValidIngredientStateUpdateRequestInput,
  ): Promise<APIResponse<ValidIngredientState>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidIngredientState>>(
        `/api/v1/valid_ingredient_states/${validIngredientStateID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidIngredientStateIngredient(
    validIngredientStateIngredientID: string,
    input: ValidIngredientStateIngredientUpdateRequestInput,
  ): Promise<APIResponse<ValidIngredientStateIngredient>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidIngredientStateIngredient>>(
        `/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidInstrument(
    validInstrumentID: string,
    input: ValidInstrumentUpdateRequestInput,
  ): Promise<APIResponse<ValidInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidInstrument>>(
        `/api/v1/valid_instruments/${validInstrumentID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidMeasurementUnit(
    validMeasurementUnitID: string,
    input: ValidMeasurementUnitUpdateRequestInput,
  ): Promise<APIResponse<ValidMeasurementUnit>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidMeasurementUnit>>(
        `/api/v1/valid_measurement_units/${validMeasurementUnitID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidMeasurementUnitConversion(
    validMeasurementUnitConversionID: string,
    input: ValidMeasurementUnitConversionUpdateRequestInput,
  ): Promise<APIResponse<ValidMeasurementUnitConversion>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidMeasurementUnitConversion>>(
        `/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidPreparation(
    validPreparationID: string,
    input: ValidPreparationUpdateRequestInput,
  ): Promise<APIResponse<ValidPreparation>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidPreparation>>(
        `/api/v1/valid_preparations/${validPreparationID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidPreparationInstrument(
    validPreparationVesselID: string,
    input: ValidPreparationInstrumentUpdateRequestInput,
  ): Promise<APIResponse<ValidPreparationInstrument>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidPreparationInstrument>>(
        `/api/v1/valid_preparation_instruments/${validPreparationVesselID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidPreparationVessel(
    validPreparationVesselID: string,
    input: ValidPreparationVesselUpdateRequestInput,
  ): Promise<APIResponse<ValidPreparationVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidPreparationVessel>>(
        `/api/v1/valid_preparation_vessels/${validPreparationVesselID}`,
        input,
      );
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async updateValidVessel(
    validVesselID: string,
    input: ValidVesselUpdateRequestInput,
  ): Promise<APIResponse<ValidVessel>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.put<APIResponse<ValidVessel>>(`/api/v1/valid_vessels/${validVesselID}`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async uploadUserAvatar(input: AvatarUpdateInput): Promise<APIResponse<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<User>>(`/api/v1/users/avatar/upload`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async validMeasurementUnitConversionsToUnit(
    validMeasurementUnitID: string,
  ): Promise<APIResponse<ValidMeasurementUnitConversion>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.get<APIResponse<ValidMeasurementUnitConversion>>(
        `/api/v1/valid_measurement_conversions/to_unit/${validMeasurementUnitID}`,
        {},
      );

      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async verifyEmailAddress(input: EmailAddressVerificationRequestInput): Promise<APIResponse<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<User>>(`/users/email_address/verify`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async verifyTOTPSecret(input: TOTPSecretVerificationInput): Promise<APIResponse<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<User>>(`/users/totp_secret/verify`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }

  async verifyUserEmailAddress(input: EmailAddressVerificationRequestInput): Promise<APIResponse<User>> {
    return new Promise(async function (resolve, reject) {
      const response = await this.client.post<APIResponse<User>>(`/api/v1/users/email_address_verification`, input);
      if (response.data.error) {
        reject(new Error(response.data.error.message));
      }

      resolve(response.data);
    });
  }
}
