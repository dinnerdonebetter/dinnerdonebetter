// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserCreationResponse {
  createdUserID: string;
  firstName: string;
  lastName: string;
  qrCode: string;
  twoFactorSecret: string;
  username: string;
  accountStatus: string;
  avatar?: string;
  birthday?: string;
  createdAt: string;
  emailAddress: string;
  isAdmin: boolean;
}

export class UserCreationResponse implements IUserCreationResponse {
  createdUserID: string;
  firstName: string;
  lastName: string;
  qrCode: string;
  twoFactorSecret: string;
  username: string;
  accountStatus: string;
  avatar?: string;
  birthday?: string;
  createdAt: string;
  emailAddress: string;
  isAdmin: boolean;
  constructor(input: Partial<UserCreationResponse> = {}) {
    this.createdUserID = input.createdUserID = '';
    this.firstName = input.firstName = '';
    this.lastName = input.lastName = '';
    this.qrCode = input.qrCode = '';
    this.twoFactorSecret = input.twoFactorSecret = '';
    this.username = input.username = '';
    this.accountStatus = input.accountStatus = '';
    this.avatar = input.avatar;
    this.birthday = input.birthday;
    this.createdAt = input.createdAt = '';
    this.emailAddress = input.emailAddress = '';
    this.isAdmin = input.isAdmin = false;
  }
}
