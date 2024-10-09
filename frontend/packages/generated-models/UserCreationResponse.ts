// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserCreationResponse {
  username: string;
  accountStatus: string;
  avatar?: string;
  createdAt: string;
  createdUserID: string;
  emailAddress: string;
  isAdmin: boolean;
  lastName: string;
  birthday?: string;
  firstName: string;
  qrCode: string;
  twoFactorSecret: string;
}

export class UserCreationResponse implements IUserCreationResponse {
  username: string;
  accountStatus: string;
  avatar?: string;
  createdAt: string;
  createdUserID: string;
  emailAddress: string;
  isAdmin: boolean;
  lastName: string;
  birthday?: string;
  firstName: string;
  qrCode: string;
  twoFactorSecret: string;
  constructor(input: Partial<UserCreationResponse> = {}) {
    this.username = input.username = '';
    this.accountStatus = input.accountStatus = '';
    this.avatar = input.avatar;
    this.createdAt = input.createdAt = '';
    this.createdUserID = input.createdUserID = '';
    this.emailAddress = input.emailAddress = '';
    this.isAdmin = input.isAdmin = false;
    this.lastName = input.lastName = '';
    this.birthday = input.birthday;
    this.firstName = input.firstName = '';
    this.qrCode = input.qrCode = '';
    this.twoFactorSecret = input.twoFactorSecret = '';
  }
}
