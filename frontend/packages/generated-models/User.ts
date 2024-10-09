// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUser {
  emailAddressVerifiedAt?: string;
  id: string;
  lastAcceptedPrivacyPolicy?: string;
  requiresPasswordChange: boolean;
  accountStatus: string;
  archivedAt?: string;
  birthday?: string;
  lastName: string;
  lastUpdatedAt?: string;
  avatar?: string;
  createdAt: string;
  emailAddress: string;
  lastAcceptedTOS?: string;
  serviceRoles: string;
  accountStatusExplanation: string;
  firstName: string;
  passwordLastChangedAt?: string;
  twoFactorSecretVerifiedAt?: string;
  username: string;
}

export class User implements IUser {
  emailAddressVerifiedAt?: string;
  id: string;
  lastAcceptedPrivacyPolicy?: string;
  requiresPasswordChange: boolean;
  accountStatus: string;
  archivedAt?: string;
  birthday?: string;
  lastName: string;
  lastUpdatedAt?: string;
  avatar?: string;
  createdAt: string;
  emailAddress: string;
  lastAcceptedTOS?: string;
  serviceRoles: string;
  accountStatusExplanation: string;
  firstName: string;
  passwordLastChangedAt?: string;
  twoFactorSecretVerifiedAt?: string;
  username: string;
  constructor(input: Partial<User> = {}) {
    this.emailAddressVerifiedAt = input.emailAddressVerifiedAt;
    this.id = input.id = '';
    this.lastAcceptedPrivacyPolicy = input.lastAcceptedPrivacyPolicy;
    this.requiresPasswordChange = input.requiresPasswordChange = false;
    this.accountStatus = input.accountStatus = '';
    this.archivedAt = input.archivedAt;
    this.birthday = input.birthday;
    this.lastName = input.lastName = '';
    this.lastUpdatedAt = input.lastUpdatedAt;
    this.avatar = input.avatar;
    this.createdAt = input.createdAt = '';
    this.emailAddress = input.emailAddress = '';
    this.lastAcceptedTOS = input.lastAcceptedTOS;
    this.serviceRoles = input.serviceRoles = '';
    this.accountStatusExplanation = input.accountStatusExplanation = '';
    this.firstName = input.firstName = '';
    this.passwordLastChangedAt = input.passwordLastChangedAt;
    this.twoFactorSecretVerifiedAt = input.twoFactorSecretVerifiedAt;
    this.username = input.username = '';
  }
}
