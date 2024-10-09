// GENERATED CODE, DO NOT EDIT MANUALLY

export interface IJWTResponse {
  token: string;
  userID: string;
  householdID: string;
}

export class JWTResponse implements IJWTResponse {
  token: string;
  userID: string;
  householdID: string;
  constructor(input: Partial<JWTResponse> = {}) {
    this.token = input.token = '';
    this.userID = input.userID = '';
    this.householdID = input.householdID = '';
  }
}
