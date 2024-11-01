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
  UserDataCollection,
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

    const expectedError = buildObligatoryError('acceptHouseholdInvitation user error');
    const exampleResponse = new APIResponse<HouseholdInvitation>(expectedError);
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/accept`).reply(200, exampleResponse);

    expect(client.acceptHouseholdInvitation(householdInvitationID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Accepts a received household invitation', () => {
    let householdInvitationID = fakeID();

    const exampleInput = new HouseholdInvitationUpdateRequestInput();

    const expectedError = buildObligatoryError('acceptHouseholdInvitation service error');
    const exampleResponse = new APIResponse<HouseholdInvitation>(expectedError);
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/accept`).reply(500, exampleResponse);

    expect(client.acceptHouseholdInvitation(householdInvitationID, exampleInput)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('aggregateUserDataReport user error');
    const exampleResponse = new APIResponse<UserDataCollectionResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/data_privacy/disclose`).reply(201, exampleResponse);

    expect(client.aggregateUserDataReport()).rejects.toEqual(expectedError.error);
  });

  it("should appropriately raise service errors when they occur during Aggregates a user's data into a big disclosure blob", () => {
    const expectedError = buildObligatoryError('aggregateUserDataReport service error');
    const exampleResponse = new APIResponse<UserDataCollectionResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/data_privacy/disclose`).reply(500, exampleResponse);

    expect(client.aggregateUserDataReport()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveUserMembership user error');
    const exampleResponse = new APIResponse<HouseholdUserMembership>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/households/${householdID}/members/${userID}`).reply(202, exampleResponse);

    expect(client.archiveUserMembership(householdID, userID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Archive a household user membership', () => {
    let householdID = fakeID();
    let userID = fakeID();

    const expectedError = buildObligatoryError('archiveUserMembership service error');
    const exampleResponse = new APIResponse<HouseholdUserMembership>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/households/${householdID}/members/${userID}`).reply(500, exampleResponse);

    expect(client.archiveUserMembership(householdID, userID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('checkPermissions user error');
    const exampleResponse = new APIResponse<UserPermissionsResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/users/permissions/check`).reply(201, exampleResponse);

    expect(client.checkPermissions(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it("should appropriately raise service errors when they occur during Checks a user's permissions", () => {
    const exampleInput = new UserPermissionsRequestInput();

    const expectedError = buildObligatoryError('checkPermissions service error');
    const exampleResponse = new APIResponse<UserPermissionsResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/users/permissions/check`).reply(500, exampleResponse);

    expect(client.checkPermissions(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('cloneRecipe user error');
    const exampleResponse = new APIResponse<Recipe>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/clone`).reply(201, exampleResponse);

    expect(client.cloneRecipe(recipeID)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Clones a recipe', () => {
    let recipeID = fakeID();

    const expectedError = buildObligatoryError('cloneRecipe service error');
    const exampleResponse = new APIResponse<Recipe>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/clone`).reply(500, exampleResponse);

    expect(client.cloneRecipe(recipeID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createHouseholdInvitation user error');
    const exampleResponse = new APIResponse<HouseholdInvitation>(expectedError);
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/invite`).reply(201, exampleResponse);

    expect(client.createHouseholdInvitation(householdID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Create a household invitation', () => {
    let householdID = fakeID();

    const exampleInput = new HouseholdInvitationCreationRequestInput();

    const expectedError = buildObligatoryError('createHouseholdInvitation service error');
    const exampleResponse = new APIResponse<HouseholdInvitation>(expectedError);
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/invite`).reply(500, exampleResponse);

    expect(client.createHouseholdInvitation(householdID, exampleInput)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('destroyAllUserData user error');
    const exampleResponse = new APIResponse<DataDeletionResponse>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/data_privacy/destroy`).reply(202, exampleResponse);

    expect(client.destroyAllUserData()).rejects.toEqual(expectedError.error);
  });

  it("should raise service errors appropriately when trying to Destroys a user's data", () => {
    const expectedError = buildObligatoryError('destroyAllUserData service error');
    const exampleResponse = new APIResponse<DataDeletionResponse>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/data_privacy/destroy`).reply(500, exampleResponse);

    expect(client.destroyAllUserData()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('finalizeMealPlan user error');
    const exampleResponse = new APIResponse<FinalizeMealPlansResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/finalize`).reply(201, exampleResponse);

    expect(client.finalizeMealPlan(mealPlanID)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Finalizes a meal plan', () => {
    let mealPlanID = fakeID();

    const expectedError = buildObligatoryError('finalizeMealPlan service error');
    const exampleResponse = new APIResponse<FinalizeMealPlansResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/finalize`).reply(500, exampleResponse);

    expect(client.finalizeMealPlan(mealPlanID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('adminLoginForJWT user error');
    const exampleResponse = new APIResponse<JWTResponse>(expectedError);
    mock.onPost(`${baseURL}/users/login/jwt/admin`).reply(201, exampleResponse);

    expect(client.adminLoginForJWT(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Operation for creating JWTResponse', () => {
    const exampleInput = new UserLoginInput();

    const expectedError = buildObligatoryError('adminLoginForJWT service error');
    const exampleResponse = new APIResponse<JWTResponse>(expectedError);
    mock.onPost(`${baseURL}/users/login/jwt/admin`).reply(500, exampleResponse);

    expect(client.adminLoginForJWT(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('loginForJWT user error');
    const exampleResponse = new APIResponse<JWTResponse>(expectedError);
    mock.onPost(`${baseURL}/users/login/jwt`).reply(201, exampleResponse);

    expect(client.loginForJWT(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Operation for creating JWTResponse', () => {
    const exampleInput = new UserLoginInput();

    const expectedError = buildObligatoryError('loginForJWT service error');
    const exampleResponse = new APIResponse<JWTResponse>(expectedError);
    mock.onPost(`${baseURL}/users/login/jwt`).reply(500, exampleResponse);

    expect(client.loginForJWT(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('requestPasswordResetToken user error');
    const exampleResponse = new APIResponse<PasswordResetToken>(expectedError);
    mock.onPost(`${baseURL}/users/password/reset`).reply(201, exampleResponse);

    expect(client.requestPasswordResetToken(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Operation for creating PasswordResetToken', () => {
    const exampleInput = new PasswordResetTokenCreationRequestInput();

    const expectedError = buildObligatoryError('requestPasswordResetToken service error');
    const exampleResponse = new APIResponse<PasswordResetToken>(expectedError);
    mock.onPost(`${baseURL}/users/password/reset`).reply(500, exampleResponse);

    expect(client.requestPasswordResetToken(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('verifyEmailAddress user error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPost(`${baseURL}/users/email_address/verify`).reply(201, exampleResponse);

    expect(client.verifyEmailAddress(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Operation for creating User', () => {
    const exampleInput = new EmailAddressVerificationRequestInput();

    const expectedError = buildObligatoryError('verifyEmailAddress service error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPost(`${baseURL}/users/email_address/verify`).reply(500, exampleResponse);

    expect(client.verifyEmailAddress(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('verifyTOTPSecret user error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPost(`${baseURL}/users/totp_secret/verify`).reply(201, exampleResponse);

    expect(client.verifyTOTPSecret(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Operation for creating User', () => {
    const exampleInput = new TOTPSecretVerificationInput();

    const expectedError = buildObligatoryError('verifyTOTPSecret service error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPost(`${baseURL}/users/totp_secret/verify`).reply(500, exampleResponse);

    expect(client.verifyTOTPSecret(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('requestUsernameReminder user error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPost(`${baseURL}/users/username/reminder`).reply(201, exampleResponse);

    expect(client.requestUsernameReminder(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Operation for creating User', () => {
    const exampleInput = new UsernameReminderRequestInput();

    const expectedError = buildObligatoryError('requestUsernameReminder service error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPost(`${baseURL}/users/username/reminder`).reply(500, exampleResponse);

    expect(client.requestUsernameReminder(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it("should Reads a user's data report from storage", () => {
    let userDataAggregationReportID = fakeID();

    const exampleResponse = new APIResponse<UserDataCollection>();
    mock.onGet(`${baseURL}/api/v1/data_privacy/reports/${userDataAggregationReportID}`).reply(200, exampleResponse);

    client
      .fetchUserDataReport(userDataAggregationReportID)
      .then((response: APIResponse<UserDataCollection>) => {
        expect(response).toEqual(exampleResponse);
      })
      .then(() => {
        expect(mock.history.get.length).toBe(1);
        expect(mock.history.get[0].headers).toHaveProperty('Authorization');
        expect((mock.history.get[0].headers || {})['Authorization']).toBe(`Bearer test-token`);
      });
  });

  it("should raise errors appropriately when trying to Reads a user's data report from storage", () => {
    let userDataAggregationReportID = fakeID();

    const expectedError = buildObligatoryError('fetchUserDataReport user error');
    const exampleResponse = new APIResponse<UserDataCollection>(expectedError);
    mock.onGet(`${baseURL}/api/v1/data_privacy/reports/${userDataAggregationReportID}`).reply(200, exampleResponse);

    expect(client.fetchUserDataReport(userDataAggregationReportID)).rejects.toEqual(expectedError.error);
  });

  it("should raise service errors appropriately when trying to Reads a user's data report from storage", () => {
    let userDataAggregationReportID = fakeID();

    const expectedError = buildObligatoryError('fetchUserDataReport service error');
    const exampleResponse = new APIResponse<UserDataCollection>(expectedError);
    mock.onGet(`${baseURL}/api/v1/data_privacy/reports/${userDataAggregationReportID}`).reply(500, exampleResponse);

    expect(client.fetchUserDataReport(userDataAggregationReportID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('redeemPasswordResetToken user error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPost(`${baseURL}/users/password/reset/redeem`).reply(201, exampleResponse);

    expect(client.redeemPasswordResetToken(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Redeems a password reset token', () => {
    const exampleInput = new PasswordResetTokenRedemptionRequestInput();

    const expectedError = buildObligatoryError('redeemPasswordResetToken service error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPost(`${baseURL}/users/password/reset/redeem`).reply(500, exampleResponse);

    expect(client.redeemPasswordResetToken(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('refreshTOTPSecret user error');
    const exampleResponse = new APIResponse<TOTPSecretRefreshResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/users/totp_secret/new`).reply(201, exampleResponse);

    expect(client.refreshTOTPSecret(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it("should appropriately raise service errors when they occur during Refreshes a user's TOTP secret", () => {
    const exampleInput = new TOTPSecretRefreshInput();

    const expectedError = buildObligatoryError('refreshTOTPSecret service error');
    const exampleResponse = new APIResponse<TOTPSecretRefreshResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/users/totp_secret/new`).reply(500, exampleResponse);

    expect(client.refreshTOTPSecret(exampleInput)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('requestEmailVerificationEmail user error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPost(`${baseURL}/api/v1/users/email_address_verification`).reply(201, exampleResponse);

    expect(client.requestEmailVerificationEmail()).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Requests an email address verification email', () => {
    const expectedError = buildObligatoryError('requestEmailVerificationEmail service error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPost(`${baseURL}/api/v1/users/email_address_verification`).reply(500, exampleResponse);

    expect(client.requestEmailVerificationEmail()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getAuditLogEntryByID user error');
    const exampleResponse = new APIResponse<AuditLogEntry>(expectedError);
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/${auditLogEntryID}`).reply(200, exampleResponse);

    expect(client.getAuditLogEntryByID(auditLogEntryID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Retrieves an audit log entry by ID', () => {
    let auditLogEntryID = fakeID();

    const expectedError = buildObligatoryError('getAuditLogEntryByID service error');
    const exampleResponse = new APIResponse<AuditLogEntry>(expectedError);
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/${auditLogEntryID}`).reply(500, exampleResponse);

    expect(client.getAuditLogEntryByID(auditLogEntryID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getAuditLogEntriesForHousehold user error');
    const exampleResponse = new APIResponse<Array<AuditLogEntry>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/for_household`).reply(200, exampleResponse);

    expect(client.getAuditLogEntriesForHousehold()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Retrieves audit log entries for a household', () => {
    const expectedError = buildObligatoryError('getAuditLogEntriesForHousehold service error');
    const exampleResponse = new APIResponse<Array<AuditLogEntry>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/for_household`).reply(500, exampleResponse);

    expect(client.getAuditLogEntriesForHousehold()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getAuditLogEntriesForUser user error');
    const exampleResponse = new APIResponse<Array<AuditLogEntry>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/for_user`).reply(200, exampleResponse);

    expect(client.getAuditLogEntriesForUser()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Retrieves audit log entries for a user', () => {
    const expectedError = buildObligatoryError('getAuditLogEntriesForUser service error');
    const exampleResponse = new APIResponse<Array<AuditLogEntry>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/for_user`).reply(500, exampleResponse);

    expect(client.getAuditLogEntriesForUser()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('runFinalizeMealPlanWorker user error');
    const exampleResponse = new APIResponse<FinalizeMealPlansResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/workers/finalize_meal_plans`).reply(201, exampleResponse);

    expect(client.runFinalizeMealPlanWorker(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Runs the Finalize Meal Plans worker', () => {
    const exampleInput = new FinalizeMealPlansRequest();

    const expectedError = buildObligatoryError('runFinalizeMealPlanWorker service error');
    const exampleResponse = new APIResponse<FinalizeMealPlansResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/workers/finalize_meal_plans`).reply(500, exampleResponse);

    expect(client.runFinalizeMealPlanWorker(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('runMealPlanGroceryListInitializerWorker user error');
    const exampleResponse = new APIResponse<InitializeMealPlanGroceryListResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/workers/meal_plan_grocery_list_init`).reply(201, exampleResponse);

    expect(client.runMealPlanGroceryListInitializerWorker(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Runs the meal plan grocery list initialization worker', () => {
    const exampleInput = new InitializeMealPlanGroceryListRequest();

    const expectedError = buildObligatoryError('runMealPlanGroceryListInitializerWorker service error');
    const exampleResponse = new APIResponse<InitializeMealPlanGroceryListResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/workers/meal_plan_grocery_list_init`).reply(500, exampleResponse);

    expect(client.runMealPlanGroceryListInitializerWorker(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('runMealPlanTaskCreatorWorker user error');
    const exampleResponse = new APIResponse<CreateMealPlanTasksResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/workers/meal_plan_tasks`).reply(201, exampleResponse);

    expect(client.runMealPlanTaskCreatorWorker(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Runs the meal plan task creator worker', () => {
    const exampleInput = new CreateMealPlanTasksRequest();

    const expectedError = buildObligatoryError('runMealPlanTaskCreatorWorker service error');
    const exampleResponse = new APIResponse<CreateMealPlanTasksResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/workers/meal_plan_tasks`).reply(500, exampleResponse);

    expect(client.runMealPlanTaskCreatorWorker(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchForMeals user error');
    const exampleResponse = new APIResponse<Array<Meal>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meals/search`).reply(200, exampleResponse);

    expect(client.searchForMeals(q)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for Meals', () => {
    let q = fakeID();

    const expectedError = buildObligatoryError('searchForMeals service error');
    const exampleResponse = new APIResponse<Array<Meal>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meals/search`).reply(500, exampleResponse);

    expect(client.searchForMeals(q)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchForRecipes user error');
    const exampleResponse = new APIResponse<Array<Recipe>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/search`).reply(200, exampleResponse);

    expect(client.searchForRecipes(q)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for Recipes', () => {
    let q = fakeID();

    const expectedError = buildObligatoryError('searchForRecipes service error');
    const exampleResponse = new APIResponse<Array<Recipe>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/search`).reply(500, exampleResponse);

    expect(client.searchForRecipes(q)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchForServiceSettings user error');
    const exampleResponse = new APIResponse<Array<ServiceSetting>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/settings/search`).reply(200, exampleResponse);

    expect(client.searchForServiceSettings(q)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for ServiceSettings', () => {
    let q = fakeID();

    const expectedError = buildObligatoryError('searchForServiceSettings service error');
    const exampleResponse = new APIResponse<Array<ServiceSetting>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/settings/search`).reply(500, exampleResponse);

    expect(client.searchForServiceSettings(q)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchForUsers user error');
    const exampleResponse = new APIResponse<Array<User>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/users/search`).reply(200, exampleResponse);

    expect(client.searchForUsers(q)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for Users', () => {
    let q = fakeID();

    const expectedError = buildObligatoryError('searchForUsers service error');
    const exampleResponse = new APIResponse<Array<User>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/users/search`).reply(500, exampleResponse);

    expect(client.searchForUsers(q)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchForValidIngredientGroups user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientGroup>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/search`).reply(200, exampleResponse);

    expect(client.searchForValidIngredientGroups(q)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for ValidIngredientGroups', () => {
    let q = fakeID();

    const expectedError = buildObligatoryError('searchForValidIngredientGroups service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientGroup>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/search`).reply(500, exampleResponse);

    expect(client.searchForValidIngredientGroups(q)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchForValidIngredientStates user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientState>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/search`).reply(200, exampleResponse);

    expect(client.searchForValidIngredientStates(q)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for ValidIngredientStates', () => {
    let q = fakeID();

    const expectedError = buildObligatoryError('searchForValidIngredientStates service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientState>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/search`).reply(500, exampleResponse);

    expect(client.searchForValidIngredientStates(q)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchValidIngredientsByPreparation user error');
    const exampleResponse = new APIResponse<Array<ValidIngredient>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/by_preparation/${validPreparationID}`).reply(200, exampleResponse);

    expect(client.searchValidIngredientsByPreparation(q, validPreparationID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for ValidIngredients by a ValidPreparation ID', () => {
    let q = fakeID();
    let validPreparationID = fakeID();

    const expectedError = buildObligatoryError('searchValidIngredientsByPreparation service error');
    const exampleResponse = new APIResponse<Array<ValidIngredient>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/by_preparation/${validPreparationID}`).reply(500, exampleResponse);

    expect(client.searchValidIngredientsByPreparation(q, validPreparationID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchForValidIngredients user error');
    const exampleResponse = new APIResponse<Array<ValidIngredient>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/search`).reply(200, exampleResponse);

    expect(client.searchForValidIngredients(q)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for ValidIngredients', () => {
    let q = fakeID();

    const expectedError = buildObligatoryError('searchForValidIngredients service error');
    const exampleResponse = new APIResponse<Array<ValidIngredient>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/search`).reply(500, exampleResponse);

    expect(client.searchForValidIngredients(q)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchForValidInstruments user error');
    const exampleResponse = new APIResponse<Array<ValidInstrument>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_instruments/search`).reply(200, exampleResponse);

    expect(client.searchForValidInstruments(q)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for ValidInstruments', () => {
    let q = fakeID();

    const expectedError = buildObligatoryError('searchForValidInstruments service error');
    const exampleResponse = new APIResponse<Array<ValidInstrument>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_instruments/search`).reply(500, exampleResponse);

    expect(client.searchForValidInstruments(q)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchValidMeasurementUnitsByIngredient user error');
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_units/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.searchValidMeasurementUnitsByIngredient(q, validIngredientID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for ValidMeasurementUnits by a ValidIngredient ID', () => {
    let q = fakeID();
    let validIngredientID = fakeID();

    const expectedError = buildObligatoryError('searchValidMeasurementUnitsByIngredient service error');
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_units/by_ingredient/${validIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.searchValidMeasurementUnitsByIngredient(q, validIngredientID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchForValidMeasurementUnits user error');
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/search`).reply(200, exampleResponse);

    expect(client.searchForValidMeasurementUnits(q)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for ValidMeasurementUnits', () => {
    let q = fakeID();

    const expectedError = buildObligatoryError('searchForValidMeasurementUnits service error');
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/search`).reply(500, exampleResponse);

    expect(client.searchForValidMeasurementUnits(q)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchForValidPreparations user error');
    const exampleResponse = new APIResponse<Array<ValidPreparation>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparations/search`).reply(200, exampleResponse);

    expect(client.searchForValidPreparations(q)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for ValidPreparations', () => {
    let q = fakeID();

    const expectedError = buildObligatoryError('searchForValidPreparations service error');
    const exampleResponse = new APIResponse<Array<ValidPreparation>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparations/search`).reply(500, exampleResponse);

    expect(client.searchForValidPreparations(q)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('searchForValidVessels user error');
    const exampleResponse = new APIResponse<Array<ValidVessel>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_vessels/search`).reply(200, exampleResponse);

    expect(client.searchForValidVessels(q)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to Searches for ValidVessels', () => {
    let q = fakeID();

    const expectedError = buildObligatoryError('searchForValidVessels service error');
    const exampleResponse = new APIResponse<Array<ValidVessel>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_vessels/search`).reply(500, exampleResponse);

    expect(client.searchForValidVessels(q)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('transferHouseholdOwnership user error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/transfer`).reply(201, exampleResponse);

    expect(client.transferHouseholdOwnership(householdID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Transfer household ownership to another user', () => {
    let householdID = fakeID();

    const exampleInput = new HouseholdOwnershipTransferInput();

    const expectedError = buildObligatoryError('transferHouseholdOwnership service error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/transfer`).reply(500, exampleResponse);

    expect(client.transferHouseholdOwnership(householdID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateHouseholdMemberPermissions user error');
    const exampleResponse = new APIResponse<UserPermissionsResponse>(expectedError);
    mock
      .onPatch(`${baseURL}/api/v1/households/${householdID}/members/${userID}/permissions`)
      .reply(200, exampleResponse);

    expect(client.updateHouseholdMemberPermissions(householdID, userID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it("should appropriately raise service errors when they occur during Update a household member's household permissions", () => {
    let householdID = fakeID();
    let userID = fakeID();

    const exampleInput = new ModifyUserPermissionsInput();

    const expectedError = buildObligatoryError('updateHouseholdMemberPermissions service error');
    const exampleResponse = new APIResponse<UserPermissionsResponse>(expectedError);
    mock
      .onPatch(`${baseURL}/api/v1/households/${householdID}/members/${userID}/permissions`)
      .reply(500, exampleResponse);

    expect(client.updateHouseholdMemberPermissions(householdID, userID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('uploadUserAvatar user error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPost(`${baseURL}/api/v1/users/avatar/upload`).reply(201, exampleResponse);

    expect(client.uploadUserAvatar(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during Uploads a new user avatar', () => {
    const exampleInput = new AvatarUpdateInput();

    const expectedError = buildObligatoryError('uploadUserAvatar service error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPost(`${baseURL}/api/v1/users/avatar/upload`).reply(500, exampleResponse);

    expect(client.uploadUserAvatar(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveMeal user error');
    const exampleResponse = new APIResponse<Meal>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/meals/${mealID}`).reply(202, exampleResponse);

    expect(client.archiveMeal(mealID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a Meal', () => {
    let mealID = fakeID();

    const expectedError = buildObligatoryError('archiveMeal service error');
    const exampleResponse = new APIResponse<Meal>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/meals/${mealID}`).reply(500, exampleResponse);

    expect(client.archiveMeal(mealID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveMealPlan user error');
    const exampleResponse = new APIResponse<MealPlan>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(202, exampleResponse);

    expect(client.archiveMealPlan(mealPlanID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a MealPlan', () => {
    let mealPlanID = fakeID();

    const expectedError = buildObligatoryError('archiveMealPlan service error');
    const exampleResponse = new APIResponse<MealPlan>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(500, exampleResponse);

    expect(client.archiveMealPlan(mealPlanID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveMealPlanEvent user error');
    const exampleResponse = new APIResponse<MealPlanEvent>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(202, exampleResponse);

    expect(client.archiveMealPlanEvent(mealPlanID, mealPlanEventID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a MealPlanEvent', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const expectedError = buildObligatoryError('archiveMealPlanEvent service error');
    const exampleResponse = new APIResponse<MealPlanEvent>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(500, exampleResponse);

    expect(client.archiveMealPlanEvent(mealPlanID, mealPlanEventID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveMealPlanGroceryListItem user error');
    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(202, exampleResponse);

    expect(client.archiveMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to archive a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();
    let mealPlanGroceryListItemID = fakeID();

    const expectedError = buildObligatoryError('archiveMealPlanGroceryListItem service error');
    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(500, exampleResponse);

    expect(client.archiveMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('archiveMealPlanOption user error');
    const exampleResponse = new APIResponse<MealPlanOption>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(202, exampleResponse);

    expect(client.archiveMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to archive a MealPlanOption', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const expectedError = buildObligatoryError('archiveMealPlanOption service error');
    const exampleResponse = new APIResponse<MealPlanOption>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(500, exampleResponse);

    expect(client.archiveMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('archiveOAuth2Client user error');
    const exampleResponse = new APIResponse<OAuth2Client>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/oauth2_clients/${oauth2ClientID}`).reply(202, exampleResponse);

    expect(client.archiveOAuth2Client(oauth2ClientID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a OAuth2Client', () => {
    let oauth2ClientID = fakeID();

    const expectedError = buildObligatoryError('archiveOAuth2Client service error');
    const exampleResponse = new APIResponse<OAuth2Client>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/oauth2_clients/${oauth2ClientID}`).reply(500, exampleResponse);

    expect(client.archiveOAuth2Client(oauth2ClientID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveRecipe user error');
    const exampleResponse = new APIResponse<Recipe>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}`).reply(202, exampleResponse);

    expect(client.archiveRecipe(recipeID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a Recipe', () => {
    let recipeID = fakeID();

    const expectedError = buildObligatoryError('archiveRecipe service error');
    const exampleResponse = new APIResponse<Recipe>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}`).reply(500, exampleResponse);

    expect(client.archiveRecipe(recipeID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveRecipePrepTask user error');
    const exampleResponse = new APIResponse<RecipePrepTask>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(202, exampleResponse);

    expect(client.archiveRecipePrepTask(recipeID, recipePrepTaskID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a RecipePrepTask', () => {
    let recipeID = fakeID();
    let recipePrepTaskID = fakeID();

    const expectedError = buildObligatoryError('archiveRecipePrepTask service error');
    const exampleResponse = new APIResponse<RecipePrepTask>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(500, exampleResponse);

    expect(client.archiveRecipePrepTask(recipeID, recipePrepTaskID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveRecipeRating user error');
    const exampleResponse = new APIResponse<RecipeRating>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(202, exampleResponse);

    expect(client.archiveRecipeRating(recipeID, recipeRatingID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a RecipeRating', () => {
    let recipeID = fakeID();
    let recipeRatingID = fakeID();

    const expectedError = buildObligatoryError('archiveRecipeRating service error');
    const exampleResponse = new APIResponse<RecipeRating>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(500, exampleResponse);

    expect(client.archiveRecipeRating(recipeID, recipeRatingID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveRecipeStep user error');
    const exampleResponse = new APIResponse<RecipeStep>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(202, exampleResponse);

    expect(client.archiveRecipeStep(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a RecipeStep', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const expectedError = buildObligatoryError('archiveRecipeStep service error');
    const exampleResponse = new APIResponse<RecipeStep>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(500, exampleResponse);

    expect(client.archiveRecipeStep(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveRecipeStepCompletionCondition user error');
    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(expectedError);
    mock
      .onDelete(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(202, exampleResponse);

    expect(
      client.archiveRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID),
    ).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepCompletionConditionID = fakeID();

    const expectedError = buildObligatoryError('archiveRecipeStepCompletionCondition service error');
    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(expectedError);
    mock
      .onDelete(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(500, exampleResponse);

    expect(
      client.archiveRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID),
    ).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveRecipeStepIngredient user error');
    const exampleResponse = new APIResponse<RecipeStepIngredient>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(202, exampleResponse);

    expect(client.archiveRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to archive a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepIngredientID = fakeID();

    const expectedError = buildObligatoryError('archiveRecipeStepIngredient service error');
    const exampleResponse = new APIResponse<RecipeStepIngredient>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.archiveRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('archiveRecipeStepInstrument user error');
    const exampleResponse = new APIResponse<RecipeStepInstrument>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(202, exampleResponse);

    expect(client.archiveRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to archive a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepInstrumentID = fakeID();

    const expectedError = buildObligatoryError('archiveRecipeStepInstrument service error');
    const exampleResponse = new APIResponse<RecipeStepInstrument>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(500, exampleResponse);

    expect(client.archiveRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('archiveRecipeStepProduct user error');
    const exampleResponse = new APIResponse<RecipeStepProduct>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(202, exampleResponse);

    expect(client.archiveRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to archive a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepProductID = fakeID();

    const expectedError = buildObligatoryError('archiveRecipeStepProduct service error');
    const exampleResponse = new APIResponse<RecipeStepProduct>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(500, exampleResponse);

    expect(client.archiveRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('archiveRecipeStepVessel user error');
    const exampleResponse = new APIResponse<RecipeStepVessel>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(202, exampleResponse);

    expect(client.archiveRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to archive a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepVesselID = fakeID();

    const expectedError = buildObligatoryError('archiveRecipeStepVessel service error');
    const exampleResponse = new APIResponse<RecipeStepVessel>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(500, exampleResponse);

    expect(client.archiveRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('archiveServiceSetting user error');
    const exampleResponse = new APIResponse<ServiceSetting>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/settings/${serviceSettingID}`).reply(202, exampleResponse);

    expect(client.archiveServiceSetting(serviceSettingID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a ServiceSetting', () => {
    let serviceSettingID = fakeID();

    const expectedError = buildObligatoryError('archiveServiceSetting service error');
    const exampleResponse = new APIResponse<ServiceSetting>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/settings/${serviceSettingID}`).reply(500, exampleResponse);

    expect(client.archiveServiceSetting(serviceSettingID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveServiceSettingConfiguration user error');
    const exampleResponse = new APIResponse<ServiceSettingConfiguration>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/settings/configurations/${serviceSettingConfigurationID}`)
      .reply(202, exampleResponse);

    expect(client.archiveServiceSettingConfiguration(serviceSettingConfigurationID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to archive a ServiceSettingConfiguration', () => {
    let serviceSettingConfigurationID = fakeID();

    const expectedError = buildObligatoryError('archiveServiceSettingConfiguration service error');
    const exampleResponse = new APIResponse<ServiceSettingConfiguration>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/settings/configurations/${serviceSettingConfigurationID}`)
      .reply(500, exampleResponse);

    expect(client.archiveServiceSettingConfiguration(serviceSettingConfigurationID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('archiveUser user error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/users/${userID}`).reply(202, exampleResponse);

    expect(client.archiveUser(userID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a User', () => {
    let userID = fakeID();

    const expectedError = buildObligatoryError('archiveUser service error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/users/${userID}`).reply(500, exampleResponse);

    expect(client.archiveUser(userID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveUserIngredientPreference user error');
    const exampleResponse = new APIResponse<UserIngredientPreference>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`)
      .reply(202, exampleResponse);

    expect(client.archiveUserIngredientPreference(userIngredientPreferenceID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a UserIngredientPreference', () => {
    let userIngredientPreferenceID = fakeID();

    const expectedError = buildObligatoryError('archiveUserIngredientPreference service error');
    const exampleResponse = new APIResponse<UserIngredientPreference>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`)
      .reply(500, exampleResponse);

    expect(client.archiveUserIngredientPreference(userIngredientPreferenceID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveValidIngredient user error');
    const exampleResponse = new APIResponse<ValidIngredient>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(202, exampleResponse);

    expect(client.archiveValidIngredient(validIngredientID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a ValidIngredient', () => {
    let validIngredientID = fakeID();

    const expectedError = buildObligatoryError('archiveValidIngredient service error');
    const exampleResponse = new APIResponse<ValidIngredient>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(500, exampleResponse);

    expect(client.archiveValidIngredient(validIngredientID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveValidIngredientGroup user error');
    const exampleResponse = new APIResponse<ValidIngredientGroup>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(202, exampleResponse);

    expect(client.archiveValidIngredientGroup(validIngredientGroupID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a ValidIngredientGroup', () => {
    let validIngredientGroupID = fakeID();

    const expectedError = buildObligatoryError('archiveValidIngredientGroup service error');
    const exampleResponse = new APIResponse<ValidIngredientGroup>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(500, exampleResponse);

    expect(client.archiveValidIngredientGroup(validIngredientGroupID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveValidIngredientMeasurementUnit user error');
    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(202, exampleResponse);

    expect(client.archiveValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = fakeID();

    const expectedError = buildObligatoryError('archiveValidIngredientMeasurementUnit service error');
    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(500, exampleResponse);

    expect(client.archiveValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('archiveValidIngredientPreparation user error');
    const exampleResponse = new APIResponse<ValidIngredientPreparation>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(202, exampleResponse);

    expect(client.archiveValidIngredientPreparation(validIngredientPreparationID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a ValidIngredientPreparation', () => {
    let validIngredientPreparationID = fakeID();

    const expectedError = buildObligatoryError('archiveValidIngredientPreparation service error');
    const exampleResponse = new APIResponse<ValidIngredientPreparation>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(500, exampleResponse);

    expect(client.archiveValidIngredientPreparation(validIngredientPreparationID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveValidIngredientState user error');
    const exampleResponse = new APIResponse<ValidIngredientState>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(202, exampleResponse);

    expect(client.archiveValidIngredientState(validIngredientStateID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a ValidIngredientState', () => {
    let validIngredientStateID = fakeID();

    const expectedError = buildObligatoryError('archiveValidIngredientState service error');
    const exampleResponse = new APIResponse<ValidIngredientState>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(500, exampleResponse);

    expect(client.archiveValidIngredientState(validIngredientStateID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveValidIngredientStateIngredient user error');
    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(202, exampleResponse);

    expect(client.archiveValidIngredientStateIngredient(validIngredientStateIngredientID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = fakeID();

    const expectedError = buildObligatoryError('archiveValidIngredientStateIngredient service error');
    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.archiveValidIngredientStateIngredient(validIngredientStateIngredientID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('archiveValidInstrument user error');
    const exampleResponse = new APIResponse<ValidInstrument>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(202, exampleResponse);

    expect(client.archiveValidInstrument(validInstrumentID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a ValidInstrument', () => {
    let validInstrumentID = fakeID();

    const expectedError = buildObligatoryError('archiveValidInstrument service error');
    const exampleResponse = new APIResponse<ValidInstrument>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(500, exampleResponse);

    expect(client.archiveValidInstrument(validInstrumentID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveValidMeasurementUnit user error');
    const exampleResponse = new APIResponse<ValidMeasurementUnit>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(202, exampleResponse);

    expect(client.archiveValidMeasurementUnit(validMeasurementUnitID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a ValidMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const expectedError = buildObligatoryError('archiveValidMeasurementUnit service error');
    const exampleResponse = new APIResponse<ValidMeasurementUnit>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(500, exampleResponse);

    expect(client.archiveValidMeasurementUnit(validMeasurementUnitID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveValidMeasurementUnitConversion user error');
    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(202, exampleResponse);

    expect(client.archiveValidMeasurementUnitConversion(validMeasurementUnitConversionID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to archive a ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = fakeID();

    const expectedError = buildObligatoryError('archiveValidMeasurementUnitConversion service error');
    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(500, exampleResponse);

    expect(client.archiveValidMeasurementUnitConversion(validMeasurementUnitConversionID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('archiveValidPreparation user error');
    const exampleResponse = new APIResponse<ValidPreparation>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(202, exampleResponse);

    expect(client.archiveValidPreparation(validPreparationID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a ValidPreparation', () => {
    let validPreparationID = fakeID();

    const expectedError = buildObligatoryError('archiveValidPreparation service error');
    const exampleResponse = new APIResponse<ValidPreparation>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(500, exampleResponse);

    expect(client.archiveValidPreparation(validPreparationID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveValidPreparationInstrument user error');
    const exampleResponse = new APIResponse<ValidPreparationInstrument>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(202, exampleResponse);

    expect(client.archiveValidPreparationInstrument(validPreparationVesselID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a ValidPreparationInstrument', () => {
    let validPreparationVesselID = fakeID();

    const expectedError = buildObligatoryError('archiveValidPreparationInstrument service error');
    const exampleResponse = new APIResponse<ValidPreparationInstrument>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(500, exampleResponse);

    expect(client.archiveValidPreparationInstrument(validPreparationVesselID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveValidPreparationVessel user error');
    const exampleResponse = new APIResponse<ValidPreparationVessel>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`)
      .reply(202, exampleResponse);

    expect(client.archiveValidPreparationVessel(validPreparationVesselID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a ValidPreparationVessel', () => {
    let validPreparationVesselID = fakeID();

    const expectedError = buildObligatoryError('archiveValidPreparationVessel service error');
    const exampleResponse = new APIResponse<ValidPreparationVessel>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`)
      .reply(500, exampleResponse);

    expect(client.archiveValidPreparationVessel(validPreparationVesselID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveValidVessel user error');
    const exampleResponse = new APIResponse<ValidVessel>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(202, exampleResponse);

    expect(client.archiveValidVessel(validVesselID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a ValidVessel', () => {
    let validVesselID = fakeID();

    const expectedError = buildObligatoryError('archiveValidVessel service error');
    const exampleResponse = new APIResponse<ValidVessel>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(500, exampleResponse);

    expect(client.archiveValidVessel(validVesselID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveWebhook user error');
    const exampleResponse = new APIResponse<Webhook>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/webhooks/${webhookID}`).reply(202, exampleResponse);

    expect(client.archiveWebhook(webhookID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a Webhook', () => {
    let webhookID = fakeID();

    const expectedError = buildObligatoryError('archiveWebhook service error');
    const exampleResponse = new APIResponse<Webhook>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/webhooks/${webhookID}`).reply(500, exampleResponse);

    expect(client.archiveWebhook(webhookID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveWebhookTriggerEvent user error');
    const exampleResponse = new APIResponse<WebhookTriggerEvent>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/webhooks/${webhookID}/trigger_events/${webhookTriggerEventID}`)
      .reply(202, exampleResponse);

    expect(client.archiveWebhookTriggerEvent(webhookID, webhookTriggerEventID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a WebhookTriggerEvent', () => {
    let webhookID = fakeID();
    let webhookTriggerEventID = fakeID();

    const expectedError = buildObligatoryError('archiveWebhookTriggerEvent service error');
    const exampleResponse = new APIResponse<WebhookTriggerEvent>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/webhooks/${webhookID}/trigger_events/${webhookTriggerEventID}`)
      .reply(500, exampleResponse);

    expect(client.archiveWebhookTriggerEvent(webhookID, webhookTriggerEventID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveHouseholdInstrumentOwnership user error');
    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(202, exampleResponse);

    expect(client.archiveHouseholdInstrumentOwnership(householdInstrumentOwnershipID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to archive a household instrument ownership', () => {
    let householdInstrumentOwnershipID = fakeID();

    const expectedError = buildObligatoryError('archiveHouseholdInstrumentOwnership service error');
    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(expectedError);
    mock
      .onDelete(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(500, exampleResponse);

    expect(client.archiveHouseholdInstrumentOwnership(householdInstrumentOwnershipID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('archiveHousehold user error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/households/${householdID}`).reply(202, exampleResponse);

    expect(client.archiveHousehold(householdID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a household', () => {
    let householdID = fakeID();

    const expectedError = buildObligatoryError('archiveHousehold service error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onDelete(`${baseURL}/api/v1/households/${householdID}`).reply(500, exampleResponse);

    expect(client.archiveHousehold(householdID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('archiveMealPlanOptionVote user error');
    const exampleResponse = new APIResponse<MealPlanOptionVote>(expectedError);
    mock
      .onDelete(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(202, exampleResponse);

    expect(
      client.archiveMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID),
    ).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to archive a meal plan option vote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();
    let mealPlanOptionVoteID = fakeID();

    const expectedError = buildObligatoryError('archiveMealPlanOptionVote service error');
    const exampleResponse = new APIResponse<MealPlanOptionVote>(expectedError);
    mock
      .onDelete(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(500, exampleResponse);

    expect(
      client.archiveMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID),
    ).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('cancelHouseholdInvitation user error');
    const exampleResponse = new APIResponse<HouseholdInvitation>(expectedError);
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/cancel`).reply(200, exampleResponse);

    expect(client.cancelHouseholdInvitation(householdInvitationID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during cancel a sent household invitation', () => {
    let householdInvitationID = fakeID();

    const exampleInput = new HouseholdInvitationUpdateRequestInput();

    const expectedError = buildObligatoryError('cancelHouseholdInvitation service error');
    const exampleResponse = new APIResponse<HouseholdInvitation>(expectedError);
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/cancel`).reply(500, exampleResponse);

    expect(client.cancelHouseholdInvitation(householdInvitationID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createMeal user error');
    const exampleResponse = new APIResponse<Meal>(expectedError);
    mock.onPost(`${baseURL}/api/v1/meals`).reply(201, exampleResponse);

    expect(client.createMeal(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a Meal', () => {
    const exampleInput = new MealCreationRequestInput();

    const expectedError = buildObligatoryError('createMeal service error');
    const exampleResponse = new APIResponse<Meal>(expectedError);
    mock.onPost(`${baseURL}/api/v1/meals`).reply(500, exampleResponse);

    expect(client.createMeal(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createMealPlan user error');
    const exampleResponse = new APIResponse<MealPlan>(expectedError);
    mock.onPost(`${baseURL}/api/v1/meal_plans`).reply(201, exampleResponse);

    expect(client.createMealPlan(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a MealPlan', () => {
    const exampleInput = new MealPlanCreationRequestInput();

    const expectedError = buildObligatoryError('createMealPlan service error');
    const exampleResponse = new APIResponse<MealPlan>(expectedError);
    mock.onPost(`${baseURL}/api/v1/meal_plans`).reply(500, exampleResponse);

    expect(client.createMealPlan(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createMealPlanEvent user error');
    const exampleResponse = new APIResponse<MealPlanEvent>(expectedError);
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events`).reply(201, exampleResponse);

    expect(client.createMealPlanEvent(mealPlanID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a MealPlanEvent', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanEventCreationRequestInput();

    const expectedError = buildObligatoryError('createMealPlanEvent service error');
    const exampleResponse = new APIResponse<MealPlanEvent>(expectedError);
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events`).reply(500, exampleResponse);

    expect(client.createMealPlanEvent(mealPlanID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createMealPlanGroceryListItem user error');
    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(expectedError);
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items`).reply(201, exampleResponse);

    expect(client.createMealPlanGroceryListItem(mealPlanID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanGroceryListItemCreationRequestInput();

    const expectedError = buildObligatoryError('createMealPlanGroceryListItem service error');
    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(expectedError);
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items`).reply(500, exampleResponse);

    expect(client.createMealPlanGroceryListItem(mealPlanID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createMealPlanOptionVote user error');
    const exampleResponse = new APIResponse<Array<MealPlanOptionVote>>(expectedError);
    mock
      .onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/vote`)
      .reply(201, exampleResponse);

    expect(client.createMealPlanOptionVote(mealPlanID, mealPlanEventID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during create a MealPlanOptionVote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleInput = new MealPlanOptionVoteCreationRequestInput();

    const expectedError = buildObligatoryError('createMealPlanOptionVote service error');
    const exampleResponse = new APIResponse<Array<MealPlanOptionVote>>(expectedError);
    mock
      .onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/vote`)
      .reply(500, exampleResponse);

    expect(client.createMealPlanOptionVote(mealPlanID, mealPlanEventID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('createMealPlanTask user error');
    const exampleResponse = new APIResponse<MealPlanTask>(expectedError);
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks`).reply(201, exampleResponse);

    expect(client.createMealPlanTask(mealPlanID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a MealPlanTask', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanTaskCreationRequestInput();

    const expectedError = buildObligatoryError('createMealPlanTask service error');
    const exampleResponse = new APIResponse<MealPlanTask>(expectedError);
    mock.onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks`).reply(500, exampleResponse);

    expect(client.createMealPlanTask(mealPlanID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createOAuth2Client user error');
    const exampleResponse = new APIResponse<OAuth2ClientCreationResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/oauth2_clients`).reply(201, exampleResponse);

    expect(client.createOAuth2Client(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a OAuth2Client', () => {
    const exampleInput = new OAuth2ClientCreationRequestInput();

    const expectedError = buildObligatoryError('createOAuth2Client service error');
    const exampleResponse = new APIResponse<OAuth2ClientCreationResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/oauth2_clients`).reply(500, exampleResponse);

    expect(client.createOAuth2Client(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createRecipe user error');
    const exampleResponse = new APIResponse<Recipe>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes`).reply(201, exampleResponse);

    expect(client.createRecipe(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a Recipe', () => {
    const exampleInput = new RecipeCreationRequestInput();

    const expectedError = buildObligatoryError('createRecipe service error');
    const exampleResponse = new APIResponse<Recipe>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes`).reply(500, exampleResponse);

    expect(client.createRecipe(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createRecipePrepTask user error');
    const exampleResponse = new APIResponse<RecipePrepTask>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks`).reply(201, exampleResponse);

    expect(client.createRecipePrepTask(recipeID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a RecipePrepTask', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipePrepTaskCreationRequestInput();

    const expectedError = buildObligatoryError('createRecipePrepTask service error');
    const exampleResponse = new APIResponse<RecipePrepTask>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks`).reply(500, exampleResponse);

    expect(client.createRecipePrepTask(recipeID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createRecipeRating user error');
    const exampleResponse = new APIResponse<RecipeRating>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/ratings`).reply(201, exampleResponse);

    expect(client.createRecipeRating(recipeID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a RecipeRating', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipeRatingCreationRequestInput();

    const expectedError = buildObligatoryError('createRecipeRating service error');
    const exampleResponse = new APIResponse<RecipeRating>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/ratings`).reply(500, exampleResponse);

    expect(client.createRecipeRating(recipeID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createRecipeStep user error');
    const exampleResponse = new APIResponse<RecipeStep>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps`).reply(201, exampleResponse);

    expect(client.createRecipeStep(recipeID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a RecipeStep', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipeStepCreationRequestInput();

    const expectedError = buildObligatoryError('createRecipeStep service error');
    const exampleResponse = new APIResponse<RecipeStep>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps`).reply(500, exampleResponse);

    expect(client.createRecipeStep(recipeID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createRecipeStepCompletionCondition user error');
    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(expectedError);
    mock
      .onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`)
      .reply(201, exampleResponse);

    expect(client.createRecipeStepCompletionCondition(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during create a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepCompletionConditionForExistingRecipeCreationRequestInput();

    const expectedError = buildObligatoryError('createRecipeStepCompletionCondition service error');
    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(expectedError);
    mock
      .onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`)
      .reply(500, exampleResponse);

    expect(client.createRecipeStepCompletionCondition(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('createRecipeStepIngredient user error');
    const exampleResponse = new APIResponse<RecipeStepIngredient>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`).reply(201, exampleResponse);

    expect(client.createRecipeStepIngredient(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during create a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepIngredientCreationRequestInput();

    const expectedError = buildObligatoryError('createRecipeStepIngredient service error');
    const exampleResponse = new APIResponse<RecipeStepIngredient>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`).reply(500, exampleResponse);

    expect(client.createRecipeStepIngredient(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('createRecipeStepInstrument user error');
    const exampleResponse = new APIResponse<RecipeStepInstrument>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`).reply(201, exampleResponse);

    expect(client.createRecipeStepInstrument(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during create a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepInstrumentCreationRequestInput();

    const expectedError = buildObligatoryError('createRecipeStepInstrument service error');
    const exampleResponse = new APIResponse<RecipeStepInstrument>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`).reply(500, exampleResponse);

    expect(client.createRecipeStepInstrument(recipeID, recipeStepID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('createRecipeStepProduct user error');
    const exampleResponse = new APIResponse<RecipeStepProduct>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`).reply(201, exampleResponse);

    expect(client.createRecipeStepProduct(recipeID, recipeStepID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepProductCreationRequestInput();

    const expectedError = buildObligatoryError('createRecipeStepProduct service error');
    const exampleResponse = new APIResponse<RecipeStepProduct>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`).reply(500, exampleResponse);

    expect(client.createRecipeStepProduct(recipeID, recipeStepID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createRecipeStepVessel user error');
    const exampleResponse = new APIResponse<RecipeStepVessel>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`).reply(201, exampleResponse);

    expect(client.createRecipeStepVessel(recipeID, recipeStepID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepVesselCreationRequestInput();

    const expectedError = buildObligatoryError('createRecipeStepVessel service error');
    const exampleResponse = new APIResponse<RecipeStepVessel>(expectedError);
    mock.onPost(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`).reply(500, exampleResponse);

    expect(client.createRecipeStepVessel(recipeID, recipeStepID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createServiceSetting user error');
    const exampleResponse = new APIResponse<ServiceSetting>(expectedError);
    mock.onPost(`${baseURL}/api/v1/settings`).reply(201, exampleResponse);

    expect(client.createServiceSetting(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ServiceSetting', () => {
    const exampleInput = new ServiceSettingCreationRequestInput();

    const expectedError = buildObligatoryError('createServiceSetting service error');
    const exampleResponse = new APIResponse<ServiceSetting>(expectedError);
    mock.onPost(`${baseURL}/api/v1/settings`).reply(500, exampleResponse);

    expect(client.createServiceSetting(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createServiceSettingConfiguration user error');
    const exampleResponse = new APIResponse<ServiceSettingConfiguration>(expectedError);
    mock.onPost(`${baseURL}/api/v1/settings/configurations`).reply(201, exampleResponse);

    expect(client.createServiceSettingConfiguration(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ServiceSettingConfiguration', () => {
    const exampleInput = new ServiceSettingConfigurationCreationRequestInput();

    const expectedError = buildObligatoryError('createServiceSettingConfiguration service error');
    const exampleResponse = new APIResponse<ServiceSettingConfiguration>(expectedError);
    mock.onPost(`${baseURL}/api/v1/settings/configurations`).reply(500, exampleResponse);

    expect(client.createServiceSettingConfiguration(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createUserIngredientPreference user error');
    const exampleResponse = new APIResponse<Array<UserIngredientPreference>>(expectedError);
    mock.onPost(`${baseURL}/api/v1/user_ingredient_preferences`).reply(201, exampleResponse);

    expect(client.createUserIngredientPreference(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a UserIngredientPreference', () => {
    const exampleInput = new UserIngredientPreferenceCreationRequestInput();

    const expectedError = buildObligatoryError('createUserIngredientPreference service error');
    const exampleResponse = new APIResponse<Array<UserIngredientPreference>>(expectedError);
    mock.onPost(`${baseURL}/api/v1/user_ingredient_preferences`).reply(500, exampleResponse);

    expect(client.createUserIngredientPreference(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createUserNotification user error');
    const exampleResponse = new APIResponse<UserNotification>(expectedError);
    mock.onPost(`${baseURL}/api/v1/user_notifications`).reply(201, exampleResponse);

    expect(client.createUserNotification(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a UserNotification', () => {
    const exampleInput = new UserNotificationCreationRequestInput();

    const expectedError = buildObligatoryError('createUserNotification service error');
    const exampleResponse = new APIResponse<UserNotification>(expectedError);
    mock.onPost(`${baseURL}/api/v1/user_notifications`).reply(500, exampleResponse);

    expect(client.createUserNotification(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidIngredient user error');
    const exampleResponse = new APIResponse<ValidIngredient>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_ingredients`).reply(201, exampleResponse);

    expect(client.createValidIngredient(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidIngredient', () => {
    const exampleInput = new ValidIngredientCreationRequestInput();

    const expectedError = buildObligatoryError('createValidIngredient service error');
    const exampleResponse = new APIResponse<ValidIngredient>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_ingredients`).reply(500, exampleResponse);

    expect(client.createValidIngredient(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidIngredientGroup user error');
    const exampleResponse = new APIResponse<ValidIngredientGroup>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_groups`).reply(201, exampleResponse);

    expect(client.createValidIngredientGroup(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidIngredientGroup', () => {
    const exampleInput = new ValidIngredientGroupCreationRequestInput();

    const expectedError = buildObligatoryError('createValidIngredientGroup service error');
    const exampleResponse = new APIResponse<ValidIngredientGroup>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_groups`).reply(500, exampleResponse);

    expect(client.createValidIngredientGroup(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidIngredientMeasurementUnit user error');
    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_measurement_units`).reply(201, exampleResponse);

    expect(client.createValidIngredientMeasurementUnit(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidIngredientMeasurementUnit', () => {
    const exampleInput = new ValidIngredientMeasurementUnitCreationRequestInput();

    const expectedError = buildObligatoryError('createValidIngredientMeasurementUnit service error');
    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_measurement_units`).reply(500, exampleResponse);

    expect(client.createValidIngredientMeasurementUnit(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidIngredientPreparation user error');
    const exampleResponse = new APIResponse<ValidIngredientPreparation>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_preparations`).reply(201, exampleResponse);

    expect(client.createValidIngredientPreparation(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidIngredientPreparation', () => {
    const exampleInput = new ValidIngredientPreparationCreationRequestInput();

    const expectedError = buildObligatoryError('createValidIngredientPreparation service error');
    const exampleResponse = new APIResponse<ValidIngredientPreparation>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_preparations`).reply(500, exampleResponse);

    expect(client.createValidIngredientPreparation(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidIngredientState user error');
    const exampleResponse = new APIResponse<ValidIngredientState>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_states`).reply(201, exampleResponse);

    expect(client.createValidIngredientState(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidIngredientState', () => {
    const exampleInput = new ValidIngredientStateCreationRequestInput();

    const expectedError = buildObligatoryError('createValidIngredientState service error');
    const exampleResponse = new APIResponse<ValidIngredientState>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_states`).reply(500, exampleResponse);

    expect(client.createValidIngredientState(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidIngredientStateIngredient user error');
    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_state_ingredients`).reply(201, exampleResponse);

    expect(client.createValidIngredientStateIngredient(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidIngredientStateIngredient', () => {
    const exampleInput = new ValidIngredientStateIngredientCreationRequestInput();

    const expectedError = buildObligatoryError('createValidIngredientStateIngredient service error');
    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_ingredient_state_ingredients`).reply(500, exampleResponse);

    expect(client.createValidIngredientStateIngredient(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidInstrument user error');
    const exampleResponse = new APIResponse<ValidInstrument>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_instruments`).reply(201, exampleResponse);

    expect(client.createValidInstrument(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidInstrument', () => {
    const exampleInput = new ValidInstrumentCreationRequestInput();

    const expectedError = buildObligatoryError('createValidInstrument service error');
    const exampleResponse = new APIResponse<ValidInstrument>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_instruments`).reply(500, exampleResponse);

    expect(client.createValidInstrument(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidMeasurementUnit user error');
    const exampleResponse = new APIResponse<ValidMeasurementUnit>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_measurement_units`).reply(201, exampleResponse);

    expect(client.createValidMeasurementUnit(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidMeasurementUnit', () => {
    const exampleInput = new ValidMeasurementUnitCreationRequestInput();

    const expectedError = buildObligatoryError('createValidMeasurementUnit service error');
    const exampleResponse = new APIResponse<ValidMeasurementUnit>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_measurement_units`).reply(500, exampleResponse);

    expect(client.createValidMeasurementUnit(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidMeasurementUnitConversion user error');
    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_measurement_conversions`).reply(201, exampleResponse);

    expect(client.createValidMeasurementUnitConversion(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidMeasurementUnitConversion', () => {
    const exampleInput = new ValidMeasurementUnitConversionCreationRequestInput();

    const expectedError = buildObligatoryError('createValidMeasurementUnitConversion service error');
    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_measurement_conversions`).reply(500, exampleResponse);

    expect(client.createValidMeasurementUnitConversion(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidPreparation user error');
    const exampleResponse = new APIResponse<ValidPreparation>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_preparations`).reply(201, exampleResponse);

    expect(client.createValidPreparation(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidPreparation', () => {
    const exampleInput = new ValidPreparationCreationRequestInput();

    const expectedError = buildObligatoryError('createValidPreparation service error');
    const exampleResponse = new APIResponse<ValidPreparation>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_preparations`).reply(500, exampleResponse);

    expect(client.createValidPreparation(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidPreparationInstrument user error');
    const exampleResponse = new APIResponse<ValidPreparationInstrument>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_preparation_instruments`).reply(201, exampleResponse);

    expect(client.createValidPreparationInstrument(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidPreparationInstrument', () => {
    const exampleInput = new ValidPreparationInstrumentCreationRequestInput();

    const expectedError = buildObligatoryError('createValidPreparationInstrument service error');
    const exampleResponse = new APIResponse<ValidPreparationInstrument>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_preparation_instruments`).reply(500, exampleResponse);

    expect(client.createValidPreparationInstrument(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidPreparationVessel user error');
    const exampleResponse = new APIResponse<ValidPreparationVessel>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_preparation_vessels`).reply(201, exampleResponse);

    expect(client.createValidPreparationVessel(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidPreparationVessel', () => {
    const exampleInput = new ValidPreparationVesselCreationRequestInput();

    const expectedError = buildObligatoryError('createValidPreparationVessel service error');
    const exampleResponse = new APIResponse<ValidPreparationVessel>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_preparation_vessels`).reply(500, exampleResponse);

    expect(client.createValidPreparationVessel(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createValidVessel user error');
    const exampleResponse = new APIResponse<ValidVessel>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_vessels`).reply(201, exampleResponse);

    expect(client.createValidVessel(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a ValidVessel', () => {
    const exampleInput = new ValidVesselCreationRequestInput();

    const expectedError = buildObligatoryError('createValidVessel service error');
    const exampleResponse = new APIResponse<ValidVessel>(expectedError);
    mock.onPost(`${baseURL}/api/v1/valid_vessels`).reply(500, exampleResponse);

    expect(client.createValidVessel(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createWebhook user error');
    const exampleResponse = new APIResponse<Webhook>(expectedError);
    mock.onPost(`${baseURL}/api/v1/webhooks`).reply(201, exampleResponse);

    expect(client.createWebhook(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a Webhook', () => {
    const exampleInput = new WebhookCreationRequestInput();

    const expectedError = buildObligatoryError('createWebhook service error');
    const exampleResponse = new APIResponse<Webhook>(expectedError);
    mock.onPost(`${baseURL}/api/v1/webhooks`).reply(500, exampleResponse);

    expect(client.createWebhook(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createWebhookTriggerEvent user error');
    const exampleResponse = new APIResponse<WebhookTriggerEvent>(expectedError);
    mock.onPost(`${baseURL}/api/v1/webhooks/${webhookID}/trigger_events`).reply(201, exampleResponse);

    expect(client.createWebhookTriggerEvent(webhookID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a WebhookTriggerEvent', () => {
    let webhookID = fakeID();

    const exampleInput = new WebhookTriggerEventCreationRequestInput();

    const expectedError = buildObligatoryError('createWebhookTriggerEvent service error');
    const exampleResponse = new APIResponse<WebhookTriggerEvent>(expectedError);
    mock.onPost(`${baseURL}/api/v1/webhooks/${webhookID}/trigger_events`).reply(500, exampleResponse);

    expect(client.createWebhookTriggerEvent(webhookID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createHouseholdInstrumentOwnership user error');
    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(expectedError);
    mock.onPost(`${baseURL}/api/v1/households/instruments`).reply(201, exampleResponse);

    expect(client.createHouseholdInstrumentOwnership(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a household instrument ownership', () => {
    const exampleInput = new HouseholdInstrumentOwnershipCreationRequestInput();

    const expectedError = buildObligatoryError('createHouseholdInstrumentOwnership service error');
    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(expectedError);
    mock.onPost(`${baseURL}/api/v1/households/instruments`).reply(500, exampleResponse);

    expect(client.createHouseholdInstrumentOwnership(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createHousehold user error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onPost(`${baseURL}/api/v1/households`).reply(201, exampleResponse);

    expect(client.createHousehold(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a household', () => {
    const exampleInput = new HouseholdCreationRequestInput();

    const expectedError = buildObligatoryError('createHousehold service error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onPost(`${baseURL}/api/v1/households`).reply(500, exampleResponse);

    expect(client.createHousehold(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createMealPlanOption user error');
    const exampleResponse = new APIResponse<MealPlanOption>(expectedError);
    mock
      .onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`)
      .reply(201, exampleResponse);

    expect(client.createMealPlanOption(mealPlanID, mealPlanEventID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a meal plan option', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleInput = new MealPlanOptionCreationRequestInput();

    const expectedError = buildObligatoryError('createMealPlanOption service error');
    const exampleResponse = new APIResponse<MealPlanOption>(expectedError);
    mock
      .onPost(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`)
      .reply(500, exampleResponse);

    expect(client.createMealPlanOption(mealPlanID, mealPlanEventID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('createUser user error');
    const exampleResponse = new APIResponse<UserCreationResponse>(expectedError);
    mock.onPost(`${baseURL}/users`).reply(201, exampleResponse);

    expect(client.createUser(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during create a new user', () => {
    const exampleInput = new UserRegistrationInput();

    const expectedError = buildObligatoryError('createUser service error');
    const exampleResponse = new APIResponse<UserCreationResponse>(expectedError);
    mock.onPost(`${baseURL}/users`).reply(500, exampleResponse);

    expect(client.createUser(exampleInput)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getHouseholdInstrumentOwnerships user error');
    const exampleResponse = new APIResponse<Array<HouseholdInstrumentOwnership>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/households/instruments`).reply(200, exampleResponse);

    expect(client.getHouseholdInstrumentOwnerships()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to creates a household instrument', () => {
    const expectedError = buildObligatoryError('getHouseholdInstrumentOwnerships service error');
    const exampleResponse = new APIResponse<Array<HouseholdInstrumentOwnership>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/households/instruments`).reply(500, exampleResponse);

    expect(client.getHouseholdInstrumentOwnerships()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getAuthStatus user error');
    const exampleResponse = new APIResponse<UserStatusResponse>(expectedError);
    mock.onGet(`${baseURL}/auth/status`).reply(200, exampleResponse);

    expect(client.getAuthStatus()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a AuthStatus', () => {
    const expectedError = buildObligatoryError('getAuthStatus service error');
    const exampleResponse = new APIResponse<UserStatusResponse>(expectedError);
    mock.onGet(`${baseURL}/auth/status`).reply(500, exampleResponse);

    expect(client.getAuthStatus()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getMeal user error');
    const exampleResponse = new APIResponse<Meal>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meals/${mealID}`).reply(200, exampleResponse);

    expect(client.getMeal(mealID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a Meal', () => {
    let mealID = fakeID();

    const expectedError = buildObligatoryError('getMeal service error');
    const exampleResponse = new APIResponse<Meal>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meals/${mealID}`).reply(500, exampleResponse);

    expect(client.getMeal(mealID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getMealPlan user error');
    const exampleResponse = new APIResponse<MealPlan>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(200, exampleResponse);

    expect(client.getMealPlan(mealPlanID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a MealPlan', () => {
    let mealPlanID = fakeID();

    const expectedError = buildObligatoryError('getMealPlan service error');
    const exampleResponse = new APIResponse<MealPlan>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(500, exampleResponse);

    expect(client.getMealPlan(mealPlanID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getMealPlanEvent user error');
    const exampleResponse = new APIResponse<MealPlanEvent>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(200, exampleResponse);

    expect(client.getMealPlanEvent(mealPlanID, mealPlanEventID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a MealPlanEvent', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const expectedError = buildObligatoryError('getMealPlanEvent service error');
    const exampleResponse = new APIResponse<MealPlanEvent>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(500, exampleResponse);

    expect(client.getMealPlanEvent(mealPlanID, mealPlanEventID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getMealPlanEvents user error');
    const exampleResponse = new APIResponse<Array<MealPlanEvent>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events`).reply(200, exampleResponse);

    expect(client.getMealPlanEvents(mealPlanID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a MealPlanEvents', () => {
    let mealPlanID = fakeID();

    const expectedError = buildObligatoryError('getMealPlanEvents service error');
    const exampleResponse = new APIResponse<Array<MealPlanEvent>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events`).reply(500, exampleResponse);

    expect(client.getMealPlanEvents(mealPlanID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getMealPlanGroceryListItem user error');
    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(200, exampleResponse);

    expect(client.getMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();
    let mealPlanGroceryListItemID = fakeID();

    const expectedError = buildObligatoryError('getMealPlanGroceryListItem service error');
    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(500, exampleResponse);

    expect(client.getMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('getMealPlanGroceryListItemsForMealPlan user error');
    const exampleResponse = new APIResponse<Array<MealPlanGroceryListItem>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items`).reply(200, exampleResponse);

    expect(client.getMealPlanGroceryListItemsForMealPlan(mealPlanID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a MealPlanGroceryListItemsForMealPlan', () => {
    let mealPlanID = fakeID();

    const expectedError = buildObligatoryError('getMealPlanGroceryListItemsForMealPlan service error');
    const exampleResponse = new APIResponse<Array<MealPlanGroceryListItem>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items`).reply(500, exampleResponse);

    expect(client.getMealPlanGroceryListItemsForMealPlan(mealPlanID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getMealPlanTask user error');
    const exampleResponse = new APIResponse<MealPlanTask>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`).reply(200, exampleResponse);

    expect(client.getMealPlanTask(mealPlanID, mealPlanTaskID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a MealPlanTask', () => {
    let mealPlanID = fakeID();
    let mealPlanTaskID = fakeID();

    const expectedError = buildObligatoryError('getMealPlanTask service error');
    const exampleResponse = new APIResponse<MealPlanTask>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`).reply(500, exampleResponse);

    expect(client.getMealPlanTask(mealPlanID, mealPlanTaskID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getMealPlanTasks user error');
    const exampleResponse = new APIResponse<Array<MealPlanTask>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks`).reply(200, exampleResponse);

    expect(client.getMealPlanTasks(mealPlanID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a MealPlanTasks', () => {
    let mealPlanID = fakeID();

    const expectedError = buildObligatoryError('getMealPlanTasks service error');
    const exampleResponse = new APIResponse<Array<MealPlanTask>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks`).reply(500, exampleResponse);

    expect(client.getMealPlanTasks(mealPlanID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getMealPlansForHousehold user error');
    const exampleResponse = new APIResponse<Array<MealPlan>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans`).reply(200, exampleResponse);

    expect(client.getMealPlansForHousehold()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a MealPlans', () => {
    const expectedError = buildObligatoryError('getMealPlansForHousehold service error');
    const exampleResponse = new APIResponse<Array<MealPlan>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meal_plans`).reply(500, exampleResponse);

    expect(client.getMealPlansForHousehold()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getMeals user error');
    const exampleResponse = new APIResponse<Array<Meal>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meals`).reply(200, exampleResponse);

    expect(client.getMeals()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a Meals', () => {
    const expectedError = buildObligatoryError('getMeals service error');
    const exampleResponse = new APIResponse<Array<Meal>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/meals`).reply(500, exampleResponse);

    expect(client.getMeals()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getMermaidDiagramForRecipe user error');
    const exampleResponse = expectedError;
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/mermaid`).reply(200, exampleResponse);

    expect(client.getMermaidDiagramForRecipe(recipeID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a MermaidDiagramForRecipe', () => {
    let recipeID = fakeID();

    const expectedError = buildObligatoryError('getMermaidDiagramForRecipe service error');
    const exampleResponse = expectedError;
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/mermaid`).reply(500, exampleResponse);

    expect(client.getMermaidDiagramForRecipe(recipeID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getOAuth2Client user error');
    const exampleResponse = new APIResponse<OAuth2Client>(expectedError);
    mock.onGet(`${baseURL}/api/v1/oauth2_clients/${oauth2ClientID}`).reply(200, exampleResponse);

    expect(client.getOAuth2Client(oauth2ClientID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a OAuth2Client', () => {
    let oauth2ClientID = fakeID();

    const expectedError = buildObligatoryError('getOAuth2Client service error');
    const exampleResponse = new APIResponse<OAuth2Client>(expectedError);
    mock.onGet(`${baseURL}/api/v1/oauth2_clients/${oauth2ClientID}`).reply(500, exampleResponse);

    expect(client.getOAuth2Client(oauth2ClientID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getOAuth2Clients user error');
    const exampleResponse = new APIResponse<Array<OAuth2Client>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/oauth2_clients`).reply(200, exampleResponse);

    expect(client.getOAuth2Clients()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a OAuth2Clients', () => {
    const expectedError = buildObligatoryError('getOAuth2Clients service error');
    const exampleResponse = new APIResponse<Array<OAuth2Client>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/oauth2_clients`).reply(500, exampleResponse);

    expect(client.getOAuth2Clients()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getRandomValidIngredient user error');
    const exampleResponse = new APIResponse<ValidIngredient>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/random`).reply(200, exampleResponse);

    expect(client.getRandomValidIngredient()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RandomValidIngredient', () => {
    const expectedError = buildObligatoryError('getRandomValidIngredient service error');
    const exampleResponse = new APIResponse<ValidIngredient>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/random`).reply(500, exampleResponse);

    expect(client.getRandomValidIngredient()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getRandomValidInstrument user error');
    const exampleResponse = new APIResponse<ValidInstrument>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_instruments/random`).reply(200, exampleResponse);

    expect(client.getRandomValidInstrument()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RandomValidInstrument', () => {
    const expectedError = buildObligatoryError('getRandomValidInstrument service error');
    const exampleResponse = new APIResponse<ValidInstrument>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_instruments/random`).reply(500, exampleResponse);

    expect(client.getRandomValidInstrument()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getRandomValidPreparation user error');
    const exampleResponse = new APIResponse<ValidPreparation>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparations/random`).reply(200, exampleResponse);

    expect(client.getRandomValidPreparation()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RandomValidPreparation', () => {
    const expectedError = buildObligatoryError('getRandomValidPreparation service error');
    const exampleResponse = new APIResponse<ValidPreparation>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparations/random`).reply(500, exampleResponse);

    expect(client.getRandomValidPreparation()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getRandomValidVessel user error');
    const exampleResponse = new APIResponse<ValidVessel>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_vessels/random`).reply(200, exampleResponse);

    expect(client.getRandomValidVessel()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RandomValidVessel', () => {
    const expectedError = buildObligatoryError('getRandomValidVessel service error');
    const exampleResponse = new APIResponse<ValidVessel>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_vessels/random`).reply(500, exampleResponse);

    expect(client.getRandomValidVessel()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipe user error');
    const exampleResponse = new APIResponse<Recipe>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}`).reply(200, exampleResponse);

    expect(client.getRecipe(recipeID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a Recipe', () => {
    let recipeID = fakeID();

    const expectedError = buildObligatoryError('getRecipe service error');
    const exampleResponse = new APIResponse<Recipe>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}`).reply(500, exampleResponse);

    expect(client.getRecipe(recipeID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipeMealPlanTasks user error');
    const exampleResponse = new APIResponse<RecipePrepTaskStep>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_steps`).reply(200, exampleResponse);

    expect(client.getRecipeMealPlanTasks(recipeID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipeMealPlanTasks', () => {
    let recipeID = fakeID();

    const expectedError = buildObligatoryError('getRecipeMealPlanTasks service error');
    const exampleResponse = new APIResponse<RecipePrepTaskStep>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_steps`).reply(500, exampleResponse);

    expect(client.getRecipeMealPlanTasks(recipeID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipePrepTask user error');
    const exampleResponse = new APIResponse<RecipePrepTask>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(200, exampleResponse);

    expect(client.getRecipePrepTask(recipeID, recipePrepTaskID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipePrepTask', () => {
    let recipeID = fakeID();
    let recipePrepTaskID = fakeID();

    const expectedError = buildObligatoryError('getRecipePrepTask service error');
    const exampleResponse = new APIResponse<RecipePrepTask>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(500, exampleResponse);

    expect(client.getRecipePrepTask(recipeID, recipePrepTaskID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipePrepTasks user error');
    const exampleResponse = new APIResponse<Array<RecipePrepTask>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks`).reply(200, exampleResponse);

    expect(client.getRecipePrepTasks(recipeID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipePrepTasks', () => {
    let recipeID = fakeID();

    const expectedError = buildObligatoryError('getRecipePrepTasks service error');
    const exampleResponse = new APIResponse<Array<RecipePrepTask>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks`).reply(500, exampleResponse);

    expect(client.getRecipePrepTasks(recipeID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipeRating user error');
    const exampleResponse = new APIResponse<RecipeRating>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(200, exampleResponse);

    expect(client.getRecipeRating(recipeID, recipeRatingID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipeRating', () => {
    let recipeID = fakeID();
    let recipeRatingID = fakeID();

    const expectedError = buildObligatoryError('getRecipeRating service error');
    const exampleResponse = new APIResponse<RecipeRating>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(500, exampleResponse);

    expect(client.getRecipeRating(recipeID, recipeRatingID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipeRatingsForRecipe user error');
    const exampleResponse = new APIResponse<Array<RecipeRating>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings`).reply(200, exampleResponse);

    expect(client.getRecipeRatingsForRecipe(recipeID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipeRatings', () => {
    let recipeID = fakeID();

    const expectedError = buildObligatoryError('getRecipeRatingsForRecipe service error');
    const exampleResponse = new APIResponse<Array<RecipeRating>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings`).reply(500, exampleResponse);

    expect(client.getRecipeRatingsForRecipe(recipeID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipeStep user error');
    const exampleResponse = new APIResponse<RecipeStep>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(200, exampleResponse);

    expect(client.getRecipeStep(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStep', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const expectedError = buildObligatoryError('getRecipeStep service error');
    const exampleResponse = new APIResponse<RecipeStep>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(500, exampleResponse);

    expect(client.getRecipeStep(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipeStepCompletionCondition user error');
    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(expectedError);
    mock
      .onGet(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(200, exampleResponse);

    expect(
      client.getRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID),
    ).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepCompletionConditionID = fakeID();

    const expectedError = buildObligatoryError('getRecipeStepCompletionCondition service error');
    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(expectedError);
    mock
      .onGet(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(500, exampleResponse);

    expect(
      client.getRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID),
    ).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipeStepCompletionConditions user error');
    const exampleResponse = new APIResponse<Array<RecipeStepCompletionCondition>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`)
      .reply(200, exampleResponse);

    expect(client.getRecipeStepCompletionConditions(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepCompletionConditions', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const expectedError = buildObligatoryError('getRecipeStepCompletionConditions service error');
    const exampleResponse = new APIResponse<Array<RecipeStepCompletionCondition>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`)
      .reply(500, exampleResponse);

    expect(client.getRecipeStepCompletionConditions(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipeStepIngredient user error');
    const exampleResponse = new APIResponse<RecipeStepIngredient>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.getRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepIngredientID = fakeID();

    const expectedError = buildObligatoryError('getRecipeStepIngredient service error');
    const exampleResponse = new APIResponse<RecipeStepIngredient>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.getRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('getRecipeStepIngredients user error');
    const exampleResponse = new APIResponse<Array<RecipeStepIngredient>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`).reply(200, exampleResponse);

    expect(client.getRecipeStepIngredients(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepIngredients', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const expectedError = buildObligatoryError('getRecipeStepIngredients service error');
    const exampleResponse = new APIResponse<Array<RecipeStepIngredient>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`).reply(500, exampleResponse);

    expect(client.getRecipeStepIngredients(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipeStepInstrument user error');
    const exampleResponse = new APIResponse<RecipeStepInstrument>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(200, exampleResponse);

    expect(client.getRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepInstrumentID = fakeID();

    const expectedError = buildObligatoryError('getRecipeStepInstrument service error');
    const exampleResponse = new APIResponse<RecipeStepInstrument>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(500, exampleResponse);

    expect(client.getRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('getRecipeStepInstruments user error');
    const exampleResponse = new APIResponse<Array<RecipeStepInstrument>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`).reply(200, exampleResponse);

    expect(client.getRecipeStepInstruments(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepInstruments', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const expectedError = buildObligatoryError('getRecipeStepInstruments service error');
    const exampleResponse = new APIResponse<Array<RecipeStepInstrument>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`).reply(500, exampleResponse);

    expect(client.getRecipeStepInstruments(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipeStepProduct user error');
    const exampleResponse = new APIResponse<RecipeStepProduct>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(200, exampleResponse);

    expect(client.getRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepProductID = fakeID();

    const expectedError = buildObligatoryError('getRecipeStepProduct service error');
    const exampleResponse = new APIResponse<RecipeStepProduct>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(500, exampleResponse);

    expect(client.getRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('getRecipeStepProducts user error');
    const exampleResponse = new APIResponse<Array<RecipeStepProduct>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`).reply(200, exampleResponse);

    expect(client.getRecipeStepProducts(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepProducts', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const expectedError = buildObligatoryError('getRecipeStepProducts service error');
    const exampleResponse = new APIResponse<Array<RecipeStepProduct>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`).reply(500, exampleResponse);

    expect(client.getRecipeStepProducts(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipeStepVessel user error');
    const exampleResponse = new APIResponse<RecipeStepVessel>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(200, exampleResponse);

    expect(client.getRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepVesselID = fakeID();

    const expectedError = buildObligatoryError('getRecipeStepVessel service error');
    const exampleResponse = new APIResponse<RecipeStepVessel>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(500, exampleResponse);

    expect(client.getRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipeStepVessels user error');
    const exampleResponse = new APIResponse<Array<RecipeStepVessel>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`).reply(200, exampleResponse);

    expect(client.getRecipeStepVessels(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipeStepVessels', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const expectedError = buildObligatoryError('getRecipeStepVessels service error');
    const exampleResponse = new APIResponse<Array<RecipeStepVessel>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`).reply(500, exampleResponse);

    expect(client.getRecipeStepVessels(recipeID, recipeStepID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getRecipeSteps user error');
    const exampleResponse = new APIResponse<Array<RecipeStep>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps`).reply(200, exampleResponse);

    expect(client.getRecipeSteps(recipeID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a RecipeSteps', () => {
    let recipeID = fakeID();

    const expectedError = buildObligatoryError('getRecipeSteps service error');
    const exampleResponse = new APIResponse<Array<RecipeStep>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps`).reply(500, exampleResponse);

    expect(client.getRecipeSteps(recipeID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getRecipes user error');
    const exampleResponse = new APIResponse<Array<Recipe>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes`).reply(200, exampleResponse);

    expect(client.getRecipes()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a Recipes', () => {
    const expectedError = buildObligatoryError('getRecipes service error');
    const exampleResponse = new APIResponse<Array<Recipe>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/recipes`).reply(500, exampleResponse);

    expect(client.getRecipes()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getSelf user error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onGet(`${baseURL}/api/v1/users/self`).reply(200, exampleResponse);

    expect(client.getSelf()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a Self', () => {
    const expectedError = buildObligatoryError('getSelf service error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onGet(`${baseURL}/api/v1/users/self`).reply(500, exampleResponse);

    expect(client.getSelf()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getServiceSetting user error');
    const exampleResponse = new APIResponse<ServiceSetting>(expectedError);
    mock.onGet(`${baseURL}/api/v1/settings/${serviceSettingID}`).reply(200, exampleResponse);

    expect(client.getServiceSetting(serviceSettingID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ServiceSetting', () => {
    let serviceSettingID = fakeID();

    const expectedError = buildObligatoryError('getServiceSetting service error');
    const exampleResponse = new APIResponse<ServiceSetting>(expectedError);
    mock.onGet(`${baseURL}/api/v1/settings/${serviceSettingID}`).reply(500, exampleResponse);

    expect(client.getServiceSetting(serviceSettingID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getServiceSettingConfigurationByName user error');
    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/settings/configurations/user/${serviceSettingConfigurationName}`)
      .reply(200, exampleResponse);

    expect(client.getServiceSettingConfigurationByName(serviceSettingConfigurationName)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a ServiceSettingConfigurationByName', () => {
    let serviceSettingConfigurationName = fakeID();

    const expectedError = buildObligatoryError('getServiceSettingConfigurationByName service error');
    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/settings/configurations/user/${serviceSettingConfigurationName}`)
      .reply(500, exampleResponse);

    expect(client.getServiceSettingConfigurationByName(serviceSettingConfigurationName)).rejects.toEqual(
      expectedError.error,
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
    const expectedError = buildObligatoryError('getServiceSettingConfigurationsForHousehold user error');
    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/settings/configurations/household`).reply(200, exampleResponse);

    expect(client.getServiceSettingConfigurationsForHousehold()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ServiceSettingConfigurationsForHousehold', () => {
    const expectedError = buildObligatoryError('getServiceSettingConfigurationsForHousehold service error');
    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/settings/configurations/household`).reply(500, exampleResponse);

    expect(client.getServiceSettingConfigurationsForHousehold()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getServiceSettingConfigurationsForUser user error');
    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/settings/configurations/user`).reply(200, exampleResponse);

    expect(client.getServiceSettingConfigurationsForUser()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ServiceSettingConfigurationsForUser', () => {
    const expectedError = buildObligatoryError('getServiceSettingConfigurationsForUser service error');
    const exampleResponse = new APIResponse<Array<ServiceSettingConfiguration>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/settings/configurations/user`).reply(500, exampleResponse);

    expect(client.getServiceSettingConfigurationsForUser()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getServiceSettings user error');
    const exampleResponse = new APIResponse<Array<ServiceSetting>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/settings`).reply(200, exampleResponse);

    expect(client.getServiceSettings()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ServiceSettings', () => {
    const expectedError = buildObligatoryError('getServiceSettings service error');
    const exampleResponse = new APIResponse<Array<ServiceSetting>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/settings`).reply(500, exampleResponse);

    expect(client.getServiceSettings()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getUser user error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onGet(`${baseURL}/api/v1/users/${userID}`).reply(200, exampleResponse);

    expect(client.getUser(userID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a User', () => {
    let userID = fakeID();

    const expectedError = buildObligatoryError('getUser service error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onGet(`${baseURL}/api/v1/users/${userID}`).reply(500, exampleResponse);

    expect(client.getUser(userID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getUserIngredientPreferences user error');
    const exampleResponse = new APIResponse<Array<UserIngredientPreference>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/user_ingredient_preferences`).reply(200, exampleResponse);

    expect(client.getUserIngredientPreferences()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a UserIngredientPreferences', () => {
    const expectedError = buildObligatoryError('getUserIngredientPreferences service error');
    const exampleResponse = new APIResponse<Array<UserIngredientPreference>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/user_ingredient_preferences`).reply(500, exampleResponse);

    expect(client.getUserIngredientPreferences()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getUserNotification user error');
    const exampleResponse = new APIResponse<UserNotification>(expectedError);
    mock.onGet(`${baseURL}/api/v1/user_notifications/${userNotificationID}`).reply(200, exampleResponse);

    expect(client.getUserNotification(userNotificationID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a UserNotification', () => {
    let userNotificationID = fakeID();

    const expectedError = buildObligatoryError('getUserNotification service error');
    const exampleResponse = new APIResponse<UserNotification>(expectedError);
    mock.onGet(`${baseURL}/api/v1/user_notifications/${userNotificationID}`).reply(500, exampleResponse);

    expect(client.getUserNotification(userNotificationID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getUserNotifications user error');
    const exampleResponse = new APIResponse<Array<UserNotification>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/user_notifications`).reply(200, exampleResponse);

    expect(client.getUserNotifications()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a UserNotifications', () => {
    const expectedError = buildObligatoryError('getUserNotifications service error');
    const exampleResponse = new APIResponse<Array<UserNotification>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/user_notifications`).reply(500, exampleResponse);

    expect(client.getUserNotifications()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getUsers user error');
    const exampleResponse = new APIResponse<Array<User>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/users`).reply(200, exampleResponse);

    expect(client.getUsers()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a Users', () => {
    const expectedError = buildObligatoryError('getUsers service error');
    const exampleResponse = new APIResponse<Array<User>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/users`).reply(500, exampleResponse);

    expect(client.getUsers()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidIngredient user error');
    const exampleResponse = new APIResponse<ValidIngredient>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(200, exampleResponse);

    expect(client.getValidIngredient(validIngredientID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredient', () => {
    let validIngredientID = fakeID();

    const expectedError = buildObligatoryError('getValidIngredient service error');
    const exampleResponse = new APIResponse<ValidIngredient>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(500, exampleResponse);

    expect(client.getValidIngredient(validIngredientID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidIngredientGroup user error');
    const exampleResponse = new APIResponse<ValidIngredientGroup>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(200, exampleResponse);

    expect(client.getValidIngredientGroup(validIngredientGroupID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientGroup', () => {
    let validIngredientGroupID = fakeID();

    const expectedError = buildObligatoryError('getValidIngredientGroup service error');
    const exampleResponse = new APIResponse<ValidIngredientGroup>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(500, exampleResponse);

    expect(client.getValidIngredientGroup(validIngredientGroupID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getValidIngredientGroups user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientGroup>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups`).reply(200, exampleResponse);

    expect(client.getValidIngredientGroups()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientGroups', () => {
    const expectedError = buildObligatoryError('getValidIngredientGroups service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientGroup>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups`).reply(500, exampleResponse);

    expect(client.getValidIngredientGroups()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidIngredientMeasurementUnit user error');
    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = fakeID();

    const expectedError = buildObligatoryError('getValidIngredientMeasurementUnit service error');
    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)).rejects.toEqual(
      expectedError.error,
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
    const expectedError = buildObligatoryError('getValidIngredientMeasurementUnits user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units`).reply(200, exampleResponse);

    expect(client.getValidIngredientMeasurementUnits()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientMeasurementUnits', () => {
    const expectedError = buildObligatoryError('getValidIngredientMeasurementUnits service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units`).reply(500, exampleResponse);

    expect(client.getValidIngredientMeasurementUnits()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidIngredientMeasurementUnitsByIngredient user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientMeasurementUnitsByIngredient(validIngredientID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientMeasurementUnitsByIngredient', () => {
    let validIngredientID = fakeID();

    const expectedError = buildObligatoryError('getValidIngredientMeasurementUnitsByIngredient service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/by_ingredient/${validIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientMeasurementUnitsByIngredient(validIngredientID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('getValidIngredientMeasurementUnitsByMeasurementUnit user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/by_measurement_unit/${validMeasurementUnitID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientMeasurementUnitsByMeasurementUnit(validMeasurementUnitID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientMeasurementUnitsByMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const expectedError = buildObligatoryError('getValidIngredientMeasurementUnitsByMeasurementUnit service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientMeasurementUnit>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/by_measurement_unit/${validMeasurementUnitID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientMeasurementUnitsByMeasurementUnit(validMeasurementUnitID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('getValidIngredientPreparation user error');
    const exampleResponse = new APIResponse<ValidIngredientPreparation>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientPreparation(validIngredientPreparationID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientPreparation', () => {
    let validIngredientPreparationID = fakeID();

    const expectedError = buildObligatoryError('getValidIngredientPreparation service error');
    const exampleResponse = new APIResponse<ValidIngredientPreparation>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientPreparation(validIngredientPreparationID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getValidIngredientPreparations user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_preparations`).reply(200, exampleResponse);

    expect(client.getValidIngredientPreparations()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientPreparations', () => {
    const expectedError = buildObligatoryError('getValidIngredientPreparations service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_preparations`).reply(500, exampleResponse);

    expect(client.getValidIngredientPreparations()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidIngredientPreparationsByIngredient user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientPreparationsByIngredient(validIngredientID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientPreparationsByIngredient', () => {
    let validIngredientID = fakeID();

    const expectedError = buildObligatoryError('getValidIngredientPreparationsByIngredient service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/by_ingredient/${validIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientPreparationsByIngredient(validIngredientID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidIngredientPreparationsByPreparation user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/by_preparation/${validPreparationID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientPreparationsByPreparation(validPreparationID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientPreparationsByPreparation', () => {
    let validPreparationID = fakeID();

    const expectedError = buildObligatoryError('getValidIngredientPreparationsByPreparation service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientPreparation>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/by_preparation/${validPreparationID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientPreparationsByPreparation(validPreparationID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidIngredientState user error');
    const exampleResponse = new APIResponse<ValidIngredientState>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(200, exampleResponse);

    expect(client.getValidIngredientState(validIngredientStateID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientState', () => {
    let validIngredientStateID = fakeID();

    const expectedError = buildObligatoryError('getValidIngredientState service error');
    const exampleResponse = new APIResponse<ValidIngredientState>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(500, exampleResponse);

    expect(client.getValidIngredientState(validIngredientStateID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidIngredientStateIngredient user error');
    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientStateIngredient(validIngredientStateIngredientID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = fakeID();

    const expectedError = buildObligatoryError('getValidIngredientStateIngredient service error');
    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientStateIngredient(validIngredientStateIngredientID)).rejects.toEqual(
      expectedError.error,
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
    const expectedError = buildObligatoryError('getValidIngredientStateIngredients user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients`).reply(200, exampleResponse);

    expect(client.getValidIngredientStateIngredients()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientStateIngredients', () => {
    const expectedError = buildObligatoryError('getValidIngredientStateIngredients service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients`).reply(500, exampleResponse);

    expect(client.getValidIngredientStateIngredients()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidIngredientStateIngredientsByIngredient user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientStateIngredientsByIngredient(validIngredientID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientStateIngredientsByIngredient', () => {
    let validIngredientID = fakeID();

    const expectedError = buildObligatoryError('getValidIngredientStateIngredientsByIngredient service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/by_ingredient/${validIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientStateIngredientsByIngredient(validIngredientID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('getValidIngredientStateIngredientsByIngredientState user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/${validIngredientStateID}`)
      .reply(200, exampleResponse);

    expect(client.getValidIngredientStateIngredientsByIngredientState(validIngredientStateID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientStateIngredientsByIngredientState', () => {
    let validIngredientStateID = fakeID();

    const expectedError = buildObligatoryError('getValidIngredientStateIngredientsByIngredientState service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientStateIngredient>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/${validIngredientStateID}`)
      .reply(500, exampleResponse);

    expect(client.getValidIngredientStateIngredientsByIngredientState(validIngredientStateID)).rejects.toEqual(
      expectedError.error,
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
    const expectedError = buildObligatoryError('getValidIngredientStates user error');
    const exampleResponse = new APIResponse<Array<ValidIngredientState>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states`).reply(200, exampleResponse);

    expect(client.getValidIngredientStates()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredientStates', () => {
    const expectedError = buildObligatoryError('getValidIngredientStates service error');
    const exampleResponse = new APIResponse<Array<ValidIngredientState>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states`).reply(500, exampleResponse);

    expect(client.getValidIngredientStates()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getValidIngredients user error');
    const exampleResponse = new APIResponse<Array<ValidIngredient>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredients`).reply(200, exampleResponse);

    expect(client.getValidIngredients()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidIngredients', () => {
    const expectedError = buildObligatoryError('getValidIngredients service error');
    const exampleResponse = new APIResponse<Array<ValidIngredient>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_ingredients`).reply(500, exampleResponse);

    expect(client.getValidIngredients()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidInstrument user error');
    const exampleResponse = new APIResponse<ValidInstrument>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(200, exampleResponse);

    expect(client.getValidInstrument(validInstrumentID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidInstrument', () => {
    let validInstrumentID = fakeID();

    const expectedError = buildObligatoryError('getValidInstrument service error');
    const exampleResponse = new APIResponse<ValidInstrument>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(500, exampleResponse);

    expect(client.getValidInstrument(validInstrumentID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getValidInstruments user error');
    const exampleResponse = new APIResponse<Array<ValidInstrument>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_instruments`).reply(200, exampleResponse);

    expect(client.getValidInstruments()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidInstruments', () => {
    const expectedError = buildObligatoryError('getValidInstruments service error');
    const exampleResponse = new APIResponse<Array<ValidInstrument>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_instruments`).reply(500, exampleResponse);

    expect(client.getValidInstruments()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidMeasurementUnit user error');
    const exampleResponse = new APIResponse<ValidMeasurementUnit>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(200, exampleResponse);

    expect(client.getValidMeasurementUnit(validMeasurementUnitID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const expectedError = buildObligatoryError('getValidMeasurementUnit service error');
    const exampleResponse = new APIResponse<ValidMeasurementUnit>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(500, exampleResponse);

    expect(client.getValidMeasurementUnit(validMeasurementUnitID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidMeasurementUnitConversion user error');
    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(200, exampleResponse);

    expect(client.getValidMeasurementUnitConversion(validMeasurementUnitConversionID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = fakeID();

    const expectedError = buildObligatoryError('getValidMeasurementUnitConversion service error');
    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(500, exampleResponse);

    expect(client.getValidMeasurementUnitConversion(validMeasurementUnitConversionID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('getValidMeasurementUnitConversionsFromUnit user error');
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnitConversion>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/from_unit/${validMeasurementUnitID}`)
      .reply(200, exampleResponse);

    expect(client.getValidMeasurementUnitConversionsFromUnit(validMeasurementUnitID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidMeasurementUnitConversionsFromUnit', () => {
    let validMeasurementUnitID = fakeID();

    const expectedError = buildObligatoryError('getValidMeasurementUnitConversionsFromUnit service error');
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnitConversion>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/from_unit/${validMeasurementUnitID}`)
      .reply(500, exampleResponse);

    expect(client.getValidMeasurementUnitConversionsFromUnit(validMeasurementUnitID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('getValidMeasurementUnitConversionsToUnit user error');
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnitConversion>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/to_unit/${validMeasurementUnitID}`)
      .reply(200, exampleResponse);

    expect(client.getValidMeasurementUnitConversionsToUnit(validMeasurementUnitID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a ValidMeasurementUnitConversionsToUnit', () => {
    let validMeasurementUnitID = fakeID();

    const expectedError = buildObligatoryError('getValidMeasurementUnitConversionsToUnit service error');
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnitConversion>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/to_unit/${validMeasurementUnitID}`)
      .reply(500, exampleResponse);

    expect(client.getValidMeasurementUnitConversionsToUnit(validMeasurementUnitID)).rejects.toEqual(
      expectedError.error,
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
    const expectedError = buildObligatoryError('getValidMeasurementUnits user error');
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units`).reply(200, exampleResponse);

    expect(client.getValidMeasurementUnits()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidMeasurementUnits', () => {
    const expectedError = buildObligatoryError('getValidMeasurementUnits service error');
    const exampleResponse = new APIResponse<Array<ValidMeasurementUnit>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units`).reply(500, exampleResponse);

    expect(client.getValidMeasurementUnits()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidPreparation user error');
    const exampleResponse = new APIResponse<ValidPreparation>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(200, exampleResponse);

    expect(client.getValidPreparation(validPreparationID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparation', () => {
    let validPreparationID = fakeID();

    const expectedError = buildObligatoryError('getValidPreparation service error');
    const exampleResponse = new APIResponse<ValidPreparation>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(500, exampleResponse);

    expect(client.getValidPreparation(validPreparationID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidPreparationInstrument user error');
    const exampleResponse = new APIResponse<ValidPreparationInstrument>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(200, exampleResponse);

    expect(client.getValidPreparationInstrument(validPreparationVesselID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationInstrument', () => {
    let validPreparationVesselID = fakeID();

    const expectedError = buildObligatoryError('getValidPreparationInstrument service error');
    const exampleResponse = new APIResponse<ValidPreparationInstrument>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(500, exampleResponse);

    expect(client.getValidPreparationInstrument(validPreparationVesselID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getValidPreparationInstruments user error');
    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparation_instruments`).reply(200, exampleResponse);

    expect(client.getValidPreparationInstruments()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationInstruments', () => {
    const expectedError = buildObligatoryError('getValidPreparationInstruments service error');
    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparation_instruments`).reply(500, exampleResponse);

    expect(client.getValidPreparationInstruments()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidPreparationInstrumentsByInstrument user error');
    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/by_instrument/${validInstrumentID}`)
      .reply(200, exampleResponse);

    expect(client.getValidPreparationInstrumentsByInstrument(validInstrumentID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationInstrumentsByInstrument', () => {
    let validInstrumentID = fakeID();

    const expectedError = buildObligatoryError('getValidPreparationInstrumentsByInstrument service error');
    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/by_instrument/${validInstrumentID}`)
      .reply(500, exampleResponse);

    expect(client.getValidPreparationInstrumentsByInstrument(validInstrumentID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidPreparationInstrumentsByPreparation user error');
    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/by_preparation/${validPreparationID}`)
      .reply(200, exampleResponse);

    expect(client.getValidPreparationInstrumentsByPreparation(validPreparationID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationInstrumentsByPreparation', () => {
    let validPreparationID = fakeID();

    const expectedError = buildObligatoryError('getValidPreparationInstrumentsByPreparation service error');
    const exampleResponse = new APIResponse<Array<ValidPreparationInstrument>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/by_preparation/${validPreparationID}`)
      .reply(500, exampleResponse);

    expect(client.getValidPreparationInstrumentsByPreparation(validPreparationID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidPreparationVessel user error');
    const exampleResponse = new APIResponse<ValidPreparationVessel>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(200, exampleResponse);

    expect(client.getValidPreparationVessel(validPreparationVesselID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationVessel', () => {
    let validPreparationVesselID = fakeID();

    const expectedError = buildObligatoryError('getValidPreparationVessel service error');
    const exampleResponse = new APIResponse<ValidPreparationVessel>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(500, exampleResponse);

    expect(client.getValidPreparationVessel(validPreparationVesselID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getValidPreparationVessels user error');
    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels`).reply(200, exampleResponse);

    expect(client.getValidPreparationVessels()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationVessels', () => {
    const expectedError = buildObligatoryError('getValidPreparationVessels service error');
    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels`).reply(500, exampleResponse);

    expect(client.getValidPreparationVessels()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidPreparationVesselsByPreparation user error');
    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_vessels/by_preparation/${validPreparationID}`)
      .reply(200, exampleResponse);

    expect(client.getValidPreparationVesselsByPreparation(validPreparationID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationVesselsByPreparation', () => {
    let validPreparationID = fakeID();

    const expectedError = buildObligatoryError('getValidPreparationVesselsByPreparation service error');
    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_vessels/by_preparation/${validPreparationID}`)
      .reply(500, exampleResponse);

    expect(client.getValidPreparationVesselsByPreparation(validPreparationID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidPreparationVesselsByVessel user error');
    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/by_vessel/${ValidVesselID}`).reply(200, exampleResponse);

    expect(client.getValidPreparationVesselsByVessel(ValidVesselID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparationVesselsByVessel', () => {
    let ValidVesselID = fakeID();

    const expectedError = buildObligatoryError('getValidPreparationVesselsByVessel service error');
    const exampleResponse = new APIResponse<Array<ValidPreparationVessel>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/by_vessel/${ValidVesselID}`).reply(500, exampleResponse);

    expect(client.getValidPreparationVesselsByVessel(ValidVesselID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getValidPreparations user error');
    const exampleResponse = new APIResponse<Array<ValidPreparation>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparations`).reply(200, exampleResponse);

    expect(client.getValidPreparations()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidPreparations', () => {
    const expectedError = buildObligatoryError('getValidPreparations service error');
    const exampleResponse = new APIResponse<Array<ValidPreparation>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_preparations`).reply(500, exampleResponse);

    expect(client.getValidPreparations()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getValidVessel user error');
    const exampleResponse = new APIResponse<ValidVessel>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(200, exampleResponse);

    expect(client.getValidVessel(validVesselID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidVessel', () => {
    let validVesselID = fakeID();

    const expectedError = buildObligatoryError('getValidVessel service error');
    const exampleResponse = new APIResponse<ValidVessel>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(500, exampleResponse);

    expect(client.getValidVessel(validVesselID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getValidVessels user error');
    const exampleResponse = new APIResponse<Array<ValidVessel>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_vessels`).reply(200, exampleResponse);

    expect(client.getValidVessels()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a ValidVessels', () => {
    const expectedError = buildObligatoryError('getValidVessels service error');
    const exampleResponse = new APIResponse<Array<ValidVessel>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/valid_vessels`).reply(500, exampleResponse);

    expect(client.getValidVessels()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getWebhook user error');
    const exampleResponse = new APIResponse<Webhook>(expectedError);
    mock.onGet(`${baseURL}/api/v1/webhooks/${webhookID}`).reply(200, exampleResponse);

    expect(client.getWebhook(webhookID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a Webhook', () => {
    let webhookID = fakeID();

    const expectedError = buildObligatoryError('getWebhook service error');
    const exampleResponse = new APIResponse<Webhook>(expectedError);
    mock.onGet(`${baseURL}/api/v1/webhooks/${webhookID}`).reply(500, exampleResponse);

    expect(client.getWebhook(webhookID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getWebhooks user error');
    const exampleResponse = new APIResponse<Array<Webhook>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/webhooks`).reply(200, exampleResponse);

    expect(client.getWebhooks()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a Webhooks', () => {
    const expectedError = buildObligatoryError('getWebhooks service error');
    const exampleResponse = new APIResponse<Array<Webhook>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/webhooks`).reply(500, exampleResponse);

    expect(client.getWebhooks()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getHouseholdInstrumentOwnership user error');
    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(200, exampleResponse);

    expect(client.getHouseholdInstrumentOwnership(householdInstrumentOwnershipID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a household instrument ownership', () => {
    let householdInstrumentOwnershipID = fakeID();

    const expectedError = buildObligatoryError('getHouseholdInstrumentOwnership service error');
    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(500, exampleResponse);

    expect(client.getHouseholdInstrumentOwnership(householdInstrumentOwnershipID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getHouseholdInvitationByID user error');
    const exampleResponse = new APIResponse<HouseholdInvitation>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/households/${householdID}/invitations/${householdInvitationID}`)
      .reply(200, exampleResponse);

    expect(client.getHouseholdInvitationByID(householdID, householdInvitationID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a household invitation by its ID', () => {
    let householdID = fakeID();
    let householdInvitationID = fakeID();

    const expectedError = buildObligatoryError('getHouseholdInvitationByID service error');
    const exampleResponse = new APIResponse<HouseholdInvitation>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/households/${householdID}/invitations/${householdInvitationID}`)
      .reply(500, exampleResponse);

    expect(client.getHouseholdInvitationByID(householdID, householdInvitationID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getHousehold user error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onGet(`${baseURL}/api/v1/households/${householdID}`).reply(200, exampleResponse);

    expect(client.getHousehold(householdID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a household', () => {
    let householdID = fakeID();

    const expectedError = buildObligatoryError('getHousehold service error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onGet(`${baseURL}/api/v1/households/${householdID}`).reply(500, exampleResponse);

    expect(client.getHousehold(householdID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getHouseholds user error');
    const exampleResponse = new APIResponse<Array<Household>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/households`).reply(200, exampleResponse);

    expect(client.getHouseholds()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a list of households', () => {
    const expectedError = buildObligatoryError('getHouseholds service error');
    const exampleResponse = new APIResponse<Array<Household>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/households`).reply(500, exampleResponse);

    expect(client.getHouseholds()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getMealPlanOption user error');
    const exampleResponse = new APIResponse<MealPlanOption>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(200, exampleResponse);

    expect(client.getMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should raise service errors appropriately when trying to fetch a meal plan option by its ID', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const expectedError = buildObligatoryError('getMealPlanOption service error');
    const exampleResponse = new APIResponse<MealPlanOption>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(500, exampleResponse);

    expect(client.getMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('getMealPlanOptionVote user error');
    const exampleResponse = new APIResponse<MealPlanOptionVote>(expectedError);
    mock
      .onGet(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(200, exampleResponse);

    expect(
      client.getMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID),
    ).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch a meal plan option vote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();
    let mealPlanOptionVoteID = fakeID();

    const expectedError = buildObligatoryError('getMealPlanOptionVote service error');
    const exampleResponse = new APIResponse<MealPlanOptionVote>(expectedError);
    mock
      .onGet(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(500, exampleResponse);

    expect(
      client.getMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID),
    ).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getMealPlanOptionVotes user error');
    const exampleResponse = new APIResponse<Array<MealPlanOptionVote>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes`)
      .reply(200, exampleResponse);

    expect(client.getMealPlanOptionVotes(mealPlanID, mealPlanEventID, mealPlanOptionID)).rejects.toEqual(
      expectedError.error,
    );
  });

  it("should raise service errors appropriately when trying to fetch a meal plan option's votes", () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const expectedError = buildObligatoryError('getMealPlanOptionVotes service error');
    const exampleResponse = new APIResponse<Array<MealPlanOptionVote>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes`)
      .reply(500, exampleResponse);

    expect(client.getMealPlanOptionVotes(mealPlanID, mealPlanEventID, mealPlanOptionID)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('getMealPlanOptions user error');
    const exampleResponse = new APIResponse<Array<MealPlanOption>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`)
      .reply(200, exampleResponse);

    expect(client.getMealPlanOptions(mealPlanID, mealPlanEventID)).rejects.toEqual(expectedError.error);
  });

  it("should raise service errors appropriately when trying to fetch a meal plan's options", () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const expectedError = buildObligatoryError('getMealPlanOptions service error');
    const exampleResponse = new APIResponse<Array<MealPlanOption>>(expectedError);
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`)
      .reply(500, exampleResponse);

    expect(client.getMealPlanOptions(mealPlanID, mealPlanEventID)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('getHouseholdInvitation user error');
    const exampleResponse = new APIResponse<HouseholdInvitation>(expectedError);
    mock.onGet(`${baseURL}/api/v1/household_invitations/${householdInvitationID}`).reply(200, exampleResponse);

    expect(client.getHouseholdInvitation(householdInvitationID)).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch household invitations for a given household', () => {
    let householdInvitationID = fakeID();

    const expectedError = buildObligatoryError('getHouseholdInvitation service error');
    const exampleResponse = new APIResponse<HouseholdInvitation>(expectedError);
    mock.onGet(`${baseURL}/api/v1/household_invitations/${householdInvitationID}`).reply(500, exampleResponse);

    expect(client.getHouseholdInvitation(householdInvitationID)).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getReceivedHouseholdInvitations user error');
    const exampleResponse = new APIResponse<Array<HouseholdInvitation>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/household_invitations/received`).reply(200, exampleResponse);

    expect(client.getReceivedHouseholdInvitations()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch received household invitations', () => {
    const expectedError = buildObligatoryError('getReceivedHouseholdInvitations service error');
    const exampleResponse = new APIResponse<Array<HouseholdInvitation>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/household_invitations/received`).reply(500, exampleResponse);

    expect(client.getReceivedHouseholdInvitations()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getSentHouseholdInvitations user error');
    const exampleResponse = new APIResponse<Array<HouseholdInvitation>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/household_invitations/sent`).reply(200, exampleResponse);

    expect(client.getSentHouseholdInvitations()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch sent household invitations', () => {
    const expectedError = buildObligatoryError('getSentHouseholdInvitations service error');
    const exampleResponse = new APIResponse<Array<HouseholdInvitation>>(expectedError);
    mock.onGet(`${baseURL}/api/v1/household_invitations/sent`).reply(500, exampleResponse);

    expect(client.getSentHouseholdInvitations()).rejects.toEqual(expectedError.error);
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
    const expectedError = buildObligatoryError('getActiveHousehold user error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onGet(`${baseURL}/api/v1/households/current`).reply(200, exampleResponse);

    expect(client.getActiveHousehold()).rejects.toEqual(expectedError.error);
  });

  it('should raise service errors appropriately when trying to fetch the currently active household', () => {
    const expectedError = buildObligatoryError('getActiveHousehold service error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onGet(`${baseURL}/api/v1/households/current`).reply(500, exampleResponse);

    expect(client.getActiveHousehold()).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('rejectHouseholdInvitation user error');
    const exampleResponse = new APIResponse<HouseholdInvitation>(expectedError);
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/reject`).reply(200, exampleResponse);

    expect(client.rejectHouseholdInvitation(householdInvitationID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during reject a received household invitation', () => {
    let householdInvitationID = fakeID();

    const exampleInput = new HouseholdInvitationUpdateRequestInput();

    const expectedError = buildObligatoryError('rejectHouseholdInvitation service error');
    const exampleResponse = new APIResponse<HouseholdInvitation>(expectedError);
    mock.onPut(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/reject`).reply(500, exampleResponse);

    expect(client.rejectHouseholdInvitation(householdInvitationID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateMealPlan user error');
    const exampleResponse = new APIResponse<MealPlan>(expectedError);
    mock.onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(200, exampleResponse);

    expect(client.updateMealPlan(mealPlanID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a MealPlan', () => {
    let mealPlanID = fakeID();

    const exampleInput = new MealPlanUpdateRequestInput();

    const expectedError = buildObligatoryError('updateMealPlan service error');
    const exampleResponse = new APIResponse<MealPlan>(expectedError);
    mock.onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(500, exampleResponse);

    expect(client.updateMealPlan(mealPlanID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateMealPlanEvent user error');
    const exampleResponse = new APIResponse<MealPlanEvent>(expectedError);
    mock.onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(200, exampleResponse);

    expect(client.updateMealPlanEvent(mealPlanID, mealPlanEventID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a MealPlanEvent', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();

    const exampleInput = new MealPlanEventUpdateRequestInput();

    const expectedError = buildObligatoryError('updateMealPlanEvent service error');
    const exampleResponse = new APIResponse<MealPlanEvent>(expectedError);
    mock.onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(500, exampleResponse);

    expect(client.updateMealPlanEvent(mealPlanID, mealPlanEventID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateMealPlanGroceryListItem user error');
    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(200, exampleResponse);

    expect(client.updateMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a MealPlanGroceryListItem', () => {
    let mealPlanID = fakeID();
    let mealPlanGroceryListItemID = fakeID();

    const exampleInput = new MealPlanGroceryListItemUpdateRequestInput();

    const expectedError = buildObligatoryError('updateMealPlanGroceryListItem service error');
    const exampleResponse = new APIResponse<MealPlanGroceryListItem>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(500, exampleResponse);

    expect(client.updateMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateMealPlanTaskStatus user error');
    const exampleResponse = new APIResponse<MealPlanTask>(expectedError);
    mock.onPatch(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`).reply(200, exampleResponse);

    expect(client.updateMealPlanTaskStatus(mealPlanID, mealPlanTaskID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a MealPlanTaskStatus', () => {
    let mealPlanID = fakeID();
    let mealPlanTaskID = fakeID();

    const exampleInput = new MealPlanTaskStatusChangeRequestInput();

    const expectedError = buildObligatoryError('updateMealPlanTaskStatus service error');
    const exampleResponse = new APIResponse<MealPlanTask>(expectedError);
    mock.onPatch(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`).reply(500, exampleResponse);

    expect(client.updateMealPlanTaskStatus(mealPlanID, mealPlanTaskID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updatePassword user error');
    const exampleResponse = new APIResponse<PasswordResetResponse>(expectedError);
    mock.onPut(`${baseURL}/api/v1/users/password/new`).reply(200, exampleResponse);

    expect(client.updatePassword(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a Password', () => {
    const exampleInput = new PasswordUpdateInput();

    const expectedError = buildObligatoryError('updatePassword service error');
    const exampleResponse = new APIResponse<PasswordResetResponse>(expectedError);
    mock.onPut(`${baseURL}/api/v1/users/password/new`).reply(500, exampleResponse);

    expect(client.updatePassword(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateRecipe user error');
    const exampleResponse = new APIResponse<Recipe>(expectedError);
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}`).reply(200, exampleResponse);

    expect(client.updateRecipe(recipeID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a Recipe', () => {
    let recipeID = fakeID();

    const exampleInput = new RecipeUpdateRequestInput();

    const expectedError = buildObligatoryError('updateRecipe service error');
    const exampleResponse = new APIResponse<Recipe>(expectedError);
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}`).reply(500, exampleResponse);

    expect(client.updateRecipe(recipeID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateRecipePrepTask user error');
    const exampleResponse = new APIResponse<RecipePrepTask>(expectedError);
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(200, exampleResponse);

    expect(client.updateRecipePrepTask(recipeID, recipePrepTaskID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a RecipePrepTask', () => {
    let recipeID = fakeID();
    let recipePrepTaskID = fakeID();

    const exampleInput = new RecipePrepTaskUpdateRequestInput();

    const expectedError = buildObligatoryError('updateRecipePrepTask service error');
    const exampleResponse = new APIResponse<RecipePrepTask>(expectedError);
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(500, exampleResponse);

    expect(client.updateRecipePrepTask(recipeID, recipePrepTaskID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateRecipeRating user error');
    const exampleResponse = new APIResponse<RecipeRating>(expectedError);
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(200, exampleResponse);

    expect(client.updateRecipeRating(recipeID, recipeRatingID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a RecipeRating', () => {
    let recipeID = fakeID();
    let recipeRatingID = fakeID();

    const exampleInput = new RecipeRatingUpdateRequestInput();

    const expectedError = buildObligatoryError('updateRecipeRating service error');
    const exampleResponse = new APIResponse<RecipeRating>(expectedError);
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(500, exampleResponse);

    expect(client.updateRecipeRating(recipeID, recipeRatingID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateRecipeStep user error');
    const exampleResponse = new APIResponse<RecipeStep>(expectedError);
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(200, exampleResponse);

    expect(client.updateRecipeStep(recipeID, recipeStepID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a RecipeStep', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();

    const exampleInput = new RecipeStepUpdateRequestInput();

    const expectedError = buildObligatoryError('updateRecipeStep service error');
    const exampleResponse = new APIResponse<RecipeStep>(expectedError);
    mock.onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(500, exampleResponse);

    expect(client.updateRecipeStep(recipeID, recipeStepID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateRecipeStepCompletionCondition user error');
    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(expectedError);
    mock
      .onPut(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(200, exampleResponse);

    expect(
      client.updateRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID, exampleInput),
    ).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a RecipeStepCompletionCondition', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepCompletionConditionID = fakeID();

    const exampleInput = new RecipeStepCompletionConditionUpdateRequestInput();

    const expectedError = buildObligatoryError('updateRecipeStepCompletionCondition service error');
    const exampleResponse = new APIResponse<RecipeStepCompletionCondition>(expectedError);
    mock
      .onPut(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(500, exampleResponse);

    expect(
      client.updateRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID, exampleInput),
    ).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateRecipeStepIngredient user error');
    const exampleResponse = new APIResponse<RecipeStepIngredient>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(200, exampleResponse);

    expect(
      client.updateRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID, exampleInput),
    ).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a RecipeStepIngredient', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepIngredientID = fakeID();

    const exampleInput = new RecipeStepIngredientUpdateRequestInput();

    const expectedError = buildObligatoryError('updateRecipeStepIngredient service error');
    const exampleResponse = new APIResponse<RecipeStepIngredient>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(500, exampleResponse);

    expect(
      client.updateRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID, exampleInput),
    ).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateRecipeStepInstrument user error');
    const exampleResponse = new APIResponse<RecipeStepInstrument>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(200, exampleResponse);

    expect(
      client.updateRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID, exampleInput),
    ).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a RecipeStepInstrument', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepInstrumentID = fakeID();

    const exampleInput = new RecipeStepInstrumentUpdateRequestInput();

    const expectedError = buildObligatoryError('updateRecipeStepInstrument service error');
    const exampleResponse = new APIResponse<RecipeStepInstrument>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(500, exampleResponse);

    expect(
      client.updateRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID, exampleInput),
    ).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateRecipeStepProduct user error');
    const exampleResponse = new APIResponse<RecipeStepProduct>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(200, exampleResponse);

    expect(client.updateRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a RecipeStepProduct', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepProductID = fakeID();

    const exampleInput = new RecipeStepProductUpdateRequestInput();

    const expectedError = buildObligatoryError('updateRecipeStepProduct service error');
    const exampleResponse = new APIResponse<RecipeStepProduct>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(500, exampleResponse);

    expect(client.updateRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateRecipeStepVessel user error');
    const exampleResponse = new APIResponse<RecipeStepVessel>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(200, exampleResponse);

    expect(client.updateRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a RecipeStepVessel', () => {
    let recipeID = fakeID();
    let recipeStepID = fakeID();
    let recipeStepVesselID = fakeID();

    const exampleInput = new RecipeStepVesselUpdateRequestInput();

    const expectedError = buildObligatoryError('updateRecipeStepVessel service error');
    const exampleResponse = new APIResponse<RecipeStepVessel>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(500, exampleResponse);

    expect(client.updateRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateServiceSettingConfiguration user error');
    const exampleResponse = new APIResponse<ServiceSettingConfiguration>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/settings/configurations/${serviceSettingConfigurationID}`)
      .reply(200, exampleResponse);

    expect(client.updateServiceSettingConfiguration(serviceSettingConfigurationID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a ServiceSettingConfiguration', () => {
    let serviceSettingConfigurationID = fakeID();

    const exampleInput = new ServiceSettingConfigurationUpdateRequestInput();

    const expectedError = buildObligatoryError('updateServiceSettingConfiguration service error');
    const exampleResponse = new APIResponse<ServiceSettingConfiguration>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/settings/configurations/${serviceSettingConfigurationID}`)
      .reply(500, exampleResponse);

    expect(client.updateServiceSettingConfiguration(serviceSettingConfigurationID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateUserDetails user error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPut(`${baseURL}/api/v1/users/details`).reply(200, exampleResponse);

    expect(client.updateUserDetails(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a UserDetails', () => {
    const exampleInput = new UserDetailsUpdateRequestInput();

    const expectedError = buildObligatoryError('updateUserDetails service error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPut(`${baseURL}/api/v1/users/details`).reply(500, exampleResponse);

    expect(client.updateUserDetails(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateUserEmailAddress user error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPut(`${baseURL}/api/v1/users/email_address`).reply(200, exampleResponse);

    expect(client.updateUserEmailAddress(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a UserEmailAddress', () => {
    const exampleInput = new UserEmailAddressUpdateInput();

    const expectedError = buildObligatoryError('updateUserEmailAddress service error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPut(`${baseURL}/api/v1/users/email_address`).reply(500, exampleResponse);

    expect(client.updateUserEmailAddress(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateUserIngredientPreference user error');
    const exampleResponse = new APIResponse<UserIngredientPreference>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`)
      .reply(200, exampleResponse);

    expect(client.updateUserIngredientPreference(userIngredientPreferenceID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a UserIngredientPreference', () => {
    let userIngredientPreferenceID = fakeID();

    const exampleInput = new UserIngredientPreferenceUpdateRequestInput();

    const expectedError = buildObligatoryError('updateUserIngredientPreference service error');
    const exampleResponse = new APIResponse<UserIngredientPreference>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`)
      .reply(500, exampleResponse);

    expect(client.updateUserIngredientPreference(userIngredientPreferenceID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateUserNotification user error');
    const exampleResponse = new APIResponse<UserNotification>(expectedError);
    mock.onPatch(`${baseURL}/api/v1/user_notifications/${userNotificationID}`).reply(200, exampleResponse);

    expect(client.updateUserNotification(userNotificationID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a UserNotification', () => {
    let userNotificationID = fakeID();

    const exampleInput = new UserNotificationUpdateRequestInput();

    const expectedError = buildObligatoryError('updateUserNotification service error');
    const exampleResponse = new APIResponse<UserNotification>(expectedError);
    mock.onPatch(`${baseURL}/api/v1/user_notifications/${userNotificationID}`).reply(500, exampleResponse);

    expect(client.updateUserNotification(userNotificationID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateUserUsername user error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPut(`${baseURL}/api/v1/users/username`).reply(200, exampleResponse);

    expect(client.updateUserUsername(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a UserUsername', () => {
    const exampleInput = new UsernameUpdateInput();

    const expectedError = buildObligatoryError('updateUserUsername service error');
    const exampleResponse = new APIResponse<User>(expectedError);
    mock.onPut(`${baseURL}/api/v1/users/username`).reply(500, exampleResponse);

    expect(client.updateUserUsername(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateValidIngredient user error');
    const exampleResponse = new APIResponse<ValidIngredient>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(200, exampleResponse);

    expect(client.updateValidIngredient(validIngredientID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a ValidIngredient', () => {
    let validIngredientID = fakeID();

    const exampleInput = new ValidIngredientUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidIngredient service error');
    const exampleResponse = new APIResponse<ValidIngredient>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(500, exampleResponse);

    expect(client.updateValidIngredient(validIngredientID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateValidIngredientGroup user error');
    const exampleResponse = new APIResponse<ValidIngredientGroup>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(200, exampleResponse);

    expect(client.updateValidIngredientGroup(validIngredientGroupID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidIngredientGroup', () => {
    let validIngredientGroupID = fakeID();

    const exampleInput = new ValidIngredientGroupUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidIngredientGroup service error');
    const exampleResponse = new APIResponse<ValidIngredientGroup>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(500, exampleResponse);

    expect(client.updateValidIngredientGroup(validIngredientGroupID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateValidIngredientMeasurementUnit user error');
    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(200, exampleResponse);

    expect(client.updateValidIngredientMeasurementUnit(validIngredientMeasurementUnitID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = fakeID();

    const exampleInput = new ValidIngredientMeasurementUnitUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidIngredientMeasurementUnit service error');
    const exampleResponse = new APIResponse<ValidIngredientMeasurementUnit>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(500, exampleResponse);

    expect(client.updateValidIngredientMeasurementUnit(validIngredientMeasurementUnitID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateValidIngredientPreparation user error');
    const exampleResponse = new APIResponse<ValidIngredientPreparation>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(200, exampleResponse);

    expect(client.updateValidIngredientPreparation(validIngredientPreparationID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidIngredientPreparation', () => {
    let validIngredientPreparationID = fakeID();

    const exampleInput = new ValidIngredientPreparationUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidIngredientPreparation service error');
    const exampleResponse = new APIResponse<ValidIngredientPreparation>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(500, exampleResponse);

    expect(client.updateValidIngredientPreparation(validIngredientPreparationID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateValidIngredientState user error');
    const exampleResponse = new APIResponse<ValidIngredientState>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(200, exampleResponse);

    expect(client.updateValidIngredientState(validIngredientStateID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidIngredientState', () => {
    let validIngredientStateID = fakeID();

    const exampleInput = new ValidIngredientStateUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidIngredientState service error');
    const exampleResponse = new APIResponse<ValidIngredientState>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(500, exampleResponse);

    expect(client.updateValidIngredientState(validIngredientStateID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateValidIngredientStateIngredient user error');
    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(200, exampleResponse);

    expect(client.updateValidIngredientStateIngredient(validIngredientStateIngredientID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = fakeID();

    const exampleInput = new ValidIngredientStateIngredientUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidIngredientStateIngredient service error');
    const exampleResponse = new APIResponse<ValidIngredientStateIngredient>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(500, exampleResponse);

    expect(client.updateValidIngredientStateIngredient(validIngredientStateIngredientID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateValidInstrument user error');
    const exampleResponse = new APIResponse<ValidInstrument>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(200, exampleResponse);

    expect(client.updateValidInstrument(validInstrumentID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a ValidInstrument', () => {
    let validInstrumentID = fakeID();

    const exampleInput = new ValidInstrumentUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidInstrument service error');
    const exampleResponse = new APIResponse<ValidInstrument>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(500, exampleResponse);

    expect(client.updateValidInstrument(validInstrumentID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateValidMeasurementUnit user error');
    const exampleResponse = new APIResponse<ValidMeasurementUnit>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(200, exampleResponse);

    expect(client.updateValidMeasurementUnit(validMeasurementUnitID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidMeasurementUnit', () => {
    let validMeasurementUnitID = fakeID();

    const exampleInput = new ValidMeasurementUnitUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidMeasurementUnit service error');
    const exampleResponse = new APIResponse<ValidMeasurementUnit>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(500, exampleResponse);

    expect(client.updateValidMeasurementUnit(validMeasurementUnitID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateValidMeasurementUnitConversion user error');
    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(200, exampleResponse);

    expect(client.updateValidMeasurementUnitConversion(validMeasurementUnitConversionID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = fakeID();

    const exampleInput = new ValidMeasurementUnitConversionUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidMeasurementUnitConversion service error');
    const exampleResponse = new APIResponse<ValidMeasurementUnitConversion>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(500, exampleResponse);

    expect(client.updateValidMeasurementUnitConversion(validMeasurementUnitConversionID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateValidPreparation user error');
    const exampleResponse = new APIResponse<ValidPreparation>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(200, exampleResponse);

    expect(client.updateValidPreparation(validPreparationID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a ValidPreparation', () => {
    let validPreparationID = fakeID();

    const exampleInput = new ValidPreparationUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidPreparation service error');
    const exampleResponse = new APIResponse<ValidPreparation>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(500, exampleResponse);

    expect(client.updateValidPreparation(validPreparationID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateValidPreparationInstrument user error');
    const exampleResponse = new APIResponse<ValidPreparationInstrument>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(200, exampleResponse);

    expect(client.updateValidPreparationInstrument(validPreparationVesselID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidPreparationInstrument', () => {
    let validPreparationVesselID = fakeID();

    const exampleInput = new ValidPreparationInstrumentUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidPreparationInstrument service error');
    const exampleResponse = new APIResponse<ValidPreparationInstrument>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(500, exampleResponse);

    expect(client.updateValidPreparationInstrument(validPreparationVesselID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateValidPreparationVessel user error');
    const exampleResponse = new APIResponse<ValidPreparationVessel>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(200, exampleResponse);

    expect(client.updateValidPreparationVessel(validPreparationVesselID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a ValidPreparationVessel', () => {
    let validPreparationVesselID = fakeID();

    const exampleInput = new ValidPreparationVesselUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidPreparationVessel service error');
    const exampleResponse = new APIResponse<ValidPreparationVessel>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(500, exampleResponse);

    expect(client.updateValidPreparationVessel(validPreparationVesselID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateValidVessel user error');
    const exampleResponse = new APIResponse<ValidVessel>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(200, exampleResponse);

    expect(client.updateValidVessel(validVesselID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a ValidVessel', () => {
    let validVesselID = fakeID();

    const exampleInput = new ValidVesselUpdateRequestInput();

    const expectedError = buildObligatoryError('updateValidVessel service error');
    const exampleResponse = new APIResponse<ValidVessel>(expectedError);
    mock.onPut(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(500, exampleResponse);

    expect(client.updateValidVessel(validVesselID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateHouseholdInstrumentOwnership user error');
    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(200, exampleResponse);

    expect(client.updateHouseholdInstrumentOwnership(householdInstrumentOwnershipID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a household instrument ownership', () => {
    let householdInstrumentOwnershipID = fakeID();

    const exampleInput = new HouseholdInstrumentOwnershipUpdateRequestInput();

    const expectedError = buildObligatoryError('updateHouseholdInstrumentOwnership service error');
    const exampleResponse = new APIResponse<HouseholdInstrumentOwnership>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(500, exampleResponse);

    expect(client.updateHouseholdInstrumentOwnership(householdInstrumentOwnershipID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('updateHousehold user error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onPut(`${baseURL}/api/v1/households/${householdID}`).reply(200, exampleResponse);

    expect(client.updateHousehold(householdID, exampleInput)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a household', () => {
    let householdID = fakeID();

    const exampleInput = new HouseholdUpdateRequestInput();

    const expectedError = buildObligatoryError('updateHousehold service error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onPut(`${baseURL}/api/v1/households/${householdID}`).reply(500, exampleResponse);

    expect(client.updateHousehold(householdID, exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateMealPlanOptionVote user error');
    const exampleResponse = new APIResponse<MealPlanOptionVote>(expectedError);
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
    ).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update a meal plan option vote', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();
    let mealPlanOptionVoteID = fakeID();

    const exampleInput = new MealPlanOptionVoteUpdateRequestInput();

    const expectedError = buildObligatoryError('updateMealPlanOptionVote service error');
    const exampleResponse = new APIResponse<MealPlanOptionVote>(expectedError);
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
    ).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('updateMealPlanOption user error');
    const exampleResponse = new APIResponse<MealPlanOption>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(200, exampleResponse);

    expect(client.updateMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID, exampleInput)).rejects.toEqual(
      expectedError.error,
    );
  });

  it('should appropriately raise service errors when they occur during update a meal plan option', () => {
    let mealPlanID = fakeID();
    let mealPlanEventID = fakeID();
    let mealPlanOptionID = fakeID();

    const exampleInput = new MealPlanOptionUpdateRequestInput();

    const expectedError = buildObligatoryError('updateMealPlanOption service error');
    const exampleResponse = new APIResponse<MealPlanOption>(expectedError);
    mock
      .onPut(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(500, exampleResponse);

    expect(client.updateMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID, exampleInput)).rejects.toEqual(
      expectedError.error,
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

    const expectedError = buildObligatoryError('adminUpdateUserStatus user error');
    const exampleResponse = new APIResponse<UserStatusResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/admin/users/status`).reply(201, exampleResponse);

    expect(client.adminUpdateUserStatus(exampleInput)).rejects.toEqual(expectedError.error);
  });

  it("should appropriately raise service errors when they occur during update a user's account status", () => {
    const exampleInput = new UserAccountStatusUpdateInput();

    const expectedError = buildObligatoryError('adminUpdateUserStatus service error');
    const exampleResponse = new APIResponse<UserStatusResponse>(expectedError);
    mock.onPost(`${baseURL}/api/v1/admin/users/status`).reply(500, exampleResponse);

    expect(client.adminUpdateUserStatus(exampleInput)).rejects.toEqual(expectedError.error);
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

    const expectedError = buildObligatoryError('setDefaultHousehold user error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/default`).reply(201, exampleResponse);

    expect(client.setDefaultHousehold(householdID)).rejects.toEqual(expectedError.error);
  });

  it('should appropriately raise service errors when they occur during update the default household assigned at login', () => {
    let householdID = fakeID();

    const expectedError = buildObligatoryError('setDefaultHousehold service error');
    const exampleResponse = new APIResponse<Household>(expectedError);
    mock.onPost(`${baseURL}/api/v1/households/${householdID}/default`).reply(500, exampleResponse);

    expect(client.setDefaultHousehold(householdID)).rejects.toEqual(expectedError.error);
  });
});
