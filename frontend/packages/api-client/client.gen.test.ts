// GENERATED CODE, DO NOT EDIT MANUALLY

import axios from 'axios';
import AxiosMockAdapter from 'axios-mock-adapter';
import { faker } from '@faker-js/faker';

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

import { DinnerDoneBetterAPIClient } from './client.gen';

const mock = new AxiosMockAdapter(axios);
const baseURL = 'http://things.stuff';
const client = new DinnerDoneBetterAPIClient(baseURL, 'test-token');

beforeEach(() => mock.reset());

describe('basic', () => {
  it('should Operation for fetching MealPlanTask', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanTaskID = faker.string.uuid();

    const exampleResponse = new MealPlanTask({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`).reply(200, exampleResponse);

    client.getMealPlanTask(mealPlanID, mealPlanTaskID).then((response: APIResponse<MealPlanTask>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidPreparationInstrument', () => {
    const exampleResponse = new ValidPreparationInstrument({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparation_instruments`).reply(200, exampleResponse);

    client.getValidPreparationInstruments().then((response: APIResponse<Array<ValidPreparationInstrument>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredient', () => {
    const exampleResponse = new ValidIngredient({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/random`).reply(200, exampleResponse);

    client.getRandomValidIngredient().then((response: APIResponse<ValidIngredient>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating WebhookTriggerEvent', () => {
    let webhookID = faker.string.uuid();

    const exampleResponse = new WebhookTriggerEvent({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/webhooks/${webhookID}/trigger_events`).reply(200, exampleResponse);

    client.createWebhookTriggerEvent(webhookID).then((response: APIResponse<WebhookTriggerEvent>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching User', () => {
    let q = faker.string.uuid();

    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users/search`).reply(200, exampleResponse);

    client.searchForUsers(q).then((response: APIResponse<Array<User>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidVessel', () => {
    let q = faker.string.uuid();

    const exampleResponse = new ValidVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_vessels/search`).reply(200, exampleResponse);

    client.searchForValidVessels(q).then((response: APIResponse<Array<ValidVessel>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating User', () => {
    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users/email_address_verification`).reply(200, exampleResponse);

    client.requestEmailVerificationEmail().then((response: APIResponse<User>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching MealPlanGroceryListItem', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanGroceryListItemID = faker.string.uuid();

    const exampleResponse = new MealPlanGroceryListItem({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(200, exampleResponse);

    client
      .getMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)
      .then((response: APIResponse<MealPlanGroceryListItem>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for archiving User', () => {
    let userID = faker.string.uuid();

    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users/${userID}`).reply(200, exampleResponse);

    client.archiveUser(userID).then((response: APIResponse<User>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating ValidInstrument', () => {
    const exampleResponse = new ValidInstrument({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_instruments`).reply(200, exampleResponse);

    client.createValidInstrument().then((response: APIResponse<ValidInstrument>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving WebhookTriggerEvent', () => {
    let webhookID = faker.string.uuid();
    let webhookTriggerEventID = faker.string.uuid();

    const exampleResponse = new WebhookTriggerEvent({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/webhooks/${webhookID}/trigger_events/${webhookTriggerEventID}`)
      .reply(200, exampleResponse);

    client
      .archiveWebhookTriggerEvent(webhookID, webhookTriggerEventID)
      .then((response: APIResponse<WebhookTriggerEvent>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating RecipeStepIngredient', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepIngredientID = faker.string.uuid();

    const exampleResponse = new RecipeStepIngredient({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(200, exampleResponse);

    client
      .updateRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)
      .then((response: APIResponse<RecipeStepIngredient>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating FinalizeMealPlansResponse', () => {
    const exampleResponse = new FinalizeMealPlansResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/workers/finalize_meal_plans`).reply(200, exampleResponse);

    client.runFinalizeMealPlanWorker().then((response: APIResponse<FinalizeMealPlansResponse>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredient', () => {
    let q = faker.string.uuid();

    const exampleResponse = new ValidIngredient({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/search`).reply(200, exampleResponse);

    client.searchForValidIngredients(q).then((response: APIResponse<Array<ValidIngredient>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching RecipeStepInstrument', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepInstrumentID = faker.string.uuid();

    const exampleResponse = new RecipeStepInstrument({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(200, exampleResponse);

    client
      .getRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)
      .then((response: APIResponse<RecipeStepInstrument>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidIngredient', () => {
    const exampleResponse = new ValidIngredient({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredients`).reply(200, exampleResponse);

    client.getValidIngredients().then((response: APIResponse<Array<ValidIngredient>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating Recipe', () => {
    let recipeID = faker.string.uuid();

    const exampleResponse = new Recipe({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}`).reply(200, exampleResponse);

    client.updateRecipe(recipeID).then((response: APIResponse<Recipe>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating MealPlanEvent', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();

    const exampleResponse = new MealPlanEvent({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(200, exampleResponse);

    client.updateMealPlanEvent(mealPlanID, mealPlanEventID).then((response: APIResponse<MealPlanEvent>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving Meal', () => {
    let mealID = faker.string.uuid();

    const exampleResponse = new Meal({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meals/${mealID}`).reply(200, exampleResponse);

    client.archiveMeal(mealID).then((response: APIResponse<Meal>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating UserStatusResponse', () => {
    const exampleResponse = new UserStatusResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/admin/users/status`).reply(200, exampleResponse);

    client.adminUpdateUserStatus().then((response: APIResponse<UserStatusResponse>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating MealPlanGroceryListItem', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanGroceryListItemID = faker.string.uuid();

    const exampleResponse = new MealPlanGroceryListItem({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(200, exampleResponse);

    client
      .updateMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)
      .then((response: APIResponse<MealPlanGroceryListItem>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating OAuth2ClientCreationResponse', () => {
    const exampleResponse = new OAuth2ClientCreationResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/oauth2_clients`).reply(200, exampleResponse);

    client.createOAuth2Client().then((response: APIResponse<OAuth2ClientCreationResponse>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving MealPlanEvent', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();

    const exampleResponse = new MealPlanEvent({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(200, exampleResponse);

    client.archiveMealPlanEvent(mealPlanID, mealPlanEventID).then((response: APIResponse<MealPlanEvent>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidVessel', () => {
    let validVesselID = faker.string.uuid();

    const exampleResponse = new ValidVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(200, exampleResponse);

    client.getValidVessel(validVesselID).then((response: APIResponse<ValidVessel>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating RecipeStepVessel', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStepVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`).reply(200, exampleResponse);

    client.createRecipeStepVessel(recipeID, recipeStepID).then((response: APIResponse<RecipeStepVessel>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching HouseholdInstrumentOwnership', () => {
    const exampleResponse = new HouseholdInstrumentOwnership({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households/instruments`).reply(200, exampleResponse);

    client.getHouseholdInstrumentOwnerships().then((response: APIResponse<Array<HouseholdInstrumentOwnership>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating MealPlanEvent', () => {
    let mealPlanID = faker.string.uuid();

    const exampleResponse = new MealPlanEvent({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events`).reply(200, exampleResponse);

    client.createMealPlanEvent(mealPlanID).then((response: APIResponse<MealPlanEvent>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching MealPlanOption', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();
    let mealPlanOptionID = faker.string.uuid();

    const exampleResponse = new MealPlanOption({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(200, exampleResponse);

    client
      .getMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)
      .then((response: APIResponse<MealPlanOption>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating MealPlanOption', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();
    let mealPlanOptionID = faker.string.uuid();

    const exampleResponse = new MealPlanOption({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(200, exampleResponse);

    client
      .updateMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)
      .then((response: APIResponse<MealPlanOption>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidIngredientPreparation', () => {
    let validPreparationID = faker.string.uuid();

    const exampleResponse = new ValidIngredientPreparation({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/by_preparation/${validPreparationID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientPreparationsByPreparation(validPreparationID)
      .then((response: APIResponse<Array<ValidIngredientPreparation>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidPreparationVessel', () => {
    let validPreparationVesselID = faker.string.uuid();

    const exampleResponse = new ValidPreparationVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(200, exampleResponse);

    client.getValidPreparationVessel(validPreparationVesselID).then((response: APIResponse<ValidPreparationVessel>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching AuditLogEntry', () => {
    let auditLogEntryID = faker.string.uuid();

    const exampleResponse = new AuditLogEntry({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/${auditLogEntryID}`).reply(200, exampleResponse);

    client.getAuditLogEntryByID(auditLogEntryID).then((response: APIResponse<AuditLogEntry>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching Meal', () => {
    let q = faker.string.uuid();

    const exampleResponse = new Meal({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meals/search`).reply(200, exampleResponse);

    client.searchForMeals(q).then((response: APIResponse<Array<Meal>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving Recipe', () => {
    let recipeID = faker.string.uuid();

    const exampleResponse = new Recipe({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}`).reply(200, exampleResponse);

    client.archiveRecipe(recipeID).then((response: APIResponse<Recipe>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating MealPlan', () => {
    let mealPlanID = faker.string.uuid();

    const exampleResponse = new MealPlan({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(200, exampleResponse);

    client.updateMealPlan(mealPlanID).then((response: APIResponse<MealPlan>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching OAuth2Client', () => {
    let oauth2ClientID = faker.string.uuid();

    const exampleResponse = new OAuth2Client({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/oauth2_clients/${oauth2ClientID}`).reply(200, exampleResponse);

    client.getOAuth2Client(oauth2ClientID).then((response: APIResponse<OAuth2Client>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating ValidIngredientPreparation', () => {
    let validIngredientPreparationID = faker.string.uuid();

    const exampleResponse = new ValidIngredientPreparation({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(200, exampleResponse);

    client
      .updateValidIngredientPreparation(validIngredientPreparationID)
      .then((response: APIResponse<ValidIngredientPreparation>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating ValidIngredientState', () => {
    const exampleResponse = new ValidIngredientState({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states`).reply(200, exampleResponse);

    client.createValidIngredientState().then((response: APIResponse<ValidIngredientState>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching RecipeStepVessel', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepVesselID = faker.string.uuid();

    const exampleResponse = new RecipeStepVessel({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(200, exampleResponse);

    client
      .getRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)
      .then((response: APIResponse<RecipeStepVessel>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating TOTPSecretRefreshResponse', () => {
    const exampleResponse = new TOTPSecretRefreshResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users/totp_secret/new`).reply(200, exampleResponse);

    client.refreshTOTPSecret().then((response: APIResponse<TOTPSecretRefreshResponse>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving HouseholdInstrumentOwnership', () => {
    let householdInstrumentOwnershipID = faker.string.uuid();

    const exampleResponse = new HouseholdInstrumentOwnership({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(200, exampleResponse);

    client
      .archiveHouseholdInstrumentOwnership(householdInstrumentOwnershipID)
      .then((response: APIResponse<HouseholdInstrumentOwnership>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = faker.string.uuid();

    const exampleResponse = new ValidMeasurementUnitConversion({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(200, exampleResponse);

    client
      .getValidMeasurementUnitConversion(validMeasurementUnitConversionID)
      .then((response: APIResponse<ValidMeasurementUnitConversion>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching AuditLogEntry', () => {
    const exampleResponse = new AuditLogEntry({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/for_user`).reply(200, exampleResponse);

    client.getAuditLogEntriesForUser().then((response: APIResponse<Array<AuditLogEntry>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating RecipeRating', () => {
    let recipeID = faker.string.uuid();

    const exampleResponse = new RecipeRating({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings`).reply(200, exampleResponse);

    client.createRecipeRating(recipeID).then((response: APIResponse<RecipeRating>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = faker.string.uuid();

    const exampleResponse = new ValidIngredientStateIngredient({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(200, exampleResponse);

    client
      .updateValidIngredientStateIngredient(validIngredientStateIngredientID)
      .then((response: APIResponse<ValidIngredientStateIngredient>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for archiving Household', () => {
    let householdID = faker.string.uuid();

    const exampleResponse = new Household({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households/${householdID}`).reply(200, exampleResponse);

    client.archiveHousehold(householdID).then((response: APIResponse<Household>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching Recipe', () => {
    const exampleResponse = new Recipe({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes`).reply(200, exampleResponse);

    client.getRecipes().then((response: APIResponse<Array<Recipe>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating UserIngredientPreference', () => {
    const exampleResponse = new UserIngredientPreference({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/user_ingredient_preferences`).reply(200, exampleResponse);

    client.createUserIngredientPreference().then((response: APIResponse<Array<UserIngredientPreference>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitID = faker.string.uuid();

    const exampleResponse = new ValidMeasurementUnitConversion({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/from_unit/${validMeasurementUnitID}`)
      .reply(200, exampleResponse);

    client
      .getValidMeasurementUnitConversionsFromUnit(validMeasurementUnitID)
      .then((response: APIResponse<Array<ValidMeasurementUnitConversion>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating ServiceSetting', () => {
    const exampleResponse = new ServiceSetting({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/settings`).reply(200, exampleResponse);

    client.createServiceSetting().then((response: APIResponse<ServiceSetting>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating HouseholdInvitation', () => {
    let householdID = faker.string.uuid();

    const exampleResponse = new HouseholdInvitation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households/${householdID}/invite`).reply(200, exampleResponse);

    client.createHouseholdInvitation(householdID).then((response: APIResponse<HouseholdInvitation>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating ValidMeasurementUnit', () => {
    const exampleResponse = new ValidMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units`).reply(200, exampleResponse);

    client.createValidMeasurementUnit().then((response: APIResponse<ValidMeasurementUnit>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating ValidIngredientGroup', () => {
    let validIngredientGroupID = faker.string.uuid();

    const exampleResponse = new ValidIngredientGroup({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(200, exampleResponse);

    client.updateValidIngredientGroup(validIngredientGroupID).then((response: APIResponse<ValidIngredientGroup>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ServiceSettingConfiguration', () => {
    const exampleResponse = new ServiceSettingConfiguration({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/settings/configurations/household`).reply(200, exampleResponse);

    client
      .getServiceSettingConfigurationsForHousehold()
      .then((response: APIResponse<Array<ServiceSettingConfiguration>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating RecipeStepCompletionCondition', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepCompletionConditionID = faker.string.uuid();

    const exampleResponse = new RecipeStepCompletionCondition({
      id: faker.string.uuid(),
    });
    mock
      .onGet(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(200, exampleResponse);

    client
      .updateRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID)
      .then((response: APIResponse<RecipeStepCompletionCondition>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching RecipeStepIngredient', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStepIngredient({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`).reply(200, exampleResponse);

    client
      .getRecipeStepIngredients(recipeID, recipeStepID)
      .then((response: APIResponse<Array<RecipeStepIngredient>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching User', () => {
    let userID = faker.string.uuid();

    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users/${userID}`).reply(200, exampleResponse);

    client.getUser(userID).then((response: APIResponse<User>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating HouseholdInvitation', () => {
    let householdInvitationID = faker.string.uuid();

    const exampleResponse = new HouseholdInvitation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/reject`).reply(200, exampleResponse);

    client.rejectHouseholdInvitation(householdInvitationID).then((response: APIResponse<HouseholdInvitation>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating Webhook', () => {
    const exampleResponse = new Webhook({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/webhooks`).reply(200, exampleResponse);

    client.createWebhook().then((response: APIResponse<Webhook>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching Recipe', () => {
    let recipeID = faker.string.uuid();

    const exampleResponse = new Recipe({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}`).reply(200, exampleResponse);

    client.getRecipe(recipeID).then((response: APIResponse<Recipe>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating UserNotification', () => {
    let userNotificationID = faker.string.uuid();

    const exampleResponse = new UserNotification({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/user_notifications/${userNotificationID}`).reply(200, exampleResponse);

    client.updateUserNotification(userNotificationID).then((response: APIResponse<UserNotification>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching RecipePrepTaskStep', () => {
    let recipeID = faker.string.uuid();

    const exampleResponse = new RecipePrepTaskStep({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_steps`).reply(200, exampleResponse);

    client.getRecipeMealPlanTasks(recipeID).then((response: APIResponse<RecipePrepTaskStep>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidInstrument', () => {
    const exampleResponse = new ValidInstrument({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_instruments/random`).reply(200, exampleResponse);

    client.getRandomValidInstrument().then((response: APIResponse<ValidInstrument>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ServiceSetting', () => {
    const exampleResponse = new ServiceSetting({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/settings`).reply(200, exampleResponse);

    client.getServiceSettings().then((response: APIResponse<Array<ServiceSetting>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching User', () => {
    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users/self`).reply(200, exampleResponse);

    client.getSelf().then((response: APIResponse<User>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating ValidInstrument', () => {
    let validInstrumentID = faker.string.uuid();

    const exampleResponse = new ValidInstrument({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(200, exampleResponse);

    client.updateValidInstrument(validInstrumentID).then((response: APIResponse<ValidInstrument>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating ValidVessel', () => {
    const exampleResponse = new ValidVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_vessels`).reply(200, exampleResponse);

    client.createValidVessel().then((response: APIResponse<ValidVessel>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating ValidPreparation', () => {
    const exampleResponse = new ValidPreparation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparations`).reply(200, exampleResponse);

    client.createValidPreparation().then((response: APIResponse<ValidPreparation>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredientStateIngredient', () => {
    let validIngredientStateID = faker.string.uuid();

    const exampleResponse = new ValidIngredientStateIngredient({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/by_ingredient_state/${validIngredientStateID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientStateIngredientsByIngredientState(validIngredientStateID)
      .then((response: APIResponse<Array<ValidIngredientStateIngredient>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ServiceSettingConfiguration', () => {
    let serviceSettingConfigurationName = faker.string.uuid();

    const exampleResponse = new ServiceSettingConfiguration({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/settings/configurations/user/${serviceSettingConfigurationName}`)
      .reply(200, exampleResponse);

    client
      .getServiceSettingConfigurationByName(serviceSettingConfigurationName)
      .then((response: APIResponse<Array<ServiceSettingConfiguration>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching HouseholdInvitation', () => {
    const exampleResponse = new HouseholdInvitation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/household_invitations/received`).reply(200, exampleResponse);

    client.getReceivedHouseholdInvitations().then((response: APIResponse<Array<HouseholdInvitation>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching HouseholdInvitation', () => {
    let householdID = faker.string.uuid();
    let householdInvitationID = faker.string.uuid();

    const exampleResponse = new HouseholdInvitation({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/households/${householdID}/invitations/${householdInvitationID}`)
      .reply(200, exampleResponse);

    client
      .getHouseholdInvitationByID(householdID, householdInvitationID)
      .then((response: APIResponse<HouseholdInvitation>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for archiving RecipePrepTask', () => {
    let recipeID = faker.string.uuid();
    let recipePrepTaskID = faker.string.uuid();

    const exampleResponse = new RecipePrepTask({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(200, exampleResponse);

    client.archiveRecipePrepTask(recipeID, recipePrepTaskID).then((response: APIResponse<RecipePrepTask>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating ValidIngredientMeasurementUnit', () => {
    const exampleResponse = new ValidIngredientMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units`).reply(200, exampleResponse);

    client.createValidIngredientMeasurementUnit().then((response: APIResponse<ValidIngredientMeasurementUnit>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredientStateIngredient', () => {
    let validIngredientID = faker.string.uuid();

    const exampleResponse = new ValidIngredientStateIngredient({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientStateIngredientsByIngredient(validIngredientID)
      .then((response: APIResponse<Array<ValidIngredientStateIngredient>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ServiceSetting', () => {
    let q = faker.string.uuid();

    const exampleResponse = new ServiceSetting({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/settings/search`).reply(200, exampleResponse);

    client.searchForServiceSettings(q).then((response: APIResponse<Array<ServiceSetting>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating UserPermissionsResponse', () => {
    const exampleResponse = new UserPermissionsResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users/permissions/check`).reply(200, exampleResponse);

    client.checkPermissions().then((response: APIResponse<UserPermissionsResponse>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = faker.string.uuid();

    const exampleResponse = new ValidIngredientMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)
      .then((response: APIResponse<ValidIngredientMeasurementUnit>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidInstrument', () => {
    const exampleResponse = new ValidInstrument({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_instruments`).reply(200, exampleResponse);

    client.getValidInstruments().then((response: APIResponse<Array<ValidInstrument>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating ValidPreparationInstrument', () => {
    let validPreparationVesselID = faker.string.uuid();

    const exampleResponse = new ValidPreparationInstrument({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(200, exampleResponse);

    client
      .updateValidPreparationInstrument(validPreparationVesselID)
      .then((response: APIResponse<ValidPreparationInstrument>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidMeasurementUnit', () => {
    const exampleResponse = new ValidMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units`).reply(200, exampleResponse);

    client.getValidMeasurementUnits().then((response: APIResponse<Array<ValidMeasurementUnit>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching Household', () => {
    const exampleResponse = new Household({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households/current`).reply(200, exampleResponse);

    client.getActiveHousehold().then((response: APIResponse<Household>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating MealPlanTask', () => {
    let mealPlanID = faker.string.uuid();

    const exampleResponse = new MealPlanTask({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks`).reply(200, exampleResponse);

    client.createMealPlanTask(mealPlanID).then((response: APIResponse<MealPlanTask>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating Household', () => {
    let householdID = faker.string.uuid();

    const exampleResponse = new Household({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households/${householdID}/transfer`).reply(200, exampleResponse);

    client.transferHouseholdOwnership(householdID).then((response: APIResponse<Household>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching UserNotification', () => {
    let userNotificationID = faker.string.uuid();

    const exampleResponse = new UserNotification({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/user_notifications/${userNotificationID}`).reply(200, exampleResponse);

    client.getUserNotification(userNotificationID).then((response: APIResponse<UserNotification>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidInstrument', () => {
    let q = faker.string.uuid();

    const exampleResponse = new ValidInstrument({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_instruments/search`).reply(200, exampleResponse);

    client.searchForValidInstruments(q).then((response: APIResponse<Array<ValidInstrument>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating User', () => {
    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/users/username/reminder`).reply(200, exampleResponse);

    client.requestUsernameReminder().then((response: APIResponse<User>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ValidIngredient', () => {
    let validIngredientID = faker.string.uuid();

    const exampleResponse = new ValidIngredient({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(200, exampleResponse);

    client.archiveValidIngredient(validIngredientID).then((response: APIResponse<ValidIngredient>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating MealPlanTask', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanTaskID = faker.string.uuid();

    const exampleResponse = new MealPlanTask({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks/${mealPlanTaskID}`).reply(200, exampleResponse);

    client.updateMealPlanTaskStatus(mealPlanID, mealPlanTaskID).then((response: APIResponse<MealPlanTask>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching Household', () => {
    const exampleResponse = new Household({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households`).reply(200, exampleResponse);

    client.getHouseholds().then((response: APIResponse<Array<Household>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching MealPlanOption', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();

    const exampleResponse = new MealPlanOption({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`)
      .reply(200, exampleResponse);

    client.getMealPlanOptions(mealPlanID, mealPlanEventID).then((response: APIResponse<Array<MealPlanOption>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidPreparationInstrument', () => {
    let validPreparationVesselID = faker.string.uuid();

    const exampleResponse = new ValidPreparationInstrument({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(200, exampleResponse);

    client
      .getValidPreparationInstrument(validPreparationVesselID)
      .then((response: APIResponse<ValidPreparationInstrument>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching Meal', () => {
    let mealID = faker.string.uuid();

    const exampleResponse = new Meal({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meals/${mealID}`).reply(200, exampleResponse);

    client.getMeal(mealID).then((response: APIResponse<Meal>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating Household', () => {
    const exampleResponse = new Household({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households`).reply(200, exampleResponse);

    client.createHousehold().then((response: APIResponse<Household>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = faker.string.uuid();

    const exampleResponse = new ValidIngredientMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(200, exampleResponse);

    client
      .archiveValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)
      .then((response: APIResponse<ValidIngredientMeasurementUnit>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating PasswordResetResponse', () => {
    const exampleResponse = new PasswordResetResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users/password/new`).reply(200, exampleResponse);

    client.updatePassword().then((response: APIResponse<PasswordResetResponse>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching RecipeStepVessel', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStepVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels`).reply(200, exampleResponse);

    client.getRecipeStepVessels(recipeID, recipeStepID).then((response: APIResponse<Array<RecipeStepVessel>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = faker.string.uuid();

    const exampleResponse = new ValidIngredientStateIngredient({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(200, exampleResponse);

    client
      .archiveValidIngredientStateIngredient(validIngredientStateIngredientID)
      .then((response: APIResponse<ValidIngredientStateIngredient>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating ValidIngredientGroup', () => {
    const exampleResponse = new ValidIngredientGroup({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups`).reply(200, exampleResponse);

    client.createValidIngredientGroup().then((response: APIResponse<ValidIngredientGroup>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving MealPlanOptionVote', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();
    let mealPlanOptionID = faker.string.uuid();
    let mealPlanOptionVoteID = faker.string.uuid();

    const exampleResponse = new MealPlanOptionVote({
      id: faker.string.uuid(),
    });
    mock
      .onGet(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(200, exampleResponse);

    client
      .archiveMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
      .then((response: APIResponse<MealPlanOptionVote>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating RecipeStepInstrument', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepInstrumentID = faker.string.uuid();

    const exampleResponse = new RecipeStepInstrument({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(200, exampleResponse);

    client
      .updateRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)
      .then((response: APIResponse<RecipeStepInstrument>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating MealPlan', () => {
    const exampleResponse = new MealPlan({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans`).reply(200, exampleResponse);

    client.createMealPlan().then((response: APIResponse<MealPlan>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching HouseholdInvitation', () => {
    let householdInvitationID = faker.string.uuid();

    const exampleResponse = new HouseholdInvitation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/household_invitations/${householdInvitationID}`).reply(200, exampleResponse);

    client.getHouseholdInvitation(householdInvitationID).then((response: APIResponse<HouseholdInvitation>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving OAuth2Client', () => {
    let oauth2ClientID = faker.string.uuid();

    const exampleResponse = new OAuth2Client({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/oauth2_clients/${oauth2ClientID}`).reply(200, exampleResponse);

    client.archiveOAuth2Client(oauth2ClientID).then((response: APIResponse<OAuth2Client>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidPreparationVessel', () => {
    const exampleResponse = new ValidPreparationVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels`).reply(200, exampleResponse);

    client.getValidPreparationVessels().then((response: APIResponse<Array<ValidPreparationVessel>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidMeasurementUnit', () => {
    let q = faker.string.uuid();
    let validIngredientID = faker.string.uuid();

    const exampleResponse = new ValidMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_units/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    client
      .searchValidMeasurementUnitsByIngredient(q, validIngredientID)
      .then((response: APIResponse<Array<ValidMeasurementUnit>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating CreateMealPlanTasksResponse', () => {
    const exampleResponse = new CreateMealPlanTasksResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/workers/meal_plan_tasks`).reply(200, exampleResponse);

    client.runMealPlanTaskCreatorWorker().then((response: APIResponse<CreateMealPlanTasksResponse>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidMeasurementUnit', () => {
    let q = faker.string.uuid();

    const exampleResponse = new ValidMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/search`).reply(200, exampleResponse);

    client.searchForValidMeasurementUnits(q).then((response: APIResponse<Array<ValidMeasurementUnit>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating ValidIngredientState', () => {
    let validIngredientStateID = faker.string.uuid();

    const exampleResponse = new ValidIngredientState({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(200, exampleResponse);

    client.updateValidIngredientState(validIngredientStateID).then((response: APIResponse<ValidIngredientState>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching Webhook', () => {
    const exampleResponse = new Webhook({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/webhooks`).reply(200, exampleResponse);

    client.getWebhooks().then((response: APIResponse<Array<Webhook>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating ValidMeasurementUnit', () => {
    let validMeasurementUnitID = faker.string.uuid();

    const exampleResponse = new ValidMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(200, exampleResponse);

    client.updateValidMeasurementUnit(validMeasurementUnitID).then((response: APIResponse<ValidMeasurementUnit>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating RecipeStepVessel', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepVesselID = faker.string.uuid();

    const exampleResponse = new RecipeStepVessel({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(200, exampleResponse);

    client
      .updateRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)
      .then((response: APIResponse<RecipeStepVessel>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching MealPlanGroceryListItem', () => {
    let mealPlanID = faker.string.uuid();

    const exampleResponse = new MealPlanGroceryListItem({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items`).reply(200, exampleResponse);

    client
      .getMealPlanGroceryListItemsForMealPlan(mealPlanID)
      .then((response: APIResponse<Array<MealPlanGroceryListItem>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidPreparationVessel', () => {
    let validPreparationID = faker.string.uuid();

    const exampleResponse = new ValidPreparationVessel({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_vessels/by_preparation/${validPreparationID}`)
      .reply(200, exampleResponse);

    client
      .getValidPreparationVesselsByPreparation(validPreparationID)
      .then((response: APIResponse<Array<ValidPreparationVessel>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching Household', () => {
    let householdID = faker.string.uuid();

    const exampleResponse = new Household({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households/${householdID}`).reply(200, exampleResponse);

    client.getHousehold(householdID).then((response: APIResponse<Household>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating RecipeStepProduct', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepProductID = faker.string.uuid();

    const exampleResponse = new RecipeStepProduct({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(200, exampleResponse);

    client
      .updateRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)
      .then((response: APIResponse<RecipeStepProduct>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching RecipeRating', () => {
    let recipeID = faker.string.uuid();

    const exampleResponse = new RecipeRating({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings`).reply(200, exampleResponse);

    client.getRecipeRatings(recipeID).then((response: APIResponse<Array<RecipeRating>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating ValidIngredientStateIngredient', () => {
    const exampleResponse = new ValidIngredientStateIngredient({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients`).reply(200, exampleResponse);

    client.createValidIngredientStateIngredient().then((response: APIResponse<ValidIngredientStateIngredient>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving MealPlan', () => {
    let mealPlanID = faker.string.uuid();

    const exampleResponse = new MealPlan({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(200, exampleResponse);

    client.archiveMealPlan(mealPlanID).then((response: APIResponse<MealPlan>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ServiceSetting', () => {
    let serviceSettingID = faker.string.uuid();

    const exampleResponse = new ServiceSetting({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/settings/${serviceSettingID}`).reply(200, exampleResponse);

    client.getServiceSetting(serviceSettingID).then((response: APIResponse<ServiceSetting>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ValidIngredientPreparation', () => {
    let validIngredientPreparationID = faker.string.uuid();

    const exampleResponse = new ValidIngredientPreparation({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(200, exampleResponse);

    client
      .archiveValidIngredientPreparation(validIngredientPreparationID)
      .then((response: APIResponse<ValidIngredientPreparation>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching MealPlanOptionVote', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();
    let mealPlanOptionID = faker.string.uuid();

    const exampleResponse = new MealPlanOptionVote({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes`)
      .reply(200, exampleResponse);

    client
      .getMealPlanOptionVotes(mealPlanID, mealPlanEventID, mealPlanOptionID)
      .then((response: APIResponse<Array<MealPlanOptionVote>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating Household', () => {
    let householdID = faker.string.uuid();

    const exampleResponse = new Household({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households/${householdID}`).reply(200, exampleResponse);

    client.updateHousehold(householdID).then((response: APIResponse<Household>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating PasswordResetToken', () => {
    const exampleResponse = new PasswordResetToken({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/users/password/reset`).reply(200, exampleResponse);

    client.requestPasswordResetToken().then((response: APIResponse<PasswordResetToken>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ServiceSettingConfiguration', () => {
    let serviceSettingConfigurationID = faker.string.uuid();

    const exampleResponse = new ServiceSettingConfiguration({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/settings/configurations/${serviceSettingConfigurationID}`)
      .reply(200, exampleResponse);

    client
      .archiveServiceSettingConfiguration(serviceSettingConfigurationID)
      .then((response: APIResponse<ServiceSettingConfiguration>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching MealPlan', () => {
    let mealPlanID = faker.string.uuid();

    const exampleResponse = new MealPlan({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}`).reply(200, exampleResponse);

    client.getMealPlan(mealPlanID).then((response: APIResponse<MealPlan>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating UserPermissionsResponse', () => {
    let householdID = faker.string.uuid();
    let userID = faker.string.uuid();

    const exampleResponse = new UserPermissionsResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households/${householdID}/members/${userID}/permissions`).reply(200, exampleResponse);

    client
      .updateHouseholdMemberPermissions(householdID, userID)
      .then((response: APIResponse<UserPermissionsResponse>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidIngredientGroup', () => {
    let validIngredientGroupID = faker.string.uuid();

    const exampleResponse = new ValidIngredientGroup({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(200, exampleResponse);

    client.getValidIngredientGroup(validIngredientGroupID).then((response: APIResponse<ValidIngredientGroup>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidPreparation', () => {
    let validPreparationID = faker.string.uuid();

    const exampleResponse = new ValidPreparation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(200, exampleResponse);

    client.getValidPreparation(validPreparationID).then((response: APIResponse<ValidPreparation>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating RecipeStepProduct', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStepProduct({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`).reply(200, exampleResponse);

    client.createRecipeStepProduct(recipeID, recipeStepID).then((response: APIResponse<RecipeStepProduct>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredientPreparation', () => {
    let validIngredientPreparationID = faker.string.uuid();

    const exampleResponse = new ValidIngredientPreparation({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/${validIngredientPreparationID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientPreparation(validIngredientPreparationID)
      .then((response: APIResponse<ValidIngredientPreparation>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating Household', () => {
    let householdID = faker.string.uuid();

    const exampleResponse = new Household({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households/${householdID}/default`).reply(200, exampleResponse);

    client.setDefaultHousehold(householdID).then((response: APIResponse<Household>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating User', () => {
    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users/username`).reply(200, exampleResponse);

    client.updateUserUsername().then((response: APIResponse<User>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for redeeming a password reset token', () => {
    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/users/password/reset/redeem`).reply(200, exampleResponse);

    client.redeemPasswordResetToken().then((response: APIResponse<User>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching RecipeStepIngredient', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepIngredientID = faker.string.uuid();

    const exampleResponse = new RecipeStepIngredient({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(200, exampleResponse);

    client
      .getRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)
      .then((response: APIResponse<RecipeStepIngredient>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidPreparation', () => {
    const exampleResponse = new ValidPreparation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparations`).reply(200, exampleResponse);

    client.getValidPreparations().then((response: APIResponse<Array<ValidPreparation>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredientMeasurementUnit', () => {
    const exampleResponse = new ValidIngredientMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units`).reply(200, exampleResponse);

    client.getValidIngredientMeasurementUnits().then((response: APIResponse<Array<ValidIngredientMeasurementUnit>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating User', () => {
    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users/details`).reply(200, exampleResponse);

    client.updateUserDetails().then((response: APIResponse<User>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating ValidPreparation', () => {
    let validPreparationID = faker.string.uuid();

    const exampleResponse = new ValidPreparation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(200, exampleResponse);

    client.updateValidPreparation(validPreparationID).then((response: APIResponse<ValidPreparation>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating ValidMeasurementUnitConversion', () => {
    const exampleResponse = new ValidMeasurementUnitConversion({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_measurement_conversions`).reply(200, exampleResponse);

    client.createValidMeasurementUnitConversion().then((response: APIResponse<ValidMeasurementUnitConversion>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating Recipe', () => {
    const exampleResponse = new Recipe({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes`).reply(200, exampleResponse);

    client.createRecipe().then((response: APIResponse<Recipe>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating ValidPreparationVessel', () => {
    let validPreparationVesselID = faker.string.uuid();

    const exampleResponse = new ValidPreparationVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(200, exampleResponse);

    client
      .updateValidPreparationVessel(validPreparationVesselID)
      .then((response: APIResponse<ValidPreparationVessel>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating UserIngredientPreference', () => {
    let userIngredientPreferenceID = faker.string.uuid();

    const exampleResponse = new UserIngredientPreference({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`)
      .reply(200, exampleResponse);

    client
      .updateUserIngredientPreference(userIngredientPreferenceID)
      .then((response: APIResponse<UserIngredientPreference>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating User', () => {
    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users/email_address`).reply(200, exampleResponse);

    client.updateUserEmailAddress().then((response: APIResponse<User>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching User', () => {
    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users`).reply(200, exampleResponse);

    client.getUsers().then((response: APIResponse<Array<User>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating User', () => {
    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/users/avatar/upload`).reply(200, exampleResponse);

    client.uploadUserAvatar().then((response: APIResponse<User>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredient', () => {
    let validIngredientID = faker.string.uuid();

    const exampleResponse = new ValidIngredient({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(200, exampleResponse);

    client.getValidIngredient(validIngredientID).then((response: APIResponse<ValidIngredient>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ValidPreparation', () => {
    let validPreparationID = faker.string.uuid();

    const exampleResponse = new ValidPreparation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparations/${validPreparationID}`).reply(200, exampleResponse);

    client.archiveValidPreparation(validPreparationID).then((response: APIResponse<ValidPreparation>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching HouseholdInstrumentOwnership', () => {
    let householdInstrumentOwnershipID = faker.string.uuid();

    const exampleResponse = new HouseholdInstrumentOwnership({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(200, exampleResponse);

    client
      .getHouseholdInstrumentOwnership(householdInstrumentOwnershipID)
      .then((response: APIResponse<HouseholdInstrumentOwnership>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating RecipeStep', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStep({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(200, exampleResponse);

    client.updateRecipeStep(recipeID, recipeStepID).then((response: APIResponse<RecipeStep>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving HouseholdUserMembership', () => {
    let householdID = faker.string.uuid();
    let userID = faker.string.uuid();

    const exampleResponse = new HouseholdUserMembership({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households/${householdID}/members/${userID}`).reply(200, exampleResponse);

    client.archiveUserMembership(householdID, userID).then((response: APIResponse<HouseholdUserMembership>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating RecipeRating', () => {
    let recipeID = faker.string.uuid();
    let recipeRatingID = faker.string.uuid();

    const exampleResponse = new RecipeRating({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(200, exampleResponse);

    client.updateRecipeRating(recipeID, recipeRatingID).then((response: APIResponse<RecipeRating>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving RecipeRating', () => {
    let recipeID = faker.string.uuid();
    let recipeRatingID = faker.string.uuid();

    const exampleResponse = new RecipeRating({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(200, exampleResponse);

    client.archiveRecipeRating(recipeID, recipeRatingID).then((response: APIResponse<RecipeRating>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitID = faker.string.uuid();

    const exampleResponse = new ValidMeasurementUnitConversion({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/to_unit/${validMeasurementUnitID}`)
      .reply(200, exampleResponse);

    client
      .getValidMeasurementUnitConversionsToUnit(validMeasurementUnitID)
      .then((response: APIResponse<Array<ValidMeasurementUnitConversion>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating User', () => {
    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/users/email_address/verify`).reply(200, exampleResponse);

    client.verifyEmailAddress().then((response: APIResponse<User>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating MealPlanOption', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();

    const exampleResponse = new MealPlanOption({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options`)
      .reply(200, exampleResponse);

    client.createMealPlanOption(mealPlanID, mealPlanEventID).then((response: APIResponse<MealPlanOption>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching RecipeStepProduct', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStepProduct({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products`).reply(200, exampleResponse);

    client.getRecipeStepProducts(recipeID, recipeStepID).then((response: APIResponse<Array<RecipeStepProduct>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving RecipeStepVessel', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepVesselID = faker.string.uuid();

    const exampleResponse = new RecipeStepVessel({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/vessels/${recipeStepVesselID}`)
      .reply(200, exampleResponse);

    client
      .archiveRecipeStepVessel(recipeID, recipeStepID, recipeStepVesselID)
      .then((response: APIResponse<RecipeStepVessel>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidIngredientStateIngredient', () => {
    let validIngredientStateIngredientID = faker.string.uuid();

    const exampleResponse = new ValidIngredientStateIngredient({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients/${validIngredientStateIngredientID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientStateIngredient(validIngredientStateIngredientID)
      .then((response: APIResponse<ValidIngredientStateIngredient>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating MealPlanGroceryListItem', () => {
    let mealPlanID = faker.string.uuid();

    const exampleResponse = new MealPlanGroceryListItem({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items`).reply(200, exampleResponse);

    client.createMealPlanGroceryListItem(mealPlanID).then((response: APIResponse<MealPlanGroceryListItem>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidInstrument', () => {
    let validInstrumentID = faker.string.uuid();

    const exampleResponse = new ValidInstrument({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(200, exampleResponse);

    client.getValidInstrument(validInstrumentID).then((response: APIResponse<ValidInstrument>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching OAuth2Client', () => {
    const exampleResponse = new OAuth2Client({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/oauth2_clients`).reply(200, exampleResponse);

    client.getOAuth2Clients().then((response: APIResponse<Array<OAuth2Client>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredientState', () => {
    let q = faker.string.uuid();

    const exampleResponse = new ValidIngredientState({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/search`).reply(200, exampleResponse);

    client.searchForValidIngredientStates(q).then((response: APIResponse<Array<ValidIngredientState>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating JWTResponse', () => {
    const exampleResponse = new JWTResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/users/login/jwt`).reply(200, exampleResponse);

    client.loginForJWT().then((response: APIResponse<JWTResponse>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidPreparationVessel', () => {
    let ValidVesselID = faker.string.uuid();

    const exampleResponse = new ValidPreparationVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/by_vessel/${ValidVesselID}`).reply(200, exampleResponse);

    client
      .getValidPreparationVesselsByVessel(ValidVesselID)
      .then((response: APIResponse<Array<ValidPreparationVessel>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching Recipe', () => {
    let q = faker.string.uuid();

    const exampleResponse = new Recipe({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/search`).reply(200, exampleResponse);

    client.searchForRecipes(q).then((response: APIResponse<Array<Recipe>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching RecipeStepCompletionCondition', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStepCompletionCondition({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`)
      .reply(200, exampleResponse);

    client
      .getRecipeStepCompletionConditions(recipeID, recipeStepID)
      .then((response: APIResponse<Array<RecipeStepCompletionCondition>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating RecipeStepInstrument', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStepInstrument({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`).reply(200, exampleResponse);

    client.createRecipeStepInstrument(recipeID, recipeStepID).then((response: APIResponse<RecipeStepInstrument>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating MealPlanOptionVote', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();

    const exampleResponse = new MealPlanOptionVote({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/vote`).reply(200, exampleResponse);

    client
      .createMealPlanOptionVote(mealPlanID, mealPlanEventID)
      .then((response: APIResponse<Array<MealPlanOptionVote>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching RecipePrepTask', () => {
    let recipeID = faker.string.uuid();
    let recipePrepTaskID = faker.string.uuid();

    const exampleResponse = new RecipePrepTask({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(200, exampleResponse);

    client.getRecipePrepTask(recipeID, recipePrepTaskID).then((response: APIResponse<RecipePrepTask>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching AuditLogEntry', () => {
    const exampleResponse = new AuditLogEntry({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/audit_log_entries/for_household`).reply(200, exampleResponse);

    client.getAuditLogEntriesForHousehold().then((response: APIResponse<Array<AuditLogEntry>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving MealPlanGroceryListItem', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanGroceryListItemID = faker.string.uuid();

    const exampleResponse = new MealPlanGroceryListItem({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/grocery_list_items/${mealPlanGroceryListItemID}`)
      .reply(200, exampleResponse);

    client
      .archiveMealPlanGroceryListItem(mealPlanID, mealPlanGroceryListItemID)
      .then((response: APIResponse<MealPlanGroceryListItem>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating ValidPreparationVessel', () => {
    const exampleResponse = new ValidPreparationVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels`).reply(200, exampleResponse);

    client.createValidPreparationVessel().then((response: APIResponse<ValidPreparationVessel>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating JWTResponse', () => {
    const exampleResponse = new JWTResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/users/login/jwt/admin`).reply(200, exampleResponse);

    client.adminLoginForJWT().then((response: APIResponse<JWTResponse>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching RecipeStepCompletionCondition', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepCompletionConditionID = faker.string.uuid();

    const exampleResponse = new RecipeStepCompletionCondition({
      id: faker.string.uuid(),
    });
    mock
      .onGet(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(200, exampleResponse);

    client
      .getRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID)
      .then((response: APIResponse<RecipeStepCompletionCondition>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating HouseholdInvitation', () => {
    let householdInvitationID = faker.string.uuid();

    const exampleResponse = new HouseholdInvitation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/cancel`).reply(200, exampleResponse);

    client.cancelHouseholdInvitation(householdInvitationID).then((response: APIResponse<HouseholdInvitation>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating InitializeMealPlanGroceryListResponse', () => {
    const exampleResponse = new InitializeMealPlanGroceryListResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/workers/meal_plan_grocery_list_init`).reply(200, exampleResponse);

    client
      .runMealPlanGroceryListInitializerWorker()
      .then((response: APIResponse<InitializeMealPlanGroceryListResponse>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching UserNotification', () => {
    const exampleResponse = new UserNotification({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/user_notifications`).reply(200, exampleResponse);

    client.getUserNotifications().then((response: APIResponse<Array<UserNotification>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating User', () => {
    const exampleResponse = new User({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/users/totp_secret/verify`).reply(200, exampleResponse);

    client.verifyTOTPSecret().then((response: APIResponse<User>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating ServiceSettingConfiguration', () => {
    let serviceSettingConfigurationID = faker.string.uuid();

    const exampleResponse = new ServiceSettingConfiguration({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/settings/configurations/${serviceSettingConfigurationID}`)
      .reply(200, exampleResponse);

    client
      .updateServiceSettingConfiguration(serviceSettingConfigurationID)
      .then((response: APIResponse<ServiceSettingConfiguration>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating ValidIngredientPreparation', () => {
    const exampleResponse = new ValidIngredientPreparation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_preparations`).reply(200, exampleResponse);

    client.createValidIngredientPreparation().then((response: APIResponse<ValidIngredientPreparation>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching RecipeStep', () => {
    let recipeID = faker.string.uuid();

    const exampleResponse = new RecipeStep({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps`).reply(200, exampleResponse);

    client.getRecipeSteps(recipeID).then((response: APIResponse<Array<RecipeStep>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating FinalizeMealPlansResponse', () => {
    let mealPlanID = faker.string.uuid();

    const exampleResponse = new FinalizeMealPlansResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/finalize`).reply(200, exampleResponse);

    client.finalizeMealPlan(mealPlanID).then((response: APIResponse<FinalizeMealPlansResponse>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredient', () => {
    let q = faker.string.uuid();
    let validPreparationID = faker.string.uuid();

    const exampleResponse = new ValidIngredient({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/by_preparation/${validPreparationID}`).reply(200, exampleResponse);

    client
      .searchValidIngredientsByPreparation(q, validPreparationID)
      .then((response: APIResponse<Array<ValidIngredient>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching string', () => {
    let recipeID = faker.string.uuid();

    const exampleResponse = new string({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/mermaid`).reply(200, exampleResponse);

    client.getMermaidDiagramForRecipe(recipeID).then((response: APIResponse<string>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidVessel', () => {
    const exampleResponse = new ValidVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_vessels/random`).reply(200, exampleResponse);

    client.getRandomValidVessel().then((response: APIResponse<ValidVessel>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving RecipeStepIngredient', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepIngredientID = faker.string.uuid();

    const exampleResponse = new RecipeStepIngredient({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients/${recipeStepIngredientID}`)
      .reply(200, exampleResponse);

    client
      .archiveRecipeStepIngredient(recipeID, recipeStepID, recipeStepIngredientID)
      .then((response: APIResponse<RecipeStepIngredient>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching RecipeStep', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStep({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(200, exampleResponse);

    client.getRecipeStep(recipeID, recipeStepID).then((response: APIResponse<RecipeStep>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating ServiceSettingConfiguration', () => {
    const exampleResponse = new ServiceSettingConfiguration({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/settings/configurations`).reply(200, exampleResponse);

    client.createServiceSettingConfiguration().then((response: APIResponse<ServiceSettingConfiguration>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ValidMeasurementUnit', () => {
    let validMeasurementUnitID = faker.string.uuid();

    const exampleResponse = new ValidMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(200, exampleResponse);

    client.archiveValidMeasurementUnit(validMeasurementUnitID).then((response: APIResponse<ValidMeasurementUnit>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidPreparationInstrument', () => {
    let validPreparationID = faker.string.uuid();

    const exampleResponse = new ValidPreparationInstrument({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/by_preparation/${validPreparationID}`)
      .reply(200, exampleResponse);

    client
      .getValidPreparationInstrumentsByPreparation(validPreparationID)
      .then((response: APIResponse<Array<ValidPreparationInstrument>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for archiving ValidIngredientState', () => {
    let validIngredientStateID = faker.string.uuid();

    const exampleResponse = new ValidIngredientState({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(200, exampleResponse);

    client.archiveValidIngredientState(validIngredientStateID).then((response: APIResponse<ValidIngredientState>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidPreparation', () => {
    const exampleResponse = new ValidPreparation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparations/random`).reply(200, exampleResponse);

    client.getRandomValidPreparation().then((response: APIResponse<ValidPreparation>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidPreparation', () => {
    let q = faker.string.uuid();

    const exampleResponse = new ValidPreparation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparations/search`).reply(200, exampleResponse);

    client.searchForValidPreparations(q).then((response: APIResponse<Array<ValidPreparation>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredientState', () => {
    const exampleResponse = new ValidIngredientState({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states`).reply(200, exampleResponse);

    client.getValidIngredientStates().then((response: APIResponse<Array<ValidIngredientState>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ValidIngredientGroup', () => {
    let validIngredientGroupID = faker.string.uuid();

    const exampleResponse = new ValidIngredientGroup({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/${validIngredientGroupID}`).reply(200, exampleResponse);

    client.archiveValidIngredientGroup(validIngredientGroupID).then((response: APIResponse<ValidIngredientGroup>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating ValidIngredient', () => {
    let validIngredientID = faker.string.uuid();

    const exampleResponse = new ValidIngredient({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredients/${validIngredientID}`).reply(200, exampleResponse);

    client.updateValidIngredient(validIngredientID).then((response: APIResponse<ValidIngredient>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating ValidVessel', () => {
    let validVesselID = faker.string.uuid();

    const exampleResponse = new ValidVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(200, exampleResponse);

    client.updateValidVessel(validVesselID).then((response: APIResponse<ValidVessel>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = faker.string.uuid();

    const exampleResponse = new ValidMeasurementUnitConversion({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(200, exampleResponse);

    client
      .archiveValidMeasurementUnitConversion(validMeasurementUnitConversionID)
      .then((response: APIResponse<ValidMeasurementUnitConversion>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching RecipePrepTask', () => {
    let recipeID = faker.string.uuid();

    const exampleResponse = new RecipePrepTask({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks`).reply(200, exampleResponse);

    client.getRecipePrepTasks(recipeID).then((response: APIResponse<Array<RecipePrepTask>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching UserStatusResponse', () => {
    const exampleResponse = new UserStatusResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/auth/status`).reply(200, exampleResponse);

    client.getAuthStatus().then((response: APIResponse<UserStatusResponse>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching RecipeStepProduct', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepProductID = faker.string.uuid();

    const exampleResponse = new RecipeStepProduct({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(200, exampleResponse);

    client
      .getRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)
      .then((response: APIResponse<RecipeStepProduct>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating RecipeStepIngredient', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStepIngredient({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/ingredients`).reply(200, exampleResponse);

    client.createRecipeStepIngredient(recipeID, recipeStepID).then((response: APIResponse<RecipeStepIngredient>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating HouseholdInstrumentOwnership', () => {
    const exampleResponse = new HouseholdInstrumentOwnership({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/households/instruments`).reply(200, exampleResponse);

    client.createHouseholdInstrumentOwnership().then((response: APIResponse<HouseholdInstrumentOwnership>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating ValidIngredient', () => {
    const exampleResponse = new ValidIngredient({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredients`).reply(200, exampleResponse);

    client.createValidIngredient().then((response: APIResponse<ValidIngredient>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidPreparationInstrument', () => {
    let validInstrumentID = faker.string.uuid();

    const exampleResponse = new ValidPreparationInstrument({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/by_instrument/${validInstrumentID}`)
      .reply(200, exampleResponse);

    client
      .getValidPreparationInstrumentsByInstrument(validInstrumentID)
      .then((response: APIResponse<Array<ValidPreparationInstrument>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating RecipePrepTask', () => {
    let recipeID = faker.string.uuid();
    let recipePrepTaskID = faker.string.uuid();

    const exampleResponse = new RecipePrepTask({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks/${recipePrepTaskID}`).reply(200, exampleResponse);

    client.updateRecipePrepTask(recipeID, recipePrepTaskID).then((response: APIResponse<RecipePrepTask>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching MealPlanEvent', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();

    const exampleResponse = new MealPlanEvent({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}`).reply(200, exampleResponse);

    client.getMealPlanEvent(mealPlanID, mealPlanEventID).then((response: APIResponse<MealPlanEvent>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ValidVessel', () => {
    let validVesselID = faker.string.uuid();

    const exampleResponse = new ValidVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_vessels/${validVesselID}`).reply(200, exampleResponse);

    client.archiveValidVessel(validVesselID).then((response: APIResponse<ValidVessel>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredientMeasurementUnit', () => {
    let validIngredientID = faker.string.uuid();

    const exampleResponse = new ValidIngredientMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientMeasurementUnitsByIngredient(validIngredientID)
      .then((response: APIResponse<Array<ValidIngredientMeasurementUnit>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for archiving RecipeStepInstrument', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepInstrumentID = faker.string.uuid();

    const exampleResponse = new RecipeStepInstrument({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments/${recipeStepInstrumentID}`)
      .reply(200, exampleResponse);

    client
      .archiveRecipeStepInstrument(recipeID, recipeStepID, recipeStepInstrumentID)
      .then((response: APIResponse<RecipeStepInstrument>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating HouseholdInstrumentOwnership', () => {
    let householdInstrumentOwnershipID = faker.string.uuid();

    const exampleResponse = new HouseholdInstrumentOwnership({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/households/instruments/${householdInstrumentOwnershipID}`)
      .reply(200, exampleResponse);

    client
      .updateHouseholdInstrumentOwnership(householdInstrumentOwnershipID)
      .then((response: APIResponse<HouseholdInstrumentOwnership>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidIngredientGroup', () => {
    let q = faker.string.uuid();

    const exampleResponse = new ValidIngredientGroup({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups/search`).reply(200, exampleResponse);

    client.searchForValidIngredientGroups(q).then((response: APIResponse<Array<ValidIngredientGroup>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving UserIngredientPreference', () => {
    let userIngredientPreferenceID = faker.string.uuid();

    const exampleResponse = new UserIngredientPreference({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/user_ingredient_preferences/${userIngredientPreferenceID}`)
      .reply(200, exampleResponse);

    client
      .archiveUserIngredientPreference(userIngredientPreferenceID)
      .then((response: APIResponse<UserIngredientPreference>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating HouseholdInvitation', () => {
    let householdInvitationID = faker.string.uuid();

    const exampleResponse = new HouseholdInvitation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/household_invitations/${householdInvitationID}/accept`).reply(200, exampleResponse);

    client.acceptHouseholdInvitation(householdInvitationID).then((response: APIResponse<HouseholdInvitation>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredientPreparation', () => {
    let validIngredientID = faker.string.uuid();

    const exampleResponse = new ValidIngredientPreparation({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_preparations/by_ingredient/${validIngredientID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientPreparationsByIngredient(validIngredientID)
      .then((response: APIResponse<Array<ValidIngredientPreparation>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching HouseholdInvitation', () => {
    const exampleResponse = new HouseholdInvitation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/household_invitations/sent`).reply(200, exampleResponse);

    client.getSentHouseholdInvitations().then((response: APIResponse<Array<HouseholdInvitation>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredientPreparation', () => {
    const exampleResponse = new ValidIngredientPreparation({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_preparations`).reply(200, exampleResponse);

    client.getValidIngredientPreparations().then((response: APIResponse<Array<ValidIngredientPreparation>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating RecipeStep', () => {
    let recipeID = faker.string.uuid();

    const exampleResponse = new RecipeStep({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps`).reply(200, exampleResponse);

    client.createRecipeStep(recipeID).then((response: APIResponse<RecipeStep>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ValidPreparationVessel', () => {
    let validPreparationVesselID = faker.string.uuid();

    const exampleResponse = new ValidPreparationVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparation_vessels/${validPreparationVesselID}`).reply(200, exampleResponse);

    client
      .archiveValidPreparationVessel(validPreparationVesselID)
      .then((response: APIResponse<ValidPreparationVessel>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating Meal', () => {
    const exampleResponse = new Meal({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meals`).reply(200, exampleResponse);

    client.createMeal().then((response: APIResponse<Meal>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredientStateIngredient', () => {
    const exampleResponse = new ValidIngredientStateIngredient({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_state_ingredients`).reply(200, exampleResponse);

    client.getValidIngredientStateIngredients().then((response: APIResponse<Array<ValidIngredientStateIngredient>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching Meal', () => {
    const exampleResponse = new Meal({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meals`).reply(200, exampleResponse);

    client.getMeals().then((response: APIResponse<Array<Meal>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ValidInstrument', () => {
    let validInstrumentID = faker.string.uuid();

    const exampleResponse = new ValidInstrument({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_instruments/${validInstrumentID}`).reply(200, exampleResponse);

    client.archiveValidInstrument(validInstrumentID).then((response: APIResponse<ValidInstrument>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating RecipePrepTask', () => {
    let recipeID = faker.string.uuid();

    const exampleResponse = new RecipePrepTask({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/prep_tasks`).reply(200, exampleResponse);

    client.createRecipePrepTask(recipeID).then((response: APIResponse<RecipePrepTask>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Creates a new user', () => {
    const exampleResponse = new UserCreationResponse({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/users`).reply(200, exampleResponse);

    client.createUser().then((response: APIResponse<UserCreationResponse>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidVessel', () => {
    const exampleResponse = new ValidVessel({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_vessels`).reply(200, exampleResponse);

    client.getValidVessels().then((response: APIResponse<Array<ValidVessel>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving Webhook', () => {
    let webhookID = faker.string.uuid();

    const exampleResponse = new Webhook({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/webhooks/${webhookID}`).reply(200, exampleResponse);

    client.archiveWebhook(webhookID).then((response: APIResponse<Webhook>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving RecipeStepProduct', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepProductID = faker.string.uuid();

    const exampleResponse = new RecipeStepProduct({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/products/${recipeStepProductID}`)
      .reply(200, exampleResponse);

    client
      .archiveRecipeStepProduct(recipeID, recipeStepID, recipeStepProductID)
      .then((response: APIResponse<RecipeStepProduct>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ServiceSettingConfiguration', () => {
    const exampleResponse = new ServiceSettingConfiguration({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/settings/configurations/user`).reply(200, exampleResponse);

    client
      .getServiceSettingConfigurationsForUser()
      .then((response: APIResponse<Array<ServiceSettingConfiguration>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for archiving MealPlanOption', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();
    let mealPlanOptionID = faker.string.uuid();

    const exampleResponse = new MealPlanOption({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}`)
      .reply(200, exampleResponse);

    client
      .archiveMealPlanOption(mealPlanID, mealPlanEventID, mealPlanOptionID)
      .then((response: APIResponse<MealPlanOption>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidIngredientState', () => {
    let validIngredientStateID = faker.string.uuid();

    const exampleResponse = new ValidIngredientState({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_states/${validIngredientStateID}`).reply(200, exampleResponse);

    client.getValidIngredientState(validIngredientStateID).then((response: APIResponse<ValidIngredientState>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating UserNotification', () => {
    const exampleResponse = new UserNotification({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/user_notifications`).reply(200, exampleResponse);

    client.createUserNotification().then((response: APIResponse<UserNotification>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching MealPlanTask', () => {
    let mealPlanID = faker.string.uuid();

    const exampleResponse = new MealPlanTask({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/tasks`).reply(200, exampleResponse);

    client.getMealPlanTasks(mealPlanID).then((response: APIResponse<Array<MealPlanTask>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching Webhook', () => {
    let webhookID = faker.string.uuid();

    const exampleResponse = new Webhook({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/webhooks/${webhookID}`).reply(200, exampleResponse);

    client.getWebhook(webhookID).then((response: APIResponse<Webhook>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving RecipeStepCompletionCondition', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();
    let recipeStepCompletionConditionID = faker.string.uuid();

    const exampleResponse = new RecipeStepCompletionCondition({
      id: faker.string.uuid(),
    });
    mock
      .onGet(
        `${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions/${recipeStepCompletionConditionID}`,
      )
      .reply(200, exampleResponse);

    client
      .archiveRecipeStepCompletionCondition(recipeID, recipeStepID, recipeStepCompletionConditionID)
      .then((response: APIResponse<RecipeStepCompletionCondition>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching ValidIngredientGroup', () => {
    const exampleResponse = new ValidIngredientGroup({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_ingredient_groups`).reply(200, exampleResponse);

    client.getValidIngredientGroups().then((response: APIResponse<Array<ValidIngredientGroup>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching MealPlan', () => {
    const exampleResponse = new MealPlan({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans`).reply(200, exampleResponse);

    client.getMealPlans().then((response: APIResponse<Array<MealPlan>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for creating Recipe', () => {
    let recipeID = faker.string.uuid();

    const exampleResponse = new Recipe({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/clone`).reply(200, exampleResponse);

    client.cloneRecipe(recipeID).then((response: APIResponse<Recipe>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching MealPlanEvent', () => {
    let mealPlanID = faker.string.uuid();

    const exampleResponse = new MealPlanEvent({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/meal_plans/${mealPlanID}/events`).reply(200, exampleResponse);

    client.getMealPlanEvents(mealPlanID).then((response: APIResponse<Array<MealPlanEvent>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for archiving ValidPreparationInstrument', () => {
    let validPreparationVesselID = faker.string.uuid();

    const exampleResponse = new ValidPreparationInstrument({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_preparation_instruments/${validPreparationVesselID}`)
      .reply(200, exampleResponse);

    client
      .archiveValidPreparationInstrument(validPreparationVesselID)
      .then((response: APIResponse<ValidPreparationInstrument>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating ValidMeasurementUnitConversion', () => {
    let validMeasurementUnitConversionID = faker.string.uuid();

    const exampleResponse = new ValidMeasurementUnitConversion({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_measurement_conversions/${validMeasurementUnitConversionID}`)
      .reply(200, exampleResponse);

    client
      .updateValidMeasurementUnitConversion(validMeasurementUnitConversionID)
      .then((response: APIResponse<ValidMeasurementUnitConversion>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating ValidPreparationInstrument', () => {
    const exampleResponse = new ValidPreparationInstrument({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_preparation_instruments`).reply(200, exampleResponse);

    client.createValidPreparationInstrument().then((response: APIResponse<ValidPreparationInstrument>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidMeasurementUnit', () => {
    let validMeasurementUnitID = faker.string.uuid();

    const exampleResponse = new ValidMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/valid_measurement_units/${validMeasurementUnitID}`).reply(200, exampleResponse);

    client.getValidMeasurementUnit(validMeasurementUnitID).then((response: APIResponse<ValidMeasurementUnit>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for updating MealPlanOptionVote', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();
    let mealPlanOptionID = faker.string.uuid();
    let mealPlanOptionVoteID = faker.string.uuid();

    const exampleResponse = new MealPlanOptionVote({
      id: faker.string.uuid(),
    });
    mock
      .onGet(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(200, exampleResponse);

    client
      .updateMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
      .then((response: APIResponse<MealPlanOptionVote>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for updating ValidIngredientMeasurementUnit', () => {
    let validIngredientMeasurementUnitID = faker.string.uuid();

    const exampleResponse = new ValidIngredientMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/${validIngredientMeasurementUnitID}`)
      .reply(200, exampleResponse);

    client
      .updateValidIngredientMeasurementUnit(validIngredientMeasurementUnitID)
      .then((response: APIResponse<ValidIngredientMeasurementUnit>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for archiving ServiceSetting', () => {
    let serviceSettingID = faker.string.uuid();

    const exampleResponse = new ServiceSetting({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/settings/${serviceSettingID}`).reply(200, exampleResponse);

    client.archiveServiceSetting(serviceSettingID).then((response: APIResponse<ServiceSetting>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching RecipeRating', () => {
    let recipeID = faker.string.uuid();
    let recipeRatingID = faker.string.uuid();

    const exampleResponse = new RecipeRating({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/ratings/${recipeRatingID}`).reply(200, exampleResponse);

    client.getRecipeRating(recipeID, recipeRatingID).then((response: APIResponse<RecipeRating>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching ValidIngredientMeasurementUnit', () => {
    let validMeasurementUnitID = faker.string.uuid();

    const exampleResponse = new ValidIngredientMeasurementUnit({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/valid_ingredient_measurement_units/by_measurement_unit/${validMeasurementUnitID}`)
      .reply(200, exampleResponse);

    client
      .getValidIngredientMeasurementUnitsByMeasurementUnit(validMeasurementUnitID)
      .then((response: APIResponse<Array<ValidIngredientMeasurementUnit>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for creating RecipeStepCompletionCondition', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStepCompletionCondition({
      id: faker.string.uuid(),
    });
    mock
      .onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/completion_conditions`)
      .reply(200, exampleResponse);

    client
      .createRecipeStepCompletionCondition(recipeID, recipeStepID)
      .then((response: APIResponse<RecipeStepCompletionCondition>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for archiving RecipeStep', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStep({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}`).reply(200, exampleResponse);

    client.archiveRecipeStep(recipeID, recipeStepID).then((response: APIResponse<RecipeStep>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching UserIngredientPreference', () => {
    const exampleResponse = new UserIngredientPreference({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/user_ingredient_preferences`).reply(200, exampleResponse);

    client.getUserIngredientPreferences().then((response: APIResponse<Array<UserIngredientPreference>>) => {
      expect(response).toEqual(exampleResponse);
    });
  });

  it('should Operation for fetching RecipeStepInstrument', () => {
    let recipeID = faker.string.uuid();
    let recipeStepID = faker.string.uuid();

    const exampleResponse = new RecipeStepInstrument({
      id: faker.string.uuid(),
    });
    mock.onGet(`${baseURL}/api/v1/recipes/${recipeID}/steps/${recipeStepID}/instruments`).reply(200, exampleResponse);

    client
      .getRecipeStepInstruments(recipeID, recipeStepID)
      .then((response: APIResponse<Array<RecipeStepInstrument>>) => {
        expect(response).toEqual(exampleResponse);
      });
  });

  it('should Operation for fetching MealPlanOptionVote', () => {
    let mealPlanID = faker.string.uuid();
    let mealPlanEventID = faker.string.uuid();
    let mealPlanOptionID = faker.string.uuid();
    let mealPlanOptionVoteID = faker.string.uuid();

    const exampleResponse = new MealPlanOptionVote({
      id: faker.string.uuid(),
    });
    mock
      .onGet(
        `${baseURL}/api/v1/meal_plans/${mealPlanID}/events/${mealPlanEventID}/options/${mealPlanOptionID}/votes/${mealPlanOptionVoteID}`,
      )
      .reply(200, exampleResponse);

    client
      .getMealPlanOptionVote(mealPlanID, mealPlanEventID, mealPlanOptionID, mealPlanOptionVoteID)
      .then((response: APIResponse<MealPlanOptionVote>) => {
        expect(response).toEqual(exampleResponse);
      });
  });
});
