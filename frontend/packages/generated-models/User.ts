// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUser {
  lastAcceptedTOS?: string;
  lastUpdatedAt?: string;
  accountStatus: string;
  archivedAt?: string;
  birthday?: string;
  emailAddressVerifiedAt?: string;
  firstName: string;
  username: string;
  accountStatusExplanation: string;
  createdAt: string;
  emailAddress: string;
  passwordLastChangedAt?: string;
  serviceRoles: string;
  avatar?: string;
  id: string;
  twoFactorSecretVerifiedAt?: string;
  lastAcceptedPrivacyPolicy?: string;
  lastName: string;
  requiresPasswordChange: boolean;
}

export class User implements IUser {
  lastAcceptedTOS?: string;
  lastUpdatedAt?: string;
  accountStatus: string;
  archivedAt?: string;
  birthday?: string;
  emailAddressVerifiedAt?: string;
  firstName: string;
  username: string;
  accountStatusExplanation: string;
  createdAt: string;
  emailAddress: string;
  passwordLastChangedAt?: string;
  serviceRoles: string;
  avatar?: string;
  id: string;
  twoFactorSecretVerifiedAt?: string;
  lastAcceptedPrivacyPolicy?: string;
  lastName: string;
  requiresPasswordChange: boolean;
  constructor(input: Partial<User> = {}) {
    this.lastAcceptedTOS = input.lastAcceptedTOS;
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.accountStatus = input.accountStatus = '';
    this.archivedAt = input.archivedAt;
    this.birthday = input.birthday;
    this.emailAddressVerifiedAt = input.emailAddressVerifiedAt;
    this.firstName = input.firstName = '';
    this.username = input.username = '';
    this.accountStatusExplanation = input.accountStatusExplanation = '';
    this.createdAt = input.createdAt = '';
    this.emailAddress = input.emailAddress = '';
    this.passwordLastChangedAt = input.passwordLastChangedAt;
    this.serviceRoles = input.serviceRoles = '';
    this.avatar = input.avatar;
    this.id = input.id = '';
    this.twoFactorSecretVerifiedAt = input.twoFactorSecretVerifiedAt;
    this.lastAcceptedPrivacyPolicy = input.lastAcceptedPrivacyPolicy;
    this.lastName = input.lastName = '';
    this.requiresPasswordChange = input.requiresPasswordChange = false;
  }
}
