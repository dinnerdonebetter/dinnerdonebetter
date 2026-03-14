/**
 * Module-level gRPC client singletons. Reused across requests.
 * Per-request auth is passed via call-level metadata, not per-client credentials.
 */

import { env } from '$env/dynamic/private';
import * as grpc from '@grpc/grpc-js';
import { AuthServiceClient } from '$lib/generated/auth/auth_service';
import type {
	LoginForTokenRequest,
	LoginForTokenResponse,
	GetSelfRequest,
	GetSelfResponse,
	GetActiveAccountRequest,
	GetActiveAccountResponse,
	RequestPasswordResetTokenRequest,
	RequestPasswordResetTokenResponse,
	RedeemPasswordResetTokenRequest,
	RedeemPasswordResetTokenResponse,
	VerifyEmailAddressRequest,
	VerifyEmailAddressResponse,
	ListPasskeysRequest,
	ListPasskeysResponse,
	ArchivePasskeyRequest,
	ArchivePasskeyResponse,
	BeginPasskeyRegistrationRequest,
	BeginPasskeyRegistrationResponse,
	FinishPasskeyRegistrationRequest,
	FinishPasskeyRegistrationResponse,
	BeginPasskeyAuthenticationRequest,
	BeginPasskeyAuthenticationResponse,
	FinishPasskeyAuthenticationRequest
} from '$lib/generated/auth/auth_service_types';
import { IdentityServiceClient } from '$lib/generated/identity/identity_service';
import type {
	GetAccountsForUserRequest,
	GetAccountsForUserResponse,
	GetSentAccountInvitationsRequest,
	GetSentAccountInvitationsResponse,
	CreateAccountInvitationRequest,
	CreateAccountInvitationResponse,
	CancelAccountInvitationRequest,
	CancelAccountInvitationResponse,
	UpdateAccountMemberPermissionsRequest,
	UpdateAccountMemberPermissionsResponse,
	UpdateAccountRequest,
	UpdateAccountResponse,
	UpdateUserUsernameRequest,
	UpdateUserUsernameResponse,
	UpdateUserDetailsRequest,
	UpdateUserDetailsResponse
} from '$lib/generated/identity/identity_service_types';
import { SettingsServiceClient } from '$lib/generated/settings/settings_service';
import type {
	GetServiceSettingsRequest,
	GetServiceSettingsResponse,
	GetServiceSettingConfigurationsForUserRequest,
	GetServiceSettingConfigurationsForUserResponse,
	CreateServiceSettingConfigurationRequest,
	CreateServiceSettingConfigurationResponse,
	UpdateServiceSettingConfigurationRequest,
	UpdateServiceSettingConfigurationResponse
} from '$lib/generated/settings/settings_service_types';
import { MealPlanningServiceClient } from '$lib/generated/mealplanning/mealplanning_service';
import { AnalyticsServiceClient } from '$lib/generated/analytics/analytics_service';
import type {
	SearchForValidPreparationsRequest,
	SearchForValidPreparationsResponse,
	GetValidPreparationInstrumentsByPreparationRequest,
	GetValidPreparationInstrumentsByPreparationResponse,
	GetValidPreparationVesselsByPreparationRequest,
	GetValidPreparationVesselsByPreparationResponse,
	SearchValidIngredientsByPreparationRequest,
	SearchValidIngredientsByPreparationResponse,
	SearchValidMeasurementUnitsByIngredientRequest,
	SearchValidMeasurementUnitsByIngredientResponse,
	GetValidIngredientMeasurementUnitsByIngredientRequest,
	GetValidIngredientMeasurementUnitsByIngredientResponse,
	GetValidIngredientPreparationsByPreparationRequest,
	GetValidIngredientPreparationsByPreparationResponse,
	SearchForValidMeasurementUnitsRequest,
	SearchForValidMeasurementUnitsResponse,
	SearchForValidVesselsRequest,
	SearchForValidVesselsResponse,
	SearchForValidIngredientStatesRequest,
	SearchForValidIngredientStatesResponse,
	CreateRecipeRequest,
	CreateRecipeResponse
} from '$lib/generated/mealplanning/mealplanning_service_types';
import type { Metadata } from '@grpc/grpc-js';

let authClient: AuthServiceClient | null = null;
let identityClient: IdentityServiceClient | null = null;
let settingsClient: SettingsServiceClient | null = null;
let mealplanningClient: MealPlanningServiceClient | null = null;
let analyticsClient: AnalyticsServiceClient | null = null;

function getGrpcServerUrl(): string {
	return env.GRPC_API_SERVER_URL ?? 'localhost:50051';
}

function getCredentials(): grpc.ChannelCredentials {
	if (env.DEVELOPING_LOCALLY === 'true') {
		return grpc.credentials.createInsecure();
	}
	return grpc.credentials.createSsl();
}

function getAuthClient(): AuthServiceClient {
	if (!authClient) {
		authClient = new AuthServiceClient(getGrpcServerUrl(), getCredentials());
	}
	return authClient;
}

function getIdentityClient(): IdentityServiceClient {
	if (!identityClient) {
		identityClient = new IdentityServiceClient(getGrpcServerUrl(), getCredentials());
	}
	return identityClient;
}

function getSettingsClient(): SettingsServiceClient {
	if (!settingsClient) {
		settingsClient = new SettingsServiceClient(getGrpcServerUrl(), getCredentials());
	}
	return settingsClient;
}

function getMealplanningClient(): MealPlanningServiceClient {
	if (!mealplanningClient) {
		mealplanningClient = new MealPlanningServiceClient(getGrpcServerUrl(), getCredentials());
	}
	return mealplanningClient;
}

function getAnalyticsClient(): AnalyticsServiceClient {
	if (!analyticsClient) {
		analyticsClient = new AnalyticsServiceClient(getGrpcServerUrl(), getCredentials());
	}
	return analyticsClient;
}

/**
 * Metadata with Bearer token for authenticated gRPC calls.
 */
export function authMetadata(oauth2AccessToken: string): Metadata {
	const metadata = new grpc.Metadata();
	metadata.add('authorization', `Bearer ${oauth2AccessToken}`);
	return metadata;
}

function promisifyUnary<TRequest, TResponse>(
	call: (
		req: TRequest,
		metadata: Metadata,
		callback: (err: grpc.ServiceError | null, res: TResponse) => void
	) => grpc.ClientUnaryCall
): (req: TRequest, metadata: Metadata) => Promise<TResponse> {
	return (req, metadata) =>
		new Promise((resolve, reject) => {
			call(req, metadata, (err, res) => {
				if (err) reject(err);
				else if (res) resolve(res);
				else reject(new Error('No response'));
			});
		});
}

const emptyMetadata = new grpc.Metadata();

// --- Auth (unauthenticated) ---

export async function loginForToken(
	request: LoginForTokenRequest
): Promise<LoginForTokenResponse> {
	const client = getAuthClient();
	return promisifyUnary<LoginForTokenRequest, LoginForTokenResponse>(
		client.loginForToken.bind(client)
	)(request, emptyMetadata);
}

export async function requestPasswordResetToken(
	request: RequestPasswordResetTokenRequest
): Promise<RequestPasswordResetTokenResponse> {
	const client = getAuthClient();
	return promisifyUnary<RequestPasswordResetTokenRequest, RequestPasswordResetTokenResponse>(
		client.requestPasswordResetToken.bind(client)
	)(request, emptyMetadata);
}

export async function redeemPasswordResetToken(
	request: RedeemPasswordResetTokenRequest
): Promise<RedeemPasswordResetTokenResponse> {
	const client = getAuthClient();
	return promisifyUnary<RedeemPasswordResetTokenRequest, RedeemPasswordResetTokenResponse>(
		client.redeemPasswordResetToken.bind(client)
	)(request, emptyMetadata);
}

export async function verifyEmailAddress(
	request: VerifyEmailAddressRequest
): Promise<VerifyEmailAddressResponse> {
	const client = getAuthClient();
	return promisifyUnary<VerifyEmailAddressRequest, VerifyEmailAddressResponse>(
		client.verifyEmailAddress.bind(client)
	)(request, emptyMetadata);
}

export async function beginPasskeyRegistration(
	oauth2Token: string
): Promise<BeginPasskeyRegistrationResponse> {
	const client = getAuthClient();
	return promisifyUnary<BeginPasskeyRegistrationRequest, BeginPasskeyRegistrationResponse>(
		client.beginPasskeyRegistration.bind(client)
	)({}, authMetadata(oauth2Token));
}

export async function finishPasskeyRegistration(
	oauth2Token: string,
	request: FinishPasskeyRegistrationRequest
): Promise<FinishPasskeyRegistrationResponse> {
	const client = getAuthClient();
	return promisifyUnary<FinishPasskeyRegistrationRequest, FinishPasskeyRegistrationResponse>(
		client.finishPasskeyRegistration.bind(client)
	)(request, authMetadata(oauth2Token));
}

export async function beginPasskeyAuthentication(
	request: BeginPasskeyAuthenticationRequest
): Promise<BeginPasskeyAuthenticationResponse> {
	const client = getAuthClient();
	return promisifyUnary<
		BeginPasskeyAuthenticationRequest,
		BeginPasskeyAuthenticationResponse
	>(client.beginPasskeyAuthentication.bind(client))(request, emptyMetadata);
}

export async function finishPasskeyAuthentication(
	request: FinishPasskeyAuthenticationRequest
): Promise<LoginForTokenResponse> {
	const client = getAuthClient();
	return promisifyUnary<FinishPasskeyAuthenticationRequest, LoginForTokenResponse>(
		client.finishPasskeyAuthentication.bind(client)
	)(request, emptyMetadata);
}

// --- Auth (authenticated) ---

export async function getSelf(oauth2Token: string): Promise<GetSelfResponse> {
	const client = getAuthClient();
	return promisifyUnary<GetSelfRequest, GetSelfResponse>(
		client.getSelf.bind(client)
	)({}, authMetadata(oauth2Token));
}

export async function getActiveAccount(
	oauth2Token: string
): Promise<GetActiveAccountResponse> {
	const client = getAuthClient();
	return promisifyUnary<GetActiveAccountRequest, GetActiveAccountResponse>(
		client.getActiveAccount.bind(client)
	)({}, authMetadata(oauth2Token));
}

export async function listPasskeys(
	oauth2Token: string
): Promise<ListPasskeysResponse> {
	const client = getAuthClient();
	return promisifyUnary<ListPasskeysRequest, ListPasskeysResponse>(
		client.listPasskeys.bind(client)
	)({}, authMetadata(oauth2Token));
}

export async function archivePasskey(
	oauth2Token: string,
	request: ArchivePasskeyRequest
): Promise<ArchivePasskeyResponse> {
	const client = getAuthClient();
	return promisifyUnary<ArchivePasskeyRequest, ArchivePasskeyResponse>(
		client.archivePasskey.bind(client)
	)(request, authMetadata(oauth2Token));
}

// --- Identity (authenticated) ---

export async function getAccountsForUser(
	oauth2Token: string,
	request: GetAccountsForUserRequest
): Promise<GetAccountsForUserResponse> {
	const client = getIdentityClient();
	return promisifyUnary<GetAccountsForUserRequest, GetAccountsForUserResponse>(
		client.getAccountsForUser.bind(client)
	)(request, authMetadata(oauth2Token));
}

export async function getSentAccountInvitations(
	oauth2Token: string,
	request: GetSentAccountInvitationsRequest
): Promise<GetSentAccountInvitationsResponse> {
	const client = getIdentityClient();
	return promisifyUnary<
		GetSentAccountInvitationsRequest,
		GetSentAccountInvitationsResponse
	>(client.getSentAccountInvitations.bind(client))(request, authMetadata(oauth2Token));
}

export async function createAccountInvitation(
	oauth2Token: string,
	request: CreateAccountInvitationRequest
): Promise<CreateAccountInvitationResponse> {
	const client = getIdentityClient();
	return promisifyUnary<
		CreateAccountInvitationRequest,
		CreateAccountInvitationResponse
	>(client.createAccountInvitation.bind(client))(request, authMetadata(oauth2Token));
}

export async function cancelAccountInvitation(
	oauth2Token: string,
	request: CancelAccountInvitationRequest
): Promise<CancelAccountInvitationResponse> {
	const client = getIdentityClient();
	return promisifyUnary<
		CancelAccountInvitationRequest,
		CancelAccountInvitationResponse
	>(client.cancelAccountInvitation.bind(client))(request, authMetadata(oauth2Token));
}

export async function updateAccountMemberPermissions(
	oauth2Token: string,
	request: UpdateAccountMemberPermissionsRequest
): Promise<UpdateAccountMemberPermissionsResponse> {
	const client = getIdentityClient();
	return promisifyUnary<
		UpdateAccountMemberPermissionsRequest,
		UpdateAccountMemberPermissionsResponse
	>(client.updateAccountMemberPermissions.bind(client))(
		request,
		authMetadata(oauth2Token)
	);
}

export async function updateAccount(
	oauth2Token: string,
	request: UpdateAccountRequest
): Promise<UpdateAccountResponse> {
	const client = getIdentityClient();
	return promisifyUnary<UpdateAccountRequest, UpdateAccountResponse>(
		client.updateAccount.bind(client)
	)(request, authMetadata(oauth2Token));
}

export async function updateUserUsername(
	oauth2Token: string,
	request: UpdateUserUsernameRequest
): Promise<UpdateUserUsernameResponse> {
	const client = getIdentityClient();
	return promisifyUnary<UpdateUserUsernameRequest, UpdateUserUsernameResponse>(
		client.updateUserUsername.bind(client)
	)(request, authMetadata(oauth2Token));
}

export async function updateUserDetails(
	oauth2Token: string,
	request: UpdateUserDetailsRequest
): Promise<UpdateUserDetailsResponse> {
	const client = getIdentityClient();
	return promisifyUnary<UpdateUserDetailsRequest, UpdateUserDetailsResponse>(
		client.updateUserDetails.bind(client)
	)(request, authMetadata(oauth2Token));
}

// --- Settings (authenticated) ---

export async function getServiceSettings(
	oauth2Token: string,
	request: GetServiceSettingsRequest
): Promise<GetServiceSettingsResponse> {
	const client = getSettingsClient();
	return promisifyUnary<GetServiceSettingsRequest, GetServiceSettingsResponse>(
		client.getServiceSettings.bind(client)
	)(request, authMetadata(oauth2Token));
}

export async function getServiceSettingConfigurationsForUser(
	oauth2Token: string,
	request: GetServiceSettingConfigurationsForUserRequest
): Promise<GetServiceSettingConfigurationsForUserResponse> {
	const client = getSettingsClient();
	return promisifyUnary<
		GetServiceSettingConfigurationsForUserRequest,
		GetServiceSettingConfigurationsForUserResponse
	>(client.getServiceSettingConfigurationsForUser.bind(client))(
		request,
		authMetadata(oauth2Token)
	);
}

export async function createServiceSettingConfiguration(
	oauth2Token: string,
	request: CreateServiceSettingConfigurationRequest
): Promise<CreateServiceSettingConfigurationResponse> {
	const client = getSettingsClient();
	return promisifyUnary<
		CreateServiceSettingConfigurationRequest,
		CreateServiceSettingConfigurationResponse
	>(client.createServiceSettingConfiguration.bind(client))(
		request,
		authMetadata(oauth2Token)
	);
}

export async function updateServiceSettingConfiguration(
	oauth2Token: string,
	request: UpdateServiceSettingConfigurationRequest
): Promise<UpdateServiceSettingConfigurationResponse> {
	const client = getSettingsClient();
	return promisifyUnary<
		UpdateServiceSettingConfigurationRequest,
		UpdateServiceSettingConfigurationResponse
	>(client.updateServiceSettingConfiguration.bind(client))(
		request,
		authMetadata(oauth2Token)
	);
}

// --- MealPlanning (authenticated) ---

const defaultSearchFilter = { maxResponseSize: 20 };

export async function searchForValidPreparations(
	oauth2Token: string,
	request: Omit<SearchForValidPreparationsRequest, 'filter'> & { filter?: SearchForValidPreparationsRequest['filter'] }
): Promise<SearchForValidPreparationsResponse> {
	const client = getMealplanningClient();
	return promisifyUnary<SearchForValidPreparationsRequest, SearchForValidPreparationsResponse>(
		client.searchForValidPreparations.bind(client)
	)(
		{ ...request, filter: request.filter ?? defaultSearchFilter },
		authMetadata(oauth2Token)
	);
}

export async function getValidPreparationInstrumentsByPreparation(
	oauth2Token: string,
	request: GetValidPreparationInstrumentsByPreparationRequest
): Promise<GetValidPreparationInstrumentsByPreparationResponse> {
	const client = getMealplanningClient();
	return promisifyUnary<
		GetValidPreparationInstrumentsByPreparationRequest,
		GetValidPreparationInstrumentsByPreparationResponse
	>(client.getValidPreparationInstrumentsByPreparation.bind(client))(
		{ ...request, filter: request.filter ?? defaultSearchFilter },
		authMetadata(oauth2Token)
	);
}

export async function getValidPreparationVesselsByPreparation(
	oauth2Token: string,
	request: GetValidPreparationVesselsByPreparationRequest
): Promise<GetValidPreparationVesselsByPreparationResponse> {
	const client = getMealplanningClient();
	return promisifyUnary<
		GetValidPreparationVesselsByPreparationRequest,
		GetValidPreparationVesselsByPreparationResponse
	>(client.getValidPreparationVesselsByPreparation.bind(client))(
		{ ...request, filter: request.filter ?? defaultSearchFilter },
		authMetadata(oauth2Token)
	);
}

export async function searchValidIngredientsByPreparation(
	oauth2Token: string,
	request: SearchValidIngredientsByPreparationRequest
): Promise<SearchValidIngredientsByPreparationResponse> {
	const client = getMealplanningClient();
	return promisifyUnary<
		SearchValidIngredientsByPreparationRequest,
		SearchValidIngredientsByPreparationResponse
	>(client.searchValidIngredientsByPreparation.bind(client))(
		{ ...request, filter: request.filter ?? defaultSearchFilter },
		authMetadata(oauth2Token)
	);
}

export async function searchValidMeasurementUnitsByIngredient(
	oauth2Token: string,
	request: SearchValidMeasurementUnitsByIngredientRequest
): Promise<SearchValidMeasurementUnitsByIngredientResponse> {
	const client = getMealplanningClient();
	return promisifyUnary<
		SearchValidMeasurementUnitsByIngredientRequest,
		SearchValidMeasurementUnitsByIngredientResponse
	>(client.searchValidMeasurementUnitsByIngredient.bind(client))(
		{ ...request, filter: request.filter ?? defaultSearchFilter },
		authMetadata(oauth2Token)
	);
}

export async function getValidIngredientMeasurementUnitsByIngredient(
	oauth2Token: string,
	request: GetValidIngredientMeasurementUnitsByIngredientRequest
): Promise<GetValidIngredientMeasurementUnitsByIngredientResponse> {
	const client = getMealplanningClient();
	return promisifyUnary<
		GetValidIngredientMeasurementUnitsByIngredientRequest,
		GetValidIngredientMeasurementUnitsByIngredientResponse
	>(client.getValidIngredientMeasurementUnitsByIngredient.bind(client))(
		{ ...request, filter: request.filter ?? defaultSearchFilter },
		authMetadata(oauth2Token)
	);
}

export async function getValidIngredientPreparationsByPreparation(
	oauth2Token: string,
	request: GetValidIngredientPreparationsByPreparationRequest
): Promise<GetValidIngredientPreparationsByPreparationResponse> {
	const client = getMealplanningClient();
	return promisifyUnary<
		GetValidIngredientPreparationsByPreparationRequest,
		GetValidIngredientPreparationsByPreparationResponse
	>(client.getValidIngredientPreparationsByPreparation.bind(client))(
		{ ...request, filter: request.filter ?? defaultSearchFilter },
		authMetadata(oauth2Token)
	);
}

export async function searchForValidMeasurementUnits(
	oauth2Token: string,
	request: Omit<SearchForValidMeasurementUnitsRequest, 'filter'> & { filter?: SearchForValidMeasurementUnitsRequest['filter'] }
): Promise<SearchForValidMeasurementUnitsResponse> {
	const client = getMealplanningClient();
	return promisifyUnary<SearchForValidMeasurementUnitsRequest, SearchForValidMeasurementUnitsResponse>(
		client.searchForValidMeasurementUnits.bind(client)
	)(
		{ ...request, filter: request.filter ?? defaultSearchFilter },
		authMetadata(oauth2Token)
	);
}

export async function searchForValidVessels(
	oauth2Token: string,
	request: Omit<SearchForValidVesselsRequest, 'filter'> & { filter?: SearchForValidVesselsRequest['filter'] }
): Promise<SearchForValidVesselsResponse> {
	const client = getMealplanningClient();
	return promisifyUnary<SearchForValidVesselsRequest, SearchForValidVesselsResponse>(
		client.searchForValidVessels.bind(client)
	)(
		{ ...request, filter: request.filter ?? defaultSearchFilter },
		authMetadata(oauth2Token)
	);
}

export async function searchForValidIngredientStates(
	oauth2Token: string,
	request: Omit<SearchForValidIngredientStatesRequest, 'filter'> & { filter?: SearchForValidIngredientStatesRequest['filter'] }
): Promise<SearchForValidIngredientStatesResponse> {
	const client = getMealplanningClient();
	return promisifyUnary<SearchForValidIngredientStatesRequest, SearchForValidIngredientStatesResponse>(
		client.searchForValidIngredientStates.bind(client)
	)(
		{ ...request, filter: request.filter ?? defaultSearchFilter },
		authMetadata(oauth2Token)
	);
}

export async function createRecipe(
	oauth2Token: string,
	request: CreateRecipeRequest
): Promise<CreateRecipeResponse> {
	const client = getMealplanningClient();
	return promisifyUnary<CreateRecipeRequest, CreateRecipeResponse>(
		client.createRecipe.bind(client)
	)(request, authMetadata(oauth2Token));
}

// --- Analytics (unauthenticated) ---

export async function trackEvent(
	event: string,
	properties: Record<string, string> = {}
): Promise<void> {
	const client = getAnalyticsClient();
	await promisifyUnary(
		client.trackEvent.bind(client)
	)(
		{ source: 'web', event, properties },
		emptyMetadata
	);
}

export async function trackAnonymousEvent(
	event: string,
	anonymousId: string,
	properties: Record<string, string> = {}
): Promise<void> {
	const client = getAnalyticsClient();
	await promisifyUnary(
		client.trackAnonymousEvent.bind(client)
	)(
		{ source: 'web', event, anonymousId, properties },
		emptyMetadata
	);
}
