// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IJWTResponse {
  householdID: string;
  token: string;
  userID: string;
}

export class JWTResponse implements IJWTResponse {
  householdID: string;
  token: string;
  userID: string;
  constructor(input: Partial<JWTResponse> = {}) {
    this.householdID = input.householdID || '';
    this.token = input.token || '';
    this.userID = input.userID || '';
  }
}
