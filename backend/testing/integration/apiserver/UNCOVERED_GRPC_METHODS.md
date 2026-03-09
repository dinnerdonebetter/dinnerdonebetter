# Uncovered gRPC Methods

The following gRPC methods are not tested in the integration tests.

## AuthService

- EvaluateInt64FeatureFlag
- ExchangeToken
- GetSelf
- RequestEmailVerificationEmail
- RequestUsernameReminder
- VerifyEmailAddress

## IdentityService

- GetAccountInvitation
- GetAccountsForUser
- GetUsers
- GetUsersForAccount
- UpdateUserDetails
- UpdateUserEmailAddress
- UpdateUserUsername
- UploadUserAvatar (streaming)

## AuditService

- GetAuditLogEntryByID

## MealPlanningService

- ArchiveMealPlanGroceryListItem
- ArchiveMealPlanRecipeOptionSelection
- ArchiveValidIngredientMeasurementUnit
- ArchiveValidIngredientPreparation
- ArchiveValidIngredientStateIngredient
- ArchiveValidMeasurementUnitConversion
- ArchiveValidPreparationInstrument
- ArchiveValidPreparationVessel
- CreateMealPlanRecipeOptionSelection
- EstimateRecipePrepTasks
- FinalizeMealPlan
- GetMealPlanGroceryListItem
- GetMealPlanRecipeOptionSelection
- GetMealPlanRecipeOptionSelectionsForMealPlanOption
- GetMeasurementUnitConversionMismatches
- GetRecipes
- RunMealPlanTaskCreatorWorker
- SearchForMealEligibleRecipes
- SearchForRecipesWithInstrumentOwnership
- SearchValidIngredientsByPreparation
- SearchValidMeasurementUnitsByIngredient
- UpdateAccountInstrumentOwnership
- UpdateMealPlan
- UpdateMealPlanGroceryListItem
- UpdateMealPlanRecipeOptionSelection
- UpdateMealPlanTaskStatus
- UpdateUserIngredientPreference
- UpdateValidIngredientMeasurementUnit
- UpdateValidIngredientPreparation
- UpdateValidIngredientStateIngredient
- UpdateValidMeasurementUnitConversion
- UpdateValidPreparationInstrument
- UpdateValidPreparationVessel
- UploadMealImage (streaming)
- UploadRecipeImage (streaming)

## WebhooksService

- ArchiveWebhookTriggerEvent
- GetWebhookTriggerEvent
- GetWebhookTriggerEvents
- UpdateWebhookTriggerEvent

## OAuthService

- ArchiveOAuth2Client
- CreateOAuth2Client
- GetOAuth2Client
- GetOAuth2Clients

## SettingsService

- UpdateServiceSettingConfiguration

## DataPrivacyService

- AggregateUserDataReport
- DestroyAllUserData
- FetchUserDataReport

## InternalOperations

- TestQueueMessage

## UploadedMediaService

- Upload (streaming)
