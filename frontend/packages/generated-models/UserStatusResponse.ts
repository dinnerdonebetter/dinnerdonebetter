// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IUserStatusResponse {
  accountStatus: string;
  accountStatusExplanation: string;
  activeHousehold: string;
  isAuthenticated: boolean;
  userID: string;
}

export class UserStatusResponse implements IUserStatusResponse {
  accountStatus: string;
  accountStatusExplanation: string;
  activeHousehold: string;
  isAuthenticated: boolean;
  userID: string;
  constructor(input: Partial<UserStatusResponse> = {}) {
    this.accountStatus = input.accountStatus = '';
    this.accountStatusExplanation = input.accountStatusExplanation = '';
    this.activeHousehold = input.activeHousehold = '';
    this.isAuthenticated = input.isAuthenticated = false;
    this.userID = input.userID = '';
  }
}
