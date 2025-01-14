// GENERATED CODE, DO NOT EDIT MANUALLY

export interface ITokenResponse {
  householdID: string;
  token: string;
  userID: string;
}

export class TokenResponse implements ITokenResponse {
  householdID: string;
  token: string;
  userID: string;
  constructor(input: Partial<TokenResponse> = {}) {
    this.householdID = input.householdID || '';
    this.token = input.token || '';
    this.userID = input.userID || '';
  }
}
