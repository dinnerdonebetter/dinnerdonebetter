// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUser {
   requiresPasswordChange: boolean;
 serviceRoles: string;
 username: string;
 accountStatusExplanation: string;
 birthday?: string;
 emailAddressVerifiedAt?: string;
 lastAcceptedTOS?: string;
 lastName: string;
 archivedAt?: string;
 emailAddress: string;
 firstName: string;
 id: string;
 lastAcceptedPrivacyPolicy?: string;
 twoFactorSecretVerifiedAt?: string;
 avatar?: string;
 createdAt: string;
 lastUpdatedAt?: string;
 passwordLastChangedAt?: string;
 accountStatus: string;

}

export class User implements IUser {
   requiresPasswordChange: boolean;
 serviceRoles: string;
 username: string;
 accountStatusExplanation: string;
 birthday?: string;
 emailAddressVerifiedAt?: string;
 lastAcceptedTOS?: string;
 lastName: string;
 archivedAt?: string;
 emailAddress: string;
 firstName: string;
 id: string;
 lastAcceptedPrivacyPolicy?: string;
 twoFactorSecretVerifiedAt?: string;
 avatar?: string;
 createdAt: string;
 lastUpdatedAt?: string;
 passwordLastChangedAt?: string;
 accountStatus: string;
constructor(input: Partial<User> = {}) {
	 this.requiresPasswordChange = input.requiresPasswordChange = false;
 this.serviceRoles = input.serviceRoles = '';
 this.username = input.username = '';
 this.accountStatusExplanation = input.accountStatusExplanation = '';
 this.birthday = input.birthday;
 this.emailAddressVerifiedAt = input.emailAddressVerifiedAt;
 this.lastAcceptedTOS = input.lastAcceptedTOS;
 this.lastName = input.lastName = '';
 this.archivedAt = input.archivedAt;
 this.emailAddress = input.emailAddress = '';
 this.firstName = input.firstName = '';
 this.id = input.id = '';
 this.lastAcceptedPrivacyPolicy = input.lastAcceptedPrivacyPolicy;
 this.twoFactorSecretVerifiedAt = input.twoFactorSecretVerifiedAt;
 this.avatar = input.avatar;
 this.createdAt = input.createdAt = '';
 this.lastUpdatedAt = input.lastUpdatedAt;
 this.passwordLastChangedAt = input.passwordLastChangedAt;
 this.accountStatus = input.accountStatus = '';
}
}