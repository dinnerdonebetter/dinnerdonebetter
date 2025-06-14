// GENERATED CODE, DO NOT EDIT MANUALLY

export interface ITokenResponse {
  accountID: string;
  token: string;
  userID: string;
}

export class TokenResponse implements ITokenResponse {
  accountID: string;
  token: string;
  userID: string;
  constructor(input: Partial<TokenResponse> = {}) {
    this.accountID = input.accountID || '';
    this.token = input.token || '';
    this.userID = input.userID || '';
  }
}
