// GENERATED CODE, DO NOT EDIT MANUALLY

import axios, { AxiosResponse } from 'axios';
import AxiosMockAdapter from 'axios-mock-adapter';

import {
  IAPIError,
  ResponseDetails,
  APIResponse,
  AuditLogEntry,
  AvatarUpdateInput,
  CreateMealPlanTasksRequest,
  CreateMealPlanTasksResponse,
  DataDeletionResponse,
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
  UserDataCollectionResponse,
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

import { DinnerDoneBetterAPIClient } from './client.gen';

const mock = new AxiosMockAdapter(axios, { onNoMatch: 'throwException' });
const baseURL = 'http://things.stuff';
const fakeToken = 'test-token';
const client = new DinnerDoneBetterAPIClient(baseURL, fakeToken);

beforeEach(() => mock.reset());

type responsePartial = {
  error?: IAPIError;
  details: ResponseDetails;
};

function buildObligatoryError(msg: string): responsePartial {
  return {
    details: {
      currentHouseholdID: 'test',
      traceID: 'test',
    },
    error: {
      message: msg,
      code: 'E999',
    },
  };
}

function fakeID(): string {
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';

  for (let i = 0; i < 20; i++) {
    result += characters.charAt(Math.floor(Math.random() * characters.length));
  }

  return result;
}

describe('basic', () => {
  it('should Accepts a received household invitation', () => {
    let householdInvitationID = fakeID();

    const exampleInput = new HouseholdInvitationUpdateRequestInput();

    const exampleResponse = new APIResponse<HouseholdInvitation>();
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/accept`).reply(200, exampleResponse);

    client
      .acceptHouseholdInvitation(householdInvitationID, exampleInput)
      .then((response: APIResponse<HouseholdInvitation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Accepts a received household invitation', () => {
    let householdInvitationID = fakeID();

    const exampleInput = new HouseholdInvitationUpdateRequestInput();

    const exampleResponse = new APIResponse<HouseholdInvitation>(
      buildObligatoryError('acceptHouseholdInvitation user error'),
    );
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/accept`).reply(200, exampleResponse);

    expect(client.acceptHouseholdInvitation(householdInvitationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during Accepts a received household invitation', () => {
    let householdInvitationID = fakeID();

    const exampleInput = new HouseholdInvitationUpdateRequestInput();

    const exampleResponse = new APIResponse<HouseholdInvitation>(
      buildObligatoryError('acceptHouseholdInvitation service error'),
    );
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/accept`).reply(500, exampleResponse);

    expect(client.acceptHouseholdInvitation(householdInvitationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it("should Aggregates a user's data into a big disclosure blob", () => {
    const exampleResponse = new APIResponse<UserDataCollectionResponse>();
    mock.onPost(`${baseURL}/api/v1/data_privacy/disclose`).reply(201, exampleResponse);

    client
      .aggregateUserDataReport()
      .then((response: APIResponse<UserDataCollectionResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);

        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it("should appropriately raise errors when they occur during Aggregates a user's data into a big disclosure blob", () => {
    const exampleResponse = new APIResponse<UserDataCollectionResponse>(
      buildObligatoryError('aggregateUserDataReport user error'),
    );
    mock.onPost(`${baseURL}/api/v1/data_privacy/disclose`).reply(201, exampleResponse);

    expect(client.aggregateUserDataReport()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it("should appropriately raise service errors when they occur during Aggregates a user's data into a big disclosure blob", () => {
    const exampleResponse = new APIResponse<UserDataCollectionResponse>(
      buildObligatoryError('aggregateUserDataReport service error'),
    );
    mock.onPost(`${baseURL}/api/v1/data_privacy/disclose`).reply(500, exampleResponse);

    expect(client.aggregateUserDataReport()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Archive a household user membership', () => {
    let householdID = fakeID();
    let userID = fakeID();

    const exampleResponse = new APIResponse<HouseholdUserMembership>();
    mock.onDelete(`${baseURL}/api/v1/households/${householdID}/members/${userID}`).reply(202, exampleResponse);

    client
      .archiveUserMembership(householdID, userID)
      .then((response: APIResponse<HouseholdUserMembership>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Archive a household user membership', () => {
    let householdID = fakeID();
    let userID = fakeID();

    const exampleResponse = new APIResponse<HouseholdUserMembership>(buildObligatoryError('archiveUserMembership'));
    mock.onDelete(`${baseURL}/api/v1/households/${householdID}/members/${userID}`).reply(202, exampleResponse);

    expect(client.archiveUserMembership(householdID, userID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to Archive a household user membership', () => {
    let householdID = fakeID();
    let userID = fakeID();

    const exampleResponse = new APIResponse<HouseholdUserMembership>(buildObligatoryError('archiveUserMembership'));
    mock.onDelete(`${baseURL}/api/v1/households/${householdID}/members/${userID}`).reply(500, exampleResponse);

    expect(client.archiveUserMembership(householdID, userID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it("should Checks a user's permissions", () => {
    const exampleInput = new UserPermissionsRequestInput();

    const exampleResponse = new APIResponse<UserPermissionsResponse>();
    mock.onPost(`${baseURL}/api/v1/users/permissions/check`).reply(201, exampleResponse);

    client
      .checkPermissions(exampleInput)
      .then((response: APIResponse<UserPermissionsResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it("should appropriately raise errors when they occur during Checks a user's permissions", () => {
    const exampleInput = new UserPermissionsRequestInput();

    const exampleResponse = new APIResponse<UserPermissionsResponse>(
      buildObligatoryError('checkPermissions user error'),
    );
    mock.onPost(`${baseURL}/api/v1/users/permissions/check`).reply(201, exampleResponse);

    expect(client.checkPermissions(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it("should appropriately raise service errors when they occur during Checks a user's permissions", () => {
    const exampleInput = new UserPermissionsRequestInput();

    const exampleResponse = new APIResponse<UserPermissionsResponse>(
      buildObligatoryError('checkPermissions service error'),
    );
    mock.onPost(`${baseURL}/api/v1/users/permissions/check`).reply(500, exampleResponse);

    expect(client.checkPermissions(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Clones a recipe', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Recipe>();
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/clone`).reply(201, exampleResponse);

    client
      .cloneRecipe(recipeID)
      .then((response: APIResponse<Recipe>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);

        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Clones a recipe', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Recipe>(buildObligatoryError('cloneRecipe user error'));
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/clone`).reply(201, exampleResponse);

    expect(client.cloneRecipe(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during Clones a recipe', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Recipe>(buildObligatoryError('cloneRecipe service error'));
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/clone`).reply(500, exampleResponse);

    expect(client.cloneRecipe(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Create a household invitation', () => {
    let householdID = fakeID();

    const exampleInput = new HouseholdInvitationCreationRequestInput();

    const exampleResponse = new APIResponse<HouseholdInvitation>();
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/invite`).reply(201, exampleResponse);

    client
      .createHouseholdInvitation(householdID, exampleInput)
      .then((response: APIResponse<HouseholdInvitation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Create a household invitation', () => {
    let householdID = fakeID();

    const exampleInput = new HouseholdInvitationCreationRequestInput();

    const exampleResponse = new APIResponse<HouseholdInvitation>(
      buildObligatoryError('createHouseholdInvitation user error'),
    );
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/invite`).reply(201, exampleResponse);

    expect(client.createHouseholdInvitation(householdID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during Create a household invitation', () => {
    let householdID = fakeID();

    const exampleInput = new HouseholdInvitationCreationRequestInput();

    const exampleResponse = new APIResponse<HouseholdInvitation>(
      buildObligatoryError('createHouseholdInvitation service error'),
    );
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/invite`).reply(500, exampleResponse);

    expect(client.createHouseholdInvitation(householdID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it("should Destroys a user's data", () => {
    const exampleResponse = new APIResponse<DataDeletionResponse>();
    mock.onDelete(`${baseURL}/api/v1/data_privacy/destroy`).reply(202, exampleResponse);

    client
      .destroyAllUserData()
      .then((response: APIResponse<DataDeletionResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it("should raise errors appropriately when trying to Destroys a user's data", () => {
    const exampleResponse = new APIResponse<DataDeletionResponse>(buildObligatoryError('destroyAllUserData'));
    mock.onDelete(`${baseURL}/api/v1/data_privacy/destroy`).reply(202, exampleResponse);

    expect(client.destroyAllUserData()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it("should raise service errors appropriately when trying to Destroys a user's data", () => {
    const exampleResponse = new APIResponse<DataDeletionResponse>(buildObligatoryError('destroyAllUserData'));
    mock.onDelete(`${baseURL}/api/v1/data_privacy/destroy`).reply(500, exampleResponse);

    expect(client.destroyAllUserData()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Finalizes a meal plan', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<FinalizeMealPlansResponse>();
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/finalize`).reply(201, exampleResponse);

    client
      .finalizeMealPlan(mealPlanID)
      .then((response: APIResponse<FinalizeMealPlansResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);

        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Finalizes a meal plan', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<FinalizeMealPlansResponse>(
      buildObligatoryError('finalizeMealPlan user error'),
    );
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/finalize`).reply(201, exampleResponse);

    expect(client.finalizeMealPlan(mealPlanID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during Finalizes a meal plan', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<FinalizeMealPlansResponse>(
      buildObligatoryError('finalizeMealPlan service error'),
    );
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/finalize`).reply(500, exampleResponse);

    expect(client.finalizeMealPlan(mealPlanID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Operation for creating JWTResponse', () => {
    const exampleInput = new UserLoginInput();

    const exampleResponse = new APIResponse<JWTResponse>();
    mock.onPost(`${baseURL}/users/login/jwt/admin`).reply(201, exampleResponse);

    client
      .adminLoginForJWT(exampleInput)
      .then((response: AxiosResponse<APIResponse<JWTResponse>>) => {
        expect(response.data).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Operation for creating JWTResponse', () => {
    const exampleInput = new UserLoginInput();

    const exampleResponse = new APIResponse<JWTResponse>(buildObligatoryError('adminLoginForJWT user error'));
    mock.onPost(`${baseURL}/users/login/jwt/admin`).reply(201, exampleResponse);

    expect(client.adminLoginForJWT(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during Operation for creating JWTResponse', () => {
    const exampleInput = new UserLoginInput();

    const exampleResponse = new APIResponse<JWTResponse>(buildObligatoryError('adminLoginForJWT service error'));
    mock.onPost(`${baseURL}/users/login/jwt/admin`).reply(500, exampleResponse);

    expect(client.adminLoginForJWT(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Operation for creating JWTResponse', () => {
    const exampleInput = new UserLoginInput();

    const exampleResponse = new APIResponse<JWTResponse>();
    mock.onPost(`${baseURL}/users/login/jwt`).reply(201, exampleResponse);

    client
      .loginForJWT(exampleInput)
      .then((response: AxiosResponse<APIResponse<JWTResponse>>) => {
        expect(response.data).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Operation for creating JWTResponse', () => {
    const exampleInput = new UserLoginInput();

    const exampleResponse = new APIResponse<JWTResponse>(buildObligatoryError('loginForJWT user error'));
    mock.onPost(`${baseURL}/users/login/jwt`).reply(201, exampleResponse);

    expect(client.loginForJWT(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during Operation for creating JWTResponse', () => {
    const exampleInput = new UserLoginInput();

    const exampleResponse = new APIResponse<JWTResponse>(buildObligatoryError('loginForJWT service error'));
    mock.onPost(`${baseURL}/users/login/jwt`).reply(500, exampleResponse);

    expect(client.loginForJWT(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Operation for creating PasswordResetToken', () => {
    const exampleInput = new PasswordResetTokenCreationRequestInput();

    const exampleResponse = new APIResponse<PasswordResetToken>();
    mock.onPost(`${baseURL}/users/password/reset`).reply(201, exampleResponse);

    client
      .requestPasswordResetToken(exampleInput)
      .then((response: APIResponse<PasswordResetToken>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Operation for creating PasswordResetToken', () => {
    const exampleInput = new PasswordResetTokenCreationRequestInput();

    const exampleResponse = new APIResponse<PasswordResetToken>(
      buildObligatoryError('requestPasswordResetToken user error'),
    );
    mock.onPost(`${baseURL}/users/password/reset`).reply(201, exampleResponse);

    expect(client.requestPasswordResetToken(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during Operation for creating PasswordResetToken', () => {
    const exampleInput = new PasswordResetTokenCreationRequestInput();

    const exampleResponse = new APIResponse<PasswordResetToken>(
      buildObligatoryError('requestPasswordResetToken service error'),
    );
    mock.onPost(`${baseURL}/users/password/reset`).reply(500, exampleResponse);

    expect(client.requestPasswordResetToken(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Operation for creating User', () => {
    const exampleInput = new EmailAddressVerificationRequestInput();

    const exampleResponse = new APIResponse<User>();
    mock.onPost(`${baseURL}/users/email_address/verify`).reply(201, exampleResponse);

    client
      .verifyEmailAddress(exampleInput)
      .then((response: APIResponse<User>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Operation for creating User', () => {
    const exampleInput = new EmailAddressVerificationRequestInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('verifyEmailAddress user error'));
    mock.onPost(`${baseURL}/users/email_address/verify`).reply(201, exampleResponse);

    expect(client.verifyEmailAddress(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during Operation for creating User', () => {
    const exampleInput = new EmailAddressVerificationRequestInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('verifyEmailAddress service error'));
    mock.onPost(`${baseURL}/users/email_address/verify`).reply(500, exampleResponse);

    expect(client.verifyEmailAddress(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Operation for creating User', () => {
    const exampleInput = new TOTPSecretVerificationInput();

    const exampleResponse = new APIResponse<User>();
    mock.onPost(`${baseURL}/users/totp_secret/verify`).reply(201, exampleResponse);

    client
      .verifyTOTPSecret(exampleInput)
      .then((response: APIResponse<User>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Operation for creating User', () => {
    const exampleInput = new TOTPSecretVerificationInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('verifyTOTPSecret user error'));
    mock.onPost(`${baseURL}/users/totp_secret/verify`).reply(201, exampleResponse);

    expect(client.verifyTOTPSecret(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during Operation for creating User', () => {
    const exampleInput = new TOTPSecretVerificationInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('verifyTOTPSecret service error'));
    mock.onPost(`${baseURL}/users/totp_secret/verify`).reply(500, exampleResponse);

    expect(client.verifyTOTPSecret(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Operation for creating User', () => {
    const exampleInput = new UsernameReminderRequestInput();

    const exampleResponse = new APIResponse<User>();
    mock.onPost(`${baseURL}/users/username/reminder`).reply(201, exampleResponse);

    client
      .requestUsernameReminder(exampleInput)
      .then((response: APIResponse<User>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Operation for creating User', () => {
    const exampleInput = new UsernameReminderRequestInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('requestUsernameReminder user error'));
    mock.onPost(`${baseURL}/users/username/reminder`).reply(201, exampleResponse);

    expect(client.requestUsernameReminder(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during Operation for creating User', () => {
    const exampleInput = new UsernameReminderRequestInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('requestUsernameReminder service error'));
    mock.onPost(`${baseURL}/users/username/reminder`).reply(500, exampleResponse);

    expect(client.requestUsernameReminder(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Redeems a password reset token', () => {
    const exampleInput = new PasswordResetTokenRedemptionRequestInput();

    const exampleResponse = new APIResponse<User>();
    mock.onPost(`${baseURL}/users/password/reset/redeem`).reply(201, exampleResponse);

    client
      .redeemPasswordResetToken(exampleInput)
      .then((response: APIResponse<User>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Redeems a password reset token', () => {
    const exampleInput = new PasswordResetTokenRedemptionRequestInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('redeemPasswordResetToken user error'));
    mock.onPost(`${baseURL}/users/password/reset/redeem`).reply(201, exampleResponse);

    expect(client.redeemPasswordResetToken(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during Redeems a password reset token', () => {
    const exampleInput = new PasswordResetTokenRedemptionRequestInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('redeemPasswordResetToken service error'));
    mock.onPost(`${baseURL}/users/password/reset/redeem`).reply(500, exampleResponse);

    expect(client.redeemPasswordResetToken(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it("should Refreshes a user's TOTP secret", () => {
    const exampleInput = new TOTPSecretRefreshInput();

    const exampleResponse = new APIResponse<TOTPSecretRefreshResponse>();
    mock.onPost(`${baseURL}/api/v1/users/totp_secret/new`).reply(201, exampleResponse);

    client
      .refreshTOTPSecret(exampleInput)
      .then((response: APIResponse<TOTPSecretRefreshResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it("should appropriately raise errors when they occur during Refreshes a user's TOTP secret", () => {
    const exampleInput = new TOTPSecretRefreshInput();

    const exampleResponse = new APIResponse<TOTPSecretRefreshResponse>(
      buildObligatoryError('refreshTOTPSecret user error'),
    );
    mock.onPost(`${baseURL}/api/v1/users/totp_secret/new`).reply(201, exampleResponse);

    expect(client.refreshTOTPSecret(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it("should appropriately raise service errors when they occur during Refreshes a user's TOTP secret", () => {
    const exampleInput = new TOTPSecretRefreshInput();

    const exampleResponse = new APIResponse<TOTPSecretRefreshResponse>(
      buildObligatoryError('refreshTOTPSecret service error'),
    );
    mock.onPost(`${baseURL}/api/v1/users/totp_secret/new`).reply(500, exampleResponse);

    expect(client.refreshTOTPSecret(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Requests an email address verification email', () => {
    const exampleResponse = new APIResponse<User>();
    mock.onPost(`${baseURL}/api/v1/users/email_address_verification`).reply(201, exampleResponse);

    client
      .requestEmailVerificationEmail()
      .then((response: APIResponse<User>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);

        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Requests an email address verification email', () => {
    const exampleResponse = new APIResponse<User>(buildObligatoryError('requestEmailVerificationEmail user error'));
    mock.onPost(`${baseURL}/api/v1/users/email_address_verification`).reply(201, exampleResponse);

    expect(client.requestEmailVerificationEmail()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during Requests an email address verification email', () => {
    const exampleResponse = new APIResponse<User>(buildObligatoryError('requestEmailVerificationEmail service error'));
    mock.onPost(`${baseURL}/api/v1/users/email_address_verification`).reply(500, exampleResponse);

    expect(client.requestEmailVerificationEmail()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Retrieves an audit log entry by ID', () => {
    let auditLogEntryID = fakeID();

    const exampleResponse = new APIResponse<AuditLogEntry>();
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/${auditLogEntryID}`).reply(200, exampleResponse);

    client
      .getAuditLogEntryByID(auditLogEntryID)
      .then((response: APIResponse<AuditLogEntry>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Retrieves an audit log entry by ID', () => {
    let auditLogEntryID = fakeID();

    const exampleResponse = new APIResponse<AuditLogEntry>(buildObligatoryError('getAuditLogEntryByID'));
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/${auditLogEntryID}`).reply(200, exampleResponse);

    expect(client.getAuditLogEntryByID(auditLogEntryID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Retrieves an audit log entry by ID', () => {
    let auditLogEntryID = fakeID();

    const exampleResponse = new APIResponse<AuditLogEntry>(buildObligatoryError('getAuditLogEntryByID'));
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/${auditLogEntryID}`).reply(500, exampleResponse);

    expect(client.getAuditLogEntryByID(auditLogEntryID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Retrieves audit log entries for a household', () => {
    const exampleResponse = new APIResponse<Array<AuditLogEntry>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new AuditLogEntry()],
    });
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/for_household`).reply(200, exampleResponse);

    client
      .getAuditLogEntriesForHousehold()
      .then((response: QueryFilteredResult<AuditLogEntry>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Retrieves audit log entries for a household', () => {
    const exampleResponse = new APIResponse<Array<AuditLogEntry>>(
      buildObligatoryError('getAuditLogEntriesForHousehold user error'),
    );
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/for_household`).reply(200, exampleResponse);

    expect(client.getAuditLogEntriesForHousehold()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Retrieves audit log entries for a household', () => {
    const exampleResponse = new APIResponse<Array<AuditLogEntry>>(
      buildObligatoryError('getAuditLogEntriesForHousehold service error'),
    );
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/for_household`).reply(500, exampleResponse);

    expect(client.getAuditLogEntriesForHousehold()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Retrieves audit log entries for a user', () => {
    const exampleResponse = new APIResponse<Array<AuditLogEntry>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new AuditLogEntry()],
    });
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/for_user`).reply(200, exampleResponse);

    client
      .getAuditLogEntriesForUser()
      .then((response: QueryFilteredResult<AuditLogEntry>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Retrieves audit log entries for a user', () => {
    const exampleResponse = new APIResponse<Array<AuditLogEntry>>(
      buildObligatoryError('getAuditLogEntriesForUser user error'),
    );
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/for_user`).reply(200, exampleResponse);

    expect(client.getAuditLogEntriesForUser()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Retrieves audit log entries for a user', () => {
    const exampleResponse = new APIResponse<Array<AuditLogEntry>>(
      buildObligatoryError('getAuditLogEntriesForUser service error'),
    );
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/for_user`).reply(500, exampleResponse);

    expect(client.getAuditLogEntriesForUser()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Runs the Finalize Meal Plans worker', () => {
    const exampleInput = new FinalizeMealPlansRequest();

    const exampleResponse = new APIResponse<FinalizeMealPlansResponse>();
    mock.onPost(`${baseURL}/api/v1/workers/finalize_meal_plans`).reply(201, exampleResponse);

    client
      .runFinalizeMealPlanWorker(exampleInput)
      .then((response: APIResponse<FinalizeMealPlansResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Runs the Finalize Meal Plans worker', () => {
    const exampleInput = new FinalizeMealPlansRequest();

    const exampleResponse = new APIResponse<FinalizeMealPlansResponse>(
      buildObligatoryError('runFinalizeMealPlanWorker user error'),
    );
    mock.onPost(`${baseURL}/api/v1/workers/finalize_meal_plans`).reply(201, exampleResponse);

    expect(client.runFinalizeMealPlanWorker(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during Runs the Finalize Meal Plans worker', () => {
    const exampleInput = new FinalizeMealPlansRequest();

    const exampleResponse = new APIResponse<FinalizeMealPlansResponse>(
      buildObligatoryError('runFinalizeMealPlanWorker service error'),
    );
    mock.onPost(`${baseURL}/api/v1/workers/finalize_meal_plans`).reply(500, exampleResponse);

    expect(client.runFinalizeMealPlanWorker(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Runs the meal plan grocery list initialization worker', () => {
    const exampleInput = new InitializeMealPlanGroceryListRequest();

    const exampleResponse = new APIResponse<InitializeMealPlanGroceryListResponse>();
    mock.onPost(`${baseURL}/api/v1/workers/meal_plan_grocery_list_init`).reply(201, exampleResponse);

    client
      .runMealPlanGroceryListInitializerWorker(exampleInput)
      .then((response: APIResponse<InitializeMealPlanGroceryListResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Runs the meal plan grocery list initialization worker', () => {
    const exampleInput = new InitializeMealPlanGroceryListRequest();

    const exampleResponse = new APIResponse<InitializeMealPlanGroceryListResponse>(
      buildObligatoryError('runMealPlanGroceryListInitializerWorker user error'),
    );
    mock.onPost(`${baseURL}/api/v1/workers/meal_plan_grocery_list_init`).reply(201, exampleResponse);

    expect(client.runMealPlanGroceryListInitializerWorker(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during Runs the meal plan grocery list initialization worker', () => {
    const exampleInput = new InitializeMealPlanGroceryListRequest();

    const exampleResponse = new APIResponse<InitializeMealPlanGroceryListResponse>(
      buildObligatoryError('runMealPlanGroceryListInitializerWorker service error'),
    );
    mock.onPost(`${baseURL}/api/v1/workers/meal_plan_grocery_list_init`).reply(500, exampleResponse);

    expect(client.runMealPlanGroceryListInitializerWorker(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should Runs the meal plan task creator worker', () => {
    const exampleInput = new CreateMealPlanTasksRequest();

    const exampleResponse = new APIResponse<CreateMealPlanTasksResponse>();
    mock.onPost(`${baseURL}/api/v1/workers/meal_plan_tasks`).reply(201, exampleResponse);

    client
      .runMealPlanTaskCreatorWorker(exampleInput)
      .then((response: APIResponse<CreateMealPlanTasksResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Runs the meal plan task creator worker', () => {
    const exampleInput = new CreateMealPlanTasksRequest();

    const exampleResponse = new APIResponse<CreateMealPlanTasksResponse>(
      buildObligatoryError('runMealPlanTaskCreatorWorker user error'),
    );
    mock.onPost(`${baseURL}/api/v1/workers/meal_plan_tasks`).reply(201, exampleResponse);

    expect(client.runMealPlanTaskCreatorWorker(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during Runs the meal plan task creator worker', () => {
    const exampleInput = new CreateMealPlanTasksRequest();

    const exampleResponse = new APIResponse<CreateMealPlanTasksResponse>(
      buildObligatoryError('runMealPlanTaskCreatorWorker service error'),
    );
    mock.onPost(`${baseURL}/api/v1/workers/meal_plan_tasks`).reply(500, exampleResponse);

    expect(client.runMealPlanTaskCreatorWorker(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should Searches for Meals', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<Meal>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new Meal()],
    });
    mock.onGet(`${baseURL}/api/v1/meals/search`).reply(200, exampleResponse);

    client
      .searchForMeals(q)
      .then((response: QueryFilteredResult<Meal>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for Meals', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<Meal>>(buildObligatoryError('searchForMeals user error'));
    mock.onGet(`${baseURL}/api/v1/meals/search`).reply(200, exampleResponse);

    expect(client.searchForMeals(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Searches for Meals', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<Meal>>(buildObligatoryError('searchForMeals service error'));
    mock.onGet(`${baseURL}/api/v1/meals/search`).reply(500, exampleResponse);

    expect(client.searchForMeals(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Searches for Recipes', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<Recipe>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new Recipe()],
    });
    mock.onGet(`${baseURL}/api/v1/recipes/search`).reply(200, exampleResponse);

    client
      .searchForRecipes(q)
      .then((response: QueryFilteredResult<Recipe>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for Recipes', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<Recipe>>(buildObligatoryError('searchForRecipes user error'));
    mock.onGet(`${baseURL}/api/v1/recipes/search`).reply(200, exampleResponse);

    expect(client.searchForRecipes(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Searches for Recipes', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<Recipe>>(buildObligatoryError('searchForRecipes service error'));
    mock.onGet(`${baseURL}/api/v1/recipes/search`).reply(500, exampleResponse);

    expect(client.searchForRecipes(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Searches for ServiceSettings', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ServiceSetting>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ServiceSetting()],
    });
    mock.onGet(`${baseURL}/api/v1/settings/search`).reply(200, exampleResponse);

    client
      .searchForServiceSettings(q)
      .then((response: QueryFilteredResult<ServiceSetting>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for ServiceSettings', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ServiceSetting>>(
      buildObligatoryError('searchForServiceSettings user error'),
    );
    mock.onGet(`${baseURL}/api/v1/settings/search`).reply(200, exampleResponse);

    expect(client.searchForServiceSettings(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Searches for ServiceSettings', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ServiceSetting>>(
      buildObligatoryError('searchForServiceSettings service error'),
    );
    mock.onGet(`${baseURL}/api/v1/settings/search`).reply(500, exampleResponse);

    expect(client.searchForServiceSettings(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Searches for Users', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<User>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new User()],
    });
    mock.onGet(`${baseURL}/api/v1/users/search`).reply(200, exampleResponse);

    client
      .searchForUsers(q)
      .then((response: QueryFilteredResult<User>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for Users', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<User>>(buildObligatoryError('searchForUsers user error'));
    mock.onGet(`${baseURL}/api/v1/users/search`).reply(200, exampleResponse);

    expect(client.searchForUsers(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Searches for Users', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<User>>(buildObligatoryError('searchForUsers service error'));
    mock.onGet(`${baseURL}/api/v1/users/search`).reply(500, exampleResponse);

    expect(client.searchForUsers(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Searches for ValidIngredientGroups', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientGroup>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientGroup()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/search`).reply(200, exampleResponse);

    client
      .searchForValidIngredientGroups(q)
      .then((response: QueryFilteredResult<ValidIngredientGroup>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for ValidIngredientGroups', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientGroup>>(
      buildObligatoryError('searchForValidIngredientGroups user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/search`).reply(200, exampleResponse);

    expect(client.searchForValidIngredientGroups(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Searches for ValidIngredientGroups', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientGroup>>(
      buildObligatoryError('searchForValidIngredientGroups service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/search`).reply(500, exampleResponse);

    expect(client.searchForValidIngredientGroups(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Searches for ValidIngredientStates', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientState>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientState()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/search`).reply(200, exampleResponse);

    client
      .searchForValidIngredientStates(q)
      .then((response: QueryFilteredResult<ValidIngredientState>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for ValidIngredientStates', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientState>>(
      buildObligatoryError('searchForValidIngredientStates user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/search`).reply(200, exampleResponse);

    expect(client.searchForValidIngredientStates(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Searches for ValidIngredientStates', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientState>>(
      buildObligatoryError('searchForValidIngredientStates service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/search`).reply(500, exampleResponse);

    expect(client.searchForValidIngredientStates(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Searches for ValidIngredients by a ValidPreparation ID', () => {
    let q = fakeID();
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredient>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredient()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/by_preparation/${validPreparationID}`).reply(200, exampleResponse);

    client
      .searchValidIngredientsByPreparation(q, validPreparationID)
      .then((response: QueryFilteredResult<ValidIngredient>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for ValidIngredients by a ValidPreparation ID', () => {
    let q = fakeID();
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredient>>(
      buildObligatoryError('searchValidIngredientsByPreparation user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/by_preparation/${validPreparationID}`).reply(200, exampleResponse);

    expect(client.searchValidIngredientsByPreparation(q, validPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to Searches for ValidIngredients by a ValidPreparation ID', () => {
    let q = fakeID();
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredient>>(
      buildObligatoryError('searchValidIngredientsByPreparation service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/by_preparation/${validPreparationID}`).reply(500, exampleResponse);

    expect(client.searchValidIngredientsByPreparation(q, validPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should Searches for ValidIngredients', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredient>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredient()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/search`).reply(200, exampleResponse);

    client
      .searchForValidIngredients(q)
      .then((response: QueryFilteredResult<ValidIngredient>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for ValidIngredients', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredient>>(
      buildObligatoryError('searchForValidIngredients user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/search`).reply(200, exampleResponse);

    expect(client.searchForValidIngredients(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Searches for ValidIngredients', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredient>>(
      buildObligatoryError('searchForValidIngredients service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/search`).reply(500, exampleResponse);

    expect(client.searchForValidIngredients(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Searches for ValidInstruments', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidInstrument>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidInstrument()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_instruments/search`).reply(200, exampleResponse);

    client
      .searchForValidInstruments(q)
      .then((response: QueryFilteredResult<ValidInstrument>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for ValidInstruments', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidInstrument>>(
      buildObligatoryError('searchForValidInstruments user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_instruments/search`).reply(200, exampleResponse);

    expect(client.searchForValidInstruments(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Searches for ValidInstruments', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidInstrument>>(
      buildObligatoryError('searchForValidInstruments service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_instruments/search`).reply(500, exampleResponse);

    expect(client.searchForValidInstruments(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Searches for ValidMeasurementUnits by a ValidIngredient ID', () => {
    let q = fakeID();
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidMeasurementUnit()],
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_units/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    client
      .searchValidMeasurementUnitsByIngredient(q, validIngredientID)
      .then((response: QueryFilteredResult<ValidMeasurementUnit>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for ValidMeasurementUnits by a ValidIngredient ID', () => {
    let q = fakeID();
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>(
      buildObligatoryError('searchValidMeasurementUnitsByIngredient user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_units/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.searchValidMeasurementUnitsByIngredient(q, validIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to Searches for ValidMeasurementUnits by a ValidIngredient ID', () => {
    let q = fakeID();
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>(
      buildObligatoryError('searchValidMeasurementUnitsByIngredient service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_units/by_ingredient/${validIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.searchValidMeasurementUnitsByIngredient(q, validIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should Searches for ValidMeasurementUnits', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidMeasurementUnit()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/search`).reply(200, exampleResponse);

    client
      .searchForValidMeasurementUnits(q)
      .then((response: QueryFilteredResult<ValidMeasurementUnit>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for ValidMeasurementUnits', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>(
      buildObligatoryError('searchForValidMeasurementUnits user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/search`).reply(200, exampleResponse);

    expect(client.searchForValidMeasurementUnits(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Searches for ValidMeasurementUnits', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>(
      buildObligatoryError('searchForValidMeasurementUnits service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/search`).reply(500, exampleResponse);

    expect(client.searchForValidMeasurementUnits(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Searches for ValidPreparations', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparation>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidPreparation()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparations/search`).reply(200, exampleResponse);

    client
      .searchForValidPreparations(q)
      .then((response: QueryFilteredResult<ValidPreparation>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for ValidPreparations', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparation>>(
      buildObligatoryError('searchForValidPreparations user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_preparations/search`).reply(200, exampleResponse);

    expect(client.searchForValidPreparations(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Searches for ValidPreparations', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparation>>(
      buildObligatoryError('searchForValidPreparations service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_preparations/search`).reply(500, exampleResponse);

    expect(client.searchForValidPreparations(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Searches for ValidVessels', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidVessel>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidVessel()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_vessels/search`).reply(200, exampleResponse);

    client
      .searchForValidVessels(q)
      .then((response: QueryFilteredResult<ValidVessel>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to Searches for ValidVessels', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidVessel>>(
      buildObligatoryError('searchForValidVessels user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_vessels/search`).reply(200, exampleResponse);

    expect(client.searchForValidVessels(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to Searches for ValidVessels', () => {
    let q = fakeID();

    const exampleResponse = new APIResponse<Array<ValidVessel>>(
      buildObligatoryError('searchForValidVessels service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_vessels/search`).reply(500, exampleResponse);

    expect(client.searchForValidVessels(q)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should Transfer household ownership to another user', () => {
    let householdID = fakeID();

    const exampleInput = new HouseholdOwnershipTransferInput();

    const exampleResponse = new APIResponse<Household>();
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/transfer`).reply(201, exampleResponse);

    client
      .transferHouseholdOwnership(householdID, exampleInput)
      .then((response: APIResponse<Household>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Transfer household ownership to another user', () => {
    let householdID = fakeID();

    const exampleInput = new HouseholdOwnershipTransferInput();

    const exampleResponse = new APIResponse<Household>(buildObligatoryError('transferHouseholdOwnership user error'));
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/transfer`).reply(201, exampleResponse);

    expect(client.transferHouseholdOwnership(householdID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during Transfer household ownership to another user', () => {
    let householdID = fakeID();

    const exampleInput = new HouseholdOwnershipTransferInput();

    const exampleResponse = new APIResponse<Household>(
      buildObligatoryError('transferHouseholdOwnership service error'),
    );
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/transfer`).reply(500, exampleResponse);

    expect(client.transferHouseholdOwnership(householdID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it("should Update a household member's household permissions", () => {
    let householdID = fakeID();
    let userID = fakeID();

    const exampleInput = new ModifyUserPermissionsInput();

    const exampleResponse = new APIResponse<UserPermissionsResponse>();
    mock
      .onPatch(`${baseURL}/api/v1/households/${householdID}/members/${userID}/permissions`)
      .reply(200, exampleResponse);

    client
      .updateHouseholdMemberPermissions(householdID, userID, exampleInput)
      .then((response: APIResponse<UserPermissionsResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.patch.length).toBe(1);
        expect(mock.history.patch[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.patch[0].headers).toHaveProperty('Authorization');
        expect((mock.history.patch[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it("should appropriately raise errors when they occur during Update a household member's household permissions", () => {
    let householdID = fakeID();
    let userID = fakeID();

    const exampleInput = new ModifyUserPermissionsInput();

    const exampleResponse = new APIResponse<UserPermissionsResponse>(
      buildObligatoryError('updateHouseholdMemberPermissions user error'),
    );
    mock
      .onPatch(`${baseURL}/api/v1/households/${householdID}/members/${userID}/permissions`)
      .reply(200, exampleResponse);

    expect(client.updateHouseholdMemberPermissions(householdID, userID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it("should appropriately raise service errors when they occur during Update a household member's household permissions", () => {
    let householdID = fakeID();
    let userID = fakeID();

    const exampleInput = new ModifyUserPermissionsInput();

    const exampleResponse = new APIResponse<UserPermissionsResponse>(
      buildObligatoryError('updateHouseholdMemberPermissions service error'),
    );
    mock
      .onPatch(`${baseURL}/api/v1/households/${householdID}/members/${userID}/permissions`)
      .reply(500, exampleResponse);

    expect(client.updateHouseholdMemberPermissions(householdID, userID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should Uploads a new user avatar', () => {
    const exampleInput = new AvatarUpdateInput();

    const exampleResponse = new APIResponse<User>();
    mock.onPost(`${baseURL}/api/v1/users/avatar/upload`).reply(201, exampleResponse);

    client
      .uploadUserAvatar(exampleInput)
      .then((response: APIResponse<User>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during Uploads a new user avatar', () => {
    const exampleInput = new AvatarUpdateInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('uploadUserAvatar user error'));
    mock.onPost(`${baseURL}/api/v1/users/avatar/upload`).reply(201, exampleResponse);

    expect(client.uploadUserAvatar(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during Uploads a new user avatar', () => {
    const exampleInput = new AvatarUpdateInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('uploadUserAvatar service error'));
    mock.onPost(`${baseURL}/api/v1/users/avatar/upload`).reply(500, exampleResponse);

    expect(client.uploadUserAvatar(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a Meal', () => {
    let mealID = fakeID();

    const exampleResponse = new APIResponse<Meal>();
    mock.onDelete(`${baseURL}/api/v1/meals/${mealID}`).reply(202, exampleResponse);

    client
      .archiveMeal(mealID)
      .then((response: APIResponse<Meal>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a Meal', () => {
    let mealID = fakeID();

    const exampleResponse = new APIResponse<Meal>(buildObligatoryError('archiveMeal'));
    mock.onDelete(`${baseURL}/api/v1/meals/${mealID}`).reply(202, exampleResponse);

    expect(client.archiveMeal(mealID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a Meal', () => {
    let mealID = fakeID();

    const exampleResponse = new APIResponse<Meal>(buildObligatoryError('archiveMeal'));
    mock.onDelete(`${baseURL}/api/v1/meals/${mealID}`).reply(500, exampleResponse);

    expect(client.archiveMeal(mealID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a MealPlan', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<MealPlan>();
    mock.onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(202, exampleResponse);

    client
      .archiveMealPlan(mealPlanID)
      .then((response: APIResponse<MealPlan>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a MealPlan', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<MealPlan>(buildObligatoryError('archiveMealPlan'));
    mock.onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(202, exampleResponse);

    expect(client.archiveMealPlan(mealPlanID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a MealPlan', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<MealPlan>(buildObligatoryError('archiveMealPlan'));
    mock.onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(500, exampleResponse);

    expect(client.archiveMealPlan(mealPlanID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a MealPlanEvent', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleResponse = new APIResponse<MealPlanEvent>();
    mock.onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(202, exampleResponse);

    client
      .archiveMealPlanEvent(mealPlanID, mealPlanEventID)
      .then((response: APIResponse<MealPlanEvent>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a MealPlanEvent', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleResponse = new APIResponse<MealPlanEvent>(buildObligatoryError('archiveMealPlanEvent'));
    mock.onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(202, exampleResponse);

    expect(client.archiveMealPlanEvent(mealPlanID, mealPlanEventID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a MealPlanEvent', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleResponse = new APIResponse<MealPlanEvent>(buildObligatoryError('archiveMealPlanEvent'));
    mock.onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(500, exampleResponse);

    expect(client.archiveMealPlanEvent(mealPlanID, mealPlanEventID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();
    let mealPlanGroceryListItemID = fakeID();

    const exampleResponse = new APIResponse<MealPlanGroceryListItem>();
    mock
      .onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(202, exampleResponse);

    client
      .archiveMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)
      .then((response: APIResponse<MealPlanGroceryListItem>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();
    let mealPlanGroceryListItemID = fakeID();

    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(
      buildObligatoryError('archiveMealPlanGroceryListItem'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(202, exampleResponse);

    expect(client.archiveMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();
    let mealPlanGroceryListItemID = fakeID();

    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(
      buildObligatoryError('archiveMealPlanGroceryListItem'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(500, exampleResponse);

    expect(client.archiveMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a MealPlanOption', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleResponse = new APIResponse<MealPlanOption>();
    mock
      .onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(202, exampleResponse);

    client
      .archiveMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)
      .then((response: APIResponse<MealPlanOption>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a MealPlanOption', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleResponse = new APIResponse<MealPlanOption>(buildObligatoryError('archiveMealPlanOption'));
    mock
      .onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(202, exampleResponse);

    expect(client.archiveMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a MealPlanOption', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleResponse = new APIResponse<MealPlanOption>(buildObligatoryError('archiveMealPlanOption'));
    mock
      .onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(500, exampleResponse);

    expect(client.archiveMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a OAuth2Client', () => {
    let oauth2ClientID = fakeID();

    const exampleResponse = new APIResponse<OAuth2Client>();
    mock.onDelete(`${baseURL}/api/v1/oauth2_clients/${oauth2ClientID}`).reply(202, exampleResponse);

    client
      .archiveOAuth2Client(oauth2ClientID)
      .then((response: APIResponse<OAuth2Client>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a OAuth2Client', () => {
    let oauth2ClientID = fakeID();

    const exampleResponse = new APIResponse<OAuth2Client>(buildObligatoryError('archiveOAuth2Client'));
    mock.onDelete(`${baseURL}/api/v1/oauth2_clients/${oauth2ClientID}`).reply(202, exampleResponse);

    expect(client.archiveOAuth2Client(oauth2ClientID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a OAuth2Client', () => {
    let oauth2ClientID = fakeID();

    const exampleResponse = new APIResponse<OAuth2Client>(buildObligatoryError('archiveOAuth2Client'));
    mock.onDelete(`${baseURL}/api/v1/oauth2_clients/${oauth2ClientID}`).reply(500, exampleResponse);

    expect(client.archiveOAuth2Client(oauth2ClientID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a Recipe', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Recipe>();
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}`).reply(202, exampleResponse);

    client
      .archiveRecipe(recipeID)
      .then((response: APIResponse<Recipe>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a Recipe', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Recipe>(buildObligatoryError('archiveRecipe'));
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}`).reply(202, exampleResponse);

    expect(client.archiveRecipe(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a Recipe', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Recipe>(buildObligatoryError('archiveRecipe'));
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}`).reply(500, exampleResponse);

    expect(client.archiveRecipe(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a RecipePrepTask', () => {
    let recipeID = fakeID();
    let recipePrepTaskID = fakeID();

    const exampleResponse = new APIResponse<RecipePrepTask>();
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(202, exampleResponse);

    client
      .archiveRecipePrepTask(recipeID, recipePrepTaskID)
      .then((response: APIResponse<RecipePrepTask>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a RecipePrepTask', () => {
    let recipeID = fakeID();
    let recipePrepTaskID = fakeID();

    const exampleResponse = new APIResponse<RecipePrepTask>(buildObligatoryError('archiveRecipePrepTask'));
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(202, exampleResponse);

    expect(client.archiveRecipePrepTask(recipeID, recipePrepTaskID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a RecipePrepTask', () => {
    let recipeID = fakeID();
    let recipePrepTaskID = fakeID();

    const exampleResponse = new APIResponse<RecipePrepTask>(buildObligatoryError('archiveRecipePrepTask'));
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(500, exampleResponse);

    expect(client.archiveRecipePrepTask(recipeID, recipePrepTaskID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a RecipeRating', () => {
    let recipeID = fakeID();
    let recipeRatingID = fakeID();

    const exampleResponse = new APIResponse<RecipeRating>();
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(202, exampleResponse);

    client
      .archiveRecipeRating(recipeID, recipeRatingID)
      .then((response: APIResponse<RecipeRating>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a RecipeRating', () => {
    let recipeID = fakeID();
    let recipeRatingID = fakeID();

    const exampleResponse = new APIResponse<RecipeRating>(buildObligatoryError('archiveRecipeRating'));
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(202, exampleResponse);

    expect(client.archiveRecipeRating(recipeID, recipeRatingID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a RecipeRating', () => {
    let recipeID = fakeID();
    let recipeRatingID = fakeID();

    const exampleResponse = new APIResponse<RecipeRating>(buildObligatoryError('archiveRecipeRating'));
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(500, exampleResponse);

    expect(client.archiveRecipeRating(recipeID, recipeRatingID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a RecipeStep', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<RecipeStep>();
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(202, exampleResponse);

    client
      .archiveRecipeStep(recipeID, recipeStepID)
      .then((response: APIResponse<RecipeStep>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a RecipeStep', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<RecipeStep>(buildObligatoryError('archiveRecipeStep'));
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(202, exampleResponse);

    expect(client.archiveRecipeStep(recipeID, recipeStepID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a RecipeStep', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<RecipeStep>(buildObligatoryError('archiveRecipeStep'));
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(500, exampleResponse);

    expect(client.archiveRecipeStep(recipeID, recipeStepID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepCompletionConditionID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>();
    mock
      .onDelete(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(202, exampleResponse);

    client
      .archiveRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID)
      .then((response: APIResponse<RecipeStepCompletionCondition>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepCompletionConditionID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(
      buildObligatoryError('archiveRecipeStepCompletionCondition'),
    );
    mock
      .onDelete(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(202, exampleResponse);

    expect(
      client.archiveRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepCompletionConditionID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(
      buildObligatoryError('archiveRecipeStepCompletionCondition'),
    );
    mock
      .onDelete(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(500, exampleResponse);

    expect(
      client.archiveRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepIngredientID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepIngredient>();
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(202, exampleResponse);

    client
      .archiveRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)
      .then((response: APIResponse<RecipeStepIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepIngredientID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepIngredient>(buildObligatoryError('archiveRecipeStepIngredient'));
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(202, exampleResponse);

    expect(client.archiveRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepIngredientID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepIngredient>(buildObligatoryError('archiveRecipeStepIngredient'));
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.archiveRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepInstrumentID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepInstrument>();
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(202, exampleResponse);

    client
      .archiveRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)
      .then((response: APIResponse<RecipeStepInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepInstrumentID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepInstrument>(buildObligatoryError('archiveRecipeStepInstrument'));
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(202, exampleResponse);

    expect(client.archiveRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepInstrumentID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepInstrument>(buildObligatoryError('archiveRecipeStepInstrument'));
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(500, exampleResponse);

    expect(client.archiveRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepProductID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepProduct>();
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(202, exampleResponse);

    client
      .archiveRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)
      .then((response: APIResponse<RecipeStepProduct>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepProductID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepProduct>(buildObligatoryError('archiveRecipeStepProduct'));
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(202, exampleResponse);

    expect(client.archiveRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepProductID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepProduct>(buildObligatoryError('archiveRecipeStepProduct'));
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(500, exampleResponse);

    expect(client.archiveRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepVesselID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepVessel>();
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(202, exampleResponse);

    client
      .archiveRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)
      .then((response: APIResponse<RecipeStepVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepVesselID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepVessel>(buildObligatoryError('archiveRecipeStepVessel'));
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(202, exampleResponse);

    expect(client.archiveRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepVesselID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepVessel>(buildObligatoryError('archiveRecipeStepVessel'));
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(500, exampleResponse);

    expect(client.archiveRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a ServiceSetting', () => {
    let serviceSettingID = fakeID();

    const exampleResponse = new APIResponse<ServiceSetting>();
    mock.onDelete(`${baseURL}/api/v1/settings/${serviceSettingID}`).reply(202, exampleResponse);

    client
      .archiveServiceSetting(serviceSettingID)
      .then((response: APIResponse<ServiceSetting>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ServiceSetting', () => {
    let serviceSettingID = fakeID();

    const exampleResponse = new APIResponse<ServiceSetting>(buildObligatoryError('archiveServiceSetting'));
    mock.onDelete(`${baseURL}/api/v1/settings/${serviceSettingID}`).reply(202, exampleResponse);

    expect(client.archiveServiceSetting(serviceSettingID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a ServiceSetting', () => {
    let serviceSettingID = fakeID();

    const exampleResponse = new APIResponse<ServiceSetting>(buildObligatoryError('archiveServiceSetting'));
    mock.onDelete(`${baseURL}/api/v1/settings/${serviceSettingID}`).reply(500, exampleResponse);

    expect(client.archiveServiceSetting(serviceSettingID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a ServiceSettingConfiguration', () => {
    let serviceSettingConfigurationID = fakeID();

    const exampleResponse = new APIResponse<ServiceSettingConfiguration>();
    mock
      .onDelete(`${baseURL}/api/v1/settings/configurations/${serviceSettingConfigurationID}`)
      .reply(202, exampleResponse);

    client
      .archiveServiceSettingConfiguration(serviceSettingConfigurationID)
      .then((response: APIResponse<ServiceSettingConfiguration>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ServiceSettingConfiguration', () => {
    let serviceSettingConfigurationID = fakeID();

    const exampleResponse = new APIResponse<ServiceSettingConfiguration>(
      buildObligatoryError('archiveServiceSettingConfiguration'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/settings/configurations/${serviceSettingConfigurationID}`)
      .reply(202, exampleResponse);

    expect(client.archiveServiceSettingConfiguration(serviceSettingConfigurationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a ServiceSettingConfiguration', () => {
    let serviceSettingConfigurationID = fakeID();

    const exampleResponse = new APIResponse<ServiceSettingConfiguration>(
      buildObligatoryError('archiveServiceSettingConfiguration'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/settings/configurations/${serviceSettingConfigurationID}`)
      .reply(500, exampleResponse);

    expect(client.archiveServiceSettingConfiguration(serviceSettingConfigurationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a User', () => {
    let userID = fakeID();

    const exampleResponse = new APIResponse<User>();
    mock.onDelete(`${baseURL}/api/v1/users/${userID}`).reply(202, exampleResponse);

    client
      .archiveUser(userID)
      .then((response: APIResponse<User>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a User', () => {
    let userID = fakeID();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('archiveUser'));
    mock.onDelete(`${baseURL}/api/v1/users/${userID}`).reply(202, exampleResponse);

    expect(client.archiveUser(userID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a User', () => {
    let userID = fakeID();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('archiveUser'));
    mock.onDelete(`${baseURL}/api/v1/users/${userID}`).reply(500, exampleResponse);

    expect(client.archiveUser(userID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a UserIngredientPreference', () => {
    let userIngredientPreferenceID = fakeID();

    const exampleResponse = new APIResponse<UserIngredientPreference>();
    mock
      .onDelete(`${baseURL}/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`)
      .reply(202, exampleResponse);

    client
      .archiveUserIngredientPreference(userIngredientPreferenceID)
      .then((response: APIResponse<UserIngredientPreference>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a UserIngredientPreference', () => {
    let userIngredientPreferenceID = fakeID();

    const exampleResponse = new APIResponse<UserIngredientPreference>(
      buildObligatoryError('archiveUserIngredientPreference'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`)
      .reply(202, exampleResponse);

    expect(client.archiveUserIngredientPreference(userIngredientPreferenceID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a UserIngredientPreference', () => {
    let userIngredientPreferenceID = fakeID();

    const exampleResponse = new APIResponse<UserIngredientPreference>(
      buildObligatoryError('archiveUserIngredientPreference'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`)
      .reply(500, exampleResponse);

    expect(client.archiveUserIngredientPreference(userIngredientPreferenceID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a ValidIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredient>();
    mock.onDelete(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(202, exampleResponse);

    client
      .archiveValidIngredient(validIngredientID)
      .then((response: APIResponse<ValidIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredient>(buildObligatoryError('archiveValidIngredient'));
    mock.onDelete(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(202, exampleResponse);

    expect(client.archiveValidIngredient(validIngredientID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a ValidIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredient>(buildObligatoryError('archiveValidIngredient'));
    mock.onDelete(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(500, exampleResponse);

    expect(client.archiveValidIngredient(validIngredientID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a ValidIngredientGroup', () => {
    let validIngredientGroupID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientGroup>();
    mock.onDelete(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(202, exampleResponse);

    client
      .archiveValidIngredientGroup(validIngredientGroupID)
      .then((response: APIResponse<ValidIngredientGroup>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidIngredientGroup', () => {
    let validIngredientGroupID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientGroup>(buildObligatoryError('archiveValidIngredientGroup'));
    mock.onDelete(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(202, exampleResponse);

    expect(client.archiveValidIngredientGroup(validIngredientGroupID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidIngredientGroup', () => {
    let validIngredientGroupID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientGroup>(buildObligatoryError('archiveValidIngredientGroup'));
    mock.onDelete(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(500, exampleResponse);

    expect(client.archiveValidIngredientGroup(validIngredientGroupID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>();
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(202, exampleResponse);

    client
      .archiveValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)
      .then((response: APIResponse<ValidIngredientMeasurementUnit>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(
      buildObligatoryError('archiveValidIngredientMeasurementUnit'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(202, exampleResponse);

    expect(client.archiveValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(
      buildObligatoryError('archiveValidIngredientMeasurementUnit'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(500, exampleResponse);

    expect(client.archiveValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a ValidIngredientPreparation', () => {
    let validIngredientPreparationID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientPreparation>();
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(202, exampleResponse);

    client
      .archiveValidIngredientPreparation(validIngredientPreparationID)
      .then((response: APIResponse<ValidIngredientPreparation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidIngredientPreparation', () => {
    let validIngredientPreparationID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientPreparation>(
      buildObligatoryError('archiveValidIngredientPreparation'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(202, exampleResponse);

    expect(client.archiveValidIngredientPreparation(validIngredientPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidIngredientPreparation', () => {
    let validIngredientPreparationID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientPreparation>(
      buildObligatoryError('archiveValidIngredientPreparation'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(500, exampleResponse);

    expect(client.archiveValidIngredientPreparation(validIngredientPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a ValidIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientState>();
    mock.onDelete(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(202, exampleResponse);

    client
      .archiveValidIngredientState(validIngredientStateID)
      .then((response: APIResponse<ValidIngredientState>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientState>(buildObligatoryError('archiveValidIngredientState'));
    mock.onDelete(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(202, exampleResponse);

    expect(client.archiveValidIngredientState(validIngredientStateID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientState>(buildObligatoryError('archiveValidIngredientState'));
    mock.onDelete(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(500, exampleResponse);

    expect(client.archiveValidIngredientState(validIngredientStateID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>();
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(202, exampleResponse);

    client
      .archiveValidIngredientStateIngredient(validIngredientStateIngredientID)
      .then((response: APIResponse<ValidIngredientStateIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(
      buildObligatoryError('archiveValidIngredientStateIngredient'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(202, exampleResponse);

    expect(client.archiveValidIngredientStateIngredient(validIngredientStateIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(
      buildObligatoryError('archiveValidIngredientStateIngredient'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.archiveValidIngredientStateIngredient(validIngredientStateIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a ValidInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleResponse = new APIResponse<ValidInstrument>();
    mock.onDelete(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(202, exampleResponse);

    client
      .archiveValidInstrument(validInstrumentID)
      .then((response: APIResponse<ValidInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleResponse = new APIResponse<ValidInstrument>(buildObligatoryError('archiveValidInstrument'));
    mock.onDelete(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(202, exampleResponse);

    expect(client.archiveValidInstrument(validInstrumentID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a ValidInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleResponse = new APIResponse<ValidInstrument>(buildObligatoryError('archiveValidInstrument'));
    mock.onDelete(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(500, exampleResponse);

    expect(client.archiveValidInstrument(validInstrumentID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a ValidMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<ValidMeasurementUnit>();
    mock.onDelete(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(202, exampleResponse);

    client
      .archiveValidMeasurementUnit(validMeasurementUnitID)
      .then((response: APIResponse<ValidMeasurementUnit>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<ValidMeasurementUnit>(buildObligatoryError('archiveValidMeasurementUnit'));
    mock.onDelete(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(202, exampleResponse);

    expect(client.archiveValidMeasurementUnit(validMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<ValidMeasurementUnit>(buildObligatoryError('archiveValidMeasurementUnit'));
    mock.onDelete(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(500, exampleResponse);

    expect(client.archiveValidMeasurementUnit(validMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = fakeID();

    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>();
    mock
      .onDelete(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(202, exampleResponse);

    client
      .archiveValidMeasurementUnitConversion(validMeasurementUnitConversionID)
      .then((response: APIResponse<ValidMeasurementUnitConversion>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = fakeID();

    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(
      buildObligatoryError('archiveValidMeasurementUnitConversion'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(202, exampleResponse);

    expect(client.archiveValidMeasurementUnitConversion(validMeasurementUnitConversionID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = fakeID();

    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(
      buildObligatoryError('archiveValidMeasurementUnitConversion'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(500, exampleResponse);

    expect(client.archiveValidMeasurementUnitConversion(validMeasurementUnitConversionID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a ValidPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparation>();
    mock.onDelete(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(202, exampleResponse);

    client
      .archiveValidPreparation(validPreparationID)
      .then((response: APIResponse<ValidPreparation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparation>(buildObligatoryError('archiveValidPreparation'));
    mock.onDelete(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(202, exampleResponse);

    expect(client.archiveValidPreparation(validPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparation>(buildObligatoryError('archiveValidPreparation'));
    mock.onDelete(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(500, exampleResponse);

    expect(client.archiveValidPreparation(validPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a ValidPreparationInstrument', () => {
    let validPreparationVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparationInstrument>();
    mock
      .onDelete(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(202, exampleResponse);

    client
      .archiveValidPreparationInstrument(validPreparationVesselID)
      .then((response: APIResponse<ValidPreparationInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidPreparationInstrument', () => {
    let validPreparationVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparationInstrument>(
      buildObligatoryError('archiveValidPreparationInstrument'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(202, exampleResponse);

    expect(client.archiveValidPreparationInstrument(validPreparationVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidPreparationInstrument', () => {
    let validPreparationVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparationInstrument>(
      buildObligatoryError('archiveValidPreparationInstrument'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(500, exampleResponse);

    expect(client.archiveValidPreparationInstrument(validPreparationVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a ValidPreparationVessel', () => {
    let validPreparationVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparationVessel>();
    mock
      .onDelete(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`)
      .reply(202, exampleResponse);

    client
      .archiveValidPreparationVessel(validPreparationVesselID)
      .then((response: APIResponse<ValidPreparationVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidPreparationVessel', () => {
    let validPreparationVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparationVessel>(
      buildObligatoryError('archiveValidPreparationVessel'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`)
      .reply(202, exampleResponse);

    expect(client.archiveValidPreparationVessel(validPreparationVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidPreparationVessel', () => {
    let validPreparationVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparationVessel>(
      buildObligatoryError('archiveValidPreparationVessel'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`)
      .reply(500, exampleResponse);

    expect(client.archiveValidPreparationVessel(validPreparationVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a ValidVessel', () => {
    let validVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidVessel>();
    mock.onDelete(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(202, exampleResponse);

    client
      .archiveValidVessel(validVesselID)
      .then((response: APIResponse<ValidVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a ValidVessel', () => {
    let validVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidVessel>(buildObligatoryError('archiveValidVessel'));
    mock.onDelete(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(202, exampleResponse);

    expect(client.archiveValidVessel(validVesselID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a ValidVessel', () => {
    let validVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidVessel>(buildObligatoryError('archiveValidVessel'));
    mock.onDelete(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(500, exampleResponse);

    expect(client.archiveValidVessel(validVesselID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a Webhook', () => {
    let webhookID = fakeID();

    const exampleResponse = new APIResponse<Webhook>();
    mock.onDelete(`${baseURL}/api/v1/webhooks/${webhookID}`).reply(202, exampleResponse);

    client
      .archiveWebhook(webhookID)
      .then((response: APIResponse<Webhook>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a Webhook', () => {
    let webhookID = fakeID();

    const exampleResponse = new APIResponse<Webhook>(buildObligatoryError('archiveWebhook'));
    mock.onDelete(`${baseURL}/api/v1/webhooks/${webhookID}`).reply(202, exampleResponse);

    expect(client.archiveWebhook(webhookID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a Webhook', () => {
    let webhookID = fakeID();

    const exampleResponse = new APIResponse<Webhook>(buildObligatoryError('archiveWebhook'));
    mock.onDelete(`${baseURL}/api/v1/webhooks/${webhookID}`).reply(500, exampleResponse);

    expect(client.archiveWebhook(webhookID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a WebhookTriggerEvent', () => {
    let webhookID = fakeID();
    let webhookTriggerEventID = fakeID();

    const exampleResponse = new APIResponse<WebhookTriggerEvent>();
    mock
      .onDelete(`${baseURL}/api/v1/webhooks/${webhookID}/trigger_events/${webhookTriggerEventID}`)
      .reply(202, exampleResponse);

    client
      .archiveWebhookTriggerEvent(webhookID, webhookTriggerEventID)
      .then((response: APIResponse<WebhookTriggerEvent>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a WebhookTriggerEvent', () => {
    let webhookID = fakeID();
    let webhookTriggerEventID = fakeID();

    const exampleResponse = new APIResponse<WebhookTriggerEvent>(buildObligatoryError('archiveWebhookTriggerEvent'));
    mock
      .onDelete(`${baseURL}/api/v1/webhooks/${webhookID}/trigger_events/${webhookTriggerEventID}`)
      .reply(202, exampleResponse);

    expect(client.archiveWebhookTriggerEvent(webhookID, webhookTriggerEventID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a WebhookTriggerEvent', () => {
    let webhookID = fakeID();
    let webhookTriggerEventID = fakeID();

    const exampleResponse = new APIResponse<WebhookTriggerEvent>(buildObligatoryError('archiveWebhookTriggerEvent'));
    mock
      .onDelete(`${baseURL}/api/v1/webhooks/${webhookID}/trigger_events/${webhookTriggerEventID}`)
      .reply(500, exampleResponse);

    expect(client.archiveWebhookTriggerEvent(webhookID, webhookTriggerEventID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a household instrument ownership', () => {
    let householdInstrumentOwnershipID = fakeID();

    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>();
    mock
      .onDelete(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(202, exampleResponse);

    client
      .archiveHouseholdInstrumentOwnership(householdInstrumentOwnershipID)
      .then((response: APIResponse<HouseholdInstrumentOwnership>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a household instrument ownership', () => {
    let householdInstrumentOwnershipID = fakeID();

    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(
      buildObligatoryError('archiveHouseholdInstrumentOwnership'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(202, exampleResponse);

    expect(client.archiveHouseholdInstrumentOwnership(householdInstrumentOwnershipID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to archive a household instrument ownership', () => {
    let householdInstrumentOwnershipID = fakeID();

    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(
      buildObligatoryError('archiveHouseholdInstrumentOwnership'),
    );
    mock
      .onDelete(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(500, exampleResponse);

    expect(client.archiveHouseholdInstrumentOwnership(householdInstrumentOwnershipID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should archive a household', () => {
    let householdID = fakeID();

    const exampleResponse = new APIResponse<Household>();
    mock.onDelete(`${baseURL}/api/v1/households/${householdID}`).reply(202, exampleResponse);

    client
      .archiveHousehold(householdID)
      .then((response: APIResponse<Household>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a household', () => {
    let householdID = fakeID();

    const exampleResponse = new APIResponse<Household>(buildObligatoryError('archiveHousehold'));
    mock.onDelete(`${baseURL}/api/v1/households/${householdID}`).reply(202, exampleResponse);

    expect(client.archiveHousehold(householdID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a household', () => {
    let householdID = fakeID();

    const exampleResponse = new APIResponse<Household>(buildObligatoryError('archiveHousehold'));
    mock.onDelete(`${baseURL}/api/v1/households/${householdID}`).reply(500, exampleResponse);

    expect(client.archiveHousehold(householdID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should archive a meal plan option vote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();
    let mealPlanOptionVoteID = fakeID();

    const exampleResponse = new APIResponse<MealPlanOptionVote>();
    mock
      .onDelete(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(202, exampleResponse);

    client
      .archiveMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
      .then((response: APIResponse<MealPlanOptionVote>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.delete.length).toBe(1);
        expect(mock.history.delete[0].headers).toHaveProperty('Authorization');
        expect((mock.history.delete[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to archive a meal plan option vote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();
    let mealPlanOptionVoteID = fakeID();

    const exampleResponse = new APIResponse<MealPlanOptionVote>(buildObligatoryError('archiveMealPlanOptionVote'));
    mock
      .onDelete(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(202, exampleResponse);

    expect(
      client.archiveMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to archive a meal plan option vote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();
    let mealPlanOptionVoteID = fakeID();

    const exampleResponse = new APIResponse<MealPlanOptionVote>(buildObligatoryError('archiveMealPlanOptionVote'));
    mock
      .onDelete(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(500, exampleResponse);

    expect(
      client.archiveMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should cancel a sent household invitation', () => {
    let householdInvitationID = fakeID();

    const exampleInput = new HouseholdInvitationUpdateRequestInput();

    const exampleResponse = new APIResponse<HouseholdInvitation>();
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/cancel`).reply(200, exampleResponse);

    client
      .cancelHouseholdInvitation(householdInvitationID, exampleInput)
      .then((response: APIResponse<HouseholdInvitation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during cancel a sent household invitation', () => {
    let householdInvitationID = fakeID();

    const exampleInput = new HouseholdInvitationUpdateRequestInput();

    const exampleResponse = new APIResponse<HouseholdInvitation>(
      buildObligatoryError('cancelHouseholdInvitation user error'),
    );
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/cancel`).reply(200, exampleResponse);

    expect(client.cancelHouseholdInvitation(householdInvitationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during cancel a sent household invitation', () => {
    let householdInvitationID = fakeID();

    const exampleInput = new HouseholdInvitationUpdateRequestInput();

    const exampleResponse = new APIResponse<HouseholdInvitation>(
      buildObligatoryError('cancelHouseholdInvitation service error'),
    );
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/cancel`).reply(500, exampleResponse);

    expect(client.cancelHouseholdInvitation(householdInvitationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a Meal', () => {
    const exampleInput = new MealCreationRequestInput();

    const exampleResponse = new APIResponse<Meal>();
    mock.onPost(`${baseURL}/api/v1/meals`).reply(201, exampleResponse);

    client
      .createMeal(exampleInput)
      .then((response: APIResponse<Meal>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a Meal', () => {
    const exampleInput = new MealCreationRequestInput();

    const exampleResponse = new APIResponse<Meal>(buildObligatoryError('createMeal user error'));
    mock.onPost(`${baseURL}/api/v1/meals`).reply(201, exampleResponse);

    expect(client.createMeal(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a Meal', () => {
    const exampleInput = new MealCreationRequestInput();

    const exampleResponse = new APIResponse<Meal>(buildObligatoryError('createMeal service error'));
    mock.onPost(`${baseURL}/api/v1/meals`).reply(500, exampleResponse);

    expect(client.createMeal(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a MealPlan', () => {
    const exampleInput = new MealPlanCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlan>();
    mock.onPost(`${baseURL}/api/v1/meal_plans`).reply(201, exampleResponse);

    client
      .createMealPlan(exampleInput)
      .then((response: APIResponse<MealPlan>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a MealPlan', () => {
    const exampleInput = new MealPlanCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlan>(buildObligatoryError('createMealPlan user error'));
    mock.onPost(`${baseURL}/api/v1/meal_plans`).reply(201, exampleResponse);

    expect(client.createMealPlan(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a MealPlan', () => {
    const exampleInput = new MealPlanCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlan>(buildObligatoryError('createMealPlan service error'));
    mock.onPost(`${baseURL}/api/v1/meal_plans`).reply(500, exampleResponse);

    expect(client.createMealPlan(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a MealPlanEvent', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanEventCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlanEvent>();
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events`).reply(201, exampleResponse);

    client
      .createMealPlanEvent(mealPlanID, exampleInput)
      .then((response: APIResponse<MealPlanEvent>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a MealPlanEvent', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanEventCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlanEvent>(buildObligatoryError('createMealPlanEvent user error'));
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events`).reply(201, exampleResponse);

    expect(client.createMealPlanEvent(mealPlanID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a MealPlanEvent', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanEventCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlanEvent>(buildObligatoryError('createMealPlanEvent service error'));
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events`).reply(500, exampleResponse);

    expect(client.createMealPlanEvent(mealPlanID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanGroceryListItemCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlanGroceryListItem>();
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items`).reply(201, exampleResponse);

    client
      .createMealPlanGroceryListItem(mealPlanID, exampleInput)
      .then((response: APIResponse<MealPlanGroceryListItem>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanGroceryListItemCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(
      buildObligatoryError('createMealPlanGroceryListItem user error'),
    );
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items`).reply(201, exampleResponse);

    expect(client.createMealPlanGroceryListItem(mealPlanID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanGroceryListItemCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(
      buildObligatoryError('createMealPlanGroceryListItem service error'),
    );
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items`).reply(500, exampleResponse);

    expect(client.createMealPlanGroceryListItem(mealPlanID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a MealPlanOptionVote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleInput = new MealPlanOptionVoteCreationRequestInput();

    const exampleResponse = new APIResponse<Array<MealPlanOptionVote>>();
    mock
      .onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/vote`)
      .reply(201, exampleResponse);

    client
      .createMealPlanOptionVote(mealPlanID, mealPlanEventID, exampleInput)
      .then((response: APIResponse<Array<MealPlanOptionVote>>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a MealPlanOptionVote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleInput = new MealPlanOptionVoteCreationRequestInput();

    const exampleResponse = new APIResponse<Array<MealPlanOptionVote>>(
      buildObligatoryError('createMealPlanOptionVote user error'),
    );
    mock
      .onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/vote`)
      .reply(201, exampleResponse);

    expect(client.createMealPlanOptionVote(mealPlanID, mealPlanEventID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a MealPlanOptionVote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleInput = new MealPlanOptionVoteCreationRequestInput();

    const exampleResponse = new APIResponse<Array<MealPlanOptionVote>>(
      buildObligatoryError('createMealPlanOptionVote service error'),
    );
    mock
      .onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/vote`)
      .reply(500, exampleResponse);

    expect(client.createMealPlanOptionVote(mealPlanID, mealPlanEventID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a MealPlanTask', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanTaskCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlanTask>();
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks`).reply(201, exampleResponse);

    client
      .createMealPlanTask(mealPlanID, exampleInput)
      .then((response: APIResponse<MealPlanTask>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a MealPlanTask', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanTaskCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlanTask>(buildObligatoryError('createMealPlanTask user error'));
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks`).reply(201, exampleResponse);

    expect(client.createMealPlanTask(mealPlanID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a MealPlanTask', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanTaskCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlanTask>(buildObligatoryError('createMealPlanTask service error'));
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks`).reply(500, exampleResponse);

    expect(client.createMealPlanTask(mealPlanID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a OAuth2Client', () => {
    const exampleInput = new OAuth2ClientCreationRequestInput();

    const exampleResponse = new APIResponse<OAuth2ClientCreationResponse>();
    mock.onPost(`${baseURL}/api/v1/oauth2_clients`).reply(201, exampleResponse);

    client
      .createOAuth2Client(exampleInput)
      .then((response: APIResponse<OAuth2ClientCreationResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a OAuth2Client', () => {
    const exampleInput = new OAuth2ClientCreationRequestInput();

    const exampleResponse = new APIResponse<OAuth2ClientCreationResponse>(
      buildObligatoryError('createOAuth2Client user error'),
    );
    mock.onPost(`${baseURL}/api/v1/oauth2_clients`).reply(201, exampleResponse);

    expect(client.createOAuth2Client(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a OAuth2Client', () => {
    const exampleInput = new OAuth2ClientCreationRequestInput();

    const exampleResponse = new APIResponse<OAuth2ClientCreationResponse>(
      buildObligatoryError('createOAuth2Client service error'),
    );
    mock.onPost(`${baseURL}/api/v1/oauth2_clients`).reply(500, exampleResponse);

    expect(client.createOAuth2Client(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a Recipe', () => {
    const exampleInput = new RecipeCreationRequestInput();

    const exampleResponse = new APIResponse<Recipe>();
    mock.onPost(`${baseURL}/api/v1/recipes`).reply(201, exampleResponse);

    client
      .createRecipe(exampleInput)
      .then((response: APIResponse<Recipe>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a Recipe', () => {
    const exampleInput = new RecipeCreationRequestInput();

    const exampleResponse = new APIResponse<Recipe>(buildObligatoryError('createRecipe user error'));
    mock.onPost(`${baseURL}/api/v1/recipes`).reply(201, exampleResponse);

    expect(client.createRecipe(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a Recipe', () => {
    const exampleInput = new RecipeCreationRequestInput();

    const exampleResponse = new APIResponse<Recipe>(buildObligatoryError('createRecipe service error'));
    mock.onPost(`${baseURL}/api/v1/recipes`).reply(500, exampleResponse);

    expect(client.createRecipe(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a RecipePrepTask', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipePrepTaskCreationRequestInput();

    const exampleResponse = new APIResponse<RecipePrepTask>();
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks`).reply(201, exampleResponse);

    client
      .createRecipePrepTask(recipeID, exampleInput)
      .then((response: APIResponse<RecipePrepTask>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a RecipePrepTask', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipePrepTaskCreationRequestInput();

    const exampleResponse = new APIResponse<RecipePrepTask>(buildObligatoryError('createRecipePrepTask user error'));
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks`).reply(201, exampleResponse);

    expect(client.createRecipePrepTask(recipeID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a RecipePrepTask', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipePrepTaskCreationRequestInput();

    const exampleResponse = new APIResponse<RecipePrepTask>(buildObligatoryError('createRecipePrepTask service error'));
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks`).reply(500, exampleResponse);

    expect(client.createRecipePrepTask(recipeID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a RecipeRating', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipeRatingCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeRating>();
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/ratings`).reply(201, exampleResponse);

    client
      .createRecipeRating(recipeID, exampleInput)
      .then((response: APIResponse<RecipeRating>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a RecipeRating', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipeRatingCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeRating>(buildObligatoryError('createRecipeRating user error'));
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/ratings`).reply(201, exampleResponse);

    expect(client.createRecipeRating(recipeID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a RecipeRating', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipeRatingCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeRating>(buildObligatoryError('createRecipeRating service error'));
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/ratings`).reply(500, exampleResponse);

    expect(client.createRecipeRating(recipeID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a RecipeStep', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipeStepCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStep>();
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps`).reply(201, exampleResponse);

    client
      .createRecipeStep(recipeID, exampleInput)
      .then((response: APIResponse<RecipeStep>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a RecipeStep', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipeStepCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStep>(buildObligatoryError('createRecipeStep user error'));
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps`).reply(201, exampleResponse);

    expect(client.createRecipeStep(recipeID, exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a RecipeStep', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipeStepCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStep>(buildObligatoryError('createRecipeStep service error'));
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps`).reply(500, exampleResponse);

    expect(client.createRecipeStep(recipeID, exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepCompletionConditionForExistingRecipeCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>();
    mock
      .onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`)
      .reply(201, exampleResponse);

    client
      .createRecipeStepCompletionCondition(recipeID, recipeStepID, exampleInput)
      .then((response: APIResponse<RecipeStepCompletionCondition>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepCompletionConditionForExistingRecipeCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(
      buildObligatoryError('createRecipeStepCompletionCondition user error'),
    );
    mock
      .onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`)
      .reply(201, exampleResponse);

    expect(client.createRecipeStepCompletionCondition(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepCompletionConditionForExistingRecipeCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(
      buildObligatoryError('createRecipeStepCompletionCondition service error'),
    );
    mock
      .onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`)
      .reply(500, exampleResponse);

    expect(client.createRecipeStepCompletionCondition(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepIngredientCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepIngredient>();
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`).reply(201, exampleResponse);

    client
      .createRecipeStepIngredient(recipeID, recipeStepID, exampleInput)
      .then((response: APIResponse<RecipeStepIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepIngredientCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepIngredient>(
      buildObligatoryError('createRecipeStepIngredient user error'),
    );
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`).reply(201, exampleResponse);

    expect(client.createRecipeStepIngredient(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepIngredientCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepIngredient>(
      buildObligatoryError('createRecipeStepIngredient service error'),
    );
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`).reply(500, exampleResponse);

    expect(client.createRecipeStepIngredient(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepInstrumentCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepInstrument>();
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`).reply(201, exampleResponse);

    client
      .createRecipeStepInstrument(recipeID, recipeStepID, exampleInput)
      .then((response: APIResponse<RecipeStepInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepInstrumentCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepInstrument>(
      buildObligatoryError('createRecipeStepInstrument user error'),
    );
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`).reply(201, exampleResponse);

    expect(client.createRecipeStepInstrument(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepInstrumentCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepInstrument>(
      buildObligatoryError('createRecipeStepInstrument service error'),
    );
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`).reply(500, exampleResponse);

    expect(client.createRecipeStepInstrument(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepProductCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepProduct>();
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`).reply(201, exampleResponse);

    client
      .createRecipeStepProduct(recipeID, recipeStepID, exampleInput)
      .then((response: APIResponse<RecipeStepProduct>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepProductCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepProduct>(
      buildObligatoryError('createRecipeStepProduct user error'),
    );
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`).reply(201, exampleResponse);

    expect(client.createRecipeStepProduct(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepProductCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepProduct>(
      buildObligatoryError('createRecipeStepProduct service error'),
    );
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`).reply(500, exampleResponse);

    expect(client.createRecipeStepProduct(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepVesselCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepVessel>();
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`).reply(201, exampleResponse);

    client
      .createRecipeStepVessel(recipeID, recipeStepID, exampleInput)
      .then((response: APIResponse<RecipeStepVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepVesselCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepVessel>(
      buildObligatoryError('createRecipeStepVessel user error'),
    );
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`).reply(201, exampleResponse);

    expect(client.createRecipeStepVessel(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepVesselCreationRequestInput();

    const exampleResponse = new APIResponse<RecipeStepVessel>(
      buildObligatoryError('createRecipeStepVessel service error'),
    );
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`).reply(500, exampleResponse);

    expect(client.createRecipeStepVessel(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a ServiceSetting', () => {
    const exampleInput = new ServiceSettingCreationRequestInput();

    const exampleResponse = new APIResponse<ServiceSetting>();
    mock.onPost(`${baseURL}/api/v1/settings`).reply(201, exampleResponse);

    client
      .createServiceSetting(exampleInput)
      .then((response: APIResponse<ServiceSetting>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ServiceSetting', () => {
    const exampleInput = new ServiceSettingCreationRequestInput();

    const exampleResponse = new APIResponse<ServiceSetting>(buildObligatoryError('createServiceSetting user error'));
    mock.onPost(`${baseURL}/api/v1/settings`).reply(201, exampleResponse);

    expect(client.createServiceSetting(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a ServiceSetting', () => {
    const exampleInput = new ServiceSettingCreationRequestInput();

    const exampleResponse = new APIResponse<ServiceSetting>(buildObligatoryError('createServiceSetting service error'));
    mock.onPost(`${baseURL}/api/v1/settings`).reply(500, exampleResponse);

    expect(client.createServiceSetting(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a ServiceSettingConfiguration', () => {
    const exampleInput = new ServiceSettingConfigurationCreationRequestInput();

    const exampleResponse = new APIResponse<ServiceSettingConfiguration>();
    mock.onPost(`${baseURL}/api/v1/settings/configurations`).reply(201, exampleResponse);

    client
      .createServiceSettingConfiguration(exampleInput)
      .then((response: APIResponse<ServiceSettingConfiguration>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ServiceSettingConfiguration', () => {
    const exampleInput = new ServiceSettingConfigurationCreationRequestInput();

    const exampleResponse = new APIResponse<ServiceSettingConfiguration>(
      buildObligatoryError('createServiceSettingConfiguration user error'),
    );
    mock.onPost(`${baseURL}/api/v1/settings/configurations`).reply(201, exampleResponse);

    expect(client.createServiceSettingConfiguration(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a ServiceSettingConfiguration', () => {
    const exampleInput = new ServiceSettingConfigurationCreationRequestInput();

    const exampleResponse = new APIResponse<ServiceSettingConfiguration>(
      buildObligatoryError('createServiceSettingConfiguration service error'),
    );
    mock.onPost(`${baseURL}/api/v1/settings/configurations`).reply(500, exampleResponse);

    expect(client.createServiceSettingConfiguration(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a UserIngredientPreference', () => {
    const exampleInput = new UserIngredientPreferenceCreationRequestInput();

    const exampleResponse = new APIResponse<Array<UserIngredientPreference>>();
    mock.onPost(`${baseURL}/api/v1/user_ingredient_preferences`).reply(201, exampleResponse);

    client
      .createUserIngredientPreference(exampleInput)
      .then((response: APIResponse<Array<UserIngredientPreference>>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a UserIngredientPreference', () => {
    const exampleInput = new UserIngredientPreferenceCreationRequestInput();

    const exampleResponse = new APIResponse<Array<UserIngredientPreference>>(
      buildObligatoryError('createUserIngredientPreference user error'),
    );
    mock.onPost(`${baseURL}/api/v1/user_ingredient_preferences`).reply(201, exampleResponse);

    expect(client.createUserIngredientPreference(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a UserIngredientPreference', () => {
    const exampleInput = new UserIngredientPreferenceCreationRequestInput();

    const exampleResponse = new APIResponse<Array<UserIngredientPreference>>(
      buildObligatoryError('createUserIngredientPreference service error'),
    );
    mock.onPost(`${baseURL}/api/v1/user_ingredient_preferences`).reply(500, exampleResponse);

    expect(client.createUserIngredientPreference(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a UserNotification', () => {
    const exampleInput = new UserNotificationCreationRequestInput();

    const exampleResponse = new APIResponse<UserNotification>();
    mock.onPost(`${baseURL}/api/v1/user_notifications`).reply(201, exampleResponse);

    client
      .createUserNotification(exampleInput)
      .then((response: APIResponse<UserNotification>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a UserNotification', () => {
    const exampleInput = new UserNotificationCreationRequestInput();

    const exampleResponse = new APIResponse<UserNotification>(
      buildObligatoryError('createUserNotification user error'),
    );
    mock.onPost(`${baseURL}/api/v1/user_notifications`).reply(201, exampleResponse);

    expect(client.createUserNotification(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a UserNotification', () => {
    const exampleInput = new UserNotificationCreationRequestInput();

    const exampleResponse = new APIResponse<UserNotification>(
      buildObligatoryError('createUserNotification service error'),
    );
    mock.onPost(`${baseURL}/api/v1/user_notifications`).reply(500, exampleResponse);

    expect(client.createUserNotification(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a ValidIngredient', () => {
    const exampleInput = new ValidIngredientCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredient>();
    mock.onPost(`${baseURL}/api/v1/valid_ingredients`).reply(201, exampleResponse);

    client
      .createValidIngredient(exampleInput)
      .then((response: APIResponse<ValidIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidIngredient', () => {
    const exampleInput = new ValidIngredientCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredient>(buildObligatoryError('createValidIngredient user error'));
    mock.onPost(`${baseURL}/api/v1/valid_ingredients`).reply(201, exampleResponse);

    expect(client.createValidIngredient(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a ValidIngredient', () => {
    const exampleInput = new ValidIngredientCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredient>(
      buildObligatoryError('createValidIngredient service error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_ingredients`).reply(500, exampleResponse);

    expect(client.createValidIngredient(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a ValidIngredientGroup', () => {
    const exampleInput = new ValidIngredientGroupCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientGroup>();
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_groups`).reply(201, exampleResponse);

    client
      .createValidIngredientGroup(exampleInput)
      .then((response: APIResponse<ValidIngredientGroup>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidIngredientGroup', () => {
    const exampleInput = new ValidIngredientGroupCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientGroup>(
      buildObligatoryError('createValidIngredientGroup user error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_groups`).reply(201, exampleResponse);

    expect(client.createValidIngredientGroup(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a ValidIngredientGroup', () => {
    const exampleInput = new ValidIngredientGroupCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientGroup>(
      buildObligatoryError('createValidIngredientGroup service error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_groups`).reply(500, exampleResponse);

    expect(client.createValidIngredientGroup(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a ValidIngredientMeasurementUnit', () => {
    const exampleInput = new ValidIngredientMeasurementUnitCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>();
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_measurement_units`).reply(201, exampleResponse);

    client
      .createValidIngredientMeasurementUnit(exampleInput)
      .then((response: APIResponse<ValidIngredientMeasurementUnit>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidIngredientMeasurementUnit', () => {
    const exampleInput = new ValidIngredientMeasurementUnitCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(
      buildObligatoryError('createValidIngredientMeasurementUnit user error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_measurement_units`).reply(201, exampleResponse);

    expect(client.createValidIngredientMeasurementUnit(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a ValidIngredientMeasurementUnit', () => {
    const exampleInput = new ValidIngredientMeasurementUnitCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(
      buildObligatoryError('createValidIngredientMeasurementUnit service error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_measurement_units`).reply(500, exampleResponse);

    expect(client.createValidIngredientMeasurementUnit(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a ValidIngredientPreparation', () => {
    const exampleInput = new ValidIngredientPreparationCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientPreparation>();
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_preparations`).reply(201, exampleResponse);

    client
      .createValidIngredientPreparation(exampleInput)
      .then((response: APIResponse<ValidIngredientPreparation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidIngredientPreparation', () => {
    const exampleInput = new ValidIngredientPreparationCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientPreparation>(
      buildObligatoryError('createValidIngredientPreparation user error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_preparations`).reply(201, exampleResponse);

    expect(client.createValidIngredientPreparation(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a ValidIngredientPreparation', () => {
    const exampleInput = new ValidIngredientPreparationCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientPreparation>(
      buildObligatoryError('createValidIngredientPreparation service error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_preparations`).reply(500, exampleResponse);

    expect(client.createValidIngredientPreparation(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a ValidIngredientState', () => {
    const exampleInput = new ValidIngredientStateCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientState>();
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_states`).reply(201, exampleResponse);

    client
      .createValidIngredientState(exampleInput)
      .then((response: APIResponse<ValidIngredientState>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidIngredientState', () => {
    const exampleInput = new ValidIngredientStateCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientState>(
      buildObligatoryError('createValidIngredientState user error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_states`).reply(201, exampleResponse);

    expect(client.createValidIngredientState(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a ValidIngredientState', () => {
    const exampleInput = new ValidIngredientStateCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientState>(
      buildObligatoryError('createValidIngredientState service error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_states`).reply(500, exampleResponse);

    expect(client.createValidIngredientState(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a ValidIngredientStateIngredient', () => {
    const exampleInput = new ValidIngredientStateIngredientCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>();
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_state_ingredients`).reply(201, exampleResponse);

    client
      .createValidIngredientStateIngredient(exampleInput)
      .then((response: APIResponse<ValidIngredientStateIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidIngredientStateIngredient', () => {
    const exampleInput = new ValidIngredientStateIngredientCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(
      buildObligatoryError('createValidIngredientStateIngredient user error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_state_ingredients`).reply(201, exampleResponse);

    expect(client.createValidIngredientStateIngredient(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a ValidIngredientStateIngredient', () => {
    const exampleInput = new ValidIngredientStateIngredientCreationRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(
      buildObligatoryError('createValidIngredientStateIngredient service error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_state_ingredients`).reply(500, exampleResponse);

    expect(client.createValidIngredientStateIngredient(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a ValidInstrument', () => {
    const exampleInput = new ValidInstrumentCreationRequestInput();

    const exampleResponse = new APIResponse<ValidInstrument>();
    mock.onPost(`${baseURL}/api/v1/valid_instruments`).reply(201, exampleResponse);

    client
      .createValidInstrument(exampleInput)
      .then((response: APIResponse<ValidInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidInstrument', () => {
    const exampleInput = new ValidInstrumentCreationRequestInput();

    const exampleResponse = new APIResponse<ValidInstrument>(buildObligatoryError('createValidInstrument user error'));
    mock.onPost(`${baseURL}/api/v1/valid_instruments`).reply(201, exampleResponse);

    expect(client.createValidInstrument(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a ValidInstrument', () => {
    const exampleInput = new ValidInstrumentCreationRequestInput();

    const exampleResponse = new APIResponse<ValidInstrument>(
      buildObligatoryError('createValidInstrument service error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_instruments`).reply(500, exampleResponse);

    expect(client.createValidInstrument(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a ValidMeasurementUnit', () => {
    const exampleInput = new ValidMeasurementUnitCreationRequestInput();

    const exampleResponse = new APIResponse<ValidMeasurementUnit>();
    mock.onPost(`${baseURL}/api/v1/valid_measurement_units`).reply(201, exampleResponse);

    client
      .createValidMeasurementUnit(exampleInput)
      .then((response: APIResponse<ValidMeasurementUnit>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidMeasurementUnit', () => {
    const exampleInput = new ValidMeasurementUnitCreationRequestInput();

    const exampleResponse = new APIResponse<ValidMeasurementUnit>(
      buildObligatoryError('createValidMeasurementUnit user error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_measurement_units`).reply(201, exampleResponse);

    expect(client.createValidMeasurementUnit(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a ValidMeasurementUnit', () => {
    const exampleInput = new ValidMeasurementUnitCreationRequestInput();

    const exampleResponse = new APIResponse<ValidMeasurementUnit>(
      buildObligatoryError('createValidMeasurementUnit service error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_measurement_units`).reply(500, exampleResponse);

    expect(client.createValidMeasurementUnit(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a ValidMeasurementUnitConversion', () => {
    const exampleInput = new ValidMeasurementUnitConversionCreationRequestInput();

    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>();
    mock.onPost(`${baseURL}/api/v1/valid_measurement_conversions`).reply(201, exampleResponse);

    client
      .createValidMeasurementUnitConversion(exampleInput)
      .then((response: APIResponse<ValidMeasurementUnitConversion>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidMeasurementUnitConversion', () => {
    const exampleInput = new ValidMeasurementUnitConversionCreationRequestInput();

    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(
      buildObligatoryError('createValidMeasurementUnitConversion user error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_measurement_conversions`).reply(201, exampleResponse);

    expect(client.createValidMeasurementUnitConversion(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a ValidMeasurementUnitConversion', () => {
    const exampleInput = new ValidMeasurementUnitConversionCreationRequestInput();

    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(
      buildObligatoryError('createValidMeasurementUnitConversion service error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_measurement_conversions`).reply(500, exampleResponse);

    expect(client.createValidMeasurementUnitConversion(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a ValidPreparation', () => {
    const exampleInput = new ValidPreparationCreationRequestInput();

    const exampleResponse = new APIResponse<ValidPreparation>();
    mock.onPost(`${baseURL}/api/v1/valid_preparations`).reply(201, exampleResponse);

    client
      .createValidPreparation(exampleInput)
      .then((response: APIResponse<ValidPreparation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidPreparation', () => {
    const exampleInput = new ValidPreparationCreationRequestInput();

    const exampleResponse = new APIResponse<ValidPreparation>(
      buildObligatoryError('createValidPreparation user error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_preparations`).reply(201, exampleResponse);

    expect(client.createValidPreparation(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a ValidPreparation', () => {
    const exampleInput = new ValidPreparationCreationRequestInput();

    const exampleResponse = new APIResponse<ValidPreparation>(
      buildObligatoryError('createValidPreparation service error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_preparations`).reply(500, exampleResponse);

    expect(client.createValidPreparation(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a ValidPreparationInstrument', () => {
    const exampleInput = new ValidPreparationInstrumentCreationRequestInput();

    const exampleResponse = new APIResponse<ValidPreparationInstrument>();
    mock.onPost(`${baseURL}/api/v1/valid_preparation_instruments`).reply(201, exampleResponse);

    client
      .createValidPreparationInstrument(exampleInput)
      .then((response: APIResponse<ValidPreparationInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidPreparationInstrument', () => {
    const exampleInput = new ValidPreparationInstrumentCreationRequestInput();

    const exampleResponse = new APIResponse<ValidPreparationInstrument>(
      buildObligatoryError('createValidPreparationInstrument user error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_preparation_instruments`).reply(201, exampleResponse);

    expect(client.createValidPreparationInstrument(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a ValidPreparationInstrument', () => {
    const exampleInput = new ValidPreparationInstrumentCreationRequestInput();

    const exampleResponse = new APIResponse<ValidPreparationInstrument>(
      buildObligatoryError('createValidPreparationInstrument service error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_preparation_instruments`).reply(500, exampleResponse);

    expect(client.createValidPreparationInstrument(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a ValidPreparationVessel', () => {
    const exampleInput = new ValidPreparationVesselCreationRequestInput();

    const exampleResponse = new APIResponse<ValidPreparationVessel>();
    mock.onPost(`${baseURL}/api/v1/valid_preparation_vessels`).reply(201, exampleResponse);

    client
      .createValidPreparationVessel(exampleInput)
      .then((response: APIResponse<ValidPreparationVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidPreparationVessel', () => {
    const exampleInput = new ValidPreparationVesselCreationRequestInput();

    const exampleResponse = new APIResponse<ValidPreparationVessel>(
      buildObligatoryError('createValidPreparationVessel user error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_preparation_vessels`).reply(201, exampleResponse);

    expect(client.createValidPreparationVessel(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a ValidPreparationVessel', () => {
    const exampleInput = new ValidPreparationVesselCreationRequestInput();

    const exampleResponse = new APIResponse<ValidPreparationVessel>(
      buildObligatoryError('createValidPreparationVessel service error'),
    );
    mock.onPost(`${baseURL}/api/v1/valid_preparation_vessels`).reply(500, exampleResponse);

    expect(client.createValidPreparationVessel(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a ValidVessel', () => {
    const exampleInput = new ValidVesselCreationRequestInput();

    const exampleResponse = new APIResponse<ValidVessel>();
    mock.onPost(`${baseURL}/api/v1/valid_vessels`).reply(201, exampleResponse);

    client
      .createValidVessel(exampleInput)
      .then((response: APIResponse<ValidVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a ValidVessel', () => {
    const exampleInput = new ValidVesselCreationRequestInput();

    const exampleResponse = new APIResponse<ValidVessel>(buildObligatoryError('createValidVessel user error'));
    mock.onPost(`${baseURL}/api/v1/valid_vessels`).reply(201, exampleResponse);

    expect(client.createValidVessel(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a ValidVessel', () => {
    const exampleInput = new ValidVesselCreationRequestInput();

    const exampleResponse = new APIResponse<ValidVessel>(buildObligatoryError('createValidVessel service error'));
    mock.onPost(`${baseURL}/api/v1/valid_vessels`).reply(500, exampleResponse);

    expect(client.createValidVessel(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a Webhook', () => {
    const exampleInput = new WebhookCreationRequestInput();

    const exampleResponse = new APIResponse<Webhook>();
    mock.onPost(`${baseURL}/api/v1/webhooks`).reply(201, exampleResponse);

    client
      .createWebhook(exampleInput)
      .then((response: APIResponse<Webhook>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a Webhook', () => {
    const exampleInput = new WebhookCreationRequestInput();

    const exampleResponse = new APIResponse<Webhook>(buildObligatoryError('createWebhook user error'));
    mock.onPost(`${baseURL}/api/v1/webhooks`).reply(201, exampleResponse);

    expect(client.createWebhook(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a Webhook', () => {
    const exampleInput = new WebhookCreationRequestInput();

    const exampleResponse = new APIResponse<Webhook>(buildObligatoryError('createWebhook service error'));
    mock.onPost(`${baseURL}/api/v1/webhooks`).reply(500, exampleResponse);

    expect(client.createWebhook(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a WebhookTriggerEvent', () => {
    let webhookID = fakeID();

    const exampleInput = new WebhookTriggerEventCreationRequestInput();

    const exampleResponse = new APIResponse<WebhookTriggerEvent>();
    mock.onPost(`${baseURL}/api/v1/webhooks/${webhookID}/trigger_events`).reply(201, exampleResponse);

    client
      .createWebhookTriggerEvent(webhookID, exampleInput)
      .then((response: APIResponse<WebhookTriggerEvent>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a WebhookTriggerEvent', () => {
    let webhookID = fakeID();

    const exampleInput = new WebhookTriggerEventCreationRequestInput();

    const exampleResponse = new APIResponse<WebhookTriggerEvent>(
      buildObligatoryError('createWebhookTriggerEvent user error'),
    );
    mock.onPost(`${baseURL}/api/v1/webhooks/${webhookID}/trigger_events`).reply(201, exampleResponse);

    expect(client.createWebhookTriggerEvent(webhookID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a WebhookTriggerEvent', () => {
    let webhookID = fakeID();

    const exampleInput = new WebhookTriggerEventCreationRequestInput();

    const exampleResponse = new APIResponse<WebhookTriggerEvent>(
      buildObligatoryError('createWebhookTriggerEvent service error'),
    );
    mock.onPost(`${baseURL}/api/v1/webhooks/${webhookID}/trigger_events`).reply(500, exampleResponse);

    expect(client.createWebhookTriggerEvent(webhookID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a household instrument ownership', () => {
    const exampleInput = new HouseholdInstrumentOwnershipCreationRequestInput();

    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>();
    mock.onPost(`${baseURL}/api/v1/households/instruments`).reply(201, exampleResponse);

    client
      .createHouseholdInstrumentOwnership(exampleInput)
      .then((response: APIResponse<HouseholdInstrumentOwnership>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a household instrument ownership', () => {
    const exampleInput = new HouseholdInstrumentOwnershipCreationRequestInput();

    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(
      buildObligatoryError('createHouseholdInstrumentOwnership user error'),
    );
    mock.onPost(`${baseURL}/api/v1/households/instruments`).reply(201, exampleResponse);

    expect(client.createHouseholdInstrumentOwnership(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a household instrument ownership', () => {
    const exampleInput = new HouseholdInstrumentOwnershipCreationRequestInput();

    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(
      buildObligatoryError('createHouseholdInstrumentOwnership service error'),
    );
    mock.onPost(`${baseURL}/api/v1/households/instruments`).reply(500, exampleResponse);

    expect(client.createHouseholdInstrumentOwnership(exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a household', () => {
    const exampleInput = new HouseholdCreationRequestInput();

    const exampleResponse = new APIResponse<Household>();
    mock.onPost(`${baseURL}/api/v1/households`).reply(201, exampleResponse);

    client
      .createHousehold(exampleInput)
      .then((response: APIResponse<Household>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a household', () => {
    const exampleInput = new HouseholdCreationRequestInput();

    const exampleResponse = new APIResponse<Household>(buildObligatoryError('createHousehold user error'));
    mock.onPost(`${baseURL}/api/v1/households`).reply(201, exampleResponse);

    expect(client.createHousehold(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a household', () => {
    const exampleInput = new HouseholdCreationRequestInput();

    const exampleResponse = new APIResponse<Household>(buildObligatoryError('createHousehold service error'));
    mock.onPost(`${baseURL}/api/v1/households`).reply(500, exampleResponse);

    expect(client.createHousehold(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should create a meal plan option', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleInput = new MealPlanOptionCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlanOption>();
    mock
      .onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`)
      .reply(201, exampleResponse);

    client
      .createMealPlanOption(mealPlanID, mealPlanEventID, exampleInput)
      .then((response: APIResponse<MealPlanOption>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a meal plan option', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleInput = new MealPlanOptionCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlanOption>(buildObligatoryError('createMealPlanOption user error'));
    mock
      .onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`)
      .reply(201, exampleResponse);

    expect(client.createMealPlanOption(mealPlanID, mealPlanEventID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during create a meal plan option', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleInput = new MealPlanOptionCreationRequestInput();

    const exampleResponse = new APIResponse<MealPlanOption>(buildObligatoryError('createMealPlanOption service error'));
    mock
      .onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`)
      .reply(500, exampleResponse);

    expect(client.createMealPlanOption(mealPlanID, mealPlanEventID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should create a new user', () => {
    const exampleInput = new UserRegistrationInput();

    const exampleResponse = new APIResponse<UserCreationResponse>();
    mock.onPost(`${baseURL}/users`).reply(201, exampleResponse);

    client
      .createUser(exampleInput)
      .then((response: APIResponse<UserCreationResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during create a new user', () => {
    const exampleInput = new UserRegistrationInput();

    const exampleResponse = new APIResponse<UserCreationResponse>(buildObligatoryError('createUser user error'));
    mock.onPost(`${baseURL}/users`).reply(201, exampleResponse);

    expect(client.createUser(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during create a new user', () => {
    const exampleInput = new UserRegistrationInput();

    const exampleResponse = new APIResponse<UserCreationResponse>(buildObligatoryError('createUser service error'));
    mock.onPost(`${baseURL}/users`).reply(500, exampleResponse);

    expect(client.createUser(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should creates a household instrument', () => {
    const exampleResponse = new APIResponse<Array<HouseholdInstrumentOwnership>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new HouseholdInstrumentOwnership()],
    });
    mock.onGet(`${baseURL}/api/v1/households/instruments`).reply(200, exampleResponse);

    client
      .getHouseholdInstrumentOwnerships()
      .then((response: QueryFilteredResult<HouseholdInstrumentOwnership>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to creates a household instrument', () => {
    const exampleResponse = new APIResponse<Array<HouseholdInstrumentOwnership>>(
      buildObligatoryError('getHouseholdInstrumentOwnerships user error'),
    );
    mock.onGet(`${baseURL}/api/v1/households/instruments`).reply(200, exampleResponse);

    expect(client.getHouseholdInstrumentOwnerships()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to creates a household instrument', () => {
    const exampleResponse = new APIResponse<Array<HouseholdInstrumentOwnership>>(
      buildObligatoryError('getHouseholdInstrumentOwnerships service error'),
    );
    mock.onGet(`${baseURL}/api/v1/households/instruments`).reply(500, exampleResponse);

    expect(client.getHouseholdInstrumentOwnerships()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a AuthStatus', () => {
    const exampleResponse = new APIResponse<UserStatusResponse>();
    mock.onGet(`${baseURL}/auth/status`).reply(200, exampleResponse);

    client
      .getAuthStatus()
      .then((response: APIResponse<UserStatusResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a AuthStatus', () => {
    const exampleResponse = new APIResponse<UserStatusResponse>(buildObligatoryError('getAuthStatus'));
    mock.onGet(`${baseURL}/auth/status`).reply(200, exampleResponse);

    expect(client.getAuthStatus()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a AuthStatus', () => {
    const exampleResponse = new APIResponse<UserStatusResponse>(buildObligatoryError('getAuthStatus'));
    mock.onGet(`${baseURL}/auth/status`).reply(500, exampleResponse);

    expect(client.getAuthStatus()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a Meal', () => {
    let mealID = fakeID();

    const exampleResponse = new APIResponse<Meal>();
    mock.onGet(`${baseURL}/api/v1/meals/${mealID}`).reply(200, exampleResponse);

    client
      .getMeal(mealID)
      .then((response: APIResponse<Meal>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a Meal', () => {
    let mealID = fakeID();

    const exampleResponse = new APIResponse<Meal>(buildObligatoryError('getMeal'));
    mock.onGet(`${baseURL}/api/v1/meals/${mealID}`).reply(200, exampleResponse);

    expect(client.getMeal(mealID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a Meal', () => {
    let mealID = fakeID();

    const exampleResponse = new APIResponse<Meal>(buildObligatoryError('getMeal'));
    mock.onGet(`${baseURL}/api/v1/meals/${mealID}`).reply(500, exampleResponse);

    expect(client.getMeal(mealID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a MealPlan', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<MealPlan>();
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(200, exampleResponse);

    client
      .getMealPlan(mealPlanID)
      .then((response: APIResponse<MealPlan>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a MealPlan', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<MealPlan>(buildObligatoryError('getMealPlan'));
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(200, exampleResponse);

    expect(client.getMealPlan(mealPlanID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a MealPlan', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<MealPlan>(buildObligatoryError('getMealPlan'));
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(500, exampleResponse);

    expect(client.getMealPlan(mealPlanID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a MealPlanEvent', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleResponse = new APIResponse<MealPlanEvent>();
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(200, exampleResponse);

    client
      .getMealPlanEvent(mealPlanID, mealPlanEventID)
      .then((response: APIResponse<MealPlanEvent>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a MealPlanEvent', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleResponse = new APIResponse<MealPlanEvent>(buildObligatoryError('getMealPlanEvent'));
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(200, exampleResponse);

    expect(client.getMealPlanEvent(mealPlanID, mealPlanEventID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a MealPlanEvent', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleResponse = new APIResponse<MealPlanEvent>(buildObligatoryError('getMealPlanEvent'));
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(500, exampleResponse);

    expect(client.getMealPlanEvent(mealPlanID, mealPlanEventID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a MealPlanEvents', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanEvent>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new MealPlanEvent()],
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events`).reply(200, exampleResponse);

    client
      .getMealPlanEvents(mealPlanID)
      .then((response: QueryFilteredResult<MealPlanEvent>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a MealPlanEvents', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanEvent>>(buildObligatoryError('getMealPlanEvents user error'));
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events`).reply(200, exampleResponse);

    expect(client.getMealPlanEvents(mealPlanID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a MealPlanEvents', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanEvent>>(
      buildObligatoryError('getMealPlanEvents service error'),
    );
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events`).reply(500, exampleResponse);

    expect(client.getMealPlanEvents(mealPlanID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();
    let mealPlanGroceryListItemID = fakeID();

    const exampleResponse = new APIResponse<MealPlanGroceryListItem>();
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(200, exampleResponse);

    client
      .getMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)
      .then((response: APIResponse<MealPlanGroceryListItem>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();
    let mealPlanGroceryListItemID = fakeID();

    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(
      buildObligatoryError('getMealPlanGroceryListItem'),
    );
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(200, exampleResponse);

    expect(client.getMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();
    let mealPlanGroceryListItemID = fakeID();

    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(
      buildObligatoryError('getMealPlanGroceryListItem'),
    );
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(500, exampleResponse);

    expect(client.getMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a MealPlanGroceryListItemsForMealPlan', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanGroceryListItem>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new MealPlanGroceryListItem()],
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items`).reply(200, exampleResponse);

    client
      .getMealPlanGroceryListItemsForMealPlan(mealPlanID)
      .then((response: QueryFilteredResult<MealPlanGroceryListItem>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a MealPlanGroceryListItemsForMealPlan', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanGroceryListItem>>(
      buildObligatoryError('getMealPlanGroceryListItemsForMealPlan user error'),
    );
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items`).reply(200, exampleResponse);

    expect(client.getMealPlanGroceryListItemsForMealPlan(mealPlanID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a MealPlanGroceryListItemsForMealPlan', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanGroceryListItem>>(
      buildObligatoryError('getMealPlanGroceryListItemsForMealPlan service error'),
    );
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items`).reply(500, exampleResponse);

    expect(client.getMealPlanGroceryListItemsForMealPlan(mealPlanID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a MealPlanTask', () => {
    let mealPlanID = fakeID();
    let mealPlanTaskID = fakeID();

    const exampleResponse = new APIResponse<MealPlanTask>();
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`).reply(200, exampleResponse);

    client
      .getMealPlanTask(mealPlanID, mealPlanTaskID)
      .then((response: APIResponse<MealPlanTask>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a MealPlanTask', () => {
    let mealPlanID = fakeID();
    let mealPlanTaskID = fakeID();

    const exampleResponse = new APIResponse<MealPlanTask>(buildObligatoryError('getMealPlanTask'));
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`).reply(200, exampleResponse);

    expect(client.getMealPlanTask(mealPlanID, mealPlanTaskID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a MealPlanTask', () => {
    let mealPlanID = fakeID();
    let mealPlanTaskID = fakeID();

    const exampleResponse = new APIResponse<MealPlanTask>(buildObligatoryError('getMealPlanTask'));
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`).reply(500, exampleResponse);

    expect(client.getMealPlanTask(mealPlanID, mealPlanTaskID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a MealPlanTasks', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanTask>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new MealPlanTask()],
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks`).reply(200, exampleResponse);

    client
      .getMealPlanTasks(mealPlanID)
      .then((response: QueryFilteredResult<MealPlanTask>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a MealPlanTasks', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanTask>>(buildObligatoryError('getMealPlanTasks user error'));
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks`).reply(200, exampleResponse);

    expect(client.getMealPlanTasks(mealPlanID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a MealPlanTasks', () => {
    let mealPlanID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanTask>>(
      buildObligatoryError('getMealPlanTasks service error'),
    );
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks`).reply(500, exampleResponse);

    expect(client.getMealPlanTasks(mealPlanID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a MealPlans', () => {
    const exampleResponse = new APIResponse<Array<MealPlan>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new MealPlan()],
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans`).reply(200, exampleResponse);

    client
      .getMealPlansForHousehold()
      .then((response: QueryFilteredResult<MealPlan>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a MealPlans', () => {
    const exampleResponse = new APIResponse<Array<MealPlan>>(
      buildObligatoryError('getMealPlansForHousehold user error'),
    );
    mock.onGet(`${baseURL}/api/v1/meal_plans`).reply(200, exampleResponse);

    expect(client.getMealPlansForHousehold()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a MealPlans', () => {
    const exampleResponse = new APIResponse<Array<MealPlan>>(
      buildObligatoryError('getMealPlansForHousehold service error'),
    );
    mock.onGet(`${baseURL}/api/v1/meal_plans`).reply(500, exampleResponse);

    expect(client.getMealPlansForHousehold()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a Meals', () => {
    const exampleResponse = new APIResponse<Array<Meal>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new Meal()],
    });
    mock.onGet(`${baseURL}/api/v1/meals`).reply(200, exampleResponse);

    client
      .getMeals()
      .then((response: QueryFilteredResult<Meal>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a Meals', () => {
    const exampleResponse = new APIResponse<Array<Meal>>(buildObligatoryError('getMeals user error'));
    mock.onGet(`${baseURL}/api/v1/meals`).reply(200, exampleResponse);

    expect(client.getMeals()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a Meals', () => {
    const exampleResponse = new APIResponse<Array<Meal>>(buildObligatoryError('getMeals service error'));
    mock.onGet(`${baseURL}/api/v1/meals`).reply(500, exampleResponse);

    expect(client.getMeals()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a MermaidDiagramForRecipe', () => {
    let recipeID = fakeID();

    const exampleResponse = '';
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/mermaid`).reply(200, exampleResponse);

    client
      .getMermaidDiagramForRecipe(recipeID)
      .then((response: APIResponse<string>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a MermaidDiagramForRecipe', () => {
    let recipeID = fakeID();

    const exampleResponse = buildObligatoryError('getMermaidDiagramForRecipe user error');
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/mermaid`).reply(200, exampleResponse);

    expect(client.getMermaidDiagramForRecipe(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a MermaidDiagramForRecipe', () => {
    let recipeID = fakeID();

    const exampleResponse = buildObligatoryError('getMermaidDiagramForRecipe service error');
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/mermaid`).reply(500, exampleResponse);

    expect(client.getMermaidDiagramForRecipe(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a OAuth2Client', () => {
    let oauth2ClientID = fakeID();

    const exampleResponse = new APIResponse<OAuth2Client>();
    mock.onGet(`${baseURL}/api/v1/oauth2_clients/${oauth2ClientID}`).reply(200, exampleResponse);

    client
      .getOAuth2Client(oauth2ClientID)
      .then((response: APIResponse<OAuth2Client>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a OAuth2Client', () => {
    let oauth2ClientID = fakeID();

    const exampleResponse = new APIResponse<OAuth2Client>(buildObligatoryError('getOAuth2Client'));
    mock.onGet(`${baseURL}/api/v1/oauth2_clients/${oauth2ClientID}`).reply(200, exampleResponse);

    expect(client.getOAuth2Client(oauth2ClientID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a OAuth2Client', () => {
    let oauth2ClientID = fakeID();

    const exampleResponse = new APIResponse<OAuth2Client>(buildObligatoryError('getOAuth2Client'));
    mock.onGet(`${baseURL}/api/v1/oauth2_clients/${oauth2ClientID}`).reply(500, exampleResponse);

    expect(client.getOAuth2Client(oauth2ClientID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a OAuth2Clients', () => {
    const exampleResponse = new APIResponse<Array<OAuth2Client>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new OAuth2Client()],
    });
    mock.onGet(`${baseURL}/api/v1/oauth2_clients`).reply(200, exampleResponse);

    client
      .getOAuth2Clients()
      .then((response: QueryFilteredResult<OAuth2Client>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a OAuth2Clients', () => {
    const exampleResponse = new APIResponse<Array<OAuth2Client>>(buildObligatoryError('getOAuth2Clients user error'));
    mock.onGet(`${baseURL}/api/v1/oauth2_clients`).reply(200, exampleResponse);

    expect(client.getOAuth2Clients()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a OAuth2Clients', () => {
    const exampleResponse = new APIResponse<Array<OAuth2Client>>(
      buildObligatoryError('getOAuth2Clients service error'),
    );
    mock.onGet(`${baseURL}/api/v1/oauth2_clients`).reply(500, exampleResponse);

    expect(client.getOAuth2Clients()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a RandomValidIngredient', () => {
    const exampleResponse = new APIResponse<ValidIngredient>();
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/random`).reply(200, exampleResponse);

    client
      .getRandomValidIngredient()
      .then((response: APIResponse<ValidIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RandomValidIngredient', () => {
    const exampleResponse = new APIResponse<ValidIngredient>(buildObligatoryError('getRandomValidIngredient'));
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/random`).reply(200, exampleResponse);

    expect(client.getRandomValidIngredient()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a RandomValidIngredient', () => {
    const exampleResponse = new APIResponse<ValidIngredient>(buildObligatoryError('getRandomValidIngredient'));
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/random`).reply(500, exampleResponse);

    expect(client.getRandomValidIngredient()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a RandomValidInstrument', () => {
    const exampleResponse = new APIResponse<ValidInstrument>();
    mock.onGet(`${baseURL}/api/v1/valid_instruments/random`).reply(200, exampleResponse);

    client
      .getRandomValidInstrument()
      .then((response: APIResponse<ValidInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RandomValidInstrument', () => {
    const exampleResponse = new APIResponse<ValidInstrument>(buildObligatoryError('getRandomValidInstrument'));
    mock.onGet(`${baseURL}/api/v1/valid_instruments/random`).reply(200, exampleResponse);

    expect(client.getRandomValidInstrument()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a RandomValidInstrument', () => {
    const exampleResponse = new APIResponse<ValidInstrument>(buildObligatoryError('getRandomValidInstrument'));
    mock.onGet(`${baseURL}/api/v1/valid_instruments/random`).reply(500, exampleResponse);

    expect(client.getRandomValidInstrument()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a RandomValidPreparation', () => {
    const exampleResponse = new APIResponse<ValidPreparation>();
    mock.onGet(`${baseURL}/api/v1/valid_preparations/random`).reply(200, exampleResponse);

    client
      .getRandomValidPreparation()
      .then((response: APIResponse<ValidPreparation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RandomValidPreparation', () => {
    const exampleResponse = new APIResponse<ValidPreparation>(buildObligatoryError('getRandomValidPreparation'));
    mock.onGet(`${baseURL}/api/v1/valid_preparations/random`).reply(200, exampleResponse);

    expect(client.getRandomValidPreparation()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a RandomValidPreparation', () => {
    const exampleResponse = new APIResponse<ValidPreparation>(buildObligatoryError('getRandomValidPreparation'));
    mock.onGet(`${baseURL}/api/v1/valid_preparations/random`).reply(500, exampleResponse);

    expect(client.getRandomValidPreparation()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a RandomValidVessel', () => {
    const exampleResponse = new APIResponse<ValidVessel>();
    mock.onGet(`${baseURL}/api/v1/valid_vessels/random`).reply(200, exampleResponse);

    client
      .getRandomValidVessel()
      .then((response: APIResponse<ValidVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RandomValidVessel', () => {
    const exampleResponse = new APIResponse<ValidVessel>(buildObligatoryError('getRandomValidVessel'));
    mock.onGet(`${baseURL}/api/v1/valid_vessels/random`).reply(200, exampleResponse);

    expect(client.getRandomValidVessel()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a RandomValidVessel', () => {
    const exampleResponse = new APIResponse<ValidVessel>(buildObligatoryError('getRandomValidVessel'));
    mock.onGet(`${baseURL}/api/v1/valid_vessels/random`).reply(500, exampleResponse);

    expect(client.getRandomValidVessel()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a Recipe', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Recipe>();
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}`).reply(200, exampleResponse);

    client
      .getRecipe(recipeID)
      .then((response: APIResponse<Recipe>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a Recipe', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Recipe>(buildObligatoryError('getRecipe'));
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}`).reply(200, exampleResponse);

    expect(client.getRecipe(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a Recipe', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Recipe>(buildObligatoryError('getRecipe'));
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}`).reply(500, exampleResponse);

    expect(client.getRecipe(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a RecipeMealPlanTasks', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<RecipePrepTaskStep>();
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_steps`).reply(200, exampleResponse);

    client
      .getRecipeMealPlanTasks(recipeID)
      .then((response: APIResponse<RecipePrepTaskStep>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeMealPlanTasks', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<RecipePrepTaskStep>(buildObligatoryError('getRecipeMealPlanTasks'));
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_steps`).reply(200, exampleResponse);

    expect(client.getRecipeMealPlanTasks(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a RecipeMealPlanTasks', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<RecipePrepTaskStep>(buildObligatoryError('getRecipeMealPlanTasks'));
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_steps`).reply(500, exampleResponse);

    expect(client.getRecipeMealPlanTasks(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a RecipePrepTask', () => {
    let recipeID = fakeID();
    let recipePrepTaskID = fakeID();

    const exampleResponse = new APIResponse<RecipePrepTask>();
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(200, exampleResponse);

    client
      .getRecipePrepTask(recipeID, recipePrepTaskID)
      .then((response: APIResponse<RecipePrepTask>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipePrepTask', () => {
    let recipeID = fakeID();
    let recipePrepTaskID = fakeID();

    const exampleResponse = new APIResponse<RecipePrepTask>(buildObligatoryError('getRecipePrepTask'));
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(200, exampleResponse);

    expect(client.getRecipePrepTask(recipeID, recipePrepTaskID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipePrepTask', () => {
    let recipeID = fakeID();
    let recipePrepTaskID = fakeID();

    const exampleResponse = new APIResponse<RecipePrepTask>(buildObligatoryError('getRecipePrepTask'));
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(500, exampleResponse);

    expect(client.getRecipePrepTask(recipeID, recipePrepTaskID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a RecipePrepTasks', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipePrepTask>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new RecipePrepTask()],
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks`).reply(200, exampleResponse);

    client
      .getRecipePrepTasks(recipeID)
      .then((response: QueryFilteredResult<RecipePrepTask>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipePrepTasks', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipePrepTask>>(
      buildObligatoryError('getRecipePrepTasks user error'),
    );
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks`).reply(200, exampleResponse);

    expect(client.getRecipePrepTasks(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a RecipePrepTasks', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipePrepTask>>(
      buildObligatoryError('getRecipePrepTasks service error'),
    );
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks`).reply(500, exampleResponse);

    expect(client.getRecipePrepTasks(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a RecipeRating', () => {
    let recipeID = fakeID();
    let recipeRatingID = fakeID();

    const exampleResponse = new APIResponse<RecipeRating>();
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(200, exampleResponse);

    client
      .getRecipeRating(recipeID, recipeRatingID)
      .then((response: APIResponse<RecipeRating>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeRating', () => {
    let recipeID = fakeID();
    let recipeRatingID = fakeID();

    const exampleResponse = new APIResponse<RecipeRating>(buildObligatoryError('getRecipeRating'));
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(200, exampleResponse);

    expect(client.getRecipeRating(recipeID, recipeRatingID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a RecipeRating', () => {
    let recipeID = fakeID();
    let recipeRatingID = fakeID();

    const exampleResponse = new APIResponse<RecipeRating>(buildObligatoryError('getRecipeRating'));
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(500, exampleResponse);

    expect(client.getRecipeRating(recipeID, recipeRatingID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a RecipeRatings', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeRating>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new RecipeRating()],
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings`).reply(200, exampleResponse);

    client
      .getRecipeRatingsForRecipe(recipeID)
      .then((response: QueryFilteredResult<RecipeRating>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeRatings', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeRating>>(
      buildObligatoryError('getRecipeRatingsForRecipe user error'),
    );
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings`).reply(200, exampleResponse);

    expect(client.getRecipeRatingsForRecipe(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a RecipeRatings', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeRating>>(
      buildObligatoryError('getRecipeRatingsForRecipe service error'),
    );
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings`).reply(500, exampleResponse);

    expect(client.getRecipeRatingsForRecipe(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a RecipeStep', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<RecipeStep>();
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(200, exampleResponse);

    client
      .getRecipeStep(recipeID, recipeStepID)
      .then((response: APIResponse<RecipeStep>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeStep', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<RecipeStep>(buildObligatoryError('getRecipeStep'));
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(200, exampleResponse);

    expect(client.getRecipeStep(recipeID, recipeStepID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStep', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<RecipeStep>(buildObligatoryError('getRecipeStep'));
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(500, exampleResponse);

    expect(client.getRecipeStep(recipeID, recipeStepID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepCompletionConditionID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>();
    mock
      .onGet(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(200, exampleResponse);

    client
      .getRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID)
      .then((response: APIResponse<RecipeStepCompletionCondition>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepCompletionConditionID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(
      buildObligatoryError('getRecipeStepCompletionCondition'),
    );
    mock
      .onGet(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(200, exampleResponse);

    expect(
      client.getRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepCompletionConditionID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(
      buildObligatoryError('getRecipeStepCompletionCondition'),
    );
    mock
      .onGet(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(500, exampleResponse);

    expect(
      client.getRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a RecipeStepCompletionConditions', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepCompletionCondition>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new RecipeStepCompletionCondition()],
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`)
      .reply(200, exampleResponse);

    client
      .getRecipeStepCompletionConditions(recipeID, recipeStepID)
      .then((response: QueryFilteredResult<RecipeStepCompletionCondition>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeStepCompletionConditions', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepCompletionCondition>>(
      buildObligatoryError('getRecipeStepCompletionConditions user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`)
      .reply(200, exampleResponse);

    expect(client.getRecipeStepCompletionConditions(recipeID, recipeStepID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepCompletionConditions', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepCompletionCondition>>(
      buildObligatoryError('getRecipeStepCompletionConditions service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`)
      .reply(500, exampleResponse);

    expect(client.getRecipeStepCompletionConditions(recipeID, recipeStepID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepIngredientID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepIngredient>();
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(200, exampleResponse);

    client
      .getRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)
      .then((response: APIResponse<RecipeStepIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepIngredientID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepIngredient>(buildObligatoryError('getRecipeStepIngredient'));
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.getRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepIngredientID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepIngredient>(buildObligatoryError('getRecipeStepIngredient'));
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.getRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a RecipeStepIngredients', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepIngredient>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new RecipeStepIngredient()],
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`).reply(200, exampleResponse);

    client
      .getRecipeStepIngredients(recipeID, recipeStepID)
      .then((response: QueryFilteredResult<RecipeStepIngredient>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeStepIngredients', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepIngredient>>(
      buildObligatoryError('getRecipeStepIngredients user error'),
    );
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`).reply(200, exampleResponse);

    expect(client.getRecipeStepIngredients(recipeID, recipeStepID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepIngredients', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepIngredient>>(
      buildObligatoryError('getRecipeStepIngredients service error'),
    );
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`).reply(500, exampleResponse);

    expect(client.getRecipeStepIngredients(recipeID, recipeStepID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepInstrumentID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepInstrument>();
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(200, exampleResponse);

    client
      .getRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)
      .then((response: APIResponse<RecipeStepInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepInstrumentID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepInstrument>(buildObligatoryError('getRecipeStepInstrument'));
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(200, exampleResponse);

    expect(client.getRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepInstrumentID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepInstrument>(buildObligatoryError('getRecipeStepInstrument'));
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(500, exampleResponse);

    expect(client.getRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a RecipeStepInstruments', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepInstrument>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new RecipeStepInstrument()],
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`).reply(200, exampleResponse);

    client
      .getRecipeStepInstruments(recipeID, recipeStepID)
      .then((response: QueryFilteredResult<RecipeStepInstrument>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeStepInstruments', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepInstrument>>(
      buildObligatoryError('getRecipeStepInstruments user error'),
    );
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`).reply(200, exampleResponse);

    expect(client.getRecipeStepInstruments(recipeID, recipeStepID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepInstruments', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepInstrument>>(
      buildObligatoryError('getRecipeStepInstruments service error'),
    );
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`).reply(500, exampleResponse);

    expect(client.getRecipeStepInstruments(recipeID, recipeStepID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepProductID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepProduct>();
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(200, exampleResponse);

    client
      .getRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)
      .then((response: APIResponse<RecipeStepProduct>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepProductID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepProduct>(buildObligatoryError('getRecipeStepProduct'));
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(200, exampleResponse);

    expect(client.getRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepProductID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepProduct>(buildObligatoryError('getRecipeStepProduct'));
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(500, exampleResponse);

    expect(client.getRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a RecipeStepProducts', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepProduct>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new RecipeStepProduct()],
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`).reply(200, exampleResponse);

    client
      .getRecipeStepProducts(recipeID, recipeStepID)
      .then((response: QueryFilteredResult<RecipeStepProduct>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeStepProducts', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepProduct>>(
      buildObligatoryError('getRecipeStepProducts user error'),
    );
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`).reply(200, exampleResponse);

    expect(client.getRecipeStepProducts(recipeID, recipeStepID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepProducts', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepProduct>>(
      buildObligatoryError('getRecipeStepProducts service error'),
    );
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`).reply(500, exampleResponse);

    expect(client.getRecipeStepProducts(recipeID, recipeStepID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepVesselID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepVessel>();
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(200, exampleResponse);

    client
      .getRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)
      .then((response: APIResponse<RecipeStepVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepVesselID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepVessel>(buildObligatoryError('getRecipeStepVessel'));
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(200, exampleResponse);

    expect(client.getRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepVesselID = fakeID();

    const exampleResponse = new APIResponse<RecipeStepVessel>(buildObligatoryError('getRecipeStepVessel'));
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(500, exampleResponse);

    expect(client.getRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a RecipeStepVessels', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepVessel>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new RecipeStepVessel()],
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`).reply(200, exampleResponse);

    client
      .getRecipeStepVessels(recipeID, recipeStepID)
      .then((response: QueryFilteredResult<RecipeStepVessel>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeStepVessels', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepVessel>>(
      buildObligatoryError('getRecipeStepVessels user error'),
    );
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`).reply(200, exampleResponse);

    expect(client.getRecipeStepVessels(recipeID, recipeStepID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepVessels', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStepVessel>>(
      buildObligatoryError('getRecipeStepVessels service error'),
    );
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`).reply(500, exampleResponse);

    expect(client.getRecipeStepVessels(recipeID, recipeStepID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a RecipeSteps', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStep>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new RecipeStep()],
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps`).reply(200, exampleResponse);

    client
      .getRecipeSteps(recipeID)
      .then((response: QueryFilteredResult<RecipeStep>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a RecipeSteps', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStep>>(buildObligatoryError('getRecipeSteps user error'));
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps`).reply(200, exampleResponse);

    expect(client.getRecipeSteps(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a RecipeSteps', () => {
    let recipeID = fakeID();

    const exampleResponse = new APIResponse<Array<RecipeStep>>(buildObligatoryError('getRecipeSteps service error'));
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps`).reply(500, exampleResponse);

    expect(client.getRecipeSteps(recipeID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a Recipes', () => {
    const exampleResponse = new APIResponse<Array<Recipe>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new Recipe()],
    });
    mock.onGet(`${baseURL}/api/v1/recipes`).reply(200, exampleResponse);

    client
      .getRecipes()
      .then((response: QueryFilteredResult<Recipe>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a Recipes', () => {
    const exampleResponse = new APIResponse<Array<Recipe>>(buildObligatoryError('getRecipes user error'));
    mock.onGet(`${baseURL}/api/v1/recipes`).reply(200, exampleResponse);

    expect(client.getRecipes()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a Recipes', () => {
    const exampleResponse = new APIResponse<Array<Recipe>>(buildObligatoryError('getRecipes service error'));
    mock.onGet(`${baseURL}/api/v1/recipes`).reply(500, exampleResponse);

    expect(client.getRecipes()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a Self', () => {
    const exampleResponse = new APIResponse<User>();
    mock.onGet(`${baseURL}/api/v1/users/self`).reply(200, exampleResponse);

    client
      .getSelf()
      .then((response: APIResponse<User>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a Self', () => {
    const exampleResponse = new APIResponse<User>(buildObligatoryError('getSelf'));
    mock.onGet(`${baseURL}/api/v1/users/self`).reply(200, exampleResponse);

    expect(client.getSelf()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a Self', () => {
    const exampleResponse = new APIResponse<User>(buildObligatoryError('getSelf'));
    mock.onGet(`${baseURL}/api/v1/users/self`).reply(500, exampleResponse);

    expect(client.getSelf()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ServiceSetting', () => {
    let serviceSettingID = fakeID();

    const exampleResponse = new APIResponse<ServiceSetting>();
    mock.onGet(`${baseURL}/api/v1/settings/${serviceSettingID}`).reply(200, exampleResponse);

    client
      .getServiceSetting(serviceSettingID)
      .then((response: APIResponse<ServiceSetting>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ServiceSetting', () => {
    let serviceSettingID = fakeID();

    const exampleResponse = new APIResponse<ServiceSetting>(buildObligatoryError('getServiceSetting'));
    mock.onGet(`${baseURL}/api/v1/settings/${serviceSettingID}`).reply(200, exampleResponse);

    expect(client.getServiceSetting(serviceSettingID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ServiceSetting', () => {
    let serviceSettingID = fakeID();

    const exampleResponse = new APIResponse<ServiceSetting>(buildObligatoryError('getServiceSetting'));
    mock.onGet(`${baseURL}/api/v1/settings/${serviceSettingID}`).reply(500, exampleResponse);

    expect(client.getServiceSetting(serviceSettingID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ServiceSettingConfigurationByName', () => {
    let serviceSettingConfigurationName = fakeID();

    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ServiceSettingConfiguration()],
    });
    mock
      .onGet(`${baseURL}/api/v1/settings/configurations/user/${serviceSettingConfigurationName}`)
      .reply(200, exampleResponse);

    client
      .getServiceSettingConfigurationByName(serviceSettingConfigurationName)
      .then((response: QueryFilteredResult<ServiceSettingConfiguration>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ServiceSettingConfigurationByName', () => {
    let serviceSettingConfigurationName = fakeID();

    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>(
      buildObligatoryError('getServiceSettingConfigurationByName user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/settings/configurations/user/${serviceSettingConfigurationName}`)
      .reply(200, exampleResponse);

    expect(client.getServiceSettingConfigurationByName(serviceSettingConfigurationName)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ServiceSettingConfigurationByName', () => {
    let serviceSettingConfigurationName = fakeID();

    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>(
      buildObligatoryError('getServiceSettingConfigurationByName service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/settings/configurations/user/${serviceSettingConfigurationName}`)
      .reply(500, exampleResponse);

    expect(client.getServiceSettingConfigurationByName(serviceSettingConfigurationName)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ServiceSettingConfigurationsForHousehold', () => {
    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ServiceSettingConfiguration()],
    });
    mock.onGet(`${baseURL}/api/v1/settings/configurations/household`).reply(200, exampleResponse);

    client
      .getServiceSettingConfigurationsForHousehold()
      .then((response: QueryFilteredResult<ServiceSettingConfiguration>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ServiceSettingConfigurationsForHousehold', () => {
    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>(
      buildObligatoryError('getServiceSettingConfigurationsForHousehold user error'),
    );
    mock.onGet(`${baseURL}/api/v1/settings/configurations/household`).reply(200, exampleResponse);

    expect(client.getServiceSettingConfigurationsForHousehold()).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ServiceSettingConfigurationsForHousehold', () => {
    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>(
      buildObligatoryError('getServiceSettingConfigurationsForHousehold service error'),
    );
    mock.onGet(`${baseURL}/api/v1/settings/configurations/household`).reply(500, exampleResponse);

    expect(client.getServiceSettingConfigurationsForHousehold()).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ServiceSettingConfigurationsForUser', () => {
    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ServiceSettingConfiguration()],
    });
    mock.onGet(`${baseURL}/api/v1/settings/configurations/user`).reply(200, exampleResponse);

    client
      .getServiceSettingConfigurationsForUser()
      .then((response: QueryFilteredResult<ServiceSettingConfiguration>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ServiceSettingConfigurationsForUser', () => {
    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>(
      buildObligatoryError('getServiceSettingConfigurationsForUser user error'),
    );
    mock.onGet(`${baseURL}/api/v1/settings/configurations/user`).reply(200, exampleResponse);

    expect(client.getServiceSettingConfigurationsForUser()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ServiceSettingConfigurationsForUser', () => {
    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>(
      buildObligatoryError('getServiceSettingConfigurationsForUser service error'),
    );
    mock.onGet(`${baseURL}/api/v1/settings/configurations/user`).reply(500, exampleResponse);

    expect(client.getServiceSettingConfigurationsForUser()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ServiceSettings', () => {
    const exampleResponse = new APIResponse<Array<ServiceSetting>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ServiceSetting()],
    });
    mock.onGet(`${baseURL}/api/v1/settings`).reply(200, exampleResponse);

    client
      .getServiceSettings()
      .then((response: QueryFilteredResult<ServiceSetting>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ServiceSettings', () => {
    const exampleResponse = new APIResponse<Array<ServiceSetting>>(
      buildObligatoryError('getServiceSettings user error'),
    );
    mock.onGet(`${baseURL}/api/v1/settings`).reply(200, exampleResponse);

    expect(client.getServiceSettings()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ServiceSettings', () => {
    const exampleResponse = new APIResponse<Array<ServiceSetting>>(
      buildObligatoryError('getServiceSettings service error'),
    );
    mock.onGet(`${baseURL}/api/v1/settings`).reply(500, exampleResponse);

    expect(client.getServiceSettings()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a User', () => {
    let userID = fakeID();

    const exampleResponse = new APIResponse<User>();
    mock.onGet(`${baseURL}/api/v1/users/${userID}`).reply(200, exampleResponse);

    client
      .getUser(userID)
      .then((response: APIResponse<User>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a User', () => {
    let userID = fakeID();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('getUser'));
    mock.onGet(`${baseURL}/api/v1/users/${userID}`).reply(200, exampleResponse);

    expect(client.getUser(userID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a User', () => {
    let userID = fakeID();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('getUser'));
    mock.onGet(`${baseURL}/api/v1/users/${userID}`).reply(500, exampleResponse);

    expect(client.getUser(userID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a UserIngredientPreferences', () => {
    const exampleResponse = new APIResponse<Array<UserIngredientPreference>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new UserIngredientPreference()],
    });
    mock.onGet(`${baseURL}/api/v1/user_ingredient_preferences`).reply(200, exampleResponse);

    client
      .getUserIngredientPreferences()
      .then((response: QueryFilteredResult<UserIngredientPreference>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a UserIngredientPreferences', () => {
    const exampleResponse = new APIResponse<Array<UserIngredientPreference>>(
      buildObligatoryError('getUserIngredientPreferences user error'),
    );
    mock.onGet(`${baseURL}/api/v1/user_ingredient_preferences`).reply(200, exampleResponse);

    expect(client.getUserIngredientPreferences()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a UserIngredientPreferences', () => {
    const exampleResponse = new APIResponse<Array<UserIngredientPreference>>(
      buildObligatoryError('getUserIngredientPreferences service error'),
    );
    mock.onGet(`${baseURL}/api/v1/user_ingredient_preferences`).reply(500, exampleResponse);

    expect(client.getUserIngredientPreferences()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a UserNotification', () => {
    let userNotificationID = fakeID();

    const exampleResponse = new APIResponse<UserNotification>();
    mock.onGet(`${baseURL}/api/v1/user_notifications/${userNotificationID}`).reply(200, exampleResponse);

    client
      .getUserNotification(userNotificationID)
      .then((response: APIResponse<UserNotification>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a UserNotification', () => {
    let userNotificationID = fakeID();

    const exampleResponse = new APIResponse<UserNotification>(buildObligatoryError('getUserNotification'));
    mock.onGet(`${baseURL}/api/v1/user_notifications/${userNotificationID}`).reply(200, exampleResponse);

    expect(client.getUserNotification(userNotificationID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a UserNotification', () => {
    let userNotificationID = fakeID();

    const exampleResponse = new APIResponse<UserNotification>(buildObligatoryError('getUserNotification'));
    mock.onGet(`${baseURL}/api/v1/user_notifications/${userNotificationID}`).reply(500, exampleResponse);

    expect(client.getUserNotification(userNotificationID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a UserNotifications', () => {
    const exampleResponse = new APIResponse<Array<UserNotification>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new UserNotification()],
    });
    mock.onGet(`${baseURL}/api/v1/user_notifications`).reply(200, exampleResponse);

    client
      .getUserNotifications()
      .then((response: QueryFilteredResult<UserNotification>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a UserNotifications', () => {
    const exampleResponse = new APIResponse<Array<UserNotification>>(
      buildObligatoryError('getUserNotifications user error'),
    );
    mock.onGet(`${baseURL}/api/v1/user_notifications`).reply(200, exampleResponse);

    expect(client.getUserNotifications()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a UserNotifications', () => {
    const exampleResponse = new APIResponse<Array<UserNotification>>(
      buildObligatoryError('getUserNotifications service error'),
    );
    mock.onGet(`${baseURL}/api/v1/user_notifications`).reply(500, exampleResponse);

    expect(client.getUserNotifications()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a Users', () => {
    const exampleResponse = new APIResponse<Array<User>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new User()],
    });
    mock.onGet(`${baseURL}/api/v1/users`).reply(200, exampleResponse);

    client
      .getUsers()
      .then((response: QueryFilteredResult<User>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a Users', () => {
    const exampleResponse = new APIResponse<Array<User>>(buildObligatoryError('getUsers user error'));
    mock.onGet(`${baseURL}/api/v1/users`).reply(200, exampleResponse);

    expect(client.getUsers()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a Users', () => {
    const exampleResponse = new APIResponse<Array<User>>(buildObligatoryError('getUsers service error'));
    mock.onGet(`${baseURL}/api/v1/users`).reply(500, exampleResponse);

    expect(client.getUsers()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredient>();
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(200, exampleResponse);

    client
      .getValidIngredient(validIngredientID)
      .then((response: APIResponse<ValidIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredient>(buildObligatoryError('getValidIngredient'));
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(200, exampleResponse);

    expect(client.getValidIngredient(validIngredientID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredient>(buildObligatoryError('getValidIngredient'));
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(500, exampleResponse);

    expect(client.getValidIngredient(validIngredientID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidIngredientGroup', () => {
    let validIngredientGroupID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientGroup>();
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(200, exampleResponse);

    client
      .getValidIngredientGroup(validIngredientGroupID)
      .then((response: APIResponse<ValidIngredientGroup>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientGroup', () => {
    let validIngredientGroupID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientGroup>(buildObligatoryError('getValidIngredientGroup'));
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(200, exampleResponse);

    expect(client.getValidIngredientGroup(validIngredientGroupID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientGroup', () => {
    let validIngredientGroupID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientGroup>(buildObligatoryError('getValidIngredientGroup'));
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(500, exampleResponse);

    expect(client.getValidIngredientGroup(validIngredientGroupID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidIngredientGroups', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientGroup>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientGroup()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups`).reply(200, exampleResponse);

    client
      .getValidIngredientGroups()
      .then((response: QueryFilteredResult<ValidIngredientGroup>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientGroups', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientGroup>>(
      buildObligatoryError('getValidIngredientGroups user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups`).reply(200, exampleResponse);

    expect(client.getValidIngredientGroups()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientGroups', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientGroup>>(
      buildObligatoryError('getValidIngredientGroups service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups`).reply(500, exampleResponse);

    expect(client.getValidIngredientGroups()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>();
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)
      .then((response: APIResponse<ValidIngredientMeasurementUnit>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(
      buildObligatoryError('getValidIngredientMeasurementUnit'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(
      buildObligatoryError('getValidIngredientMeasurementUnit'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidIngredientMeasurementUnits', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientMeasurementUnit()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units`).reply(200, exampleResponse);

    client
      .getValidIngredientMeasurementUnits()
      .then((response: QueryFilteredResult<ValidIngredientMeasurementUnit>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientMeasurementUnits', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>(
      buildObligatoryError('getValidIngredientMeasurementUnits user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units`).reply(200, exampleResponse);

    expect(client.getValidIngredientMeasurementUnits()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientMeasurementUnits', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>(
      buildObligatoryError('getValidIngredientMeasurementUnits service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units`).reply(500, exampleResponse);

    expect(client.getValidIngredientMeasurementUnits()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidIngredientMeasurementUnitsByIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientMeasurementUnit()],
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientMeasurementUnitsByIngredient(validIngredientID)
      .then((response: QueryFilteredResult<ValidIngredientMeasurementUnit>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientMeasurementUnitsByIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>(
      buildObligatoryError('getValidIngredientMeasurementUnitsByIngredient user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientMeasurementUnitsByIngredient(validIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientMeasurementUnitsByIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>(
      buildObligatoryError('getValidIngredientMeasurementUnitsByIngredient service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/by_ingredient/${validIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientMeasurementUnitsByIngredient(validIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidIngredientMeasurementUnitsByMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientMeasurementUnit()],
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/by_measurement_unit/${validMeasurementUnitID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientMeasurementUnitsByMeasurementUnit(validMeasurementUnitID)
      .then((response: QueryFilteredResult<ValidIngredientMeasurementUnit>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientMeasurementUnitsByMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>(
      buildObligatoryError('getValidIngredientMeasurementUnitsByMeasurementUnit user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/by_measurement_unit/${validMeasurementUnitID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientMeasurementUnitsByMeasurementUnit(validMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientMeasurementUnitsByMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>(
      buildObligatoryError('getValidIngredientMeasurementUnitsByMeasurementUnit service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/by_measurement_unit/${validMeasurementUnitID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientMeasurementUnitsByMeasurementUnit(validMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidIngredientPreparation', () => {
    let validIngredientPreparationID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientPreparation>();
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientPreparation(validIngredientPreparationID)
      .then((response: APIResponse<ValidIngredientPreparation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientPreparation', () => {
    let validIngredientPreparationID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientPreparation>(
      buildObligatoryError('getValidIngredientPreparation'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientPreparation(validIngredientPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientPreparation', () => {
    let validIngredientPreparationID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientPreparation>(
      buildObligatoryError('getValidIngredientPreparation'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientPreparation(validIngredientPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidIngredientPreparations', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientPreparation()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_preparations`).reply(200, exampleResponse);

    client
      .getValidIngredientPreparations()
      .then((response: QueryFilteredResult<ValidIngredientPreparation>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientPreparations', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>(
      buildObligatoryError('getValidIngredientPreparations user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_preparations`).reply(200, exampleResponse);

    expect(client.getValidIngredientPreparations()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientPreparations', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>(
      buildObligatoryError('getValidIngredientPreparations service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_preparations`).reply(500, exampleResponse);

    expect(client.getValidIngredientPreparations()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidIngredientPreparationsByIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientPreparation()],
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientPreparationsByIngredient(validIngredientID)
      .then((response: QueryFilteredResult<ValidIngredientPreparation>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientPreparationsByIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>(
      buildObligatoryError('getValidIngredientPreparationsByIngredient user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientPreparationsByIngredient(validIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientPreparationsByIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>(
      buildObligatoryError('getValidIngredientPreparationsByIngredient service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/by_ingredient/${validIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientPreparationsByIngredient(validIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidIngredientPreparationsByPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientPreparation()],
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/by_preparation/${validPreparationID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientPreparationsByPreparation(validPreparationID)
      .then((response: QueryFilteredResult<ValidIngredientPreparation>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientPreparationsByPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>(
      buildObligatoryError('getValidIngredientPreparationsByPreparation user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/by_preparation/${validPreparationID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientPreparationsByPreparation(validPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientPreparationsByPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>(
      buildObligatoryError('getValidIngredientPreparationsByPreparation service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/by_preparation/${validPreparationID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientPreparationsByPreparation(validPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientState>();
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(200, exampleResponse);

    client
      .getValidIngredientState(validIngredientStateID)
      .then((response: APIResponse<ValidIngredientState>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientState>(buildObligatoryError('getValidIngredientState'));
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(200, exampleResponse);

    expect(client.getValidIngredientState(validIngredientStateID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientState>(buildObligatoryError('getValidIngredientState'));
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(500, exampleResponse);

    expect(client.getValidIngredientState(validIngredientStateID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>();
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientStateIngredient(validIngredientStateIngredientID)
      .then((response: APIResponse<ValidIngredientStateIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(
      buildObligatoryError('getValidIngredientStateIngredient'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientStateIngredient(validIngredientStateIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = fakeID();

    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(
      buildObligatoryError('getValidIngredientStateIngredient'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientStateIngredient(validIngredientStateIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidIngredientStateIngredients', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientStateIngredient()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients`).reply(200, exampleResponse);

    client
      .getValidIngredientStateIngredients()
      .then((response: QueryFilteredResult<ValidIngredientStateIngredient>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientStateIngredients', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>(
      buildObligatoryError('getValidIngredientStateIngredients user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients`).reply(200, exampleResponse);

    expect(client.getValidIngredientStateIngredients()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientStateIngredients', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>(
      buildObligatoryError('getValidIngredientStateIngredients service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients`).reply(500, exampleResponse);

    expect(client.getValidIngredientStateIngredients()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidIngredientStateIngredientsByIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientStateIngredient()],
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientStateIngredientsByIngredient(validIngredientID)
      .then((response: QueryFilteredResult<ValidIngredientStateIngredient>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientStateIngredientsByIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>(
      buildObligatoryError('getValidIngredientStateIngredientsByIngredient user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientStateIngredientsByIngredient(validIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientStateIngredientsByIngredient', () => {
    let validIngredientID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>(
      buildObligatoryError('getValidIngredientStateIngredientsByIngredient service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/by_ingredient/${validIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientStateIngredientsByIngredient(validIngredientID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidIngredientStateIngredientsByIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientStateIngredient()],
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/${validIngredientStateID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientStateIngredientsByIngredientState(validIngredientStateID)
      .then((response: QueryFilteredResult<ValidIngredientStateIngredient>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientStateIngredientsByIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>(
      buildObligatoryError('getValidIngredientStateIngredientsByIngredientState user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/${validIngredientStateID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientStateIngredientsByIngredientState(validIngredientStateID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientStateIngredientsByIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>(
      buildObligatoryError('getValidIngredientStateIngredientsByIngredientState service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/${validIngredientStateID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientStateIngredientsByIngredientState(validIngredientStateID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidIngredientStates', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientState>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredientState()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states`).reply(200, exampleResponse);

    client
      .getValidIngredientStates()
      .then((response: QueryFilteredResult<ValidIngredientState>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredientStates', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientState>>(
      buildObligatoryError('getValidIngredientStates user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states`).reply(200, exampleResponse);

    expect(client.getValidIngredientStates()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientStates', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredientState>>(
      buildObligatoryError('getValidIngredientStates service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states`).reply(500, exampleResponse);

    expect(client.getValidIngredientStates()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidIngredients', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredient>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidIngredient()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredients`).reply(200, exampleResponse);

    client
      .getValidIngredients()
      .then((response: QueryFilteredResult<ValidIngredient>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidIngredients', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredient>>(
      buildObligatoryError('getValidIngredients user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredients`).reply(200, exampleResponse);

    expect(client.getValidIngredients()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredients', () => {
    const exampleResponse = new APIResponse<Array<ValidIngredient>>(
      buildObligatoryError('getValidIngredients service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_ingredients`).reply(500, exampleResponse);

    expect(client.getValidIngredients()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleResponse = new APIResponse<ValidInstrument>();
    mock.onGet(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(200, exampleResponse);

    client
      .getValidInstrument(validInstrumentID)
      .then((response: APIResponse<ValidInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleResponse = new APIResponse<ValidInstrument>(buildObligatoryError('getValidInstrument'));
    mock.onGet(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(200, exampleResponse);

    expect(client.getValidInstrument(validInstrumentID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleResponse = new APIResponse<ValidInstrument>(buildObligatoryError('getValidInstrument'));
    mock.onGet(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(500, exampleResponse);

    expect(client.getValidInstrument(validInstrumentID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidInstruments', () => {
    const exampleResponse = new APIResponse<Array<ValidInstrument>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidInstrument()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_instruments`).reply(200, exampleResponse);

    client
      .getValidInstruments()
      .then((response: QueryFilteredResult<ValidInstrument>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidInstruments', () => {
    const exampleResponse = new APIResponse<Array<ValidInstrument>>(
      buildObligatoryError('getValidInstruments user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_instruments`).reply(200, exampleResponse);

    expect(client.getValidInstruments()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidInstruments', () => {
    const exampleResponse = new APIResponse<Array<ValidInstrument>>(
      buildObligatoryError('getValidInstruments service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_instruments`).reply(500, exampleResponse);

    expect(client.getValidInstruments()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<ValidMeasurementUnit>();
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(200, exampleResponse);

    client
      .getValidMeasurementUnit(validMeasurementUnitID)
      .then((response: APIResponse<ValidMeasurementUnit>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<ValidMeasurementUnit>(buildObligatoryError('getValidMeasurementUnit'));
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(200, exampleResponse);

    expect(client.getValidMeasurementUnit(validMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<ValidMeasurementUnit>(buildObligatoryError('getValidMeasurementUnit'));
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(500, exampleResponse);

    expect(client.getValidMeasurementUnit(validMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = fakeID();

    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>();
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(200, exampleResponse);

    client
      .getValidMeasurementUnitConversion(validMeasurementUnitConversionID)
      .then((response: APIResponse<ValidMeasurementUnitConversion>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = fakeID();

    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(
      buildObligatoryError('getValidMeasurementUnitConversion'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(200, exampleResponse);

    expect(client.getValidMeasurementUnitConversion(validMeasurementUnitConversionID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = fakeID();

    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(
      buildObligatoryError('getValidMeasurementUnitConversion'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(500, exampleResponse);

    expect(client.getValidMeasurementUnitConversion(validMeasurementUnitConversionID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidMeasurementUnitConversionsFromUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidMeasurementUnitConversion>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidMeasurementUnitConversion()],
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/from_unit/${validMeasurementUnitID}`)
      .reply(200, exampleResponse);

    client
      .getValidMeasurementUnitConversionsFromUnit(validMeasurementUnitID)
      .then((response: QueryFilteredResult<ValidMeasurementUnitConversion>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidMeasurementUnitConversionsFromUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidMeasurementUnitConversion>>(
      buildObligatoryError('getValidMeasurementUnitConversionsFromUnit user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/from_unit/${validMeasurementUnitID}`)
      .reply(200, exampleResponse);

    expect(client.getValidMeasurementUnitConversionsFromUnit(validMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidMeasurementUnitConversionsFromUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidMeasurementUnitConversion>>(
      buildObligatoryError('getValidMeasurementUnitConversionsFromUnit service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/from_unit/${validMeasurementUnitID}`)
      .reply(500, exampleResponse);

    expect(client.getValidMeasurementUnitConversionsFromUnit(validMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidMeasurementUnitConversionsToUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidMeasurementUnitConversion>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidMeasurementUnitConversion()],
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/to_unit/${validMeasurementUnitID}`)
      .reply(200, exampleResponse);

    client
      .getValidMeasurementUnitConversionsToUnit(validMeasurementUnitID)
      .then((response: QueryFilteredResult<ValidMeasurementUnitConversion>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidMeasurementUnitConversionsToUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidMeasurementUnitConversion>>(
      buildObligatoryError('getValidMeasurementUnitConversionsToUnit user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/to_unit/${validMeasurementUnitID}`)
      .reply(200, exampleResponse);

    expect(client.getValidMeasurementUnitConversionsToUnit(validMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidMeasurementUnitConversionsToUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidMeasurementUnitConversion>>(
      buildObligatoryError('getValidMeasurementUnitConversionsToUnit service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/to_unit/${validMeasurementUnitID}`)
      .reply(500, exampleResponse);

    expect(client.getValidMeasurementUnitConversionsToUnit(validMeasurementUnitID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidMeasurementUnits', () => {
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidMeasurementUnit()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units`).reply(200, exampleResponse);

    client
      .getValidMeasurementUnits()
      .then((response: QueryFilteredResult<ValidMeasurementUnit>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidMeasurementUnits', () => {
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>(
      buildObligatoryError('getValidMeasurementUnits user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units`).reply(200, exampleResponse);

    expect(client.getValidMeasurementUnits()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidMeasurementUnits', () => {
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>(
      buildObligatoryError('getValidMeasurementUnits service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units`).reply(500, exampleResponse);

    expect(client.getValidMeasurementUnits()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparation>();
    mock.onGet(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(200, exampleResponse);

    client
      .getValidPreparation(validPreparationID)
      .then((response: APIResponse<ValidPreparation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparation>(buildObligatoryError('getValidPreparation'));
    mock.onGet(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(200, exampleResponse);

    expect(client.getValidPreparation(validPreparationID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparation>(buildObligatoryError('getValidPreparation'));
    mock.onGet(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(500, exampleResponse);

    expect(client.getValidPreparation(validPreparationID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidPreparationInstrument', () => {
    let validPreparationVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparationInstrument>();
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(200, exampleResponse);

    client
      .getValidPreparationInstrument(validPreparationVesselID)
      .then((response: APIResponse<ValidPreparationInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidPreparationInstrument', () => {
    let validPreparationVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparationInstrument>(
      buildObligatoryError('getValidPreparationInstrument'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(200, exampleResponse);

    expect(client.getValidPreparationInstrument(validPreparationVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationInstrument', () => {
    let validPreparationVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparationInstrument>(
      buildObligatoryError('getValidPreparationInstrument'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(500, exampleResponse);

    expect(client.getValidPreparationInstrument(validPreparationVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidPreparationInstruments', () => {
    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidPreparationInstrument()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparation_instruments`).reply(200, exampleResponse);

    client
      .getValidPreparationInstruments()
      .then((response: QueryFilteredResult<ValidPreparationInstrument>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidPreparationInstruments', () => {
    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>(
      buildObligatoryError('getValidPreparationInstruments user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_preparation_instruments`).reply(200, exampleResponse);

    expect(client.getValidPreparationInstruments()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationInstruments', () => {
    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>(
      buildObligatoryError('getValidPreparationInstruments service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_preparation_instruments`).reply(500, exampleResponse);

    expect(client.getValidPreparationInstruments()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidPreparationInstrumentsByInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidPreparationInstrument()],
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/by_instrument/${validInstrumentID}`)
      .reply(200, exampleResponse);

    client
      .getValidPreparationInstrumentsByInstrument(validInstrumentID)
      .then((response: QueryFilteredResult<ValidPreparationInstrument>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidPreparationInstrumentsByInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>(
      buildObligatoryError('getValidPreparationInstrumentsByInstrument user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/by_instrument/${validInstrumentID}`)
      .reply(200, exampleResponse);

    expect(client.getValidPreparationInstrumentsByInstrument(validInstrumentID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationInstrumentsByInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>(
      buildObligatoryError('getValidPreparationInstrumentsByInstrument service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/by_instrument/${validInstrumentID}`)
      .reply(500, exampleResponse);

    expect(client.getValidPreparationInstrumentsByInstrument(validInstrumentID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidPreparationInstrumentsByPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidPreparationInstrument()],
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/by_preparation/${validPreparationID}`)
      .reply(200, exampleResponse);

    client
      .getValidPreparationInstrumentsByPreparation(validPreparationID)
      .then((response: QueryFilteredResult<ValidPreparationInstrument>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidPreparationInstrumentsByPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>(
      buildObligatoryError('getValidPreparationInstrumentsByPreparation user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/by_preparation/${validPreparationID}`)
      .reply(200, exampleResponse);

    expect(client.getValidPreparationInstrumentsByPreparation(validPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationInstrumentsByPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>(
      buildObligatoryError('getValidPreparationInstrumentsByPreparation service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/by_preparation/${validPreparationID}`)
      .reply(500, exampleResponse);

    expect(client.getValidPreparationInstrumentsByPreparation(validPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidPreparationVessel', () => {
    let validPreparationVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparationVessel>();
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(200, exampleResponse);

    client
      .getValidPreparationVessel(validPreparationVesselID)
      .then((response: APIResponse<ValidPreparationVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidPreparationVessel', () => {
    let validPreparationVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparationVessel>(buildObligatoryError('getValidPreparationVessel'));
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(200, exampleResponse);

    expect(client.getValidPreparationVessel(validPreparationVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationVessel', () => {
    let validPreparationVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidPreparationVessel>(buildObligatoryError('getValidPreparationVessel'));
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(500, exampleResponse);

    expect(client.getValidPreparationVessel(validPreparationVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidPreparationVessels', () => {
    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidPreparationVessel()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels`).reply(200, exampleResponse);

    client
      .getValidPreparationVessels()
      .then((response: QueryFilteredResult<ValidPreparationVessel>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidPreparationVessels', () => {
    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>(
      buildObligatoryError('getValidPreparationVessels user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels`).reply(200, exampleResponse);

    expect(client.getValidPreparationVessels()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationVessels', () => {
    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>(
      buildObligatoryError('getValidPreparationVessels service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels`).reply(500, exampleResponse);

    expect(client.getValidPreparationVessels()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidPreparationVesselsByPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidPreparationVessel()],
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_vessels/by_preparation/${validPreparationID}`)
      .reply(200, exampleResponse);

    client
      .getValidPreparationVesselsByPreparation(validPreparationID)
      .then((response: QueryFilteredResult<ValidPreparationVessel>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidPreparationVesselsByPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>(
      buildObligatoryError('getValidPreparationVesselsByPreparation user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_vessels/by_preparation/${validPreparationID}`)
      .reply(200, exampleResponse);

    expect(client.getValidPreparationVesselsByPreparation(validPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationVesselsByPreparation', () => {
    let validPreparationID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>(
      buildObligatoryError('getValidPreparationVesselsByPreparation service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_vessels/by_preparation/${validPreparationID}`)
      .reply(500, exampleResponse);

    expect(client.getValidPreparationVesselsByPreparation(validPreparationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidPreparationVesselsByVessel', () => {
    let ValidVesselID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidPreparationVessel()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/by_vessel/${ValidVesselID}`).reply(200, exampleResponse);

    client
      .getValidPreparationVesselsByVessel(ValidVesselID)
      .then((response: QueryFilteredResult<ValidPreparationVessel>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidPreparationVesselsByVessel', () => {
    let ValidVesselID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>(
      buildObligatoryError('getValidPreparationVesselsByVessel user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/by_vessel/${ValidVesselID}`).reply(200, exampleResponse);

    expect(client.getValidPreparationVesselsByVessel(ValidVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationVesselsByVessel', () => {
    let ValidVesselID = fakeID();

    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>(
      buildObligatoryError('getValidPreparationVesselsByVessel service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/by_vessel/${ValidVesselID}`).reply(500, exampleResponse);

    expect(client.getValidPreparationVesselsByVessel(ValidVesselID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a ValidPreparations', () => {
    const exampleResponse = new APIResponse<Array<ValidPreparation>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidPreparation()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparations`).reply(200, exampleResponse);

    client
      .getValidPreparations()
      .then((response: QueryFilteredResult<ValidPreparation>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidPreparations', () => {
    const exampleResponse = new APIResponse<Array<ValidPreparation>>(
      buildObligatoryError('getValidPreparations user error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_preparations`).reply(200, exampleResponse);

    expect(client.getValidPreparations()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparations', () => {
    const exampleResponse = new APIResponse<Array<ValidPreparation>>(
      buildObligatoryError('getValidPreparations service error'),
    );
    mock.onGet(`${baseURL}/api/v1/valid_preparations`).reply(500, exampleResponse);

    expect(client.getValidPreparations()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidVessel', () => {
    let validVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidVessel>();
    mock.onGet(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(200, exampleResponse);

    client
      .getValidVessel(validVesselID)
      .then((response: APIResponse<ValidVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidVessel', () => {
    let validVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidVessel>(buildObligatoryError('getValidVessel'));
    mock.onGet(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(200, exampleResponse);

    expect(client.getValidVessel(validVesselID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidVessel', () => {
    let validVesselID = fakeID();

    const exampleResponse = new APIResponse<ValidVessel>(buildObligatoryError('getValidVessel'));
    mock.onGet(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(500, exampleResponse);

    expect(client.getValidVessel(validVesselID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a ValidVessels', () => {
    const exampleResponse = new APIResponse<Array<ValidVessel>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new ValidVessel()],
    });
    mock.onGet(`${baseURL}/api/v1/valid_vessels`).reply(200, exampleResponse);

    client
      .getValidVessels()
      .then((response: QueryFilteredResult<ValidVessel>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a ValidVessels', () => {
    const exampleResponse = new APIResponse<Array<ValidVessel>>(buildObligatoryError('getValidVessels user error'));
    mock.onGet(`${baseURL}/api/v1/valid_vessels`).reply(200, exampleResponse);

    expect(client.getValidVessels()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a ValidVessels', () => {
    const exampleResponse = new APIResponse<Array<ValidVessel>>(buildObligatoryError('getValidVessels service error'));
    mock.onGet(`${baseURL}/api/v1/valid_vessels`).reply(500, exampleResponse);

    expect(client.getValidVessels()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a Webhook', () => {
    let webhookID = fakeID();

    const exampleResponse = new APIResponse<Webhook>();
    mock.onGet(`${baseURL}/api/v1/webhooks/${webhookID}`).reply(200, exampleResponse);

    client
      .getWebhook(webhookID)
      .then((response: APIResponse<Webhook>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a Webhook', () => {
    let webhookID = fakeID();

    const exampleResponse = new APIResponse<Webhook>(buildObligatoryError('getWebhook'));
    mock.onGet(`${baseURL}/api/v1/webhooks/${webhookID}`).reply(200, exampleResponse);

    expect(client.getWebhook(webhookID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a Webhook', () => {
    let webhookID = fakeID();

    const exampleResponse = new APIResponse<Webhook>(buildObligatoryError('getWebhook'));
    mock.onGet(`${baseURL}/api/v1/webhooks/${webhookID}`).reply(500, exampleResponse);

    expect(client.getWebhook(webhookID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a Webhooks', () => {
    const exampleResponse = new APIResponse<Array<Webhook>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new Webhook()],
    });
    mock.onGet(`${baseURL}/api/v1/webhooks`).reply(200, exampleResponse);

    client
      .getWebhooks()
      .then((response: QueryFilteredResult<Webhook>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a Webhooks', () => {
    const exampleResponse = new APIResponse<Array<Webhook>>(buildObligatoryError('getWebhooks user error'));
    mock.onGet(`${baseURL}/api/v1/webhooks`).reply(200, exampleResponse);

    expect(client.getWebhooks()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a Webhooks', () => {
    const exampleResponse = new APIResponse<Array<Webhook>>(buildObligatoryError('getWebhooks service error'));
    mock.onGet(`${baseURL}/api/v1/webhooks`).reply(500, exampleResponse);

    expect(client.getWebhooks()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a household instrument ownership', () => {
    let householdInstrumentOwnershipID = fakeID();

    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>();
    mock
      .onGet(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(200, exampleResponse);

    client
      .getHouseholdInstrumentOwnership(householdInstrumentOwnershipID)
      .then((response: APIResponse<HouseholdInstrumentOwnership>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a household instrument ownership', () => {
    let householdInstrumentOwnershipID = fakeID();

    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(
      buildObligatoryError('getHouseholdInstrumentOwnership'),
    );
    mock
      .onGet(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(200, exampleResponse);

    expect(client.getHouseholdInstrumentOwnership(householdInstrumentOwnershipID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a household instrument ownership', () => {
    let householdInstrumentOwnershipID = fakeID();

    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(
      buildObligatoryError('getHouseholdInstrumentOwnership'),
    );
    mock
      .onGet(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(500, exampleResponse);

    expect(client.getHouseholdInstrumentOwnership(householdInstrumentOwnershipID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a household invitation by its ID', () => {
    let householdID = fakeID();
    let householdInvitationID = fakeID();

    const exampleResponse = new APIResponse<HouseholdInvitation>();
    mock
      .onGet(`${baseURL}/api/v1/households/${householdID}/invitations/${householdInvitationID}`)
      .reply(200, exampleResponse);

    client
      .getHouseholdInvitationByID(householdID, householdInvitationID)
      .then((response: APIResponse<HouseholdInvitation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a household invitation by its ID', () => {
    let householdID = fakeID();
    let householdInvitationID = fakeID();

    const exampleResponse = new APIResponse<HouseholdInvitation>(buildObligatoryError('getHouseholdInvitationByID'));
    mock
      .onGet(`${baseURL}/api/v1/households/${householdID}/invitations/${householdInvitationID}`)
      .reply(200, exampleResponse);

    expect(client.getHouseholdInvitationByID(householdID, householdInvitationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a household invitation by its ID', () => {
    let householdID = fakeID();
    let householdInvitationID = fakeID();

    const exampleResponse = new APIResponse<HouseholdInvitation>(buildObligatoryError('getHouseholdInvitationByID'));
    mock
      .onGet(`${baseURL}/api/v1/households/${householdID}/invitations/${householdInvitationID}`)
      .reply(500, exampleResponse);

    expect(client.getHouseholdInvitationByID(householdID, householdInvitationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a household', () => {
    let householdID = fakeID();

    const exampleResponse = new APIResponse<Household>();
    mock.onGet(`${baseURL}/api/v1/households/${householdID}`).reply(200, exampleResponse);

    client
      .getHousehold(householdID)
      .then((response: APIResponse<Household>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a household', () => {
    let householdID = fakeID();

    const exampleResponse = new APIResponse<Household>(buildObligatoryError('getHousehold'));
    mock.onGet(`${baseURL}/api/v1/households/${householdID}`).reply(200, exampleResponse);

    expect(client.getHousehold(householdID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a household', () => {
    let householdID = fakeID();

    const exampleResponse = new APIResponse<Household>(buildObligatoryError('getHousehold'));
    mock.onGet(`${baseURL}/api/v1/households/${householdID}`).reply(500, exampleResponse);

    expect(client.getHousehold(householdID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a list of households', () => {
    const exampleResponse = new APIResponse<Array<Household>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new Household()],
    });
    mock.onGet(`${baseURL}/api/v1/households`).reply(200, exampleResponse);

    client
      .getHouseholds()
      .then((response: QueryFilteredResult<Household>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a list of households', () => {
    const exampleResponse = new APIResponse<Array<Household>>(buildObligatoryError('getHouseholds user error'));
    mock.onGet(`${baseURL}/api/v1/households`).reply(200, exampleResponse);

    expect(client.getHouseholds()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a list of households', () => {
    const exampleResponse = new APIResponse<Array<Household>>(buildObligatoryError('getHouseholds service error'));
    mock.onGet(`${baseURL}/api/v1/households`).reply(500, exampleResponse);

    expect(client.getHouseholds()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch a meal plan option by its ID', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleResponse = new APIResponse<MealPlanOption>();
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(200, exampleResponse);

    client
      .getMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)
      .then((response: APIResponse<MealPlanOption>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a meal plan option by its ID', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleResponse = new APIResponse<MealPlanOption>(buildObligatoryError('getMealPlanOption'));
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(200, exampleResponse);

    expect(client.getMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch a meal plan option by its ID', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleResponse = new APIResponse<MealPlanOption>(buildObligatoryError('getMealPlanOption'));
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(500, exampleResponse);

    expect(client.getMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch a meal plan option vote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();
    let mealPlanOptionVoteID = fakeID();

    const exampleResponse = new APIResponse<MealPlanOptionVote>();
    mock
      .onGet(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(200, exampleResponse);

    client
      .getMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
      .then((response: APIResponse<MealPlanOptionVote>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch a meal plan option vote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();
    let mealPlanOptionVoteID = fakeID();

    const exampleResponse = new APIResponse<MealPlanOptionVote>(buildObligatoryError('getMealPlanOptionVote'));
    mock
      .onGet(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(200, exampleResponse);

    expect(
      client.getMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch a meal plan option vote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();
    let mealPlanOptionVoteID = fakeID();

    const exampleResponse = new APIResponse<MealPlanOptionVote>(buildObligatoryError('getMealPlanOptionVote'));
    mock
      .onGet(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(500, exampleResponse);

    expect(
      client.getMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it("should fetch a meal plan option's votes", () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanOptionVote>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new MealPlanOptionVote()],
    });
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes`)
      .reply(200, exampleResponse);

    client
      .getMealPlanOptionVotes(mealPlanID, mealPlanEventID, mealPlanOptionID)
      .then((response: QueryFilteredResult<MealPlanOptionVote>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it("should raise errors appropriately when trying to fetch a meal plan option's votes", () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanOptionVote>>(
      buildObligatoryError('getMealPlanOptionVotes user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes`)
      .reply(200, exampleResponse);

    expect(client.getMealPlanOptionVotes(mealPlanID, mealPlanEventID, mealPlanOptionID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it("should raise service errors appropriately when trying to fetch a meal plan option's votes", () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanOptionVote>>(
      buildObligatoryError('getMealPlanOptionVotes service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes`)
      .reply(500, exampleResponse);

    expect(client.getMealPlanOptionVotes(mealPlanID, mealPlanEventID, mealPlanOptionID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it("should fetch a meal plan's options", () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanOption>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new MealPlanOption()],
    });
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`)
      .reply(200, exampleResponse);

    client
      .getMealPlanOptions(mealPlanID, mealPlanEventID)
      .then((response: QueryFilteredResult<MealPlanOption>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it("should raise errors appropriately when trying to fetch a meal plan's options", () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanOption>>(
      buildObligatoryError('getMealPlanOptions user error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`)
      .reply(200, exampleResponse);

    expect(client.getMealPlanOptions(mealPlanID, mealPlanEventID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it("should raise service errors appropriately when trying to fetch a meal plan's options", () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleResponse = new APIResponse<Array<MealPlanOption>>(
      buildObligatoryError('getMealPlanOptions service error'),
    );
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`)
      .reply(500, exampleResponse);

    expect(client.getMealPlanOptions(mealPlanID, mealPlanEventID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch household invitations for a given household', () => {
    let householdInvitationID = fakeID();

    const exampleResponse = new APIResponse<HouseholdInvitation>();
    mock.onGet(`${baseURL}/api/v1/household_invitations/${householdInvitationID}`).reply(200, exampleResponse);

    client
      .getHouseholdInvitation(householdInvitationID)
      .then((response: APIResponse<HouseholdInvitation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch household invitations for a given household', () => {
    let householdInvitationID = fakeID();

    const exampleResponse = new APIResponse<HouseholdInvitation>(buildObligatoryError('getHouseholdInvitation'));
    mock.onGet(`${baseURL}/api/v1/household_invitations/${householdInvitationID}`).reply(200, exampleResponse);

    expect(client.getHouseholdInvitation(householdInvitationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should raise service errors appropriately when trying to fetch household invitations for a given household', () => {
    let householdInvitationID = fakeID();

    const exampleResponse = new APIResponse<HouseholdInvitation>(buildObligatoryError('getHouseholdInvitation'));
    mock.onGet(`${baseURL}/api/v1/household_invitations/${householdInvitationID}`).reply(500, exampleResponse);

    expect(client.getHouseholdInvitation(householdInvitationID)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should fetch received household invitations', () => {
    const exampleResponse = new APIResponse<Array<HouseholdInvitation>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new HouseholdInvitation()],
    });
    mock.onGet(`${baseURL}/api/v1/household_invitations/received`).reply(200, exampleResponse);

    client
      .getReceivedHouseholdInvitations()
      .then((response: QueryFilteredResult<HouseholdInvitation>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch received household invitations', () => {
    const exampleResponse = new APIResponse<Array<HouseholdInvitation>>(
      buildObligatoryError('getReceivedHouseholdInvitations user error'),
    );
    mock.onGet(`${baseURL}/api/v1/household_invitations/received`).reply(200, exampleResponse);

    expect(client.getReceivedHouseholdInvitations()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch received household invitations', () => {
    const exampleResponse = new APIResponse<Array<HouseholdInvitation>>(
      buildObligatoryError('getReceivedHouseholdInvitations service error'),
    );
    mock.onGet(`${baseURL}/api/v1/household_invitations/received`).reply(500, exampleResponse);

    expect(client.getReceivedHouseholdInvitations()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch sent household invitations', () => {
    const exampleResponse = new APIResponse<Array<HouseholdInvitation>>({
      details: {
        currentHouseholdID: 'test',
        traceID: 'test',
      },
      pagination: QueryFilter.Default().toPagination(),
      data: [new HouseholdInvitation()],
    });
    mock.onGet(`${baseURL}/api/v1/household_invitations/sent`).reply(200, exampleResponse);

    client
      .getSentHouseholdInvitations()
      .then((response: QueryFilteredResult<HouseholdInvitation>) => {
        expect(response.data).toEqual(exampleResponse.data);
        expect(response.page).toEqual(exampleResponse.pagination?.page);
        expect(response.limit).toEqual(exampleResponse.pagination?.limit);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch sent household invitations', () => {
    const exampleResponse = new APIResponse<Array<HouseholdInvitation>>(
      buildObligatoryError('getSentHouseholdInvitations user error'),
    );
    mock.onGet(`${baseURL}/api/v1/household_invitations/sent`).reply(200, exampleResponse);

    expect(client.getSentHouseholdInvitations()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch sent household invitations', () => {
    const exampleResponse = new APIResponse<Array<HouseholdInvitation>>(
      buildObligatoryError('getSentHouseholdInvitations service error'),
    );
    mock.onGet(`${baseURL}/api/v1/household_invitations/sent`).reply(500, exampleResponse);

    expect(client.getSentHouseholdInvitations()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should fetch the currently active household', () => {
    const exampleResponse = new APIResponse<Household>();
    mock.onGet(`${baseURL}/api/v1/households/current`).reply(200, exampleResponse);

    client
      .getActiveHousehold()
      .then((response: APIResponse<Household>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should raise errors appropriately when trying to fetch the currently active household', () => {
    const exampleResponse = new APIResponse<Household>(buildObligatoryError('getActiveHousehold'));
    mock.onGet(`${baseURL}/api/v1/households/current`).reply(200, exampleResponse);

    expect(client.getActiveHousehold()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should raise service errors appropriately when trying to fetch the currently active household', () => {
    const exampleResponse = new APIResponse<Household>(buildObligatoryError('getActiveHousehold'));
    mock.onGet(`${baseURL}/api/v1/households/current`).reply(500, exampleResponse);

    expect(client.getActiveHousehold()).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should reject a received household invitation', () => {
    let householdInvitationID = fakeID();

    const exampleInput = new HouseholdInvitationUpdateRequestInput();

    const exampleResponse = new APIResponse<HouseholdInvitation>();
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/reject`).reply(200, exampleResponse);

    client
      .rejectHouseholdInvitation(householdInvitationID, exampleInput)
      .then((response: APIResponse<HouseholdInvitation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during reject a received household invitation', () => {
    let householdInvitationID = fakeID();

    const exampleInput = new HouseholdInvitationUpdateRequestInput();

    const exampleResponse = new APIResponse<HouseholdInvitation>(
      buildObligatoryError('rejectHouseholdInvitation user error'),
    );
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/reject`).reply(200, exampleResponse);

    expect(client.rejectHouseholdInvitation(householdInvitationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during reject a received household invitation', () => {
    let householdInvitationID = fakeID();

    const exampleInput = new HouseholdInvitationUpdateRequestInput();

    const exampleResponse = new APIResponse<HouseholdInvitation>(
      buildObligatoryError('rejectHouseholdInvitation service error'),
    );
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/reject`).reply(500, exampleResponse);

    expect(client.rejectHouseholdInvitation(householdInvitationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a MealPlan', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlan>();
    mock.onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(200, exampleResponse);

    client
      .updateMealPlan(mealPlanID, exampleInput)
      .then((response: APIResponse<MealPlan>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a MealPlan', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlan>(buildObligatoryError('updateMealPlan user error'));
    mock.onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(200, exampleResponse);

    expect(client.updateMealPlan(mealPlanID, exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during update a MealPlan', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlan>(buildObligatoryError('updateMealPlan service error'));
    mock.onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(500, exampleResponse);

    expect(client.updateMealPlan(mealPlanID, exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should update a MealPlanEvent', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleInput = new MealPlanEventUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlanEvent>();
    mock.onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(200, exampleResponse);

    client
      .updateMealPlanEvent(mealPlanID, mealPlanEventID, exampleInput)
      .then((response: APIResponse<MealPlanEvent>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a MealPlanEvent', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleInput = new MealPlanEventUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlanEvent>(buildObligatoryError('updateMealPlanEvent user error'));
    mock.onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(200, exampleResponse);

    expect(client.updateMealPlanEvent(mealPlanID, mealPlanEventID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a MealPlanEvent', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleInput = new MealPlanEventUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlanEvent>(buildObligatoryError('updateMealPlanEvent service error'));
    mock.onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(500, exampleResponse);

    expect(client.updateMealPlanEvent(mealPlanID, mealPlanEventID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();
    let mealPlanGroceryListItemID = fakeID();

    const exampleInput = new MealPlanGroceryListItemUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlanGroceryListItem>();
    mock
      .onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(200, exampleResponse);

    client
      .updateMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID, exampleInput)
      .then((response: APIResponse<MealPlanGroceryListItem>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();
    let mealPlanGroceryListItemID = fakeID();

    const exampleInput = new MealPlanGroceryListItemUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(
      buildObligatoryError('updateMealPlanGroceryListItem user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(200, exampleResponse);

    expect(client.updateMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();
    let mealPlanGroceryListItemID = fakeID();

    const exampleInput = new MealPlanGroceryListItemUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(
      buildObligatoryError('updateMealPlanGroceryListItem service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(500, exampleResponse);

    expect(client.updateMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a MealPlanTaskStatus', () => {
    let mealPlanID = fakeID();
    let mealPlanTaskID = fakeID();

    const exampleInput = new MealPlanTaskStatusChangeRequestInput();

    const exampleResponse = new APIResponse<MealPlanTask>();
    mock.onPatch(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`).reply(200, exampleResponse);

    client
      .updateMealPlanTaskStatus(mealPlanID, mealPlanTaskID, exampleInput)
      .then((response: APIResponse<MealPlanTask>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.patch.length).toBe(1);
        expect(mock.history.patch[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.patch[0].headers).toHaveProperty('Authorization');
        expect((mock.history.patch[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a MealPlanTaskStatus', () => {
    let mealPlanID = fakeID();
    let mealPlanTaskID = fakeID();

    const exampleInput = new MealPlanTaskStatusChangeRequestInput();

    const exampleResponse = new APIResponse<MealPlanTask>(buildObligatoryError('updateMealPlanTaskStatus user error'));
    mock.onPatch(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`).reply(200, exampleResponse);

    expect(client.updateMealPlanTaskStatus(mealPlanID, mealPlanTaskID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a MealPlanTaskStatus', () => {
    let mealPlanID = fakeID();
    let mealPlanTaskID = fakeID();

    const exampleInput = new MealPlanTaskStatusChangeRequestInput();

    const exampleResponse = new APIResponse<MealPlanTask>(
      buildObligatoryError('updateMealPlanTaskStatus service error'),
    );
    mock.onPatch(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`).reply(500, exampleResponse);

    expect(client.updateMealPlanTaskStatus(mealPlanID, mealPlanTaskID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a Password', () => {
    const exampleInput = new PasswordUpdateInput();

    const exampleResponse = new APIResponse<PasswordResetResponse>();
    mock.onPut(`${baseURL}/api/v1/users/password/new`).reply(200, exampleResponse);

    client
      .updatePassword(exampleInput)
      .then((response: AxiosResponse<APIResponse<PasswordResetResponse>>) => {
        expect(response.data).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a Password', () => {
    const exampleInput = new PasswordUpdateInput();

    const exampleResponse = new APIResponse<PasswordResetResponse>(buildObligatoryError('updatePassword user error'));
    mock.onPut(`${baseURL}/api/v1/users/password/new`).reply(200, exampleResponse);

    expect(client.updatePassword(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during update a Password', () => {
    const exampleInput = new PasswordUpdateInput();

    const exampleResponse = new APIResponse<PasswordResetResponse>(
      buildObligatoryError('updatePassword service error'),
    );
    mock.onPut(`${baseURL}/api/v1/users/password/new`).reply(500, exampleResponse);

    expect(client.updatePassword(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should update a Recipe', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipeUpdateRequestInput();

    const exampleResponse = new APIResponse<Recipe>();
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}`).reply(200, exampleResponse);

    client
      .updateRecipe(recipeID, exampleInput)
      .then((response: APIResponse<Recipe>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a Recipe', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipeUpdateRequestInput();

    const exampleResponse = new APIResponse<Recipe>(buildObligatoryError('updateRecipe user error'));
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}`).reply(200, exampleResponse);

    expect(client.updateRecipe(recipeID, exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during update a Recipe', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipeUpdateRequestInput();

    const exampleResponse = new APIResponse<Recipe>(buildObligatoryError('updateRecipe service error'));
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}`).reply(500, exampleResponse);

    expect(client.updateRecipe(recipeID, exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should update a RecipePrepTask', () => {
    let recipeID = fakeID();
    let recipePrepTaskID = fakeID();

    const exampleInput = new RecipePrepTaskUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipePrepTask>();
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(200, exampleResponse);

    client
      .updateRecipePrepTask(recipeID, recipePrepTaskID, exampleInput)
      .then((response: APIResponse<RecipePrepTask>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a RecipePrepTask', () => {
    let recipeID = fakeID();
    let recipePrepTaskID = fakeID();

    const exampleInput = new RecipePrepTaskUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipePrepTask>(buildObligatoryError('updateRecipePrepTask user error'));
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(200, exampleResponse);

    expect(client.updateRecipePrepTask(recipeID, recipePrepTaskID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a RecipePrepTask', () => {
    let recipeID = fakeID();
    let recipePrepTaskID = fakeID();

    const exampleInput = new RecipePrepTaskUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipePrepTask>(buildObligatoryError('updateRecipePrepTask service error'));
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(500, exampleResponse);

    expect(client.updateRecipePrepTask(recipeID, recipePrepTaskID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a RecipeRating', () => {
    let recipeID = fakeID();
    let recipeRatingID = fakeID();

    const exampleInput = new RecipeRatingUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeRating>();
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(200, exampleResponse);

    client
      .updateRecipeRating(recipeID, recipeRatingID, exampleInput)
      .then((response: APIResponse<RecipeRating>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a RecipeRating', () => {
    let recipeID = fakeID();
    let recipeRatingID = fakeID();

    const exampleInput = new RecipeRatingUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeRating>(buildObligatoryError('updateRecipeRating user error'));
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(200, exampleResponse);

    expect(client.updateRecipeRating(recipeID, recipeRatingID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a RecipeRating', () => {
    let recipeID = fakeID();
    let recipeRatingID = fakeID();

    const exampleInput = new RecipeRatingUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeRating>(buildObligatoryError('updateRecipeRating service error'));
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(500, exampleResponse);

    expect(client.updateRecipeRating(recipeID, recipeRatingID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a RecipeStep', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStep>();
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(200, exampleResponse);

    client
      .updateRecipeStep(recipeID, recipeStepID, exampleInput)
      .then((response: APIResponse<RecipeStep>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a RecipeStep', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStep>(buildObligatoryError('updateRecipeStep user error'));
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(200, exampleResponse);

    expect(client.updateRecipeStep(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a RecipeStep', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStep>(buildObligatoryError('updateRecipeStep service error'));
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(500, exampleResponse);

    expect(client.updateRecipeStep(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepCompletionConditionID = fakeID();

    const exampleInput = new RecipeStepCompletionConditionUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>();
    mock
      .onPut(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(200, exampleResponse);

    client
      .updateRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID, exampleInput)
      .then((response: APIResponse<RecipeStepCompletionCondition>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepCompletionConditionID = fakeID();

    const exampleInput = new RecipeStepCompletionConditionUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(
      buildObligatoryError('updateRecipeStepCompletionCondition user error'),
    );
    mock
      .onPut(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(200, exampleResponse);

    expect(
      client.updateRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID, exampleInput),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during update a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepCompletionConditionID = fakeID();

    const exampleInput = new RecipeStepCompletionConditionUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(
      buildObligatoryError('updateRecipeStepCompletionCondition service error'),
    );
    mock
      .onPut(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(500, exampleResponse);

    expect(
      client.updateRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID, exampleInput),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should update a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepIngredientID = fakeID();

    const exampleInput = new RecipeStepIngredientUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepIngredient>();
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(200, exampleResponse);

    client
      .updateRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID, exampleInput)
      .then((response: APIResponse<RecipeStepIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepIngredientID = fakeID();

    const exampleInput = new RecipeStepIngredientUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepIngredient>(
      buildObligatoryError('updateRecipeStepIngredient user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(200, exampleResponse);

    expect(
      client.updateRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID, exampleInput),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during update a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepIngredientID = fakeID();

    const exampleInput = new RecipeStepIngredientUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepIngredient>(
      buildObligatoryError('updateRecipeStepIngredient service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(500, exampleResponse);

    expect(
      client.updateRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID, exampleInput),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should update a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepInstrumentID = fakeID();

    const exampleInput = new RecipeStepInstrumentUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepInstrument>();
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(200, exampleResponse);

    client
      .updateRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID, exampleInput)
      .then((response: APIResponse<RecipeStepInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepInstrumentID = fakeID();

    const exampleInput = new RecipeStepInstrumentUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepInstrument>(
      buildObligatoryError('updateRecipeStepInstrument user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(200, exampleResponse);

    expect(
      client.updateRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID, exampleInput),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during update a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepInstrumentID = fakeID();

    const exampleInput = new RecipeStepInstrumentUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepInstrument>(
      buildObligatoryError('updateRecipeStepInstrument service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(500, exampleResponse);

    expect(
      client.updateRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID, exampleInput),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should update a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepProductID = fakeID();

    const exampleInput = new RecipeStepProductUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepProduct>();
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(200, exampleResponse);

    client
      .updateRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID, exampleInput)
      .then((response: APIResponse<RecipeStepProduct>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepProductID = fakeID();

    const exampleInput = new RecipeStepProductUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepProduct>(
      buildObligatoryError('updateRecipeStepProduct user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(200, exampleResponse);

    expect(client.updateRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepProductID = fakeID();

    const exampleInput = new RecipeStepProductUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepProduct>(
      buildObligatoryError('updateRecipeStepProduct service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(500, exampleResponse);

    expect(client.updateRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepVesselID = fakeID();

    const exampleInput = new RecipeStepVesselUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepVessel>();
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(200, exampleResponse);

    client
      .updateRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID, exampleInput)
      .then((response: APIResponse<RecipeStepVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepVesselID = fakeID();

    const exampleInput = new RecipeStepVesselUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepVessel>(
      buildObligatoryError('updateRecipeStepVessel user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(200, exampleResponse);

    expect(client.updateRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepVesselID = fakeID();

    const exampleInput = new RecipeStepVesselUpdateRequestInput();

    const exampleResponse = new APIResponse<RecipeStepVessel>(
      buildObligatoryError('updateRecipeStepVessel service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(500, exampleResponse);

    expect(client.updateRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ServiceSettingConfiguration', () => {
    let serviceSettingConfigurationID = fakeID();

    const exampleInput = new ServiceSettingConfigurationUpdateRequestInput();

    const exampleResponse = new APIResponse<ServiceSettingConfiguration>();
    mock
      .onPut(`${baseURL}/api/v1/settings/configurations/${serviceSettingConfigurationID}`)
      .reply(200, exampleResponse);

    client
      .updateServiceSettingConfiguration(serviceSettingConfigurationID, exampleInput)
      .then((response: APIResponse<ServiceSettingConfiguration>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ServiceSettingConfiguration', () => {
    let serviceSettingConfigurationID = fakeID();

    const exampleInput = new ServiceSettingConfigurationUpdateRequestInput();

    const exampleResponse = new APIResponse<ServiceSettingConfiguration>(
      buildObligatoryError('updateServiceSettingConfiguration user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/settings/configurations/${serviceSettingConfigurationID}`)
      .reply(200, exampleResponse);

    expect(client.updateServiceSettingConfiguration(serviceSettingConfigurationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ServiceSettingConfiguration', () => {
    let serviceSettingConfigurationID = fakeID();

    const exampleInput = new ServiceSettingConfigurationUpdateRequestInput();

    const exampleResponse = new APIResponse<ServiceSettingConfiguration>(
      buildObligatoryError('updateServiceSettingConfiguration service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/settings/configurations/${serviceSettingConfigurationID}`)
      .reply(500, exampleResponse);

    expect(client.updateServiceSettingConfiguration(serviceSettingConfigurationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a UserDetails', () => {
    const exampleInput = new UserDetailsUpdateRequestInput();

    const exampleResponse = new APIResponse<User>();
    mock.onPut(`${baseURL}/api/v1/users/details`).reply(200, exampleResponse);

    client
      .updateUserDetails(exampleInput)
      .then((response: APIResponse<User>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a UserDetails', () => {
    const exampleInput = new UserDetailsUpdateRequestInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('updateUserDetails user error'));
    mock.onPut(`${baseURL}/api/v1/users/details`).reply(200, exampleResponse);

    expect(client.updateUserDetails(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during update a UserDetails', () => {
    const exampleInput = new UserDetailsUpdateRequestInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('updateUserDetails service error'));
    mock.onPut(`${baseURL}/api/v1/users/details`).reply(500, exampleResponse);

    expect(client.updateUserDetails(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should update a UserEmailAddress', () => {
    const exampleInput = new UserEmailAddressUpdateInput();

    const exampleResponse = new APIResponse<User>();
    mock.onPut(`${baseURL}/api/v1/users/email_address`).reply(200, exampleResponse);

    client
      .updateUserEmailAddress(exampleInput)
      .then((response: APIResponse<User>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a UserEmailAddress', () => {
    const exampleInput = new UserEmailAddressUpdateInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('updateUserEmailAddress user error'));
    mock.onPut(`${baseURL}/api/v1/users/email_address`).reply(200, exampleResponse);

    expect(client.updateUserEmailAddress(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during update a UserEmailAddress', () => {
    const exampleInput = new UserEmailAddressUpdateInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('updateUserEmailAddress service error'));
    mock.onPut(`${baseURL}/api/v1/users/email_address`).reply(500, exampleResponse);

    expect(client.updateUserEmailAddress(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should update a UserIngredientPreference', () => {
    let userIngredientPreferenceID = fakeID();

    const exampleInput = new UserIngredientPreferenceUpdateRequestInput();

    const exampleResponse = new APIResponse<UserIngredientPreference>();
    mock
      .onPut(`${baseURL}/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`)
      .reply(200, exampleResponse);

    client
      .updateUserIngredientPreference(userIngredientPreferenceID, exampleInput)
      .then((response: APIResponse<UserIngredientPreference>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a UserIngredientPreference', () => {
    let userIngredientPreferenceID = fakeID();

    const exampleInput = new UserIngredientPreferenceUpdateRequestInput();

    const exampleResponse = new APIResponse<UserIngredientPreference>(
      buildObligatoryError('updateUserIngredientPreference user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`)
      .reply(200, exampleResponse);

    expect(client.updateUserIngredientPreference(userIngredientPreferenceID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a UserIngredientPreference', () => {
    let userIngredientPreferenceID = fakeID();

    const exampleInput = new UserIngredientPreferenceUpdateRequestInput();

    const exampleResponse = new APIResponse<UserIngredientPreference>(
      buildObligatoryError('updateUserIngredientPreference service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`)
      .reply(500, exampleResponse);

    expect(client.updateUserIngredientPreference(userIngredientPreferenceID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a UserNotification', () => {
    let userNotificationID = fakeID();

    const exampleInput = new UserNotificationUpdateRequestInput();

    const exampleResponse = new APIResponse<UserNotification>();
    mock.onPatch(`${baseURL}/api/v1/user_notifications/${userNotificationID}`).reply(200, exampleResponse);

    client
      .updateUserNotification(userNotificationID, exampleInput)
      .then((response: APIResponse<UserNotification>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.patch.length).toBe(1);
        expect(mock.history.patch[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.patch[0].headers).toHaveProperty('Authorization');
        expect((mock.history.patch[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a UserNotification', () => {
    let userNotificationID = fakeID();

    const exampleInput = new UserNotificationUpdateRequestInput();

    const exampleResponse = new APIResponse<UserNotification>(
      buildObligatoryError('updateUserNotification user error'),
    );
    mock.onPatch(`${baseURL}/api/v1/user_notifications/${userNotificationID}`).reply(200, exampleResponse);

    expect(client.updateUserNotification(userNotificationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a UserNotification', () => {
    let userNotificationID = fakeID();

    const exampleInput = new UserNotificationUpdateRequestInput();

    const exampleResponse = new APIResponse<UserNotification>(
      buildObligatoryError('updateUserNotification service error'),
    );
    mock.onPatch(`${baseURL}/api/v1/user_notifications/${userNotificationID}`).reply(500, exampleResponse);

    expect(client.updateUserNotification(userNotificationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a UserUsername', () => {
    const exampleInput = new UsernameUpdateInput();

    const exampleResponse = new APIResponse<User>();
    mock.onPut(`${baseURL}/api/v1/users/username`).reply(200, exampleResponse);

    client
      .updateUserUsername(exampleInput)
      .then((response: APIResponse<User>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a UserUsername', () => {
    const exampleInput = new UsernameUpdateInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('updateUserUsername user error'));
    mock.onPut(`${baseURL}/api/v1/users/username`).reply(200, exampleResponse);

    expect(client.updateUserUsername(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during update a UserUsername', () => {
    const exampleInput = new UsernameUpdateInput();

    const exampleResponse = new APIResponse<User>(buildObligatoryError('updateUserUsername service error'));
    mock.onPut(`${baseURL}/api/v1/users/username`).reply(500, exampleResponse);

    expect(client.updateUserUsername(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should update a ValidIngredient', () => {
    let validIngredientID = fakeID();

    const exampleInput = new ValidIngredientUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredient>();
    mock.onPut(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(200, exampleResponse);

    client
      .updateValidIngredient(validIngredientID, exampleInput)
      .then((response: APIResponse<ValidIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidIngredient', () => {
    let validIngredientID = fakeID();

    const exampleInput = new ValidIngredientUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredient>(buildObligatoryError('updateValidIngredient user error'));
    mock.onPut(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(200, exampleResponse);

    expect(client.updateValidIngredient(validIngredientID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidIngredient', () => {
    let validIngredientID = fakeID();

    const exampleInput = new ValidIngredientUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredient>(
      buildObligatoryError('updateValidIngredient service error'),
    );
    mock.onPut(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(500, exampleResponse);

    expect(client.updateValidIngredient(validIngredientID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ValidIngredientGroup', () => {
    let validIngredientGroupID = fakeID();

    const exampleInput = new ValidIngredientGroupUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientGroup>();
    mock.onPut(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(200, exampleResponse);

    client
      .updateValidIngredientGroup(validIngredientGroupID, exampleInput)
      .then((response: APIResponse<ValidIngredientGroup>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidIngredientGroup', () => {
    let validIngredientGroupID = fakeID();

    const exampleInput = new ValidIngredientGroupUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientGroup>(
      buildObligatoryError('updateValidIngredientGroup user error'),
    );
    mock.onPut(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(200, exampleResponse);

    expect(client.updateValidIngredientGroup(validIngredientGroupID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidIngredientGroup', () => {
    let validIngredientGroupID = fakeID();

    const exampleInput = new ValidIngredientGroupUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientGroup>(
      buildObligatoryError('updateValidIngredientGroup service error'),
    );
    mock.onPut(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(500, exampleResponse);

    expect(client.updateValidIngredientGroup(validIngredientGroupID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = fakeID();

    const exampleInput = new ValidIngredientMeasurementUnitUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>();
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(200, exampleResponse);

    client
      .updateValidIngredientMeasurementUnit(validIngredientMeasurementUnitID, exampleInput)
      .then((response: APIResponse<ValidIngredientMeasurementUnit>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = fakeID();

    const exampleInput = new ValidIngredientMeasurementUnitUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(
      buildObligatoryError('updateValidIngredientMeasurementUnit user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(200, exampleResponse);

    expect(client.updateValidIngredientMeasurementUnit(validIngredientMeasurementUnitID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = fakeID();

    const exampleInput = new ValidIngredientMeasurementUnitUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(
      buildObligatoryError('updateValidIngredientMeasurementUnit service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(500, exampleResponse);

    expect(client.updateValidIngredientMeasurementUnit(validIngredientMeasurementUnitID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ValidIngredientPreparation', () => {
    let validIngredientPreparationID = fakeID();

    const exampleInput = new ValidIngredientPreparationUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientPreparation>();
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(200, exampleResponse);

    client
      .updateValidIngredientPreparation(validIngredientPreparationID, exampleInput)
      .then((response: APIResponse<ValidIngredientPreparation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidIngredientPreparation', () => {
    let validIngredientPreparationID = fakeID();

    const exampleInput = new ValidIngredientPreparationUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientPreparation>(
      buildObligatoryError('updateValidIngredientPreparation user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(200, exampleResponse);

    expect(client.updateValidIngredientPreparation(validIngredientPreparationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidIngredientPreparation', () => {
    let validIngredientPreparationID = fakeID();

    const exampleInput = new ValidIngredientPreparationUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientPreparation>(
      buildObligatoryError('updateValidIngredientPreparation service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(500, exampleResponse);

    expect(client.updateValidIngredientPreparation(validIngredientPreparationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ValidIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleInput = new ValidIngredientStateUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientState>();
    mock.onPut(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(200, exampleResponse);

    client
      .updateValidIngredientState(validIngredientStateID, exampleInput)
      .then((response: APIResponse<ValidIngredientState>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleInput = new ValidIngredientStateUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientState>(
      buildObligatoryError('updateValidIngredientState user error'),
    );
    mock.onPut(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(200, exampleResponse);

    expect(client.updateValidIngredientState(validIngredientStateID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleInput = new ValidIngredientStateUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientState>(
      buildObligatoryError('updateValidIngredientState service error'),
    );
    mock.onPut(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(500, exampleResponse);

    expect(client.updateValidIngredientState(validIngredientStateID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = fakeID();

    const exampleInput = new ValidIngredientStateIngredientUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>();
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(200, exampleResponse);

    client
      .updateValidIngredientStateIngredient(validIngredientStateIngredientID, exampleInput)
      .then((response: APIResponse<ValidIngredientStateIngredient>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = fakeID();

    const exampleInput = new ValidIngredientStateIngredientUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(
      buildObligatoryError('updateValidIngredientStateIngredient user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.updateValidIngredientStateIngredient(validIngredientStateIngredientID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = fakeID();

    const exampleInput = new ValidIngredientStateIngredientUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(
      buildObligatoryError('updateValidIngredientStateIngredient service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.updateValidIngredientStateIngredient(validIngredientStateIngredientID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ValidInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleInput = new ValidInstrumentUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidInstrument>();
    mock.onPut(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(200, exampleResponse);

    client
      .updateValidInstrument(validInstrumentID, exampleInput)
      .then((response: APIResponse<ValidInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleInput = new ValidInstrumentUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidInstrument>(buildObligatoryError('updateValidInstrument user error'));
    mock.onPut(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(200, exampleResponse);

    expect(client.updateValidInstrument(validInstrumentID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleInput = new ValidInstrumentUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidInstrument>(
      buildObligatoryError('updateValidInstrument service error'),
    );
    mock.onPut(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(500, exampleResponse);

    expect(client.updateValidInstrument(validInstrumentID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ValidMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleInput = new ValidMeasurementUnitUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidMeasurementUnit>();
    mock.onPut(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(200, exampleResponse);

    client
      .updateValidMeasurementUnit(validMeasurementUnitID, exampleInput)
      .then((response: APIResponse<ValidMeasurementUnit>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleInput = new ValidMeasurementUnitUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidMeasurementUnit>(
      buildObligatoryError('updateValidMeasurementUnit user error'),
    );
    mock.onPut(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(200, exampleResponse);

    expect(client.updateValidMeasurementUnit(validMeasurementUnitID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleInput = new ValidMeasurementUnitUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidMeasurementUnit>(
      buildObligatoryError('updateValidMeasurementUnit service error'),
    );
    mock.onPut(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(500, exampleResponse);

    expect(client.updateValidMeasurementUnit(validMeasurementUnitID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = fakeID();

    const exampleInput = new ValidMeasurementUnitConversionUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>();
    mock
      .onPut(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(200, exampleResponse);

    client
      .updateValidMeasurementUnitConversion(validMeasurementUnitConversionID, exampleInput)
      .then((response: APIResponse<ValidMeasurementUnitConversion>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = fakeID();

    const exampleInput = new ValidMeasurementUnitConversionUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(
      buildObligatoryError('updateValidMeasurementUnitConversion user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(200, exampleResponse);

    expect(client.updateValidMeasurementUnitConversion(validMeasurementUnitConversionID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = fakeID();

    const exampleInput = new ValidMeasurementUnitConversionUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(
      buildObligatoryError('updateValidMeasurementUnitConversion service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(500, exampleResponse);

    expect(client.updateValidMeasurementUnitConversion(validMeasurementUnitConversionID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ValidPreparation', () => {
    let validPreparationID = fakeID();

    const exampleInput = new ValidPreparationUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidPreparation>();
    mock.onPut(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(200, exampleResponse);

    client
      .updateValidPreparation(validPreparationID, exampleInput)
      .then((response: APIResponse<ValidPreparation>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidPreparation', () => {
    let validPreparationID = fakeID();

    const exampleInput = new ValidPreparationUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidPreparation>(
      buildObligatoryError('updateValidPreparation user error'),
    );
    mock.onPut(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(200, exampleResponse);

    expect(client.updateValidPreparation(validPreparationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidPreparation', () => {
    let validPreparationID = fakeID();

    const exampleInput = new ValidPreparationUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidPreparation>(
      buildObligatoryError('updateValidPreparation service error'),
    );
    mock.onPut(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(500, exampleResponse);

    expect(client.updateValidPreparation(validPreparationID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ValidPreparationInstrument', () => {
    let validPreparationVesselID = fakeID();

    const exampleInput = new ValidPreparationInstrumentUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidPreparationInstrument>();
    mock
      .onPut(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(200, exampleResponse);

    client
      .updateValidPreparationInstrument(validPreparationVesselID, exampleInput)
      .then((response: APIResponse<ValidPreparationInstrument>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidPreparationInstrument', () => {
    let validPreparationVesselID = fakeID();

    const exampleInput = new ValidPreparationInstrumentUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidPreparationInstrument>(
      buildObligatoryError('updateValidPreparationInstrument user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(200, exampleResponse);

    expect(client.updateValidPreparationInstrument(validPreparationVesselID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidPreparationInstrument', () => {
    let validPreparationVesselID = fakeID();

    const exampleInput = new ValidPreparationInstrumentUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidPreparationInstrument>(
      buildObligatoryError('updateValidPreparationInstrument service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(500, exampleResponse);

    expect(client.updateValidPreparationInstrument(validPreparationVesselID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ValidPreparationVessel', () => {
    let validPreparationVesselID = fakeID();

    const exampleInput = new ValidPreparationVesselUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidPreparationVessel>();
    mock.onPut(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(200, exampleResponse);

    client
      .updateValidPreparationVessel(validPreparationVesselID, exampleInput)
      .then((response: APIResponse<ValidPreparationVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidPreparationVessel', () => {
    let validPreparationVesselID = fakeID();

    const exampleInput = new ValidPreparationVesselUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidPreparationVessel>(
      buildObligatoryError('updateValidPreparationVessel user error'),
    );
    mock.onPut(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(200, exampleResponse);

    expect(client.updateValidPreparationVessel(validPreparationVesselID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidPreparationVessel', () => {
    let validPreparationVesselID = fakeID();

    const exampleInput = new ValidPreparationVesselUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidPreparationVessel>(
      buildObligatoryError('updateValidPreparationVessel service error'),
    );
    mock.onPut(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(500, exampleResponse);

    expect(client.updateValidPreparationVessel(validPreparationVesselID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a ValidVessel', () => {
    let validVesselID = fakeID();

    const exampleInput = new ValidVesselUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidVessel>();
    mock.onPut(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(200, exampleResponse);

    client
      .updateValidVessel(validVesselID, exampleInput)
      .then((response: APIResponse<ValidVessel>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a ValidVessel', () => {
    let validVesselID = fakeID();

    const exampleInput = new ValidVesselUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidVessel>(buildObligatoryError('updateValidVessel user error'));
    mock.onPut(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(200, exampleResponse);

    expect(client.updateValidVessel(validVesselID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidVessel', () => {
    let validVesselID = fakeID();

    const exampleInput = new ValidVesselUpdateRequestInput();

    const exampleResponse = new APIResponse<ValidVessel>(buildObligatoryError('updateValidVessel service error'));
    mock.onPut(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(500, exampleResponse);

    expect(client.updateValidVessel(validVesselID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a household instrument ownership', () => {
    let householdInstrumentOwnershipID = fakeID();

    const exampleInput = new HouseholdInstrumentOwnershipUpdateRequestInput();

    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>();
    mock
      .onPut(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(200, exampleResponse);

    client
      .updateHouseholdInstrumentOwnership(householdInstrumentOwnershipID, exampleInput)
      .then((response: APIResponse<HouseholdInstrumentOwnership>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a household instrument ownership', () => {
    let householdInstrumentOwnershipID = fakeID();

    const exampleInput = new HouseholdInstrumentOwnershipUpdateRequestInput();

    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(
      buildObligatoryError('updateHouseholdInstrumentOwnership user error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(200, exampleResponse);

    expect(client.updateHouseholdInstrumentOwnership(householdInstrumentOwnershipID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a household instrument ownership', () => {
    let householdInstrumentOwnershipID = fakeID();

    const exampleInput = new HouseholdInstrumentOwnershipUpdateRequestInput();

    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(
      buildObligatoryError('updateHouseholdInstrumentOwnership service error'),
    );
    mock
      .onPut(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(500, exampleResponse);

    expect(client.updateHouseholdInstrumentOwnership(householdInstrumentOwnershipID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a household', () => {
    let householdID = fakeID();

    const exampleInput = new HouseholdUpdateRequestInput();

    const exampleResponse = new APIResponse<Household>();
    mock.onPut(`${baseURL}/api/v1/households/${householdID}`).reply(200, exampleResponse);

    client
      .updateHousehold(householdID, exampleInput)
      .then((response: APIResponse<Household>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a household', () => {
    let householdID = fakeID();

    const exampleInput = new HouseholdUpdateRequestInput();

    const exampleResponse = new APIResponse<Household>(buildObligatoryError('updateHousehold user error'));
    mock.onPut(`${baseURL}/api/v1/households/${householdID}`).reply(200, exampleResponse);

    expect(client.updateHousehold(householdID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a household', () => {
    let householdID = fakeID();

    const exampleInput = new HouseholdUpdateRequestInput();

    const exampleResponse = new APIResponse<Household>(buildObligatoryError('updateHousehold service error'));
    mock.onPut(`${baseURL}/api/v1/households/${householdID}`).reply(500, exampleResponse);

    expect(client.updateHousehold(householdID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should update a meal plan option vote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();
    let mealPlanOptionVoteID = fakeID();

    const exampleInput = new MealPlanOptionVoteUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlanOptionVote>();
    mock
      .onPut(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(200, exampleResponse);

    client
      .updateMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID, exampleInput)
      .then((response: APIResponse<MealPlanOptionVote>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a meal plan option vote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();
    let mealPlanOptionVoteID = fakeID();

    const exampleInput = new MealPlanOptionVoteUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlanOptionVote>(
      buildObligatoryError('updateMealPlanOptionVote user error'),
    );
    mock
      .onPut(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(200, exampleResponse);

    expect(
      client.updateMealPlanOptionVote(
        mealPlanID,
        mealPlanEventID,
        mealPlanOptionID,
        mealPlanOptionVoteID,
        exampleInput,
      ),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during update a meal plan option vote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();
    let mealPlanOptionVoteID = fakeID();

    const exampleInput = new MealPlanOptionVoteUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlanOptionVote>(
      buildObligatoryError('updateMealPlanOptionVote service error'),
    );
    mock
      .onPut(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(500, exampleResponse);

    expect(
      client.updateMealPlanOptionVote(
        mealPlanID,
        mealPlanEventID,
        mealPlanOptionID,
        mealPlanOptionVoteID,
        exampleInput,
      ),
    ).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should update a meal plan option', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleInput = new MealPlanOptionUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlanOption>();
    mock
      .onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(200, exampleResponse);

    client
      .updateMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID, exampleInput)
      .then((response: APIResponse<MealPlanOption>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.put.length).toBe(1);
        expect(mock.history.put[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.put[0].headers).toHaveProperty('Authorization');
        expect((mock.history.put[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update a meal plan option', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleInput = new MealPlanOptionUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlanOption>(buildObligatoryError('updateMealPlanOption user error'));
    mock
      .onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(200, exampleResponse);

    expect(client.updateMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it('should appropriately raise service errors when they occur during update a meal plan option', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleInput = new MealPlanOptionUpdateRequestInput();

    const exampleResponse = new APIResponse<MealPlanOption>(buildObligatoryError('updateMealPlanOption service error'));
    mock
      .onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(500, exampleResponse);

    expect(client.updateMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID, exampleInput)).rejects.toEqual(
      new Error(exampleResponse.error?.message),
    );
  });

  it("should update a user's account status", () => {
    const exampleInput = new UserAccountStatusUpdateInput();

    const exampleResponse = new APIResponse<UserStatusResponse>();
    mock.onPost(`${baseURL}/api/v1/admin/users/status`).reply(201, exampleResponse);

    client
      .adminUpdateUserStatus(exampleInput)
      .then((response: APIResponse<UserStatusResponse>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);
        expect(mock.history.post[0].data).toBe(JSON.stringify(exampleInput));
        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it("should appropriately raise errors when they occur during update a user's account status", () => {
    const exampleInput = new UserAccountStatusUpdateInput();

    const exampleResponse = new APIResponse<UserStatusResponse>(
      buildObligatoryError('adminUpdateUserStatus user error'),
    );
    mock.onPost(`${baseURL}/api/v1/admin/users/status`).reply(201, exampleResponse);

    expect(client.adminUpdateUserStatus(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it("should appropriately raise service errors when they occur during update a user's account status", () => {
    const exampleInput = new UserAccountStatusUpdateInput();

    const exampleResponse = new APIResponse<UserStatusResponse>(
      buildObligatoryError('adminUpdateUserStatus service error'),
    );
    mock.onPost(`${baseURL}/api/v1/admin/users/status`).reply(500, exampleResponse);

    expect(client.adminUpdateUserStatus(exampleInput)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should update the default household assigned at login', () => {
    let householdID = fakeID();

    const exampleResponse = new APIResponse<Household>();
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/default`).reply(201, exampleResponse);

    client
      .setDefaultHousehold(householdID)
      .then((response: APIResponse<Household>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.post.length).toBe(1);

        expect(mock.history.post[0].headers).toHaveProperty('Authorization');
        expect((mock.history.post[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it('should appropriately raise errors when they occur during update the default household assigned at login', () => {
    let householdID = fakeID();

    const exampleResponse = new APIResponse<Household>(buildObligatoryError('setDefaultHousehold user error'));
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/default`).reply(201, exampleResponse);

    expect(client.setDefaultHousehold(householdID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });

  it('should appropriately raise service errors when they occur during update the default household assigned at login', () => {
    let householdID = fakeID();

    const exampleResponse = new APIResponse<Household>(buildObligatoryError('setDefaultHousehold service error'));
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/default`).reply(500, exampleResponse);

    expect(client.setDefaultHousehold(householdID)).rejects.toEqual(new Error(exampleResponse.error?.message));
  });
});
