// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserCreationResponse {
   accountStatus: string;
 avatar: string;
 birthday: string;
 createdAt: string;
 createdUserID: string;
 emailAddress: string;
 firstName: string;
 isAdmin: boolean;
 lastName: string;
 qrCode: string;
 twoFactorSecret: string;
 username: string;

}

export class UserCreationResponse implements IUserCreationResponse {
   accountStatus: string;
 avatar: string;
 birthday: string;
 createdAt: string;
 createdUserID: string;
 emailAddress: string;
 firstName: string;
 isAdmin: boolean;
 lastName: string;
 qrCode: string;
 twoFactorSecret: string;
 username: string;
constructor(input: Partial<UserCreationResponse> = {}) {
	 this.accountStatus = input.accountStatus || '';
 this.avatar = input.avatar || '';
 this.birthday = input.birthday || '';
 this.createdAt = input.createdAt || '';
 this.createdUserID = input.createdUserID || '';
 this.emailAddress = input.emailAddress || '';
 this.firstName = input.firstName || '';
 this.isAdmin = input.isAdmin || false;
 this.lastName = input.lastName || '';
 this.qrCode = input.qrCode || '';
 this.twoFactorSecret = input.twoFactorSecret || '';
 this.username = input.username || '';
}
}