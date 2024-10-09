// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserCreationResponse {
  avatar?: string;
  createdUserID: string;
  emailAddress: string;
  twoFactorSecret: string;
  username: string;
  accountStatus: string;
  birthday?: string;
  createdAt: string;
  firstName: string;
  isAdmin: boolean;
  lastName: string;
  qrCode: string;
}

export class UserCreationResponse implements IUserCreationResponse {
  avatar?: string;
  createdUserID: string;
  emailAddress: string;
  twoFactorSecret: string;
  username: string;
  accountStatus: string;
  birthday?: string;
  createdAt: string;
  firstName: string;
  isAdmin: boolean;
  lastName: string;
  qrCode: string;
  constructor(input: Partial<UserCreationResponse> = {}) {
    this.avatar = input.avatar;
    this.createdUserID = input.createdUserID = '';
    this.emailAddress = input.emailAddress = '';
    this.twoFactorSecret = input.twoFactorSecret = '';
    this.username = input.username = '';
    this.accountStatus = input.accountStatus = '';
    this.birthday = input.birthday;
    this.createdAt = input.createdAt = '';
    this.firstName = input.firstName = '';
    this.isAdmin = input.isAdmin = false;
    this.lastName = input.lastName = '';
    this.qrCode = input.qrCode = '';
  }
}
