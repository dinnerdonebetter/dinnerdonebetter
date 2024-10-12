// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUser {
   accountStatus: string;
 accountStatusExplanation: string;
 archivedAt: string;
 avatar: string;
 birthday: string;
 createdAt: string;
 emailAddress: string;
 emailAddressVerifiedAt: string;
 firstName: string;
 id: string;
 lastAcceptedPrivacyPolicy: string;
 lastAcceptedTOS: string;
 lastName: string;
 lastUpdatedAt: string;
 passwordLastChangedAt: string;
 requiresPasswordChange: boolean;
 serviceRoles: string;
 twoFactorSecretVerifiedAt: string;
 username: string;

}

export class User implements IUser {
   accountStatus: string;
 accountStatusExplanation: string;
 archivedAt: string;
 avatar: string;
 birthday: string;
 createdAt: string;
 emailAddress: string;
 emailAddressVerifiedAt: string;
 firstName: string;
 id: string;
 lastAcceptedPrivacyPolicy: string;
 lastAcceptedTOS: string;
 lastName: string;
 lastUpdatedAt: string;
 passwordLastChangedAt: string;
 requiresPasswordChange: boolean;
 serviceRoles: string;
 twoFactorSecretVerifiedAt: string;
 username: string;
constructor(input: Partial<User> = {}) {
	 this.accountStatus = input.accountStatus || '';
 this.accountStatusExplanation = input.accountStatusExplanation || '';
 this.archivedAt = input.archivedAt || '';
 this.avatar = input.avatar || '';
 this.birthday = input.birthday || '';
 this.createdAt = input.createdAt || '';
 this.emailAddress = input.emailAddress || '';
 this.emailAddressVerifiedAt = input.emailAddressVerifiedAt || '';
 this.firstName = input.firstName || '';
 this.id = input.id || '';
 this.lastAcceptedPrivacyPolicy = input.lastAcceptedPrivacyPolicy || '';
 this.lastAcceptedTOS = input.lastAcceptedTOS || '';
 this.lastName = input.lastName || '';
 this.lastUpdatedAt = input.lastUpdatedAt || '';
 this.passwordLastChangedAt = input.passwordLastChangedAt || '';
 this.requiresPasswordChange = input.requiresPasswordChange || false;
 this.serviceRoles = input.serviceRoles || '';
 this.twoFactorSecretVerifiedAt = input.twoFactorSecretVerifiedAt || '';
 this.username = input.username || '';
}
}