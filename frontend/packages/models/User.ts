// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUser {
  lastAcceptedTOS: string;
  lastUpdatedAt: string;
  requiresPasswordChange: boolean;
  twoFactorSecretVerifiedAt: string;
  username: string;
  accountStatus: string;
  birthday: string;
  emailAddressVerifiedAt: string;
  avatar: string;
  createdAt: string;
  emailAddress: string;
  firstName: string;
  lastName: string;
  lastAcceptedPrivacyPolicy: string;
  passwordLastChangedAt: string;
  serviceRoles: string;
  accountStatusExplanation: string;
  archivedAt: string;
  id: string;
}

export class User implements IUser {
  lastAcceptedTOS: string;
  lastUpdatedAt: string;
  requiresPasswordChange: boolean;
  twoFactorSecretVerifiedAt: string;
  username: string;
  accountStatus: string;
  birthday: string;
  emailAddressVerifiedAt: string;
  avatar: string;
  createdAt: string;
  emailAddress: string;
  firstName: string;
  lastName: string;
  lastAcceptedPrivacyPolicy: string;
  passwordLastChangedAt: string;
  serviceRoles: string;
  accountStatusExplanation: string;
  archivedAt: string;
  id: string;
  constructor(input: Partial<User> = {}) {
    this.lastAcceptedTOS = input.lastAcceptedTOS || '';
    this.lastUpdatedAt = input.lastUpdatedAt || '';
    this.requiresPasswordChange = input.requiresPasswordChange || false;
    this.twoFactorSecretVerifiedAt = input.twoFactorSecretVerifiedAt || '';
    this.username = input.username || '';
    this.accountStatus = input.accountStatus || '';
    this.birthday = input.birthday || '';
    this.emailAddressVerifiedAt = input.emailAddressVerifiedAt || '';
    this.avatar = input.avatar || '';
    this.createdAt = input.createdAt || '';
    this.emailAddress = input.emailAddress || '';
    this.firstName = input.firstName || '';
    this.lastName = input.lastName || '';
    this.lastAcceptedPrivacyPolicy = input.lastAcceptedPrivacyPolicy || '';
    this.passwordLastChangedAt = input.passwordLastChangedAt || '';
    this.serviceRoles = input.serviceRoles || '';
    this.accountStatusExplanation = input.accountStatusExplanation || '';
    this.archivedAt = input.archivedAt || '';
    this.id = input.id || '';
  }
}
