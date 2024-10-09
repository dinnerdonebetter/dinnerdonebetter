// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserStatusResponse {
  activeHousehold: string;
  isAuthenticated: boolean;
  userID: string;
  accountStatus: string;
  accountStatusExplanation: string;
}

export class UserStatusResponse implements IUserStatusResponse {
  activeHousehold: string;
  isAuthenticated: boolean;
  userID: string;
  accountStatus: string;
  accountStatusExplanation: string;
  constructor(input: Partial<UserStatusResponse> = {}) {
    this.activeHousehold = input.activeHousehold = '';
    this.isAuthenticated = input.isAuthenticated = false;
    this.userID = input.userID = '';
    this.accountStatus = input.accountStatus = '';
    this.accountStatusExplanation = input.accountStatusExplanation = '';
  }
}
