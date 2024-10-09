// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUser {
  createdAt: string;
  firstName: string;
  passwordLastChangedAt?: string;
  emailAddressVerifiedAt?: string;
  lastName: string;
  twoFactorSecretVerifiedAt?: string;
  accountStatus: string;
  birthday?: string;
  emailAddress: string;
  lastAcceptedTOS?: string;
  lastUpdatedAt?: string;
  requiresPasswordChange: boolean;
  username: string;
  accountStatusExplanation: string;
  archivedAt?: string;
  avatar?: string;
  id: string;
  lastAcceptedPrivacyPolicy?: string;
  serviceRoles: string;
}

export class User implements IUser {
  createdAt: string;
  firstName: string;
  passwordLastChangedAt?: string;
  emailAddressVerifiedAt?: string;
  lastName: string;
  twoFactorSecretVerifiedAt?: string;
  accountStatus: string;
  birthday?: string;
  emailAddress: string;
  lastAcceptedTOS?: string;
  lastUpdatedAt?: string;
  requiresPasswordChange: boolean;
  username: string;
  accountStatusExplanation: string;
  archivedAt?: string;
  avatar?: string;
  id: string;
  lastAcceptedPrivacyPolicy?: string;
  serviceRoles: string;
  constructor(input: Partial<User> = {}) {
    this.createdAt = input.createdAt = '';
    this.firstName = input.firstName = '';
    this.passwordLastChangedAt = input.passwordLastChangedAt;
    this.emailAddressVerifiedAt = input.emailAddressVerifiedAt;
    this.lastName = input.lastName = '';
    this.twoFactorSecretVerifiedAt = input.twoFactorSecretVerifiedAt;
    this.accountStatus = input.accountStatus = '';
    this.birthday = input.birthday;
    this.emailAddress = input.emailAddress = '';
    this.lastAcceptedTOS = input.lastAcceptedTOS;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.requiresPasswordChange = input.requiresPasswordChange = false;
    this.username = input.username = '';
    this.accountStatusExplanation = input.accountStatusExplanation = '';
    this.archivedAt = input.archivedAt;
    this.avatar = input.avatar;
    this.id = input.id = '';
    this.lastAcceptedPrivacyPolicy = input.lastAcceptedPrivacyPolicy;
    this.serviceRoles = input.serviceRoles = '';
  }
}
