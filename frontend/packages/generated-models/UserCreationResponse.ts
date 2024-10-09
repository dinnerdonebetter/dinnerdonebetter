// GENERATED CODE, DO NOT EDIT MANUALLY



export interface IUserCreationResponse {
   isAdmin: boolean;
 qrCode: string;
 username: string;
 accountStatus: string;
 avatar?: string;
 createdUserID: string;
 emailAddress: string;
 twoFactorSecret: string;
 birthday?: string;
 createdAt: string;
 firstName: string;
 lastName: string;

}

export class UserCreationResponse implements IUserCreationResponse {
   isAdmin: boolean;
 qrCode: string;
 username: string;
 accountStatus: string;
 avatar?: string;
 createdUserID: string;
 emailAddress: string;
 twoFactorSecret: string;
 birthday?: string;
 createdAt: string;
 firstName: string;
 lastName: string;
constructor(input: Partial<UserCreationResponse> = {}) {
	 this.isAdmin = input.isAdmin = false;
 this.qrCode = input.qrCode = '';
 this.username = input.username = '';
 this.accountStatus = input.accountStatus = '';
 this.avatar = input.avatar;
 this.createdUserID = input.createdUserID = '';
 this.emailAddress = input.emailAddress = '';
 this.twoFactorSecret = input.twoFactorSecret = '';
 this.birthday = input.birthday;
 this.createdAt = input.createdAt = '';
 this.firstName = input.firstName = '';
 this.lastName = input.lastName = '';
}
}