// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserStatusResponse {
  userID: string;
  accountStatus: string;
  accountStatusExplanation: string;
  activeHousehold: string;
  isAuthenticated: boolean;
}

export class UserStatusResponse implements IUserStatusResponse {
  userID: string;
  accountStatus: string;
  accountStatusExplanation: string;
  activeHousehold: string;
  isAuthenticated: boolean;
  constructor(input: Partial<UserStatusResponse> = {}) {
    this.userID = input.userID = '';
    this.accountStatus = input.accountStatus = '';
    this.accountStatusExplanation = input.accountStatusExplanation = '';
    this.activeHousehold = input.activeHousehold = '';
    this.isAuthenticated = input.isAuthenticated = false;
  }
}
