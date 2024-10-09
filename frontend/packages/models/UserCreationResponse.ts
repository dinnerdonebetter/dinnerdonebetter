// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserCreationResponse {
  emailAddress: string;
  twoFactorSecret: string;
  username: string;
  accountStatus: string;
  createdAt: string;
  createdUserID: string;
  firstName: string;
  isAdmin: boolean;
  lastName: string;
  qrCode: string;
  avatar: string;
  birthday: string;
}

export class UserCreationResponse implements IUserCreationResponse {
  emailAddress: string;
  twoFactorSecret: string;
  username: string;
  accountStatus: string;
  createdAt: string;
  createdUserID: string;
  firstName: string;
  isAdmin: boolean;
  lastName: string;
  qrCode: string;
  avatar: string;
  birthday: string;
  constructor(input: Partial<UserCreationResponse> = {}) {
    this.emailAddress = input.emailAddress || '';
    this.twoFactorSecret = input.twoFactorSecret || '';
    this.username = input.username || '';
    this.accountStatus = input.accountStatus || '';
    this.createdAt = input.createdAt || '';
    this.createdUserID = input.createdUserID || '';
    this.firstName = input.firstName || '';
    this.isAdmin = input.isAdmin || false;
    this.lastName = input.lastName || '';
    this.qrCode = input.qrCode || '';
    this.avatar = input.avatar || '';
    this.birthday = input.birthday || '';
  }
}
