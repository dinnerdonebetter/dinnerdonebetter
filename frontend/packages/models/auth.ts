// Code generated by gen_typescript. DO NOT EDIT.

export interface IChangeActiveHouseholdInput {
  householdID: NonNullable<string>;
}

export class ChangeActiveHouseholdInput implements IChangeActiveHouseholdInput {
  householdID: NonNullable<string> = '';

  constructor(input: Partial<ChangeActiveHouseholdInput> = {}) {
    this.householdID = input.householdID ?? '';
  }
}

export interface IPasswordResetToken {
  createdAt: NonNullable<string>;
  expiresAt: NonNullable<string>;
  archivedAt?: string;
  lastUpdatedAt?: string;
  id: NonNullable<string>;
  token: NonNullable<string>;
  belongsToUser: NonNullable<string>;
}

export class PasswordResetToken implements IPasswordResetToken {
  createdAt: NonNullable<string> = '1970-01-01T00:00:00Z';
  expiresAt: NonNullable<string> = '1970-01-01T00:00:00Z';
  archivedAt?: string;
  lastUpdatedAt?: string;
  id: NonNullable<string> = '';
  token: NonNullable<string> = '';
  belongsToUser: NonNullable<string> = '';

  constructor(input: Partial<PasswordResetToken> = {}) {
    this.createdAt = input.createdAt ?? '1970-01-01T00:00:00Z';
    this.expiresAt = input.expiresAt ?? '1970-01-01T00:00:00Z';
    this.archivedAt = input.archivedAt;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.id = input.id ?? '';
    this.token = input.token ?? '';
    this.belongsToUser = input.belongsToUser ?? '';
  }
}

export interface IPasswordResetTokenCreationRequestInput {
  emailAddress: NonNullable<string>;
}

export class PasswordResetTokenCreationRequestInput implements IPasswordResetTokenCreationRequestInput {
  emailAddress: NonNullable<string> = '';

  constructor(input: Partial<PasswordResetTokenCreationRequestInput> = {}) {
    this.emailAddress = input.emailAddress ?? '';
  }
}

export interface IPasswordResetTokenRedemptionRequestInput {
  token: NonNullable<string>;
  newPassword: NonNullable<string>;
}

export class PasswordResetTokenRedemptionRequestInput implements IPasswordResetTokenRedemptionRequestInput {
  token: NonNullable<string> = '';
  newPassword: NonNullable<string> = '';

  constructor(input: Partial<PasswordResetTokenRedemptionRequestInput> = {}) {
    this.token = input.token ?? '';
    this.newPassword = input.newPassword ?? '';
  }
}

export interface ITOTPSecretRefreshInput {
  currentPassword: NonNullable<string>;
  totpToken: NonNullable<string>;
}

export class TOTPSecretRefreshInput implements ITOTPSecretRefreshInput {
  currentPassword: NonNullable<string> = '';
  totpToken: NonNullable<string> = '';

  constructor(input: Partial<TOTPSecretRefreshInput> = {}) {
    this.currentPassword = input.currentPassword ?? '';
    this.totpToken = input.totpToken ?? '';
  }
}

export interface ITOTPSecretVerificationInput {
  totpToken: NonNullable<string>;
  userID: NonNullable<string>;
}

export class TOTPSecretVerificationInput implements ITOTPSecretVerificationInput {
  totpToken: NonNullable<string> = '';
  userID: NonNullable<string> = '';

  constructor(input: Partial<TOTPSecretVerificationInput> = {}) {
    this.totpToken = input.totpToken ?? '';
    this.userID = input.userID ?? '';
  }
}

export interface ITOTPSecretRefreshResponse {
  qrCode: NonNullable<string>;
  twoFactorSecret: NonNullable<string>;
}

export class TOTPSecretRefreshResponse implements ITOTPSecretRefreshResponse {
  qrCode: NonNullable<string> = '';
  twoFactorSecret: NonNullable<string> = '';

  constructor(input: Partial<TOTPSecretRefreshResponse> = {}) {
    this.qrCode = input.qrCode ?? '';
    this.twoFactorSecret = input.twoFactorSecret ?? '';
  }
}

export interface IPasswordUpdateInput {
  newPassword: NonNullable<string>;
  currentPassword: NonNullable<string>;
  totpToken: NonNullable<string>;
}

export class PasswordUpdateInput implements IPasswordUpdateInput {
  newPassword: NonNullable<string> = '';
  currentPassword: NonNullable<string> = '';
  totpToken: NonNullable<string> = '';

  constructor(input: Partial<PasswordUpdateInput> = {}) {
    this.newPassword = input.newPassword ?? '';
    this.currentPassword = input.currentPassword ?? '';
    this.totpToken = input.totpToken ?? '';
  }
}

export interface IJWTResponse {
  userID: NonNullable<string>;
  householdID: NonNullable<string>;
  token: NonNullable<string>;
}

export class JWTResponse implements IJWTResponse {
  userID: NonNullable<string> = '';
  householdID: NonNullable<string> = '';
  token: NonNullable<string> = '';

  constructor(input: Partial<JWTResponse> = {}) {
    this.userID = input.userID ?? '';
    this.householdID = input.householdID ?? '';
    this.token = input.token ?? '';
  }
}
