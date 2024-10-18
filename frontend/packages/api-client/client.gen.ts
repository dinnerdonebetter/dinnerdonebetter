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
  oauth2Token: string;
  requestInterceptorID: number;
  responseInterceptorID: number;
  logger: LoggerType = buildServerSideLogger('api_client');

  constructor(baseURL: string = '', oauth2Token?: string, clientName: string = 'DDB-Service-Client') {
    this.baseURL = baseURL;
    this.oauth2Token = '';

    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'X-Request-Source': 'webapp',
      'X-Service-Client': clientName,
    };

    // because this client is used both in the browser and on the server, we can't mandate oauth2 tokens
    if (oauth2Token) {
      this.oauth2Token = oauth2Token;
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
        // console.log(`${_curlFromAxiosConfig(request)}`);

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
        // console.log(_curlFromAxiosConfig(request));

        if (spanContext.traceId) {
          request.headers.set('traceparent', spanLogDetails.traceID);
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
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<HouseholdInvitation>>(`/api/v1/household_invitations/${householdInvitationID}/accept`, input)
        .then((res: AxiosResponse<APIResponse<HouseholdInvitation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<HouseholdInvitation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async adminLoginForJWT(input: UserLoginInput): Promise<AxiosResponse<APIResponse<JWTResponse>>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<JWTResponse>>(`/users/login/jwt/admin`, input)
        .then((res: AxiosResponse<APIResponse<JWTResponse>>) => {
          if (res.data.error && res.data.error.message.toLowerCase() != 'totp required') {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res);
          }
        })
        .catch((error: AxiosError<APIResponse<JWTResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async adminUpdateUserStatus(input: UserAccountStatusUpdateInput): Promise<APIResponse<UserStatusResponse>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<UserStatusResponse>>(`/api/v1/admin/users/status`, input)
        .then((res: AxiosResponse<APIResponse<UserStatusResponse>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<UserStatusResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveHousehold(householdID: string): Promise<APIResponse<Household>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<Household>>(`/api/v1/households/${householdID}`)
        .then((res: AxiosResponse<APIResponse<Household>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Household>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveHouseholdInstrumentOwnership(
    householdInstrumentOwnershipID: string,
  ): Promise<APIResponse<HouseholdInstrumentOwnership>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<HouseholdInstrumentOwnership>>(
          `/api/v1/households/instruments/${householdInstrumentOwnershipID}`,
        )
        .then((res: AxiosResponse<APIResponse<HouseholdInstrumentOwnership>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<HouseholdInstrumentOwnership>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveMeal(mealID: string): Promise<APIResponse<Meal>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<Meal>>(`/api/v1/meals/${mealID}`)
        .then((res: AxiosResponse<APIResponse<Meal>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Meal>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveMealPlan(mealPlanID: string): Promise<APIResponse<MealPlan>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<MealPlan>>(`/api/v1/meal_plans/${mealPlanID}`)
        .then((res: AxiosResponse<APIResponse<MealPlan>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlan>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveMealPlanEvent(mealPlanID: string, mealPlanEventID: string): Promise<APIResponse<MealPlanEvent>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<MealPlanEvent>>(`/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`)
        .then((res: AxiosResponse<APIResponse<MealPlanEvent>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanEvent>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveMealPlanGroceryListItem(
    mealPlanID: string,
    mealPlanGroceryListItemID: string,
  ): Promise<APIResponse<MealPlanGroceryListItem>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<MealPlanGroceryListItem>>(
          `/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`,
        )
        .then((res: AxiosResponse<APIResponse<MealPlanGroceryListItem>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanGroceryListItem>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveMealPlanOption(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
  ): Promise<APIResponse<MealPlanOption>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<MealPlanOption>>(
          `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`,
        )
        .then((res: AxiosResponse<APIResponse<MealPlanOption>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanOption>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveMealPlanOptionVote(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
    mealPlanOptionVoteID: string,
  ): Promise<APIResponse<MealPlanOptionVote>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<MealPlanOptionVote>>(
          `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
        )
        .then((res: AxiosResponse<APIResponse<MealPlanOptionVote>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanOptionVote>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveOAuth2Client(oauth2ClientID: string): Promise<APIResponse<OAuth2Client>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<OAuth2Client>>(`/api/v1/oauth2_clients/${oauth2ClientID}`)
        .then((res: AxiosResponse<APIResponse<OAuth2Client>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<OAuth2Client>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveRecipe(recipeID: string): Promise<APIResponse<Recipe>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<Recipe>>(`/api/v1/recipes/${recipeID}`)
        .then((res: AxiosResponse<APIResponse<Recipe>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Recipe>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveRecipePrepTask(recipeID: string, recipePrepTaskID: string): Promise<APIResponse<RecipePrepTask>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<RecipePrepTask>>(`/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`)
        .then((res: AxiosResponse<APIResponse<RecipePrepTask>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipePrepTask>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveRecipeRating(recipeID: string, recipeRatingID: string): Promise<APIResponse<RecipeRating>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<RecipeRating>>(`/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`)
        .then((res: AxiosResponse<APIResponse<RecipeRating>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeRating>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveRecipeStep(recipeID: string, recipeStepID: string): Promise<APIResponse<RecipeStep>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<RecipeStep>>(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}`)
        .then((res: AxiosResponse<APIResponse<RecipeStep>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStep>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveRecipeStepCompletionCondition(
    recipeID: string,
    recipeStepID: string,
    recipeStepCompletionConditionID: string,
  ): Promise<APIResponse<RecipeStepCompletionCondition>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<RecipeStepCompletionCondition>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepCompletionCondition>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepCompletionCondition>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveRecipeStepIngredient(
    recipeID: string,
    recipeStepID: string,
    recipeStepIngredientID: string,
  ): Promise<APIResponse<RecipeStepIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<RecipeStepIngredient>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveRecipeStepInstrument(
    recipeID: string,
    recipeStepID: string,
    recipeStepInstrumentID: string,
  ): Promise<APIResponse<RecipeStepInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<RecipeStepInstrument>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveRecipeStepProduct(
    recipeID: string,
    recipeStepID: string,
    recipeStepProductID: string,
  ): Promise<APIResponse<RecipeStepProduct>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<RecipeStepProduct>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepProduct>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepProduct>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveRecipeStepVessel(
    recipeID: string,
    recipeStepID: string,
    recipeStepVesselID: string,
  ): Promise<APIResponse<RecipeStepVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<RecipeStepVessel>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveServiceSetting(serviceSettingID: string): Promise<APIResponse<ServiceSetting>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ServiceSetting>>(`/api/v1/settings/${serviceSettingID}`)
        .then((res: AxiosResponse<APIResponse<ServiceSetting>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ServiceSetting>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveServiceSettingConfiguration(
    serviceSettingConfigurationID: string,
  ): Promise<APIResponse<ServiceSettingConfiguration>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ServiceSettingConfiguration>>(
          `/api/v1/settings/configurations/${serviceSettingConfigurationID}`,
        )
        .then((res: AxiosResponse<APIResponse<ServiceSettingConfiguration>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ServiceSettingConfiguration>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveUser(userID: string): Promise<APIResponse<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<User>>(`/api/v1/users/${userID}`)
        .then((res: AxiosResponse<APIResponse<User>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<User>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveUserIngredientPreference(
    userIngredientPreferenceID: string,
  ): Promise<APIResponse<UserIngredientPreference>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<UserIngredientPreference>>(
          `/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`,
        )
        .then((res: AxiosResponse<APIResponse<UserIngredientPreference>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<UserIngredientPreference>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveUserMembership(householdID: string, userID: string): Promise<APIResponse<HouseholdUserMembership>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<HouseholdUserMembership>>(`/api/v1/households/${householdID}/members/${userID}`)
        .then((res: AxiosResponse<APIResponse<HouseholdUserMembership>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<HouseholdUserMembership>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidIngredient(validIngredientID: string): Promise<APIResponse<ValidIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidIngredient>>(`/api/v1/valid_ingredients/${validIngredientID}`)
        .then((res: AxiosResponse<APIResponse<ValidIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidIngredientGroup(validIngredientGroupID: string): Promise<APIResponse<ValidIngredientGroup>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidIngredientGroup>>(`/api/v1/valid_ingredient_groups/${validIngredientGroupID}`)
        .then((res: AxiosResponse<APIResponse<ValidIngredientGroup>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientGroup>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidIngredientMeasurementUnit(
    validIngredientMeasurementUnitID: string,
  ): Promise<APIResponse<ValidIngredientMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidIngredientMeasurementUnit>>(
          `/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`,
        )
        .then((res: AxiosResponse<APIResponse<ValidIngredientMeasurementUnit>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientMeasurementUnit>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidIngredientPreparation(
    validIngredientPreparationID: string,
  ): Promise<APIResponse<ValidIngredientPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidIngredientPreparation>>(
          `/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`,
        )
        .then((res: AxiosResponse<APIResponse<ValidIngredientPreparation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientPreparation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidIngredientState(validIngredientStateID: string): Promise<APIResponse<ValidIngredientState>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidIngredientState>>(`/api/v1/valid_ingredient_states/${validIngredientStateID}`)
        .then((res: AxiosResponse<APIResponse<ValidIngredientState>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientState>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidIngredientStateIngredient(
    validIngredientStateIngredientID: string,
  ): Promise<APIResponse<ValidIngredientStateIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidIngredientStateIngredient>>(
          `/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`,
        )
        .then((res: AxiosResponse<APIResponse<ValidIngredientStateIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientStateIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidInstrument(validInstrumentID: string): Promise<APIResponse<ValidInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidInstrument>>(`/api/v1/valid_instruments/${validInstrumentID}`)
        .then((res: AxiosResponse<APIResponse<ValidInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidMeasurementUnit(validMeasurementUnitID: string): Promise<APIResponse<ValidMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidMeasurementUnit>>(`/api/v1/valid_measurement_units/${validMeasurementUnitID}`)
        .then((res: AxiosResponse<APIResponse<ValidMeasurementUnit>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidMeasurementUnit>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidMeasurementUnitConversion(
    validMeasurementUnitConversionID: string,
  ): Promise<APIResponse<ValidMeasurementUnitConversion>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidMeasurementUnitConversion>>(
          `/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`,
        )
        .then((res: AxiosResponse<APIResponse<ValidMeasurementUnitConversion>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidMeasurementUnitConversion>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidPreparation(validPreparationID: string): Promise<APIResponse<ValidPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidPreparation>>(`/api/v1/valid_preparations/${validPreparationID}`)
        .then((res: AxiosResponse<APIResponse<ValidPreparation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidPreparationInstrument(
    validPreparationVesselID: string,
  ): Promise<APIResponse<ValidPreparationInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidPreparationInstrument>>(
          `/api/v1/valid_preparation_instruments/${validPreparationVesselID}`,
        )
        .then((res: AxiosResponse<APIResponse<ValidPreparationInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparationInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidPreparationVessel(validPreparationVesselID: string): Promise<APIResponse<ValidPreparationVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidPreparationVessel>>(`/api/v1/valid_preparation_vessels/${validPreparationVesselID}`)
        .then((res: AxiosResponse<APIResponse<ValidPreparationVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparationVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveValidVessel(validVesselID: string): Promise<APIResponse<ValidVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<ValidVessel>>(`/api/v1/valid_vessels/${validVesselID}`)
        .then((res: AxiosResponse<APIResponse<ValidVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveWebhook(webhookID: string): Promise<APIResponse<Webhook>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<Webhook>>(`/api/v1/webhooks/${webhookID}`)
        .then((res: AxiosResponse<APIResponse<Webhook>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Webhook>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async archiveWebhookTriggerEvent(
    webhookID: string,
    webhookTriggerEventID: string,
  ): Promise<APIResponse<WebhookTriggerEvent>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .delete<APIResponse<WebhookTriggerEvent>>(
          `/api/v1/webhooks/${webhookID}/trigger_events/${webhookTriggerEventID}`,
        )
        .then((res: AxiosResponse<APIResponse<WebhookTriggerEvent>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<WebhookTriggerEvent>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async cancelHouseholdInvitation(
    householdInvitationID: string,
    input: HouseholdInvitationUpdateRequestInput,
  ): Promise<APIResponse<HouseholdInvitation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<HouseholdInvitation>>(`/api/v1/household_invitations/${householdInvitationID}/cancel`, input)
        .then((res: AxiosResponse<APIResponse<HouseholdInvitation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<HouseholdInvitation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async checkPermissions(input: UserPermissionsRequestInput): Promise<APIResponse<UserPermissionsResponse>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<UserPermissionsResponse>>(`/api/v1/users/permissions/check`, input)
        .then((res: AxiosResponse<APIResponse<UserPermissionsResponse>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<UserPermissionsResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async cloneRecipe(recipeID: string): Promise<APIResponse<Recipe>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<Recipe>>(`/api/v1/recipes/${recipeID}/clone`)
        .then((res: AxiosResponse<APIResponse<Recipe>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Recipe>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createHousehold(input: HouseholdCreationRequestInput): Promise<APIResponse<Household>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<Household>>(`/api/v1/households`, input)
        .then((res: AxiosResponse<APIResponse<Household>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Household>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createHouseholdInstrumentOwnership(
    input: HouseholdInstrumentOwnershipCreationRequestInput,
  ): Promise<APIResponse<HouseholdInstrumentOwnership>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<HouseholdInstrumentOwnership>>(`/api/v1/households/instruments`, input)
        .then((res: AxiosResponse<APIResponse<HouseholdInstrumentOwnership>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<HouseholdInstrumentOwnership>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createHouseholdInvitation(
    householdID: string,
    input: HouseholdInvitationCreationRequestInput,
  ): Promise<APIResponse<HouseholdInvitation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<HouseholdInvitation>>(`/api/v1/households/${householdID}/invite`, input)
        .then((res: AxiosResponse<APIResponse<HouseholdInvitation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<HouseholdInvitation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createMeal(input: MealCreationRequestInput): Promise<APIResponse<Meal>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<Meal>>(`/api/v1/meals`, input)
        .then((res: AxiosResponse<APIResponse<Meal>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Meal>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createMealPlan(input: MealPlanCreationRequestInput): Promise<APIResponse<MealPlan>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<MealPlan>>(`/api/v1/meal_plans`, input)
        .then((res: AxiosResponse<APIResponse<MealPlan>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlan>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createMealPlanEvent(
    mealPlanID: string,
    input: MealPlanEventCreationRequestInput,
  ): Promise<APIResponse<MealPlanEvent>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<MealPlanEvent>>(`/api/v1/meal_plans/${mealPlanID}/events`, input)
        .then((res: AxiosResponse<APIResponse<MealPlanEvent>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanEvent>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createMealPlanGroceryListItem(
    mealPlanID: string,
    input: MealPlanGroceryListItemCreationRequestInput,
  ): Promise<APIResponse<MealPlanGroceryListItem>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<MealPlanGroceryListItem>>(`/api/v1/meal_plans/${mealPlanID}/grocery_list_items`, input)
        .then((res: AxiosResponse<APIResponse<MealPlanGroceryListItem>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanGroceryListItem>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createMealPlanOption(
    mealPlanID: string,
    mealPlanEventID: string,
    input: MealPlanOptionCreationRequestInput,
  ): Promise<APIResponse<MealPlanOption>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<MealPlanOption>>(`/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`, input)
        .then((res: AxiosResponse<APIResponse<MealPlanOption>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanOption>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createMealPlanOptionVote(
    mealPlanID: string,
    mealPlanEventID: string,
    input: MealPlanOptionVoteCreationRequestInput,
  ): Promise<APIResponse<Array<MealPlanOptionVote>>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<Array<MealPlanOptionVote>>>(
          `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/vote`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<Array<MealPlanOptionVote>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Array<MealPlanOptionVote>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createMealPlanTask(
    mealPlanID: string,
    input: MealPlanTaskCreationRequestInput,
  ): Promise<APIResponse<MealPlanTask>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<MealPlanTask>>(`/api/v1/meal_plans/${mealPlanID}/tasks`, input)
        .then((res: AxiosResponse<APIResponse<MealPlanTask>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanTask>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createOAuth2Client(
    input: OAuth2ClientCreationRequestInput,
  ): Promise<APIResponse<OAuth2ClientCreationResponse>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<OAuth2ClientCreationResponse>>(`/api/v1/oauth2_clients`, input)
        .then((res: AxiosResponse<APIResponse<OAuth2ClientCreationResponse>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<OAuth2ClientCreationResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createRecipe(input: RecipeCreationRequestInput): Promise<APIResponse<Recipe>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<Recipe>>(`/api/v1/recipes`, input)
        .then((res: AxiosResponse<APIResponse<Recipe>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Recipe>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createRecipePrepTask(
    recipeID: string,
    input: RecipePrepTaskCreationRequestInput,
  ): Promise<APIResponse<RecipePrepTask>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<RecipePrepTask>>(`/api/v1/recipes/${recipeID}/prep_tasks`, input)
        .then((res: AxiosResponse<APIResponse<RecipePrepTask>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipePrepTask>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createRecipeRating(
    recipeID: string,
    input: RecipeRatingCreationRequestInput,
  ): Promise<APIResponse<RecipeRating>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<RecipeRating>>(`/api/v1/recipes/${recipeID}/ratings`, input)
        .then((res: AxiosResponse<APIResponse<RecipeRating>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeRating>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createRecipeStep(recipeID: string, input: RecipeStepCreationRequestInput): Promise<APIResponse<RecipeStep>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<RecipeStep>>(`/api/v1/recipes/${recipeID}/steps`, input)
        .then((res: AxiosResponse<APIResponse<RecipeStep>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStep>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createRecipeStepCompletionCondition(
    recipeID: string,
    recipeStepID: string,
    input: RecipeStepCompletionConditionForExistingRecipeCreationRequestInput,
  ): Promise<APIResponse<RecipeStepCompletionCondition>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<RecipeStepCompletionCondition>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepCompletionCondition>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepCompletionCondition>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createRecipeStepIngredient(
    recipeID: string,
    recipeStepID: string,
    input: RecipeStepIngredientCreationRequestInput,
  ): Promise<APIResponse<RecipeStepIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<RecipeStepIngredient>>(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`, input)
        .then((res: AxiosResponse<APIResponse<RecipeStepIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createRecipeStepInstrument(
    recipeID: string,
    recipeStepID: string,
    input: RecipeStepInstrumentCreationRequestInput,
  ): Promise<APIResponse<RecipeStepInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<RecipeStepInstrument>>(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`, input)
        .then((res: AxiosResponse<APIResponse<RecipeStepInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createRecipeStepProduct(
    recipeID: string,
    recipeStepID: string,
    input: RecipeStepProductCreationRequestInput,
  ): Promise<APIResponse<RecipeStepProduct>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<RecipeStepProduct>>(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`, input)
        .then((res: AxiosResponse<APIResponse<RecipeStepProduct>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepProduct>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createRecipeStepVessel(
    recipeID: string,
    recipeStepID: string,
    input: RecipeStepVesselCreationRequestInput,
  ): Promise<APIResponse<RecipeStepVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<RecipeStepVessel>>(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`, input)
        .then((res: AxiosResponse<APIResponse<RecipeStepVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createServiceSetting(input: ServiceSettingCreationRequestInput): Promise<APIResponse<ServiceSetting>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ServiceSetting>>(`/api/v1/settings`, input)
        .then((res: AxiosResponse<APIResponse<ServiceSetting>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ServiceSetting>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createServiceSettingConfiguration(
    input: ServiceSettingConfigurationCreationRequestInput,
  ): Promise<APIResponse<ServiceSettingConfiguration>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ServiceSettingConfiguration>>(`/api/v1/settings/configurations`, input)
        .then((res: AxiosResponse<APIResponse<ServiceSettingConfiguration>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ServiceSettingConfiguration>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createUser(input: UserRegistrationInput): Promise<APIResponse<UserCreationResponse>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<UserCreationResponse>>(`/users`, input)
        .then((res: AxiosResponse<APIResponse<UserCreationResponse>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<UserCreationResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createUserIngredientPreference(
    input: UserIngredientPreferenceCreationRequestInput,
  ): Promise<APIResponse<Array<UserIngredientPreference>>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<Array<UserIngredientPreference>>>(`/api/v1/user_ingredient_preferences`, input)
        .then((res: AxiosResponse<APIResponse<Array<UserIngredientPreference>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Array<UserIngredientPreference>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createUserNotification(input: UserNotificationCreationRequestInput): Promise<APIResponse<UserNotification>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<UserNotification>>(`/api/v1/user_notifications`, input)
        .then((res: AxiosResponse<APIResponse<UserNotification>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<UserNotification>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidIngredient(input: ValidIngredientCreationRequestInput): Promise<APIResponse<ValidIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidIngredient>>(`/api/v1/valid_ingredients`, input)
        .then((res: AxiosResponse<APIResponse<ValidIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidIngredientGroup(
    input: ValidIngredientGroupCreationRequestInput,
  ): Promise<APIResponse<ValidIngredientGroup>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidIngredientGroup>>(`/api/v1/valid_ingredient_groups`, input)
        .then((res: AxiosResponse<APIResponse<ValidIngredientGroup>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientGroup>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidIngredientMeasurementUnit(
    input: ValidIngredientMeasurementUnitCreationRequestInput,
  ): Promise<APIResponse<ValidIngredientMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidIngredientMeasurementUnit>>(`/api/v1/valid_ingredient_measurement_units`, input)
        .then((res: AxiosResponse<APIResponse<ValidIngredientMeasurementUnit>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientMeasurementUnit>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidIngredientPreparation(
    input: ValidIngredientPreparationCreationRequestInput,
  ): Promise<APIResponse<ValidIngredientPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidIngredientPreparation>>(`/api/v1/valid_ingredient_preparations`, input)
        .then((res: AxiosResponse<APIResponse<ValidIngredientPreparation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientPreparation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidIngredientState(
    input: ValidIngredientStateCreationRequestInput,
  ): Promise<APIResponse<ValidIngredientState>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidIngredientState>>(`/api/v1/valid_ingredient_states`, input)
        .then((res: AxiosResponse<APIResponse<ValidIngredientState>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientState>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidIngredientStateIngredient(
    input: ValidIngredientStateIngredientCreationRequestInput,
  ): Promise<APIResponse<ValidIngredientStateIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidIngredientStateIngredient>>(`/api/v1/valid_ingredient_state_ingredients`, input)
        .then((res: AxiosResponse<APIResponse<ValidIngredientStateIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientStateIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidInstrument(input: ValidInstrumentCreationRequestInput): Promise<APIResponse<ValidInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidInstrument>>(`/api/v1/valid_instruments`, input)
        .then((res: AxiosResponse<APIResponse<ValidInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidMeasurementUnit(
    input: ValidMeasurementUnitCreationRequestInput,
  ): Promise<APIResponse<ValidMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidMeasurementUnit>>(`/api/v1/valid_measurement_units`, input)
        .then((res: AxiosResponse<APIResponse<ValidMeasurementUnit>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidMeasurementUnit>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidMeasurementUnitConversion(
    input: ValidMeasurementUnitConversionCreationRequestInput,
  ): Promise<APIResponse<ValidMeasurementUnitConversion>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidMeasurementUnitConversion>>(`/api/v1/valid_measurement_conversions`, input)
        .then((res: AxiosResponse<APIResponse<ValidMeasurementUnitConversion>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidMeasurementUnitConversion>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidPreparation(input: ValidPreparationCreationRequestInput): Promise<APIResponse<ValidPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidPreparation>>(`/api/v1/valid_preparations`, input)
        .then((res: AxiosResponse<APIResponse<ValidPreparation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidPreparationInstrument(
    input: ValidPreparationInstrumentCreationRequestInput,
  ): Promise<APIResponse<ValidPreparationInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidPreparationInstrument>>(`/api/v1/valid_preparation_instruments`, input)
        .then((res: AxiosResponse<APIResponse<ValidPreparationInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparationInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidPreparationVessel(
    input: ValidPreparationVesselCreationRequestInput,
  ): Promise<APIResponse<ValidPreparationVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidPreparationVessel>>(`/api/v1/valid_preparation_vessels`, input)
        .then((res: AxiosResponse<APIResponse<ValidPreparationVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparationVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createValidVessel(input: ValidVesselCreationRequestInput): Promise<APIResponse<ValidVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<ValidVessel>>(`/api/v1/valid_vessels`, input)
        .then((res: AxiosResponse<APIResponse<ValidVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createWebhook(input: WebhookCreationRequestInput): Promise<APIResponse<Webhook>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<Webhook>>(`/api/v1/webhooks`, input)
        .then((res: AxiosResponse<APIResponse<Webhook>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Webhook>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async createWebhookTriggerEvent(
    webhookID: string,
    input: WebhookTriggerEventCreationRequestInput,
  ): Promise<APIResponse<WebhookTriggerEvent>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<WebhookTriggerEvent>>(`/api/v1/webhooks/${webhookID}/trigger_events`, input)
        .then((res: AxiosResponse<APIResponse<WebhookTriggerEvent>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<WebhookTriggerEvent>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async finalizeMealPlan(mealPlanID: string): Promise<APIResponse<FinalizeMealPlansResponse>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<FinalizeMealPlansResponse>>(`/api/v1/meal_plans/${mealPlanID}/finalize`)
        .then((res: AxiosResponse<APIResponse<FinalizeMealPlansResponse>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<FinalizeMealPlansResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getActiveHousehold(): Promise<APIResponse<Household>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Household>>(`/api/v1/households/current`)
        .then((res: AxiosResponse<APIResponse<Household>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Household>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getAuditLogEntriesForHousehold(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<AuditLogEntry>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<AuditLogEntry>>>(`/api/v1/audit_log_entries/for_household`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<AuditLogEntry>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<AuditLogEntry>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<AuditLogEntry>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getAuditLogEntriesForUser(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<AuditLogEntry>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<AuditLogEntry>>>(`/api/v1/audit_log_entries/for_user`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<AuditLogEntry>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<AuditLogEntry>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<AuditLogEntry>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getAuditLogEntryByID(auditLogEntryID: string): Promise<APIResponse<AuditLogEntry>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<AuditLogEntry>>(`/api/v1/audit_log_entries/${auditLogEntryID}`)
        .then((res: AxiosResponse<APIResponse<AuditLogEntry>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<AuditLogEntry>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getAuthStatus(): Promise<APIResponse<UserStatusResponse>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<UserStatusResponse>>(`/auth/status`)
        .then((res: AxiosResponse<APIResponse<UserStatusResponse>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<UserStatusResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getHousehold(householdID: string): Promise<APIResponse<Household>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Household>>(`/api/v1/households/${householdID}`)
        .then((res: AxiosResponse<APIResponse<Household>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Household>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getHouseholdInstrumentOwnership(
    householdInstrumentOwnershipID: string,
  ): Promise<APIResponse<HouseholdInstrumentOwnership>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<HouseholdInstrumentOwnership>>(
          `/api/v1/households/instruments/${householdInstrumentOwnershipID}`,
        )
        .then((res: AxiosResponse<APIResponse<HouseholdInstrumentOwnership>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<HouseholdInstrumentOwnership>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getHouseholdInstrumentOwnerships(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<HouseholdInstrumentOwnership>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<HouseholdInstrumentOwnership>>>(`/api/v1/households/instruments`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<HouseholdInstrumentOwnership>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<HouseholdInstrumentOwnership>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<HouseholdInstrumentOwnership>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getHouseholdInvitation(householdInvitationID: string): Promise<APIResponse<HouseholdInvitation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<HouseholdInvitation>>(`/api/v1/household_invitations/${householdInvitationID}`)
        .then((res: AxiosResponse<APIResponse<HouseholdInvitation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<HouseholdInvitation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getHouseholdInvitationByID(
    householdID: string,
    householdInvitationID: string,
  ): Promise<APIResponse<HouseholdInvitation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<HouseholdInvitation>>(`/api/v1/households/${householdID}/invitations/${householdInvitationID}`)
        .then((res: AxiosResponse<APIResponse<HouseholdInvitation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<HouseholdInvitation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getHouseholds(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Household>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<Household>>>(`/api/v1/households`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<Household>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<Household>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<Household>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMeal(mealID: string): Promise<APIResponse<Meal>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Meal>>(`/api/v1/meals/${mealID}`)
        .then((res: AxiosResponse<APIResponse<Meal>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Meal>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMealPlan(mealPlanID: string): Promise<APIResponse<MealPlan>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<MealPlan>>(`/api/v1/meal_plans/${mealPlanID}`)
        .then((res: AxiosResponse<APIResponse<MealPlan>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlan>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMealPlanEvent(mealPlanID: string, mealPlanEventID: string): Promise<APIResponse<MealPlanEvent>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<MealPlanEvent>>(`/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`)
        .then((res: AxiosResponse<APIResponse<MealPlanEvent>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanEvent>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMealPlanEvents(
    mealPlanID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<MealPlanEvent>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<MealPlanEvent>>>(`/api/v1/meal_plans/${mealPlanID}/events`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<MealPlanEvent>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<MealPlanEvent>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<MealPlanEvent>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMealPlanGroceryListItem(
    mealPlanID: string,
    mealPlanGroceryListItemID: string,
  ): Promise<APIResponse<MealPlanGroceryListItem>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<MealPlanGroceryListItem>>(
          `/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`,
        )
        .then((res: AxiosResponse<APIResponse<MealPlanGroceryListItem>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanGroceryListItem>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMealPlanGroceryListItemsForMealPlan(
    mealPlanID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<MealPlanGroceryListItem>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<MealPlanGroceryListItem>>>(`/api/v1/meal_plans/${mealPlanID}/grocery_list_items`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<MealPlanGroceryListItem>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<MealPlanGroceryListItem>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<MealPlanGroceryListItem>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMealPlanOption(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
  ): Promise<APIResponse<MealPlanOption>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<MealPlanOption>>(
          `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`,
        )
        .then((res: AxiosResponse<APIResponse<MealPlanOption>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanOption>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMealPlanOptionVote(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
    mealPlanOptionVoteID: string,
  ): Promise<APIResponse<MealPlanOptionVote>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<MealPlanOptionVote>>(
          `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
        )
        .then((res: AxiosResponse<APIResponse<MealPlanOptionVote>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanOptionVote>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMealPlanOptionVotes(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<MealPlanOptionVote>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<MealPlanOptionVote>>>(
          `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<MealPlanOptionVote>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<MealPlanOptionVote>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<MealPlanOptionVote>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMealPlanOptions(
    mealPlanID: string,
    mealPlanEventID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<MealPlanOption>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<MealPlanOption>>>(`/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<MealPlanOption>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<MealPlanOption>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<MealPlanOption>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMealPlanTask(mealPlanID: string, mealPlanTaskID: string): Promise<APIResponse<MealPlanTask>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<MealPlanTask>>(`/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`)
        .then((res: AxiosResponse<APIResponse<MealPlanTask>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanTask>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMealPlanTasks(
    mealPlanID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<MealPlanTask>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<MealPlanTask>>>(`/api/v1/meal_plans/${mealPlanID}/tasks`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<MealPlanTask>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<MealPlanTask>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<MealPlanTask>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMealPlans(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<MealPlan>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<MealPlan>>>(`/api/v1/meal_plans`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<MealPlan>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<MealPlan>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<MealPlan>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMeals(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Meal>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<Meal>>>(`/api/v1/meals`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<Meal>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<Meal>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<Meal>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getMermaidDiagramForRecipe(recipeID: string): Promise<APIResponse<string>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<string>>(`/api/v1/recipes/${recipeID}/mermaid`)
        .then((res: AxiosResponse<APIResponse<string>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<string>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getOAuth2Client(oauth2ClientID: string): Promise<APIResponse<OAuth2Client>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<OAuth2Client>>(`/api/v1/oauth2_clients/${oauth2ClientID}`)
        .then((res: AxiosResponse<APIResponse<OAuth2Client>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<OAuth2Client>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getOAuth2Clients(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<OAuth2Client>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<OAuth2Client>>>(`/api/v1/oauth2_clients`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<OAuth2Client>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<OAuth2Client>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<OAuth2Client>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRandomValidIngredient(): Promise<APIResponse<ValidIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidIngredient>>(`/api/v1/valid_ingredients/random`)
        .then((res: AxiosResponse<APIResponse<ValidIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRandomValidInstrument(): Promise<APIResponse<ValidInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidInstrument>>(`/api/v1/valid_instruments/random`)
        .then((res: AxiosResponse<APIResponse<ValidInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRandomValidPreparation(): Promise<APIResponse<ValidPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidPreparation>>(`/api/v1/valid_preparations/random`)
        .then((res: AxiosResponse<APIResponse<ValidPreparation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRandomValidVessel(): Promise<APIResponse<ValidVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidVessel>>(`/api/v1/valid_vessels/random`)
        .then((res: AxiosResponse<APIResponse<ValidVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getReceivedHouseholdInvitations(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<HouseholdInvitation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<HouseholdInvitation>>>(`/api/v1/household_invitations/received`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<HouseholdInvitation>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<HouseholdInvitation>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<HouseholdInvitation>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipe(recipeID: string): Promise<APIResponse<Recipe>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Recipe>>(`/api/v1/recipes/${recipeID}`)
        .then((res: AxiosResponse<APIResponse<Recipe>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Recipe>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeMealPlanTasks(recipeID: string): Promise<APIResponse<RecipePrepTaskStep>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<RecipePrepTaskStep>>(`/api/v1/recipes/${recipeID}/prep_steps`)
        .then((res: AxiosResponse<APIResponse<RecipePrepTaskStep>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipePrepTaskStep>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipePrepTask(recipeID: string, recipePrepTaskID: string): Promise<APIResponse<RecipePrepTask>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<RecipePrepTask>>(`/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`)
        .then((res: AxiosResponse<APIResponse<RecipePrepTask>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipePrepTask>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipePrepTasks(
    recipeID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<RecipePrepTask>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<RecipePrepTask>>>(`/api/v1/recipes/${recipeID}/prep_tasks`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<RecipePrepTask>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<RecipePrepTask>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<RecipePrepTask>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeRating(recipeID: string, recipeRatingID: string): Promise<APIResponse<RecipeRating>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<RecipeRating>>(`/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`)
        .then((res: AxiosResponse<APIResponse<RecipeRating>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeRating>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeRatings(
    recipeID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<RecipeRating>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<RecipeRating>>>(`/api/v1/recipes/${recipeID}/ratings`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<RecipeRating>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<RecipeRating>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<RecipeRating>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeStep(recipeID: string, recipeStepID: string): Promise<APIResponse<RecipeStep>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<RecipeStep>>(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}`)
        .then((res: AxiosResponse<APIResponse<RecipeStep>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStep>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeStepCompletionCondition(
    recipeID: string,
    recipeStepID: string,
    recipeStepCompletionConditionID: string,
  ): Promise<APIResponse<RecipeStepCompletionCondition>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<RecipeStepCompletionCondition>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepCompletionCondition>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepCompletionCondition>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeStepCompletionConditions(
    recipeID: string,
    recipeStepID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<RecipeStepCompletionCondition>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<RecipeStepCompletionCondition>>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<RecipeStepCompletionCondition>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<RecipeStepCompletionCondition>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<RecipeStepCompletionCondition>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeStepIngredient(
    recipeID: string,
    recipeStepID: string,
    recipeStepIngredientID: string,
  ): Promise<APIResponse<RecipeStepIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<RecipeStepIngredient>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeStepIngredients(
    recipeID: string,
    recipeStepID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<RecipeStepIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<RecipeStepIngredient>>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<RecipeStepIngredient>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<RecipeStepIngredient>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<RecipeStepIngredient>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeStepInstrument(
    recipeID: string,
    recipeStepID: string,
    recipeStepInstrumentID: string,
  ): Promise<APIResponse<RecipeStepInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<RecipeStepInstrument>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeStepInstruments(
    recipeID: string,
    recipeStepID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<RecipeStepInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<RecipeStepInstrument>>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<RecipeStepInstrument>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<RecipeStepInstrument>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<RecipeStepInstrument>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeStepProduct(
    recipeID: string,
    recipeStepID: string,
    recipeStepProductID: string,
  ): Promise<APIResponse<RecipeStepProduct>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<RecipeStepProduct>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepProduct>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepProduct>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeStepProducts(
    recipeID: string,
    recipeStepID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<RecipeStepProduct>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<RecipeStepProduct>>>(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<RecipeStepProduct>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<RecipeStepProduct>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<RecipeStepProduct>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeStepVessel(
    recipeID: string,
    recipeStepID: string,
    recipeStepVesselID: string,
  ): Promise<APIResponse<RecipeStepVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<RecipeStepVessel>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeStepVessels(
    recipeID: string,
    recipeStepID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<RecipeStepVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<RecipeStepVessel>>>(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<RecipeStepVessel>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<RecipeStepVessel>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<RecipeStepVessel>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipeSteps(
    recipeID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<RecipeStep>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<RecipeStep>>>(`/api/v1/recipes/${recipeID}/steps`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<RecipeStep>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<RecipeStep>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<RecipeStep>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getRecipes(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Recipe>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<Recipe>>>(`/api/v1/recipes`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<Recipe>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<Recipe>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<Recipe>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getSelf(): Promise<APIResponse<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<User>>(`/api/v1/users/self`)
        .then((res: AxiosResponse<APIResponse<User>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<User>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getSentHouseholdInvitations(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<HouseholdInvitation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<HouseholdInvitation>>>(`/api/v1/household_invitations/sent`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<HouseholdInvitation>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<HouseholdInvitation>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<HouseholdInvitation>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getServiceSetting(serviceSettingID: string): Promise<APIResponse<ServiceSetting>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ServiceSetting>>(`/api/v1/settings/${serviceSettingID}`)
        .then((res: AxiosResponse<APIResponse<ServiceSetting>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ServiceSetting>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getServiceSettingConfigurationByName(
    serviceSettingConfigurationName: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ServiceSettingConfiguration>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ServiceSettingConfiguration>>>(
          `/api/v1/settings/configurations/user/${serviceSettingConfigurationName}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ServiceSettingConfiguration>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ServiceSettingConfiguration>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ServiceSettingConfiguration>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getServiceSettingConfigurationsForHousehold(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ServiceSettingConfiguration>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ServiceSettingConfiguration>>>(`/api/v1/settings/configurations/household`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ServiceSettingConfiguration>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ServiceSettingConfiguration>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ServiceSettingConfiguration>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getServiceSettingConfigurationsForUser(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ServiceSettingConfiguration>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ServiceSettingConfiguration>>>(`/api/v1/settings/configurations/user`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ServiceSettingConfiguration>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ServiceSettingConfiguration>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ServiceSettingConfiguration>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getServiceSettings(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<ServiceSetting>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ServiceSetting>>>(`/api/v1/settings`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ServiceSetting>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ServiceSetting>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ServiceSetting>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getUser(userID: string): Promise<APIResponse<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<User>>(`/api/v1/users/${userID}`)
        .then((res: AxiosResponse<APIResponse<User>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<User>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getUserIngredientPreferences(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<UserIngredientPreference>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<UserIngredientPreference>>>(`/api/v1/user_ingredient_preferences`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<UserIngredientPreference>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<UserIngredientPreference>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<UserIngredientPreference>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getUserNotification(userNotificationID: string): Promise<APIResponse<UserNotification>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<UserNotification>>(`/api/v1/user_notifications/${userNotificationID}`)
        .then((res: AxiosResponse<APIResponse<UserNotification>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<UserNotification>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getUserNotifications(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<UserNotification>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<UserNotification>>>(`/api/v1/user_notifications`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<UserNotification>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<UserNotification>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<UserNotification>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getUsers(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<User>>>(`/api/v1/users`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<User>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<User>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<User>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredient(validIngredientID: string): Promise<APIResponse<ValidIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidIngredient>>(`/api/v1/valid_ingredients/${validIngredientID}`)
        .then((res: AxiosResponse<APIResponse<ValidIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientGroup(validIngredientGroupID: string): Promise<APIResponse<ValidIngredientGroup>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidIngredientGroup>>(`/api/v1/valid_ingredient_groups/${validIngredientGroupID}`)
        .then((res: AxiosResponse<APIResponse<ValidIngredientGroup>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientGroup>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientGroups(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientGroup>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientGroup>>>(`/api/v1/valid_ingredient_groups`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientGroup>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientGroup>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientGroup>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientMeasurementUnit(
    validIngredientMeasurementUnitID: string,
  ): Promise<APIResponse<ValidIngredientMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidIngredientMeasurementUnit>>(
          `/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`,
        )
        .then((res: AxiosResponse<APIResponse<ValidIngredientMeasurementUnit>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientMeasurementUnit>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientMeasurementUnits(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientMeasurementUnit>>>(`/api/v1/valid_ingredient_measurement_units`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientMeasurementUnit>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientMeasurementUnit>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientMeasurementUnit>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientMeasurementUnitsByIngredient(
    validIngredientID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientMeasurementUnit>>>(
          `/api/v1/valid_ingredient_measurement_units/by_ingredient/${validIngredientID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientMeasurementUnit>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientMeasurementUnit>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientMeasurementUnit>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientMeasurementUnitsByMeasurementUnit(
    validMeasurementUnitID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientMeasurementUnit>>>(
          `/api/v1/valid_ingredient_measurement_units/by_measurement_unit/${validMeasurementUnitID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientMeasurementUnit>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientMeasurementUnit>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientMeasurementUnit>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientPreparation(
    validIngredientPreparationID: string,
  ): Promise<APIResponse<ValidIngredientPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidIngredientPreparation>>(
          `/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`,
        )
        .then((res: AxiosResponse<APIResponse<ValidIngredientPreparation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientPreparation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientPreparations(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientPreparation>>>(`/api/v1/valid_ingredient_preparations`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientPreparation>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientPreparation>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientPreparation>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientPreparationsByIngredient(
    validIngredientID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientPreparation>>>(
          `/api/v1/valid_ingredient_preparations/by_ingredient/${validIngredientID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientPreparation>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientPreparation>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientPreparation>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientPreparationsByPreparation(
    validPreparationID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientPreparation>>>(
          `/api/v1/valid_ingredient_preparations/by_preparation/${validPreparationID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientPreparation>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientPreparation>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientPreparation>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientState(validIngredientStateID: string): Promise<APIResponse<ValidIngredientState>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidIngredientState>>(`/api/v1/valid_ingredient_states/${validIngredientStateID}`)
        .then((res: AxiosResponse<APIResponse<ValidIngredientState>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientState>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientStateIngredient(
    validIngredientStateIngredientID: string,
  ): Promise<APIResponse<ValidIngredientStateIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidIngredientStateIngredient>>(
          `/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`,
        )
        .then((res: AxiosResponse<APIResponse<ValidIngredientStateIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientStateIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientStateIngredients(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientStateIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientStateIngredient>>>(`/api/v1/valid_ingredient_state_ingredients`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientStateIngredient>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientStateIngredient>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientStateIngredient>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientStateIngredientsByIngredient(
    validIngredientID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientStateIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientStateIngredient>>>(
          `/api/v1/valid_ingredient_state_ingredients/by_ingredient/${validIngredientID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientStateIngredient>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientStateIngredient>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientStateIngredient>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientStateIngredientsByIngredientState(
    validIngredientStateID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientStateIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientStateIngredient>>>(
          `/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/${validIngredientStateID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientStateIngredient>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientStateIngredient>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientStateIngredient>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredientStates(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientState>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientState>>>(`/api/v1/valid_ingredient_states`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientState>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientState>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientState>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidIngredients(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredient>>>(`/api/v1/valid_ingredients`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredient>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredient>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredient>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidInstrument(validInstrumentID: string): Promise<APIResponse<ValidInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidInstrument>>(`/api/v1/valid_instruments/${validInstrumentID}`)
        .then((res: AxiosResponse<APIResponse<ValidInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidInstruments(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidInstrument>>>(`/api/v1/valid_instruments`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidInstrument>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidInstrument>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidInstrument>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidMeasurementUnit(validMeasurementUnitID: string): Promise<APIResponse<ValidMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidMeasurementUnit>>(`/api/v1/valid_measurement_units/${validMeasurementUnitID}`)
        .then((res: AxiosResponse<APIResponse<ValidMeasurementUnit>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidMeasurementUnit>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidMeasurementUnitConversion(
    validMeasurementUnitConversionID: string,
  ): Promise<APIResponse<ValidMeasurementUnitConversion>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidMeasurementUnitConversion>>(
          `/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`,
        )
        .then((res: AxiosResponse<APIResponse<ValidMeasurementUnitConversion>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidMeasurementUnitConversion>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidMeasurementUnitConversionsFromUnit(
    validMeasurementUnitID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidMeasurementUnitConversion>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidMeasurementUnitConversion>>>(
          `/api/v1/valid_measurement_conversions/from_unit/${validMeasurementUnitID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidMeasurementUnitConversion>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidMeasurementUnitConversion>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidMeasurementUnitConversion>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidMeasurementUnitConversionsToUnit(
    validMeasurementUnitID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidMeasurementUnitConversion>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidMeasurementUnitConversion>>>(
          `/api/v1/valid_measurement_conversions/to_unit/${validMeasurementUnitID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidMeasurementUnitConversion>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidMeasurementUnitConversion>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidMeasurementUnitConversion>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidMeasurementUnits(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidMeasurementUnit>>>(`/api/v1/valid_measurement_units`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidMeasurementUnit>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidMeasurementUnit>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidMeasurementUnit>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidPreparation(validPreparationID: string): Promise<APIResponse<ValidPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidPreparation>>(`/api/v1/valid_preparations/${validPreparationID}`)
        .then((res: AxiosResponse<APIResponse<ValidPreparation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidPreparationInstrument(
    validPreparationVesselID: string,
  ): Promise<APIResponse<ValidPreparationInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidPreparationInstrument>>(
          `/api/v1/valid_preparation_instruments/${validPreparationVesselID}`,
        )
        .then((res: AxiosResponse<APIResponse<ValidPreparationInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparationInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidPreparationInstruments(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidPreparationInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidPreparationInstrument>>>(`/api/v1/valid_preparation_instruments`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidPreparationInstrument>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidPreparationInstrument>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidPreparationInstrument>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidPreparationInstrumentsByInstrument(
    validInstrumentID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidPreparationInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidPreparationInstrument>>>(
          `/api/v1/valid_preparation_instruments/by_instrument/${validInstrumentID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidPreparationInstrument>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidPreparationInstrument>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidPreparationInstrument>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidPreparationInstrumentsByPreparation(
    validPreparationID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidPreparationInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidPreparationInstrument>>>(
          `/api/v1/valid_preparation_instruments/by_preparation/${validPreparationID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidPreparationInstrument>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidPreparationInstrument>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidPreparationInstrument>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidPreparationVessel(validPreparationVesselID: string): Promise<APIResponse<ValidPreparationVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidPreparationVessel>>(`/api/v1/valid_preparation_vessels/${validPreparationVesselID}`)
        .then((res: AxiosResponse<APIResponse<ValidPreparationVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparationVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidPreparationVessels(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidPreparationVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidPreparationVessel>>>(`/api/v1/valid_preparation_vessels`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidPreparationVessel>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidPreparationVessel>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidPreparationVessel>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidPreparationVesselsByPreparation(
    validPreparationID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidPreparationVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidPreparationVessel>>>(
          `/api/v1/valid_preparation_vessels/by_preparation/${validPreparationID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidPreparationVessel>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidPreparationVessel>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidPreparationVessel>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidPreparationVesselsByVessel(
    ValidVesselID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidPreparationVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidPreparationVessel>>>(
          `/api/v1/valid_preparation_vessels/by_vessel/${ValidVesselID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidPreparationVessel>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidPreparationVessel>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidPreparationVessel>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidPreparations(
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidPreparation>>>(`/api/v1/valid_preparations`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidPreparation>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidPreparation>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidPreparation>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidVessel(validVesselID: string): Promise<APIResponse<ValidVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<ValidVessel>>(`/api/v1/valid_vessels/${validVesselID}`)
        .then((res: AxiosResponse<APIResponse<ValidVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getValidVessels(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<ValidVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidVessel>>>(`/api/v1/valid_vessels`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidVessel>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidVessel>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidVessel>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getWebhook(webhookID: string): Promise<APIResponse<Webhook>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Webhook>>(`/api/v1/webhooks/${webhookID}`)
        .then((res: AxiosResponse<APIResponse<Webhook>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Webhook>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async getWebhooks(filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Webhook>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<Webhook>>>(`/api/v1/webhooks`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<Webhook>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<Webhook>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<Webhook>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async loginForJWT(input: UserLoginInput): Promise<AxiosResponse<APIResponse<JWTResponse>>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<JWTResponse>>(`/users/login/jwt`, input)
        .then((res: AxiosResponse<APIResponse<JWTResponse>>) => {
          if (res.data.error && res.data.error.message.toLowerCase() != 'totp required') {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res);
          }
        })
        .catch((error: AxiosError<APIResponse<JWTResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async redeemPasswordResetToken(input: PasswordResetTokenRedemptionRequestInput): Promise<APIResponse<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<User>>(`/users/password/reset/redeem`, input)
        .then((res: AxiosResponse<APIResponse<User>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<User>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async refreshTOTPSecret(input: TOTPSecretRefreshInput): Promise<APIResponse<TOTPSecretRefreshResponse>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<TOTPSecretRefreshResponse>>(`/api/v1/users/totp_secret/new`, input)
        .then((res: AxiosResponse<APIResponse<TOTPSecretRefreshResponse>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<TOTPSecretRefreshResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async rejectHouseholdInvitation(
    householdInvitationID: string,
    input: HouseholdInvitationUpdateRequestInput,
  ): Promise<APIResponse<HouseholdInvitation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<HouseholdInvitation>>(`/api/v1/household_invitations/${householdInvitationID}/reject`, input)
        .then((res: AxiosResponse<APIResponse<HouseholdInvitation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<HouseholdInvitation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async requestEmailVerificationEmail(): Promise<APIResponse<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<User>>(`/api/v1/users/email_address_verification`)
        .then((res: AxiosResponse<APIResponse<User>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<User>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async requestPasswordResetToken(
    input: PasswordResetTokenCreationRequestInput,
  ): Promise<APIResponse<PasswordResetToken>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<PasswordResetToken>>(`/users/password/reset`, input)
        .then((res: AxiosResponse<APIResponse<PasswordResetToken>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<PasswordResetToken>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async requestUsernameReminder(input: UsernameReminderRequestInput): Promise<APIResponse<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<User>>(`/users/username/reminder`, input)
        .then((res: AxiosResponse<APIResponse<User>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<User>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async runFinalizeMealPlanWorker(input: FinalizeMealPlansRequest): Promise<APIResponse<FinalizeMealPlansResponse>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<FinalizeMealPlansResponse>>(`/api/v1/workers/finalize_meal_plans`, input)
        .then((res: AxiosResponse<APIResponse<FinalizeMealPlansResponse>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<FinalizeMealPlansResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async runMealPlanGroceryListInitializerWorker(
    input: InitializeMealPlanGroceryListRequest,
  ): Promise<APIResponse<InitializeMealPlanGroceryListResponse>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<InitializeMealPlanGroceryListResponse>>(`/api/v1/workers/meal_plan_grocery_list_init`, input)
        .then((res: AxiosResponse<APIResponse<InitializeMealPlanGroceryListResponse>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<InitializeMealPlanGroceryListResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async runMealPlanTaskCreatorWorker(
    input: CreateMealPlanTasksRequest,
  ): Promise<APIResponse<CreateMealPlanTasksResponse>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<CreateMealPlanTasksResponse>>(`/api/v1/workers/meal_plan_tasks`, input)
        .then((res: AxiosResponse<APIResponse<CreateMealPlanTasksResponse>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<CreateMealPlanTasksResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchForMeals(q: string, filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Meal>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<Meal>>>(`/api/v1/meals/search`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<Meal>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<Meal>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<Meal>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchForRecipes(q: string, filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<Recipe>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<Recipe>>>(`/api/v1/recipes/search`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<Recipe>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<Recipe>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<Recipe>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchForServiceSettings(
    q: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ServiceSetting>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ServiceSetting>>>(`/api/v1/settings/search`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ServiceSetting>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ServiceSetting>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ServiceSetting>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchForUsers(q: string, filter: QueryFilter = QueryFilter.Default()): Promise<QueryFilteredResult<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<User>>>(`/api/v1/users/search`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<User>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<User>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<User>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchForValidIngredientGroups(
    q: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientGroup>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientGroup>>>(`/api/v1/valid_ingredient_groups/search`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientGroup>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientGroup>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientGroup>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchForValidIngredientStates(
    q: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredientState>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredientState>>>(`/api/v1/valid_ingredient_states/search`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredientState>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredientState>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredientState>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchForValidIngredients(
    q: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredient>>>(`/api/v1/valid_ingredients/search`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredient>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredient>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredient>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchForValidInstruments(
    q: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidInstrument>>>(`/api/v1/valid_instruments/search`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidInstrument>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidInstrument>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidInstrument>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchForValidMeasurementUnits(
    q: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidMeasurementUnit>>>(`/api/v1/valid_measurement_units/search`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidMeasurementUnit>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidMeasurementUnit>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidMeasurementUnit>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchForValidPreparations(
    q: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidPreparation>>>(`/api/v1/valid_preparations/search`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidPreparation>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidPreparation>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidPreparation>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchForValidVessels(
    q: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidVessel>>>(`/api/v1/valid_vessels/search`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidVessel>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidVessel>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidVessel>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchValidIngredientsByPreparation(
    q: string,
    validPreparationID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidIngredient>>>(`/api/v1/valid_ingredients/by_preparation/${validPreparationID}`, {
          params: filter.asRecord(),
        })
        .then((res: AxiosResponse<APIResponse<Array<ValidIngredient>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidIngredient>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidIngredient>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async searchValidMeasurementUnitsByIngredient(
    q: string,
    validIngredientID: string,
    filter: QueryFilter = QueryFilter.Default(),
  ): Promise<QueryFilteredResult<ValidMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .get<APIResponse<Array<ValidMeasurementUnit>>>(
          `/api/v1/valid_measurement_units/by_ingredient/${validIngredientID}`,
          {
            params: filter.asRecord(),
          },
        )
        .then((res: AxiosResponse<APIResponse<Array<ValidMeasurementUnit>>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(
              new QueryFilteredResult<ValidMeasurementUnit>({
                data: res.data.data,
                totalCount: res.data.pagination?.totalCount,
                page: res.data.pagination?.page,
                limit: res.data.pagination?.limit,
              }),
            );
          }
        })
        .catch((error: AxiosError<APIResponse<Array<ValidMeasurementUnit>>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async setDefaultHousehold(householdID: string): Promise<APIResponse<Household>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<Household>>(`/api/v1/households/${householdID}/default`)
        .then((res: AxiosResponse<APIResponse<Household>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Household>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async transferHouseholdOwnership(
    householdID: string,
    input: HouseholdOwnershipTransferInput,
  ): Promise<APIResponse<Household>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<Household>>(`/api/v1/households/${householdID}/transfer`, input)
        .then((res: AxiosResponse<APIResponse<Household>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Household>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateHousehold(householdID: string, input: HouseholdUpdateRequestInput): Promise<APIResponse<Household>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<Household>>(`/api/v1/households/${householdID}`, input)
        .then((res: AxiosResponse<APIResponse<Household>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Household>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateHouseholdInstrumentOwnership(
    householdInstrumentOwnershipID: string,
    input: HouseholdInstrumentOwnershipUpdateRequestInput,
  ): Promise<APIResponse<HouseholdInstrumentOwnership>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<HouseholdInstrumentOwnership>>(
          `/api/v1/households/instruments/${householdInstrumentOwnershipID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<HouseholdInstrumentOwnership>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<HouseholdInstrumentOwnership>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateHouseholdMemberPermissions(
    householdID: string,
    userID: string,
    input: ModifyUserPermissionsInput,
  ): Promise<APIResponse<UserPermissionsResponse>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .patch<APIResponse<UserPermissionsResponse>>(
          `/api/v1/households/${householdID}/members/${userID}/permissions`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<UserPermissionsResponse>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<UserPermissionsResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateMealPlan(mealPlanID: string, input: MealPlanUpdateRequestInput): Promise<APIResponse<MealPlan>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<MealPlan>>(`/api/v1/meal_plans/${mealPlanID}`, input)
        .then((res: AxiosResponse<APIResponse<MealPlan>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlan>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateMealPlanEvent(
    mealPlanID: string,
    mealPlanEventID: string,
    input: MealPlanEventUpdateRequestInput,
  ): Promise<APIResponse<MealPlanEvent>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<MealPlanEvent>>(`/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`, input)
        .then((res: AxiosResponse<APIResponse<MealPlanEvent>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanEvent>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateMealPlanGroceryListItem(
    mealPlanID: string,
    mealPlanGroceryListItemID: string,
    input: MealPlanGroceryListItemUpdateRequestInput,
  ): Promise<APIResponse<MealPlanGroceryListItem>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<MealPlanGroceryListItem>>(
          `/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<MealPlanGroceryListItem>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanGroceryListItem>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateMealPlanOption(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
    input: MealPlanOptionUpdateRequestInput,
  ): Promise<APIResponse<MealPlanOption>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<MealPlanOption>>(
          `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<MealPlanOption>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanOption>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateMealPlanOptionVote(
    mealPlanID: string,
    mealPlanEventID: string,
    mealPlanOptionID: string,
    mealPlanOptionVoteID: string,
    input: MealPlanOptionVoteUpdateRequestInput,
  ): Promise<APIResponse<MealPlanOptionVote>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<MealPlanOptionVote>>(
          `/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<MealPlanOptionVote>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanOptionVote>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateMealPlanTaskStatus(
    mealPlanID: string,
    mealPlanTaskID: string,
    input: MealPlanTaskStatusChangeRequestInput,
  ): Promise<APIResponse<MealPlanTask>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .patch<APIResponse<MealPlanTask>>(`/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`, input)
        .then((res: AxiosResponse<APIResponse<MealPlanTask>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<MealPlanTask>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updatePassword(input: PasswordUpdateInput): Promise<AxiosResponse<APIResponse<PasswordResetResponse>>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<PasswordResetResponse>>(`/api/v1/users/password/new`, input)
        .then((res: AxiosResponse<APIResponse<PasswordResetResponse>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res);
          }
        })
        .catch((error: AxiosError<APIResponse<PasswordResetResponse>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateRecipe(recipeID: string, input: RecipeUpdateRequestInput): Promise<APIResponse<Recipe>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<Recipe>>(`/api/v1/recipes/${recipeID}`, input)
        .then((res: AxiosResponse<APIResponse<Recipe>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<Recipe>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateRecipePrepTask(
    recipeID: string,
    recipePrepTaskID: string,
    input: RecipePrepTaskUpdateRequestInput,
  ): Promise<APIResponse<RecipePrepTask>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<RecipePrepTask>>(`/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`, input)
        .then((res: AxiosResponse<APIResponse<RecipePrepTask>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipePrepTask>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateRecipeRating(
    recipeID: string,
    recipeRatingID: string,
    input: RecipeRatingUpdateRequestInput,
  ): Promise<APIResponse<RecipeRating>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<RecipeRating>>(`/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`, input)
        .then((res: AxiosResponse<APIResponse<RecipeRating>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeRating>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateRecipeStep(
    recipeID: string,
    recipeStepID: string,
    input: RecipeStepUpdateRequestInput,
  ): Promise<APIResponse<RecipeStep>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<RecipeStep>>(`/api/v1/recipes/${recipeID}/steps/${recipeStepID}`, input)
        .then((res: AxiosResponse<APIResponse<RecipeStep>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStep>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateRecipeStepCompletionCondition(
    recipeID: string,
    recipeStepID: string,
    recipeStepCompletionConditionID: string,
    input: RecipeStepCompletionConditionUpdateRequestInput,
  ): Promise<APIResponse<RecipeStepCompletionCondition>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<RecipeStepCompletionCondition>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepCompletionCondition>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepCompletionCondition>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateRecipeStepIngredient(
    recipeID: string,
    recipeStepID: string,
    recipeStepIngredientID: string,
    input: RecipeStepIngredientUpdateRequestInput,
  ): Promise<APIResponse<RecipeStepIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<RecipeStepIngredient>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateRecipeStepInstrument(
    recipeID: string,
    recipeStepID: string,
    recipeStepInstrumentID: string,
    input: RecipeStepInstrumentUpdateRequestInput,
  ): Promise<APIResponse<RecipeStepInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<RecipeStepInstrument>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateRecipeStepProduct(
    recipeID: string,
    recipeStepID: string,
    recipeStepProductID: string,
    input: RecipeStepProductUpdateRequestInput,
  ): Promise<APIResponse<RecipeStepProduct>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<RecipeStepProduct>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepProduct>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepProduct>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateRecipeStepVessel(
    recipeID: string,
    recipeStepID: string,
    recipeStepVesselID: string,
    input: RecipeStepVesselUpdateRequestInput,
  ): Promise<APIResponse<RecipeStepVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<RecipeStepVessel>>(
          `/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<RecipeStepVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<RecipeStepVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateServiceSettingConfiguration(
    serviceSettingConfigurationID: string,
    input: ServiceSettingConfigurationUpdateRequestInput,
  ): Promise<APIResponse<ServiceSettingConfiguration>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ServiceSettingConfiguration>>(
          `/api/v1/settings/configurations/${serviceSettingConfigurationID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<ServiceSettingConfiguration>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ServiceSettingConfiguration>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateUserDetails(input: UserDetailsUpdateRequestInput): Promise<APIResponse<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<User>>(`/api/v1/users/details`, input)
        .then((res: AxiosResponse<APIResponse<User>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<User>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateUserEmailAddress(input: UserEmailAddressUpdateInput): Promise<APIResponse<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<User>>(`/api/v1/users/email_address`, input)
        .then((res: AxiosResponse<APIResponse<User>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<User>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateUserIngredientPreference(
    userIngredientPreferenceID: string,
    input: UserIngredientPreferenceUpdateRequestInput,
  ): Promise<APIResponse<UserIngredientPreference>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<UserIngredientPreference>>(
          `/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<UserIngredientPreference>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<UserIngredientPreference>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateUserNotification(
    userNotificationID: string,
    input: UserNotificationUpdateRequestInput,
  ): Promise<APIResponse<UserNotification>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .patch<APIResponse<UserNotification>>(`/api/v1/user_notifications/${userNotificationID}`, input)
        .then((res: AxiosResponse<APIResponse<UserNotification>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<UserNotification>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateUserUsername(input: UsernameUpdateInput): Promise<APIResponse<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<User>>(`/api/v1/users/username`, input)
        .then((res: AxiosResponse<APIResponse<User>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<User>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidIngredient(
    validIngredientID: string,
    input: ValidIngredientUpdateRequestInput,
  ): Promise<APIResponse<ValidIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidIngredient>>(`/api/v1/valid_ingredients/${validIngredientID}`, input)
        .then((res: AxiosResponse<APIResponse<ValidIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidIngredientGroup(
    validIngredientGroupID: string,
    input: ValidIngredientGroupUpdateRequestInput,
  ): Promise<APIResponse<ValidIngredientGroup>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidIngredientGroup>>(`/api/v1/valid_ingredient_groups/${validIngredientGroupID}`, input)
        .then((res: AxiosResponse<APIResponse<ValidIngredientGroup>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientGroup>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidIngredientMeasurementUnit(
    validIngredientMeasurementUnitID: string,
    input: ValidIngredientMeasurementUnitUpdateRequestInput,
  ): Promise<APIResponse<ValidIngredientMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidIngredientMeasurementUnit>>(
          `/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<ValidIngredientMeasurementUnit>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientMeasurementUnit>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidIngredientPreparation(
    validIngredientPreparationID: string,
    input: ValidIngredientPreparationUpdateRequestInput,
  ): Promise<APIResponse<ValidIngredientPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidIngredientPreparation>>(
          `/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<ValidIngredientPreparation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientPreparation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidIngredientState(
    validIngredientStateID: string,
    input: ValidIngredientStateUpdateRequestInput,
  ): Promise<APIResponse<ValidIngredientState>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidIngredientState>>(`/api/v1/valid_ingredient_states/${validIngredientStateID}`, input)
        .then((res: AxiosResponse<APIResponse<ValidIngredientState>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientState>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidIngredientStateIngredient(
    validIngredientStateIngredientID: string,
    input: ValidIngredientStateIngredientUpdateRequestInput,
  ): Promise<APIResponse<ValidIngredientStateIngredient>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidIngredientStateIngredient>>(
          `/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<ValidIngredientStateIngredient>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidIngredientStateIngredient>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidInstrument(
    validInstrumentID: string,
    input: ValidInstrumentUpdateRequestInput,
  ): Promise<APIResponse<ValidInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidInstrument>>(`/api/v1/valid_instruments/${validInstrumentID}`, input)
        .then((res: AxiosResponse<APIResponse<ValidInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidMeasurementUnit(
    validMeasurementUnitID: string,
    input: ValidMeasurementUnitUpdateRequestInput,
  ): Promise<APIResponse<ValidMeasurementUnit>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidMeasurementUnit>>(`/api/v1/valid_measurement_units/${validMeasurementUnitID}`, input)
        .then((res: AxiosResponse<APIResponse<ValidMeasurementUnit>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidMeasurementUnit>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidMeasurementUnitConversion(
    validMeasurementUnitConversionID: string,
    input: ValidMeasurementUnitConversionUpdateRequestInput,
  ): Promise<APIResponse<ValidMeasurementUnitConversion>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidMeasurementUnitConversion>>(
          `/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<ValidMeasurementUnitConversion>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidMeasurementUnitConversion>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidPreparation(
    validPreparationID: string,
    input: ValidPreparationUpdateRequestInput,
  ): Promise<APIResponse<ValidPreparation>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidPreparation>>(`/api/v1/valid_preparations/${validPreparationID}`, input)
        .then((res: AxiosResponse<APIResponse<ValidPreparation>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparation>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidPreparationInstrument(
    validPreparationVesselID: string,
    input: ValidPreparationInstrumentUpdateRequestInput,
  ): Promise<APIResponse<ValidPreparationInstrument>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidPreparationInstrument>>(
          `/api/v1/valid_preparation_instruments/${validPreparationVesselID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<ValidPreparationInstrument>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparationInstrument>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidPreparationVessel(
    validPreparationVesselID: string,
    input: ValidPreparationVesselUpdateRequestInput,
  ): Promise<APIResponse<ValidPreparationVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidPreparationVessel>>(
          `/api/v1/valid_preparation_vessels/${validPreparationVesselID}`,
          input,
        )
        .then((res: AxiosResponse<APIResponse<ValidPreparationVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidPreparationVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async updateValidVessel(
    validVesselID: string,
    input: ValidVesselUpdateRequestInput,
  ): Promise<APIResponse<ValidVessel>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .put<APIResponse<ValidVessel>>(`/api/v1/valid_vessels/${validVesselID}`, input)
        .then((res: AxiosResponse<APIResponse<ValidVessel>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<ValidVessel>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async uploadUserAvatar(input: AvatarUpdateInput): Promise<APIResponse<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<User>>(`/api/v1/users/avatar/upload`, input)
        .then((res: AxiosResponse<APIResponse<User>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<User>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async verifyEmailAddress(input: EmailAddressVerificationRequestInput): Promise<APIResponse<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<User>>(`/users/email_address/verify`, input)
        .then((res: AxiosResponse<APIResponse<User>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<User>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }

  async verifyTOTPSecret(input: TOTPSecretVerificationInput): Promise<APIResponse<User>> {
    let self = this;
    return new Promise(async function (resolve, reject) {
      self.client
        .post<APIResponse<User>>(`/users/totp_secret/verify`, input)
        .then((res: AxiosResponse<APIResponse<User>>) => {
          if (res.data.error) {
            reject(new Error(res.data.error.message));
          } else {
            resolve(res.data);
          }
        })
        .catch((error: AxiosError<APIResponse<User>>) => {
          if (error?.response?.data?.error) {
            reject(new Error(error?.response?.data?.error?.message || 'unknown error'));
          } else {
            reject(error);
          }
        });
    });
  }
}
